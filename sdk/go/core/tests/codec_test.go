package test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	oc_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/openconfig/bgp"
	oc_bgp_types "github.com/CiscoDevNet/ydk-go/ydk/models/openconfig/types"
	"github.com/CiscoDevNet/ydk-go/ydk/providers"
	"github.com/CiscoDevNet/ydk-go/ydk/services"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

const (
	xml_bgp_payload = `<bgp xmlns="http://openconfig.net/yang/bgp">
  <global>
    <config>
      <as>65172</as>
      <router-id>1.2.3.4</router-id>
    </config>
    <afi-safis>
      <afi-safi>
        <afi-safi-name xmlns:oc-bgp-types="http://openconfig.net/yang/bgp-types">oc-bgp-types:IPV4_UNICAST</afi-safi-name>
        <config>
          <afi-safi-name xmlns:oc-bgp-types="http://openconfig.net/yang/bgp-types">oc-bgp-types:IPV4_UNICAST</afi-safi-name>
          <enabled>true</enabled>
        </config>
      </afi-safi>
    </afi-safis>
  </global>
</bgp>
`

	json_bgp_payload = `{
  "openconfig-bgp:bgp": {
    "global": {
      "config": {
        "as": 65172,
        "router-id": "1.2.3.4"
      },
      "afi-safis": {
        "afi-safi": [
          {
            "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
            "config": {
              "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
              "enabled": true
            }
          }
        ]
      }
    }
  }
}
`
)

func equalPayload(s1, s2 string, cmp func([]uint8, interface{}) error) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = cmp([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = cmp([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

type CodecTestSuite struct {
	suite.Suite
	Provider providers.CodecServiceProvider
	Codec    services.CodecService
}

func (suite *CodecTestSuite) SetupSuite() {
	suite.Codec = services.CodecService{}
	suite.Provider = providers.CodecServiceProvider{}
}

func (suite *CodecTestSuite) BeforeTest(suiteName, testName string) {
	fmt.Printf("%v: %v ...\n", suiteName, testName)
}

func (suite *CodecTestSuite) TearDownSuite() {
}

func config(bgp *oc_bgp.Bgp) {
	bgp.Global.Config.As = 65172 //types.Delete
	bgp.Global.Config.RouterId = "1.2.3.4"

	ipv6_afisafi := oc_bgp.BgpGlobalAfiSafisAfiSafi{}
	ipv6_afisafi.AfiSafiName = &oc_bgp_types.Ipv6UnicastIdentity{}
	ipv6_afisafi.Config.AfiSafiName = &oc_bgp_types.Ipv6UnicastIdentity{}
	ipv6_afisafi.Config.Enabled = true
	bgp.Global.AfiSafis.AfiSafi = append(bgp.Global.AfiSafis.AfiSafi, ipv6_afisafi)

	ipv4_afisafi := oc_bgp.BgpGlobalAfiSafisAfiSafi{}
	ipv4_afisafi.AfiSafiName = &oc_bgp_types.Ipv4UnicastIdentity{}
	ipv4_afisafi.Config.AfiSafiName = &oc_bgp_types.Ipv4UnicastIdentity{}
	ipv4_afisafi.Config.Enabled = true
	bgp.Global.AfiSafis.AfiSafi = append(bgp.Global.AfiSafis.AfiSafi, ipv4_afisafi)
}

func (suite *CodecTestSuite) TestXMLEncoding() {
	bgp := oc_bgp.Bgp{}
	config(&bgp)

	suite.Provider.Encoding = types.XML
	payload := suite.Codec.Encode(&suite.Provider, &bgp)

	fmt.Printf("In TestXMLEncoding, payload = %v", payload)

	result, err := equalPayload(payload, xml_bgp_payload, xml.Unmarshal)
	if err != nil {
		panic("JSONG mashalling failed!")
	}

	suite.Equal(result, true)
}

// func (suite *CodecTestSuite) TestXMLDecoding() {
// 	bgp := oc_bgp.Bgp{}
// 	config(&bgp)

// 	suite.Provider.Encoding = types.XML

// 	entity := suite.Codec.Decode(&suite.Provider, xml_bgp_payload)
// 	bgp_decoded := entity.(*oc_bgp.Bgp)

// 	suite.Equal(types.EntityEqual(&bgp, bgp_decoded, ), true)
// }

func (suite *CodecTestSuite) TestJSONEncoding() {
	bgp := oc_bgp.Bgp{}
	config(&bgp)

	suite.Provider.Encoding = types.JSON
	payload := suite.Codec.Encode(&suite.Provider, &bgp)

	fmt.Printf("In TestJSONEncoding, payload = %v", payload)

	result, err := equalPayload(payload, json_bgp_payload, json.Unmarshal)
	if err != nil {
		panic("JSONG mashalling failed!")
	}

	suite.Equal(result, true)
}

// func (suite *CodecTestSuite) TestJSONDecoding() {
// 	bgp := oc_bgp.Bgp{}
// 	config(&bgp)

// 	suite.Provider.Encoding = types.JSON

// 	entity := suite.Codec.Decode(&suite.Provider, json_bgp_payload)
// 	bgp_decoded := entity.(*oc_bgp.Bgp)

// 	suite.Equal(types.EntityEqual(&bgp, bgp_decoded), true)
// }

func TestCodecTestSuite(t *testing.T) {
	suite.Run(t, new(CodecTestSuite))
}
