import unittest
import functional_tests.chkit as chkit
from timeout_decorator import timeout_decorator
import time


class TestPod(unittest.TestCase):

    __test_deployment = chkit.Deployment(
        name="pod-log-test",
        replicas=1,
        containers=[
            chkit.Container(image="twentydraft/shibainfo", name="shiba", limits=chkit.Resources(cpu=50, memory=50))
        ]
    )

    @timeout_decorator.timeout(seconds=650)
    @chkit.test_account
    @chkit.with_deployment(deployment=__test_deployment)
    @chkit.ensure_pods_running(deployment=__test_deployment.name)
    def test_pod_logs(self, depl: chkit.Deployment):
        pod = [pod for pod in chkit.get_pods() if pod.deploy == depl.name][0]
        time.sleep(30)
        log_lines = chkit.pod_logs(pod=pod.name, tail=10)
        print("\ngot log lines:")
        print("\n".join(log_lines))
        self.assertGreaterEqual(len(log_lines), 3)
        self.assertLessEqual(len(log_lines), 10)

