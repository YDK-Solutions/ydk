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

const (
	xml_payload = `<bgp xmlns="http://openconfig.net/yang/bgp">
  <global>
    <afi-safis>
      <afi-safi>
        <afi-safi-name xmlns:oc-bgp-types="http://openconfig.net/yang/bgp-types">oc-bgp-types:IPV4_UNICAST</afi-safi-name>
        <config>
          <afi-safi-name xmlns:oc-bgp-types="http://openconfig.net/yang/bgp-types">oc-bgp-types:IPV4_UNICAST</afi-safi-name>
          <enabled>true</enabled>
        </config>
      </afi-safi>
    </afi-safis>
    <config>
      <as>65172</as>
      <router-id>1.2.3.4</router-id>
    </config>
  </global>
</bgp>
`
)

type GenCodeTestSuite struct {
	suite.Suite
	Provider providers.CodecServiceProvider
	Codec    services.CodecService
}

func (suite *GenCodeTestSuite) SetupSuite() {
	suite.Codec = services.CodecService{}
	suite.Provider = providers.CodecServiceProvider{}
}

func (suite *GenCodeTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *GenCodeTestSuite) TearDownSuite() {
}

func configBgp(bgp *ysanity_bgp.Bgp) {
	bgp.Global.Config.As = 65172 //types.Delete
	bgp.Global.Config.RouterId = "1.2.3.4"

	ipv6_afisafi := ysanity_bgp.Bgp_Global_AfiSafis_AfiSafi{}
	ipv6_afisafi.AfiSafiName = &ysanity_bgp_types.Ipv6_Unicast{}
	ipv6_afisafi.Config.AfiSafiName = &ysanity_bgp_types.Ipv6_Unicast{}
	ipv6_afisafi.Config.Enabled = true
	bgp.Global.AfiSafis.AfiSafi = append(bgp.Global.AfiSafis.AfiSafi, ipv6_afisafi)

	ipv4_afisafi := ysanity_bgp.Bgp_Global_AfiSafis_AfiSafi{}
	ipv4_afisafi.AfiSafiName = &ysanity_bgp_types.Ipv4_Unicast{}
	ipv4_afisafi.Config.AfiSafiName = &ysanity_bgp_types.Ipv4_Unicast{}
	ipv4_afisafi.Config.Enabled = true
	bgp.Global.AfiSafis.AfiSafi = append(bgp.Global.AfiSafis.AfiSafi, ipv4_afisafi)
}

func (suite *GenCodeTestSuite) TestGenCodeXMLEncoding() {
	bgp := ysanity_bgp.Bgp{}
	configBgp(&bgp)

	suite.Provider.Encoding = types.XML
	payload := suite.Codec.Encode(&suite.Provider, &bgp)

	fmt.Printf("In TestGenCodeXMLEncoding, payload = %v", payload)
}

func (suite *GenCodeTestSuite) TestGenCodeXMLDecoding() {
	entity := suite.Codec.Decode(&suite.Provider, xml_payload)
	bgp_decoded := entity.(*ysanity_bgp.Bgp)

	p := suite.Codec.Encode(&suite.Provider, bgp_decoded)
	fmt.Printf("In TestGenCodeXMlDecoding, payload = %v", p)
}

func TestGenCodeTestSuite(t *testing.T) {
	suite.Run(t, new(GenCodeTestSuite))
}
