Ethereum On-Boarding Testsuite
==============================

## Implement a Single Node Ethereum Test Network

1. Create an account on https://kurtosistech.com if you don't have one yet.
2. Verify that the Docker daemon is running on your local machine.
3. Clone this repository by running `git clone https://github.com/kurtosis-tech/kurtosis-onboarding-experience.git`
4. Run the empty Ethereum single node test `my_test` to verify that everything works on your local machine.
   1. Run `bash scripts/build-and-run.sh all`
   2. Verify that the output of the build-and-run.sh script indicates that one test ran (my_test) and that it passed.
5. Import the dependencies that are used in this example test suite.
   1. Run `go get github.com/ethereum/go-ethereum`
   2. Run `go get github.com/palantir/stacktrace`
   3. Run `go get github.com/sirupsen/logrus`
6. Write the Setup() method on the Ethereum single node test in order to bootstrap the network and leave it running and ready to use it.
   1. In your preferred IDE, open the Ethereum single node test `my_test` at `testsuite/testsuite_impl/my_test/my_test.go`
   2. Implement the private helper function `getEthereumServiceConfigurations` used to set the Ethereum node service inside the testsuite
      1. Add the following private helper functions `getContainerCreationConfig()`, `getRunConfigFunc()` and `getEthereumServiceConfigurations()` in [this Gist](https://gist.github.com/leoporoli/123cb1d682d74dcafe7920f01809b167) to the bottom of the file, so they can be used later. 
   3. Implement the first part of the Setup() method which starts the Ethereum single node service in the testsuite
      1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/d81577dcfd5fdc6605ccdf9f61eed81f) in the top of the Setup() method
   4. Implement the second and last part of the Setup() method to check if the service is available
      1. Add the following private helper functions `getEnodeAddress()` and `sendRpcCall` in [this Gist](https://gist.github.com/leoporoli/1a03e9500a61a20d06ed8e3827d72f5e) to the bottom of the file, so they can be used later.
      2. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/f9aacad32b2800a98f68bcf5fa32165c) in the bottom of the Setup() method
   5. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
7. Write the Run() method on the Ethereum single node test in order to test the initilization of an Ethereum Single Node Network in Dev Mode
   1. Read [the official Geth documentation for start Ethereum in Dev mode](https://geth.ethereum.org/docs/getting-started/dev-mode) to have a context as this test implements the commands that are described in this document
   3. Implement the first part of the Run() method which get the Ethereum single node service from the network, get the Ethereum client from the service and uses this to get the chain ID of the private Ethereum Network
      1. Add the private helper function `getEthClient` in [this Gist](https://gist.github.com/leoporoli/8a1641f9d78752f984ed672895c1f97c) to the bottom of the file, so it can be used later.
      2. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/628c78452e7e05549809f5cbcb62cdf6) in the top of the Run() method
      3. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
   4. Implement the second and last part of the Run() method which implements all the commands of the official Geth documentation in the `Dev mode` section
      1. 2. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/5a1539d2e8e45d3658a5b2398d9f3ba7) in the bottom of the Run() method
   5. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that the test contained business logic for an Ethereum single node network test, and that it passed.