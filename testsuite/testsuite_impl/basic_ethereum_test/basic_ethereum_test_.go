package basic_ethereum_test

import (
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/networks"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/services"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
)

const (
	ethereumNodeImage = "ethereum/client-go"
	ethereumNodeRpcPort = 8545

	node0ServiceID services.ServiceID = "node-0"

	genesisStaticFileID    services.StaticFileID = "genesis.json"
	genesisStaticFilePath                        = "/data/genesis.json"
	passwordStaticFileID   services.StaticFileID = "password.txt"
	passwordStaticFilePath                       = "/data/password.txt"
	signerKeystoreFileID   services.StaticFileID = "UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026"
	signerKeystoreFilePath                       = "/data/UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026"
)

type BasicEthereumTest struct{}

type NodeInfoResponse struct {
	Result NodeInfo `json:"result"`
}

type NodeInfo struct {
	Enode string `json:"enode"`
}

type AddPeerResponse struct {
	Result bool `json:"result"`
}

func (test BasicEthereumTest) Configure(builder *testsuite.TestConfigurationBuilder) {
	builder.WithSetupTimeoutSeconds(360).WithRunTimeoutSeconds(360)
}

func (test BasicEthereumTest) Setup(networkCtx *networks.NetworkContext) (networks.Network, error) {

	//TODO Start a single ETH node in dev mode

	//TODO Check if the Ethereum network is available

	return networkCtx, nil
}

func (test BasicEthereumTest) Run(uncastedNetwork networks.Network) error {
	//TODO Get Go Ethereum client

	//TODO Get ETH network's chain ID

	//TODO Create new ETH account

	//TODO Make ETH transfer transaction

	//TODO Get ETH account balance

	return nil
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================
//TODO Container creation & run config helper functions

//TODO Create Go Ethereum client helper function
