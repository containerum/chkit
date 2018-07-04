from functional_tests import chkit
import unittest
import time


class TestDeployment(unittest.TestCase):

    def test_base(self):
        depl = chkit.Deployment(
            name="functional-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        try:
            chkit.login(user="helpik94@yandex.com", password="12345678")
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            attempts: int
            for i in range(1, 40):
                attempts = i
                pods = chkit.get_pods()
                not_running_pods = [pod for pod in pods if pod.deploy == depl.name and pod.status.phase != "Running"]
                if len(not_running_pods) == 0:
                    break
                time.sleep(15)
            self.assertLessEqual(attempts, 40)
        finally:
            chkit.delete_deploy(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    def test_set_image(self):
        depl = chkit.Deployment(
            name="set-image-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        try:
            chkit.login(user="helpik94@yandex.com", password="12345678")
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            chkit.set_image(image="redis", container=depl.containers[0].name, deployment=depl.name)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(got_depl.containers[0].image, "redis")
        finally:
            chkit.delete_deploy(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    def test_replace_container(self):
        depl = chkit.Deployment(
            name="replace-container-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        try:
            chkit.login(user="helpik94@yandex.com", password="12345678")
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            new_container = chkit.Container(
                name=depl.containers[0].name,
                limits=chkit.Resources(cpu=15, memory=15),
                image="redis",
                env={"HELLO": "world"},
            )
            chkit.replace_container(deployment=depl.name, container=new_container)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(got_depl.containers[0].name, new_container.name)
            self.assertEqual(got_depl.containers[0].env, new_container.env)
            self.assertEqual(got_depl.containers[0].limits.cpu, new_container.limits.cpu)
            self.assertEqual(got_depl.containers[0].limits.memory, new_container.limits.memory)
            self.assertEqual(got_depl.containers[0].image, new_container.image)
        finally:
            chkit.delete_deploy(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    def test_add_container(self):
        depl = chkit.Deployment(
            name="add-container-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        try:
            chkit.login(user="helpik94@yandex.com", password="12345678")
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            new_container = chkit.Container(
                name="second",
                limits=chkit.Resources(cpu=15, memory=15),
                image="redis",
                env={"HELLO": "world"},
            )
            chkit.add_container(deployment=depl.name, container=new_container)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(len(got_depl.containers), 2)
            needed_containers = [container for container in got_depl.containers if container.name == new_container.name]
            self.assertGreater(len(needed_containers), 0)
            self.assertEqual(needed_containers[0].name, new_container.name)
            self.assertEqual(needed_containers[0].env, new_container.env)
            self.assertEqual(needed_containers[0].limits.cpu, new_container.limits.cpu)
            self.assertEqual(needed_containers[0].limits.memory, new_container.limits.memory)
            self.assertEqual(needed_containers[0].image, new_container.image)
        finally:
            chkit.delete_deploy(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    def test_delete_container(self):
        depl = chkit.Deployment(
            name="del-container-test-depl",
            replicas=1,
            containers=[
                chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10)),
                chkit.Container(
                    name="second",
                    limits=chkit.Resources(cpu=15, memory=15),
                    image="redis",
                    env={"HELLO": "world"},
                )
            ],
        )
        try:
            chkit.login(user="helpik94@yandex.com", password="12345678")
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            chkit.delete_container(depl.name, depl.containers[1].name)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(len(got_depl.containers), 1)
            self.assertEqual(got_depl.containers[0].name, depl.containers[0].name)
        finally:
            chkit.delete_deploy(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    def test_set_deploy_replicas(self):
        depl = chkit.Deployment(
            name="add-container-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        try:
            chkit.login(user="helpik94@yandex.com", password="12345678")
            chkit.create_deployment(depl)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            self.assertEqual(depl.replicas, got_depl.replicas)
            chkit.set_deploy_replicas(deploy=depl.name, replicas=2)
            got_depl = chkit.get_deployment(depl.name)
            self.assertEqual(depl.name, got_depl.name)
            self.assertEqual(got_depl.replicas, 2)
        finally:
            chkit.delete_deploy(name=depl.name)
            time.sleep(5)
            self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])
