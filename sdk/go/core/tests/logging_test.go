package test

import (
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoggingTestSuite struct {
	suite.Suite
}

func (suite *LoggingTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *LoggingTestSuite) TestLogging() {
	ydk.YLogInfo("Testing YLogInfo")
	ydk.YLogDebug("Testing YLogDebug")
	ydk.YLogWarn("Testing YLogDebug")
	ydk.YLogError("Testing YLogDebug")
}

func TestLoggingTestSuite(t *testing.T) {
	if testing.Verbose() {
		ydk.EnableLogging(ydk.Debug)
	}
	suite.Run(t, new(LoggingTestSuite))
}
