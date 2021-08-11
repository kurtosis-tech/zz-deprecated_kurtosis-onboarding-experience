package my_advanced_test

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/networks"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/services"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	ethereumGoClientImageName = "ethereum/client-go"
	ethereumGoClientPort      = 8565
	rpcPort                   = 8545
	adminInfoRpcCall          = `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`

	genesisStaticFileID   services.StaticFileID = "genesis.json"
	genesisStaticFilePath                       = "/data/genesis.json"
	bootnodeID            services.ServiceID    = "bootnode"
	gethServiceIdPrefix                         = "ethereum-node-"

	execCommandSuccessExitCode = 0

	bootnodePort  = 30304
	firstNodePort = bootnodePort + 1

	// Num Geth nodes (including bootstrapper)
	numGethNodes = 3
)

type MyAdvancedTest struct{}

type NodeInfoResponse struct {
	Result NodeInfo `json:"result"`
}

type NodeInfo struct {
	Enode string `json:"enode"`
}

func (test MyAdvancedTest) Configure(builder *testsuite.TestConfigurationBuilder) {
	builder.WithSetupTimeoutSeconds(
		1800,
	).WithRunTimeoutSeconds(
		1800,
	).WithStaticFileFilepaths(map[services.StaticFileID]string{
		genesisStaticFileID: genesisStaticFilePath,
	})
}

func (test MyAdvancedTest) Setup(networkCtx *networks.NetworkContext) (networks.Network, error) {

	bootNodeEnr, err := startEthBootnode(networkCtx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred starting the Ethereum Bootnode service")
	}

	for i := 0; i < numGethNodes; i++ {
		serviceID := services.ServiceID(gethServiceIdPrefix + strconv.Itoa(i))
		port := firstNodePort + i
		err := starEthNodeByBootnode(networkCtx, serviceID, port, bootNodeEnr)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred starting the Ethereum Node %v service", serviceID)
		}
	}

	return networkCtx, nil
}

func (test MyAdvancedTest) Run(uncastedNetwork networks.Network) error {

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
	logrus.Infof("Got service context for Ethereum Bootnode service '%v'", serviceCtx.GetServiceID())

	time.Sleep(2 * time.Second)

	exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("printf \"passphrase\\npassphrase\\n\" | geth attach data/geth.ipc --exec admin.nodeInfo.enr"),
	})
	if err != nil {
		return "", stacktrace.NewError("Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != execCommandSuccessExitCode {
		return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}

	enr := string(*logOutput)

	logrus.Infof("Bootnode enr: %+v", enr)

	return enr, nil
}

func starEthNodeByBootnode(networkCtx *networks.NetworkContext, serviceID services.ServiceID, port int, bootnodeEnr string) error {
	containerCreationConfig, runConfigFunc := getEthereumServiceConfigurationsForNode(port, bootnodeEnr)
	_, hostPortBindings, err := networkCtx.AddService(serviceID, containerCreationConfig, runConfigFunc)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred adding the Ethereum Node %v", serviceID)
	}

	logrus.Infof("Added Ethereum Node %v service with host port bindings: %+v ", serviceID, hostPortBindings)

	serviceCtx, err := networkCtx.GetServiceContext(serviceID)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the Ethereum Node %v service context", serviceID)
	}
	logrus.Infof("Got service context for Ethereum Node '%v'", serviceCtx.GetServiceID())

	time.Sleep(1 * time.Second)

	exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("printf \"passphrase\\npassphrase\\n\" | geth attach data/geth.ipc --exec admin.peers"),
	})
	if err != nil {
		return stacktrace.NewError("Executing command returned an error with logs: %+v", string(*logOutput))
	}
	if exitCode != execCommandSuccessExitCode {
		return stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}

	logrus.Infof("Current peers: %+v for Eth Node %v", string(*logOutput), serviceID)

	return nil
}

func getEthereumServiceConfigurationsForBootnode() (*services.ContainerCreationConfig, func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error)) {
	containerCreationConfig := getContainerCreationConfig()

	runConfigFunc := getRunConfigFunc()
	return containerCreationConfig, runConfigFunc
}
func getEthereumServiceConfigurationsForNode(port int, bootnodeEnr string) (*services.ContainerCreationConfig, func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error)) {
	containerCreationConfig := getContainerCreationConfig()

	runConfigFunc := getRunConfigFuncForSimpleNodes(port, bootnodeEnr)
	return containerCreationConfig, runConfigFunc
}

func getContainerCreationConfig() *services.ContainerCreationConfig {
	containerCreationConfig := services.NewContainerCreationConfigBuilder(
		ethereumGoClientImageName,
	).WithUsedPorts(
		map[string]bool{
			fmt.Sprintf("%v/tcp", ethereumGoClientPort): true,
		},
	).WithStaticFiles(map[services.StaticFileID]bool{
		genesisStaticFileID: true,
	}).Build()
	return containerCreationConfig
}

func getRunConfigFunc() func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
	runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
		genesisFilepath, found := staticFileFilepaths[genesisStaticFileID]
		if !found {
			return nil, stacktrace.NewError("No filepath found for test file 1 key '%v'; this is a bug in Kurtosis!", genesisStaticFileID)
		}

		entryPointArgs := []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("geth init --datadir data %v && geth --datadir data --networkid 15 -http --http.api admin,eth,net,rpc --http.addr %v --http.corsdomain '*' --nat extip:%v --port %v", genesisFilepath, ipAddr, ipAddr, bootnodePort),
		}

		result := services.NewContainerRunConfigBuilder().WithEntrypointOverride(entryPointArgs).Build()
		return result, nil
	}
	return runConfigFunc
}

func getRunConfigFuncForSimpleNodes(port int, bootnodeEnr string) func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
	runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
		genesisFilepath, found := staticFileFilepaths[genesisStaticFileID]
		if !found {
			return nil, stacktrace.NewError("No filepath found for test file 1 key '%v'; this is a bug in Kurtosis!", genesisStaticFileID)
		}

		entryPointArgs := []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("geth init --datadir data %v && geth --datadir data --networkid 15 -http --http.api admin,eth,net,rpc --http.addr %v --http.corsdomain '*' --nat extip:%v --port %v --bootnodes %v", genesisFilepath, ipAddr, ipAddr, port, bootnodeEnr),
		}

		result := services.NewContainerRunConfigBuilder().WithEntrypointOverride(entryPointArgs).Build()
		return result, nil
	}
	return runConfigFunc
}
