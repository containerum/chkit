import unittest
from timeout_decorator import timeout_decorator
import functional_tests.chkit as chkit
import time
import requests


class TestIngress(unittest.TestCase):
    __default_services_deployment = chkit.Deployment(
        name="default-services-test-depl",
        replicas=1,
        containers=[
            chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10)),
        ],
    )

    __default_external_service = chkit.Service(
        name="test-external-service",
        deploy=__default_services_deployment.name,
        ports=[chkit.ServicePort(name="test-external-port", target_port=80), chkit.ServicePort(name="test-external-port-2", target_port=443)]
    )

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment(deployment=__default_services_deployment)
    @chkit.with_service(service=__default_external_service)
    def test_base(self, depl: chkit.Deployment, svc: chkit.Service):
        ingr = chkit.Ingress(
            name="test-ingress",
            rules=[chkit.IngressRules(
                host="test-host",
                path=[chkit.IngressPath(
                    path="/",
                    service_name=svc.name,
                    service_port=svc.ports[0].port
                )]
            )]
        )
        try:
            got_svc = chkit.get_service(svc.name)
            ingr.rules[0].path[0].service_port = got_svc.ports[0].port
            chkit.create_ingress(ingr)
            self.assertIn(ingr.name[0], [ingr.name[0] for ingr in chkit.get_ingresses()])
            got_ingr = chkit.get_ingress(ingr.name[0])
            self.assertEqual(got_ingr.name[0], ingr.name[0])
            self.assertEqual(got_ingr.rules[0].host, ingr.rules[0].host + ".hub.containerum.io")
            self.assertEqual(got_ingr.rules[0].path[0].path, ingr.rules[0].path[0].path)
            self.assertEqual(got_ingr.rules[0].path[0].service_port, ingr.rules[0].path[0].service_port)
            self.assertEqual(got_ingr.rules[0].path[0].service_name, ingr.rules[0].path[0].service_name)
        finally:
            chkit.delete_ingress(ingr.name[0])
            time.sleep(1)
            self.assertNotIn(ingr.name[0], [ingrs.name[0] for ingrs in chkit.get_ingresses()])

    __default_internal_service = chkit.Service(
        name="test-internal-service",
        deploy=__default_services_deployment.name,
        ports=[chkit.ServicePort(name="test-int-port", port=80, target_port=80), chkit.ServicePort(name="test-int-port-2", port=443, target_port=443)]
    )

    __default_update_ingress = chkit.Ingress(
        name="test-ingress",
        rules=[chkit.IngressRules(
            host="test-host",
            path=[chkit.IngressPath(
                path="/",
                service_name=__default_internal_service.name,
                service_port=__default_internal_service.ports[0].port
            )]
        )]
    )

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment(deployment=__default_services_deployment)
    @chkit.with_service(service=__default_internal_service)
    @chkit.with_ingress(ingress=__default_update_ingress)
    def test_update(self, depl: chkit.Deployment, svc: chkit.Service, ingr: chkit.Ingress):
        self.assertIn(ingr.name[0], [ingr.name[0] for ingr in chkit.get_ingresses()])
        repl_ingr = chkit.Ingress(
            name=ingr.name[0],
            rules=[chkit.IngressRules(
                path=[chkit.IngressPath(
                    path="/test-path",
                    service_name=svc.name,
                    service_port=svc.ports[1].port
                )]
            )]
        )
        chkit.replace_ingress(repl_ingr)
        got_ingr = chkit.get_ingress(ingr.name[0])
        self.assertEqual(got_ingr.name[0], repl_ingr.name[0])
        self.assertEqual(got_ingr.rules[0].path[0].path, repl_ingr.rules[0].path[0].path)
        self.assertEqual(got_ingr.rules[0].path[0].service_port, repl_ingr.rules[0].path[0].service_port)
        self.assertEqual(got_ingr.rules[0].path[0].service_name, repl_ingr.rules[0].path[0].service_name)
