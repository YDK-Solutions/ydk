
import unittest
from compare import is_equal
from ydk.providers import NetconfServiceProvider
from ydk.services import CRUDService

from ydk.models.ydktest.ydktest_sanity import BaseIdentity_Identity
from ydk.models.ydktest.ydktest_sanity import ChildChildIdentity_Identity
from ydk.models.ydktest.ydktest_sanity import ChildIdentity_Identity
from ydk.models.ydktest.ydktest_sanity import YdkEnumTest_Enum
from ydk.models.ydktest.ydktest_sanity import Runner

class ydktest_sanityTest(unittest.TestCase):

    def setUp(self):
        self.ncc = NetconfServiceProvider(address='localhost' , 
                            username='admin', password='admin', port=12022)
        self.crud = CRUDService()

    def test_inbtw_list(self):
        inbtw_list = Runner.InbtwList()

    def test_inbtw_list_ldata(self):
        inbtw_list = Runner.InbtwList()
        ldata = inbtw_list.Ldata()
        ldata.parent = inbtw_list
        inbtw_list.ldata.append(ldata)
        ldata.inbtw_list.ldata.number = 100
        ldata.inbtw_list.ldata.name = 'Hello

    def test_inbtw_list_ldata_subc(self):
        inbtw_list = Runner.InbtwList()
        ldata = inbtw_list.Ldata()
        ldata.parent = inbtw_list
        inbtw_list.ldata.append(ldata)
        inbtw_list.ldata.subc.inbtw_list.ldata.subc.name = 'Hello
        inbtw_list.ldata.subc.inbtw_list.ldata.subc.number = 100

    def test_inbtw_list_ldata_subc_subc_subl1(self):
        inbtw_list = Runner.InbtwList()
        ldata = inbtw_list.Ldata()
        ldata.parent = inbtw_list
        inbtw_list.ldata.append(ldata)
        subc_subl1 = inbtw_list.ldata.subc.SubcSubl1()
        subc_subl1.parent = inbtw_list.ldata.subc
        inbtw_list.ldata.subc.subc_subl1.append(subc_subl1)
        subc_subl1.inbtw_list.ldata.subc.subc_subl1.number = 100
        subc_subl1.inbtw_list.ldata.subc.subc_subl1.name = 'Hello

    def test_leaf_ref(self):
        leaf_ref = Runner.LeafRef()
        leaf_ref.leaf_ref.ref_inbtw = None
        leaf_ref.leaf_ref.ref_one_name = None
        leaf_ref.leaf_ref.ref_three_sub1_sub2_number = None
        leaf_ref.leaf_ref.ref_two_sub1_number = None

    def test_leaf_ref_one(self):
        leaf_ref = Runner.LeafRef()
        leaf_ref.one.leaf_ref.one.name = 'Hello

    def test_leaf_ref_one_two(self):
        leaf_ref = Runner.LeafRef()
        leaf_ref.one.two.leaf_ref.one.two.self_ref_one_name = None

    def test_one(self):
        one = Runner.One()
        one.one.name = 'Hello
        one.one.number = 100

    def test_one_one_aug(self):
        one = Runner.One()
        one.one_aug.one.one_aug.name = 'Hello
        one.one_aug.one.one_aug.number = 100

    def test_one_list(self):
        one_list = Runner.OneList()

    def test_one_list_ldata(self):
        one_list = Runner.OneList()
        ldata = one_list.Ldata()
        ldata.parent = one_list
        one_list.ldata.append(ldata)
        ldata.one_list.ldata.number = 100
        ldata.one_list.ldata.name = 'Hello

    def test_one_list_one_aug_list(self):
        one_list = Runner.OneList()
        one_list.one_aug_list.one_list.one_aug_list.enabled = None

    def test_one_list_one_aug_list_ldata(self):
        one_list = Runner.OneList()
        ldata = one_list.one_aug_list.Ldata()
        ldata.parent = one_list.one_aug_list
        one_list.one_aug_list.ldata.append(ldata)
        ldata.one_list.one_aug_list.ldata.number = 100
        ldata.one_list.one_aug_list.ldata.name = 'Hello

    def test_runner_2(self):
        runner_2 = Runner.Runner2()
        runner_2.runner_2.some_leaf = 'Hello

    def test_three(self):
        three = Runner.Three()
        three.three.name = 'Hello
        three.three.number = 100

    def test_three_sub1(self):
        three = Runner.Three()
        three.sub1.three.sub1.number = 100

    def test_three_sub1_sub2(self):
        three = Runner.Three()
        three.sub1.sub2.three.sub1.sub2.number = 100

    def test_three_list(self):
        three_list = Runner.ThreeList()

    def test_three_list_ldata(self):
        three_list = Runner.ThreeList()
        ldata = three_list.Ldata()
        ldata.parent = three_list
        three_list.ldata.append(ldata)
        ldata.three_list.ldata.number = 100
        ldata.three_list.ldata.name = 'Hello

    def test_three_list_ldata_subl1(self):
        three_list = Runner.ThreeList()
        ldata = three_list.Ldata()
        ldata.parent = three_list
        three_list.ldata.append(ldata)
        subl1 = three_list.ldata.Subl1()
        subl1.parent = three_list.ldata
        three_list.ldata.subl1.append(subl1)
        subl1.three_list.ldata.subl1.number = 100
        subl1.three_list.ldata.subl1.name = 'Hello

    def test_three_list_ldata_subl1_sub_subl1(self):
        three_list = Runner.ThreeList()
        ldata = three_list.Ldata()
        ldata.parent = three_list
        three_list.ldata.append(ldata)
        subl1 = three_list.ldata.Subl1()
        subl1.parent = three_list.ldata
        three_list.ldata.subl1.append(subl1)
        sub_subl1 = three_list.ldata.subl1.SubSubl1()
        sub_subl1.parent = three_list.ldata.subl1
        three_list.ldata.subl1.sub_subl1.append(sub_subl1)
        sub_subl1.three_list.ldata.subl1.sub_subl1.number = 100
        sub_subl1.three_list.ldata.subl1.sub_subl1.name = 'Hello

    def test_two(self):
        two = Runner.Two()
        two.two.name = 'Hello
        two.two.number = 100

    def test_two_sub1(self):
        two = Runner.Two()
        two.sub1.two.sub1.number = 100

    def test_two_list(self):
        two_list = Runner.TwoList()

    def test_two_list_ldata(self):
        two_list = Runner.TwoList()
        ldata = two_list.Ldata()
        ldata.parent = two_list
        two_list.ldata.append(ldata)
        ldata.two_list.ldata.number = 100
        ldata.two_list.ldata.name = 'Hello

    def test_two_list_ldata_subl1(self):
        two_list = Runner.TwoList()
        ldata = two_list.Ldata()
        ldata.parent = two_list
        two_list.ldata.append(ldata)
        subl1 = two_list.ldata.Subl1()
        subl1.parent = two_list.ldata
        two_list.ldata.subl1.append(subl1)
        subl1.two_list.ldata.subl1.number = 100
        subl1.two_list.ldata.subl1.name = 'Hello

    def test_ytypes(self):
        ytypes = Runner.Ytypes()

    def test_ytypes_built_in_t(self):
        ytypes = Runner.Ytypes()
        ytypes.built_in_t.ytypes.built_in_t.bincoded = None
        ytypes.built_in_t.ytypes.built_in_t.bits_value = None
        ytypes.built_in_t.ytypes.built_in_t.bool_value = None
        ytypes.built_in_t.ytypes.built_in_t.deci64 = None
        ytypes.built_in_t.ytypes.built_in_t.emptee = None
        ytypes.built_in_t.enum_value = ytypes.built_in_t.YdkEnumTest_Enum()
        ytypes.built_in_t.ytypes.built_in_t.leaf_ref = None
        ytypes.built_in_t.ytypes.built_in_t.llstring = 'Hello
        ytypes.built_in_t.ytypes.built_in_t.name = 'Hello
        ytypes.built_in_t.ytypes.built_in_t.number16 = 100
        ytypes.built_in_t.ytypes.built_in_t.number32 = 100
        ytypes.built_in_t.ytypes.built_in_t.number64 = 100
        ytypes.built_in_t.ytypes.built_in_t.number8 = 100
        ytypes.built_in_t.ytypes.built_in_t.u_number16 = 100
        ytypes.built_in_t.ytypes.built_in_t.u_number32 = 100
        ytypes.built_in_t.ytypes.built_in_t.u_number64 = 100
        ytypes.built_in_t.ytypes.built_in_t.u_number8 = 100
        ytypes.built_in_t.ytypes.built_in_t.younion = None

    def test_ytypes_built_in_t_identity_ref_value(self):
        ytypes = Runner.Ytypes()

    def test_ytypes_derived_t(self):
        ytypes = Runner.Ytypes()

