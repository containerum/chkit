import unittest
import functional_tests.chkit as chkit


class TestConfig(unittest.TestCase):

    def test_login(self):
        users = [("test1@containerum.io", "12345678"), ("helpik94@yandex.com", "12345678")]
        for user in users:
            chkit.login(user=user[0], password=user[1])
            profile = chkit.get_profile()
            self.assertEqual(profile['Login'], user[0])

    def test_namespace_selector(self):
        chkit.login(user="helpik94@yandex.com", password="12345678")
        owner, namespace = chkit.get_default_namespace()
        self.assertIsNotNone(owner)
        self.assertIsNotNone(namespace)
        chkit.set_default_namespace(namespace="mynewns")
        owner, namespace = chkit.get_default_namespace()
        self.assertEqual(namespace, "mynewns")

    def test_api_selector(self):
        try:
            chkit.login(user="helpik94@yandex.com", password="12345678")
            self.assertEqual(chkit.get_api_url(), chkit.DEFAULT_API_URL)
            chkit.set_api_url("http://test.api.domain.com")
            self.assertEqual(chkit.get_api_url(), "http://test.api.domain.com")
        finally:
            chkit.set_api_url()
