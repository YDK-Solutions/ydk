/*  ----------------------------------------------------------------
 YDK - YANG Development Kit
 Copyright 2016 Cisco Systems. All rights reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 -------------------------------------------------------------------
 This file has been modified by Yan Gorelik, YDK Solutions.
 All modifications in original under CiscoDevNet domain
 introduced since October 2019 are copyrighted.
 All rights reserved under Apache License, Version 2.0.
 ------------------------------------------------------------------*/

package test

import (
	"fmt"
	ysanity "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/sanity"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/CiscoDevNet/ydk-go/ydk/types/ylist"
	"github.com/stretchr/testify/suite"
	"testing"
	"sort"
)

type CoreTypesTestSuite struct {
	suite.Suite
}

func (suite *CoreTypesTestSuite) SetupSuite() {
}

func (suite *CoreTypesTestSuite) TestEntityCollection() {
    // Create Data entities and access values
    runner := ysanity.Runner{}
    native := ysanity.Native{}
    runnerPath := types.GetSegmentPath(&runner)
    nativePath := types.GetSegmentPath(&native)
    suite.Equal(types.EntityToString(&runner), "Type: *sanity.Runner, Path: ydktest-sanity:runner")
    suite.Equal(types.EntityToString(&native), "Type: *sanity.Native, Path: ydktest-sanity:native")

    // Initialization
    config := types.NewEntityCollection()
    suite.Equal(config.Len(), 0)
    suite.Equal(config.String(), "EntityCollection is empty")

    config = types.NewEntityCollection(&runner)
    suite.Equal(config.Len(), 1)

    config = types.NewEntityCollection(&runner, &native)
    suite.Equal(config.Len(), 2)
    suite.Equal(config.String(), "EntityCollection [Type: *sanity.Runner, Path: ydktest-sanity:runner; Type: *sanity.Native, Path: ydktest-sanity:native]")

    // Add
    config = types.NewEntityCollection()
    config.Append([]types.Entity{&runner, &native})
    suite.Equal(config.Len(), 2)

    config.Add(&runner)
    suite.Equal(config.Len(), 2)

    // Get
    e, _ := config.Get(runnerPath)
    suite.NotNil(e)
    suite.IsType(&runner, e)

    // HasKey
    suite.Equal(config.HasKey(runnerPath), true)
    suite.Equal(config.HasKey(nativePath), true)
    suite.Equal(config.HasKey("oc_bgp"), false)

    // Get all keys
    suite.Equal(config.Keys(), []string{runnerPath, nativePath})

    // Get all entities
    ydk.YLogDebug("All entities:")
    for _, entity := range config.Entities() {
    	ydk.YLogDebug(types.EntityToString(entity))
    }

    // Delete entity
    e, _ = config.Pop(runnerPath)
    suite.NotNil(e)
    suite.IsType(&runner, e)
    suite.Equal(config.Keys(), []string{nativePath})

    ydk.YLogDebug("All entities after Runner deleted:")
    for _, key := range config.Keys() {
        en, exist := config.Get(key)
        if exist {
	        ydk.YLogDebug(fmt.Sprintf("%s", types.EntityToString(en)))
        }
    }

    // Add back and test order
    config.Add(&runner)
    suite.Equal(config.Keys(), []string{nativePath, runnerPath})

	// Getting enities by item number
	ydk.YLogDebug("Getting enities by item number:")
	for i:=0; i<3; i++ {
		e = config.GetItem(i)
		if e != nil {
			ydk.YLogDebug(fmt.Sprintf("%d:  %s", i, types.EntityToString(e)))
		} else {
			ydk.YLogDebug(fmt.Sprintf("%d:  nil\n", i))
		}
	}
    // Clear collection
    config.Clear()
    suite.Equal(config.Len(), 0)

    // Testing passing parameters and return values
    ret1 := testParams("test1", config)
    col1 := types.EntityToCollection(ret1)
    suite.Equal(col1.Len(), 0)

    ret2 := testParams("test2", &runner)
    suite.Equal(types.EntityEqual(ret2, &runner), true)

    config.Add(&runner, &native)
    ret3 := testParams("test3", config)
    col3 := types.EntityToCollection(ret3)
    suite.Equal(col3.Len(), 2)

    // Testing Config and Filter aliases
    cfg := types.NewConfig()
	cfg.Add(&runner, &native)
	suite.Equal(types.IsEntityCollection(cfg), true)
	filter := types.NewFilter(&native)
	suite.Equal(types.IsEntityCollection(filter), true)
	suite.Equal(filter.Len(), 1)
}

func (suite *CoreTypesTestSuite) TestListTwoKeys() {
	runner := ysanity.Runner{}
	l1 := ysanity.Runner_TwoKeyList{First: "f1", Second: 11}
	l2 := ysanity.Runner_TwoKeyList{First: "f2", Second: 22}
	runner.TwoKeyList = []*ysanity.Runner_TwoKeyList {&l1, &l2}

	ldataKeys := ylist.Keys(runner.TwoKeyList)
	suite.Equal(fmt.Sprintf("%v", ldataKeys), "[[f1 11] [f2 22]]")

	for _, lkey := range ldataKeys {
		_, ldata := ylist.Get(runner.TwoKeyList, lkey);
		suite.NotNil(ldata)
	}
	_, ldata := ylist.Get(runner.TwoKeyList, "f1", 11)
	suite.Equal(types.EntityEqual(ldata, &l1), true)
	_, ldata = ylist.Get(runner.TwoKeyList, "f2", 22)
	suite.Equal(types.EntityEqual(ldata, &l2), true)

	i, ldata := ylist.Get(runner.TwoKeyList, "f1", 11)
	suite.Equal(i, 0)
	suite.NotNil(ldata)
	suite.Equal(types.EntityEqual(ldata, &l1), true)
	if ldata != nil {
		runner.TwoKeyList = append(runner.TwoKeyList[:i], runner.TwoKeyList[i+1:]...)
		i, ldata = ylist.Get(runner.TwoKeyList, "f1", 11)
		suite.Nil(ldata)
	}
}

func (suite *CoreTypesTestSuite) TestIdentityList() {
	runner := ysanity.Runner{}
	i1 := ysanity.Runner_IdentityList{Name: "first"}
	i2 := ysanity.Runner_IdentityList{Name: "second"}
	i3 := ysanity.Runner_IdentityList{Name: "third"}
	runner.IdentityList = []*ysanity.Runner_IdentityList {&i1, &i2, &i3}

	ldataKeys := ylist.Keys(runner.IdentityList)
	suite.Equal(fmt.Sprintf("%v", ldataKeys), "[first second third]")

	var ldata types.Entity
	for _, lkey := range ldataKeys {
		_, ldata = ylist.Get(runner.IdentityList, lkey);
		suite.NotNil(ldata)
	}
	_, ldata = ylist.Get(runner.IdentityList, "first")
	suite.Equal(types.EntityEqual(ldata, &i1), true)
	_, ldata = ylist.Get(runner.IdentityList, "third")
	suite.Equal(types.EntityEqual(ldata, &i3), true)
}

func (suite *CoreTypesTestSuite) TestEnumList() {
	runner := ysanity.Runner{}
	i1 := ysanity.Runner_EnumList{KeyName: ysanity.YdkEnumTest_none}
	i2 := ysanity.Runner_EnumList{KeyName: ysanity.YdkEnumTest_local}
	i3 := ysanity.Runner_EnumList{KeyName: ysanity.YdkEnumTest_remote}
	runner.EnumList = []*ysanity.Runner_EnumList{&i1, &i2, &i3}

	ldataKeys := ylist.Keys(runner.EnumList)
	suite.Equal(fmt.Sprintf("%v", ldataKeys), "[none local remote]")

	var ldata types.Entity
	for _, lkey := range ldataKeys {
		_, ldata = ylist.Get(runner.EnumList, lkey);
		suite.NotNil(ldata)
	}
	_, ldata = ylist.Get(runner.EnumList, "none")
	suite.Equal(types.EntityEqual(ldata, &i1), true)
	_, ldata = ylist.Get(runner.EnumList, "remote")
	suite.Equal(types.EntityEqual(ldata, &i3), true)
}

func (suite *CoreTypesTestSuite) TestListNoKeys() {
	runner := ysanity.Runner{}
	t1 := ysanity.Runner_NoKeyList{Test: "t1"}
	t2 := ysanity.Runner_NoKeyList{Test: "t2"}
	t3 := ysanity.Runner_NoKeyList{Test: "t3"}
	runner.NoKeyList = []*ysanity.Runner_NoKeyList {&t1, &t2, &t3}

	suite.Equal(len(ylist.Keys(runner.NoKeyList)), 0)
	var count string
	for _, elem := range runner.NoKeyList {
		count = count + elem.Test.(string)
	}
	suite.Equal(count, "t1t2t3")
}

func checkEmptyStringValue(v string) string {
	if len(v) == 0 {
		v = "Exists"
	}
	return v
}

func printDictionary(legend string, entdict map[string]string) {
	fmt.Printf("\n------> DICTIONARY%s\n", legend)
	var keys []string
	for k := range entdict {
	        keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, checkEmptyStringValue(entdict[k]))
	}
}

func printDiff(suite *CoreTypesTestSuite, diff map[string]types.StringPair, entLeft, entRight types.Entity ) {
	fmt.Printf("\n------> DIFFS:\n")
	var keys []string
	for k := range diff {
	        keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		valuePair := diff[key]
		fmt.Printf("%s: %s vs %s\n", key,
			checkEmptyStringValue(valuePair.Left),
			checkEmptyStringValue(valuePair.Right))
		ent := entLeft
		if valuePair.Left == "None" {
			ent = entRight
		}
		suite.NotNil(types.PathToEntity(ent, key))
	}
}

func (suite *CoreTypesTestSuite) TestEntityDiff_TwoKeys() {
	runner1 := ysanity.Runner{}
	l1 := ysanity.Runner_TwoKeyList{First: "f1", Second: 11, Property: 82}
	l2 := ysanity.Runner_TwoKeyList{First: "f2", Second: 22, Property: 83}
	runner1.TwoKeyList = []*ysanity.Runner_TwoKeyList {&l1, &l2}
	types.SetAllParents(&runner1)

	entDict1 := types.EntityToDict(&runner1)
	suite.Equal(len(entDict1), 4)
	printDictionary("-ORIGINAL", entDict1)

	runner2 := ysanity.Runner{}
	ll1 := ysanity.Runner_TwoKeyList{First: "f1", Second: 11, Property: 82}
	ll2 := ysanity.Runner_TwoKeyList{First: "f2", Second: 22, Property: 83}
	runner2.TwoKeyList = []*ysanity.Runner_TwoKeyList {&ll1, &ll2}
	types.SetAllParents(&runner2)

	diff := types.EntityDiff(&runner1, &runner2)
	suite.Equal(0, len(diff))

	ll1.Property = 83
	entDict2 := types.EntityToDict(&runner2)
	printDictionary("-CHANGED-LEAF", entDict2)

	diff = types.EntityDiff(&runner1, &runner2)
	suite.Equal(1, len(diff))
	printDiff(suite, diff, &runner1, &runner2)

	ll1.First = "f2"
	entDict3 := types.EntityToDict(&runner2)
	printDictionary("-CHANGED-KEY", entDict3)

	diff = types.EntityDiff(&runner1, &runner2)
	suite.Equal(2, len(diff))
	printDiff(suite, diff, &runner1, &runner2)

	clone := types.EntityClone(&runner2)
	diff = types.EntityDiff(clone, &runner2)
	suite.Equal(0, len(diff))
}

func (suite *CoreTypesTestSuite) TestEntityDiff_Presence() {
	runner := ysanity.Runner{}
	runner.Runner2.SomeLeaf = "some_leaf"
	types.SetPresenceFlag(&runner.Runner2)
	types.SetAllParents(&runner)

	entDict1 := types.EntityToDict(&runner)
	suite.Equal(len(entDict1), 2)
	printDictionary("-LEFT", entDict1)

	emptyRunner := ysanity.Runner{}
	types.SetAllParents(&emptyRunner)
	entDict2 := types.EntityToDict(&emptyRunner)
	suite.Equal(len(entDict2), 0)
	printDictionary("-RIGHT", entDict2);

	diff := types.EntityDiff(&runner, &emptyRunner)
	suite.Equal(1, len(diff))
	printDiff(suite, diff, &runner, &emptyRunner)

	clone := types.EntityClone(&runner)
	diff = types.EntityDiff(&runner, clone)
	suite.Equal(0, len(diff))
}

func (suite *CoreTypesTestSuite) TestEntityDiff_TwoListPos() {
	runner1 := ysanity.Runner{}
	elem1 := ysanity.Runner_TwoList_Ldata{}
	elem2 := ysanity.Runner_TwoList_Ldata{}
	elem1.Number = 1
	elem1.Name = "runner:twolist:ldata:1"
	elem2.Number = 2
	elem2.Name = "runner:twolist:ldata:2"

	elem11 := ysanity.Runner_TwoList_Ldata_Subl1{}
	elem12 := ysanity.Runner_TwoList_Ldata_Subl1{}
	elem11.Number = 11
	elem11.Name = "runner:twolist:ldata:1:subl1:11"
	elem12.Number = 12
	elem12.Name = "runner:twolist:ldata:1:subl1:12"

	elem1.Subl1 = append(elem1.Subl1, &elem11)
	elem1.Subl1 = append(elem1.Subl1, &elem12)

	elem21 := ysanity.Runner_TwoList_Ldata_Subl1{}
	elem22 := ysanity.Runner_TwoList_Ldata_Subl1{}
	elem21.Number = 21
	elem21.Name = "runner:twolist:ldata:2:subl1:21"
	elem22.Number = 22
	elem22.Name = "runner:twolist:ldata:2:subl1:22"

	elem2.Subl1 = append(elem2.Subl1, &elem21)
	elem2.Subl1 = append(elem2.Subl1, &elem22)

	runner1.TwoList.Ldata = append(runner1.TwoList.Ldata, &elem1)
	runner1.TwoList.Ldata = append(runner1.TwoList.Ldata, &elem2)
	types.SetAllParents(&runner1)

	entDict := types.EntityToDict(&runner1)
	suite.Equal(12, len(entDict))
	printDictionary("-LEFT", entDict)

	runner2 := ysanity.Runner{}
	elem1 = ysanity.Runner_TwoList_Ldata{}
	elem1.Number = 1
	elem1.Name = "runner:twolist:ldata:1"
	elem11 = ysanity.Runner_TwoList_Ldata_Subl1{}
	elem12 = ysanity.Runner_TwoList_Ldata_Subl1{}
	elem11.Number = 11
	elem11.Name = "runner:twolist:ldata:1:subl1:11"
	elem12.Number = 12
	elem12.Name = "runner:twolist:ldata:1:subl1:12"

	elem1.Subl1 = append(elem1.Subl1, &elem11)
	elem1.Subl1 = append(elem1.Subl1, &elem12)

	runner2.TwoList.Ldata = append(runner2.TwoList.Ldata, &elem1)
	types.SetAllParents(&runner2)

	entDict = types.EntityToDict(&runner2)
	suite.Equal(6, len(entDict))
	printDictionary("-RIGHT", entDict)

	diff := types.EntityDiff(&runner1, &runner2)
	suite.Equal(1, len(diff))
	printDiff(suite, diff, &runner1, &runner2)

	clone := types.EntityClone(&runner1)
	diff = types.EntityDiff(&runner1, clone)
	suite.Equal(0, len(diff))
}

func (suite *CoreTypesTestSuite) TestEntityDiff_ListNoKeys() {
	t1 := ysanity.Runner_NoKeyList{Test: "t1"}
	t2 := ysanity.Runner_NoKeyList{Test: "tt"}
	t3 := ysanity.Runner_NoKeyList{Test: "t3"}
	runner1 := ysanity.Runner{}
	runner1.NoKeyList = []*ysanity.Runner_NoKeyList {&t1, &t2, &t3}

	entDict := types.EntityToDict(&runner1)
	suite.Equal(6, len(entDict))
	printDictionary("-LEFT", entDict)

	t11 := ysanity.Runner_NoKeyList{Test: "t1"}
	t22 := ysanity.Runner_NoKeyList{Test: "t2"}
	runner2 := ysanity.Runner{}
	runner2.NoKeyList = []*ysanity.Runner_NoKeyList {&t11, &t22}

	entDict = types.EntityToDict(&runner2)
	suite.Equal(4, len(entDict))
	printDictionary("-RIGHT", entDict)

	diff := types.EntityDiff(&runner1, &runner2)
	suite.Equal(2, len(diff))
	printDiff(suite, diff, &runner1, &runner2)

	diff = types.EntityDiff(&runner2, &runner1)
	suite.Equal(2, len(diff))
	printDiff(suite, diff, &runner2, &runner1)

	clone := types.EntityClone(&runner1)
	diff = types.EntityDiff(&runner1, clone)
	suite.Equal(0, len(diff))
}

func (suite *CoreTypesTestSuite) TestEntityDiff_OneAugList() {
	oneAugList := ysanity.Runner_OneList_OneAugList{}
	oneAugList.Enabled = true

	elem1 := ysanity.Runner_OneList_OneAugList_Ldata{}
	elem1.Number = 1
	elem1.Name = "elem1"

	elem2 := ysanity.Runner_OneList_OneAugList_Ldata{}
	elem2.Number = 2
	elem2.Name = "elem2"

	oneAugList.Ldata = []*ysanity.Runner_OneList_OneAugList_Ldata {&elem1, &elem2}
	types.SetAllParents(&oneAugList)

	entDict := types.EntityToDict(&oneAugList)
	suite.Equal(5, len(entDict))
	printDictionary("", entDict)

	clone := types.EntityClone(&oneAugList)
	diff := types.EntityDiff(&oneAugList, clone)
	suite.Equal(0, len(diff))
}

func (suite *CoreTypesTestSuite) TestEntityDiff_EnumLeafList() {
	runner := ysanity.Runner{}
	runner.Ytypes.BuiltInT.EnumLlist = append(runner.Ytypes.BuiltInT.EnumLlist, ysanity.YdkEnumTest_local)
	runner.Ytypes.BuiltInT.EnumLlist = append(runner.Ytypes.BuiltInT.EnumLlist, ysanity.YdkEnumTest_remote)
	entDict := types.EntityToDict(&runner)
	suite.Equal(2, len(entDict))
	printDictionary("", entDict)

	clone := types.EntityClone(&runner)
	diff := types.EntityDiff(&runner, clone)
	suite.Equal(0, len(diff))
	suite.True(types.EntityEqual(&runner, clone))
}

func testParams(testName string, entity types.Entity) types.Entity {
	ec := types.EntityToCollection(entity)
	if ec == nil {
		ydk.YLogDebug(fmt.Sprintf("%s: %s\n", testName, types.EntityToString(entity)))
		return entity
	}
	ydk.YLogDebug(fmt.Sprintf("%s: %s\n", testName, ec.String()))
	return ec
}

func TestCoreTypesTestSuite(t *testing.T) {
	if testing.Verbose() {
		ydk.EnableLogging(ydk.Debug)
	}
	suite.Run(t, new(CoreTypesTestSuite))
}
