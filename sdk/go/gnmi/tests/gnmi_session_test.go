package test

import (
	"fmt"
	"strings"
	"time"
	"runtime"
	"path/filepath"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/path"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	encoding "github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GnmiSessionTestSuite struct {
	suite.Suite
	GnmiSession path.GnmiSession
	Schema      types.RootSchemaNode
}

func (suite *GnmiSessionTestSuite) SetupSuite() {
	_, callerFile, _, _ := runtime.Caller(0)
	executablePath := filepath.Dir(callerFile)
	repopath := executablePath + "/../../../cpp/core/tests/models"
	repo := types.Repository{Path: repopath}
	suite.GnmiSession = path.GnmiSession{
		Repo: repo,
		Address:  "127.0.0.1",
		Port:     50051}
	suite.GnmiSession.Connect()
	
	suite.Schema = suite.GnmiSession.GetRootSchemaNode()
}

func (suite *GnmiSessionTestSuite) TearDownSuite() {
	suite.GnmiSession.Disconnect()
}

func (suite *GnmiSessionTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *GnmiSessionTestSuite) TestGnmiRpcCaps() {
    capRpc := path.CreateRpc(suite.Schema, "ydk:gnmi-caps")
    caps := suite.GnmiSession.ExecuteRpc(capRpc)

    json := path.CodecEncode( caps, encoding.JSON, true)
    ydk.YLogDebug(fmt.Sprintf("Server capabilities:\n%v\n", json))
    suite.NotEqual(strings.Index(json, "supported-encodings"), -1)
}

func (suite *GnmiSessionTestSuite) TestGnmiRpcGetSetBgp() {
    // Configure BGP
    bgp := path.CreateRootDataNode( suite.Schema, "openconfig-bgp:bgp")
    path.CreateDataNode( bgp, "global/config/as", 65172)
    neighbor := path.CreateDataNode( bgp, "neighbors/neighbor[neighbor-address='172.16.255.2']", "")
    path.CreateDataNode( neighbor, "config/neighbor-address", "172.16.255.2")
    path.CreateDataNode( neighbor, "config/peer-as", 65172);   

    // Configure interface
    ifc  := path.CreateRootDataNode( suite.Schema, "openconfig-interfaces:interfaces")
    lo10 := path.CreateDataNode( ifc, "interface[name='Loopback10']", "")
    path.CreateDataNode( lo10, "config/name", "Loopback10")
    path.CreateDataNode( lo10, "config/description", "Test")

    // Build and execute Rpc
    setRpc := path.CreateRpc( suite.Schema, "ydk:gnmi-set")
    bgpCreatePayload := path.CodecEncode( bgp, encoding.JSON, false)
    intCreatePayload := path.CodecEncode( ifc, encoding.JSON, false)
    path.CreateDataNode( setRpc.Input, "replace[alias='bgp']/entity", bgpCreatePayload)
    path.CreateDataNode( setRpc.Input, "replace[alias='int']/entity", intCreatePayload)

    suite.GnmiSession.ExecuteRpc(setRpc);

	// Read configuration
    bgpReadDn := path.CreateRootDataNode( suite.Schema, "openconfig-bgp:bgp")
    bgpReadPayload := path.CodecEncode( bgpReadDn, encoding.JSON, false)

    intReadDn := path.CreateRootDataNode( suite.Schema, "openconfig-interfaces:interfaces")
    intReadPayload := path.CodecEncode( intReadDn, encoding.JSON, false)

    readRpc := path.CreateRpc( suite.Schema, "ydk:gnmi-get")
    path.CreateDataNode( readRpc.Input, "type", "CONFIG")

    path.CreateDataNode( readRpc.Input, "request[alias='bgp']/entity", bgpReadPayload)
    path.CreateDataNode( readRpc.Input, "request[alias='int']/entity", intReadPayload)

    readResult := suite.GnmiSession.ExecuteRpc(readRpc);
    readConfig := path.CodecEncode( readResult, encoding.JSON, false)
    suite.Equal(bgpCreatePayload+intCreatePayload, readConfig)
    
    // Delete configuration
    deleteRpc := path.CreateRpc( suite.Schema, "ydk:gnmi-set")
    path.CreateDataNode( deleteRpc.Input, "delete[alias='bgp']/entity", bgpCreatePayload)
    path.CreateDataNode( deleteRpc.Input, "delete[alias='int']/entity", intCreatePayload)
    suite.GnmiSession.ExecuteRpc(deleteRpc);
}

func configureBgp(session *path.GnmiSession, schema types.RootSchemaNode) {
    bgp := path.CreateRootDataNode( schema, "openconfig-bgp:bgp")
    path.CreateDataNode( bgp, "global/config/as", 65172)
    neighbor := path.CreateDataNode( bgp, "neighbors/neighbor[neighbor-address='172.16.255.2']", "")
    path.CreateDataNode( neighbor, "config/neighbor-address", "172.16.255.2")
    path.CreateDataNode( neighbor, "config/peer-as", 65172);   

    setRpc := path.CreateRpc( schema, "ydk:gnmi-set")
    bgpCreatePayload := path.CodecEncode( bgp, encoding.JSON, false)
    path.CreateDataNode( setRpc.Input, "replace[alias='bgp']/entity", bgpCreatePayload)

    session.ExecuteRpc(setRpc);
}

func deleteConfigBgp(session *path.GnmiSession, schema types.RootSchemaNode) {
    bgp := path.CreateRootDataNode( schema, "openconfig-bgp:bgp")
    path.CreateDataNode( bgp, "global/config/as", 65172)
    neighbor := path.CreateDataNode( bgp, "neighbors/neighbor[neighbor-address='172.16.255.2']", "")
    path.CreateDataNode( neighbor, "config/neighbor-address", "172.16.255.2")
    path.CreateDataNode( neighbor, "config/peer-as", 65172);   

    setRpc := path.CreateRpc( schema, "ydk:gnmi-set")
    bgpCreatePayload := path.CodecEncode( bgp, encoding.JSON, false)
    path.CreateDataNode( setRpc.Input, "delete[alias='bgp']/entity", bgpCreatePayload)

    session.ExecuteRpc(setRpc);	
}

func getBgpSubscription(gs *path.GnmiSession) {
    previous := ""
    for true {
        response := gs.GetLastSubscribeResponse(previous)
        if len(response) > 0 && response != previous {
            // Do anything with received response
            ydk.YLogDebug(fmt.Sprintf("%s:  Last received subscribe response:\n%s\n", time.Now().Format(time.RFC850), response))
            previous = response
        }
        if !gs.SubscribeInProgress() {
            break
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func (suite *GnmiSessionTestSuite) TestGnmiRpcSubscribeOnce() {
    // Configure BGP
    configureBgp(&suite.GnmiSession, suite.Schema)
    
    // Build subscription
    bgpReadDn := path.CreateRootDataNode( suite.Schema, "openconfig-bgp:bgp")
    bgpReadPayload := path.CodecEncode( bgpReadDn, encoding.JSON, false)

    rpc := path.CreateRpc( suite.Schema, "ydk:gnmi-subscribe")
    subscription := path.CreateDataNode( rpc.Input, "subscription", "")
    path.CreateDataNode( subscription, "mode", "ONCE")
    path.CreateDataNode( subscription, "qos", "10")
    path.CreateDataNode( subscription, "encoding", "JSON_IETF")

    bgpSubscription := path.CreateDataNode( subscription, "subscription-list[alias='bgp']", "")
    path.CreateDataNode( bgpSubscription, "entity", bgpReadPayload)
    path.CreateDataNode( bgpSubscription, "sample-interval", "20000000")

    suite.GnmiSession.ExecuteSubscribeRpc(rpc)
    getBgpSubscription(&suite.GnmiSession)
    
    // Delete configuration
    deleteConfigBgp(&suite.GnmiSession, suite.Schema)
}

func (suite *GnmiSessionTestSuite) TestGnmiRpcSubscribeStream() {
    // Configure BGP
    configureBgp(&suite.GnmiSession, suite.Schema)
    
    // Build subscription
    bgpReadDn := path.CreateRootDataNode( suite.Schema, "openconfig-bgp:bgp")
    bgpReadPayload := path.CodecEncode( bgpReadDn, encoding.JSON, false)

    rpc := path.CreateRpc( suite.Schema, "ydk:gnmi-subscribe")
    subscription := path.CreateDataNode( rpc.Input, "subscription", "")
    path.CreateDataNode( subscription, "mode", "STREAM")
    path.CreateDataNode( subscription, "qos", "10")
    path.CreateDataNode( subscription, "encoding", "JSON_IETF")

    bgpSubscription := path.CreateDataNode( subscription, "subscription-list[alias='bgp']", "")
    path.CreateDataNode( bgpSubscription, "entity", bgpReadPayload)
    path.CreateDataNode( bgpSubscription, "sample-interval", "2000000000")
    path.CreateDataNode( bgpSubscription, "heartbeat-interval", "10000000000")

    go suite.GnmiSession.ExecuteSubscribeRpc(rpc)
    time.Sleep(100 * time.Millisecond)
    getBgpSubscription(&suite.GnmiSession)
    
    // Delete configuration
    deleteConfigBgp(&suite.GnmiSession, suite.Schema)
}

func TestGnmiSessionTestSuite(t *testing.T) {
	suite.Run(t, new(GnmiSessionTestSuite))
}
