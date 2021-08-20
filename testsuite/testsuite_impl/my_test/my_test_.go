package my_test

import (
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/networks"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
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

	//TODO Replace with instructions for starting an Ethereum single node in dev mode

	//TODO Replace with instructions for check service availability

	return networkCtx, nil
}

func (test MyTest) Run(uncastedNetwork networks.Network) error {

	//TODO Replace with instructions for get the ETH network's chain ID

	//TODO Replace with instructions for create a new ETH account

	//TODO Replace with instructions for execute an ETH transaction

	//TODO Replace with instructions for get the account's balance

	return nil
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================
//TODO Add the private helper functions
