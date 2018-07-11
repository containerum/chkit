import unittest
from timeout_decorator import timeout_decorator
import functional_tests.chkit as chkit
import time
import requests


class TestInternalService(unittest.TestCase):
    __default_services_deployment = chkit.Deployment(
        name="default-services-test-depl",
        replicas=1,
        containers=[
            chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10)),
        ],
    )

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment(deployment=__default_services_deployment)
    def test_base(self, depl: chkit.Deployment):
        svc = chkit.Service(
            name="test-internal-service",
            deploy=depl.name,
            ports=[chkit.ServicePort(name="test-internal-port", target_port=80, port=8888)]
        )
        try:
            chkit.create_service(svc)
            self.assertIn(svc.name, [service.name for service in chkit.get_services()])
            got_svc = chkit.get_service(svc.name)
            self.assertEqual(svc.name, got_svc.name)
            self.assertFalse(got_svc.is_external())
        finally:
            chkit.delete_service(svc.name)
            time.sleep(1)
            self.assertNotIn(svc.name, [service.name for service in chkit.get_services()])

    __default_internal_service = chkit.Service(
        name="test-internal-service",
        deploy=__default_services_deployment.name,
        ports=[chkit.ServicePort(name="test-internal-port", target_port=80, port=8888)]
    )

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment(deployment=__default_services_deployment)
    @chkit.with_service(service=__default_internal_service)
    def test_update(self, depl: chkit.Deployment, svc: chkit.Service):
        new_svc = chkit.Service(
            name=svc.name,
            deploy=depl.name,
            ports=[chkit.ServicePort(name="test-internal-port-update", target_port=443, port=9999)],
        )
        chkit.replace_service(new_svc)
        got_svc = chkit.get_service(svc.name)
        self.assertEqual(got_svc.name, new_svc.name)
        self.assertEqual(got_svc.deploy, new_svc.deploy)
        self.assertEqual(got_svc.ports[0].name, new_svc.ports[0].name)
        self.assertEqual(got_svc.ports[0].target_port, new_svc.ports[0].target_port)
        self.assertEqual(got_svc.ports[0].port, new_svc.ports[0].port)


class TestExternalService(unittest.TestCase):
    __default_services_deployment = chkit.Deployment(
        name="default-services-test-depl",
        replicas=1,
        containers=[
            chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10)),
        ],
    )

    @timeout_decorator.timeout(seconds=650*2)
    @chkit.test_account
    @chkit.with_deployment(deployment=__default_services_deployment)
    @chkit.ensure_pods_running(deployment=__default_services_deployment.name)
    def test_base(self, depl: chkit.Deployment):
        svc = chkit.Service(
            name="test-external-service",
            deploy=depl.name,
            ports=[chkit.ServicePort(name="test-external-port", target_port=80)]
        )
        try:
            chkit.create_service(svc)
            self.assertIn(svc.name, [service.name for service in chkit.get_services()])
            got_svc = chkit.get_service(svc.name)
            self.assertEqual(svc.name, got_svc.name)
            self.assertTrue(got_svc.is_external())
            attempts, max_attempts = 1, 40
            while attempts <= max_attempts:
                try:
                    response = requests.get(f"http://{got_svc.ips[0]}:{got_svc.ports[0].port}",
                                            headers={"Host": got_svc.domain}, timeout=1)
                    response.raise_for_status()
                    if response.status_code == 200:
                        break
                except requests.exceptions.ConnectionError:
                    pass
                time.sleep(15)
                attempts += 1
            self.assertLessEqual(attempts, max_attempts)
        finally:
            chkit.delete_service(svc.name)
            time.sleep(1)
            self.assertNotIn(svc.name, [service.name for service in chkit.get_services()])

    __default_external_service = chkit.Service(
        name="test-external-service",
        deploy=__default_services_deployment.name,
        ports=[chkit.ServicePort(name="test-external-port", target_port=80)]
    )

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment(deployment=__default_services_deployment)
    @chkit.with_service(service=__default_external_service)
    def test_update(self, depl: chkit.Deployment, svc: chkit.Service):
        new_svc = chkit.Service(
            name=svc.name,
            deploy=depl.name,
            ports=[chkit.ServicePort(name="test-external-port-update", target_port=443)]
        )
        chkit.replace_service(service=new_svc, file=True)
        got_svc = chkit.get_service(service=new_svc.name)
        self.assertEqual(new_svc.name, got_svc.name)
        self.assertEqual(len(got_svc.ports), 1)
        self.assertEqual(got_svc.ports[0].name, new_svc.ports[0].name)
        self.assertEqual(got_svc.ports[0].target_port, new_svc.ports[0].target_port)
