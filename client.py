#!/usr/bin/python3
import os
import json
import yaml
import datetime
from data import kinds, output_formats, deployment_json
from parser import *
from tcp_handler import TcpHandler
from api_handler import ApiHandler
from bcolors import BColors
from config_json_handler import get_json_from_config, set_token_to_json_config
from answer_parsers import GetParser
import uuid

config_json_data = get_json_from_config()

class Client:
    def __init__(self):
        self.path = os.getcwd()
        self.version = config_json_data.get("version")
        self.parser = create_parser(kinds, output_formats, self.version)
        uuid_v4 = str(uuid.uuid4())
        self.args = vars(self.parser.parse_args())
        self.debug = self.args.get("debug")
        self.tcp_handler = TcpHandler(uuid_v4, self.args.get("debug"))
        self.api_handler = ApiHandler(uuid_v4)


    def modify_config(self):
        if self.args.get("set_token"):
            set_token_to_json_config(self.args.get("set_token"))

    def logout(self):
        set_token_to_json_config("")
        print("Bye!")

    def go(self):


        self.check_file_existence()
        self.check_arguments()

        if self.args['command'] == 'run':
            self.go_run()

        elif self.args['command'] == 'create':
            self.go_create()

        elif self.args['command'] == 'get':
            self.go_get()

        elif self.args['command'] == 'delete':
            self.go_delete()

        elif self.args['command'] == 'replace':
            self.go_replace()

        elif self.args['command'] == 'config':
            self.modify_config()

        elif self.args['command'] == 'logout':
            self.logout()

    def go_run(self):
        json_to_send = self.construct_run()
        if self.debug:
            self.log_time()

        self.tcp_connect()

        # with open(os.path.join(self.path, 'run_good.json'), 'r', encoding='utf-8') as f:
        #     json_to_send = json.load(f)

        api_result = self.api_handler.run(json_to_send)
        self.handle_api_result(api_result)

        self.get_and_handle_tcp_result('run')

        self.tcp_handler.close()

    def go_create(self):
        self.log_time()
        self.tcp_connect()

        json_to_send = self.get_json_from_file()
        kind = '{}s'.format(json_to_send.get('kind')).lower()

        if kind != 'namespaces':
            api_result = self.api_handler.create(json_to_send)
        else:
            api_result = self.api_handler.create_namespaces(json_to_send)
        self.handle_api_result(api_result)

        self.get_and_handle_tcp_result('create')

        self.tcp_handler.close()

    def go_get(self):
        kind, name = self.construct_get()
        if self.debug:
            self.log_time()
        self.tcp_connect()

        namespace = self.args['namespace']
        if kind != 'namespaces':
            api_result = self.api_handler.get(kind, name, namespace)
        else:
            api_result = self.api_handler.get_namespaces(name)
        self.handle_api_result(api_result)

        self.get_and_handle_tcp_result('get')

        self.tcp_handler.close()

    def get_and_handle_tcp_result(self, command_name, wide=False):
        try:
            tcp_result = self.tcp_handler.receive()
            if command_name == 'get':
                if not tcp_result.get('status') == 'Failure':
                    if self.args.get("debug"):

                        print('{}{}{}'.format(
                            BColors.OKBLUE,
                            'get result:\n',
                            BColors.ENDC
                        ))
                    self.print_result(tcp_result, wide)

            self.print_result_status(tcp_result, command_name)

        except RuntimeError as e:
            print('{}{}{}'.format(
                BColors.FAIL,
                e,
                BColors.ENDC
            ))

    def go_delete(self):
        kind, name = self.construct_delete()

        self.log_time()
        self.tcp_connect()

        self.args['output'] = 'yaml'
        namespace = self.args['namespace']
        if kind != 'namespaces':
            api_result = self.api_handler.delete(kind, name, namespace)
        else:
            api_result = self.api_handler.delete_namespaces(name)
        self.handle_api_result(api_result)

        self.get_and_handle_tcp_result('delete')

        self.tcp_handler.close()

    def go_replace(self):
        self.log_time()
        self.tcp_connect()

        self.args['output'] = 'yaml'
        namespace = self.args['namespace']

        json_to_send = self.get_json_from_file()
        kind = '{}s'.format(json_to_send.get('kind')).lower()

        if kind != 'namespaces':
            api_result = self.api_handler.replace(json_to_send, namespace)
        else:
            api_result = self.api_handler.replace_namespaces(json_to_send)
        self.handle_api_result(api_result)

        self.get_and_handle_tcp_result('replace')

        self.tcp_handler.close()

    def check_file_existence(self):
        if 'file' in self.args:
            if self.args.get('file') is not None:
                if not os.path.isfile(os.path.join(self.path, self.args.get('file'))):
                    self.parser.error('no such file: {}'.format(
                        os.path.join(self.path, self.args.get('file'))
                    ))

    def check_arguments(self):
        if not (self.args.get('version') or self.args.get('help') or self.args.get('command')):
            self.parser.print_help()

    def handle_api_result(self, api_result):
        if api_result.get('id') and self.debug:
            print('{}{}...{} {}OK{}'.format(
                BColors.OKBLUE,
                'api connection',
                BColors.ENDC,
                BColors.BOLD,
                BColors.ENDC
            ))
        elif 'error' in api_result:
            print('{}api error: {}{}'.format(
                BColors.FAIL,
                api_result.get('error'),
                BColors.ENDC
            ))
            self.tcp_handler.close()
            exit()

    def tcp_connect(self):
        try:
            tcp_auth_result = self.tcp_handler.connect()
            if tcp_auth_result.get('ok') and self.debug:
                # print(tcp_auth_result)
                print('{}{}...{} {}OK{}'.format(
                    BColors.OKBLUE,
                    'tcp authorization',
                    BColors.ENDC,
                    BColors.BOLD,
                    BColors.ENDC
                ))
        except RuntimeError as e:
            print('{}{}{}'.format(
                BColors.FAIL,
                e,
                BColors.ENDC
            ))

    def print_result_status(self, result, message):
        if result.get('status') == 'Failure':
            print('{}error: {}{}'.format(
                BColors.FAIL,
                result.get('message'),
                BColors.ENDC
            ))

        elif self.debug:
            print('{}{}...{} {}OK{}'.format(
                BColors.WARNING,
                message,
                BColors.ENDC,
                BColors.BOLD,
                BColors.ENDC
            ))

    def print_result(self, result , wide):
        if self.args['output'] == 'yaml':
            yaml_result = yaml.dump(result, default_flow_style=False)
            print(yaml_result)
        else:

            GetParser(result).show_human_readable_result()

    def log_time(self):
        print('{}{}{}'.format(
            BColors.WARNING,
            str(datetime.datetime.now())[11:19:],
            BColors.ENDC
        ))

    def construct_run(self):
        json_to_send = deployment_json

        json_to_send['metadata']['name'] = self.args['name']
        if self.args['namespace']:
            json_to_send['metadata']['namespace'] = self.args['namespace']
        if self.args['replicas']:
            json_to_send['spec']['replicas'] = self.args['replicas']
        json_to_send['spec']['selector']['matchLabels']['run'] = self.args['name']
        json_to_send['metadata']['labels']['run'] = self.args['name']
        json_to_send['spec']['template']['metadata']['labels']['run'] = self.args['name']

        json_to_send['spec']['template']['metadata']['name'] = self.args['name']
        json_to_send['spec']['template']['spec']['containers'][0]['name'] = self.args['name']
        json_to_send['spec']['template']['spec']['containers'][0]['image'] = self.args['image']
        if self.args['ports']:
            json_to_send['spec']['template']['spec']['containers'][0]['ports'] = [
                {
                    'containerPort': port,
                    'name': '',
                    'protocol': 'TCP'
                }
                for port in self.args['ports']
                ]
        if self.args['env']:
            json_to_send['spec']['template']['spec']['containers'][0]['env'] = [
                {
                    "name": key_value.split('=')[0],
                    "value": key_value.split('=')[1]
                }
                for key_value in self.args['env']]

        with open(os.path.join('/home/gree-gorey/Py/client/client/', 'run.json'), 'w', encoding='utf-8') as w:
            json.dump(json_to_send, w, indent=4)

        return json_to_send

    def get_json_from_file(self):
        file_name = os.path.join(self.path, self.args['file'])
        try:
            with open(file_name, 'r', encoding='utf-8') as f:
                body = json.load(f)
                return body
        except FileNotFoundError:
            self.parser.error('no such file: {}'.format(
                file_name
            ))
        except json.decoder.JSONDecodeError as e:
            self.parser.error('bad json: {}'.format(
                e
            ))

    def construct_delete(self):
        if self.args['file'] and not self.args['kind'] and not self.args['name']:
            body = self.get_json_from_file()
            name = body['metadata']['name']
            kind = '{}s'.format(body['kind'].lower())
            return kind, name

        elif not self.args['file'] and self.args['kind'] and self.args['name']:
            kind = self.args['kind']
            name = self.args['name']
            return kind, name

        elif not self.args['file'] and not self.args['kind']:
            self.parser.error(ONE_REQUIRED_ARGUMENT_ERROR)
        elif self.args['file'] and self.args['kind']:
            self.parser.error(KIND_OR_FILE_BOTH_ERROR)
        elif self.args['file'] and self.args['name']:
            self.parser.error(NAME_OR_FILE_BOTH_ERROR)
        elif self.args['kind'] and not self.args['name']:
            self.parser.error(NAME_WITH_KIND_ERROR)

    def construct_get(self):
        if self.args['file'] and not self.args['kind'] and not self.args['name']:
            body = self.get_json_from_file()
            name = body['metadata']['name']
            kind = '{}s'.format(body['kind'].lower())
            return kind, name

        elif not self.args['file'] and self.args['kind']:
            kind = self.args['kind']
            name = self.args['name']
            return kind, name

        elif not self.args['file'] and not self.args['kind']:
            self.parser.error(ONE_REQUIRED_ARGUMENT_ERROR)
        elif self.args['file'] and self.args['kind']:
            self.parser.error(KIND_OR_FILE_BOTH_ERROR)
        elif self.args['file'] and self.args['name']:
            self.parser.error(NAME_OR_FILE_BOTH_ERROR)


def main():
    client = Client()
    client.go()


if __name__ == '__main__':
    main()
