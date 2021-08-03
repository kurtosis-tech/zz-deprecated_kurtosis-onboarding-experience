package testsuite_impl

import (
	"github.com/kurtosis-tech/kurtosis-client/golang/lib/services"
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

func (suite MyTestsuite) GetNetworkWidthBits() uint32 {
	return 8
}

func (suite MyTestsuite) GetStaticFiles() map[services.StaticFileID]string {
	return map[services.StaticFileID]string{}
}

