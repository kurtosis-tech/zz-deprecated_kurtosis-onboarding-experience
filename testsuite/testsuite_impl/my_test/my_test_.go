package my_test

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/networks"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/services"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

const (
	lambdaID                                = "eth-lambda"
	ethereumKurtosisLambdaImage             = "kurtosistech/ethereum-kurtosis-lambda"
	emptyJsonParams                         = "{}"
	rpcPort                                 = 8545
	execCommandSuccessExitCode              = 0
)

type MyTest struct{}

type EthereumKurtosisLambdaResult struct {
	BootnodeServiceID     services.ServiceID   `json:"bootnode_service_id"`
	NodeServiceIDs        []services.ServiceID `json:"node_service_ids"`
	RpcPort               uint32               `json:"rpc_port"`
	SignerKeystoreContent string               `json:"signer_keystore_content"`
	SignerAccountPassword string               `json:"signer_account_password"`
}

func (test MyTest) Configure(builder *testsuite.TestConfigurationBuilder) {
	builder.WithSetupTimeoutSeconds(240).WithRunTimeoutSeconds(240)
}

func (test MyTest) Setup(networkCtx *networks.NetworkContext) (networks.Network, error) {

	_, err := networkCtx.LoadLambda(lambdaID, ethereumKurtosisLambdaImage, emptyJsonParams)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred loading the Ethereum Kurtosis Lambda in the test")
	}

	logrus.Info("The Ethereum Kurtosis Lambda has been successfully added in the test")

	return networkCtx, nil
}

func (test MyTest) Run(uncastedNetwork networks.Network) error {
	// Necessary because Go doesn't have generics
	castedNetwork := uncastedNetwork.(*networks.NetworkContext)

	ethLambdaCtx, err := castedNetwork.GetLambdaContext(lambdaID)

	respJsonStr, err := ethLambdaCtx.Execute(emptyJsonParams)
	if err != nil {
		return stacktrace.Propagate(err, "And error occurred executing the Ethereum Kurtosis Lambda")
	}
	ethResult := new(EthereumKurtosisLambdaResult)
	if err := json.Unmarshal([]byte(respJsonStr), ethResult); err != nil {
		return stacktrace.Propagate(err, "An error occurred deserializing the Lambda response")
	}

	bootnodeServiceCtx, err := castedNetwork.GetServiceContext(ethResult.BootnodeServiceID)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the Ethereum bootnode service context")
	}
	logrus.Infof("Got service context for Ethereum bootnode service '%v'", bootnodeServiceCtx.GetServiceID())

	//Execute a Geth command, inside the Ethereum bootnode, to get the accounts list
	gethCmd := "geth attach data/geth.ipc --exec eth.accounts"
	listAccountsCmd := []string{
		"/bin/sh",
		"-c",
		gethCmd,
	}

	exitCode, logOutput, err := bootnodeServiceCtx.ExecCommand(listAccountsCmd)
	if err != nil {
		return stacktrace.Propagate(err, "Executing command '%v' returned an error", gethCmd)
	}
	if exitCode != execCommandSuccessExitCode {
		return stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
	}

	logrus.Infof("Account list: %v", string(*logOutput))

	//Instantiate Geth Client
	url := fmt.Sprintf("http://%v:%v", bootnodeServiceCtx.GetIPAddress(), rpcPort)
	gethClient, err := ethclient.Dial(url)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred getting the Golang Ethereum client")
	}
	defer gethClient.Close()

	return nil
}
