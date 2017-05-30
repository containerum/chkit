#!/usr/bin/python3
from client import Client

VERSION = "1.2.4"


def main():
    client = Client(VERSION)
    client.go()


if __name__ == '__main__':
    main()
