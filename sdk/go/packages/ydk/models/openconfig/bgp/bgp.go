package bgp

import (
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/models/openconfig"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"reflect"
)

func init() {
	fmt.Println("Registering top level entities for package bgp...")
	ydk.RegisterEntity("{http://openconfig.net/yang/bgp bgp}", reflect.TypeOf(Bgp{}))
	ydk.RegisterEntity("openconfig-bgp:bgp", reflect.TypeOf(Bgp{}))
}

//////////////////////////////////////////////////////////////////////////
// Bgp
//////////////////////////////////////////////////////////////////////////
type Bgp struct {
	Filter types.YFilter

	Global BgpGlobal
}

func (bgp *Bgp) GetEntityPath(entity types.Entity) types.EntityPath {
	return types.EntityPath{Path: bgp.GetSegmentPath()}
}

func (bgp *Bgp) GetSegmentPath() string {
	return "openconfig-bgp:bgp"
}

func (bgp *Bgp) HasDataOrFilter() bool {
	return bgp.Global.HasDataOrFilter()
}

func (bgp *Bgp) SetValue(value_path, value string) {

}

func (bgp *Bgp) GetChildByName(child_yang_name, segment_path string) types.Entity {
	if child_yang_name == "global" {
		return &bgp.Global
	}
	return nil
}

func (bgp *Bgp) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["global"] = &bgp.Global

	return children
}

func (bgp *Bgp) SetParent(parent types.Entity) {

}

func (bgp *Bgp) GetParent() types.Entity {
	return nil
}

func (bgp *Bgp) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return openconfig.OpenconfigAugmentLookupTables
}

func (bgp *Bgp) GetBundleYangModelsLocation() string {
	return openconfig.YdkOpenconfigModelsPath
}

func (bgp *Bgp) GetBundleName() string {
	return "openconfig"
}

func (bgp *Bgp) GetYangName() string {
	return "bgp"
}

func (bgp *Bgp) GetParentYangName() string {
	return "openconfigbgp"
}

func (bgp *Bgp) GetFilter() types.YFilter {
	return bgp.Filter
}

//////////////////////////////////////////////////////////////////////////
// BgpGlobal
//////////////////////////////////////////////////////////////////////////
type BgpGlobal struct {
	parent types.Entity
	Filter types.YFilter

	Config   BgpGlobalConfig
	AfiSafis BgpGlobalAfiSafis
}

func (global *BgpGlobal) GetEntityPath(entity types.Entity) types.EntityPath {
	return types.EntityPath{Path: global.GetSegmentPath()}
}

func (global *BgpGlobal) GetSegmentPath() string {
	return "global"
}

func (global *BgpGlobal) HasDataOrFilter() bool {
	return global.Config.HasDataOrFilter()
}

func (global *BgpGlobal) SetValue(value_path, value string) {

}

func (global *BgpGlobal) GetChildByName(child_yang_name, segment_path string) types.Entity {
	if child_yang_name == "config" {
		return &global.Config
	} else if child_yang_name == "afi-safis" {
		return &global.AfiSafis
	}
	return nil
}

func (global *BgpGlobal) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["global"] = &global.Config
	children["afi-safis"] = &global.AfiSafis

	return children
}

func (global *BgpGlobal) SetParent(parent types.Entity) {
	global.parent = parent
}

func (global *BgpGlobal) GetParent() types.Entity {
	return global.parent
}

func (global *BgpGlobal) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return openconfig.OpenconfigAugmentLookupTables
}

func (global *BgpGlobal) GetBundleYangModelsLocation() string {
	return openconfig.YdkOpenconfigModelsPath
}

func (global *BgpGlobal) GetBundleName() string {
	return "openconfig"
}

func (global *BgpGlobal) GetYangName() string {
	return "global"
}

func (global *BgpGlobal) GetParentYangName() string {
	return "bgp"
}

func (global *BgpGlobal) GetFilter() types.YFilter {
	return global.Filter
}

//////////////////////////////////////////////////////////////////////////
// BgpGlobalConfig
//////////////////////////////////////////////////////////////////////////
type BgpGlobalConfig struct {
	parent types.Entity
	Filter types.YFilter

	As       interface{} // uint32
	RouterId interface{} //string
}

func (config *BgpGlobalConfig) GetEntityPath(entity types.Entity) types.EntityPath {
	var leaf_data types.LeafData
	e := types.EntityPath{Path: config.GetSegmentPath()}

	if config.As != nil {
		switch config.As.(type) {
		case types.YFilter:
			leaf_data = types.LeafData{IsSet: false, Filter: config.As.(types.YFilter)}
		default:
			leaf_data = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", config.As)}
		}
		e.ValuePaths = append(e.ValuePaths, types.NameLeafData{Name: "as", Data: leaf_data})
	}

	if config.RouterId != nil {
		switch config.RouterId.(type) {
		case types.YFilter:
			leaf_data = types.LeafData{IsSet: false, Filter: config.RouterId.(types.YFilter)}
		default:
			leaf_data = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", config.RouterId)}
		}
		e.ValuePaths = append(e.ValuePaths, types.NameLeafData{Name: "router-id", Data: leaf_data})
	}

	return e
}

func (config *BgpGlobalConfig) GetSegmentPath() string {
	return "config"
}

func (config *BgpGlobalConfig) HasDataOrFilter() bool {
	return config.As != nil || config.RouterId != nil
}

func (config *BgpGlobalConfig) SetValue(value_path, value string) {
	if value_path == "as" {
		config.As = value
	} else if value_path == "router-id" {
		config.RouterId = value
	}
}

func (config *BgpGlobalConfig) GetChildByName(child_yang_name, segment_path string) types.Entity {
	return nil
}

func (config *BgpGlobalConfig) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (config *BgpGlobalConfig) SetParent(parent types.Entity) {
	config.parent = parent
}

func (config *BgpGlobalConfig) GetParent() types.Entity {
	return config.parent
}

func (config *BgpGlobalConfig) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return openconfig.OpenconfigAugmentLookupTables
}

func (config *BgpGlobalConfig) GetBundleYangModelsLocation() string {
	return openconfig.YdkOpenconfigModelsPath
}

func (config *BgpGlobalConfig) GetBundleName() string {
	return "openconfig"
}

func (config *BgpGlobalConfig) GetYangName() string {
	return "global"
}

func (config *BgpGlobalConfig) GetParentYangName() string {
	return "bgp"
}

func (config *BgpGlobalConfig) GetFilter() types.YFilter {
	return config.Filter
}

//////////////////////////////////////////////////////////////////////////
// BgpGlobalAfiSafis
//////////////////////////////////////////////////////////////////////////
type BgpGlobalAfiSafis struct {
	parent types.Entity
	Filter types.YFilter

	AfiSafi []BgpGlobalAfiSafisAfiSafi
}

func (afisafis *BgpGlobalAfiSafis) GetEntityPath(entity types.Entity) types.EntityPath {
	return types.EntityPath{Path: afisafis.GetSegmentPath()}
	//ValuePaths:{ types.NameLeafData{Name:"as", Data:types.LeafData{IsSet: true, Value:config.As.Value}}},
	//types.NameLeafData{Name:"router-id", Data:types.LeafData{IsSet: true, Value:config.RouterId.Value}}
}

func (afisafis *BgpGlobalAfiSafis) GetSegmentPath() string {
	return "afi-safis"
}

func (afisafis *BgpGlobalAfiSafis) HasDataOrFilter() bool {
	for _, afisafi := range afisafis.AfiSafi {
		if afisafi.HasDataOrFilter() {
			return true
		}
	}
	return false
}

func (afisafis *BgpGlobalAfiSafis) SetValue(value_path, value string) {

}

func (afisafis *BgpGlobalAfiSafis) GetChildByName(child_yang_name, segment_path string) types.Entity {
	if child_yang_name == "afi-safi" {
		for _, afisafi := range afisafis.AfiSafi {
			if afisafi.GetSegmentPath() == segment_path {
				return &afisafi
			}
		}
		afi_safi := BgpGlobalAfiSafisAfiSafi{}
		afisafis.AfiSafi = append(afisafis.AfiSafi, afi_safi)
		return &afi_safi
	}
	return nil
}

func (afisafis *BgpGlobalAfiSafis) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for _, afisafi := range afisafis.AfiSafi {
		children[afisafi.GetSegmentPath()] = &afisafi
	}

	return children
}

func (afisafis *BgpGlobalAfiSafis) SetParent(parent types.Entity) {
	afisafis.parent = parent
}

func (afisafis *BgpGlobalAfiSafis) GetParent() types.Entity {
	return afisafis.parent
}

func (afisafis *BgpGlobalAfiSafis) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return openconfig.OpenconfigAugmentLookupTables
}

func (afisafis *BgpGlobalAfiSafis) GetBundleYangModelsLocation() string {
	return openconfig.YdkOpenconfigModelsPath
}

func (afisafis *BgpGlobalAfiSafis) GetBundleName() string {
	return "openconfig"
}

func (afisafis *BgpGlobalAfiSafis) GetYangName() string {
	return "global"
}

func (afisafis *BgpGlobalAfiSafis) GetParentYangName() string {
	return "bgp"
}

func (afisafis *BgpGlobalAfiSafis) GetFilter() types.YFilter {
	return afisafis.Filter
}

//////////////////////////////////////////////////////////////////////////
// BgpGlobalAfiSafisAfiSafi
//////////////////////////////////////////////////////////////////////////
type BgpGlobalAfiSafisAfiSafi struct {
	parent types.Entity
	Filter types.YFilter

	AfiSafiName interface{} // string

	Config BgpGlobalAfiSafisAfiSafiConfig
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetEntityPath(entity types.Entity) types.EntityPath {
	r := types.EntityPath{Path: afisafi.GetSegmentPath()}
	r.ValuePaths = append(r.ValuePaths, types.NameLeafData{Name: "afi-safi-name", Data: types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", afisafi.AfiSafiName)}})
	return r
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetSegmentPath() string {
	return "afi-safi[afi-safi-name='" + fmt.Sprintf("%v", afisafi.AfiSafiName) + "']"
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) HasDataOrFilter() bool {
	return afisafi.AfiSafiName != nil || afisafi.Config.HasDataOrFilter()
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) SetValue(value_path, value string) {
	if value_path == "afi-safi-name" {
		afisafi.AfiSafiName = value
	}
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetChildByName(child_yang_name, segment_path string) types.Entity {
	if child_yang_name == "config" {
		return &afisafi.Config
	}
	return nil
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["config"] = &afisafi.Config

	return children
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) SetParent(parent types.Entity) {
	afisafi.parent = parent
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetParent() types.Entity {
	return afisafi.parent
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return openconfig.OpenconfigAugmentLookupTables
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetBundleYangModelsLocation() string {
	return openconfig.YdkOpenconfigModelsPath
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetBundleName() string {
	return "openconfig"
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetYangName() string {
	return "afi-safi"
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetParentYangName() string {
	return "bgp"
}

func (afisafi *BgpGlobalAfiSafisAfiSafi) GetFilter() types.YFilter {
	return afisafi.Filter
}

//////////////////////////////////////////////////////////////////////////
// BgpGlobalAfiSafisAfiSafiConfig
//////////////////////////////////////////////////////////////////////////
type BgpGlobalAfiSafisAfiSafiConfig struct {
	parent types.Entity
	Filter types.YFilter

	AfiSafiName interface{}
	Enabled     interface{}
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetEntityPath(entity types.Entity) types.EntityPath {
	r := types.EntityPath{Path: config.GetSegmentPath()}
	r.ValuePaths = append(r.ValuePaths, types.NameLeafData{Name: "afi-safi-name", Data: types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", config.AfiSafiName)}})
	r.ValuePaths = append(r.ValuePaths, types.NameLeafData{Name: "enabled", Data: types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", config.Enabled)}})
	return r
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetSegmentPath() string {
	return "config"
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) HasDataOrFilter() bool {
	return config.AfiSafiName != nil || config.Enabled != nil
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) SetValue(value_path, value string) {

	if value_path == "afi-safi-name" {
		config.AfiSafiName = value
	} else if value_path == "enabled" {
		config.Enabled = value
	}
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetChildByName(child_yang_name, segment_path string) types.Entity {
	return nil
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) SetParent(parent types.Entity) {
	config.parent = parent
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetParent() types.Entity {
	return config.parent
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return openconfig.OpenconfigAugmentLookupTables
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetBundleYangModelsLocation() string {
	return openconfig.YdkOpenconfigModelsPath
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetBundleName() string {
	return "openconfig"
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetYangName() string {
	return "config"
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetParentYangName() string {
	return "afi-safi"
}

func (config *BgpGlobalAfiSafisAfiSafiConfig) GetFilter() types.YFilter {
	return config.Filter
}
