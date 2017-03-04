import socket
import json
import os
from bcolors import BColors


class TcpHandler:
    def __init__(self):
        self.TCP_IP = '146.185.135.181'
        self.TCP_PORT = 3000
        self.BUFFER_SIZE = 1024
        self.AUTH_FORM = {
            "channel": "default",
            "login": "guest",
            "token": read_token()
        }

        self.s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    def connect(self):
        self.s.connect((self.TCP_IP, self.TCP_PORT))
        self.s.send((json.dumps(self.AUTH_FORM) + '\n').encode('utf-8'))
        data = self.s.recv(self.BUFFER_SIZE)
        if not data:
            raise RuntimeError("tcp socket connection broken")
        else:
            result = json.loads(data.decode('utf-8'))

        return result

    def receive(self):
        data = ''
        while data[-1::] != '\n':
            received = self.s.recv(self.BUFFER_SIZE).decode('utf-8')
            if not received:
                raise RuntimeError("tcp socket connection broken")
            print('{}tcp received {} bytes...{}'.format(
                BColors.OKBLUE,
                len(received),
                BColors.ENDC
            ))
            data += received
            # print(len(data))
            # print(data[-1::] == '\n')
        # data = self.s.recv(self.BUFFER_SIZE)
        # print(data)

        try:
            result = json.loads(data)
            print('{}{}...{} {}OK{}'.format(
                BColors.OKBLUE,
                'tcp complete',
                BColors.ENDC,
                BColors.BOLD,
                BColors.ENDC
            ))
        except:
            with open('received_str', 'w', encoding='utf-8') as w:
                w.write(data.decode('utf-8'))
            result = {}
        # print(json.dumps(result, indent=4))
        # pretty_json = json.dumps(result, indent=4)

        return result

    def close(self):
        self.s.close()


def read_token():
    path = os.path.dirname(os.path.realpath(__file__))
    file_name = os.path.join(path, 'tcp_token')
    with open(file_name, 'r', encoding='utf-8') as f:
        return f.read()
