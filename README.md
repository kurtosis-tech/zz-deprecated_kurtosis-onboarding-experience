Ethereum On-Boarding Testsuite
==============================

## Implement a Single Node Ethereum Test Network

1. Create an account on [https://www.kurtosistech.com/sign-up](https://www.kurtosistech.com/sign-up) if you don't have one yet.
1. Verify that the Docker daemon is running on your local machine using `docker container ls`
1. Clone this repository by running `git clone https://github.com/kurtosis-tech/kurtosis-onboarding-experience.git --branch master`
1. Run the empty Ethereum single node test `my_test` to verify that everything works on your local machine.
    1. Run `bash scripts/build-and-run.sh all`
    1. Verify that the output of the build-and-run.sh script indicates that one test ran (my_test) and that it passed.
1. Import the dependencies that are used in this example test suite.
   1. Run `go get github.com/ethereum/go-ethereum/ethclient`
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
          
            serviceCtx, hostPortBindings, err := networkCtx.AddService("my-eth-client", containerCreationConfig, runConfigFunc)
            if err != nil {
                return nil, stacktrace.Propagate(err, "An error occurred adding the Ethereum Go Client service")
            }
            ```
    1. Add an availability check to the `Setup()` method to make sure your Ethereum node is fully functional before your test starts.
        1. Write logic to check for the availability of your Ethereum node
            1. Add the following code to the bottom of your `Setup()` method in order to use Kurtosis function `networkCtx.WaitForEndpointAvailability()`
               ```
               adminInfoRpcCall  := `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`
               if err := networkCtx.WaitForEndpointAvailability("my-eth-client", kurtosis_core_rpc_api_bindings.WaitForEndpointAvailabilityArgs_POST, 8545, "", adminInfoRpcCall, 1, 30, 1, ""); err != nil {
                    return "", stacktrace.Propagate(err, "An error occurred waiting for service with ID '%v' to start", "bootnode")
               }
            
               logrus.Infof("Added Ethereum Go Client service with IP: %v andhost port bindings: %+v", serviceCtx.GetIPAddress(), hostPortBindings)
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

## Implement an Advanced Test which test and Ethereum Private Network with Multiple Nodes

1. Create a private Ethereum test network in Kurtosis with multiple nodes, that uses **Clique consensus** as proof of authority and that is previously set in the **genesis block**, with a **signer account**.
    1. Setup a bootnode first
        1. Edit the `Configure()` method of the test `my_test.go`
           1. Replace the body with the following code
           ```
           builder.WithSetupTimeoutSeconds(
               240,
           ).WithRunTimeoutSeconds(
               240,
           ).WithStaticFileFilepaths(map[services.StaticFileID]string{
               genesisStaticFileID:  genesisStaticFilePath,
               signerKeystoreFileID: signerKeystoreFilePath,
               passwordStaticFileID: passwordStaticFilePath,
           })
           ```
        1. Edit the Setup() method of the test `my_test.go` in order to start the Ethereum private network with multiple nodes
            1. In your preferred IDE, open the recent edited test `my_test` at `testsuite/testsuite_impl/my_test/my_test.go`
            1. Edit the container configuration to setup an Ethereum bootnode container
               1. Replace the `getContainerCreationConfig()` function body with this body
               ```
               rpcPort := 8545
               discoveryPort := 30303
               staticFiles := map[services.StaticFileID]bool{
                   genesisStaticFileID:  true,
                   passwordStaticFileID: true,
                   signerKeystoreFileID: true,
               }
            
               containerCreationConfig := services.NewContainerCreationConfigBuilder(
                    "ethereum/client-go",
               ).WithUsedPorts(
                   map[string]bool{
                       fmt.Sprintf("%v/tcp", rpcPort):       true,
                       fmt.Sprintf("%v/tcp", discoveryPort): true,
                       fmt.Sprintf("%v/udp", discoveryPort): true,
                   },
               ).WithStaticFiles(staticFiles).Build()
               
               return containerCreationConfig
               ```
            1. Edit the runtime configuration for the Ethereum container in your testnet.
                1. Add the following `getIPNet()` helper function to the bottom of the test file.
                ```
                   func getIPNet(ipAddr string) *net.IPNet {
                       subnetRange := "/24" 
                       cidr := ipAddr + subnetRange
                       _, ipNet, _ := net.ParseCIDR(cidr)   
                       return ipNet
                   }
                ```
                1. Replace the body of the `getRunConfigFunc()` function
                ```
                   discoveryPort := 30303
                   runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
                       genesisFilepath, found := staticFileFilepaths[genesisStaticFileID]
                       if !found {
                            return nil, stacktrace.NewError("No filepath found for key '%v'; this is a bug in Kurtosis!", genesisStaticFileID)
                       }
                
                       passwordFilepath, found := staticFileFilepaths[passwordStaticFileID]
                       if !found {
                            return nil, stacktrace.NewError("No filepath found for key '%v'; this is a bug in Kurtosis!", passwordStaticFileID)
                       }
                
                       signerKeystoreFilepath, found := staticFileFilepaths[signerKeystoreFileID]
                       if !found {
                            return nil, stacktrace.NewError("No filepath found for key '%v'; this is a bug in Kurtosis!", signerKeystoreFileID)
                       }
                
                       keystoreFolder := filepath.Dir(signerKeystoreFilepath)
                
                       ipNet := getIPNet(ipAddr)
                
                       entryPointArgs := []string{
                           "/bin/sh",
                           "-c",
                           fmt.Sprintf("geth init --datadir data %v && geth --keystore %v --datadir data --networkid 15 -http --http.api admin,eth,net,web3,miner,personal,txpool,debug --http.addr %v --http.corsdomain '*' --nat extip:%v --port %v --unlock 0x14f6136b48b74b147926c9f24323d16c1e54a026 --mine --allow-insecure-unlock --netrestrict %v --password %v", genesisFilepath, keystoreFolder, ipAddr, ipAddr, discoveryPort, ipNet, passwordFilepath),
                       }
                
                       result := services.NewContainerRunConfigBuilder().WithEntrypointOverride(entryPointArgs).Build()
                       return result, nil
                   }
                   return runConfigFunc
                ```
            1. Rename the `serviceId` parameter value to `bootnode` in the `AddService()` call inside the `Setup()` method
            1. Checks if the bootnode service is up and running
                1. Rename the `serviceId` parameter value to `bootnode` in the `WaitForEndpointAvailability()` call inside the `Setup()` method
        1. Get the bootnode's ENR address            
            1. Add the following lines before the return sentence in the `Setup()` method
            ```
            exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
               "/bin/sh",
               "-c",
               fmt.Sprintf("geth attach data/geth.ipc --exec admin.nodeInfo.enr"),
            })
            if err != nil {
               return "", stacktrace.Propagate(err,"Executing command returned an error with logs: %+v", string(*logOutput))
            }
            if exitCode != 0 {
                return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
            }
        
            bootNodeENR := string(*logOutput)
           
            logrus.Infof("Bootnode ENR address: %v", bootNodeENR)
            ```
    1. Start the remaining nodes with the help of the bootnode
       1. Set a new container configuration (used by the remaining nodes) for the Ethereum container in your testnet.
          1. Add the following `getContainerCreationConfigForETHNode()` helper function at the bottom of the test file
          ```
          func getContainerCreationConfigForETHNode() *services.ContainerCreationConfig {
               rpcPort := 8545
               discoveryPort := 30303
               staticFiles := map[services.StaticFileID]bool{
                   genesisStaticFileID:  true,
                   passwordStaticFileID: true,
               }
            
               containerCreationConfig := services.NewContainerCreationConfigBuilder(
                    "ethereum/client-go",
               ).WithUsedPorts(
                   map[string]bool{
                       fmt.Sprintf("%v/tcp", rpcPort):       true,
                       fmt.Sprintf("%v/tcp", discoveryPort): true,
                       fmt.Sprintf("%v/udp", discoveryPort): true,
                   },
               ).WithStaticFiles(staticFiles).Build()
          
               return containerCreationConfig
          }
          ```
       1. Set a new runtime configuration (used by the remaining nodes) for the Ethereum container in your testnet.
          1. Add the following `getRunConfigFuncForETHNode()` helper function at the bottom of the test file
          ```
           func getRunConfigFuncForETHNode(bootnodeEnr string) func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
                discoveryPort := 30303
                runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
                   genesisFilepath, found := staticFileFilepaths[genesisStaticFileID]
                   if !found {
                        return nil, stacktrace.NewError("No filepath found for test file 1 key '%v'; this is a bug in Kurtosis!", genesisStaticFileID)
                   }
        
                   entryPointArgs := []string{
                       "/bin/sh",
                       "-c",
                       fmt.Sprintf("geth init --datadir data %v && geth --datadir data --networkid 15 -http --http.api admin,eth,net,web3,miner,personal,txpool,debug --http.addr %v --http.corsdomain '*' --nat extip:%v --gcmode archive --syncmode full --port %v --bootnodes %v", genesisFilepath, ipAddr, ipAddr, discoveryPort, bootnodeEnr),
                   }
        
                   result := services.NewContainerRunConfigBuilder().WithEntrypointOverride(entryPointArgs).Build()
                   return result, nil
               }
               return runConfigFunc
           } 
          ```
       1. Add the following ```` helper function at the bottom of the test file
       ```
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
       1. Add the following `AddPeer()` helper function at the bottom of the test file
       ```
       func AddPeer(ipAddress string, peerEnode string) (bool, error) {
           adminAddPeerRpcCall := fmt.Sprintf(`{"jsonrpc":"2.0", "method": "admin_addPeer", "params": ["%v"], "id":70}`, peerEnode)
           addPeerResponse := new(AddPeerResponse)
           err := sendRpcCall(ipAddress, adminAddPeerRpcCall, addPeerResponse)
           if err != nil {
                return false, stacktrace.Propagate(err, "Failed to send addPeer RPC call for enode %v", peerEnode)
           }
           return addPeerResponse.Result, nil
       }
       ```
       1. Add the following `validatePeersQuantity()` helper function at the bottom of the test file
       ```
       func validatePeersQuantity(logString string, serviceID services.ServiceID, nodesEnode []string) error {
           enodePrefix := "enode://"
           handshakeProtocol := "eth: \"handshake\""
           peersQuantity := strings.Count(logString, enodePrefix) - strings.Count(logString, handshakeProtocol)
           validPeersQuantity := len(nodesEnode) + 1
           if peersQuantity != validPeersQuantity {
           return stacktrace.NewError("The amount of peers '%v' for node '%v' is not valid, should be '%v?", peersQuantity, serviceID, validPeersQuantity)
           }
           return nil
       } 
       ```
       1. Add the following `validatePeersQuantity()` helper function at the bottom of the test file
       ```
       func getEnodeAddress(ipAddress string) (string, error) {
           nodeInfoResponse := new(NodeInfoResponse)
           adminInfoRpcCall := `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`
           err := sendRpcCall(ipAddress, adminInfoRpcCall, nodeInfoResponse)
           if err != nil {
                return "", stacktrace.Propagate(err, "Failed to send admin node info RPC request to geth node with ip %v", ipAddress)
           }
           return nodeInfoResponse.Result.Enode, nil
       }
       ```   
       1. Add the following `starEthNodeByBootnode()` private helper function which start a ETH node using the bootnode and checks for its peers
       ```
       _func starEthNodeByBootnode(networkCtx *networks.NetworkContext, serviceID services.ServiceID, bootnodeEnr string, nodesEnode []string) (string, error) {
            containerCreationConfig := getContainerCreationConfigForETHNode()
            runConfigFunc := getRunConfigFuncForETHNode(bootnodeEnr)
                
            serviceCtx, hostPortBindings, err := networkCtx.AddService(serviceID, containerCreationConfig, runConfigFunc)
            if err != nil {
               return "", stacktrace.Propagate(err, "An error occurred adding the Ethereum Node %v", serviceID)
            }
        
            logrus.Infof("Added Ethereum Node %v service with host port bindings: %+v and IP address '%v'", serviceID, hostPortBindings, serviceCtx.GetIPAddress())
        
            adminInfoRpcCall := `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`
            if err := networkCtx.WaitForEndpointAvailability(serviceID, kurtosis_core_rpc_api_bindings.WaitForEndpointAvailabilityArgs_POST, 8545, "", adminInfoRpcCall, 1, 30, 1, ""); err != nil {
                return "", stacktrace.Propagate(err, "An error occurred waiting for service with ID '%v' to start", serviceID)
            }
        
            for _, enode := range nodesEnode {
               ok, err := AddPeer(serviceCtx.GetIPAddress(), enode)
               if err != nil {
                    return "", stacktrace.Propagate(err, "Failed to call addPeer endpoint to add peer with enode %v", enode)
               }
               if !ok {
                    return "", stacktrace.NewError("addPeer endpoint returned false on service %v, adding peer %v", serviceID, enode)
               }
            }
        
            exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
               "/bin/sh",
               "-c",
               fmt.Sprintf("geth attach data/geth.ipc --exec admin.peers"),
            })
            if err != nil {
                return "", stacktrace.Propagate(err, "Executing command returned an error with logs: %+v", string(*logOutput))
            }
            if exitCode != 0 {
                return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
            }
        
            if err = validatePeersQuantity(string(*logOutput), serviceID, nodesEnode); err != nil {
                return "", stacktrace.Propagate(err, "Validate peers error")
            }
        
            enode, err := getEnodeAddress(serviceCtx.GetIPAddress())
            if err != nil {
                return "", stacktrace.Propagate(err, "Failed to get enode from peer %v", serviceID)
            }
        
            return enode, nil
       }_
       ```
       1. Add the following lines before the return sentence in the `Setup()` method
       ```
       var nodesEnode []string
       numGethNodes := 3
       gethServiceIdPrefix := "ethereum-node-"
       for i := 1; i <= numGethNodes; i++ {
           serviceID := services.ServiceID(gethServiceIdPrefix + strconv.Itoa(i))
           enode, err := starEthNodeByBootnode(networkCtx, serviceID, bootNodeENR, nodesEnode)
           if err != nil {
                return nil, stacktrace.Propagate(err, "An error occurred starting the Ethereum Node '%v'", serviceID)
           }
           nodesEnode = append(nodesEnode, enode)
       }
       ```
    
1. Deploy the the `hello_world` smart contract into the private network to test an Ethereum transaction
    1. Edit the test logic inside the `Run()` method to verify advanced functionality of the private multiple node Ethereum network.
        1. Rename the `serviceId` parameter value to `bootnode` in the `GetServiceContext()` call inside the `Run()` method
        1. Remove all the current code relate to `ExecCommand` executions  //TODO improve this sentence
        1. Add the following `getPrivateKey()` helper function at the bottom of the test file in order to get the signer's private key
        ```
        func getPrivateKey(serviceCtx *services.ServiceContext) (*keystore.Key, error) {
           staticFileAbsFilepaths, err := serviceCtx.LoadStaticFiles(map[services.StaticFileID]bool{
               signerKeystoreFileID: true,
               passwordStaticFileID: true,
           })
           if err != nil {
            return nil, stacktrace.Propagate(err, "An error occurred loading the static files corresponding to keys '%v' and '%v'", signerKeystoreFileID, passwordStaticFileID)
           }
           signerKeystoreFilepath, found := staticFileAbsFilepaths[signerKeystoreFileID]
           if !found {
            return nil, stacktrace.Propagate(err, "No filepath found for key '%v'; this is a bug in Kurtosis!", signerKeystoreFilepath)
           }
        
           signerKeystoreContent, err := ioutil.ReadFile(signerKeystoreFilepath)
           if err != nil {
            return nil, stacktrace.Propagate(err, "An error happens reading file '%v'", signerKeystoreFilepath)
           }
        
           json, err := ioutil.ReadAll(strings.NewReader(string(signerKeystoreContent)))
           if err != nil {
            return nil, stacktrace.Propagate(err,"An error occurred when trying to read content for filepath '%v'", signerKeystoreFilepath)
           }
        
           passwordFilepath, found := staticFileAbsFilepaths[passwordStaticFileID]
           if !found {
            return nil, stacktrace.Propagate(err, "No filepath found for key '%v'; this is a bug in Kurtosis!", passwordFilepath)
           }
        
           passwordContent, err := ioutil.ReadFile(passwordFilepath)
           if err != nil {
            return nil, stacktrace.Propagate(err, "An error happens reading file '%v'", passwordFilepath)
           }
        
           key, err := keystore.DecryptKey(json, string(passwordContent))
           if err != nil {
            return nil, stacktrace.Propagate(err,"An error occurred when trying to decrypt the private key")
           }
           return key, nil
        }
        ```           
        1. Add the following `waitUntilTransactionMined()` helper function at the bottom of the test file in order to use it to check if a transaction has finished
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
        1. Add the following lines before the return sentence in the `Run()` method to execute the Ethereum transaction and checks if works well
        ```
           key, err := getPrivateKey(serviceCtx)
           if err != nil {
                return stacktrace.Propagate(err, "Failed to get private key")
           }
        
           transactor, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(15))
           if err != nil {
                logrus.Fatalf("Failed to create authorized transactor: %v", err)
           }
           transactor.GasPrice = big.NewInt(5)
           address, tx, helloWorld, err := bindings.DeployHelloWorld(transactor, gethClient)
           if err != nil {
                return stacktrace.Propagate(err, "An error occurred deploying the HelloWorld contract on the ETH Network")
           }
           fmt.Printf("Contract pending deploy: 0x%x\n", address)
           fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
        
           if err := waitUntilTransactionMined(gethClient, tx.Hash()); err != nil {
                return stacktrace.Propagate(err, "An error occurred waiting for the HelloWorld contract to be mined")
           }
           logrus.Info("Deployed Hello World contract")
        
           name, err := helloWorld.Greet(&bind.CallOpts{Pending: true})
           if err != nil {
                log.Fatalf("Failed to retrieve pending name: %v", err)
           }
           fmt.Println("Pending name:", name)
        
           listAccountsCmd := []string{
                "/bin/sh",
                "-c",
                fmt.Sprintf("geth attach data/geth.ipc --exec eth.accounts"),
           }
        
           exitCode, logOutput, err := serviceCtx.ExecCommand(listAccountsCmd)
           if err != nil {
                return stacktrace.Propagate(err, "Executing command returned an error with logs: %+v", string(*logOutput))
           }
           if exitCode != 0 {
                return stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
           }
        ```
1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that two test ran (myTest) is successfully executed
