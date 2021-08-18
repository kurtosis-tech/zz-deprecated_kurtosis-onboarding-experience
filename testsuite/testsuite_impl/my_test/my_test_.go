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

	//TODO Replace with the first part of the implementation

	//TODO Replace with the last part of the implementation

	return networkCtx, nil
}

func (test MyTest) Run(uncastedNetwork networks.Network) error {

	//TODO Replace with the first part of the implementation

	//TODO Replace with the last part of the implementation

	return nil
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================
