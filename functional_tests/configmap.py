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
        chkit.replace_configmap(new_configmap)
        got_cm = chkit.get_configmap(new_configmap.name)
        self.assertEqual(got_cm.name, new_configmap.name)
        self.assertEqual(got_cm.data["TESTKEY1"], new_configmap.data["TESTKEY1"])

    __default_cm_deployment = chkit.Deployment(
        name="cm-test-deploy",
        replicas=1,
        containers=[
            chkit.Container(
                name="first",
                limits=chkit.Resources(cpu=15, memory=15),
                image="nginx",
                config_maps=[chkit.DeploymentConfigMap(name="test-configmap")],
            )
        ],
    )

    @timeout_decorator.timeout(seconds=650*2)
    @chkit.test_account
    @chkit.with_cm(configmap=__default_cm)
    @chkit.with_deployment(deployment=__default_cm_deployment)
    @chkit.ensure_pods_running(deployment=__default_cm_deployment.name)
    def test_deploy_mount(self, cm: chkit.ConfigMap, depl: chkit.Deployment):
        chkit.get_configmap(cm.name)
        depl_created = chkit.get_deployment(depl.name)
        self.assertEqual(depl_created.containers[0].config_maps[0].name, cm.name)
