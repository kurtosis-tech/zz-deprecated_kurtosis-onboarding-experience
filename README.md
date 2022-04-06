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
kurtosis enclave new --id demo
```

The Kurtosis images that run the engine will take a few seconds to pull the first time, but once done you'll have a new enclave ready to be used.

To see the enclave you just created, run:
```
kurtosis enclave ls
```

You should see the enclave list, showing the new enclave:

```
EnclaveID   Status
demo        EnclaveContainersStatus_RUNNING
```

You've created your first enclave! To check the contents of the enclave:

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

Notice how the enclave doesn't yet have any user services or Kurtosis modules. This is because enclaves are created empty by default.

Now that we have an enclave, let's put something in it. Execute the following command to run a web server (using the `httpd` Docker image) inside the enclave and give it the ID `webserver`:

```
kurtosis service add demo webserver httpd --ports http=80
```

Your output should look something like this:

```
Service ID: webserver
Ports Bindings:
   http:   80/tcp -> 127.0.0.1:63825
```

Now, you can go to your browser and check that the web server is up and running through the local URL provided in the `http` port binding (e.g. `127.0.0.1:63825`). Note that the port number is ephemeral and yours will be different.

Let's inspect the enclave again:

```
kurtosis enclave inspect demo
```

Now, you should see the web server listed in the "User Services" section, like so:

```
========================================== User Services ==========================================
GUID                   ID          LocalPortBindings
webserver-1649186256   webserver   80/tcp -> 127.0.0.1:63825
```

Finally, you can remove the service with the following:

```
kurtosis service rm demo webserver
```


Step Three: Start An Ethereum Network (5 minutes)
-------------------------------------------------
Now that we know how to create an enclave and add a service, we can proceed to one of the most important features in Kurtosis: modules.

Ethereum is the most popular smart contract blockchain in the world, so let's create a private Ethereum network in Kurtosis:

```
kurtosis module exec --enclave-id demo 'kurtosistech/ethereum-kurtosis-module'
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

And if we inspect the enclave again....

```
kurtosis enclave inspect demo
```

...we'll see that our enclave now has the Kurtosis module and three Ethereum nodes with local port bindings:

```
========================================= Kurtosis Modules =========================================
GUID                                             LocalPortBindings
ethereum-kurtosis-module.1649269281-1649269281   1111/tcp -> 127.0.0.1:55863

========================================== User Services ==========================================
GUID                         ID                LocalPortBindings
bootnode-1649269284          bootnode          8546/tcp -> 127.0.0.1:55879
                                               30303/tcp -> 127.0.0.1:55877
                                               30303/udp -> 127.0.0.1:59562
                                               8545/tcp -> 127.0.0.1:55878
ethereum-node-1-1649269292   ethereum-node-1   30303/udp -> 127.0.0.1:55356
                                               8545/tcp -> 127.0.0.1:55883
                                               8546/tcp -> 127.0.0.1:55884
                                               30303/tcp -> 127.0.0.1:55885
ethereum-node-2-1649269296   ethereum-node-2   30303/udp -> 127.0.0.1:64347
                                               8545/tcp -> 127.0.0.1:56013
                                               8546/tcp -> 127.0.0.1:56014
                                               30303/tcp -> 127.0.0.1:56012
webserver-1649268254         webserver         <none>

```

But what just happened?

Starting networks is a very common task in Kurtosis, so we provide [a framework called "modules"](https://docs.kurtosistech.com/modules.html) for making it dead simple. An executable module is basically a chunk of code that responds to an "execute" command, packaged as a Docker image, that runs inside a Kurtosis enclave - sort of like Docker Compose on steroids. In the steps above, we executed `kurtosis module exec` to load [the Ethereum module](https://github.com/kurtosis-tech/ethereum-kurtosis-module) into the demo enclave. The Ethereum module doesn't take any parameters at load or execute time, but other modules do via the `--execute-params` flag.

Now that you have a pet Ethereum network, let's do something with it.

Step Four: Talk To Ethereum (5 minutes)
---------------------------------------

Now let's connect to the Ethereum bootnode to verify the network is producing blocks.

Find the Ethereum node with ID `bootnode` in the enclave contents, find its RPC port declared on `8545/tcp`, and copy the public IP and port that it's bound to on your machine (e.g. `127.0.0.1:55878`).

Then, slot it into RPC_URL_HERE in the below command:

```
curl -X POST RPC_URL_HERE --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' -H "Content-Type: application/json"
```

You should get a response similar to the following, where `result` is the current block number in hexadecimal:

```
{"jsonrpc":"2.0","id":83,"result":"0x93"}
```

And that's it! Anything doable on the testnet is now doable against your private Ethereum network running in your Kurtosis enclave.

To destroy the enclave we created and everything inside it, run:

```
kurtosis enclave rm -f demo
```

To destroy all enclaves, run:

```
kurtosis clean -a
```

Step Five: Get An Ethereum Testsuite (5 minutes)
---------------------------------------------
Manually verifying against an enclave is nice, but it'd be great if we could take our logic and run it as part of CI. Fortunately, Kurtosis was designed from the beginning with testing in mind.

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

This code uses the Kurtosis SDK, which is the same SDK that the Kurtosis CLI uses to communicate with the Kurtosis engine. The only new bits to pay attention to are the error-checking: all `EnclaveContext` and `ModuleContext` methods return [a Result object][neverthrow] (much like in Rust). If `result.isErr()`, we'll throw the result which will cause Mocha to mark the test as failed.

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
