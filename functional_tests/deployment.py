from functional_tests import chkit
import unittest
import time
import timeout_decorator
import psh.exceptions


class TestDeployment(unittest.TestCase):

    @timeout_decorator.timeout(seconds=630)
    @chkit.test_account
    def test_base(self):
        depl = chkit.Deployment(
            name="functional-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        try:
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            attempts = 1
            while attempts <= 40:
                pods = chkit.get_pods()
                deployment_pods = [pod for pod in pods if pod.deploy == depl.name]
                not_running_pods = [pod for pod in deployment_pods if pod.status.phase != "Running"]
                if len(not_running_pods) == 0 and len(deployment_pods) > 0:
                    break
                time.sleep(15)
                attempts += 1
            self.assertLessEqual(attempts, 40)
        finally:
            chkit.delete_deployment(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    def test_set_image(self, depl: chkit.Deployment):
        chkit.set_image(image="redis", container=depl.containers[0].name, deployment=depl.name)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(got_depl.containers[0].image, "redis")

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    def test_replace_container(self, depl: chkit.Deployment):
        new_container = chkit.Container(
            name=depl.containers[0].name,
            limits=chkit.Resources(cpu=15, memory=15),
            image="redis",
            env=[chkit.EnvVariable("HELLO", "world")],
        )
        chkit.replace_container(deployment=depl.name, container=new_container)
        got_depl = chkit.get_deployment(depl.name)
        needed_containers = [container for container in got_depl.containers if container.name == new_container.name]
        self.assertGreater(len(needed_containers), 0)
        self.assertEqual(needed_containers[0].env, new_container.env)
        self.assertEqual(needed_containers[0].limits.cpu, new_container.limits.cpu)
        self.assertEqual(needed_containers[0].limits.memory, new_container.limits.memory)
        self.assertEqual(needed_containers[0].image, new_container.image)

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    def test_add_container(self, depl: chkit.Deployment):
        new_container = chkit.Container(
            name="additional-container",
            limits=chkit.Resources(cpu=15, memory=15),
            image="redis",
            env=[chkit.EnvVariable("HELLO", "world")],
        )
        chkit.add_container(deployment=depl.name, container=new_container)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(len(got_depl.containers), len(depl.containers)+1)
        needed_containers = [container for container in got_depl.containers if container.name == new_container.name]
        self.assertGreater(len(needed_containers), 0)
        self.assertEqual(needed_containers[0].name, new_container.name)
        self.assertEqual(needed_containers[0].env, new_container.env)
        self.assertEqual(needed_containers[0].limits.cpu, new_container.limits.cpu)
        self.assertEqual(needed_containers[0].limits.memory, new_container.limits.memory)
        self.assertEqual(needed_containers[0].image, new_container.image)

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    @chkit.with_container()
    def test_delete_container(self, depl: chkit.Deployment, container: chkit.Container):
        chkit.delete_container(depl.name, container.name)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(len(got_depl.containers), len(depl.containers))
        for i in range(0, len(depl.containers)):
            self.assertEqual(got_depl.containers[i].name, depl.containers[i].name)

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    def test_set_deploy_replicas(self, depl: chkit.Deployment):
        chkit.set_deployment_replicas(deployment=depl.name, replicas=2)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(depl.name, got_depl.name)
        self.assertEqual(got_depl.replicas, 2)

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    @chkit.with_container()
    def test_change_deploy_version(self, depl: chkit.Deployment, container: chkit.Container):
        got_depl = chkit.get_deployment(depl.name)
        self.assertIn("2.0.0", got_depl.version)

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    @chkit.with_container()
    def test_get_deployment_versions(self, depl: chkit.Deployment, container: chkit.Container):
        deploy_versions = chkit.get_versions(deploy=depl.name)
        self.assertEqual(len(deploy_versions), 2)

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    @chkit.with_container()
    def test_run_deployment_version(self, depl: chkit.Deployment, container: chkit.Container):
        chkit.run_version(deploy=depl.name, version="1.0.0")
        time.sleep(5)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(got_depl.version, "1.0.0")
        self.assertEqual(len(got_depl.containers), len(depl.containers))

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    @chkit.with_container()
    def test_delete_active_deployment_version(self, depl: chkit.Deployment, container: chkit.Container):
        with self.assertRaisesRegex(psh.exceptions.ExecutionError, r".*(\[resource-service-19\]).*"):
            chkit.delete_version(deploy=depl.name, version="2.0.0")

    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    @chkit.with_deployment()
    @chkit.with_container()
    def test_delete_previous_deployment_version(self, depl: chkit.Deployment, container: chkit.Container):
        chkit.delete_version(deploy=depl.name, version="1.0.0")
        depl_versions = chkit.get_versions(depl.name)
        self.assertEqual(len(depl_versions), 1)
        self.assertIn("2.0.0", depl_versions[0].version)
