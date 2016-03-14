
import unittest
from compare import is_equal
from ydk.providers import NetconfServiceProvider
from ydk.services import CRUDService

from ydk.models.ydktest.ydktest_filterread import A

class ydktest_filterreadTest(unittest.TestCase):

    def setUp(self):
        self.ncc = NetconfServiceProvider(address='localhost' , 
                            username='admin', password='admin', port=12022)
        self.crud = CRUDService()

    def test_b(self):
        b = A.B()
        b.b.b1 = 'Hello
        b.b.b2 = 'Hello
        b.b.b3 = 'Hello

    def test_b_c(self):
        b = A.B()

    def test_b_d(self):
        b = A.B()
        b.d.b.d.d1 = 'Hello
        b.d.b.d.d2 = 'Hello
        b.d.b.d.d3 = 'Hello

    def test_b_d_e(self):
        b = A.B()
        b.d.e.b.d.e.e1 = 'Hello
        b.d.e.b.d.e.e2 = 'Hello

    def test_lst(self):
        lst = A.Lst()
        lst.lst.number = 100
        lst.lst.value = 'Hello

