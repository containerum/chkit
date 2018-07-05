from psh import sh
from typing import Dict, Tuple, List
import os
import json
from datetime import datetime
import time
from functional_tests.util import JSONSerialize

DEFAULT_API_URL = os.getenv("CONTAINERUM_API", "http://api.local.containerum.io")
DATETIME_FORMAT = "%Y-%m-%dT%H:%M:%SZ"

TEST_USER = os.getenv("TEST_USER", "helpik94@yandex.com")
TEST_PASSWORD = os.getenv("TEST_USER_PASSWORD", "12345678")
TEST_NAMESPACE = os.getenv("TEST_NAMESPACE", "-")

######################
# API URL MANAGEMENT #
######################


def set_api_url(api_url: str=DEFAULT_API_URL, allow_self_signed: bool=False) -> None:
    set_api_args = ["set", "containerum-api", api_url]
    if allow_self_signed:
        set_api_args.append("--allow-self-signed-certs")
    sh.chkit(*set_api_args).execute()


def get_api_url() -> str:
    return sh.chkit("get", "containerum-api").execute().stdout().splitlines()[0]

#####################
# LOGIN AND PROFILE #
#####################


def login(user: str="test", password: str="test", namespace: str="-") -> None:
    sh.chkit("login", "-u", user, "-p", password, "-n", namespace).execute()


def get_profile() -> Dict[str, str]:
    output = sh.chkit("get", "profile").execute().stdout()
    profile = (line.split(':') for line in output.splitlines())
    profile = (tuple(items) for items in profile if len(items) == 2)
    profile = {key.strip(): value.strip() for key, value in profile}
    return profile


def account(user: str, password: str, namespace: str="-"):
    def decorator(fn):
        def wrapped(*args, **kwargs):
            login(user, password, namespace)
            fn(*args, **kwargs)
        return wrapped
    return decorator


def test_account(fn):
    def wrapped(*args, **kwargs):
        login(user=TEST_USER, password=TEST_PASSWORD, namespace=TEST_NAMESPACE)
        fn(*args, **kwargs)
    return wrapped

#####################
# DEFAULT NAMESPACE #
#####################


def get_default_namespace() -> Tuple[str, str]:
    output = sh.chkit("get", "default-namespace").execute().stdout()
    line = output.splitlines()[0]
    kv = line.split("/")
    owner_login, namespace_name = kv[0], kv[1]
    return owner_login, namespace_name


def set_default_namespace(namespace: str="-") -> None:
    sh.chkit("set", "default-namespace", "-n", namespace).execute()

#########################
# DEPLOYMENT MANAGEMENT #
#########################


class DeploymentStatus:
    def __init__(self, replicas: int=None, ready_replicas: int=None, available_replicas: int=None,
                 unavailable_replicas: int=None, updated_replicas: int=None):
        super().__init__()
        self.replicas = replicas
        self.ready_replicas = ready_replicas
        self.available_replicas = available_replicas
        self.unavailable_replicas = unavailable_replicas
        self.updated_replicas = updated_replicas

    @staticmethod
    def json_decode(j):
        return DeploymentStatus(
            replicas=j.get('replicas'),
            unavailable_replicas=j.get('unavailable_replicas'),
            available_replicas=j.get('available_replicas'),
            ready_replicas=j.get('ready_replicas'),
            updated_replicas=j.get('updated_replicas'),
        )


class Resources:
    def __init__(self, cpu: int=None, memory: int=None):
        super().__init__()
        self.cpu = cpu
        self.memory = memory

    @staticmethod
    def json_decode(j):
        return Resources(
            cpu=j.get('cpu'),
            memory=j.get('memory')
        )


class Container:
    def __init__(self, image: str=None, name: str=None, limits: Resources=None, env: Dict[str, str]=None):
        super().__init__()
        self.image = image
        self.name = name
        self.limits = limits
        self.env = env

    @staticmethod
    def json_decode(j):
        return Container(
            name=j.get('name'),
            image=j.get('image'),
            limits=Resources.json_decode(j.get('limits')),
            env=j.get('env')
        )


class Deployment:
    def __init__(self, created_at: datetime=None, status: DeploymentStatus=None,
                 containers: List[Container]=None, name: str=None, replicas: int=None, total_cpu: int=None,
                 total_memory: int=None, active: bool=False, version: str=None):
        super().__init__()
        self.created_at = created_at
        self.status = status
        self.containers = containers
        self.name = name
        self.replicas = replicas
        self.total_cpu = total_cpu
        self.total_memory = total_memory
        self.active = active
        self.version = version

    @staticmethod
    def json_decode(j):
        return Deployment(
            created_at=datetime.strptime(j.get('created_at'), DATETIME_FORMAT) if j.get('created_at')
            not in (None, '') else None,
            status=DeploymentStatus.json_decode(j.get('status')),
            containers=[Container.json_decode(container) for container in j.get('containers')],
            name=j.get('name'),
            replicas=j.get('replicas'),
            total_cpu=j.get('total_cpu'),
            total_memory=j.get('total_memory'),
            active=j.get('active'),
            version=j.get('version')
        )


def get_deployment(name: str="") -> Deployment:
    output = sh.chkit("get", "deploy", name, "--output", "json").execute().stdout()
    return Deployment.json_decode(json.loads(output))


def get_deployments() -> List[Deployment]:
    output = sh.chkit("get", "deploy", "--output", "json").execute().stdout()
    return [Deployment.json_decode(j) for j in json.loads(output)]


def create_deployment(depl: Deployment, namespace: str=None, file: bool=False) -> None:
    args = ["create", "deployment", "--force"]
    if not file:
        if depl.name is not None:
            args.extend(["--name", depl.name])
        for container in depl.containers:
            args.extend(["--image", f"{container.name}@{container.image}"])
            args.extend(["--cpu", f"{container.name}@{container.limits.cpu}"])
            args.extend(["--memory", f"{container.name}@{container.limits.memory}"])
            if container.env is not None:
                for key, value in container.env.items():
                    args.extend(["--env", f"{container.name}@{key}:{value}"])
    else:
        args.extend(["--file", "-"])
    if namespace is not None:
        args.extend(["--namespace", namespace])

    if not file:
        sh.chkit(*args).execute()
    else:
        sh.chkit(*args, _stdin=json.dumps(depl, cls=JSONSerialize)).execute()


__default_deployment = Deployment(
    name="two-containers-test-depl",
    replicas=1,
    containers=[
        Container(image="nginx", name="first", limits=Resources(cpu=10, memory=10)),
        Container(
            name="second",
            limits=Resources(cpu=15, memory=15),
            image="redis",
            env={"HELLO": "world"},
        )
    ],
)


def with_deployment(deployment: Deployment=__default_deployment, namespace: str=None):
    def decorator(fn):
        def wrapper(*args, **kwargs):
            create_deployment(depl=deployment, namespace=namespace)
            try:
                args = list(args) + [deployment]
                fn(*args, **kwargs)
            finally:
                delete_deployment(name=deployment.name, namespace=namespace)
                time.sleep(5)
        return wrapper
    return decorator


def delete_deployment(name: str, namespace: str=None, concurrency: int=None) -> None:
    args = ["delete", "deploy", "--force", name]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    if concurrency is not None:
        args.extend(["--concurrency", concurrency])
    sh.chkit(*args).execute()


def set_deployment_replicas(deployment: str, replicas: int, namespace: str=None) -> None:
    args = ["set", "replicas", "--deployment", deployment, "--replicas", replicas]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    sh.chkit(*args).execute()


##################################
# DEPLOYMENT VERSIONS MANAGEMENT #
##################################


def get_versions(deploy: str, namespace: str=None) -> List[Deployment]:
    args = ["get", "deployment-versions", "--output", "json", deploy]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    return [Deployment.json_decode(j) for j in json.loads(sh.chkit(*args).execute().stdout())]


def run_version(deploy: str, version: str, namespace: str=None) -> None:
    args = ["run", "deployment-version", "--deployment", deploy, "--version", version, "--force"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    sh.chkit(*args).execute()


def delete_version(deploy: str, version: str, namespace: str=None) -> None:
    args = ["delete", "deployment-version", "--deployment", deploy, "--version", version, "--force"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    sh.chkit(*args).execute()

###################################
# DEPLOYMENT CONTAINER MANAGEMENT #
###################################


def set_image(image: str="", container: str="", deployment: str="", namespace: str=None) -> None:
    args = ["set", "image", "--image", image, "--container", container, "--deployment", deployment, "--force"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    sh.chkit(*args).execute()


def replace_container(deployment: str="", container: Container=Container(), namespace: str=None) -> None:
    args = ["replace", "container", "--container", container.name, "--deployment", deployment, "--force"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    if container.env is not None:
        for k, v in container.env.items():
            args.extend(["--env", f"{k}:{v}"])
    if container.limits is not None:
        if container.limits.cpu is not None:
            args.extend(["--cpu", container.limits.cpu])
        if container.limits.memory is not None:
            args.extend(["--memory", container.limits.memory])
        if container.image is not None:
            args.extend(["--image", container.image])
    sh.chkit(*args).execute()


def add_container(deployment: str="", container: Container=Container(), namespace: str=None, file: bool=False) -> None:
    args = ["create", "container", "--name", container.name, "--deployment", deployment, "--force"]
    if file:
        args.extend(["--file", "-"])
    else:
        args.extend(["--image", container.image, "--memory", container.limits.memory, "--cpu", container.limits.cpu])
        if container.env is not None:
            for k, v in container.env.items():
                args.extend(["--env", f"{k}:{v}"])
    if namespace is not None:
        args.extend(["--namespace", namespace])

    if file:
        sh.chkit(*args, _stdin=json.dumps(container, cls=JSONSerialize)).execute()
    else:
        sh.chkit(*args).execute()


__default_container = Container(
    name="test-container",
    limits=Resources(cpu=15, memory=15),
    image="redis",
    env={"HELLO": "world"},
)


def with_container(container: Container=__default_container, deployment: str=__default_deployment.name,
                   namespace: str=None):
    def decorator(fn):
        def wrapper(*args, **kwargs):
            add_container(deployment=deployment, container=container, namespace=namespace)
            args = list(args)+[container]
            fn(*args, **kwargs)
        return wrapper
    return decorator


def delete_container(deployment: str="", container: str="", namespace: str=None) -> None:
    args = ["delete", "container", "--deployment", deployment, "--container", container, "--force"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    sh.chkit(*args).execute()


##################
# POD MANAGEMENT #
##################


class PodStatus:
    def __init__(self, phase: str=None, restart_count: int=None, start_at: datetime=None):
        super().__init__()
        self.phase = phase
        self.restart_count = restart_count
        self.start_at = start_at

    @staticmethod
    def json_decode(j):
        return PodStatus(
            phase=j.get('phase'),
            restart_count=j.get('restart_count'),
            start_at=datetime.strptime(j.get('start_at'), DATETIME_FORMAT) if j.get('start_at')
            not in (None, '') else None
        )


class Pod:
    def __init__(self, created_at: datetime=None, name: str=None, owner: str=None, containers: List[Container]=None,
                 status: PodStatus=None, deploy: str=None, total_cpu: int=None, total_memory: int=None):
        super().__init__()
        self.created_at = created_at
        self.name = name
        self.owner = owner
        self.containers = containers
        self.status = status
        self.deploy = deploy
        self.total_cpu = total_cpu
        self.total_memory = total_memory

    @staticmethod
    def json_decode(j):
        return Pod(
            created_at=datetime.strptime(j.get('created_at'), DATETIME_FORMAT) if j.get('created_at')
            not in (None, '') else None,
            name=j.get('name'),
            owner=j.get('owner'),
            containers=[Container.json_decode(container) for container in j.get('containers')],
            status=PodStatus.json_decode(j.get('status')),
            deploy=j.get('deploy'),
            total_cpu=j.get('total_cpu'),
            total_memory=j.get('total_memory')
        )


def get_pods(namespace: str=None, status: str=None) -> List[Pod]:
    args = ["get", "pods", "--output", "json"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    if status is not None:
        args.extend(["--status", status])
    output = sh.chkit(*args).execute().stdout()

    return [Pod.json_decode(j) for j in json.loads(output)]


class PodWaitException(Exception):
    pass


def ensure_pods_running(deployment: str=__default_deployment.name, max_attempts: int=40, sleep_seconds: float=15,
                        exception_on_fail: Exception=PodWaitException):
    def decorator(fn):
        def wrapper(*args, **kwargs):
            attempts = 1
            while attempts <= max_attempts:
                pods = get_pods()
                not_running_pods = [pod for pod in pods if pod.deploy == deployment and pod.status.phase != "Running"]
                if len(not_running_pods) == 0:
                    break
                time.sleep(sleep_seconds)
                attempts += 1
            if attempts > max_attempts:
                raise exception_on_fail
            fn(*args, **kwargs)
        return wrapper
    return decorator


def pod_logs(pod: str, container: str=None, tail: int=None, namespace: str=None) -> List[str]:
    args = ["logs", pod]
    if container is not None:
        args.append(container)
    if tail is not None:
        args.extend(["--tail", tail])
    if namespace is not None:
        args.extend(["--namespace", namespace])

    return sh.chkit(*args).execute().stdout().splitlines()


#######################
# SERVICES MANAGEMENT #
#######################


class ServicePort:

    def __init__(self, name: str, target_port: int, protocol: str="TCP", port: int=None):
        super().__init__()
        self.name = name
        self.target_port = target_port
        self.protocol = protocol
        self.port = port

    @staticmethod
    def json_decode(j):
        return ServicePort(
            name=j.get("name"),
            target_port=j.get("target_port"),
            protocol=j.get("protocol"),
            port=j.get("port"),
        )


class Service:

    def __init__(self, name: str, deploy: str, ports: List[ServicePort], ips: List[str]=None, domain: str=None):
        super().__init__()
        self.name = name
        self.deploy = deploy
        self.ports = ports
        self.ips = ips
        self.domain = domain

    @staticmethod
    def json_decode(j):
        return Service(
            name=j.get("name"),
            deploy=j.get("deploy"),
            ports=[ServicePort.json_decode(port) for port in j.get("ports")],
            ips=j.get("ips"),
            domain=j.get("domain")
        )

    def is_external(self):
        return self.domain is not None or self.ips is not None


def create_service(service: Service, file: bool=False, namespace: str=None) -> None:
    args = ["create", "service", "--name", service.name, "--deploy", service.deploy, "--force"]
    if file:
        args.extend(["--input", "json"])
    else:
        port = service.ports[0]
        args.extend(["--port-name", port.name, "--target-port", port.target_port, "--protocol", port.protocol])
        if port.port is not None:
            args.extend(["--port", port.port])

    if namespace is not None:
        args.extend(["--namespace", namespace])

    if file:
        service_to_create = Service(name=service.name, deploy=service.deploy, ports=service.ports)
        sh.chkit(*args, _stdin=json.dumps(service_to_create, cls=JSONSerialize)).execute()
    else:
        sh.chkit(*args).execute()


def get_services(solution: str=None, namespace: str=None) -> List[Service]:
    args = ["get", "svc", "-o", "json"]
    if solution is not None:
        args.extend(["--solution-name", solution])
    if namespace is not None:
        args.extend(["--namespace", namespace])

    return [Service.json_decode(j) for j in json.loads(sh.chkit(*args).execute().stdout())]


def get_service(service: str, solution: str=None, namespace: str=None) -> Service:
    args = ["get", "svc", service, "-o", "json"]
    if solution is not None:
        args.extend(["--solution-name", solution])
    if namespace is not None:
        args.extend(["--namespace", namespace])

    return Service.json_decode(json.loads(sh.chkit(*args).execute().stdout()))


def replace_service(service: Service, file: bool=False, namespace: str=None) -> None:
    args = ["replace", "service", "--name", service.name, "--deploy", service.deploy, "--force"]
    if file:
        args.extend(["--input", "json"])
    else:
        port = service.ports[0]
        if port.name is not None:
            args.extend(["--port-name", port.name])
        if port.target_port is not None:
            args.extend(["--target-port", port.target_port])
        if port.target_port is not None:
            args.extend(["--protocol", port.protocol])
        if port.port is not None:
            args.extend(["--port", port.port])

    if namespace is not None:
        args.extend(["--namespace", namespace])

    if file:
        service_to_create = Service(name=service.name, deploy=service.deploy, ports=service.ports)
        sh.chkit(*args, _stdin=json.dumps(service_to_create, cls=JSONSerialize)).execute()
    else:
        sh.chkit(*args).execute()


def delete_service(service: str, namespace: str=None) -> None:
    args = ["delete", "service", service]
    if namespace is not None:
        args.extend(["--namespace", namespace])

    sh.chkit(*args).execute()


def with_service(service: Service, namespace: str=None):
    def decorator(fn):
        def wrapper(*args, **kwargs):
            create_service(service, file=True, namespace=namespace)
            try:
                args = list(args)+[service]
                fn(*args, **kwargs)
            finally:
                delete_service(service.name)
        return wrapper
    return decorator
