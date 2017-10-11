package test

import (
	"fmt"
	ysanity "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/sanity"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NetconfServiceTestSuite struct {
	suite.Suite
	Provider 	providers.NetconfServiceProvider
	Crud 		services.CrudService
	NS 			services.NetconfService
}

func (suite *NetconfServiceTestSuite) SetupSuite() {
	suite.Crud = services.CrudService{}
	suite.NS = services.NetconfService{}
	suite.Provider = providers.NetconfServiceProvider{
		Address:  "127.0.0.1",
		Username: "admin",
		Password: "admin",
		Port:     12022}
	suite.Provider.Connect()
}

func (suite *NetconfServiceTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *NetconfServiceTestSuite) BeforeTest(suiteName, testName string) {
	suite.Crud.Delete(&suite.Provider, &ysanity.Runner{})
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *NetconfServiceTestSuite) TestEditCommitGet() {
	runner := ysanity.Runner{}
	runner.One.Number = 1
	runner.One.Name = "runner:one:name"

	// Edit Config
	op := suite.NS.EditConfig(&suite.Provider, types.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	// Get Config
	readEntity := suite.NS.GetConfig(&suite.Provider, types.Candidate, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)

	// // Commit
	op = suite.NS.Commit(&suite.Provider, false, -1, -1, -1)
	suite.Equal(op, true)

	// Get
	readEntity = suite.NS.Get(&suite.Provider, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
}

func (suite *NetconfServiceTestSuite) TestLockUnlock() {
	op := suite.NS.Lock(&suite.Provider, types.Running)
	suite.Equal(op, true)

	op = suite.NS.Unlock(&suite.Provider, types.Running)
	suite.Equal(op, true)
}

func (suite *NetconfServiceTestSuite) TestLockUnlockFail() {
	op := suite.NS.Lock(&suite.Provider, types.Candidate)
	suite.Equal(op, true)

	funcDidPanic, panicValue := didPanic(func() { suite.NS.Unlock(&suite.Provider, types.Running) })
	suite.Equal(funcDidPanic, true)
	suite.Regexp("YGOServiceProviderError:", panicValue)
	errMsg := `<rpc-error>
    <error-type>application</error-type>
    <error-tag>operation-failed</error-tag>
    <error-severity>error</error-severity>
  </rpc-error>`
	suite.Regexp(errMsg, panicValue)
}

func (suite *NetconfServiceTestSuite) TestValidate() {
	op := suite.NS.Validate(&suite.Provider, nil, types.Candidate, "")
	suite.Equal(op, true)

	runner := ysanity.Runner{}
	runner.One.Number = 1
	runner.One.Name = "runner:one:name"
	op = suite.NS.Validate(&suite.Provider, &runner, -1, "")
	suite.Equal(op, true)
}

func (suite *NetconfServiceTestSuite) TestCommitDiscard(){
	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	op := suite.NS.EditConfig(&suite.Provider, types.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	op = suite.NS.DiscardChanges(&suite.Provider)
	suite.Equal(op, true)

	op = suite.NS.EditConfig(&suite.Provider, types.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	op = suite.NS.Commit(&suite.Provider, false, -1, -1, -1)
	suite.Equal(op, true)

	readEntity := suite.NS.Get(&suite.Provider, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
}

// skipped in python/cpp
func (suite *NetconfServiceTestSuite) TestConfirmedCommit() {
	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	op := suite.NS.EditConfig(&suite.Provider, types.Candidate, &runner, "", "", "")
	suite.Equal(op, true)

	op = suite.NS.Commit(&suite.Provider, true, 120, -1, -1)
	suite.Equal(op, true)

	// op = suite.NS.CancelCommit(&suite.Provider, -1)
	// suite.Equal(op, true)

	readEntity := suite.NS.Get(&suite.Provider, &ysanity.Runner{})
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
}

// TODO: Problems found with path functions:
//			readEntity changes with getFilter
//			EntityEqual(readEntity, &runner) returns true, even after the above
func (suite *NetconfServiceTestSuite) TestCopyConfig() {
	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"
	getFilter := ysanity.Runner{}

	var op bool
	var readEntity types.Entity

	// Modify Candidate via CopyConfig from runner
	op = suite.NS.CopyConfig(&suite.Provider, types.Candidate, -1, &runner, "")
	suite.Equal(op, true)

	readEntity = suite.NS.GetConfig(&suite.Provider, types.Candidate, &getFilter)
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
	getFilter = ysanity.Runner{}

	// Modify Candidate via CopyConfig from runner
	runner.Two.Name = fmt.Sprintf("%s_modified", runner.Two.Name)

	op = suite.NS.CopyConfig(&suite.Provider, types.Candidate, -1, &runner, "")
	suite.Equal(op, true)

	readEntity = suite.NS.GetConfig(&suite.Provider, types.Candidate, &getFilter)
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
	getFilter = ysanity.Runner{}

	// Modify Running via CopyConfig from Candidate
	op = suite.NS.CopyConfig(&suite.Provider, types.Running, types.Candidate, nil, "")
	suite.Equal(op, true)

	readEntity = suite.NS.GetConfig(&suite.Provider, types.Running, &getFilter)
	suite.Equal(types.EntityEqual(readEntity, &runner), true)
}

func (suite *NetconfServiceTestSuite) TestDeleteConfig() {
	suite.T().Skip("startup not enabled in ConfD")

	runner := ysanity.Runner{}
	runner.Two.Number = 2
	runner.Two.Name = "runner:two:name"

	op := suite.NS.CopyConfig(&suite.Provider, types.Startup, -1, &runner, "")
	suite.Equal(op, true)

	op = suite.NS.DeleteConfig(&suite.Provider, types.Startup, "")
	suite.Equal(op, true)
}

func (suite *NetconfServiceTestSuite) TestCloseSession() {
	op := suite.NS.CloseSession(&suite.Provider)
	suite.Equal(op, true)

	funcDidPanic, panicValue := didPanic(func() { suite.NS.Lock(&suite.Provider, types.Running) })
	suite.Equal(funcDidPanic, true)
	suite.Equal(panicValue, "YGOClientError: RPC error occured")

	suite.Provider.Connect()

	op = suite.NS.Lock(&suite.Provider, types.Running)
	op = suite.NS.Unlock(&suite.Provider, types.Running)
	suite.Equal(op, true)
}


func (suite *NetconfServiceTestSuite) TestKillSession() {
	funcDidPanic, panicValue := didPanic(func() { suite.NS.KillSession(&suite.Provider, 1) })
	suite.Equal(funcDidPanic, true)
	suite.Regexp("YGOServiceProviderError:", panicValue)
}

func TestNetconfServiceTestSuite(t *testing.T) {
	if testing.Verbose() {
		ydk.EnableLogging(ydk.Debug)
	}
	suite.Run(t, new(NetconfServiceTestSuite))
}


// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(panicTestFunc func()) (bool, string) {

	didPanic := false
	var message interface{}
	func() {

		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		// call the target function
		panicTestFunc()

	}()
	return didPanic, message.(string)
}