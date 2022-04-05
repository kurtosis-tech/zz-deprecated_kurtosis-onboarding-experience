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

Step Two: Start An Enclave And Run One User Service Inside  (3 minutes)
----------------------------------------------------------------------
The Kurtosis engine provides you isolated environments called "enclaves" to run your services inside. Let's use the CLI to create a new enclave:

```
mkdir /tmp/my-enclave
cd /tmp/my-enclave
kurtosis enclave new --id demo
```

The Kurtosis images that run the engine will take a few seconds to pull the first time, but once done you'll have a new enclave ready to be used.

If we want to see if the enclave has been created we can do this:
```
kurtosis enclave ls
```

You should see the enclave list, listing the new one:

```
EnclaveID   Status
demo        EnclaveContainersStatus_RUNNING
```

It means that you were able to create your first enclave, well done! Also, we can check the contents of our enclave.

```
kurtosis enclave inspect demo
```

You should see something like this:

```
Enclave ID:                           kt2022-04-04t17.17.03.135
Data Directory:                       /Users/username/Library/Application Support/kurtosis/engine-data/enclaves/kt2022-04-04t17.17.03.135
Enclave Status:                       EnclaveContainersStatus_RUNNING
API Container Status:                 EnclaveAPIContainerStatus_RUNNING
API Container Host GRPC Port:         127.0.0.1:55783
API Container Host GRPC Proxy Port:   127.0.0.1:55784

========================================= Kurtosis Modules =========================================
GUID   LocalPortBindings

========================================== User Services ==========================================
GUID   ID   LocalPortBindings
```

If you pain attention to the screen you will se that we have the enclave but no user service or Kurtosis module has been created yet

Now that we have an enclave, let's put something in it!, go ahead and execute the following command to run a `webserver` (based on the httpd Docker image) inside the enclave:

```
kurtosis service add demo webserver httpd --ports http=80
```

And the screen will show us something like this:

```
Service ID: webserver
Ports Bindings:
   http:   80/tcp -> 127.0.0.1:63825
```

Now, you can go to your browser and check that the webserver is up and running through the local url provided in the http port binding (e.g.: 127.0.0.1:63825)

Let's inspect the enclave again:

```
kurtosis enclave inspect demo
```

Now, you should see the webserver listed in the user service slot, like this example:

```
========================================== User Services ==========================================
GUID                   ID          LocalPortBindings
webserver-1649186256   webserver   80/tcp -> 127.0.0.1:63825
```

Finally, if you don't want to continue working with the service you can ger rid of it doing the following:

```
kurtosis service rm demo webserver
```


Step Three: Start An Ethereum Network (5 minutes)
-------------------------------------------------
Now that we know how to create and enclave and add a service, we can go one step ahead and start enjoying one of the most important features Kurtosis have, the Kurtosis modules.

Ethereum is one of the most popular blockchains in the world, so let's get a private Ethereum network running just to executing a Kurtosis module:

```
kurtosis module exec 'kurtosistech/ethereum-kurtosis-module'
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

And if we inspect the enclave we will see the Kurtosis module and the three service IDs which the respective local port bindings:

```
========================================= Kurtosis Modules =========================================
GUID                                             LocalPortBindings
ethereum-kurtosis-module.1649105153-1649105156   1111/tcp -> 127.0.0.1:55897

========================================== User Services ==========================================
GUID                         ID                LocalPortBindings
bootnode-1649105160          bootnode          30303/tcp -> 127.0.0.1:55902
                                               30303/udp -> 127.0.0.1:49555
                                               8545/tcp -> 127.0.0.1:55903
                                               8546/tcp -> 127.0.0.1:55901
ethereum-node-1-1649105164   ethereum-node-1   8545/tcp -> 127.0.0.1:55907
                                               8546/tcp -> 127.0.0.1:55908
                                               30303/tcp -> 127.0.0.1:55909
                                               30303/udp -> 127.0.0.1:56275
ethereum-node-2-1649105167   ethereum-node-2   30303/udp -> 127.0.0.1:51099
                                               8545/tcp -> 127.0.0.1:55964
                                               8546/tcp -> 127.0.0.1:55965
                                               30303/tcp -> 127.0.0.1:55963

```

But what just happened?

Starting networks is a very common task in Kurtosis, so we provide [a framework called "modules"](https://docs.kurtosistech.com/modules.html) for making it dead simple. An executable module is basically a chunk of code that responds to an "execute" command, packaged as a Docker image, that runs inside a Kurtosis enclave - sort of like Docker Compose on steroids. In the steps above, we executed `kurtosis module exec` to load [the Ethereum module](https://github.com/kurtosis-tech/ethereum-kurtosis-module) into a new enclave. The Ethereum module doesn't take any parameters at load or execute time, but other modules do so, you can add them using the `--execute-params` flag (e.g: --execute-params '{"numDatastores":2}')

Now that you have a pet Ethereum network, let's do something with it.

Step Four: Talk To Ethereum (5 minutes)
---------------------------------------
Talking to Ethereum is easily accomplished with cURL commands.

First, run the following to find the enclave our Eth network is running inside:

```
kurtosis enclave ls
```

You should see an output similar (but not identical) to the following:

```
EnclaveID                             Status
demo                                  EnclaveContainersStatus_RUNNING
ethereum-kurtosis-module.1649186380   EnclaveContainersStatus_RUNNING

```

Copy the enclave ID, where the Ethereum network is running, and slot it into `YOUR_ENCLAVE_ID_HERE` in the below command:

```
kurtosis enclave inspect YOUR_ENCLAVE_ID_HERE
```

Kurtosis will output everything it knows about your enclave, similar but not identical to the output below:

```
========================================= Kurtosis Modules =========================================
GUID                                             LocalPortBindings
ethereum-kurtosis-module.1649105153-1649105156   1111/tcp -> 127.0.0.1:55897

========================================== User Services ==========================================
GUID                         ID                LocalPortBindings
bootnode-1649105160          bootnode          30303/tcp -> 127.0.0.1:55902
                                               30303/udp -> 127.0.0.1:49555
                                               8545/tcp -> 127.0.0.1:55903
                                               8546/tcp -> 127.0.0.1:55901
ethereum-node-1-1649105164   ethereum-node-1   8545/tcp -> 127.0.0.1:55907
                                               8546/tcp -> 127.0.0.1:55908
                                               30303/tcp -> 127.0.0.1:55909
                                               30303/udp -> 127.0.0.1:56275
ethereum-node-2-1649105167   ethereum-node-2   30303/udp -> 127.0.0.1:51099
                                               8545/tcp -> 127.0.0.1:55964
                                               8546/tcp -> 127.0.0.1:55965
                                               30303/tcp -> 127.0.0.1:55963
```
Now let's connect to the http address exposed in the bootnode to retrieve Ethereum network information.

The ethereum Kurtosis module has set the 8545 private port as the http port where users can execute cURL request using the local port binding

Let's verify that our Ethereum network is producing blocks, copy the local port binging, where the http server is running, and slot it into `HTTP_URL_HERE` in the below command:

```
curl -X POST HTTP_URL_HERE --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' -H "Content-Type: application/json"
```

And you will see something like this, where the result is a hexadecimal number representing the current block number:

```
{"jsonrpc":"2.0","id":83,"result":"0x93"}
```

And that's it! Anything doable in Ethers is now doable against your private Ethereum network running in your Kurtosis enclave.

To clean up all the recent elements created you can execute the following:

```
kurtosis clean -a
```

and Kurtosis will tear down all enclaves and everything inside.

Step Five: Get An Ethereum Testsuite (5 minutes)
---------------------------------------------
Manually verifying against an enclave is nice, but it'd be great if we could take our logic and run it as part of CI. Kurtosis has a testing framework that allows us to do exactly that.

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

Pay attention to the error-checking: all `NetworkContext` methods, as well as the `Test.setup` and `Test.run` methods, return [a Result object][neverthrow] (much like in Rust). If `setup` or `run` return a non-`Ok` result, the test will be marked as failed. This allows for easy, consistent error-checking: simply propagate the error upwards.

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
if (bootnodeRpcPort === undefined) {
    throw new Error("We expected the boot node to have a public RPC port, but it was undefined");
}
const bootnodePublicIp = bootnodeServiceCtx.getMaybePublicIPAddress()


// Instantiate the Ethers client
const bootnodeRpcProvider = new ethers.providers.JsonRpcProvider(`http://${bootnodePublicIp}:${bootnodeRpcPort.number}`)
const blockNumber: number = await bootnodeRpcProvider.getBlockNumber();
if (blockNumber === 0) {
    throw new Error("We expected the Ethereum cluster to be producing blocks, but the block number is still 0");
}
log.info("Verified that block number is increasing");
```

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
