import os

path = os.path.dirname(os.path.realpath(__file__))


def get_version():
    with open(path + "/version.py", "rt") as f:
        return f.readline().split("=")[1].strip(' "\n')