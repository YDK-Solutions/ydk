package ydktest_sanity

import (
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk/models/ydktest"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
)

//////////////////////////////////////////////////////////////////////////
// BaseIdentity
//////////////////////////////////////////////////////////////////////////
type BaseIdentity struct {
	parent types.Entity
	Filter types.YFilter
}

//////////////////////////////////////////////////////////////////////////
// ChildIdentity
//////////////////////////////////////////////////////////////////////////
type ChildIdentity struct {
	parent types.Entity
	Filter types.YFilter
}

//////////////////////////////////////////////////////////////////////////
// ChildChildIdentity
//////////////////////////////////////////////////////////////////////////
type ChildChildIdentity struct {
	parent types.Entity
	Filter types.YFilter
}

//////////////////////////////////////////////////////////////////////////
// SubTest
//////////////////////////////////////////////////////////////////////////
type SubTest struct {
	parent types.Entity
	Filter types.YFilter

	OneAug SubTest_OneAug
}

func (subTest *SubTest) HasDataOrFilter() bool {
	return subTest.OneAug.HasDataOrFilter()
}

func (subTest *SubTest) GetFilter() types.YFilter {
	return subTest.Filter
}

func (subTest *SubTest) GetSegmentPath() string {
	return "ydktest-sanity:sub-test"
}

func (subTest *SubTest) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: subTest.GetSegmentPath()}
	return entityPath
}

func (subTest *SubTest) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "one-aug" {
		return &subTest.OneAug
	}
	return nil
}

func (subTest *SubTest) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["one-aug"] = &subTest.OneAug

	return children
}

func (subTest *SubTest) SetValue(value_path string, value string) {
}

func (subTest *SubTest) GetBundleName() string {
	return "ydktest"
}

func (subTest *SubTest) GetYangName() string {
	return "subTest"
}

func (subTest *SubTest) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (subTest *SubTest) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (subTest *SubTest) SetParent(parent types.Entity) {
	subTest.parent = parent
}

func (subTest *SubTest) GetParent() types.Entity {
	return subTest.parent
}

func (subTest *SubTest) GetParentYangName() string {
	return "ydktest-sanity"
}

//////////////////////////////////////////////////////////////////////////
// SubTest_OneAug
//////////////////////////////////////////////////////////////////////////
type SubTest_OneAug struct {
	parent types.Entity
	Filter types.YFilter

	Name   interface{} // str
	Number interface{} // int32
}

func (oneAug *SubTest_OneAug) HasDataOrFilter() bool {
	return oneAug.Name != nil ||
		oneAug.Number != nil
}

func (oneAug *SubTest_OneAug) GetFilter() types.YFilter {
	return oneAug.Filter
}

func (oneAug *SubTest_OneAug) GetSegmentPath() string {
	return "one-aug"
}

func (oneAug *SubTest_OneAug) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: oneAug.GetSegmentPath()}
	var leafData types.LeafData

	if oneAug.Name != nil {
		switch oneAug.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: oneAug.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", oneAug.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	if oneAug.Number != nil {
		switch oneAug.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: oneAug.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", oneAug.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (oneAug *SubTest_OneAug) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (oneAug *SubTest_OneAug) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (oneAug *SubTest_OneAug) SetValue(value_path string, value string) {
	if value_path == "name" {
		oneAug.Name = value
	}
	if value_path == "number" {
		oneAug.Number = value
	}
}

func (oneAug *SubTest_OneAug) GetBundleName() string {
	return "ydktest"
}

func (oneAug *SubTest_OneAug) GetYangName() string {
	return "oneAug"
}

func (oneAug *SubTest_OneAug) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (oneAug *SubTest_OneAug) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (oneAug *SubTest_OneAug) SetParent(parent types.Entity) {
	oneAug.parent = parent
}

func (oneAug *SubTest_OneAug) GetParent() types.Entity {
	return oneAug.parent
}

func (oneAug *SubTest_OneAug) GetParentYangName() string {
	return "sub-test"
}

//////////////////////////////////////////////////////////////////////////
// Runner
//////////////////////////////////////////////////////////////////////////
type Runner struct {
	parent types.Entity
	Filter types.YFilter

	InbtwList     Runner_InbtwList
	LeafRef       Runner_LeafRef
	NotSupported1 Runner_NotSupported1
	NotSupported2 []Runner_NotSupported2
	One           Runner_One
	OneList       Runner_OneList
	Runner2       Runner_Runner2 // presence node
	Three         Runner_Three
	ThreeList     Runner_ThreeList
	Two           Runner_Two
	TwoList       Runner_TwoList
	Ytypes        Runner_Ytypes
}

func (runner *Runner) HasDataOrFilter() bool {
	for _, child := range runner.NotSupported2 {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return runner.InbtwList.HasDataOrFilter() ||
		runner.LeafRef.HasDataOrFilter() ||
		runner.NotSupported1.HasDataOrFilter() ||
		runner.One.HasDataOrFilter() ||
		runner.OneList.HasDataOrFilter() ||
		runner.Runner2.HasDataOrFilter() ||
		runner.Three.HasDataOrFilter() ||
		runner.ThreeList.HasDataOrFilter() ||
		runner.Two.HasDataOrFilter() ||
		runner.TwoList.HasDataOrFilter() ||
		runner.Ytypes.HasDataOrFilter()
}

func (runner *Runner) GetFilter() types.YFilter {
	return runner.Filter
}

func (runner *Runner) GetSegmentPath() string {
	return "ydktest-sanity:runner"
}

func (runner *Runner) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: runner.GetSegmentPath()}
	return entityPath
}

func (runner *Runner) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "inbtw-list" {
		return &runner.InbtwList
	}
	if child_yang_name == "leaf-ref" {
		return &runner.LeafRef
	}
	if child_yang_name == "not-supported-1" {
		return &runner.NotSupported1
	}
	if child_yang_name == "not-supported-2" {
		for _, c := range runner.NotSupported2 {
			if runner.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_NotSupported2{}
		runner.NotSupported2 = append(runner.NotSupported2, child)
		return &runner.NotSupported2[len(runner.NotSupported2)-1]
	}
	if child_yang_name == "one" {
		return &runner.One
	}
	if child_yang_name == "one-list" {
		return &runner.OneList
	}
	if child_yang_name == "runner-2" {
		return &runner.Runner2
	}
	if child_yang_name == "three" {
		return &runner.Three
	}
	if child_yang_name == "three-list" {
		return &runner.ThreeList
	}
	if child_yang_name == "two" {
		return &runner.Two
	}
	if child_yang_name == "two-list" {
		return &runner.TwoList
	}
	if child_yang_name == "ytypes" {
		return &runner.Ytypes
	}
	return nil
}

func (runner *Runner) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["inbtw-list"] = &runner.InbtwList
	children["leaf-ref"] = &runner.LeafRef
	children["not-supported-1"] = &runner.NotSupported1
	for i := range runner.NotSupported2 {
		children[runner.NotSupported2[i].GetSegmentPath()] = &runner.NotSupported2[i]
	}
	children["one"] = &runner.One
	children["one-list"] = &runner.OneList
	children["runner-2"] = &runner.Runner2
	children["three"] = &runner.Three
	children["three-list"] = &runner.ThreeList
	children["two"] = &runner.Two
	children["two-list"] = &runner.TwoList
	children["ytypes"] = &runner.Ytypes

	return children
}

func (runner *Runner) SetValue(value_path string, value string) {
}

func (runner *Runner) GetBundleName() string {
	return "ydktest"
}

func (runner *Runner) GetYangName() string {
	return "runner"
}

func (runner *Runner) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (runner *Runner) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (runner *Runner) SetParent(parent types.Entity) {
	runner.parent = parent
}

func (runner *Runner) GetParent() types.Entity {
	return runner.parent
}

func (runner *Runner) GetParentYangName() string {
	return "ydktest-sanity"
}

//////////////////////////////////////////////////////////////////////////
// Runner_One
//////////////////////////////////////////////////////////////////////////
type Runner_One struct {
	parent types.Entity
	Filter types.YFilter

	AugmentedLeaf interface{} // str
	Config        interface{} // str
	Name          interface{} // str
	Number        interface{} // int32
	OneAug        Runner_One_OneAug
}

func (one *Runner_One) HasDataOrFilter() bool {
	return one.AugmentedLeaf != nil ||
		one.Config != nil ||
		one.Name != nil ||
		one.Number != nil ||
		one.OneAug.HasDataOrFilter()
}

func (one *Runner_One) GetFilter() types.YFilter {
	return one.Filter
}

func (one *Runner_One) GetSegmentPath() string {
	return "one"
}

func (one *Runner_One) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: one.GetSegmentPath()}
	var leafData types.LeafData

	if one.AugmentedLeaf != nil {
		switch one.AugmentedLeaf.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: one.AugmentedLeaf.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", one.AugmentedLeaf)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "augmented-leaf", Data: leafData})
		}
	}

	if one.Config != nil {
		switch one.Config.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: one.Config.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", one.Config)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "config", Data: leafData})
		}
	}

	if one.Name != nil {
		switch one.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: one.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", one.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	if one.Number != nil {
		switch one.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: one.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", one.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (one *Runner_One) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "one-aug" {
		return &one.OneAug
	}
	return nil
}

func (one *Runner_One) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["one-aug"] = &one.OneAug

	return children
}

func (one *Runner_One) SetValue(value_path string, value string) {
	if value_path == "augmented-leaf" {
		one.AugmentedLeaf = value
	}
	if value_path == "config" {
		one.Config = value
	}
	if value_path == "name" {
		one.Name = value
	}
	if value_path == "number" {
		one.Number = value
	}
}

func (one *Runner_One) GetBundleName() string {
	return "ydktest"
}

func (one *Runner_One) GetYangName() string {
	return "one"
}

func (one *Runner_One) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (one *Runner_One) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (one *Runner_One) SetParent(parent types.Entity) {
	one.parent = parent
}

func (one *Runner_One) GetParent() types.Entity {
	return one.parent
}

func (one *Runner_One) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_One_OneAug
//////////////////////////////////////////////////////////////////////////
type Runner_One_OneAug struct {
	parent types.Entity
	Filter types.YFilter

	Name   interface{} // str
	Number interface{} // int32
}

func (oneAug *Runner_One_OneAug) HasDataOrFilter() bool {
	return oneAug.Name != nil ||
		oneAug.Number != nil
}

func (oneAug *Runner_One_OneAug) GetFilter() types.YFilter {
	return oneAug.Filter
}

func (oneAug *Runner_One_OneAug) GetSegmentPath() string {
	return "ydktest-sanity-augm:one-aug"
}

func (oneAug *Runner_One_OneAug) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: oneAug.GetSegmentPath()}
	var leafData types.LeafData

	if oneAug.Name != nil {
		switch oneAug.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: oneAug.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", oneAug.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	if oneAug.Number != nil {
		switch oneAug.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: oneAug.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", oneAug.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (oneAug *Runner_One_OneAug) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (oneAug *Runner_One_OneAug) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (oneAug *Runner_One_OneAug) SetValue(value_path string, value string) {
	if value_path == "name" {
		oneAug.Name = value
	}
	if value_path == "number" {
		oneAug.Number = value
	}
}

func (oneAug *Runner_One_OneAug) GetBundleName() string {
	return "ydktest"
}

func (oneAug *Runner_One_OneAug) GetYangName() string {
	return "oneAug"
}

func (oneAug *Runner_One_OneAug) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (oneAug *Runner_One_OneAug) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (oneAug *Runner_One_OneAug) SetParent(parent types.Entity) {
	oneAug.parent = parent
}

func (oneAug *Runner_One_OneAug) GetParent() types.Entity {
	return oneAug.parent
}

func (oneAug *Runner_One_OneAug) GetParentYangName() string {
	return "one"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Two
//////////////////////////////////////////////////////////////////////////
type Runner_Two struct {
	parent types.Entity
	Filter types.YFilter

	Name   interface{} // str
	Number interface{} // int32
	Sub1   Runner_Two_Sub1
}

func (two *Runner_Two) HasDataOrFilter() bool {
	return two.Name != nil ||
		two.Number != nil ||
		two.Sub1.HasDataOrFilter()
}

func (two *Runner_Two) GetFilter() types.YFilter {
	return two.Filter
}

func (two *Runner_Two) GetSegmentPath() string {
	return "two"
}

func (two *Runner_Two) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: two.GetSegmentPath()}
	var leafData types.LeafData

	if two.Name != nil {
		switch two.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: two.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", two.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	if two.Number != nil {
		switch two.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: two.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", two.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (two *Runner_Two) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "sub1" {
		return &two.Sub1
	}
	return nil
}

func (two *Runner_Two) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["sub1"] = &two.Sub1

	return children
}

func (two *Runner_Two) SetValue(value_path string, value string) {
	if value_path == "name" {
		two.Name = value
	}
	if value_path == "number" {
		two.Number = value
	}
}

func (two *Runner_Two) GetBundleName() string {
	return "ydktest"
}

func (two *Runner_Two) GetYangName() string {
	return "two"
}

func (two *Runner_Two) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (two *Runner_Two) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (two *Runner_Two) SetParent(parent types.Entity) {
	two.parent = parent
}

func (two *Runner_Two) GetParent() types.Entity {
	return two.parent
}

func (two *Runner_Two) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Two_Sub1
//////////////////////////////////////////////////////////////////////////
type Runner_Two_Sub1 struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
}

func (sub1 *Runner_Two_Sub1) HasDataOrFilter() bool {
	return sub1.Number != nil
}

func (sub1 *Runner_Two_Sub1) GetFilter() types.YFilter {
	return sub1.Filter
}

func (sub1 *Runner_Two_Sub1) GetSegmentPath() string {
	return "sub1"
}

func (sub1 *Runner_Two_Sub1) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: sub1.GetSegmentPath()}
	var leafData types.LeafData

	if sub1.Number != nil {
		switch sub1.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: sub1.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", sub1.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (sub1 *Runner_Two_Sub1) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (sub1 *Runner_Two_Sub1) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (sub1 *Runner_Two_Sub1) SetValue(value_path string, value string) {
	if value_path == "number" {
		sub1.Number = value
	}
}

func (sub1 *Runner_Two_Sub1) GetBundleName() string {
	return "ydktest"
}

func (sub1 *Runner_Two_Sub1) GetYangName() string {
	return "sub1"
}

func (sub1 *Runner_Two_Sub1) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (sub1 *Runner_Two_Sub1) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (sub1 *Runner_Two_Sub1) SetParent(parent types.Entity) {
	sub1.parent = parent
}

func (sub1 *Runner_Two_Sub1) GetParent() types.Entity {
	return sub1.parent
}

func (sub1 *Runner_Two_Sub1) GetParentYangName() string {
	return "two"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Three
//////////////////////////////////////////////////////////////////////////
type Runner_Three struct {
	parent types.Entity
	Filter types.YFilter

	Name   interface{} // str
	Number interface{} // int32
	Sub1   Runner_Three_Sub1
}

func (three *Runner_Three) HasDataOrFilter() bool {
	return three.Name != nil ||
		three.Number != nil ||
		three.Sub1.HasDataOrFilter()
}

func (three *Runner_Three) GetFilter() types.YFilter {
	return three.Filter
}

func (three *Runner_Three) GetSegmentPath() string {
	return "three"
}

func (three *Runner_Three) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: three.GetSegmentPath()}
	var leafData types.LeafData

	if three.Name != nil {
		switch three.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: three.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", three.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	if three.Number != nil {
		switch three.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: three.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", three.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (three *Runner_Three) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "sub1" {
		return &three.Sub1
	}
	return nil
}

func (three *Runner_Three) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["sub1"] = &three.Sub1

	return children
}

func (three *Runner_Three) SetValue(value_path string, value string) {
	if value_path == "name" {
		three.Name = value
	}
	if value_path == "number" {
		three.Number = value
	}
}

func (three *Runner_Three) GetBundleName() string {
	return "ydktest"
}

func (three *Runner_Three) GetYangName() string {
	return "three"
}

func (three *Runner_Three) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (three *Runner_Three) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (three *Runner_Three) SetParent(parent types.Entity) {
	three.parent = parent
}

func (three *Runner_Three) GetParent() types.Entity {
	return three.parent
}

func (three *Runner_Three) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Three_Sub1
//////////////////////////////////////////////////////////////////////////
type Runner_Three_Sub1 struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Sub2   Runner_Three_Sub1_Sub2
}

func (sub1 *Runner_Three_Sub1) HasDataOrFilter() bool {
	return sub1.Number != nil ||
		sub1.Sub2.HasDataOrFilter()
}

func (sub1 *Runner_Three_Sub1) GetFilter() types.YFilter {
	return sub1.Filter
}

func (sub1 *Runner_Three_Sub1) GetSegmentPath() string {
	return "sub1"
}

func (sub1 *Runner_Three_Sub1) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: sub1.GetSegmentPath()}
	var leafData types.LeafData

	if sub1.Number != nil {
		switch sub1.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: sub1.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", sub1.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (sub1 *Runner_Three_Sub1) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "sub2" {
		return &sub1.Sub2
	}
	return nil
}

func (sub1 *Runner_Three_Sub1) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["sub2"] = &sub1.Sub2

	return children
}

func (sub1 *Runner_Three_Sub1) SetValue(value_path string, value string) {
	if value_path == "number" {
		sub1.Number = value
	}
}

func (sub1 *Runner_Three_Sub1) GetBundleName() string {
	return "ydktest"
}

func (sub1 *Runner_Three_Sub1) GetYangName() string {
	return "sub1"
}

func (sub1 *Runner_Three_Sub1) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (sub1 *Runner_Three_Sub1) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (sub1 *Runner_Three_Sub1) SetParent(parent types.Entity) {
	sub1.parent = parent
}

func (sub1 *Runner_Three_Sub1) GetParent() types.Entity {
	return sub1.parent
}

func (sub1 *Runner_Three_Sub1) GetParentYangName() string {
	return "three"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Three_Sub1_Sub2
//////////////////////////////////////////////////////////////////////////
type Runner_Three_Sub1_Sub2 struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
}

func (sub2 *Runner_Three_Sub1_Sub2) HasDataOrFilter() bool {
	return sub2.Number != nil
}

func (sub2 *Runner_Three_Sub1_Sub2) GetFilter() types.YFilter {
	return sub2.Filter
}

func (sub2 *Runner_Three_Sub1_Sub2) GetSegmentPath() string {
	return "sub2"
}

func (sub2 *Runner_Three_Sub1_Sub2) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: sub2.GetSegmentPath()}
	var leafData types.LeafData

	if sub2.Number != nil {
		switch sub2.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: sub2.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", sub2.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (sub2 *Runner_Three_Sub1_Sub2) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (sub2 *Runner_Three_Sub1_Sub2) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (sub2 *Runner_Three_Sub1_Sub2) SetValue(value_path string, value string) {
	if value_path == "number" {
		sub2.Number = value
	}
}

func (sub2 *Runner_Three_Sub1_Sub2) GetBundleName() string {
	return "ydktest"
}

func (sub2 *Runner_Three_Sub1_Sub2) GetYangName() string {
	return "sub2"
}

func (sub2 *Runner_Three_Sub1_Sub2) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (sub2 *Runner_Three_Sub1_Sub2) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (sub2 *Runner_Three_Sub1_Sub2) SetParent(parent types.Entity) {
	sub2.parent = parent
}

func (sub2 *Runner_Three_Sub1_Sub2) GetParent() types.Entity {
	return sub2.parent
}

func (sub2 *Runner_Three_Sub1_Sub2) GetParentYangName() string {
	return "sub1"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Ytypes
//////////////////////////////////////////////////////////////////////////
type Runner_Ytypes struct {
	parent types.Entity
	Filter types.YFilter

	Enabled  interface{} // empty
	BuiltInT Runner_Ytypes_BuiltInT
	DerivedT Runner_Ytypes_DerivedT
	None     Runner_Ytypes_None
}

func (ytypes *Runner_Ytypes) HasDataOrFilter() bool {
	return ytypes.Enabled != nil ||
		ytypes.BuiltInT.HasDataOrFilter() ||
		ytypes.DerivedT.HasDataOrFilter() ||
		ytypes.None.HasDataOrFilter()
}

func (ytypes *Runner_Ytypes) GetFilter() types.YFilter {
	return ytypes.Filter
}

func (ytypes *Runner_Ytypes) GetSegmentPath() string {
	return "ytypes"
}

func (ytypes *Runner_Ytypes) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: ytypes.GetSegmentPath()}
	var leafData types.LeafData

	if ytypes.Enabled != nil {
		switch ytypes.Enabled.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ytypes.Enabled.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ytypes.Enabled)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "enabled", Data: leafData})
		}
	}

	return entityPath
}

func (ytypes *Runner_Ytypes) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "built-in-t" {
		return &ytypes.BuiltInT
	}
	if child_yang_name == "derived-t" {
		return &ytypes.DerivedT
	}
	if child_yang_name == "none" {
		return &ytypes.None
	}
	return nil
}

func (ytypes *Runner_Ytypes) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["built-in-t"] = &ytypes.BuiltInT
	children["derived-t"] = &ytypes.DerivedT
	children["none"] = &ytypes.None

	return children
}

func (ytypes *Runner_Ytypes) SetValue(value_path string, value string) {
	if value_path == "enabled" {
		ytypes.Enabled = value
	}
}

func (ytypes *Runner_Ytypes) GetBundleName() string {
	return "ydktest"
}

func (ytypes *Runner_Ytypes) GetYangName() string {
	return "ytypes"
}

func (ytypes *Runner_Ytypes) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (ytypes *Runner_Ytypes) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (ytypes *Runner_Ytypes) SetParent(parent types.Entity) {
	ytypes.parent = parent
}

func (ytypes *Runner_Ytypes) GetParent() types.Entity {
	return ytypes.parent
}

func (ytypes *Runner_Ytypes) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Ytypes_None
//////////////////////////////////////////////////////////////////////////
type Runner_Ytypes_None struct {
	parent types.Entity
	Filter types.YFilter

	Test interface{} // str
}

func (none *Runner_Ytypes_None) HasDataOrFilter() bool {
	return none.Test != nil
}

func (none *Runner_Ytypes_None) GetFilter() types.YFilter {
	return none.Filter
}

func (none *Runner_Ytypes_None) GetSegmentPath() string {
	return "none"
}

func (none *Runner_Ytypes_None) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: none.GetSegmentPath()}
	var leafData types.LeafData

	if none.Test != nil {
		switch none.Test.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: none.Test.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", none.Test)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "test", Data: leafData})
		}
	}

	return entityPath
}

func (none *Runner_Ytypes_None) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (none *Runner_Ytypes_None) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (none *Runner_Ytypes_None) SetValue(value_path string, value string) {
	if value_path == "test" {
		none.Test = value
	}
}

func (none *Runner_Ytypes_None) GetBundleName() string {
	return "ydktest"
}

func (none *Runner_Ytypes_None) GetYangName() string {
	return "none"
}

func (none *Runner_Ytypes_None) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (none *Runner_Ytypes_None) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (none *Runner_Ytypes_None) SetParent(parent types.Entity) {
	none.parent = parent
}

func (none *Runner_Ytypes_None) GetParent() types.Entity {
	return none.parent
}

func (none *Runner_Ytypes_None) GetParentYangName() string {
	return "ytypes"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Ytypes_BuiltInT
//////////////////////////////////////////////////////////////////////////
type Runner_Ytypes_BuiltInT struct {
	parent types.Entity
	Filter types.YFilter

	Bincoded         interface{} // str
	BitsLlist        []interface{}
	BitsValue        interface{} // bits
	BoolValue        interface{} // boolean
	Deci64           interface{} // str
	EmbededEnum      interface{} // enumeration
	Emptee           interface{} // empty
	EnumIntValue     interface{} // str
	EnumLlist        []interface{}
	EnumValue        interface{} // enumeration
	IdentityLlist    []interface{}
	IdentityRefValue interface{} // identityref
	LeafRef          interface{} // str
	Llstring         []interface{}
	Llunion          []interface{}
	Name             interface{} // str
	Number16         interface{} // int16
	Number32         interface{} // int32
	Number64         interface{} // int64
	Number8          interface{} // int8
	Status           interface{} // enumeration
	U_Number16       interface{} // uint16
	U_Number32       interface{} // uint32
	U_Number64       interface{} // uint64
	U_Number8        interface{} // uint8
	Younion          interface{} // str
	YounionList      []interface{}
	YounionRecursive interface{} // str
}

func (builtInT *Runner_Ytypes_BuiltInT) HasDataOrFilter() bool {
	for _, leaf := range builtInT.BitsLlist {
		if leaf != nil {
			return true
		}
	}
	for _, leaf := range builtInT.EnumLlist {
		if leaf != nil {
			return true
		}
	}
	for _, leaf := range builtInT.IdentityLlist {
		if leaf != nil {
			return true
		}
	}
	for _, leaf := range builtInT.Llstring {
		if leaf != nil {
			return true
		}
	}
	for _, leaf := range builtInT.Llunion {
		if leaf != nil {
			return true
		}
	}
	for _, leaf := range builtInT.YounionList {
		if leaf != nil {
			return true
		}
	}
	return builtInT.Bincoded != nil ||
		builtInT.BitsValue != nil ||
		builtInT.BoolValue != nil ||
		builtInT.Deci64 != nil ||
		builtInT.EmbededEnum != nil ||
		builtInT.Emptee != nil ||
		builtInT.EnumIntValue != nil ||
		builtInT.EnumValue != nil ||
		builtInT.IdentityRefValue != nil ||
		builtInT.LeafRef != nil ||
		builtInT.Name != nil ||
		builtInT.Number16 != nil ||
		builtInT.Number32 != nil ||
		builtInT.Number64 != nil ||
		builtInT.Number8 != nil ||
		builtInT.Status != nil ||
		builtInT.U_Number16 != nil ||
		builtInT.U_Number32 != nil ||
		builtInT.U_Number64 != nil ||
		builtInT.U_Number8 != nil ||
		builtInT.Younion != nil ||
		builtInT.YounionRecursive != nil
}

func (builtInT *Runner_Ytypes_BuiltInT) GetFilter() types.YFilter {
	return builtInT.Filter
}

func (builtInT *Runner_Ytypes_BuiltInT) GetSegmentPath() string {
	return "built-in-t"
}

func (builtInT *Runner_Ytypes_BuiltInT) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: builtInT.GetSegmentPath()}
	var leafData types.LeafData

	if builtInT.Bincoded != nil {
		switch builtInT.Bincoded.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Bincoded.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Bincoded)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "bincoded", Data: leafData})
		}
	}

	if builtInT.BitsValue != nil {
		switch builtInT.BitsValue.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.BitsValue.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.BitsValue)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "bits-value", Data: leafData})
		}
	}

	if builtInT.BoolValue != nil {
		switch builtInT.BoolValue.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.BoolValue.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.BoolValue)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "bool-value", Data: leafData})
		}
	}

	if builtInT.Deci64 != nil {
		switch builtInT.Deci64.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Deci64.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Deci64)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "deci64", Data: leafData})
		}
	}

	if builtInT.EmbededEnum != nil {
		switch builtInT.EmbededEnum.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.EmbededEnum.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.EmbededEnum)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "embeded-enum", Data: leafData})
		}
	}

	if builtInT.Emptee != nil {
		switch builtInT.Emptee.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Emptee.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Emptee)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "emptee", Data: leafData})
		}
	}

	if builtInT.EnumIntValue != nil {
		switch builtInT.EnumIntValue.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.EnumIntValue.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.EnumIntValue)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "enum-int-value", Data: leafData})
		}
	}

	if builtInT.EnumValue != nil {
		switch builtInT.EnumValue.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.EnumValue.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.EnumValue)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "enum-value", Data: leafData})
		}
	}

	if builtInT.IdentityRefValue != nil {
		switch builtInT.IdentityRefValue.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.IdentityRefValue.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.IdentityRefValue)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "identity-ref-value", Data: leafData})
		}
	}

	if builtInT.LeafRef != nil {
		switch builtInT.LeafRef.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.LeafRef.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.LeafRef)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "leaf-ref", Data: leafData})
		}
	}

	if builtInT.Name != nil {
		switch builtInT.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	if builtInT.Number16 != nil {
		switch builtInT.Number16.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Number16.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Number16)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number16", Data: leafData})
		}
	}

	if builtInT.Number32 != nil {
		switch builtInT.Number32.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Number32.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Number32)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number32", Data: leafData})
		}
	}

	if builtInT.Number64 != nil {
		switch builtInT.Number64.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Number64.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Number64)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number64", Data: leafData})
		}
	}

	if builtInT.Number8 != nil {
		switch builtInT.Number8.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Number8.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Number8)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number8", Data: leafData})
		}
	}

	if builtInT.Status != nil {
		switch builtInT.Status.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Status.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Status)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "status", Data: leafData})
		}
	}

	if builtInT.U_Number16 != nil {
		switch builtInT.U_Number16.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.U_Number16.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.U_Number16)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "u_number16", Data: leafData})
		}
	}

	if builtInT.U_Number32 != nil {
		switch builtInT.U_Number32.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.U_Number32.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.U_Number32)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "u_number32", Data: leafData})
		}
	}

	if builtInT.U_Number64 != nil {
		switch builtInT.U_Number64.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.U_Number64.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.U_Number64)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "u_number64", Data: leafData})
		}
	}

	if builtInT.U_Number8 != nil {
		switch builtInT.U_Number8.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.U_Number8.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.U_Number8)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "u_number8", Data: leafData})
		}
	}

	if builtInT.Younion != nil {
		switch builtInT.Younion.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.Younion.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.Younion)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "younion", Data: leafData})
		}
	}

	if builtInT.YounionRecursive != nil {
		switch builtInT.YounionRecursive.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: builtInT.YounionRecursive.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", builtInT.YounionRecursive)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "younion-recursive", Data: leafData})
		}
	}

	return entityPath
}

func (builtInT *Runner_Ytypes_BuiltInT) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (builtInT *Runner_Ytypes_BuiltInT) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (builtInT *Runner_Ytypes_BuiltInT) SetValue(value_path string, value string) {
	if value_path == "bincoded" {
		builtInT.Bincoded = value
	}
	if value_path == "bits-llist" {
		builtInT.BitsLlist = append(builtInT.BitsLlist, value)
	}
	if value_path == "bits-value" {
		builtInT.BitsValue = value
	}
	if value_path == "bool-value" {
		builtInT.BoolValue = value
	}
	if value_path == "deci64" {
		builtInT.Deci64 = value
	}
	if value_path == "embeded-enum" {
		builtInT.EmbededEnum = value
	}
	if value_path == "emptee" {
		builtInT.Emptee = value
	}
	if value_path == "enum-int-value" {
		builtInT.EnumIntValue = value
	}
	if value_path == "enum-llist" {
		builtInT.EnumLlist = append(builtInT.EnumLlist, value)
	}
	if value_path == "enum-value" {
		builtInT.EnumValue = value
	}
	if value_path == "identity-llist" {
		builtInT.IdentityLlist = append(builtInT.IdentityLlist, value)
	}
	if value_path == "identity-ref-value" {
		builtInT.IdentityRefValue = value
	}
	if value_path == "leaf-ref" {
		builtInT.LeafRef = value
	}
	if value_path == "llstring" {
		builtInT.Llstring = append(builtInT.Llstring, value)
	}
	if value_path == "llunion" {
		builtInT.Llunion = append(builtInT.Llunion, value)
	}
	if value_path == "name" {
		builtInT.Name = value
	}
	if value_path == "number16" {
		builtInT.Number16 = value
	}
	if value_path == "number32" {
		builtInT.Number32 = value
	}
	if value_path == "number64" {
		builtInT.Number64 = value
	}
	if value_path == "number8" {
		builtInT.Number8 = value
	}
	if value_path == "status" {
		builtInT.Status = value
	}
	if value_path == "u_number16" {
		builtInT.U_Number16 = value
	}
	if value_path == "u_number32" {
		builtInT.U_Number32 = value
	}
	if value_path == "u_number64" {
		builtInT.U_Number64 = value
	}
	if value_path == "u_number8" {
		builtInT.U_Number8 = value
	}
	if value_path == "younion" {
		builtInT.Younion = value
	}
	if value_path == "younion-list" {
		builtInT.YounionList = append(builtInT.YounionList, value)
	}
	if value_path == "younion-recursive" {
		builtInT.YounionRecursive = value
	}
}

func (builtInT *Runner_Ytypes_BuiltInT) GetBundleName() string {
	return "ydktest"
}

func (builtInT *Runner_Ytypes_BuiltInT) GetYangName() string {
	return "builtInT"
}

func (builtInT *Runner_Ytypes_BuiltInT) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (builtInT *Runner_Ytypes_BuiltInT) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (builtInT *Runner_Ytypes_BuiltInT) SetParent(parent types.Entity) {
	builtInT.parent = parent
}

func (builtInT *Runner_Ytypes_BuiltInT) GetParent() types.Entity {
	return builtInT.parent
}

func (builtInT *Runner_Ytypes_BuiltInT) GetParentYangName() string {
	return "ytypes"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Ytypes_DerivedT
//////////////////////////////////////////////////////////////////////////
type Runner_Ytypes_DerivedT struct {
	parent types.Entity
	Filter types.YFilter
}

func (derivedT *Runner_Ytypes_DerivedT) HasDataOrFilter() bool {
	return false
}

func (derivedT *Runner_Ytypes_DerivedT) GetFilter() types.YFilter {
	return derivedT.Filter
}

func (derivedT *Runner_Ytypes_DerivedT) GetSegmentPath() string {
	return "derived-t"
}

func (derivedT *Runner_Ytypes_DerivedT) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: derivedT.GetSegmentPath()}
	return entityPath
}

func (derivedT *Runner_Ytypes_DerivedT) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (derivedT *Runner_Ytypes_DerivedT) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (derivedT *Runner_Ytypes_DerivedT) SetValue(value_path string, value string) {
}

func (derivedT *Runner_Ytypes_DerivedT) GetBundleName() string {
	return "ydktest"
}

func (derivedT *Runner_Ytypes_DerivedT) GetYangName() string {
	return "derivedT"
}

func (derivedT *Runner_Ytypes_DerivedT) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (derivedT *Runner_Ytypes_DerivedT) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (derivedT *Runner_Ytypes_DerivedT) SetParent(parent types.Entity) {
	derivedT.parent = parent
}

func (derivedT *Runner_Ytypes_DerivedT) GetParent() types.Entity {
	return derivedT.parent
}

func (derivedT *Runner_Ytypes_DerivedT) GetParentYangName() string {
	return "ytypes"
}

//////////////////////////////////////////////////////////////////////////
// Runner_OneList
//////////////////////////////////////////////////////////////////////////
type Runner_OneList struct {
	parent types.Entity
	Filter types.YFilter

	IdentityList []Runner_OneList_IdentityList
	Ldata        []Runner_OneList_Ldata
	OneAugList   Runner_OneList_OneAugList
}

func (oneList *Runner_OneList) HasDataOrFilter() bool {
	for _, child := range oneList.IdentityList {
		if child.HasDataOrFilter() {
			return true
		}
	}
	for _, child := range oneList.Ldata {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return oneList.OneAugList.HasDataOrFilter()
}

func (oneList *Runner_OneList) GetFilter() types.YFilter {
	return oneList.Filter
}

func (oneList *Runner_OneList) GetSegmentPath() string {
	return "one-list"
}

func (oneList *Runner_OneList) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: oneList.GetSegmentPath()}
	return entityPath
}

func (oneList *Runner_OneList) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "identity-list" {
		for _, c := range oneList.IdentityList {
			if oneList.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_OneList_IdentityList{}
		oneList.IdentityList = append(oneList.IdentityList, child)
		return &oneList.IdentityList[len(oneList.IdentityList)-1]
	}
	if child_yang_name == "ldata" {
		for _, c := range oneList.Ldata {
			if oneList.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_OneList_Ldata{}
		oneList.Ldata = append(oneList.Ldata, child)
		return &oneList.Ldata[len(oneList.Ldata)-1]
	}
	if child_yang_name == "one-aug-list" {
		return &oneList.OneAugList
	}
	return nil
}

func (oneList *Runner_OneList) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range oneList.IdentityList {
		children[oneList.IdentityList[i].GetSegmentPath()] = &oneList.IdentityList[i]
	}
	for i := range oneList.Ldata {
		children[oneList.Ldata[i].GetSegmentPath()] = &oneList.Ldata[i]
	}
	children["one-aug-list"] = &oneList.OneAugList

	return children
}

func (oneList *Runner_OneList) SetValue(value_path string, value string) {
}

func (oneList *Runner_OneList) GetBundleName() string {
	return "ydktest"
}

func (oneList *Runner_OneList) GetYangName() string {
	return "oneList"
}

func (oneList *Runner_OneList) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (oneList *Runner_OneList) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (oneList *Runner_OneList) SetParent(parent types.Entity) {
	oneList.parent = parent
}

func (oneList *Runner_OneList) GetParent() types.Entity {
	return oneList.parent
}

func (oneList *Runner_OneList) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_OneList_Ldata
//////////////////////////////////////////////////////////////////////////
type Runner_OneList_Ldata struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
}

func (ldata *Runner_OneList_Ldata) HasDataOrFilter() bool {
	return ldata.Number != nil ||
		ldata.Name != nil
}

func (ldata *Runner_OneList_Ldata) GetFilter() types.YFilter {
	return ldata.Filter
}

func (ldata *Runner_OneList_Ldata) GetSegmentPath() string {
	return "ldata" + "[number=\"" + fmt.Sprintf("%v", ldata.Number) + "\"]"
}

func (ldata *Runner_OneList_Ldata) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: ldata.GetSegmentPath()}
	var leafData types.LeafData

	if ldata.Number != nil {
		switch ldata.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if ldata.Name != nil {
		switch ldata.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (ldata *Runner_OneList_Ldata) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (ldata *Runner_OneList_Ldata) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (ldata *Runner_OneList_Ldata) SetValue(value_path string, value string) {
	if value_path == "number" {
		ldata.Number = value
	}
	if value_path == "name" {
		ldata.Name = value
	}
}

func (ldata *Runner_OneList_Ldata) GetBundleName() string {
	return "ydktest"
}

func (ldata *Runner_OneList_Ldata) GetYangName() string {
	return "ldata"
}

func (ldata *Runner_OneList_Ldata) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (ldata *Runner_OneList_Ldata) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (ldata *Runner_OneList_Ldata) SetParent(parent types.Entity) {
	ldata.parent = parent
}

func (ldata *Runner_OneList_Ldata) GetParent() types.Entity {
	return ldata.parent
}

func (ldata *Runner_OneList_Ldata) GetParentYangName() string {
	return "one-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_OneList_IdentityList
//////////////////////////////////////////////////////////////////////////
type Runner_OneList_IdentityList struct {
	parent types.Entity
	Filter types.YFilter

	IdRef  interface{} // identityref
	Config Runner_OneList_IdentityList_Config
}

func (identityList *Runner_OneList_IdentityList) HasDataOrFilter() bool {
	return identityList.IdRef != nil ||
		identityList.Config.HasDataOrFilter()
}

func (identityList *Runner_OneList_IdentityList) GetFilter() types.YFilter {
	return identityList.Filter
}

func (identityList *Runner_OneList_IdentityList) GetSegmentPath() string {
	return "identity-list" + "[id-ref=\"" + fmt.Sprintf("%v", identityList.IdRef) + "\"]"
}

func (identityList *Runner_OneList_IdentityList) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: identityList.GetSegmentPath()}
	var leafData types.LeafData

	if identityList.IdRef != nil {
		switch identityList.IdRef.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: identityList.IdRef.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", identityList.IdRef)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "id-ref", Data: leafData})
		}
	}

	return entityPath
}

func (identityList *Runner_OneList_IdentityList) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "config" {
		return &identityList.Config
	}
	return nil
}

func (identityList *Runner_OneList_IdentityList) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["config"] = &identityList.Config

	return children
}

func (identityList *Runner_OneList_IdentityList) SetValue(value_path string, value string) {
	if value_path == "id-ref" {
		identityList.IdRef = value
	}
}

func (identityList *Runner_OneList_IdentityList) GetBundleName() string {
	return "ydktest"
}

func (identityList *Runner_OneList_IdentityList) GetYangName() string {
	return "identityList"
}

func (identityList *Runner_OneList_IdentityList) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (identityList *Runner_OneList_IdentityList) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (identityList *Runner_OneList_IdentityList) SetParent(parent types.Entity) {
	identityList.parent = parent
}

func (identityList *Runner_OneList_IdentityList) GetParent() types.Entity {
	return identityList.parent
}

func (identityList *Runner_OneList_IdentityList) GetParentYangName() string {
	return "one-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_OneList_IdentityList_Config
//////////////////////////////////////////////////////////////////////////
type Runner_OneList_IdentityList_Config struct {
	parent types.Entity
	Filter types.YFilter

	Id interface{} // identityref
}

func (config *Runner_OneList_IdentityList_Config) HasDataOrFilter() bool {
	return config.Id != nil
}

func (config *Runner_OneList_IdentityList_Config) GetFilter() types.YFilter {
	return config.Filter
}

func (config *Runner_OneList_IdentityList_Config) GetSegmentPath() string {
	return "config"
}

func (config *Runner_OneList_IdentityList_Config) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: config.GetSegmentPath()}
	var leafData types.LeafData

	if config.Id != nil {
		switch config.Id.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: config.Id.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", config.Id)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "id", Data: leafData})
		}
	}

	return entityPath
}

func (config *Runner_OneList_IdentityList_Config) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (config *Runner_OneList_IdentityList_Config) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (config *Runner_OneList_IdentityList_Config) SetValue(value_path string, value string) {
	if value_path == "id" {
		config.Id = value
	}
}

func (config *Runner_OneList_IdentityList_Config) GetBundleName() string {
	return "ydktest"
}

func (config *Runner_OneList_IdentityList_Config) GetYangName() string {
	return "config"
}

func (config *Runner_OneList_IdentityList_Config) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (config *Runner_OneList_IdentityList_Config) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (config *Runner_OneList_IdentityList_Config) SetParent(parent types.Entity) {
	config.parent = parent
}

func (config *Runner_OneList_IdentityList_Config) GetParent() types.Entity {
	return config.parent
}

func (config *Runner_OneList_IdentityList_Config) GetParentYangName() string {
	return "identity-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_OneList_OneAugList
//////////////////////////////////////////////////////////////////////////
type Runner_OneList_OneAugList struct {
	parent types.Entity
	Filter types.YFilter

	Enabled interface{} // boolean
	Ldata   []Runner_OneList_OneAugList_Ldata
}

func (oneAugList *Runner_OneList_OneAugList) HasDataOrFilter() bool {
	for _, child := range oneAugList.Ldata {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return oneAugList.Enabled != nil
}

func (oneAugList *Runner_OneList_OneAugList) GetFilter() types.YFilter {
	return oneAugList.Filter
}

func (oneAugList *Runner_OneList_OneAugList) GetSegmentPath() string {
	return "ydktest-sanity-augm:one-aug-list"
}

func (oneAugList *Runner_OneList_OneAugList) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: oneAugList.GetSegmentPath()}
	var leafData types.LeafData

	if oneAugList.Enabled != nil {
		switch oneAugList.Enabled.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: oneAugList.Enabled.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", oneAugList.Enabled)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "enabled", Data: leafData})
		}
	}

	return entityPath
}

func (oneAugList *Runner_OneList_OneAugList) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "ldata" {
		for _, c := range oneAugList.Ldata {
			if oneAugList.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_OneList_OneAugList_Ldata{}
		oneAugList.Ldata = append(oneAugList.Ldata, child)
		return &oneAugList.Ldata[len(oneAugList.Ldata)-1]
	}
	return nil
}

func (oneAugList *Runner_OneList_OneAugList) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range oneAugList.Ldata {
		children[oneAugList.Ldata[i].GetSegmentPath()] = &oneAugList.Ldata[i]
	}

	return children
}

func (oneAugList *Runner_OneList_OneAugList) SetValue(value_path string, value string) {
	if value_path == "enabled" {
		oneAugList.Enabled = value
	}
}

func (oneAugList *Runner_OneList_OneAugList) GetBundleName() string {
	return "ydktest"
}

func (oneAugList *Runner_OneList_OneAugList) GetYangName() string {
	return "oneAugList"
}

func (oneAugList *Runner_OneList_OneAugList) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (oneAugList *Runner_OneList_OneAugList) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (oneAugList *Runner_OneList_OneAugList) SetParent(parent types.Entity) {
	oneAugList.parent = parent
}

func (oneAugList *Runner_OneList_OneAugList) GetParent() types.Entity {
	return oneAugList.parent
}

func (oneAugList *Runner_OneList_OneAugList) GetParentYangName() string {
	return "one-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_OneList_OneAugList_Ldata
//////////////////////////////////////////////////////////////////////////
type Runner_OneList_OneAugList_Ldata struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
}

func (ldata *Runner_OneList_OneAugList_Ldata) HasDataOrFilter() bool {
	return ldata.Number != nil ||
		ldata.Name != nil
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetFilter() types.YFilter {
	return ldata.Filter
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetSegmentPath() string {
	return "ldata" + "[number=\"" + fmt.Sprintf("%v", ldata.Number) + "\"]"
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: ldata.GetSegmentPath()}
	var leafData types.LeafData

	if ldata.Number != nil {
		switch ldata.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if ldata.Name != nil {
		switch ldata.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (ldata *Runner_OneList_OneAugList_Ldata) SetValue(value_path string, value string) {
	if value_path == "number" {
		ldata.Number = value
	}
	if value_path == "name" {
		ldata.Name = value
	}
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetBundleName() string {
	return "ydktest"
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetYangName() string {
	return "ldata"
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (ldata *Runner_OneList_OneAugList_Ldata) SetParent(parent types.Entity) {
	ldata.parent = parent
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetParent() types.Entity {
	return ldata.parent
}

func (ldata *Runner_OneList_OneAugList_Ldata) GetParentYangName() string {
	return "one-aug-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_TwoList
//////////////////////////////////////////////////////////////////////////
type Runner_TwoList struct {
	parent types.Entity
	Filter types.YFilter

	Ldata []Runner_TwoList_Ldata
}

func (twoList *Runner_TwoList) HasDataOrFilter() bool {
	for _, child := range twoList.Ldata {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return false
}

func (twoList *Runner_TwoList) GetFilter() types.YFilter {
	return twoList.Filter
}

func (twoList *Runner_TwoList) GetSegmentPath() string {
	return "two-list"
}

func (twoList *Runner_TwoList) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: twoList.GetSegmentPath()}
	return entityPath
}

func (twoList *Runner_TwoList) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "ldata" {
		for _, c := range twoList.Ldata {
			if twoList.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_TwoList_Ldata{}
		twoList.Ldata = append(twoList.Ldata, child)
		return &twoList.Ldata[len(twoList.Ldata)-1]
	}
	return nil
}

func (twoList *Runner_TwoList) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range twoList.Ldata {
		children[twoList.Ldata[i].GetSegmentPath()] = &twoList.Ldata[i]
	}

	return children
}

func (twoList *Runner_TwoList) SetValue(value_path string, value string) {
}

func (twoList *Runner_TwoList) GetBundleName() string {
	return "ydktest"
}

func (twoList *Runner_TwoList) GetYangName() string {
	return "twoList"
}

func (twoList *Runner_TwoList) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (twoList *Runner_TwoList) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (twoList *Runner_TwoList) SetParent(parent types.Entity) {
	twoList.parent = parent
}

func (twoList *Runner_TwoList) GetParent() types.Entity {
	return twoList.parent
}

func (twoList *Runner_TwoList) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_TwoList_Ldata
//////////////////////////////////////////////////////////////////////////
type Runner_TwoList_Ldata struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
	Subl1  []Runner_TwoList_Ldata_Subl1
}

func (ldata *Runner_TwoList_Ldata) HasDataOrFilter() bool {
	for _, child := range ldata.Subl1 {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return ldata.Number != nil ||
		ldata.Name != nil
}

func (ldata *Runner_TwoList_Ldata) GetFilter() types.YFilter {
	return ldata.Filter
}

func (ldata *Runner_TwoList_Ldata) GetSegmentPath() string {
	return "ldata" + "[number=\"" + fmt.Sprintf("%v", ldata.Number) + "\"]"
}

func (ldata *Runner_TwoList_Ldata) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: ldata.GetSegmentPath()}
	var leafData types.LeafData

	if ldata.Number != nil {
		switch ldata.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if ldata.Name != nil {
		switch ldata.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (ldata *Runner_TwoList_Ldata) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "subl1" {
		for _, c := range ldata.Subl1 {
			if ldata.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_TwoList_Ldata_Subl1{}
		ldata.Subl1 = append(ldata.Subl1, child)
		return &ldata.Subl1[len(ldata.Subl1)-1]
	}
	return nil
}

func (ldata *Runner_TwoList_Ldata) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range ldata.Subl1 {
		children[ldata.Subl1[i].GetSegmentPath()] = &ldata.Subl1[i]
	}

	return children
}

func (ldata *Runner_TwoList_Ldata) SetValue(value_path string, value string) {
	if value_path == "number" {
		ldata.Number = value
	}
	if value_path == "name" {
		ldata.Name = value
	}
}

func (ldata *Runner_TwoList_Ldata) GetBundleName() string {
	return "ydktest"
}

func (ldata *Runner_TwoList_Ldata) GetYangName() string {
	return "ldata"
}

func (ldata *Runner_TwoList_Ldata) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (ldata *Runner_TwoList_Ldata) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (ldata *Runner_TwoList_Ldata) SetParent(parent types.Entity) {
	ldata.parent = parent
}

func (ldata *Runner_TwoList_Ldata) GetParent() types.Entity {
	return ldata.parent
}

func (ldata *Runner_TwoList_Ldata) GetParentYangName() string {
	return "two-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_TwoList_Ldata_Subl1
//////////////////////////////////////////////////////////////////////////
type Runner_TwoList_Ldata_Subl1 struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
}

func (subl1 *Runner_TwoList_Ldata_Subl1) HasDataOrFilter() bool {
	return subl1.Number != nil ||
		subl1.Name != nil
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetFilter() types.YFilter {
	return subl1.Filter
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetSegmentPath() string {
	return "subl1" + "[number=\"" + fmt.Sprintf("%v", subl1.Number) + "\"]"
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: subl1.GetSegmentPath()}
	var leafData types.LeafData

	if subl1.Number != nil {
		switch subl1.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subl1.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subl1.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if subl1.Name != nil {
		switch subl1.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subl1.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subl1.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (subl1 *Runner_TwoList_Ldata_Subl1) SetValue(value_path string, value string) {
	if value_path == "number" {
		subl1.Number = value
	}
	if value_path == "name" {
		subl1.Name = value
	}
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetBundleName() string {
	return "ydktest"
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetYangName() string {
	return "subl1"
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (subl1 *Runner_TwoList_Ldata_Subl1) SetParent(parent types.Entity) {
	subl1.parent = parent
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetParent() types.Entity {
	return subl1.parent
}

func (subl1 *Runner_TwoList_Ldata_Subl1) GetParentYangName() string {
	return "ldata"
}

//////////////////////////////////////////////////////////////////////////
// Runner_ThreeList
//////////////////////////////////////////////////////////////////////////
type Runner_ThreeList struct {
	parent types.Entity
	Filter types.YFilter

	Ldata []Runner_ThreeList_Ldata
}

func (threeList *Runner_ThreeList) HasDataOrFilter() bool {
	for _, child := range threeList.Ldata {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return false
}

func (threeList *Runner_ThreeList) GetFilter() types.YFilter {
	return threeList.Filter
}

func (threeList *Runner_ThreeList) GetSegmentPath() string {
	return "three-list"
}

func (threeList *Runner_ThreeList) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: threeList.GetSegmentPath()}
	return entityPath
}

func (threeList *Runner_ThreeList) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "ldata" {
		for _, c := range threeList.Ldata {
			if threeList.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_ThreeList_Ldata{}
		threeList.Ldata = append(threeList.Ldata, child)
		return &threeList.Ldata[len(threeList.Ldata)-1]
	}
	return nil
}

func (threeList *Runner_ThreeList) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range threeList.Ldata {
		children[threeList.Ldata[i].GetSegmentPath()] = &threeList.Ldata[i]
	}

	return children
}

func (threeList *Runner_ThreeList) SetValue(value_path string, value string) {
}

func (threeList *Runner_ThreeList) GetBundleName() string {
	return "ydktest"
}

func (threeList *Runner_ThreeList) GetYangName() string {
	return "threeList"
}

func (threeList *Runner_ThreeList) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (threeList *Runner_ThreeList) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (threeList *Runner_ThreeList) SetParent(parent types.Entity) {
	threeList.parent = parent
}

func (threeList *Runner_ThreeList) GetParent() types.Entity {
	return threeList.parent
}

func (threeList *Runner_ThreeList) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_ThreeList_Ldata
//////////////////////////////////////////////////////////////////////////
type Runner_ThreeList_Ldata struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
	Subl1  []Runner_ThreeList_Ldata_Subl1
}

func (ldata *Runner_ThreeList_Ldata) HasDataOrFilter() bool {
	for _, child := range ldata.Subl1 {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return ldata.Number != nil ||
		ldata.Name != nil
}

func (ldata *Runner_ThreeList_Ldata) GetFilter() types.YFilter {
	return ldata.Filter
}

func (ldata *Runner_ThreeList_Ldata) GetSegmentPath() string {
	return "ldata" + "[number=\"" + fmt.Sprintf("%v", ldata.Number) + "\"]"
}

func (ldata *Runner_ThreeList_Ldata) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: ldata.GetSegmentPath()}
	var leafData types.LeafData

	if ldata.Number != nil {
		switch ldata.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if ldata.Name != nil {
		switch ldata.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (ldata *Runner_ThreeList_Ldata) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "subl1" {
		for _, c := range ldata.Subl1 {
			if ldata.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_ThreeList_Ldata_Subl1{}
		ldata.Subl1 = append(ldata.Subl1, child)
		return &ldata.Subl1[len(ldata.Subl1)-1]
	}
	return nil
}

func (ldata *Runner_ThreeList_Ldata) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range ldata.Subl1 {
		children[ldata.Subl1[i].GetSegmentPath()] = &ldata.Subl1[i]
	}

	return children
}

func (ldata *Runner_ThreeList_Ldata) SetValue(value_path string, value string) {
	if value_path == "number" {
		ldata.Number = value
	}
	if value_path == "name" {
		ldata.Name = value
	}
}

func (ldata *Runner_ThreeList_Ldata) GetBundleName() string {
	return "ydktest"
}

func (ldata *Runner_ThreeList_Ldata) GetYangName() string {
	return "ldata"
}

func (ldata *Runner_ThreeList_Ldata) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (ldata *Runner_ThreeList_Ldata) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (ldata *Runner_ThreeList_Ldata) SetParent(parent types.Entity) {
	ldata.parent = parent
}

func (ldata *Runner_ThreeList_Ldata) GetParent() types.Entity {
	return ldata.parent
}

func (ldata *Runner_ThreeList_Ldata) GetParentYangName() string {
	return "three-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_ThreeList_Ldata_Subl1
//////////////////////////////////////////////////////////////////////////
type Runner_ThreeList_Ldata_Subl1 struct {
	parent types.Entity
	Filter types.YFilter

	Number   interface{} // int32
	Name     interface{} // str
	SubSubl1 []Runner_ThreeList_Ldata_Subl1_SubSubl1
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) HasDataOrFilter() bool {
	for _, child := range subl1.SubSubl1 {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return subl1.Number != nil ||
		subl1.Name != nil
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetFilter() types.YFilter {
	return subl1.Filter
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetSegmentPath() string {
	return "subl1" + "[number=\"" + fmt.Sprintf("%v", subl1.Number) + "\"]"
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: subl1.GetSegmentPath()}
	var leafData types.LeafData

	if subl1.Number != nil {
		switch subl1.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subl1.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subl1.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if subl1.Name != nil {
		switch subl1.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subl1.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subl1.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "sub-subl1" {
		for _, c := range subl1.SubSubl1 {
			if subl1.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_ThreeList_Ldata_Subl1_SubSubl1{}
		subl1.SubSubl1 = append(subl1.SubSubl1, child)
		return &subl1.SubSubl1[len(subl1.SubSubl1)-1]
	}
	return nil
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range subl1.SubSubl1 {
		children[subl1.SubSubl1[i].GetSegmentPath()] = &subl1.SubSubl1[i]
	}

	return children
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) SetValue(value_path string, value string) {
	if value_path == "number" {
		subl1.Number = value
	}
	if value_path == "name" {
		subl1.Name = value
	}
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetBundleName() string {
	return "ydktest"
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetYangName() string {
	return "subl1"
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) SetParent(parent types.Entity) {
	subl1.parent = parent
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetParent() types.Entity {
	return subl1.parent
}

func (subl1 *Runner_ThreeList_Ldata_Subl1) GetParentYangName() string {
	return "ldata"
}

//////////////////////////////////////////////////////////////////////////
// Runner_ThreeList_Ldata_Subl1_SubSubl1
//////////////////////////////////////////////////////////////////////////
type Runner_ThreeList_Ldata_Subl1_SubSubl1 struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) HasDataOrFilter() bool {
	return subSubl1.Number != nil ||
		subSubl1.Name != nil
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetFilter() types.YFilter {
	return subSubl1.Filter
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetSegmentPath() string {
	return "sub-subl1" + "[number=\"" + fmt.Sprintf("%v", subSubl1.Number) + "\"]"
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: subSubl1.GetSegmentPath()}
	var leafData types.LeafData

	if subSubl1.Number != nil {
		switch subSubl1.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subSubl1.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subSubl1.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if subSubl1.Name != nil {
		switch subSubl1.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subSubl1.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subSubl1.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) SetValue(value_path string, value string) {
	if value_path == "number" {
		subSubl1.Number = value
	}
	if value_path == "name" {
		subSubl1.Name = value
	}
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetBundleName() string {
	return "ydktest"
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetYangName() string {
	return "subSubl1"
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) SetParent(parent types.Entity) {
	subSubl1.parent = parent
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetParent() types.Entity {
	return subSubl1.parent
}

func (subSubl1 *Runner_ThreeList_Ldata_Subl1_SubSubl1) GetParentYangName() string {
	return "subl1"
}

//////////////////////////////////////////////////////////////////////////
// Runner_InbtwList
//////////////////////////////////////////////////////////////////////////
type Runner_InbtwList struct {
	parent types.Entity
	Filter types.YFilter

	Ldata []Runner_InbtwList_Ldata
}

func (inbtwList *Runner_InbtwList) HasDataOrFilter() bool {
	for _, child := range inbtwList.Ldata {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return false
}

func (inbtwList *Runner_InbtwList) GetFilter() types.YFilter {
	return inbtwList.Filter
}

func (inbtwList *Runner_InbtwList) GetSegmentPath() string {
	return "inbtw-list"
}

func (inbtwList *Runner_InbtwList) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: inbtwList.GetSegmentPath()}
	return entityPath
}

func (inbtwList *Runner_InbtwList) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "ldata" {
		for _, c := range inbtwList.Ldata {
			if inbtwList.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_InbtwList_Ldata{}
		inbtwList.Ldata = append(inbtwList.Ldata, child)
		return &inbtwList.Ldata[len(inbtwList.Ldata)-1]
	}
	return nil
}

func (inbtwList *Runner_InbtwList) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range inbtwList.Ldata {
		children[inbtwList.Ldata[i].GetSegmentPath()] = &inbtwList.Ldata[i]
	}

	return children
}

func (inbtwList *Runner_InbtwList) SetValue(value_path string, value string) {
}

func (inbtwList *Runner_InbtwList) GetBundleName() string {
	return "ydktest"
}

func (inbtwList *Runner_InbtwList) GetYangName() string {
	return "inbtwList"
}

func (inbtwList *Runner_InbtwList) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (inbtwList *Runner_InbtwList) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (inbtwList *Runner_InbtwList) SetParent(parent types.Entity) {
	inbtwList.parent = parent
}

func (inbtwList *Runner_InbtwList) GetParent() types.Entity {
	return inbtwList.parent
}

func (inbtwList *Runner_InbtwList) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_InbtwList_Ldata
//////////////////////////////////////////////////////////////////////////
type Runner_InbtwList_Ldata struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
	Subc   Runner_InbtwList_Ldata_Subc
}

func (ldata *Runner_InbtwList_Ldata) HasDataOrFilter() bool {
	return ldata.Number != nil ||
		ldata.Name != nil ||
		ldata.Subc.HasDataOrFilter()
}

func (ldata *Runner_InbtwList_Ldata) GetFilter() types.YFilter {
	return ldata.Filter
}

func (ldata *Runner_InbtwList_Ldata) GetSegmentPath() string {
	return "ldata" + "[number=\"" + fmt.Sprintf("%v", ldata.Number) + "\"]"
}

func (ldata *Runner_InbtwList_Ldata) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: ldata.GetSegmentPath()}
	var leafData types.LeafData

	if ldata.Number != nil {
		switch ldata.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if ldata.Name != nil {
		switch ldata.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: ldata.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", ldata.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (ldata *Runner_InbtwList_Ldata) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "subc" {
		return &ldata.Subc
	}
	return nil
}

func (ldata *Runner_InbtwList_Ldata) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["subc"] = &ldata.Subc

	return children
}

func (ldata *Runner_InbtwList_Ldata) SetValue(value_path string, value string) {
	if value_path == "number" {
		ldata.Number = value
	}
	if value_path == "name" {
		ldata.Name = value
	}
}

func (ldata *Runner_InbtwList_Ldata) GetBundleName() string {
	return "ydktest"
}

func (ldata *Runner_InbtwList_Ldata) GetYangName() string {
	return "ldata"
}

func (ldata *Runner_InbtwList_Ldata) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (ldata *Runner_InbtwList_Ldata) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (ldata *Runner_InbtwList_Ldata) SetParent(parent types.Entity) {
	ldata.parent = parent
}

func (ldata *Runner_InbtwList_Ldata) GetParent() types.Entity {
	return ldata.parent
}

func (ldata *Runner_InbtwList_Ldata) GetParentYangName() string {
	return "inbtw-list"
}

//////////////////////////////////////////////////////////////////////////
// Runner_InbtwList_Ldata_Subc
//////////////////////////////////////////////////////////////////////////
type Runner_InbtwList_Ldata_Subc struct {
	parent types.Entity
	Filter types.YFilter

	Name      interface{} // str
	Number    interface{} // int32
	SubcSubl1 []Runner_InbtwList_Ldata_Subc_SubcSubl1
}

func (subc *Runner_InbtwList_Ldata_Subc) HasDataOrFilter() bool {
	for _, child := range subc.SubcSubl1 {
		if child.HasDataOrFilter() {
			return true
		}
	}
	return subc.Name != nil ||
		subc.Number != nil
}

func (subc *Runner_InbtwList_Ldata_Subc) GetFilter() types.YFilter {
	return subc.Filter
}

func (subc *Runner_InbtwList_Ldata_Subc) GetSegmentPath() string {
	return "subc"
}

func (subc *Runner_InbtwList_Ldata_Subc) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: subc.GetSegmentPath()}
	var leafData types.LeafData

	if subc.Name != nil {
		switch subc.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subc.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subc.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	if subc.Number != nil {
		switch subc.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subc.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subc.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (subc *Runner_InbtwList_Ldata_Subc) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "subc-subl1" {
		for _, c := range subc.SubcSubl1 {
			if subc.GetSegmentPath() == segment_path {
				return &c
			}
		}
		child := Runner_InbtwList_Ldata_Subc_SubcSubl1{}
		subc.SubcSubl1 = append(subc.SubcSubl1, child)
		return &subc.SubcSubl1[len(subc.SubcSubl1)-1]
	}
	return nil
}

func (subc *Runner_InbtwList_Ldata_Subc) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	for i := range subc.SubcSubl1 {
		children[subc.SubcSubl1[i].GetSegmentPath()] = &subc.SubcSubl1[i]
	}

	return children
}

func (subc *Runner_InbtwList_Ldata_Subc) SetValue(value_path string, value string) {
	if value_path == "name" {
		subc.Name = value
	}
	if value_path == "number" {
		subc.Number = value
	}
}

func (subc *Runner_InbtwList_Ldata_Subc) GetBundleName() string {
	return "ydktest"
}

func (subc *Runner_InbtwList_Ldata_Subc) GetYangName() string {
	return "subc"
}

func (subc *Runner_InbtwList_Ldata_Subc) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (subc *Runner_InbtwList_Ldata_Subc) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (subc *Runner_InbtwList_Ldata_Subc) SetParent(parent types.Entity) {
	subc.parent = parent
}

func (subc *Runner_InbtwList_Ldata_Subc) GetParent() types.Entity {
	return subc.parent
}

func (subc *Runner_InbtwList_Ldata_Subc) GetParentYangName() string {
	return "ldata"
}

//////////////////////////////////////////////////////////////////////////
// Runner_InbtwList_Ldata_Subc_SubcSubl1
//////////////////////////////////////////////////////////////////////////
type Runner_InbtwList_Ldata_Subc_SubcSubl1 struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
	Name   interface{} // str
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) HasDataOrFilter() bool {
	return subcSubl1.Number != nil ||
		subcSubl1.Name != nil
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetFilter() types.YFilter {
	return subcSubl1.Filter
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetSegmentPath() string {
	return "subc-subl1" + "[number=\"" + fmt.Sprintf("%v", subcSubl1.Number) + "\"]"
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: subcSubl1.GetSegmentPath()}
	var leafData types.LeafData

	if subcSubl1.Number != nil {
		switch subcSubl1.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subcSubl1.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subcSubl1.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	if subcSubl1.Name != nil {
		switch subcSubl1.Name.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: subcSubl1.Name.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", subcSubl1.Name)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name", Data: leafData})
		}
	}

	return entityPath
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) SetValue(value_path string, value string) {
	if value_path == "number" {
		subcSubl1.Number = value
	}
	if value_path == "name" {
		subcSubl1.Name = value
	}
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetBundleName() string {
	return "ydktest"
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetYangName() string {
	return "subcSubl1"
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) SetParent(parent types.Entity) {
	subcSubl1.parent = parent
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetParent() types.Entity {
	return subcSubl1.parent
}

func (subcSubl1 *Runner_InbtwList_Ldata_Subc_SubcSubl1) GetParentYangName() string {
	return "subc"
}

//////////////////////////////////////////////////////////////////////////
// Runner_LeafRef
//////////////////////////////////////////////////////////////////////////
type Runner_LeafRef struct {
	parent types.Entity
	Filter types.YFilter

	RefInbtw               interface{} // str
	RefOneName             interface{} // str
	RefThreeSub1Sub2Number interface{} // str
	RefTwoSub1Number       interface{} // str
	One                    Runner_LeafRef_One
}

func (leafRef *Runner_LeafRef) HasDataOrFilter() bool {
	return leafRef.RefInbtw != nil ||
		leafRef.RefOneName != nil ||
		leafRef.RefThreeSub1Sub2Number != nil ||
		leafRef.RefTwoSub1Number != nil ||
		leafRef.One.HasDataOrFilter()
}

func (leafRef *Runner_LeafRef) GetFilter() types.YFilter {
	return leafRef.Filter
}

func (leafRef *Runner_LeafRef) GetSegmentPath() string {
	return "leaf-ref"
}

func (leafRef *Runner_LeafRef) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: leafRef.GetSegmentPath()}
	var leafData types.LeafData

	if leafRef.RefInbtw != nil {
		switch leafRef.RefInbtw.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: leafRef.RefInbtw.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", leafRef.RefInbtw)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "ref-inbtw", Data: leafData})
		}
	}

	if leafRef.RefOneName != nil {
		switch leafRef.RefOneName.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: leafRef.RefOneName.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", leafRef.RefOneName)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "ref-one-name", Data: leafData})
		}
	}

	if leafRef.RefThreeSub1Sub2Number != nil {
		switch leafRef.RefThreeSub1Sub2Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: leafRef.RefThreeSub1Sub2Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", leafRef.RefThreeSub1Sub2Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "ref-three-sub1-sub2-number", Data: leafData})
		}
	}

	if leafRef.RefTwoSub1Number != nil {
		switch leafRef.RefTwoSub1Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: leafRef.RefTwoSub1Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", leafRef.RefTwoSub1Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "ref-two-sub1-number", Data: leafData})
		}
	}

	return entityPath
}

func (leafRef *Runner_LeafRef) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "one" {
		return &leafRef.One
	}
	return nil
}

func (leafRef *Runner_LeafRef) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["one"] = &leafRef.One

	return children
}

func (leafRef *Runner_LeafRef) SetValue(value_path string, value string) {
	if value_path == "ref-inbtw" {
		leafRef.RefInbtw = value
	}
	if value_path == "ref-one-name" {
		leafRef.RefOneName = value
	}
	if value_path == "ref-three-sub1-sub2-number" {
		leafRef.RefThreeSub1Sub2Number = value
	}
	if value_path == "ref-two-sub1-number" {
		leafRef.RefTwoSub1Number = value
	}
}

func (leafRef *Runner_LeafRef) GetBundleName() string {
	return "ydktest"
}

func (leafRef *Runner_LeafRef) GetYangName() string {
	return "leafRef"
}

func (leafRef *Runner_LeafRef) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (leafRef *Runner_LeafRef) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (leafRef *Runner_LeafRef) SetParent(parent types.Entity) {
	leafRef.parent = parent
}

func (leafRef *Runner_LeafRef) GetParent() types.Entity {
	return leafRef.parent
}

func (leafRef *Runner_LeafRef) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_LeafRef_One
//////////////////////////////////////////////////////////////////////////
type Runner_LeafRef_One struct {
	parent types.Entity
	Filter types.YFilter

	NameOfOne interface{} // str
	Two       Runner_LeafRef_One_Two
}

func (one *Runner_LeafRef_One) HasDataOrFilter() bool {
	return one.NameOfOne != nil ||
		one.Two.HasDataOrFilter()
}

func (one *Runner_LeafRef_One) GetFilter() types.YFilter {
	return one.Filter
}

func (one *Runner_LeafRef_One) GetSegmentPath() string {
	return "one"
}

func (one *Runner_LeafRef_One) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: one.GetSegmentPath()}
	var leafData types.LeafData

	if one.NameOfOne != nil {
		switch one.NameOfOne.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: one.NameOfOne.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", one.NameOfOne)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "name-of-one", Data: leafData})
		}
	}

	return entityPath
}

func (one *Runner_LeafRef_One) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "two" {
		return &one.Two
	}
	return nil
}

func (one *Runner_LeafRef_One) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["two"] = &one.Two

	return children
}

func (one *Runner_LeafRef_One) SetValue(value_path string, value string) {
	if value_path == "name-of-one" {
		one.NameOfOne = value
	}
}

func (one *Runner_LeafRef_One) GetBundleName() string {
	return "ydktest"
}

func (one *Runner_LeafRef_One) GetYangName() string {
	return "one"
}

func (one *Runner_LeafRef_One) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (one *Runner_LeafRef_One) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (one *Runner_LeafRef_One) SetParent(parent types.Entity) {
	one.parent = parent
}

func (one *Runner_LeafRef_One) GetParent() types.Entity {
	return one.parent
}

func (one *Runner_LeafRef_One) GetParentYangName() string {
	return "leaf-ref"
}

//////////////////////////////////////////////////////////////////////////
// Runner_LeafRef_One_Two
//////////////////////////////////////////////////////////////////////////
type Runner_LeafRef_One_Two struct {
	parent types.Entity
	Filter types.YFilter

	SelfRefOneName interface{} // str
}

func (two *Runner_LeafRef_One_Two) HasDataOrFilter() bool {
	return two.SelfRefOneName != nil
}

func (two *Runner_LeafRef_One_Two) GetFilter() types.YFilter {
	return two.Filter
}

func (two *Runner_LeafRef_One_Two) GetSegmentPath() string {
	return "two"
}

func (two *Runner_LeafRef_One_Two) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: two.GetSegmentPath()}
	var leafData types.LeafData

	if two.SelfRefOneName != nil {
		switch two.SelfRefOneName.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: two.SelfRefOneName.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", two.SelfRefOneName)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "self-ref-one-name", Data: leafData})
		}
	}

	return entityPath
}

func (two *Runner_LeafRef_One_Two) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (two *Runner_LeafRef_One_Two) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (two *Runner_LeafRef_One_Two) SetValue(value_path string, value string) {
	if value_path == "self-ref-one-name" {
		two.SelfRefOneName = value
	}
}

func (two *Runner_LeafRef_One_Two) GetBundleName() string {
	return "ydktest"
}

func (two *Runner_LeafRef_One_Two) GetYangName() string {
	return "two"
}

func (two *Runner_LeafRef_One_Two) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (two *Runner_LeafRef_One_Two) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (two *Runner_LeafRef_One_Two) SetParent(parent types.Entity) {
	two.parent = parent
}

func (two *Runner_LeafRef_One_Two) GetParent() types.Entity {
	return two.parent
}

func (two *Runner_LeafRef_One_Two) GetParentYangName() string {
	return "one"
}

//////////////////////////////////////////////////////////////////////////
// Runner_NotSupported1
//////////////////////////////////////////////////////////////////////////
type Runner_NotSupported1 struct {
	parent types.Entity
	Filter types.YFilter

	NotSupportedLeaf interface{} // str
	NotSupported12   Runner_NotSupported1_NotSupported12
}

func (notSupported1 *Runner_NotSupported1) HasDataOrFilter() bool {
	return notSupported1.NotSupportedLeaf != nil ||
		notSupported1.NotSupported12.HasDataOrFilter()
}

func (notSupported1 *Runner_NotSupported1) GetFilter() types.YFilter {
	return notSupported1.Filter
}

func (notSupported1 *Runner_NotSupported1) GetSegmentPath() string {
	return "not-supported-1"
}

func (notSupported1 *Runner_NotSupported1) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: notSupported1.GetSegmentPath()}
	var leafData types.LeafData

	if notSupported1.NotSupportedLeaf != nil {
		switch notSupported1.NotSupportedLeaf.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: notSupported1.NotSupportedLeaf.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", notSupported1.NotSupportedLeaf)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "not-supported-leaf", Data: leafData})
		}
	}

	return entityPath
}

func (notSupported1 *Runner_NotSupported1) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	if child_yang_name == "not-supported-1-2" {
		return &notSupported1.NotSupported12
	}
	return nil
}

func (notSupported1 *Runner_NotSupported1) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)
	children["not-supported-1-2"] = &notSupported1.NotSupported12

	return children
}

func (notSupported1 *Runner_NotSupported1) SetValue(value_path string, value string) {
	if value_path == "not-supported-leaf" {
		notSupported1.NotSupportedLeaf = value
	}
}

func (notSupported1 *Runner_NotSupported1) GetBundleName() string {
	return "ydktest"
}

func (notSupported1 *Runner_NotSupported1) GetYangName() string {
	return "notSupported1"
}

func (notSupported1 *Runner_NotSupported1) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (notSupported1 *Runner_NotSupported1) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (notSupported1 *Runner_NotSupported1) SetParent(parent types.Entity) {
	notSupported1.parent = parent
}

func (notSupported1 *Runner_NotSupported1) GetParent() types.Entity {
	return notSupported1.parent
}

func (notSupported1 *Runner_NotSupported1) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_NotSupported1_NotSupported12
//////////////////////////////////////////////////////////////////////////
type Runner_NotSupported1_NotSupported12 struct {
	parent types.Entity
	Filter types.YFilter

	SomeLeaf interface{} // str
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) HasDataOrFilter() bool {
	return notSupported12.SomeLeaf != nil
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetFilter() types.YFilter {
	return notSupported12.Filter
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetSegmentPath() string {
	return "not-supported-1-2"
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: notSupported12.GetSegmentPath()}
	var leafData types.LeafData

	if notSupported12.SomeLeaf != nil {
		switch notSupported12.SomeLeaf.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: notSupported12.SomeLeaf.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", notSupported12.SomeLeaf)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "some-leaf", Data: leafData})
		}
	}

	return entityPath
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) SetValue(value_path string, value string) {
	if value_path == "some-leaf" {
		notSupported12.SomeLeaf = value
	}
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetBundleName() string {
	return "ydktest"
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetYangName() string {
	return "notSupported12"
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) SetParent(parent types.Entity) {
	notSupported12.parent = parent
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetParent() types.Entity {
	return notSupported12.parent
}

func (notSupported12 *Runner_NotSupported1_NotSupported12) GetParentYangName() string {
	return "not-supported-1"
}

//////////////////////////////////////////////////////////////////////////
// Runner_NotSupported2
//////////////////////////////////////////////////////////////////////////
type Runner_NotSupported2 struct {
	parent types.Entity
	Filter types.YFilter

	Number interface{} // int32
}

func (notSupported2 *Runner_NotSupported2) HasDataOrFilter() bool {
	return notSupported2.Number != nil
}

func (notSupported2 *Runner_NotSupported2) GetFilter() types.YFilter {
	return notSupported2.Filter
}

func (notSupported2 *Runner_NotSupported2) GetSegmentPath() string {
	return "not-supported-2" + "[number=\"" + fmt.Sprintf("%v", notSupported2.Number) + "\"]"
}

func (notSupported2 *Runner_NotSupported2) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: notSupported2.GetSegmentPath()}
	var leafData types.LeafData

	if notSupported2.Number != nil {
		switch notSupported2.Number.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: notSupported2.Number.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", notSupported2.Number)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "number", Data: leafData})
		}
	}

	return entityPath
}

func (notSupported2 *Runner_NotSupported2) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (notSupported2 *Runner_NotSupported2) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (notSupported2 *Runner_NotSupported2) SetValue(value_path string, value string) {
	if value_path == "number" {
		notSupported2.Number = value
	}
}

func (notSupported2 *Runner_NotSupported2) GetBundleName() string {
	return "ydktest"
}

func (notSupported2 *Runner_NotSupported2) GetYangName() string {
	return "notSupported2"
}

func (notSupported2 *Runner_NotSupported2) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (notSupported2 *Runner_NotSupported2) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (notSupported2 *Runner_NotSupported2) SetParent(parent types.Entity) {
	notSupported2.parent = parent
}

func (notSupported2 *Runner_NotSupported2) GetParent() types.Entity {
	return notSupported2.parent
}

func (notSupported2 *Runner_NotSupported2) GetParentYangName() string {
	return "runner"
}

//////////////////////////////////////////////////////////////////////////
// Runner_Runner2
//////////////////////////////////////////////////////////////////////////
type Runner_Runner2 struct {
	parent types.Entity
	Filter types.YFilter

	SomeLeaf interface{} // str
}

func (runner2 *Runner_Runner2) HasDataOrFilter() bool {
	return runner2.SomeLeaf != nil
}

func (runner2 *Runner_Runner2) GetFilter() types.YFilter {
	return runner2.Filter
}

func (runner2 *Runner_Runner2) GetSegmentPath() string {
	return "runner-2"
}

func (runner2 *Runner_Runner2) GetEntityPath(entity types.Entity) types.EntityPath {
	entityPath := types.EntityPath{Path: runner2.GetSegmentPath()}
	var leafData types.LeafData

	if runner2.SomeLeaf != nil {
		switch runner2.SomeLeaf.(type) {
		case types.YFilter:
			leafData = types.LeafData{IsSet: false, Filter: runner2.SomeLeaf.(types.YFilter)}
		default:
			leafData = types.LeafData{IsSet: true, Value: fmt.Sprintf("%v", runner2.SomeLeaf)}
			entityPath.ValuePaths = append(entityPath.ValuePaths, types.NameLeafData{Name: "some-leaf", Data: leafData})
		}
	}

	return entityPath
}

func (runner2 *Runner_Runner2) GetChildByName(child_yang_name string, segment_path string) types.Entity {
	return nil
}

func (runner2 *Runner_Runner2) GetChildren() map[string]types.Entity {
	children := make(map[string]types.Entity)

	return children
}

func (runner2 *Runner_Runner2) SetValue(value_path string, value string) {
	if value_path == "some-leaf" {
		runner2.SomeLeaf = value
	}
}

func (runner2 *Runner_Runner2) GetBundleName() string {
	return "ydktest"
}

func (runner2 *Runner_Runner2) GetYangName() string {
	return "runner2"
}

func (runner2 *Runner_Runner2) GetBundleYangModelsLocation() string {
	return ydktest.YdkYdktestModelsPath
}

func (runner2 *Runner_Runner2) GetAugmentCapabilitiesFunction() types.AugmentCapabilitiesFunction {
	return ydktest.YdktestAugmentLookupTables
}

func (runner2 *Runner_Runner2) SetParent(parent types.Entity) {
	runner2.parent = parent
}

func (runner2 *Runner_Runner2) GetParent() types.Entity {
	return runner2.parent
}

func (runner2 *Runner_Runner2) GetParentYangName() string {
	return "runner"
}
