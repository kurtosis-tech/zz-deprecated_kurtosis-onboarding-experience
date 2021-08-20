Ethereum On-Boarding Testsuite
==============================

## Implement a Single Node Ethereum Test Network

1. Create an account on https://kurtosistech.com if you don't have one yet.
1. Verify that the Docker daemon is running on your local machine.
1. Clone this repository by running `git clone https://github.com/kurtosis-tech/kurtosis-onboarding-experience.git`
1. Run the empty Ethereum single node test `my_test` to verify that everything works on your local machine.
   1. Run `bash scripts/build-and-run.sh all`
   1. Verify that the output of the build-and-run.sh script indicates that one test ran (my_test) and that it passed.
1. Import the dependencies that are used in this example test suite.
   1. Run `go get github.com/ethereum/go-ethereum`
   1. Run `go get github.com/palantir/stacktrace`
   1. Run `go get github.com/sirupsen/logrus`
1. Write the Setup() method on the Ethereum single node test in order to bootstrap the network and leave it running and ready to use it.
   1. In your preferred IDE, open the Ethereum single node test `my_test` at `testsuite/testsuite_impl/my_test/my_test.go`
   1. Implement the private helper function `getEthereumServiceConfigurations` used to set the Ethereum node service inside the testsuite
      1. Add the following private helper functions `getContainerCreationConfig()`, `getRunConfigFunc()` and `getEthereumServiceConfigurations()` in [this Gist](https://gist.github.com/leoporoli/123cb1d682d74dcafe7920f01809b167) to the bottom of the file, so they can be used later. 
   1. Implement the first part of the Setup() method which starts the Ethereum single node service in the testsuite
      1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/d81577dcfd5fdc6605ccdf9f61eed81f) in the top of the Setup() method
   1. Implement the second and last part of the Setup() method to check if the service is available
      1. Add the following private helper functions `getEnodeAddress()` and `sendRpcCall` in [this Gist](https://gist.github.com/leoporoli/1a03e9500a61a20d06ed8e3827d72f5e) to the bottom of the file, so they can be used later.
      1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/f9aacad32b2800a98f68bcf5fa32165c) in the bottom of the Setup() method
   1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
1. Write the Run() method on the Ethereum single node test in order to test the initilization of an Ethereum Single Node Network in Dev Mode
   1. Read [the official Geth documentation for start Ethereum in Dev mode](https://geth.ethereum.org/docs/getting-started/dev-mode) to have a context as this test implements the commands that are described in this document
   1. Implement the first part of the Run() method which get the Ethereum single node service from the network, get the Ethereum client from the service and uses this to get the chain ID of the private Ethereum Network
      1. Add the private helper function `getEthClient` in [this Gist](https://gist.github.com/leoporoli/8a1641f9d78752f984ed672895c1f97c) to the bottom of the file, so it can be used later.
      1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/628c78452e7e05549809f5cbcb62cdf6) in the top of the Run() method
      1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
   1. Implement the second and last part of the Run() method which implements all the commands of the official Geth documentation in the `Dev mode` section
      1. 1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/5a1539d2e8e45d3658a5b2398d9f3ba7) in the bottom of the Run() method
   1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that the test contained business logic for an Ethereum single node network test, and that it passed.

## Implement an Advanced Test which test and Ethereum Private Network with Multiple Nodes

The purpose of this test could be separated into two parts

The first part of it implements a private Ethereum network with multiple nodes, that uses `Clique consensus` as proof of authority and that is previously set in the genesis block, with a signer account.

To start the private network, it will be necessary to start the bootnode first and then, the remaining nodes using the ENR address from the bootnode. 

To ensure interconnectivity of the private network it should contain a validation to validate the numbers of peers that all the nodes have. 

This part of the test follows the official Geth guide: [how to set up an Ethereum Private Network](https://geth.ethereum.org/docs/interface/private-network)

The second part of the test involves testing a transaction into the private network previously set. It will be executed deploying a simple smart contract, that was previously bonded and placed in `smart_contracts/bindings/hello_world.go`, called `Hello World`

The second part of the test consists of testing executions within the previously set private network. For this, first, the smart contract `Hello World` that is in the folder` smart_contracts / bindings / hello_world.go` must be implemented and validated that it has been successfully implemented.

And, finally, validate if the Bootnode contains that the signer account

1. Implements the Setup() method of the test `my_advanced_test_.go` in order to start the Ethereum private network with multiple nodes
   1. In your preferred IDE, open the advanced Ethereum test `my_advanced_test` at `testsuite/testsuite_impl/my_advanced_test/my_advanced_test.go`
   1. Start the private network composed by one bootnode and three simple nodes and using the custom genesis block stored at `testsuite/data/genesis.json`
      1. Start the bootnode and get its ENR address which will be important to start to others nodes
         1. Add the service to the testsuite's network
            1. Create and object of `services.ContainerCreationConfig` you can review the basic test `my_test` to check how was created on it
               1. Load the following statics files:`genesis.json`, `password.txt` and `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` using the IDs previously set in the `Configure` method, with the `WithStaticFiles()` method of the builder
               1. Set ports `8545` with `tcp` protocol and port `30303` with `tcp` and `udp` protocol, using the `WithUsedPorts()` method of the builder
            1. Implements the anonymous function that returns the `services.ContainerRunConfig` object
               1. Use the `services.NewContainerRunConfigBuilder()` function in order to create the `services.ContainerRunConfig` object to return
                  1. Get the filepath of `genesis.json` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                  1. Get the filepath of `password.txt` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                  1. Get the filepath of `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                  1. Create the entrypoint that contains the command to execute the Ethereum node into the container
                     1. Write the command to init the genesis block
                        1. Set `datadir` value
                        1. Set the custom genesis block file using the filepath generated on the previous step
                     1. Write the command to start the bootnode. It should be written after the init command and using the `&&` operator to execute them sequentially
                        1. Set the `datadir` option
                        1. Set the `keystore` option
                        1. Set the `network ID` option remember that it is defined in the `genesis.json` file
                        1. Enable `HTTP-RPC` server
                        1. Set the `IP address` of the `HTTP-RPC` server
                        1. Set the `API's offered over the HTTP-RPC interface` using these values `admin,eth,net,web3,miner,personal,txpool,debug`
                        1. Accept cross origin requests from any domain using this value `*`
                        1. Set the `IP address` of the node
                        1. Set the `port` of the node
                        1. Unlock the `signer account` to allow it to mine, remember that you can get this account address from the keystore file `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026`
                        1. Enable `mine`
                        1. Allow `insecure account unlocking` to allow signer account to mine
                        1. Set and Network IP to restrict network communication using a CIDR mask
                        1. Set the `filepath of the password file` that allows you to avoid entering the password manually when you want to execute a command
            1. Calls the function `networkCtx.AddService` and passes it and service identifier, and the others two arguments defined in the previous steps
         1. Checks if the service is up and running
            1. Add the following private helper functions `waitForStartup`, `isAvailable`, `getEnodeAddress` and `sendRpcCall` at the end of the file, so they are available for later use
               ```
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
               
               func sendRpcCall(ipAddress string, rpcJsonString string, targetStruct interface{}) error {
                  rpcPort := 8545
                  rpcRequestTimeout := 30 * time.Second
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
               ```
            1. Implements the call to the `waitForStartup` method in order to control if the service is running
         1. Get the bootnode's ENR address            
            1. Execute a geth command inside the service to get the ENR address
            ```
            exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
               "/bin/sh",
               "-c",
               fmt.Sprintf("geth attach data/geth.ipc --exec admin.nodeInfo.enr"),
            })
            if err != nil {
               return "", stacktrace.Propagate(err, "Executing command returned an error with logs: %+v", string(*logOutput))
            }
            if exitCode != execCommandSuccessExitCode {
               return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
            }
            ```
            1. Cast the logout value, receive for the command, to an string which contains the bootnode's ENR address
            
      1. Start the remaining nodes with the help of the bootnode
         1. Add the first node
            1. Implements the `services.ContainerCreationConfig` object
               1. Load the following statics files:`genesis.json` and `password.txt` using the IDs previously set in the `Configure` method, with the `WithStaticFiles()` method of the builder
               1. Set ports `8545` with `tcp` protocol and port `30303` with `tcp` and `udp` protocol, using the `WithUsedPorts()` method of the builder

            1. Implements the anonymous function that returns the `services.ContainerRunConfig` object
               1. Use the `services.NewContainerRunConfigBuilder()` function in order to create the `services.ContainerRunConfig` object to return
                  1. Get the filepath of `genesis.json` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                  1. Get the filepath of `password.txt` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                  1. Create the entrypoint that contains the command to execute the Ethereum node into the container
                     1. Write the command to init the genesis block
                        1. Set `datadir` value
                        1. Set the custom genesis block file using the filepath generated on the previous step
                     1. Write the command to start the bootnode. It should be written after the init command and using the `&&` operator to execute them sequentially
                        1. Set the `datadir` option
                        1. Set the `network ID` option remember that it is defined in the `genesis.json` file
                        1. Enable `HTTP-RPC` server
                        1. Set the `IP address` of the `HTTP-RPC` server
                        1. Set the `API's offered over the HTTP-RPC interface` using these values `admin,eth,net,web3,miner,personal,txpool,debug`
                        1. Accept cross origin requests from any domain using this value `*`
                        1. Set the `IP address` of the node
                        1. Set the `port` of the node
                        1. Set the `bootnode`using the ENR address previously get
            1. Calls the function `networkCtx.AddService` and passes it a service identifier, and the others two arguments defined in the previous steps
         1. Checks if the node is up and running 
            1. Implement the call to the `waitForStartup` method again to check if this node is running successfully
            
         1. Get the `Enode` address that will be used to connect with the remaining nodes
         1. Connect the node with peers (this must be done from the second loaded node)
            1. Connect the node manually using the command `admin_addPeer` [explained on this document](https://geth.ethereum.org/docs/rpc/ns-admin#admin_addpeer)
            1. Check the link between nodes. You can list the peers using the command `admin_peers` [explained on this document](https://geth.ethereum.org/docs/rpc/ns-admin#admin_peers)
         1. Repet the previous steps in order to start and link the remaining nodes
   1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that two tests ran (my_test and my_advanced_test) and that them passed
1. Implements the Run() method of the test
   1. Execute a transaction to deploy the `hello_world` smart contract into the private network
      1. Create an instance of the Geth client which will be necessary to deploy the smart contract
      1. Get the private key's signer account   
         1. Get the content of signer account's keystore file through the `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` file previously loaded into the testsuite
         1. Get the password through the `password.txt` file previously loaded into the testsuite
         1. Get the signer account's private key through the function `keystore.DecryptKey` and passes it the values get from the keystore file and the password file
      1. Get a transactor which will be necessary to execute the deployment, we suggest using the function `NewKeyedTransactorWithChainID` to accomplish this step
      1. Execute a transaction to deploy the smart contract using the function `DeployHelloWorld` provided by the binding `hello_world.go` which will receive the transactor and the Geth client as arguments  
      1. Check if the transaction has been successfully completed
         1. Add the private helper function `waitUntilTransactionMined` at the end of the file
         ```
         func waitUntilTransactionMined(validatorClient *ethclient.Client, transactionHash common.Hash) error {
            maxNumCheckTransactionMinedRetries := 10
            timeBetweenCheckTransactionMinedRetries := 1 * time.Second
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
         ```
         1. Use the private helper function `waitUntilTransactionMined` in order to check if the transaction has been successfully completed
      1. Check the operation of the `HelloWorld` smart contract using the function `helloWorld.Greet()`    
   1. Validate if the bootnode has the signer account, you can use the `eth.accounts` command to list the accounts
1. Add `my_advance_test` in the list, that the testsuite will be execute, inside the `my_testsuite.go` file
1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that two test ran (my_test and my_advanced_test) and both are successfully executed





