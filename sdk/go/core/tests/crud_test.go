package test

import (
	"fmt"
	ysanity_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp"
	ysanity_bgp_types "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp_types"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NETCONFTestSuite struct {
	suite.Suite
	Provider providers.NetconfServiceProvider
	CRUD     services.CrudService
}

func (suite *NETCONFTestSuite) SetupSuite() {
	suite.CRUD = services.CrudService{}
	suite.Provider = providers.NetconfServiceProvider{
		Address:  "127.0.0.1",
		Username: "admin",
		Password: "admin",
		Port:     12022}
	suite.Provider.Connect()
}

func (suite *NETCONFTestSuite) TearDownSuite() {
	suite.Provider.Disconnect()
}

func (suite *NETCONFTestSuite) TearDownTest() {
	bgpDelete := ysanity_bgp.Bgp{}
	suite.CRUD.Delete(&suite.Provider, &bgpDelete)
}

func (suite *NETCONFTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *NETCONFTestSuite) TestTemplate() {
	// add new test here
}

func (suite *NETCONFTestSuite) TestCreateRead() {
	bgpCreate := ysanity_bgp.Bgp{}
	bgpCreate.Global.Config.As = 65172
	bgpCreate.Global.Config.RouterId = "1.2.3.4"

	ipv6Afisafi := ysanity_bgp.Bgp_Global_AfiSafis_AfiSafi{}
	ipv6Afisafi.AfiSafiName = &ysanity_bgp_types.Ipv6_Unicast{}
	ipv6Afisafi.Config.AfiSafiName = &ysanity_bgp_types.Ipv6_Unicast{}
	ipv6Afisafi.Config.Enabled = true
	bgpCreate.Global.AfiSafis.AfiSafi = append(bgpCreate.Global.AfiSafis.AfiSafi, ipv6Afisafi)

	ipv4Afisafi := ysanity_bgp.Bgp_Global_AfiSafis_AfiSafi{}
	ipv4Afisafi.AfiSafiName = &ysanity_bgp_types.Ipv4_Unicast{}
	ipv4Afisafi.Config.AfiSafiName = &ysanity_bgp_types.Ipv4_Unicast{}
	ipv4Afisafi.Config.Enabled = true
	bgpCreate.Global.AfiSafis.AfiSafi = append(bgpCreate.Global.AfiSafis.AfiSafi, ipv4Afisafi)

	createResult := suite.CRUD.Create(&suite.Provider, &bgpCreate)
	suite.Equal(createResult, true)
	bgpFilter := ysanity_bgp.Bgp{}
	data := suite.CRUD.Read(&suite.Provider, &bgpFilter)
	bgpRead := data.(*ysanity_bgp.Bgp)

	suite.Equal(types.EntityEqual(bgpRead, &bgpCreate), true)
}

func TestNETCONFTestSuite(t *testing.T) {
	suite.Run(t, new(NETCONFTestSuite))
}
