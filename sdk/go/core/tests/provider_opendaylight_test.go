package test

import (
	"fmt"
	oc_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types/protocol"
	encoding "github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OpenDaylightProviderTestSuite struct {
	suite.Suite
	Provider providers.OpenDaylightServiceProvider
	CRUD     services.CrudService
}

func (suite *OpenDaylightProviderTestSuite) SetupSuite() {
	suite.CRUD = services.CrudService{}
	suite.Provider = providers.OpenDaylightServiceProvider{
		Path:           TestHome,
		Address:        "localhost",
		Username:       "admin",
		Password:       "admin",
		Port:           12306,
		EncodingFormat: encoding.JSON,
		Protocol:       protocol.Restconf}
	suite.Provider.Connect()
}

func (suite *OpenDaylightProviderTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *OpenDaylightProviderTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *OpenDaylightProviderTestSuite) TestReadODL() {
	bgp := oc_bgp.Bgp{}

	nodeIDs := suite.Provider.GetNodeIDs()

	ydk.YLogDebug(fmt.Sprintf("nodeIDs: %v\n", nodeIDs))

	provider := suite.Provider.GetNodeProvider("xr")

	entity := suite.CRUD.ReadConfig(provider, &bgp)

	bgpRead := entity.(*oc_bgp.Bgp)

	ydk.YLogDebug("BGP configuration: ")
	ydk.YLogDebug(fmt.Sprintf("AS: %v\n", bgpRead.Global.Config.As))
	ydk.YLogDebug(fmt.Sprintf("Router ID: %v\n", bgpRead.Global.Config.RouterId))

	suite.Equal(bgpRead.Global.Config.As, "65172")
	suite.Equal(bgpRead.Global.Config.RouterId, "1.2.3.4")
}

func (suite *OpenDaylightProviderTestSuite) TestCreateODL() {
	bgp := oc_bgp.Bgp{}
	bgp.Global.Config.As = 65172
	bgp.Global.Config.RouterId = "1.2.3.4"

	neighbor := oc_bgp.Bgp_Neighbors_Neighbor{}
	neighbor.NeighborAddress = "6.7.8.9"
	neighbor.Config.NeighborAddress = "6.7.8.9"
	neighbor.Config.PeerAs = 65001
	neighbor.Config.LocalAs = 65001
	neighbor.Config.PeerGroup = "IBGP"
	bgp.Neighbors.Neighbor = append(bgp.Neighbors.Neighbor, &neighbor)

	peerGroup := oc_bgp.Bgp_PeerGroups_PeerGroup{}
	peerGroup.PeerGroupName = "IBGP"
	peerGroup.Config.PeerGroupName = "IBGP"
	peerGroup.Config.Description = "test description"
	peerGroup.Config.PeerAs = 65001
	peerGroup.Config.LocalAs = 65001
	bgp.PeerGroups.PeerGroup = append(bgp.PeerGroups.PeerGroup, &peerGroup)

	provider := suite.Provider.GetNodeProvider("xr")
	result := suite.CRUD.Create(provider, &bgp)

	suite.Equal(result, true)
}

func TestOpenDaylightProviderTestSuite(t *testing.T) {
	if testing.Verbose() {
		ydk.EnableLogging(ydk.Debug)
	}
	suite.Run(t, new(OpenDaylightProviderTestSuite))
}
