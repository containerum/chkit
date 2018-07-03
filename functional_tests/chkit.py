from psh import sh
from typing import Dict, Tuple, List
import os
import json
from datetime import datetime

DEFAULT_API_URL = os.getenv("CONTAINERUM_API", "http://api.local.containerum.io")
DATETIME_FORMAT = "%Y-%m-%dT%H:%M:%SZ"

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


class DeploymentStatus(json.JSONEncoder):
    def __init__(self, replicas: int=0, ready_replicas: int=0, available_replicas: int=0,
                 unavailable_replicas: int=0, updated_replicas: int=0):
        self.replicas = replicas
        self.ready_replicas = ready_replicas
        self.available_replicas = available_replicas
        self.unavailable_replicas = unavailable_replicas
        self.updated_replicas = updated_replicas

    @staticmethod
    def json_decode(j):
        return DeploymentStatus(
            replicas=j['replicas'],
            unavailable_replicas=j['unavailable_replicas'],
            available_replicas=j['available_replicas'],
            ready_replicas=j['ready_replicas'],
            updated_replicas=j['updated_replicas'],
        )

    def default(self, o):
        return o.__dict__


class Resources(json.JSONEncoder):
    def __init__(self, cpu: int=0, memory: int=0):
        self.cpu = cpu
        self.memory = memory

    @staticmethod
    def json_decode(j):
        return Resources(
            cpu=j['cpu'],
            memory=j['memory']
        )

    def default(self, o):
        return o.__dict__


class Container(json.JSONEncoder):
    def __init__(self, image: str="", name: str="", limits: Resources=Resources(), env: Dict[str, str]=()):
        self.image = image
        self.name = name
        self.limits = limits
        self.env = env

    @staticmethod
    def json_decode(j):
        return Container(
            name=j['name'],
            image=j['image'],
            limits=Resources.json_decode(j['limits']),
            env=j.get('env')
        )

    def default(self, o):
        return o.__dict__


class Deployment(json.JSONEncoder):
    def __init__(self, created_at: datetime=datetime.now(), status: DeploymentStatus=DeploymentStatus(),
                 containers: List[Container]=(), name: str="", replicas: int=0, total_cpu: int=0, total_memory: int=0,
                 active: bool=False, version: str="0.0.0"):
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
            created_at=datetime.strptime(j['created_at'], DATETIME_FORMAT),
            status=DeploymentStatus.json_decode(j['status']),
            containers=[Container.json_decode(container) for container in j['containers']],
            name=j['name'],
            replicas=j['replicas'],
            total_cpu=j['total_cpu'],
            total_memory=j['total_memory'],
            active=j['active'],
            version=j['version']
        )

    def default(self, o):
        ret = {key: value for key, value in o.__dict__.items() if not isinstance(value, datetime)}
        ret.update({key: value.strftime(DATETIME_FORMAT)
                    for key, value in o.__dict__.items() if isinstance(value, datetime)})
        return ret


def get_deployment(name: str="") -> Deployment:
    output = sh.chkit("get", "deploy", name, "-o", "json").execute().stdout()
    return Deployment.json_decode(json.loads(output))


def get_deployments() -> List[Deployment]:
    output = sh.chkit("get", "deploy", "-o", "json").execute().stdout()
    return [Deployment.json_decode(j) for j in json.loads(output)]


def create_deployment(depl: Deployment, namespace: str=None, file: bool=False) -> None:
    args = ["create", "deployment", "-f"]
    if not file:
        if depl.name is not None:
            args.extend(["--name", depl.name])
        for container in depl.containers:
            args.extend(["--image", f"{container.name}@{container.image}"])
            args.extend(["--cpu", f"{container.name}@{container.limits.cpu}"])
            args.extend(["--memory", f"{container.name}@{container.limits.memory}"])
            for key, value in container.env:
                args.extend(["--env", f"{container.name}@{key}:{value}"])
    else:
        args.extend(["--file", "-"])
    if namespace is not None:
        args.extend(["--namespace", namespace])

    if not file:
        sh.chkit(*args).execute()
    else:
        sh.chkit(*args, _stdin=json.dumps(depl, cls=Deployment)).execute()


def delete_deploy(name: str="", namespace: str=None, concurrency: int=None) -> None:
    args = ["delete", "deploy", "-f", name]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    if concurrency is not None:
        args.extend(["--concurrency", concurrency])
    sh.chkit(*args).execute()


##################
# POD MANAGEMENT #
##################


class PodStatus(json.JSONEncoder):
    def __init__(self, phase: str="", restart_count: int=0, start_at: datetime=None):
        self.phase = phase
        self.restart_count = restart_count
        self.start_at = start_at

    @staticmethod
    def json_decode(j):
        return PodStatus(
            phase=j['phase'],
            restart_count=j['restart_count'],
            start_at=datetime.strptime(j['start_at'], DATETIME_FORMAT) if j['start_at'] is not "" else None
        )

    def default(self, o):
        ret = {key: value for key, value in o.__dict__.items() if not isinstance(value, datetime)}
        ret.update({key: value.strftime(DATETIME_FORMAT)
                    for key, value in o.__dict__.items() if isinstance(value, datetime)})
        return ret


class Pod(json.JSONEncoder):
    def __init__(self, created_at: datetime=datetime.now(), name: str="", owner: str="", containers: List[Container]=(),
                 status: PodStatus=PodStatus(), deploy: str="", total_cpu: int=0, total_memory: int=0):
        self.created_at = created_at
        self.name = name
        self.owner = owner
        self.containers=containers
        self.status = status
        self.deploy = deploy
        self.total_cpu = total_cpu
        self.total_memory = total_memory

    @staticmethod
    def json_decode(j):
        return Pod(
            created_at=datetime.strptime(j['created_at'], DATETIME_FORMAT),
            name=j['name'],
            owner=j['owner'],
            containers=[Container.json_decode(container) for container in j['containers']],
            status=PodStatus.json_decode(j['status']),
            deploy=j['deploy'],
            total_cpu=j['total_cpu'],
            total_memory=j['total_memory']
        )

    def default(self, o):
        ret = {key: value for key, value in o.__dict__.items() if not isinstance(value, datetime)}
        ret.update({key: value.strftime(DATETIME_FORMAT)
                    for key, value in o.__dict__.items() if isinstance(value, datetime)})
        return ret


def get_pods(namespace: str=None, status: str=None) -> List[Pod]:
    args = ["get", "pods", "-o", "json"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    if status is not None:
        args.extend(["--status", status])
    output = sh.chkit(*args).execute().stdout()

    return [Pod.json_decode(j) for j in json.loads(output)]


def set_image(image: str="", container: str="", deployment: str="", namespace: str=None) -> None:
    args = ["set", "image", "--image", image, "--container", container, "--deployment", deployment, "-f"]
    if namespace is not None:
        args.extend(["--namespace", namespace])
    sh.chkit(*args).execute()
