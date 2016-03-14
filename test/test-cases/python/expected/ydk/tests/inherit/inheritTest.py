
import unittest
from compare import is_equal
from ydk.providers import NetconfServiceProvider
from ydk.services import CRUDService

from ydk.models.inherit.inherit import Runner

class inheritTest(unittest.TestCase):

    def setUp(self):
        self.ncc = NetconfServiceProvider(address='localhost' , 
                            username='admin', password='admin', port=12022)
        self.crud = CRUDService()

    def test_one(self):
        one = Runner.One()
        one.one.name = 'Hello
        one.one.number = 100

