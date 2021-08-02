package my_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

var serviceIDs = []services.ServiceID{
	"service-0",
}

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
	// ================== BEGIN ETH CODE ================
	containerCreationConfig, runConfigFunc := getMyServiceConfigurations()

	serviceContext, hostPortBindings, err := networkCtx.AddService(serviceIDs[0], containerCreationConfig, runConfigFunc)
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
	return nil
}

func getMyServiceConfigurations() (*services.ContainerCreationConfig, func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error)) {
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
