package test

import (
	"fmt"
	ysanity "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/sanity"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	_ "github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SanityTypesTestSuite struct {
	suite.Suite
	Provider providers.NetconfServiceProvider
	CRUD     services.CrudService
}

func (suite *SanityTypesTestSuite) SetupSuite() {
	suite.CRUD = services.CrudService{}
	suite.Provider = providers.NetconfServiceProvider{
		Address:  "127.0.0.1",
		Username: "admin",
		Password: "admin",
		Port:     12022}
	suite.Provider.Connect()
}

func (suite *SanityTypesTestSuite) BeforeTest(suiteName, testName string) {
	suite.CRUD.Delete(&suite.Provider, &ysanity.Runner{})
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *SanityTypesTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *SanityTypesTestSuite) TestStatusEnum() {
	runner := ysanity.Runner{}
	runner.Ytypes.BuiltInT.Status = ysanity.Runner_Ytypes_BuiltInT_Status_not_connected
	suite.CRUD.Create(&suite.Provider, &runner)
	entity := suite.CRUD.Read(&suite.Provider, &ysanity.Runner{})
	runner_read := entity.(*ysanity.Runner)
	suite.Equal(runner_read.Ytypes.BuiltInT.Status, "not connected")
}

func (suite *SanityTypesTestSuite) TestEmbeddedEnum() {
	runner := ysanity.Runner{}
	runner.Ytypes.BuiltInT.EmbededEnum = ysanity.Runner_Ytypes_BuiltInT_EmbededEnum_zero
	suite.CRUD.Create(&suite.Provider, &runner)
	entity := suite.CRUD.Read(&suite.Provider, &ysanity.Runner{})
	runner_read := entity.(*ysanity.Runner)
	suite.Equal(runner_read.Ytypes.BuiltInT.EmbededEnum, "zero")
}

func (suite *SanityTypesTestSuite) TestEnum() {
	runner := ysanity.Runner{}
	runner.Ytypes.BuiltInT.EnumValue = ysanity.YdkEnumTest_not_set
	suite.CRUD.Create(&suite.Provider, &runner)
	entity := suite.CRUD.Read(&suite.Provider, &ysanity.Runner{})
	runner_read := entity.(*ysanity.Runner)
	suite.Equal(runner_read.Ytypes.BuiltInT.EnumValue, "not-set")
}

func (suite *SanityTypesTestSuite) TestBitsPos1() {
	runner := ysanity.Runner{}
	runner.Ytypes.BuiltInT.BitsValue = map[string]bool{"disable-nagle": true, "auto-sense-speed": true}
	suite.CRUD.Create(&suite.Provider, &runner)
	entity := suite.CRUD.Read(&suite.Provider, &ysanity.Runner{})
	runner_read := entity.(*ysanity.Runner)
	suite.Equal(runner_read.Ytypes.BuiltInT.BitsValue, "disable-nagle auto-sense-speed")
}

func (suite *SanityTypesTestSuite) TestBitsPos2() {
	runner := ysanity.Runner{}
	runner.Ytypes.BuiltInT.BitsValue = map[string]bool{"disable-nagle": true, "auto-sense-speed": false}
	suite.CRUD.Create(&suite.Provider, &runner)
	entity := suite.CRUD.Read(&suite.Provider, &ysanity.Runner{})
	runner_read := entity.(*ysanity.Runner)
	suite.Equal(runner_read.Ytypes.BuiltInT.BitsValue, "disable-nagle")
}

func TestSanityTypesTestSuite(t *testing.T) {
	suite.Run(t, new(SanityTypesTestSuite))
}
