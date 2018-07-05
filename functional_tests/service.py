import unittest
from timeout_decorator import timeout_decorator
import functional_tests.chkit as chkit


class TestService(unittest.TestCase):
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
    def test_base_internal(self, depl: chkit.Deployment):
        svc = chkit.Service(
            name="test-service",
            deploy=depl.name,
            ports=[chkit.ServicePort(name="test-port", target_port=80, port=8888)]
        )
        try:
            chkit.create_service(svc)
            self.assertIn(svc.name, [service.name for service in chkit.get_services()])
            got_svc = chkit.get_service(svc.name)
            self.assertEqual(svc.name, got_svc.name)
            self.assertFalse(got_svc.is_external())
        finally:
            chkit.delete_service(svc.name)
            self.assertNotIn(svc.name, [service.name for service in chkit.get_services()])

    __default_internal_service = chkit.Service(
        name="test-service",
        deploy=__default_services_deployment.name,
        ports=[chkit.ServicePort(name="test-port", target_port=80, port=8888)]
    )

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment(deployment=__default_services_deployment)
    @chkit.with_service(service=__default_internal_service)
    def test_update_internal(self, depl: chkit.Deployment, svc: chkit.Service):
        new_svc = chkit.Service(
            name=svc.name,
            deploy=depl.name,
            ports=[chkit.ServicePort(name="test-port-1", target_port=80, port=8080)],
        )
        chkit.replace_service(new_svc)
        got_svc = chkit.get_service(svc.name)
        self.assertEqual(got_svc.name, new_svc.name)
        self.assertEqual(got_svc.deploy, new_svc.deploy)
        self.assertEqual(got_svc.ports[0].name, new_svc.ports[0].name)
        self.assertEqual(got_svc.ports[0].target_port, new_svc.ports[0].target_port[0])
        self.assertEqual(got_svc.ports[0].port, new_svc.ports[0].target_port)
