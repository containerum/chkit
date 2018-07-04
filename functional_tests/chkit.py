from psh import sh
from typing import Dict, Tuple
import os

DEFAULT_API_URL = os.getenv("CONTAINERUM_API", "http://api.local.containerum.io")


def set_api_url(api_url: str=DEFAULT_API_URL, allow_self_signed: bool=False) -> None:
    set_api_args = ["set", "containerum-api", api_url]
    if allow_self_signed:
        set_api_args.append("--allow-self-signed-certs")
    sh.chkit(*set_api_args).execute()


def get_api_url() -> str:
    return sh.chkit("get", "containerum-api").execute().stdout().splitlines()[0]


def login(user: str="test", password: str="test", namespace: str="-") -> None:
    sh.chkit("login", "-u", user, "-p", password, "-n", namespace).execute()


def get_profile() -> Dict[str, str]:
    output = sh.chkit("get", "profile").execute().stdout()
    profile = (line.split(':') for line in output.splitlines())
    profile = (tuple(items) for items in profile if len(items) == 2)
    profile = {key.strip(): value.strip() for key, value in profile}
    return profile


def get_default_namespace() -> Tuple[str, str]:
    output = sh.chkit("get", "default-namespace").execute().stdout()
    line = output.splitlines()[0]
    kv = line.split("/")
    owner_login, namespace_name = kv[0], kv[1]
    return owner_login, namespace_name


def set_default_namespace(namespace: str="-") -> None:
    sh.chkit("set", "default-namespace", "-n", namespace).execute()
