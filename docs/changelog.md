# TBD
### Changes
* Pointed the user to the official installation docs for installing the Kurtosis CLI
* Updated onboarding docs for modules instead of Lambdas

### Features
* The README now shows the user `enclave ls`, `enclave inspect`, and `repl install`

# 0.3.2
### Changes
* Update Readme file in order to use the [Ethereum Kurtosis Lambda to v0.2.4](https://github.com/kurtosis-tech/ethereum-kurtosis-lambda/blob/develop/docs/changelog.md#024)

# 0.3.1
### Fixes
* Add note about hovering + clipboard for copying codeblocks
* Add `brew update` before CLI installation
* Prompt for `docker login` to prevent anonymous image pulls

# 0.3.0
### Features
* Revamp the onboarding once again to start the user in the sandbox and transition them to testing

# 0.2.0
### Features
* Created the Kurtosis project boilerplate which includes: CI configuration, build and run script and release script
* Added testsuite execution configuration with a custom configurator `my_testsuite_configurator` and a custom testsuite `my_testsuite`
* Added instructions for implementing the basic ETH test into the `README.md` file which test GETH client execution following the GETH documentation for [setting up a GETH client in DEV mode](https://geth.ethereum.org/docs/getting-started/dev-mode)
* Added instructions for implementing the advanced ETH test into the `README.md` file which test an Ethereum private network following the GETH documentation for [setting up an Ethereum private network](https://geth.ethereum.org/docs/interface/private-network)

### Changes
* Simplified the basic test instructions and fixed some bugs
* Renamed `MyTest` to `BasicEthereumTest`

### Removals
* Removed the advanced test for now

# 0.1.0
* Init commit
