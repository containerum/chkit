import unittest
from timeout_decorator import timeout_decorator
import functional_tests.chkit as chkit
import time
import requests


class TestConfigMap(unittest.TestCase):
    __default_cm = chkit.ConfigMap(
        name="test-configmap",
        data=dict(TESTKEY="TESTVALUE")
    )

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    def test_base(self):
        configmap = chkit.ConfigMap(
            name="test-configmap",
            data=dict(TESTKEY="TESTVALUE")
        )
        try:
            chkit.create_configmap(configmap)
            self.assertIn(configmap.name, [cm.name for cm in chkit.get_configmaps()])
            got_cm = chkit.get_configmap(configmap.name)
            self.assertEqual(configmap.name, got_cm.name)
            self.assertEqual(configmap.data["TESTKEY"], got_cm.data["TESTKEY"])
        finally:
            chkit.delete_configmap(configmap.name)
            time.sleep(1)
            self.assertNotIn(configmap.name, [cm.name for cm in chkit.get_configmaps()])

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_cm(configmap=__default_cm)
    def test_update(self, cm: chkit.ConfigMap):
        new_configmap = chkit.ConfigMap(
            name=cm.name,
            data=dict(TESTKEY1="TESTVALUE1")
        )
        print("TEST",new_configmap.data)
        chkit.replace_configmap(new_configmap)
        got_cm = chkit.get_configmap(new_configmap.name)
        self.assertEqual(got_cm.name, new_configmap.name)
        self.assertEqual(got_cm.data["TESTKEY1"], new_configmap.data["TESTKEY1"])
