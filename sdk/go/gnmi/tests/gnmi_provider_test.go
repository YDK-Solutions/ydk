package test

import (
	"fmt"
	"runtime"
	"path/filepath"
	ysanity "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/sanity"
	ysanity_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp"
	ysanity_int "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_interfaces"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/CiscoDevNet/ydk-go/ydk/types/ylist"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GnmiServiceProviderTestSuite struct {
	suite.Suite
	Provider providers.GnmiServiceProvider
	CRUD     services.CrudService
}

func (suite *GnmiServiceProviderTestSuite) SetupSuite() {
	_, callerFile, _, _ := runtime.Caller(0)
	executablePath := filepath.Dir(callerFile)
	repopath := executablePath + "/../../../cpp/core/tests/models"
	repo := types.Repository{Path: repopath}
	suite.Provider = providers.GnmiServiceProvider{
		Repo: repo,
		Address:  "127.0.0.1",
		Port:     50051}
	suite.Provider.Connect()
	
	suite.CRUD   = services.CrudService{}
}

func (suite *GnmiServiceProviderTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *GnmiServiceProviderTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *GnmiServiceProviderTestSuite) TestGnmiCrudSingle() {
    // Build configuration
    ifc := ysanity_int.Interfaces_Interface{}
    ifc.Name = "Loopback10"
    ifc.Config.Name = "Loopback10"

    ifcs := ysanity_int.Interfaces{}
    ifcs.Interface = append(ifcs.Interface, &ifc)

    // Set-replace Request
    reply := suite.CRUD.Create(&suite.Provider, &ifcs)
    suite.Equal(reply, true)

    ifc.Config.Description = "Test"
    reply = suite.CRUD.Update(&suite.Provider, &ifcs)
    suite.Equal(reply, true)

    filter := ysanity_int.Interfaces{}
    ifc_entity := suite.CRUD.Read(&suite.Provider, &filter)
    suite.NotNil(ifc_entity)
    ifcRead := ifc_entity.(*ysanity_int.Interfaces)
    suite.Equal(types.EntityEqual(&ifcs, ifcRead), true)
    
    filter = ysanity_int.Interfaces{}
    ifc_entity = suite.CRUD.ReadConfig(&suite.Provider, &filter)
    suite.NotNil(ifc_entity)

    ifcs = ysanity_int.Interfaces{}
    reply = suite.CRUD.Delete(&suite.Provider, &ifcs)
    suite.Equal(reply, true)
}

func (suite *GnmiServiceProviderTestSuite) TestGnmiCrudMultiple() {
    // Build interface configuration
    ifc := ysanity_int.Interfaces_Interface{}
    ifc.Name = "Loopback10"
    ifc.Config.Name = "Loopback10"
    ifcs := ysanity_int.Interfaces{}
    ifcs.Interface = append(ifcs.Interface, &ifc)
    
    // Build BGP configuration
    bgp := ysanity_bgp.Bgp{}
    bgp.Global.Config.As = 65172
    neighbor := ysanity_bgp.Bgp_Neighbors_Neighbor{}
    neighbor.NeighborAddress = "172.16.255.2"
    neighbor.Config.NeighborAddress = "172.16.255.2"
    bgp.Neighbors.Neighbor = append(bgp.Neighbors.Neighbor, &neighbor)

    // Create congfiguration on device
    configEC := types.NewConfig(&ifcs, &bgp)
    reply := suite.CRUD.Create(&suite.Provider, configEC)
    suite.Equal(reply, true)

    // Change and update configuration
    ifc.Config.Description = "Test"
    neighbor.Config.PeerAs = 65172
    reply = suite.CRUD.Update(&suite.Provider, configEC)
    suite.Equal(reply, true)

    // Read all
    filterInt := ysanity_int.Interfaces{}
    filterBgp := ysanity_bgp.Bgp{}
    filterEc := types.NewFilter(&filterInt, &filterBgp)
    readAll := suite.CRUD.Read(&suite.Provider, filterEc)
    suite.NotNil(readAll)
    suite.Equal(types.IsEntityCollection(readAll), true)
    
    // Validate response
    ec := types.EntityToCollection(readAll)
    readEntity, _ := ec.Get("openconfig-interfaces:interfaces")
    suite.NotNil(readEntity)
    ifcRead := readEntity.(*ysanity_int.Interfaces)
    suite.Equal(types.EntityEqual(&ifcs, ifcRead), true)

    // Read configuration only
    filterInt = ysanity_int.Interfaces{}
    filterBgp = ysanity_bgp.Bgp{}
    filterEc  = types.NewFilter(&filterInt, &filterBgp)
    readConfig := suite.CRUD.ReadConfig(&suite.Provider, filterEc)
    suite.NotNil(readConfig)
    suite.Equal(types.IsEntityCollection(readConfig), true)

    // Delete configuration
    filterInt = ysanity_int.Interfaces{}
    filterBgp = ysanity_bgp.Bgp{}
    filterEc  = types.NewFilter(&filterInt, &filterBgp)
    reply = suite.CRUD.Delete(&suite.Provider, filterEc)
    suite.Equal(reply, true)
}

func (suite *GnmiServiceProviderTestSuite) TestDeleteContainer() {
	// Build loopback configuration
	address := ysanity.Native_Interface_Loopback_Ipv4_Address{}
	address.Ip = "2.2.2.2"
	address.PrefixLength = 32

	loopback := ysanity.Native_Interface_Loopback{}
	loopback.Name = 2222
	loopback.Ipv4.Address = append(loopback.Ipv4.Address, &address)

	native := ysanity.Native{}
        native.Interface.Loopback = append(native.Interface.Loopback, &loopback)

	result := suite.CRUD.Create(&suite.Provider, &native)
	suite.True(result)

	// Read ipv4 configuration
	native = ysanity.Native{}
	loopback = ysanity.Native_Interface_Loopback{}
	loopback.Name = 2222
	native.Interface.Loopback = append(native.Interface.Loopback, &loopback)
	types.SetAllParents(&native)

	ipv4ConfigEnt := suite.CRUD.Read(&suite.Provider, &loopback.Ipv4)
	suite.NotNil(ipv4ConfigEnt)
	ipv4Config := ipv4ConfigEnt.(*ysanity.Native_Interface_Loopback_Ipv4)
	_, addressEnt := ylist.Get(ipv4Config.Address, "2.2.2.2")
	suite.NotNil(addressEnt)
	addressPtr := addressEnt.(*ysanity.Native_Interface_Loopback_Ipv4_Address)
	suite.Equal("32", addressPtr.PrefixLength)

	// Remove ipv4 configuration
	native = ysanity.Native{}
	loopback = ysanity.Native_Interface_Loopback{}
	loopback.Name = 2222
	native.Interface.Loopback = append(native.Interface.Loopback, &loopback)
	types.SetAllParents(&native)

	//loopback.Ipv4.YFilter = yfilter.Delete
	result = suite.CRUD.Delete(&suite.Provider, &loopback.Ipv4)
        suite.True(result)

	// Delete configuration
	native = ysanity.Native{}
	result = suite.CRUD.Delete(&suite.Provider, &native)
	suite.True(result)
}

func TestGnmiServiceProviderTestSuite(t *testing.T) {
    if testing.Verbose() {
        ydk.EnableLogging(ydk.Debug)
    }
    suite.Run(t, new(GnmiServiceProviderTestSuite))
}
