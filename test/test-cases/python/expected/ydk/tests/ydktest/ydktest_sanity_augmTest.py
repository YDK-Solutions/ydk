
import unittest
from compare import is_equal
from ydk.providers import NetconfServiceProvider
from ydk.services import CRUDService


class ydktest_sanity_augmTest(unittest.TestCase):

    def setUp(self):
        self.ncc = NetconfServiceProvider(address='localhost' , 
                            username='admin', password='admin', port=12022)
        self.crud = CRUDService()

