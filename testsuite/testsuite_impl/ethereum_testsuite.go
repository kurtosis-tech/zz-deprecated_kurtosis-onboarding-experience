package testsuite_impl

import (
	"github.com/kurtosis-tech/kurtosis-onboarding-experience/testsuite/testsuite_impl/basic_ethereum_test"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
)

type EthereumTestsuite struct {}

func (suite EthereumTestsuite) GetTests() map[string]testsuite.Test {
	tests := map[string]testsuite.Test{
		"BasicEthereumTest": &basic_ethereum_test.BasicEthereumTest{},
	}
	return tests
}


