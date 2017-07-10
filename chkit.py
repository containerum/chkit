#!/usr/bin/python3
from client import Client
from version import Version
VERSION = "1.2.0"


def main():
        client = Client(VERSION)
        client.go()


if __name__ == '__main__':
    main()
