Kurtosis Ethereum Testsuite Tutorial
====================================
This repo is an empty Kurtosis testsuite. The instructions below will walk you through creating a Kurtosis test that spins up a private Ethereum network and runs test logic against it. By the end of this tutorial, you will have seen how Kurtosis testing works.

Step One: Set Up Prerequisites (5min)
-------------------------------------
1. Create a Kurtosis account on [the signup page](https://www.kurtosistech.com/sign-up) if you don't have one yet
1. Verify that you have the Docker daemon installed and running on your local machine by running:

    ```
    docker image ls
    ```

    * If you don't have Docker installed, do so by following [the installation instructions](https://docs.docker.com/get-docker/)
    * If Docker is installed but not running, start it
1. Clone this repository by running the following (**NOTE:** you can copy this snippet in Github by hovering over it and clicking the clipboard in the top-right corner):

    ```
    git clone https://github.com/kurtosis-tech/kurtosis-onboarding-experience.git --branch master && cd kurtosis-onboarding-experience
    ```

1. Verify that the testsuite runs on your local machine with:

    ```
    bash scripts/build-and-run.sh all
    ```

1. Ensure the output indicates that one test, `BasicEthereumTest`, ran and passed

Step Two: Fill In BasicEthereumTest (15min)
-----------------------------------------------------------------
### Configure the test to launch a private Ethereum network inside Kurtosis (5min)
`BasicEthereumTest` currently doesn't do anything, so we'll configure it to instantiate an Ethereum network of one node:

1. In your preferred IDE, open `BasicEthereumTest` inside file `testsuite/testsuite_impl/basic_ethereum_test/basic_ethereum_test_.go`
1. At the bottom of the file under the `Private helper functions` header, replace the `//TODO Container creation & run config helper functions` line with the following helper functions for creating & running an Ethereum node container (**NOTE:** you can copy this entire code snippet by hovering over the block and clicking the clipboard icon in the top-right corner):

    ```golang
    func getContainerCreationConfig() *services.ContainerCreationConfig {
        containerCreationConfig := services.NewContainerCreationConfigBuilder(
            ethereumNodeImage,
        ).WithUsedPorts(
            map[string]bool{fmt.Sprintf("%v/tcp", ethereumNodeRpcPort): true},
        ).Build()
        return containerCreationConfig
    }

    func getRunConfigFunc() func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
        runConfigFunc := func(ipAddr string, generatedFileFilepaths map[string]string, staticFileFilepaths map[services.StaticFileID]string) (*services.ContainerRunConfig, error) {
            entrypointCommand := fmt.Sprintf(
                "geth --dev -http --http.api admin,eth,net,rpc --http.addr %v --http.corsdomain '*' --nat extip:%v",
                ipAddr,
                ipAddr,
            )
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

1. In the test's `Setup` method, replace the `//TODO Start a single ETH node in dev mode` line with the following code so that the test instantiates an Ethereum node as part of its setup:

    ```golang
    containerCreationConfig := getContainerCreationConfig()
    runConfigFunc := getRunConfigFunc()

    serviceCtx, _, err := networkCtx.AddService(node0ServiceID, containerCreationConfig, runConfigFunc)
    if err != nil {
        return nil, stacktrace.Propagate(err, "An error occurred adding the Ethereum node")
    }
    logrus.Infof("Added Ethereum node '%v' with IP '%v'", serviceCtx.GetServiceID(), serviceCtx.GetIPAddress())
    ```

1. In the same `Setup` method, replace `//TODO Check if the Ethereum network is available` with the following code to ensure that the test setup doesn't complete until the Ethereum node is available:

    ```golang
    adminInfoRpcCall  := `{"jsonrpc":"2.0","method": "admin_nodeInfo","params":[],"id":67}`
    if err := networkCtx.WaitForEndpointAvailability(
        serviceCtx.GetServiceID(),
        kurtosis_core_rpc_api_bindings.WaitForEndpointAvailabilityArgs_POST,
        ethereumNodeRpcPort,
        "",
        adminInfoRpcCall,
        nodeAvailabilityCheckInitialDelaySeconds,
        nodeAvailabilityCheckNumRetries,
        nodeAvailabilityCheckRetryWaitMilliseconds,
        "",
    ); err != nil {
        return "", stacktrace.Propagate(err, "An error occurred waiting for Ethereum node '%v' to become available", serviceCtx.GetServiceID())
    }
    logrus.Infof("Ethereum node '%v' is now available", serviceCtx.GetServiceID())
    ```

1. Verify that `BasicEthereumTest` is still passing with:

    ```
    bash scripts/build-and-run.sh all
    ```

1. Ensure that the output contains the following loglines indicating the Ethereum node started successfully (noting that the node's IP will be nondeterministic):

    ```
    Added Ethereum node 'node-0' with IP 'X.X.X.X'
    Ethereum node 'node-0' is now available
    ```

### Configure the test to run test logic against the private Ethereum network (5min)
Now that our test is creating an Ethereum network every time it runs, let's write some logic to interact with it it:

1. Under the `Private helper functions` section, replace the `//TODO Create Go Ethereum client helper function` line with the following:

    ```golang
    func getEthClient(ipAddress string) (*ethclient.Client, error) {
        url := fmt.Sprintf("http://%v:%v", ipAddress, ethereumNodeRpcPort)
        client, err := ethclient.Dial(url)
        if err != nil {
            return nil, stacktrace.Propagate(err, "An error occurred getting the Golang Ethereum client")
        }
        return client, nil
    }
    ```

1. Replace the `//TODO Get Go Ethereum client` line in the test's `Run` method with the following code to get a Go client for interacting with the Ethereum node:

    ```golang
    castedNetwork := uncastedNetwork.(*networks.NetworkContext)
       
    serviceCtx, err := castedNetwork.GetServiceContext(node0ServiceID)
    if err != nil {
       return stacktrace.Propagate(err, "An error occurred getting the service context for Ethereum node '%v'", node0ServiceID)
    }
       
    gethClient, err := getEthClient(serviceCtx.GetIPAddress())
    if err != nil {
       return stacktrace.Propagate(err, "Failed to create a Go client for Ethereum node '%v'", node0ServiceID)
    }
    defer gethClient.Close()
    ```

1. Replace the `//TODO Get ETH network's chain ID` line with the following code for getting the Ethereum network ID (**NOTE:** you might need to instruct your IDE which `context` to import; it should be the package from the Go standard library):

    ```golang
    networkId, err := gethClient.NetworkID(context.Background())
    if err != nil {
        return stacktrace.Propagate(err, "Failed to get Ethereum network ID")
    }
    logrus.Infof("Network ID: %v", networkId)
    ```

1. Verify that `BasicEthereumTest` still passes with:

    ```
    bash scripts/build-and-run.sh all
    ```

1. Ensure that you see the following logline in the output indicating that the test logic ran:

    ```
    Network ID: 1337
    ```

### Extend our test logic to send a transaction to the Ethereum testnet (5min)
We now know that the Ethereum network responds to requests, so let's send a transaction to it:

1. Replace the `//TODO Create new ETH account` line in the test's `Run` method with the following code that uses the Ethereum IPC commands in [the official documentation](https://geth.ethereum.org/docs/getting-started/dev-mode) to create an ETH account using the Ethereum IPC API:

    ```golang
    createAcctExitCode, createAcctLogOutput, err := serviceCtx.ExecCommand([]string{
        "/bin/sh", 
        "-c",
        "printf \"passphrase\\npassphrase\\n\" | geth attach /tmp/geth.ipc --exec 'personal.newAccount()'",
    })
    if err != nil {
       return stacktrace.Propagate(err, "Account creation command returned an error")
    }
    if createAcctExitCode != 0 {
       return stacktrace.NewError("Account creation command returned non-zero exit code with logs:\n%+v", string(*createAcctLogOutput))
    }
    logrus.Info("Account created successfully")
    ```

1. Replace the `//TODO Make ETH transfer transaction` line with the following code to create an ETH transfer transaction:

    ```golang
    sendTransferExitCode, sendTransferLogOutput, err := serviceCtx.ExecCommand([]string{
        "/bin/sh", 
        "-c",
        "geth attach /tmp/geth.ipc --exec 'eth.sendTransaction({from:eth.coinbase, to:eth.accounts[1], value: web3.toWei(0.05, \"ether\")})'",
    })
    if err != nil {
       return stacktrace.Propagate(err, "Send transfer command returned an error")
    }
    if sendTransferExitCode != 0 {
       return stacktrace.NewError("Send transfer command returned non-zero exit code with logs:\n%+v", string(*sendTransferLogOutput))
    }
    logrus.Info("ETH transfer transaction sent successfully")
    ```

1. Replace the `//TODO Get ETH account balance` line with the following code to verify that the account balance got updated:

    ```golang
    getBalanceExitCode, getBalanceLogOutput, err := serviceCtx.ExecCommand([]string{
        "/bin/sh",
        "-c",
        "geth attach /tmp/geth.ipc --exec 'eth.getBalance(eth.accounts[1])'",
    })
    if err != nil {
        return stacktrace.Propagate(err, "Get balance command returned an error")
    }
    if getBalanceExitCode != 0 {
        return stacktrace.NewError("Get balance command returned non-zero exit code with logs:\n%+v", string(*getBalanceLogOutput))
    }
    accountBalanceWeiStr := strings.TrimSpace(string(*getBalanceLogOutput))
    accountBalanceWei, err := strconv.ParseUint(accountBalanceWeiStr, 10, 64)
    if err != nil {
        return stacktrace.Propagate(err, "Couldn't parse account balance Wei string '%v' to number", accountBalanceWeiStr)
    }
    if accountBalanceWei != weiToSend {
        return stacktrace.NewError("Actual account balance '%v' != expected account balance '%v'", accountBalanceWei, weiToSend)
    }
    logrus.Infof("Account balance was increased by %v Wei as expected", weiToSend)
    ```

1. Verify that running the following still shows `BasicEthereumTest` as passing:

    ```
    bash scripts/build-and-run.sh all
    ```
1. Ensure you see the following log lines:

    ```
    Account created successfully
    ETH transfer transaction sent successfully
    Account balance was increased by 50000000000000000 Wei as expected
    ```

And that's it! You now have a Kurtosis test that spins up an Ethereum network, sends a transaction, and verifies the transaction got completed.
