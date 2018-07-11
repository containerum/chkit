import unittest
from timeout_decorator import timeout_decorator
import functional_tests.chkit as chkit
import time
import psh


class TestSolution(unittest.TestCase):
    @timeout_decorator.timeout(seconds=30)
    @chkit.test_account
    def test_base(self):
        templates = chkit.get_templates()
        random_template = templates[0].name
        chkit.get_template_envs(random_template)

        try:
            chkit.run_solution(random_template, "test-solution")
            time.sleep(1)
            self.assertIn("test-solution", [sol.name for sol in chkit.get_solutions()])

            got_sol = chkit.get_solution("test-solution")
            self.assertEqual("test-solution", got_sol.name)
            self.assertEqual(random_template, got_sol.template)

            time.sleep(5)
            svc = chkit.get_services(solution="test-solution")
            if len(svc) == 0:
                raise LookupError("No services created")
            self.assertEqual(svc[0].solution, "test-solution")
            deploy = chkit.get_deployments(solution="test-solution")
            if len(deploy) == 0:
                raise LookupError("No deployments created")
            self.assertEqual(deploy[0].solution, "test-solution")

        finally:
            chkit.delete_solution("test-solution")
            time.sleep(5)
            self.assertNotIn("test-solution", [sol.name for sol in chkit.get_solutions()])

            try:
                chkit.get_deployments(solution="test-solution")
            except psh.exceptions.ExecutionError as error:
                if error.stderr() != "[Solutions-12] Not Found Solution with this name doesn't exist":
                    raise error
            else:
                raise LookupError("Deployments still not deleted")

            try:
                chkit.get_services(solution="test-solution")
            except psh.exceptions.ExecutionError as error:
                if error.stderr() != "[Solutions-12] Not Found Solution with this name doesn't exist":
                    raise error
            else:
                raise LookupError("Services still not deleted")
