package test

import (
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk"
	ysanity "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/sanity"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

const testHome = "../../../cpp/core/tests/models"

type RestconfProviderTestSuite struct {
	suite.Suite
	Provider providers.RestconfServiceProvider
	CRUD     services.CrudService
}

func (suite *RestconfProviderTestSuite) SetupSuite() {
	suite.CRUD = services.CrudService{}
	suite.Provider = providers.RestconfServiceProvider{
		Path:     testHome,
		Address:  "localhost",
		Username: "admin",
		Password: "admin",
		Port:     12306}
	suite.Provider.Connect()
}

func (suite *RestconfProviderTestSuite) BeforeTest(suiteName, testName string) {
	suite.CRUD.Delete(&suite.Provider, &ysanity.Runner{})
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *RestconfProviderTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *RestconfProviderTestSuite) TestCreateDelRead() {
	runner := ysanity.Runner{}
	runner.Ytypes.BuiltInT.Number8 = 3
	suite.CRUD.Create(&suite.Provider, &runner)
	runnerRead := suite.CRUD.Read(&suite.Provider, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(&runner, runnerRead), true)
}

func TestRestconfProviderTestSuite(t *testing.T) {
	ydk.EnableLogging(ydk.Debug)
	suite.Run(t, new(RestconfProviderTestSuite))
}
