package my_advanced_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/networks"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/services"
	"github.com/kurtosis-tech/kurtosis-onboarding-experience/smart_contracts/bindings"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ethereumGoClientImageName = "ethereum/client-go"
	rpcPort                   = 8545
	discoveryPort             = 30303

	rpcRequestTimeout = 30 * time.Second

	genesisStaticFileID    services.StaticFileID = "genesis.json"
	genesisStaticFilePath                        = "/data/genesis.json"
	passwordStaticFileID   services.StaticFileID = "password.txt"
	passwordStaticFilePath                       = "/data/password.txt"
	signerKeystoreFileID   services.StaticFileID = "UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026"
	signerKeystoreFilePath                       = "/data/UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026"
	bootnodeID             services.ServiceID    = "bootnode"
	gethServiceIdPrefix                          = "ethereum-node-"

	subnetRange                = "/24"
	execCommandSuccessExitCode = 0

	numGethNodes = 3

	// Availability check constants
	gethTimeBetweenIsAvailablePolls = 1 * time.Second
	gethMaxIsAvailablePolls         = 30

	enodePrefix       = "enode://"
	handshakeProtocol = "eth: \"handshake\""
	adminInfoRpcCall  = `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`

	maxNumCheckTransactionMinedRetries      = 10
	timeBetweenCheckTransactionMinedRetries = 1 * time.Second

	signerAccount = "0x14f6136b48b74b147926c9f24323d16c1e54a026"
)

type MyAdvancedTest struct{}

type NodeInfoResponse struct {
	Result NodeInfo `json:"result"`
}

type NodeInfo struct {
	Enode string `json:"enode"`
}

type AddPeerResponse struct {
	Result bool `json:"result"`
}

func (test MyAdvancedTest) Configure(builder *testsuite.TestConfigurationBuilder) {
	builder.WithSetupTimeoutSeconds(
		240,
	).WithRunTimeoutSeconds(
		240,
	).WithStaticFileFilepaths(map[services.StaticFileID]string{
		genesisStaticFileID:  genesisStaticFilePath,
		signerKeystoreFileID: signerKeystoreFilePath,
		passwordStaticFileID: passwordStaticFilePath,
	})
}

func (test MyAdvancedTest) Setup(networkCtx *networks.NetworkContext) (networks.Network, error) {

	bootNodeEnr, err := startEthBootnode(networkCtx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred starting the Ethereum Bootnode")
	}
	var nodesEnode []string
	for i := 1; i <= numGethNodes; i++ {
		serviceID := services.ServiceID(gethServiceIdPrefix + strconv.Itoa(i))
		enode, err := starEthNodeByBootnode(networkCtx, serviceID, bootNodeEnr, nodesEnode)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred starting the Ethereum Node '%v'", serviceID)
		}
		nodesEnode = append(nodesEnode, enode)
	}

	return networkCtx, nil
}

func (test MyAdvancedTest) Run(uncastedNetwork networks.Network) error {
	// Necessary because Go doesn't have generics
	castedNetwork := uncastedNetwork.(*networks.NetworkContext)

	serviceCtx, err := castedNetwork.GetServiceContext(bootnodeID)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the Ethereum Go Client service context")
	}
	logrus.Infof("Got service context for Ethereum Go Client service '%v'", serviceCtx.GetServiceID())

	gethClient, err := getClient(serviceCtx.GetIPAddress())
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get a gethClient from bootnode.")
	}
	defer gethClient.Close()

	key, err := getPrivateKey(serviceCtx)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get private key")
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(15))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	transactor.GasPrice = big.NewInt(5)
	address, tx, helloWorld, err := bindings.DeployHelloWorld(transactor, gethClient)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred deploying the HelloWorld contract on the ETH Network")
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())

	if err := waitUntilTransactionMined(gethClient, tx.Hash()); err != nil {
		return stacktrace.Propagate(err, "An error occurred waiting for the HelloWorld contract to be mined")
	}
	logrus.Info("Deployed Hello World contract")

	name, err := helloWorld.Greet(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Fatalf("Failed to retrieve pending name: %v", err)
	}
	fmt.Println("Pending name:", name)

	listAccountsCmd := []string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("geth attach data/geth.ipc --exec eth.accounts"),
	}

	exitCode, logOutput, err := serviceCtx.ExecCommand(listAccountsCmd)
	if err != nil {
		return stacktrace.Propagate(err, "Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != execCommandSuccessExitCode {
		return stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}

	if err = validateAccount(string(*logOutput)); err != nil {
		return stacktrace.Propagate(err, "Validate account error")
	}

	return nil
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================
func startEthBootnode(networkCtx *networks.NetworkContext) (string, error) {
	containerCreationConfig, runConfigFunc := getEthereumServiceConfigurationsForBootnode()

	_, hostPortBindings, err := networkCtx.AddService(bootnodeID, containerCreationConfig, runConfigFunc)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred adding the Ethereum Bootnode")
	}

	logrus.Infof("Added Ethereum Bootnode service with host port bindings: %+v ", hostPortBindings)

	serviceCtx, err := networkCtx.GetServiceContext(bootnodeID)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred getting the Ethereum Bootnode service context")
	}
	logrus.Infof("Got service context for Ethereum Bootnode service '%v' and IP address '%v", serviceCtx.GetServiceID(), serviceCtx.GetIPAddress())

	if err := waitForStartup(serviceCtx.GetIPAddress(), gethTimeBetweenIsAvailablePolls, gethMaxIsAvailablePolls); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred waiting for service with ID '%v' to start", bootnodeID)
	}

	exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("geth attach data/geth.ipc --exec admin.nodeInfo.enr"),
	})
	if err != nil {
		return "", stacktrace.Propagate(err,"Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != execCommandSuccessExitCode {
		return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}

	enr := string(*logOutput)

	return enr, nil
}

func starEthNodeByBootnode(networkCtx *networks.NetworkContext, serviceID services.ServiceID, bootnodeEnr string, nodesEnode []string) (string, error) {
	containerCreationConfig, runConfigFunc := getEthereumServiceConfigurationsForNode(bootnodeEnr)
	_, hostPortBindings, err := networkCtx.AddService(serviceID, containerCreationConfig, runConfigFunc)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred adding the Ethereum Node %v", serviceID)
	}

	logrus.Infof("Added Ethereum Node %v service with host port bindings: %+v ", serviceID, hostPortBindings)

	serviceCtx, err := networkCtx.GetServiceContext(serviceID)
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred getting the Ethereum Node %v service context", serviceID)
	}
	logrus.Infof("Got service context for Ethereum Node '%v' and IP address '%v'", serviceCtx.GetServiceID(), serviceCtx.GetIPAddress())

	if err := waitForStartup(serviceCtx.GetIPAddress(), gethTimeBetweenIsAvailablePolls, gethMaxIsAvailablePolls); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred waiting for service with ID '%v' to start", serviceID)
	}

	for _, enode := range nodesEnode {
		ok, err := AddPeer(serviceCtx.GetIPAddress(), enode)
		if err != nil {
			return "", stacktrace.Propagate(err, "Failed to call addPeer endpoint to add peer with enode %v", enode)
		}
		if !ok {
			return "", stacktrace.NewError("addPeer endpoint returned false on service %v, adding peer %v", serviceID, enode)
		}
	}

	exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("geth attach data/geth.ipc --exec admin.peers"),
	})
	if err != nil {
		return "", stacktrace.Propagate(err, "Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != execCommandSuccessExitCode {
		return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}

	if err = validatePeersQuantity(string(*logOutput), serviceID, nodesEnode); err != nil {
		return "", stacktrace.Propagate(err, "Validate peers error")
	}

	enode, err := getEnodeAddress(serviceCtx.GetIPAddress())
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to get enode from peer %v", serviceID)
	}

	return enode, nil
}

func getEthereumServiceConfigurationsForBootnode() (*services.ContainerCreationConfig, func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error)) {
	containerCreationConfig := getContainerCreationConfig(true)

	runConfigFunc := getRunConfigFunc()
	return containerCreationConfig, runConfigFunc
}

func getEthereumServiceConfigurationsForNode(bootnodeEnr string) (*services.ContainerCreationConfig, func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error)) {
	containerCreationConfig := getContainerCreationConfig(false)

	runConfigFunc := getRunConfigFuncForSimpleNodes(bootnodeEnr)
	return containerCreationConfig, runConfigFunc
}

func getContainerCreationConfig(adSignerKeystore bool) *services.ContainerCreationConfig {
	staticFiles := map[services.StaticFileID]bool{
		genesisStaticFileID:  true,
		passwordStaticFileID: true,
	}

	if adSignerKeystore {
		staticFiles[signerKeystoreFileID] = true
	}

	containerCreationConfig := services.NewContainerCreationConfigBuilder(
		ethereumGoClientImageName,
	).WithUsedPorts(
		map[string]bool{
			fmt.Sprintf("%v/tcp", rpcPort):       true,
			fmt.Sprintf("%v/tcp", discoveryPort): true,
			fmt.Sprintf("%v/udp", discoveryPort): true,
		},
	).WithStaticFiles(staticFiles).Build()
	return containerCreationConfig
}

func getRunConfigFunc() func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
	runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
		genesisFilepath, found := staticFileFilepaths[genesisStaticFileID]
		if !found {
			return nil, stacktrace.NewError("No filepath found for key '%v'; this is a bug in Kurtosis!", genesisStaticFileID)
		}

		passwordFilepath, found := staticFileFilepaths[passwordStaticFileID]
		if !found {
			return nil, stacktrace.NewError("No filepath found for key '%v'; this is a bug in Kurtosis!", passwordStaticFileID)
		}

		signerKeystoreFilepath, found := staticFileFilepaths[signerKeystoreFileID]
		if !found {
			return nil, stacktrace.NewError("No filepath found for key '%v'; this is a bug in Kurtosis!", signerKeystoreFileID)
		}

		keystoreFolder := filepath.Dir(signerKeystoreFilepath)

		ipNet := getIPNet(ipAddr)

		entryPointArgs := []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("geth init --datadir data %v && geth --keystore %v --datadir data --networkid 15 -http --http.api admin,eth,net,web3,miner,personal,txpool,debug --http.addr %v --http.corsdomain '*' --nat extip:%v --port %v --unlock 0x14f6136b48b74b147926c9f24323d16c1e54a026 --mine --allow-insecure-unlock --netrestrict %v --password %v", genesisFilepath, keystoreFolder, ipAddr, ipAddr, discoveryPort, ipNet, passwordFilepath),
		}

		result := services.NewContainerRunConfigBuilder().WithEntrypointOverride(entryPointArgs).Build()
		return result, nil
	}
	return runConfigFunc
}

func getIPNet(ipAddr string) *net.IPNet {
	cidr := ipAddr + subnetRange
	_, ipNet, _ := net.ParseCIDR(cidr)
	logrus.Infof("IPNET value: %v", ipNet)
	return ipNet
}

func getRunConfigFuncForSimpleNodes(bootnodeEnr string) func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
	runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
		genesisFilepath, found := staticFileFilepaths[genesisStaticFileID]
		if !found {
			return nil, stacktrace.NewError("No filepath found for test file 1 key '%v'; this is a bug in Kurtosis!", genesisStaticFileID)
		}

		entryPointArgs := []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("geth init --datadir data %v && geth --datadir data --networkid 15 -http --http.api admin,eth,net,web3,miner,personal,txpool,debug --http.addr %v --http.corsdomain '*' --nat extip:%v --gcmode archive --syncmode full --port %v --bootnodes %v", genesisFilepath, ipAddr, ipAddr, discoveryPort, bootnodeEnr),
		}

		result := services.NewContainerRunConfigBuilder().WithEntrypointOverride(entryPointArgs).Build()
		return result, nil
	}
	return runConfigFunc
}

func getClient(ipAddress string) (*ethclient.Client, error) {
	url := fmt.Sprintf("http://%v:%v", ipAddress, rpcPort)
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the Golang Ethereum client")
	}
	return client, nil
}

func AddPeer(ipAddress string, peerEnode string) (bool, error) {
	adminAddPeerRpcCall := fmt.Sprintf(`{"jsonrpc":"2.0", "method": "admin_addPeer", "params": ["%v"], "id":70}`, peerEnode)
	logrus.Infof("Admin add peer rpc call: %v", adminAddPeerRpcCall)
	addPeerResponse := new(AddPeerResponse)
	err := sendRpcCall(ipAddress, adminAddPeerRpcCall, addPeerResponse)
	logrus.Infof("AddPeer response: %+v", addPeerResponse)
	if err != nil {
		return false, stacktrace.Propagate(err, "Failed to send addPeer RPC call for enode %v", peerEnode)
	}
	return addPeerResponse.Result, nil
}

func getPrivateKey(serviceCtx *services.ServiceContext) (*keystore.Key, error) {
	staticFileAbsFilepaths, err := serviceCtx.LoadStaticFiles(map[services.StaticFileID]bool{
		signerKeystoreFileID: true,
		passwordStaticFileID: true,
	})
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred loading the static files corresponding to keys '%v' and '%v'", signerKeystoreFileID, passwordStaticFileID)
	}
	signerKeystoreFilepath, found := staticFileAbsFilepaths[signerKeystoreFileID]
	if !found {
		return nil, stacktrace.Propagate(err, "No filepath found for key '%v'; this is a bug in Kurtosis!", signerKeystoreFilepath)
	}

	signerKeystoreContent, err := ioutil.ReadFile(signerKeystoreFilepath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error happens reading file '%v'", signerKeystoreFilepath)
	}

	json, err := ioutil.ReadAll(strings.NewReader(string(signerKeystoreContent)))
	if err != nil {
		return nil, stacktrace.Propagate(err,"An error occurred when trying to read content for filepath '%v'", signerKeystoreFilepath)
	}

	passwordFilepath, found := staticFileAbsFilepaths[passwordStaticFileID]
	if !found {
		return nil, stacktrace.Propagate(err, "No filepath found for key '%v'; this is a bug in Kurtosis!", passwordFilepath)
	}

	passwordContent, err := ioutil.ReadFile(passwordFilepath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error happens reading file '%v'", passwordFilepath)
	}

	key, err := keystore.DecryptKey(json, string(passwordContent))
	if err != nil {
		return nil, stacktrace.Propagate(err,"An error occurred when trying to decrypt the private key")
	}
	return key, nil
}

func sendRpcCall(ipAddress string, rpcJsonString string, targetStruct interface{}) error {
	url := fmt.Sprintf("http://%v:%v", ipAddress, rpcPort)
	var jsonByteArray = []byte(rpcJsonString)

	logrus.Debugf("Sending RPC call to '%v' with JSON body '%v'...", url, rpcJsonString)

	client := http.Client{
		Timeout: rpcRequestTimeout,
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonByteArray))
	if err != nil {
		return stacktrace.Propagate(err, "Failed to send RPC request to geth node with ip '%v'", ipAddress)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// For debugging
		var teeBuf bytes.Buffer
		tee := io.TeeReader(resp.Body, &teeBuf)
		bodyBytes, err := ioutil.ReadAll(tee)
		if err != nil {
			return stacktrace.Propagate(err, "Error parsing geth node response into bytes.")
		}
		bodyString := string(bodyBytes)
		logrus.Tracef("Response for RPC call %v: %v", rpcJsonString, bodyString)

		err = json.NewDecoder(&teeBuf).Decode(targetStruct)
		if err != nil {
			return stacktrace.Propagate(err, "Error parsing geth node response into target struct.")
		}
		return nil
	} else {
		return stacktrace.NewError("Received non-200 status code rom admin RPC api: %v", resp.StatusCode)
	}
}

func waitForStartup(ipAddress string, timeBetweenPolls time.Duration, maxNumRetries int) error {
	for i := 0; i < maxNumRetries; i++ {
		if isAvailable(ipAddress) {
			return nil
		}

		// Don't wait if we're on the last iteration of the loop, since we'd be waiting unnecessarily
		if i < maxNumRetries-1 {
			time.Sleep(timeBetweenPolls)
		}
	}
	return stacktrace.NewError(
		"Service with ip '%v' did not become available despite polling %v times with %v between polls",
		ipAddress,
		maxNumRetries,
		timeBetweenPolls)
}

func isAvailable(ipAddress string) bool {
	enodeAddress, err := getEnodeAddress(ipAddress)
	if err != nil {
		return false
	} else {
		return strings.HasPrefix(enodeAddress, enodePrefix)
	}
}

func getEnodeAddress(ipAddress string) (string, error) {
	nodeInfoResponse := new(NodeInfoResponse)
	err := sendRpcCall(ipAddress, adminInfoRpcCall, nodeInfoResponse)
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to send admin node info RPC request to geth node with ip %v", ipAddress)
	}
	return nodeInfoResponse.Result.Enode, nil
}

func waitUntilTransactionMined(validatorClient *ethclient.Client, transactionHash common.Hash) error {
	for i := 0; i < maxNumCheckTransactionMinedRetries; i++ {
		receipt, err := validatorClient.TransactionReceipt(context.Background(), transactionHash)
		if err == nil && receipt != nil && receipt.BlockNumber != nil {
			return nil
		}
		if i < maxNumCheckTransactionMinedRetries-1 {
			time.Sleep(timeBetweenCheckTransactionMinedRetries)
		}
	}
	return stacktrace.NewError(
		"Transaction with hash '%v' wasn't mined even after checking %v times with %v between checks",
		transactionHash.Hex(),
		maxNumCheckTransactionMinedRetries,
		timeBetweenCheckTransactionMinedRetries)
}

func validatePeersQuantity(logString string, serviceID services.ServiceID, nodesEnode []string) error {
	peersQuantity := strings.Count(logString, enodePrefix) - strings.Count(logString, handshakeProtocol)
	validPeersQuantity := len(nodesEnode) + 1
	if peersQuantity != validPeersQuantity {
		return stacktrace.NewError("The amount of peers '%v' for node '%v' is not valid, should be '%v?", peersQuantity, serviceID, validPeersQuantity)
	}
	return nil
}

func validateAccount(logString string) error {
	count := strings.Count(logString, signerAccount)
	if count < 1 {
		return stacktrace.NewError("The eth private network doesn't contains the signer account '%v'?", signerAccount)
	}
	return nil
}
