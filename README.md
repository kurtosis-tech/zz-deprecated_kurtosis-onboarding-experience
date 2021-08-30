Ethereum On-Boarding Testsuite
==============================

## Implement a basic test which tests transactions on a single node Ethereum testnet

1. Create an account on [https://www.kurtosistech.com/sign-up](https://www.kurtosistech.com/sign-up) if you don't have one yet.
1. Verify that the Docker daemon is running on your local machine using `docker container ls`
1. Clone this repository by running `git clone https://github.com/kurtosis-tech/kurtosis-onboarding-experience.git --branch master`
1. Run the empty Ethereum single node test `my_test` to verify that everything works on your local machine.
    1. Run `bash scripts/build-and-run.sh all`
    1. Verify that the output of the build-and-run.sh script indicates that one test ran (my_test) and that it passed.
1. Set up a single node Ethereum testnet in Kurtosis
    1. In your preferred IDE, open the Ethereum single node test `my_test` at `testsuite/testsuite_impl/my_test/my_test.go`
    1. Set the container configuration for the Ethereum container in your testnet.
        1. Add the following container configuration helper function to the bottom of the test file.
       ```
        func getContainerCreationConfig() *services.ContainerCreationConfig {
         containerCreationConfig := services.NewContainerCreationConfigBuilder(
            "ethereum/client-go",
         ).WithUsedPorts(
            map[string]bool{fmt.Sprintf("%v/tcp", 8545): true},
         ).Build()
         return containerCreationConfig
        }
       ```
    1. Set the runtime configuration for the Ethereum container in your testnet.
        1. Add the following runtime configuration helper function to the bottom of the test file.
       ```
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
          ```      
    1. Add an Ethereum node with your configurations to the `Setup()` method using the Kurtosis "network context" object.
        1. Add the following code to the body of the `Setup()` method.
            ```
            containerCreationConfig := getContainerCreationConfig()
            runConfigFunc := getRunConfigFunc()
          
            serviceContext, hostPortBindings, err := networkCtx.AddService("my-eth-client", containerCreationConfig, runConfigFunc)
            if err != nil {
                return nil, stacktrace.Propagate(err, "An error occurred adding the Ethereum Go Client service")
            }
            ```
    1. Add an availability check to the `Setup()` method to make sure your Ethereum node is fully functional before your test starts.
        1. Add a helper function `sendRpcCall` to send RPC calls to the Ethereum node to the bottom of the test file.
           ```
           func sendRpcCall(ipAddress string, rpcJsonString string, targetStruct interface{}) error {
             rpcPort := 8545
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
           ```
        1. Add a helper function to retrieve the enode address of an Ethereum node to the bottom of the test file.
           ```
           func getEnodeAddress(ipAddress string) (string, error) {
               nodeInfoResponse := new(NodeInfoResponse)
               adminInfoRpcCall := `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`
               err := sendRpcCall(ipAddress, adminInfoRpcCall, nodeInfoResponse)
               if err != nil {
                   return "", stacktrace.Propagate(err, "Failed to send admin node info RPC request to geth node.")
               }
               return nodeInfoResponse.Result.Enode, nil
           }
           ```
        1. Write logic to check for the availability of your Ethereum node
            1. Add the following code to the bottom of your `Setup()` method
               ```
               firstNodeUp := false
               waitForStartupMaxPolls := 90
               waitForStartupTimeBetweenPolls := 1 * time.Second
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
               ```
        1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (myTest) and that it passed
1. Write test logic in the `Run()` method to verify basic functionality of the single node Ethereum network.
    1. Write test logic to get and verify the Chain ID of the test chain
        1. Create an Eth Client private helper function `getEthClient` at the bottom of the test file, so it can be used later.
        ```
        func getEthClient(ipAddress string) (*ethclient.Client, error) {
            rpcPort := 8545
            url := fmt.Sprintf("http://%v:%v", ipAddress, rpcPort)
            client, err := ethclient.Dial(url)
            if err != nil {
                return nil, stacktrace.Propagate(err, "An error occurred getting the Golang Ethereum client")
            }
            return client, nil
        }
        ```
        1. Create an Eth Client using the Kurtosis network context in the `Run()` method
        ```
        // Necessary because Go doesn't have generics
        castedNetwork := uncastedNetwork.(*networks.NetworkContext)
           
        serviceCtx, err := castedNetwork.GetServiceContext("my-eth-client")
        if err != nil {
           return stacktrace.Propagate(err, "An error occurred getting the Ethereum Go Client service context")
        }
        logrus.Infof("Got service context for Ethereum Go Client service '%v'", serviceCtx.GetServiceID())
           
        gethClient, err := getEthClient(serviceCtx.GetIPAddress())
        if err != nil {
           return stacktrace.Propagate(err, "Failed to get a gethClient from first node.")
        }
        defer gethClient.Close()
        ```
        1. Use the Eth Client to get and print the Network ID.
        ```
        networkId, err := gethClient.NetworkID(context.Background())
        if err != nil {
           return stacktrace.Propagate(err, "Failed to get network ID")
        }
        logrus.Infof("Chain ID: %v", networkId)
        ```   
        1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (myTest) and that it passed
1. Send a transaction on the blockchain running in your single-node Ethereum testnet.
    1. Use the `ExecCommand()` on the Kurtosis service context to execute commands from the [official Geth documentation](https://geth.ethereum.org/docs/getting-started/dev-mode).
    ```
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
    ```
    1. Verify that running `bash scripts/build-and-run.sh all` shows one passing test (my_test) that contains the business logic for an Ethereum single node network

## Implement an advanced test which tests execution of a smart contract on a multiple-node Ethereum network

1. Create a multiple-node private Ethereum testnet in Kurtosis.
    1. Add the advanced test object in `my_advanced_test_.go` to the testsuite object in `my_testsuite.go`.
    1. Setup an Ethereum bootnode for the advanced test.
        1. Implements the Setup() method of the test `my_advanced_test_.go` in order to start the Ethereum private network with multiple nodes
            1. In your preferred IDE, open the advanced Ethereum test `my_advanced_test` at `testsuite/testsuite_impl/my_advanced_test/my_advanced_test.go`
            1. Add the service to the testsuite's network
                1. Set the container configuration for the Ethereum container in your testnet.
                    1. Set the container configuration's static files for an Ethereum bootnode
                        1. Set following statics files:`genesis.json`, `password.txt` and `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` using the IDs previously set in the `Configure` method, with the `WithStaticFiles()` method of the builder
                            1. Add the following code to set a variable that you can use in `WithStaticFiles()` method
                            ```
                            staticFiles := map[services.StaticFileID]bool{
                                genesisStaticFileID: true,
                                passwordStaticFileID: true,
                                signerKeystoreFileID: true
                            }
                            ```
                    1. Set the container configuration`s used ports for an Ethereum bootnode
                        1. Set ports `8545` with `tcp` protocol and port `30303` with `tcp` and `udp` protocol, using the `WithUsedPorts()` method of the builder
                            1. Add the following code to set a variable that you can use in `WithUsedPorts()` method
                            ```
                            usedPorts := map[string]bool{
                                   fmt.Sprintf("%v/tcp", 8545): true,
                                   fmt.Sprintf("%v/tcp", 30303): true,
                                   fmt.Sprintf("%v/udp", 30303): true,
                            },
                            ```
                1. Set the runtime configuration for the Ethereum container in your testnet.
                    1. Get the filepath of `genesis.json` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                        1. Add the following code into the anonymous function that defines the runtime configuration, to set a variable with the genesis's file location
                        ```
                        genesisFilepath, found := staticFileFilepaths[genesisStaticFileID]
                        if !found {
                            return nil, stacktrace.NewError("No filepath found for key '%v'; this is a bug in Kurtosis!", genesisStaticFileID)
                        }
                        ```
                    1. Get the filepath of `password.txt` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                    1. Get the filepath of `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
                    1. Create the entrypoint that contains the command to execute the Ethereum node into the container
                        1. Write the command to init the genesis block
                        1. Set `datadir` value
                        1. Set the custom genesis block file location using the `folder` of the `filepath` generated on the previous step
                    1. Write the command to start the bootnode. It should be written after the init command and using the `&&` operator to execute them sequentially
                        1. Set the `datadir` option
                        1. Set the `keystore` option
                        1. Set the `network ID` option, remember that the value has being defined in the `genesis.json` file
                        1. Enable `HTTP-RPC` server
                        1. Set the `IP address` of the `HTTP-RPC` server
                        1. Set the `API's offered over the HTTP-RPC interface` using these values `admin,eth,net,web3,miner,personal,txpool,debug`
                        1. Accept cross-origin requests from any domain using this value `*`
                        1. Set the `IP address` of the node
                        1. Set the `port` of the node
                        1. Unlock the `signer account` to allow it to mine, remember that you can get this account address from the keystore file `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026`
                        1. Enable `mine`
                        1. Allow `insecure account unlocking` to allow signer account to mine
                        1. Set and Network IP to restrict network communication using a CIDR mask
                            1. Add the private helper function `getIPNet()` 
                               ```
                               func getIPNet(ipAddr string) *net.IPNet {
                                   cidr := ipAddr + subnetRange
                                   _, ipNet, _ := net.ParseCIDR(cidr)
                                   return ipNet
                               }
                               ```
                            1. Call the private helper `getIPNet()` and passes it the bootnode's IP address
                        1. Set the `filepath of the password file` that allows you to avoid entering the password manually when you want to execute a command
                1. Calls the `networkCtx.AddService()` method and passes it the service identifier, the container configuration and the runtime configuration
            1. Verify if everything is working well on this test at this point
                1. Run `scripts/build-and-run.sh all --tests myAdvancedTest`
                1. Verify that the output of the build-and-run.sh script indicates that one test ran (myAdvancedTest) and that it passed.
            1. Checks if the bootnode service is up and running
                1. Use the `networkCtx.WaitForEndpointAvailability()` method to check availability
        1. Get the bootnode's ENR address            
            1. Execute the following Geth command inside the service to get the ENR address
            ```
            exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
               "/bin/sh",
               "-c",
               fmt.Sprintf("geth attach data/geth.ipc --exec admin.nodeInfo.enr"),
            })
            if err != nil {
               return "", stacktrace.Propagate(err, "Executing command returned an error: %v", err.Error())
            }
            if exitCode != 0 {
               return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
            }
            ```
            1. Cast the logout value, receive for the command, to an string which contains the bootnode's ENR address
    1. Start the remaining nodes with the help of the bootnode
        1. Add the first node
            1. Set the container configuration for the Ethereum node container in your testnet.
                1. Load file `genesis.json` using the ID previously set in the `Configure` method, with the `WithStaticFiles()` method of the builder
                1. Set ports `8545` with `tcp` protocol and port `30303` with `tcp` and `udp` protocol, using the `WithUsedPorts()` method of the builder
            1. Set the runtime configuration for the Ethereum node container in your testnet.
                1. Get the filepath of `genesis.json` file using the `staticFileFilepaths` map that the anonymous function receives as a parameter
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
                        1. Accept cross-origin requests from any domain using this value `*`
                        1. Set the `IP address` of the node
                        1. Set the `port` of the node
                        1. Set the `bootnode`using the ENR address previously get
            1. Calls the `networkCtx.AddService()` method and passes it the service identifier, the container configuration and the runtime configuration
            1. Verify again if everything is working well on this test at this point
               1. Run `scripts/build-and-run.sh all --tests myAdvancedTest`
               1. Verify that the output of the build-and-run.sh script indicates that one test ran (myAdvancedTest) and that it passed.
        1. Checks if the node is up and running 
            1. Use the `networkCtx.WaitForEndpointAvailability()` method to check availability
    1. Get the `Enode` address that will be used to connect with the remaining nodes
        1. Add a helper function `sendRpcCall` to send RPC calls to the Ethereum node to the bottom of the test file.
        ```
            func sendRpcCall(ipAddress string, rpcJsonString string, targetStruct interface{}) error {
                rpcPort := 8545
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
        1. Add a helper function `getEnodeAddress` to get the Ethereum node's enode address
        ```
        func getEnodeAddress(ipAddress string) (string, error) {
                nodeInfoResponse := new(NodeInfoResponse)
                err := sendRpcCall(ipAddress, adminInfoRpcCall, nodeInfoResponse)
                if err != nil {
                    return "", stacktrace.Propagate(err, "Failed to send admin node info RPC request to geth node with ip %v", ipAddress)
                }
                return nodeInfoResponse.Result.Enode, nil
            }
        ```
    1. Connect the node with peers (this must be done from the second loaded node)
        1. Connect the node manually using the command `admin_addPeer` [explained on this document](https://geth.ethereum.org/docs/rpc/ns-admin#admin_addpeer)
    1. Check the link between nodes to ensure interconnectivity of the private network it should contain a validation to validate the numbers of peers that all the nodes have.
        1. List the peers using the command `admin_peers` [explained on this document](https://geth.ethereum.org/docs/rpc/ns-admin#admin_peers)
    1. Repeat the previous steps in order to start and link the remaining nodes
    1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that two tests ran (myTest, and myAdvancedTest) and that them passed
    
1. Deploy the `hello_world` smart contract into the private network to test an Ethereum transaction
    1. Write test logic in the `Run()` method to verify advanced functionality of the private multiple node Ethereum network.
        1. Create an instance of the Geth client which will be necessary to deploy the smart contract
        1. Get the private key's signer account
            1. Get the content of signer account's keystore file through the `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` file and the content of the password through the `password.txt` file previously loaded into the testsuite
            1. Get the file paths using `serviceCtx.LoadStaticFiles()` method
            1. Get files' content using Golang `ioutil` package
        1. Get the signer account's private key through the function `keystore.DecryptKey` and passes it the values get from the keystore file, and the password file
        1. Get a transactor which will be necessary to execute the deployment, we suggest using the function `bind.NewKeyedTransactorWithChainID` to accomplish this step
        1. Set a transactor's gas price
        1. Execute an Ethereum transaction to deploy the smart contract
            1. Use function `DeployHelloWorld` provided by the binding `hello_world.go` which will receive the transactor and the Geth client as arguments  
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
1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that two test ran (myTest, and myAdvancedTest) and both are successfully executed
