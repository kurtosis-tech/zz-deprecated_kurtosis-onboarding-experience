Ethereum On-Boarding Testsuite
==============================

## Implement a Single Node Ethereum Test Network

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
        1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
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
        1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
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
