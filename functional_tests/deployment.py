from functional_tests import chkit
import unittest
import time


class TestDeployment(unittest.TestCase):

    def test_base(self):
        chkit.login(user="helpik94@yandex.com", password="12345678")
        depl = chkit.Deployment(
            name="functional-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
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
        chkit.delete_deploy(depl.name)
        time.sleep(5)
        self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])

    def test_set_image(self):
        chkit.login(user="helpik94@yandex.com", password="12345678")
        depl = chkit.Deployment(
            name="set-image-test-depl",
            replicas=1,
            containers=[chkit.Container(image="nginx", name="first", limits=chkit.Resources(cpu=10, memory=10))],
        )
        chkit.create_deployment(depl)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(depl.name, got_depl.name)
        chkit.set_image(image="redis", container=depl.containers[0].name, deployment=depl.name)
        got_depl = chkit.get_deployment(depl.name)
        self.assertEqual(got_depl.containers[0].image, "redis")
        chkit.delete_deploy(depl.name)
        time.sleep(5)
        self.assertNotIn(depl.name, [deploy.name for deploy in chkit.get_deployments()])
