
import unittest
from compare import is_equal
from ydk.providers import NetconfServiceProvider
from ydk.services import CRUDService

from ydk.models.main.main import A

class mainTest(unittest.TestCase):

    def setUp(self):
        self.ncc = NetconfServiceProvider(address='localhost' , 
                            username='admin', password='admin', port=12022)
        self.crud = CRUDService()

    def test_main_aug1_c(self):
        main_aug1_c = A.MainAug1_C()
        main_aug1_c.main_aug1_c.two = 'Hello

    def test_main_aug2_c(self):
        main_aug2_c = A.MainAug2_C()
        main_aug2_c.main_aug2_c.three = 100

    def test_main_aug2_d(self):
        main_aug2_d = A.MainAug2_D()
        main_aug2_d.main_aug2_d.poo = 100

    def test_main_aug3_c(self):
        main_aug3_c = A.MainAug3_C()
        main_aug3_c.main_aug3_c.meh = 100

    def test_main_aug3_d(self):
        main_aug3_d = A.MainAug3_D()
        main_aug3_d.main_aug3_d.buh = 'Hello

