package my_advanced_test

import (
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/networks"
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/services"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
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

	//TODO Start the bootnode with the genesis block and get its ENR address

	//TODO Checks if bootnode service is available

	//TODO Start the other three nodes using the bootnode's ENR address

	//TODO Checks if all of theme are available

	//TODO AddPeer between them and check if the amount of peer is ok

	return networkCtx, nil
}

func (test MyAdvancedTest) Run(uncastedNetwork networks.Network) error {

	//TODO Get the bootnode service from the network

	//TODO Instantiate the Geth Client

	//TODO Get the private key from the signer account stored into the bootnode

	//TODO Create a new transactor using the private key and the Chain ID

	//TODO Init a new transaction to deploy the 'HelloWorld' smart contract

	//TODO Validate that the transaction is fully mined

	//TODO Validate if the bootnode contains the signer account

	return nil
}

// ====================================================================================================
//                                       Private helper functions
// ====================================================================================================

