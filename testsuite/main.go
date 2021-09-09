package main

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-onboarding-experience/testsuite/execution_impl"
	"github.com/kurtosis-tech/kurtosis-testsuite-api-lib/golang/lib/execution"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	successExitCode = 0
	failureExitCode = 1
)

func main() {

	configurator := execution_impl.EthereumTestsuiteConfigurator{}

	suiteExecutor := execution.NewTestSuiteExecutor(configurator)
	if err := suiteExecutor.Run(); err != nil {
		logrus.Errorf("An error occurred running the test suite executor:")
		fmt.Fprintln(logrus.StandardLogger().Out, err)
		os.Exit(failureExitCode)
	}
	os.Exit(successExitCode)
}
