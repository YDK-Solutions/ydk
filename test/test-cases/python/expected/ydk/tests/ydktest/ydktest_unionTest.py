
import unittest
from compare import is_equal
from ydk.providers import NetconfServiceProvider
from ydk.services import CRUDService

from ydk.models.ydktest.ydktest_union import YdkEnumTest_Enum
from ydk.models.ydktest.ydktest_union import BuiltInT

class ydktest_unionTest(unittest.TestCase):

    def setUp(self):
        self.ncc = NetconfServiceProvider(address='localhost' , 
                            username='admin', password='admin', port=12022)
        self.crud = CRUDService()

