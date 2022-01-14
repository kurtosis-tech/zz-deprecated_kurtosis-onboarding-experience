Kurtosis Ethereum Quickstart
============================
The instructions below will walk you through spinning up an Ethereum network in a Kurtosis sandbox, interacting with it, and migrating the logic into a test. By the end of this tutorial, you will have a rudimentary Ethereum test in Typescript that you can begin to modify on your own.


Step One: Set Up Prerequisites (2 minutes)
------------------------------------------
### Install Docker
Verify that you have the Docker daemon installed and running on your local machine by running (you can copy this code by hovering over it and clicking the clipboard in the top-right corner):

```
docker image ls
```

* If you don't have Docker installed, do so by following [the installation instructions](https://docs.docker.com/get-docker/)
* If Docker is installed but not running, start it

**NOTE:** [DockerHub restricts downloads from users who aren't logged in](https://www.docker.com/blog/what-you-need-to-know-about-upcoming-docker-hub-rate-limiting/) to 100 images downloaded per 6 hours, so if at any point in this tutorial you see the following error message:

```
Error response from daemon: toomanyrequests: You have reached your pull rate limit. You may increase the limit by authenticating and upgrading: https://www.docker.com/increase-rate-limit
```

you can fix it by creating a DockerHub account (if you don't have one already) and registering it with your local Docker engine like so:

```
docker login
```

### Install the Kurtosis CLI
Follow the steps [on this installation page][installation] to install the CLI, or upgrade it to latest if it's already installed.

Step Two: Start A Sandbox Enclave (3 minutes)
---------------------------------------------
The Kurtosis engine provides you isolated environments called "enclaves" to run your services inside. Let's use the CLI to start a sandbox enclave:

```
kurtosis sandbox
```

The Kurtosis images that run the engine will take a few seconds to pull the first time, but once done you'll have a Javascript REPL with tab-complete attached to your enclave.

All interaction with a Kurtosis enclave is done via [a client library][core-documentation], whose entrypoint is the `EnclaveContext` object - a representation of the network running inside the enclave. The magically-populated `enclaveCtx` variable in your REPL is how you'll interact with your enclave.

Let's check the contents of our enclave (this entire block can be copy-pasted as-is into the REPL):

```javascript
getServicesResult = await enclaveCtx.getServices()
services = getServicesResult.value
```

We haven't started any services yet, so the enclave will be empty. Note how we called `await` on `enclaveCtx.getServices()`. This is because every `enclaveCtx` call is asynchronous and returns a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise); `await`ing blocks until that value is available.


Step Three: Start An Ethereum Network (5 minutes)
-------------------------------------------------
Now that we have an enclave, let's put something in it! Ethereum is one of the most popular blockchains in the world, so let's get a private Ethereum network running:

```javascript
loadEthModuleResult = await enclaveCtx.loadModule("eth-module", "kurtosistech/ethereum-kurtosis-module", "{}")
ethModuleCtx = loadEthModuleResult.value
executeEthModuleResult = await ethModuleCtx.execute("{}")
executeEthModuleResultObj = JSON.parse(executeEthModuleResult.value)
console.log(executeEthModuleResultObj)
```

This will take approximately a minute to run, with the majority of the time spent pulling the Ethereum images. After the final `console.log` line executes, you'll see a result with information about the services running inside your enclave:

```javascript
{
  bootnode_service_id: 'bootnode',
  node_info: {
    bootnode: {
      ip_addr_inside_network: '154.18.224.5',
      ip_addr_on_host_machine: '127.0.0.1',
      rpc_port_id: 'rpc',
      ws_port_id: 'ws',
      tcp_discovery_port_id: 'tcp-discovery',
      udp_discovery_port_id: 'udp-discovery'
    },
    'ethereum-node-1': {
      ip_addr_inside_network: '154.18.224.7',
      ip_addr_on_host_machine: '127.0.0.1',
      rpc_port_id: 'rpc',
      ws_port_id: 'ws',
      tcp_discovery_port_id: 'tcp-discovery',
      udp_discovery_port_id: 'udp-discovery'
    },
    'ethereum-node-2': {
      ip_addr_inside_network: '154.18.224.9',
      ip_addr_on_host_machine: '127.0.0.1',
      rpc_port_id: 'rpc',
      ws_port_id: 'ws',
      tcp_discovery_port_id: 'tcp-discovery',
      udp_discovery_port_id: 'udp-discovery'
    }
  },
  signer_keystore_content: '{"address":"14f6136b48b74b147926c9f24323d16c1e54a026","crypto":{"cipher":"aes-128-ctr","ciphertext":"39fb1d86c1082c0103ece1c5f394321f127bf1b65e6c471edcfb181058a3053a","cipherparams":{"iv":"c366d1eed33e8693fec7a85fad65d19f"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"f210bc3b55117197f62a7ab8d85f2172342085f1daafa31034016163b8bc7db6"},"mac":"2ff8aa24d9b73ccfdb99cfd15fcdbcc8f640aaa7861e6813d53efaf550725fac"},"id":"6c5ac271-d24a-4971-b365-49490cc4befc","version":3}',
  signer_account_password: 'passphrase'
}
```

And if we query the enclave's services again...

```javascript
getServicesResult = await enclaveCtx.getServices()
services = getServicesResult.value
```

...we see three service IDs:

```javascript
Set(3) { 'bootnode', 'ethereum-node-1', 'ethereum-node-2' }
```

But what just happened?

Starting networks is a very common task in Kurtosis, so we provide [a framework called "modules"](https://docs.kurtosistech.com/modules.html) for making it dead simple. An executable module is basically a chunk of code that responds to an "execute" command, packaged as a Docker image, that runs inside a Kurtosis enclave - sort of like Docker Compose on steroids. In the steps above, we called `enclaveCtx.loadModule` to load [the Ethereum module](https://github.com/kurtosis-tech/ethereum-kurtosis-module) into the enclave with module ID `eth-module`, and `ethModuleCtx.execute` to run it. The Ethereum module doesn't take any parameters at load or execute time (hence the `{}`), but other modules do.

Now that you have a pet Ethereum network, let's do something with it.

Step Four: Talk To Ethereum (5 minutes)
---------------------------------------
Talking to Ethereum in Javascript is easily accomplished with [the EthersJS library](https://docs.ethers.io/v5/). Your Javascript REPL is running in a Docker image (so that you don't need Javascript installed locally), so we'll need to install EthersJS on that image.

First, in a new terminal window, run the following to find the enclave our REPL is running inside:

```
kurtosis enclave ls
```

You should see an output similar (but not identical) to the following:

```
EnclaveID                   Status     
KT2021-10-17T15.46.23.438   EnclaveContainersStatus_RUNNING
```

Copy the enclave ID, and slot it into `YOUR_ENCLAVE_ID_HERE` in the below command:

```
kurtosis enclave inspect YOUR_ENCLAVE_ID_HERE
```

Kurtosis will output everything it knows about your enclave, similar but not identical to the output below:

```
Enclave ID:                KT2021-12-08T13.22.31.404
Enclave Status:            EnclaveContainersStatus_RUNNING
API Container Status:      EnclaveAPIContainerStatus_RUNNING
API Container Host Port:   127.0.0.1:50828

======================================== Interactive REPLs ========================================
GUID
1638991359

========================================== User Services ==========================================
GUID                         LocalPortBindings
bootnode_1638991384          8545/tcp -> 127.0.0.1:50845
                             8546/tcp -> 127.0.0.1:50846
                             30303/tcp -> 127.0.0.1:50847
                             30303/udp -> 127.0.0.1:60255
ethereum-node-1_1638991388   8545/tcp -> 127.0.0.1:50853
                             8546/tcp -> 127.0.0.1:50851
                             30303/tcp -> 127.0.0.1:50852
                             30303/udp -> 127.0.0.1:57127
ethereum-node-2_1638991392   8545/tcp -> 127.0.0.1:50948
                             8546/tcp -> 127.0.0.1:50949
                             30303/tcp -> 127.0.0.1:50947
                             30303/udp -> 127.0.0.1:49317

========================================= Kurtosis Modules =========================================
GUID                    LocalPortBindings
eth-module_1638991377   1111/tcp -> 127.0.0.1:50841
```

Copy the interactive REPL's GUID, and replace both `YOUR_ENCLAVE_ID_HERE` and `YOUR_REPL_GUID_HERE` in the below command with the appropriate values:

```
kurtosis repl install YOUR_ENCLAVE_ID_HERE YOUR_REPL_GUID_HERE ethers
```

When the command finishes, you can now use it in your CLI! (You can execute the next command in the interactive REPL that should be still open in the previous tab)

```javascript
const ethers = require("ethers")
```

Now let's get a connection to the node with service ID `bootnode` by getting a [JsonRpcProvider](https://docs.ethers.io/v5/api/providers/jsonrpc-provider/). First we'll grab the `ServiceContext` (Kurtosis' representation of the node), use it to get the ports that the bootnode is listening on, and then construct an Ethers client connected to the node:

```javascript
// Grab the bootnode's service context
bootnodeServiceId = executeEthModuleResultObj.bootnode_service_id
bootnodeNodeObj = executeEthModuleResultObj.node_info[bootnodeServiceId]
getBootnodeServiceCtxResult = await enclaveCtx.getServiceContext(bootnodeServiceId)
bootnodeServiceCtx = getBootnodeServiceCtxResult.value

// Get the IP & port of the bootnode, inside the enclave
bootnodeRpcPortId = bootnodeNodeObj.rpc_port_id
bootnodeRpcPort = bootnodeServiceCtx.getPrivatePorts().get(bootnodeRpcPortId)
bootnodePrivateIp = bootnodeServiceCtx.getPrivateIPAddress()

// Instantiate the Ethers client
bootnodeRpcProvider = new ethers.providers.JsonRpcProvider(`http://${bootnodePrivateIp}:${bootnodeRpcPort.number}`)
```

Notice how we used the `executeEthModuleResultObj` object containing details about the Ethereum network, which we got from executing the module at the very beginning.

Finally, let's verify that our Ethereum network is producing blocks:

```javascript
blockNumber = await bootnodeRpcProvider.getBlockNumber()
if (blockNumber > 0) { console.log("All is well!"); }
```

And that's it! Anything doable in Ethers is now doable against your private Ethereum network running in your Kurtosis enclave.

To exit out of the REPL you can enter any of:

* Ctrl-D
* Ctrl-C, twice
* `.exit`

Your enclave will stay running until you command it to stop. Run the following to tell Kurtosis to tear down your enclave, replacing `YOUR_ENCLAVE_ID` with the ID of your running enclave:

```
kurtosis enclave rm -f YOUR_ENCLAVE_ID
```

Step Five: Get An Ethereum Testsuite (5 minutes)
---------------------------------------------
Manually verifying against a sandbox network is nice, but it'd be great if we could take our logic and run it as part of CI. Fortunately, Kurtosis allows us to do exactly that. 

Normally, you'd have a project that you'd add the Kurtosis tests to. For the purposes of this onboarding though, we've created a sample Typescript project with testing ready to go. Go ahead and clone it from [here](https://github.com/kurtosis-tech/onboarding-ethereum-testsuite), and we'll take a look around.

The first thing to notice is the `test/basic_eth_test.ts` file. This contains a Mocha test that connects to the Kurtosis engine, spins up an enclave for the test, does nothing (right now), and stops it when it's done. 

The second thing to notice is the `kurtosis-engine-api-lib` dev dependency in the `package.json`. This is the client library for connecting to the Kurtosis engine for creating, manipulating, stopping, & destroying enclaves.

Now let's see the testing framework in action. From the root of the repo, run:

```
scripts/build.sh
```

The testsuite will run, and you'll see that our basic test passed!

If we go ahead and run the enclave-listing command again:

```
kurtosis enclave ls
```

you'll notice a new stopped `basic-ethereum-test_XXXXXXXXXXXXX` enclave. Our current test is set to stop enclaves after it's done with them so debugging information stays around, though the test could easily be switched to destroy the enclave instead.

Step Six: Test Ethereum (5 minutes)
-----------------------------------
We now have a test running in the testing framework, but our test doesn't currently do anything. Let's fix that.

First, inside the test, replace the `// TODO Replace with Ethereum network setup` line with the following code:

```typescript
log.info("Setting up Ethereum network...")
const loadEthModuleResult: Result<ModuleContext, Error> = await enclaveCtx.loadModule(ETH_MODULE_ID, ETH_MODULE_IMAGE, "{}");
if (loadEthModuleResult.isErr()) {
    throw loadEthModuleResult.error;
}
const ethModuleCtx: ModuleContext = loadEthModuleResult.value;

const executeEthModuleResult: Result<string, Error> = await ethModuleCtx.execute("{}")
if (executeEthModuleResult.isErr()) {
    throw executeEthModuleResult.error;
}
const executeEthModuleResultObj = JSON.parse(executeEthModuleResult.value);
log.info("Ethereum network set up successfully");
```

This is the same code we already executed in the REPL, cleaned up for Typescript. The only new bits to pay attention to are the error-checking: all `EnclaveContext` and `ModuleContext`  methods return [a Result object][neverthrow] (much like in Rust). If the is a non-`Ok` result, we'll throw the result which will cause Mocha to mark the test as failed.

Second, replace the `// TODO Replace with block number check` line with this code:

```typescript
log.info("Verifying block number is increasing...");
// Grab the bootnode's service context
const bootnodeServiceId = executeEthModuleResultObj.bootnode_service_id
const bootnodeNodeObj = executeEthModuleResultObj.node_info[bootnodeServiceId]
const getBootnodeServiceCtxResult = await enclaveCtx.getServiceContext(bootnodeServiceId)
if (getBootnodeServiceCtxResult.isErr()) {
    throw getBootnodeServiceCtxResult.error;
}
const bootnodeServiceCtx = getBootnodeServiceCtxResult.value;

// Get the IP & port of the bootnode, *outside* the enclave
const bootnodeRpcPortId = bootnodeNodeObj.rpc_port_id
const bootnodeRpcPort = bootnodeServiceCtx.getPublicPorts().get(bootnodeRpcPortId)
const bootnodePublicIp = bootnodeServiceCtx.getMaybePublicIPAddress()

if (bootnodeRpcPort === undefined) {
    throw new Error("We expected having set a RPC port in the bootnode, but it is undefined");
}

// Instantiate the Ethers client
const bootnodeRpcProvider = new ethers.providers.JsonRpcProvider(`http://${bootnodePublicIp}:${bootnodeRpcPort.number}`)
const blockNumber: number = await bootnodeRpcProvider.getBlockNumber();
if (blockNumber === 0) {
    throw new Error("We expected the Ethereum cluster to be producing blocks, but the block number is still 0");
}
log.info("Verified that block number is increasing");
```

This code is nearly the same as what we ran in the sandbox, but with one crucial difference: we're using the _public_ IP address and port now because our test is running outside the enclave. Code that runs inside an enclave (e.g. the sandbox REPL or a Kurtosis module) should use the private IP & ports, while code that runs outside should use the public IP & ports.

Finally, build and run the testsuite again:

```
scripts/build.sh
```

The test will pass, indicating that our test set up an Ethereum network and ran our block count verification logic against it!

<!-- explain static files, and show how they could be used for ETH genesis -->
<!-- TODO Link to docs and further deepdives -->

<!-- TODO explain extra flags to control testsuite execution -->
<!-- TODO explain executing the testsuite in CI -->
<!-- TODO explain Debug mode, host port bindings, and setting debug log level -->

[installation]: https://docs.kurtosistech.com/installation.html
[neverthrow]: https://www.npmjs.com/package/neverthrow
[core-documentation]: https://docs.kurtosistech.com/kurtosis-core/lib-documentation
