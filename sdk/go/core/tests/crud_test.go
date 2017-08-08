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
	bgp_delete := ysanity_bgp.Bgp{}
	suite.CRUD.Delete(&suite.Provider, &bgp_delete)
}

func (suite *NETCONFTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *NETCONFTestSuite) TestTemplate() {
	// add new test here
}

func (suite *NETCONFTestSuite) TestCreateRead() {
	bgp_create := ysanity_bgp.Bgp{}
	bgp_create.Global.Config.As = 65172 //types.Delete
	bgp_create.Global.Config.RouterId = "1.2.3.4"

	ipv6_afisafi := ysanity_bgp.Bgp_Global_AfiSafis_AfiSafi{}
	ipv6_afisafi.AfiSafiName = &ysanity_bgp_types.Ipv6_Unicast{}
	ipv6_afisafi.Config.AfiSafiName = &ysanity_bgp_types.Ipv6_Unicast{}
	ipv6_afisafi.Config.Enabled = true
	bgp_create.Global.AfiSafis.AfiSafi = append(bgp_create.Global.AfiSafis.AfiSafi, ipv6_afisafi)

	ipv4_afisafi := ysanity_bgp.Bgp_Global_AfiSafis_AfiSafi{}
	ipv4_afisafi.AfiSafiName = &ysanity_bgp_types.Ipv4_Unicast{}
	ipv4_afisafi.Config.AfiSafiName = &ysanity_bgp_types.Ipv4_Unicast{}
	ipv4_afisafi.Config.Enabled = true
	bgp_create.Global.AfiSafis.AfiSafi = append(bgp_create.Global.AfiSafis.AfiSafi, ipv4_afisafi)

	create_result := suite.CRUD.Create(&suite.Provider, &bgp_create)
	suite.Equal(create_result, true)
	bgp_filter := ysanity_bgp.Bgp{}
	data := suite.CRUD.Read(&suite.Provider, &bgp_filter)
	bgp_read := data.(*ysanity_bgp.Bgp)

	suite.Equal(types.EntityEqual(bgp_read, &bgp_create), true)
}

func TestNETCONFTestSuite(t *testing.T) {
	suite.Run(t, new(NETCONFTestSuite))
}
