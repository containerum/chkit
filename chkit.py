#!/usr/bin/python3
from client import Client
from version import Version
VERSION = "1.4.1"


def main():
    v = Version()
    is_ckecked = v.compare_current_version(VERSION)
    if is_ckecked:
        client = Client(VERSION)
        client.go()


if __name__ == '__main__':
    main()
