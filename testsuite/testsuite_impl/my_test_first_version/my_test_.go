package my_test_first_version

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/networks"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/services"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ethereumGoClientImageName = "ethereum/client-go"
	ethereumGoClientPort      = 8565
	rpcPort                   = 8545
	adminInfoRpcCall          = `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`

	waitForStartupTimeBetweenPolls = 1 * time.Second
	waitForStartupMaxPolls         = 90

	serviceID = "my-eth-client"
)

type MyTest struct{}

type NodeInfoResponse struct {
	Result NodeInfo `json:"result"`
}

type NodeInfo struct {
	Enode string `json:"enode"`
}

func (test MyTest) Configure(builder *testsuite.TestConfigurationBuilder) {
	builder.WithSetupTimeoutSeconds(360).WithRunTimeoutSeconds(360)
}

func (test MyTest) Setup(networkCtx *networks.NetworkContext) (networks.Network, error) {
	containerCreationConfig, runConfigFunc := getEthereumServiceConfigurations()

	serviceContext, hostPortBindings, err := networkCtx.AddService(serviceID, containerCreationConfig, runConfigFunc)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred adding the Ethereum Go Client service")
	}

	firstNodeUp := false
	for pollCount := 0; pollCount < waitForStartupMaxPolls; pollCount++ {
		enodeAddress, err := getEnodeAddress(serviceContext.GetIPAddress())
		if err == nil {
			firstNodeUp = true
			logrus.Infof("Enode address: %v", enodeAddress)
			break
		}
		time.Sleep(waitForStartupTimeBetweenPolls)
	}
	if !firstNodeUp {
		return nil, stacktrace.Propagate(err, "First geth node failed to come up")
	}

	logrus.Infof("Added Ethereum Go Client service with host port bindings: %+v", hostPortBindings)
	return networkCtx, nil
}

func (test MyTest) Run(uncastedNetwork networks.Network) error {
	// Necessary because Go doesn't have generics
	castedNetwork := uncastedNetwork.(*networks.NetworkContext)

	serviceCtx, err := castedNetwork.GetServiceContext(serviceID)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the Ethereum Go Client service context")
	}
	logrus.Infof("Got service context for Ethereum Go Client service '%v'", serviceCtx.GetServiceID())

	gethClient, err := getClient(serviceCtx.GetIPAddress())
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get a gethClient from first node.")
	}
	defer gethClient.Close()

	networkId, err := gethClient.NetworkID(context.Background())
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get network ID")
	}
	logrus.Infof("Chain ID: %v", networkId)

	exitCode, logOutput, err := serviceCtx.ExecCommand([]string{"/bin/sh", "-c",
		fmt.Sprintf("printf \"passphrase\\npassphrase\\n\" | geth attach /tmp/geth.ipc --exec 'personal.newAccount()'"),
	})
	if err != nil {
		return stacktrace.NewError("Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != 0 {
		return stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}
	logrus.Infof("Logs: %+v", string(*logOutput))

	exitCode, logOutput, err = serviceCtx.ExecCommand([]string{"/bin/sh", "-c",
		fmt.Sprintf("geth attach /tmp/geth.ipc --exec 'eth.sendTransaction({from:eth.coinbase, to:eth.accounts[1], value: web3.toWei(0.05, \"ether\")})'"),
	})
	if err != nil {
		return stacktrace.NewError("Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != 0 {
		return stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}
	logrus.Infof("Logs: %+v", string(*logOutput))

	exitCode, logOutput, err = serviceCtx.ExecCommand([]string{"/bin/sh", "-c",
		fmt.Sprintf("geth attach /tmp/geth.ipc --exec 'eth.getBalance(eth.accounts[1])'"),
	})
	if err != nil {
		return stacktrace.NewError("Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != 0 {
		return stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}
	logrus.Infof("Logs: %+v", string(*logOutput))
	return nil
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================
func getEthereumServiceConfigurations() (*services.ContainerCreationConfig, func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error)) {
	containerCreationConfig := getContainerCreationConfig()

	runConfigFunc := getRunConfigFunc()
	return containerCreationConfig, runConfigFunc
}

func getContainerCreationConfig() *services.ContainerCreationConfig {
	containerCreationConfig := services.NewContainerCreationConfigBuilder(
		ethereumGoClientImageName,
	).WithUsedPorts(
		map[string]bool{fmt.Sprintf("%v/tcp", ethereumGoClientPort): true},
	).Build()
	return containerCreationConfig
}

func getRunConfigFunc() func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
	runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
		entrypointCommand := fmt.Sprintf("geth --dev -http --http.api admin,eth,net,rpc --http.addr %v --http.corsdomain '*' --nat extip:%v",
			ipAddr,
			ipAddr)
		entrypointArgs := []string{
			"/bin/sh",
			"-c",
			entrypointCommand,
		}
		result := services.NewContainerRunConfigBuilder().WithEntrypointOverride(entrypointArgs).Build()
		return result, nil
	}
	return runConfigFunc
}

func getEnodeAddress(ipAddress string) (string, error) {
	nodeInfoResponse := new(NodeInfoResponse)
	err := sendRpcCall(ipAddress, adminInfoRpcCall, nodeInfoResponse)
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to send admin node info RPC request to geth node.")
	}
	return nodeInfoResponse.Result.Enode, nil
}

func sendRpcCall(ipAddress string, rpcJsonString string, targetStruct interface{}) error {
	url := fmt.Sprintf("http://%v:%v", ipAddress, rpcPort)
	var jsonByteArray = []byte(rpcJsonString)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonByteArray))
	if err != nil {
		return stacktrace.Propagate(err, "Failed to send RPC request to geth node")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return stacktrace.NewError("Received non-200 status code from admin RPC API: %v", resp.StatusCode)
	}

	// For debugging
	var teeBuf bytes.Buffer
	tee := io.TeeReader(resp.Body, &teeBuf)
	bodyBytes, err := ioutil.ReadAll(tee)
	if err != nil {
		return stacktrace.Propagate(err, "Error parsing geth node response into bytes.")
	}
	bodyString := string(bodyBytes)
	logrus.Debugf("Response for RPC call %v: %v", rpcJsonString, bodyString)

	if err = json.NewDecoder(&teeBuf).Decode(targetStruct); err != nil {
		return stacktrace.Propagate(err, "Error parsing geth node response into target struct.")
	}

	return nil

}

func getClient(ipAddress string) (*ethclient.Client, error) {
	url := fmt.Sprintf("http://%v:%v", ipAddress, rpcPort)
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the Golang Ethereum client")
	}
	return client, nil
}
