package test

import (
	"fmt"
	"runtime"
	"time"
    "path/filepath"
    "strings"
	ysanity_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp"
	ysanity_int "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_interfaces"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/CiscoDevNet/ydk-go/ydk/types/yfilter"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GnmiServiceTestSuite struct {
	suite.Suite
	Service  services.GnmiService
	Provider providers.GnmiServiceProvider
//	Schema   types.RootSchemaNode
//	CodecPro providers.CodecServiceProvider
//	Codec    services.CodecService
}

func (suite *GnmiServiceTestSuite) SetupSuite() {
	_, callerFile, _, _ := runtime.Caller(0)
	executablePath := filepath.Dir(callerFile)
	repopath := executablePath + "/../../../cpp/core/tests/models"
	repo := types.Repository{Path: repopath}
	suite.Provider = providers.GnmiServiceProvider{
		Repo: repo,
		Address:  "127.0.0.1",
		Port:     50051}
	suite.Provider.Connect()
	
//	suite.Codec = services.CodecService{}
//	suite.CodecPro = providers.CodecServiceProvider{}
//    suite.CodecPro.Encoding = encoding.JSON
}

func (suite *GnmiServiceTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *GnmiServiceTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *GnmiServiceTestSuite) TestGetCapabilities() {
	caps := suite.Service.Capabilities(&suite.Provider)
	suite.NotEqual(strings.Index(caps, "supported-encodings"), -1)
}

func (suite *GnmiServiceTestSuite) TestGetSetDeleteMultiple() {
	// Build interface configuration
    ifc := ysanity_int.Interfaces_Interface{}
    ifc.Name = "Loopback10"
    ifc.Config.Name = "Loopback10"
    ifcs := ysanity_int.Interfaces{}
    ifcs.Interface = append(ifcs.Interface, &ifc)
    ifcs.YFilter = yfilter.Replace
    
    // Build BGP configuration
    bgp := ysanity_bgp.Bgp{}
    bgp.Global.Config.As = 65172
    neighbor := ysanity_bgp.Bgp_Neighbors_Neighbor{}
    neighbor.NeighborAddress = "172.16.255.2"
    neighbor.Config.NeighborAddress = "172.16.255.2"
    bgp.Neighbors.Neighbor = append(bgp.Neighbors.Neighbor, &neighbor)
    bgp.YFilter = yfilter.Replace

    // Create congfiguration on device
    configEC := types.NewConfig(&ifcs, &bgp)
    reply := suite.Service.Set(&suite.Provider, configEC)
	suite.Equal(reply, true)

	// Change and update configuration
    ifc.Config.Description = "Test"
    ifcs.YFilter = yfilter.Update
    neighbor.Config.PeerAs = 65172
    bgp.YFilter = yfilter.Update
    reply = suite.Service.Set(&suite.Provider, configEC)
    suite.Equal(reply, true)

	// Read all
    filterInt := ysanity_int.Interfaces{}
    filterBgp := ysanity_bgp.Bgp{}
    filterEc := types.NewFilter(&filterInt, &filterBgp)
    readAll := suite.Service.Get(&suite.Provider, filterEc, "ALL")
    suite.NotNil(readAll)
    suite.Equal(types.IsEntityCollection(readAll), true)
    
    // Validate response
    ec := types.EntityToCollection(readAll)
    readEntity, _ := ec.Get("openconfig-interfaces:interfaces")
    suite.NotNil(readEntity)
//    ifcRead := readEntity.(*ysanity_int.Interfaces)
//    ifcConfigJson := suite.Codec.Encode(&suite.CodecPro, &ifcs)
//    ifcReadJson := suite.Codec.Encode(&suite.CodecPro, ifcRead)
//    suite.Equal(ifcConfigJson, ifcReadJson)

    // Read configuration only
    filterInt = ysanity_int.Interfaces{}
    filterBgp = ysanity_bgp.Bgp{}
    filterEc  = types.NewFilter(&filterInt, &filterBgp)
    readConfig := suite.Service.Get(&suite.Provider, filterEc, "CONFIG")
    suite.NotNil(readConfig)
    suite.Equal(types.IsEntityCollection(readConfig), true)

    // Delete configuration
    filterInt = ysanity_int.Interfaces{}
    filterInt.YFilter = yfilter.Delete
    filterBgp = ysanity_bgp.Bgp{}
    filterBgp.YFilter = yfilter.Delete
    filterEc  = types.NewFilter(&filterInt, &filterBgp)
    reply = suite.Service.Set(&suite.Provider, filterEc)
	suite.Equal(reply, true)
}

func configBgp(provider *providers.GnmiServiceProvider) {
    bgp := ysanity_bgp.Bgp{}
    bgp.Global.Config.As = 65172
    neighbor := ysanity_bgp.Bgp_Neighbors_Neighbor{}
    neighbor.NeighborAddress = "172.16.255.2"
    neighbor.Config.NeighborAddress = "172.16.255.2"
    neighbor.Config.PeerAs = 65172
    bgp.Neighbors.Neighbor = append(bgp.Neighbors.Neighbor, &neighbor)
    bgp.YFilter = yfilter.Replace
    
    service := services.GnmiService{}
	service.Set(provider, &bgp)
}

func deleteBgp(provider *providers.GnmiServiceProvider) {
    filterBgp := ysanity_bgp.Bgp{}
    filterBgp.YFilter = yfilter.Delete
    
    service := services.GnmiService{}
	service.Set(provider, &filterBgp)
}

func bgpSubscription(provider *providers.GnmiServiceProvider) {
	previous := ""
	gs := provider.GetSession()
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

func (suite *GnmiServiceTestSuite) TestSubscribeOnce() {
    // Build BGP configuration
	configBgp(&suite.Provider)
	
	// Build subscription request
	sub := services.GnmiSubscription{}
	bgp := ysanity_bgp.Bgp{}
	sub.Entity = &bgp
	sub.SubscriptionMode = "ON_CHANGE"
	var subList []services.GnmiSubscription
	subList = append(subList, sub)
	
	// Start subscription and result monitoring
	go suite.Service.Subscribe(&suite.Provider, subList, 10, "ONCE")
	time.Sleep(100 * time.Millisecond)
	bgpSubscription(&suite.Provider)
	
    // Delete BGP configuration
    deleteBgp(&suite.Provider)
}

func (suite *GnmiServiceTestSuite) TestSubscribeStream() {
    // Build BGP configuration
	configBgp(&suite.Provider)
	
	// Build subscription request
	sub := services.GnmiSubscription{}
	bgp := ysanity_bgp.Bgp{}
	sub.Entity = &bgp
	sub.SubscriptionMode = "ON_CHANGE"
	sub.SuppressRedundant = true
	sub.SampleInterval = 5000000000
	sub.HeartbeatInterval = 15000000000

	var subList []services.GnmiSubscription
	subList = append(subList, sub)
	
	// Start subscription and result monitoring
	go suite.Service.Subscribe(&suite.Provider, subList, 10, "STREAM")
	time.Sleep(100 * time.Millisecond)
	bgpSubscription(&suite.Provider)
	
    // Delete BGP configuration
    deleteBgp(&suite.Provider)
}

func TestGnmiServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GnmiServiceTestSuite))
}
