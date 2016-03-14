
import unittest

class ImportTest(unittest.TestCase):


    def test_inherit(self):
        from ydk.models.inherit.inherit import Runner


    def test_main(self):
        from ydk.models.main.main import A


    def test_main_aug1(self):
        pass


    def test_main_aug2(self):
        pass


    def test_main_aug3(self):
        pass


    def test_oc_pattern(self):
        from ydk.models.oc.oc_pattern import A


    def test_ydktest_filterread(self):
        from ydk.models.ydktest.ydktest_filterread import A


    def test_ydktest_sanity(self):
        from ydk.models.ydktest.ydktest_sanity import BaseIdentity_Identity
        from ydk.models.ydktest.ydktest_sanity import ChildChildIdentity_Identity
        from ydk.models.ydktest.ydktest_sanity import ChildIdentity_Identity
        from ydk.models.ydktest.ydktest_sanity import YdkEnumTest_Enum
        from ydk.models.ydktest.ydktest_sanity import Runner


    def test_ydktest_sanity_augm(self):
        pass


    def test_ydktest_union(self):
        from ydk.models.ydktest.ydktest_union import YdkEnumTest_Enum
        from ydk.models.ydktest.ydktest_union import BuiltInT


if __name__ == '__main__':
    unittest.main()

