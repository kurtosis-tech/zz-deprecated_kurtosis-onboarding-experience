Kurtosis Ethereum Testsuite Tutorial
====================================
The instructions below will walk you through spinning up an Ethereum network in Kurtosis, interacting with it, and migrating the logic into the Kurtosis testing framework. By the end of this tutorial, you will have a rudimentary Ethereum testsuite in Typescript that you can begin to modify on your own.


Step One: Set Up Prerequisites (2 minutes)
------------------------------------------
Verify that you have the Docker daemon installed and running on your local machine by running (you can copy this code by hovering over it):

```
docker image ls
```

* If you don't have Docker installed, do so by following [the installation instructions](https://docs.docker.com/get-docker/)
* If Docker is installed but not running, start it

Step Two: Start A Sandbox Enclave (3 minutes)
---------------------------------------------
The Kurtosis engine provides you isolated environments called "enclaves" to run your services inside. Let's use the CLI to start a sandbox enclave:

1. Download the Kurtosis CLI (this can also be copied by hovering):
    ```
    brew install kurtosis-tech/tap/kurtosis
    ```
1. Start a sandbox enclave:
    ```
    mkdir /tmp/my-enclave && cd /tmp/my-enclave && kurtosis sandbox
    ```

The Kurtosis images that run the engine will take a few seconds to pull the first time, but once done you'll have a Javascript REPL with tab-complete attached to your enclave.

All interaction with a Kurtosis enclave is done via [a client library][core-documentation], whose entrypoint is the `NetworkContext` object - a representation of the network running inside the enclave. The `networkCtx` variable in your REPL is how you'll interact with your enclave.

Let's check the contents of our enclave (code snippets can be copied by hovering):

```javascript
getServicesResult = await networkCtx.getServices()
services = getServicesResult.value
```

We haven't started any services yet, so the enclave will be empty. Note how we called `await` on `networkCtx.getServices()`. This is because every `networkCtx` call is asynchronous and returns a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise); `await`ing blocks until that value is available.


Step Three: Start An Ethereum Network (5 minutes)
-------------------------------------------------
Now that we have an enclave, let's put something in it! Ethereum is one of the most popular blockchains in the world, so let's get a private Ethereum network running:

```javascript
loadEthLambdaResult = await networkCtx.loadLambda("eth-lambda", "kurtosistech/ethereum-kurtosis-lambda:0.2.3", "{}")
ethLambdaCtx = loadEthLambdaResult.value
executeEthLambdaResult = await ethLambdaCtx.execute("{}")
executeEthLambdaResultObj = JSON.parse(executeEthLambdaResult.value)
console.log(executeEthLambdaResultObj)
```

After loading, you'll see a result with information about the services running inside your enclave:

```javascript
{
  bootnode_service_id: 'bootnode',
  node_info: {
    bootnode: {
      ip_addr_inside_network: '175.152.144.6',
      exposed_ports_set: [Object],
      port_bindings_on_local_machine: [Object]
    },
    'ethereum-node-1': {
      ip_addr_inside_network: '175.152.144.8',
      exposed_ports_set: [Object],
      port_bindings_on_local_machine: [Object]
    },
    'ethereum-node-2': {
      ip_addr_inside_network: '175.152.144.10',
      exposed_ports_set: [Object],
      port_bindings_on_local_machine: [Object]
    }
  },
  signer_keystore_content: '{"address":"14f6136b48b74b147926c9f24323d16c1e54a026","crypto":{"cipher":"aes-128-ctr","ciphertext":"39fb1d86c1082c0103ece1c5f394321f127bf1b65e6c471edcfb181058a3053a","cipherparams":{"iv":"c366d1eed33e8693fec7a85fad65d19f"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"f210bc3b55117197f62a7ab8d85f2172342085f1daafa31034016163b8bc7db6"},"mac":"2ff8aa24d9b73ccfdb99cfd15fcdbcc8f640aaa7861e6813d53efaf550725fac"},"id":"6c5ac271-d24a-4971-b365-49490cc4befc","version":3}',
  signer_account_password: 'passphrase'
}
```

And if we query the enclave's services again...

```javascript
getServicesResult = await networkCtx.getServices()
services = getServicesResult.value
```

...we see three service IDs:

```javascript
Set(3) { 'bootnode', 'ethereum-node-1', 'ethereum-node-2' }
```

But what just happened?

Starting networks is a very common task in Kurtosis, so we provide [a framework called "Lambdas"](https://docs.kurtosistech.com/lambdas.html) for making it dead simple. A Lambda is basically a function, packaged as a Docker image, that runs inside a Kurtosis enclave - sort of like Docker compose on steroids. In the steps above, we called `networkCtx.loadLambda` to load [the Ethereum Lambda](https://github.com/kurtosis-tech/ethereum-kurtosis-lambda) into the enclave with Lambda ID `eth-lambda`, and `networkCtx.executeLambda` to run it. The Ethereum Lambda doesn't take any arguments at load or execute time (hence the `{}`), but other Lambdas do.

Now that you have a pet Ethereum network, let's do something with it.

Step Four: Talk To Ethereum (5 minutes)
---------------------------------------
Talking to Ethereum in Javascript is easily accomplished with [the EthersJS library](https://docs.ethers.io/v5/), so let's load it. 

Even though the REPL is running in a Docker image so that you don't need to have `node` installed locally, the directory in which you ran `kurtosis sandbox` (`/tmp/my-enclave`) is available to your REPL. To load Ethers:

1. Unzip [this file](https://kurtosis-public-access.s3.us-east-1.amazonaws.com/onboarding-artifacts/ethers-js-5.4.7.zip) into `/tmp/my-enclave`
1. Run:
    ```javascript
    const ethers = require("ethers")
    ```

Now let's get a connection to the node with service ID `bootnode` by getting a [JsonRpcProvider](https://docs.ethers.io/v5/api/providers/jsonrpc-provider/):

```javascript
bootnodeServiceId = executeEthLambdaResultObj.bootnode_service_id
bootnodeIp = executeEthLambdaResultObj.node_info[bootnodeServiceId].ip_addr_inside_network
bootnodeRpcProvider = new ethers.providers.JsonRpcProvider(`http://${bootnodeIp}:8545`);
```

Notice how we used the `executeEthLambdaResultObj` object containing details about the Ethereum network, which came back from executing the Lambda.

Finally, let's verify that our Ethereum netowrk is producing blocks:

```javascript
blockNumber = await bootnodeRpcProvider.getBlockNumber()
if (blockNumber > 0) { console.log("All is well!"); }
```

And that's it! Anything doable in Ethers is now doable against your private Ethereum network running in your Kurtosis enclave.

To exit out of the REPL you can use any of:

* Ctrl-D
* Ctrl-C, twice
* `.exit`

Step Five: Get An Ethereum Testsuite (5 minutes)
---------------------------------------------
Manually verifying against a sandbox network is nice, but it'd be great if we could take our logic above and run it as part of CI. Kurtosis has a testing framework that allows us to do exactly this. 

Normally, we'd bootstrap from [the Testsuite Starter Pack](https://github.com/kurtosis-tech/kurtosis-testsuite-starter-pack) and use [the same Kurtosis engine documentation][core-documentation] with [the testing framework documentation](https://docs.kurtosistech.com/kurtosis-testsuite-api-lib/lib-documentation) to write our tests.

For the purposes of this onboarding though, we've already done the busywork of copying the previous steps into a Typescript testsuite. Go ahead and clone it from [here](https://github.com/kurtosis-tech/kurtosis-testsuite-starter-pack) now, and we'll take a look around.

The first thing to notice is the `testsuite/Dockerfile`. Testsuites in Kurtosis are simply packages of tests bundled in Docker images, which the testing framework will instantiate to run tests.

The second thing to notice is the `testsuite/testsuite_impl/eth_testsuite.ts` file. This is where tests are defined, and this testsuite already has a single test - `basicEthTest`.

Now open `testsuite/testsuite_impl/basic_eth_test/basic_eth_test.ts`. You'll see that a test is really just a class with three function: `configure`, `setup`, and `run`. Like most testing frameworks, `setup` is where we place the prep work that executes before the `run` method while `run` is where we make our test assertions. The `configure` method is where timeouts for both `setup` and `run` are configured, among other things.

The last thing to notice is how a `NetworkContext` is passed in as an argument to `setup`. Every Kurtosis test runs inside its own enclave to prevent cross-test interference, and you can use [the exact same `NetworkContext` APIs][core-documentation] for both your sandbox and testing needs.

Now let's see the testing framework in action. From the root of the repo, run:

```
scripts/build-and-run.sh all    # The 'all' tells Kurtosis to build your testsuite into a Docker image AND run it
```

You'll see a prompt to create a Kurtosis account, which we use for gating advanced features (don't worry, we won't sign you up for any email lists!). Follow the instructions, and click the device verification link once you have your account.

The testsuite will run, and you'll see that our `basicEthTest` passed!

Under the covers, Kurtosis is creating an enclave 

Step Six: Test Ethereum (5 minutes)
-----------------------------------
We have Ethereum running inside a test, but our test doesn't currently do anything. Let's fix that.

First, replace the `// TODO Replace with Ethereum network setup` line in the `setup` method with the following code:

```typescript
log.info("Setting up Ethereum network...")
const loadEthLambdaResult: Result<LambdaContext, Error> = await networkCtx.loadLambda(ETH_LAMBDA_ID, ETH_LAMBDA_IMAGE, "{}");
if (loadEthLambdaResult.isErr()) {
    return err(loadEthLambdaResult.error);
}
const ethLambdaCtx: LambdaContext = loadEthLambdaResult.value;

const executeEthLambdaResult: Result<string, Error> = await ethLambdaCtx.execute("{}")
if (executeEthLambdaResult.isErr()) {
    return err(executeEthLambdaResult.error);
}
this.executeEthLambdaResultObj = JSON.parse(executeEthLambdaResult.value);
log.info("Ethereum network set up successfully");
```

This is the same code we already executed in the REPL, cleaned up for Typescript. The only new bits to pay attention to are the error-checking: all `NetworkContext` methods, as well as the `Test.setup` and `Test.run` methods, return [a Result object][neverthrow] (much like in Rust). If `setup` or `run` return a non-`Ok` result, the test will be marked as failed. This allows for easy, consistent error-checking: simply propagate the error upwards.

Second, replace the `// TODO Replace with block number check` line with this code:

```typescript
log.info("Verifying block number is increasing...");
const bootnodeServiceId: ServiceID = this.executeEthLambdaResultObj.bootnode_service_id;
const bootnodeIp: string = this.executeEthLambdaResultObj.node_info[bootnodeServiceId].ip_addr_inside_network
const bootnodeRpcProvider: ethers.providers.JsonRpcProvider = new ethers.providers.JsonRpcProvider(`http://${bootnodeIp}:8545`);
const blockNumber: number = await bootnodeRpcProvider.getBlockNumber();
if (blockNumber === 0) {
    return err(new Error(""))
}
log.info("Verified that block number is increasing");
```

Finally, build and run the testsuite again:

```
scripts/build-and-run.sh all
```

You'll see our setup and verification logic running, and the tests passing!

<!-- TODO Link to docs and further deepdives -->

<!-- TODO explain extra flags to control testsuite execution -->
<!-- TODO explain executing the testsuite in CI -->
<!-- TODO explain Debug mode, host port bindings, and setting debug log level -->

[neverthrow]: https://www.npmjs.com/package/neverthrow
[core-documentation]: https://docs.kurtosistech.com/kurtosis-client/lib-documentation
