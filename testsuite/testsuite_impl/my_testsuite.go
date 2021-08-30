package testsuite_impl

import (
	"github.com/kurtosis-tech/kurtosis-onboarding-experience/testsuite/testsuite_impl/my_test"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/testsuite"
)

type MyTestsuite struct {}

func (suite MyTestsuite) GetTests() map[string]testsuite.Test {
	tests := map[string]testsuite.Test{
		"myTest": &my_test.MyTest{},
	}
	return tests
}


