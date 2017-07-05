import os
import json
import yaml
import re
from cmd_arg_parser import create_parser
from tcp_handler import TcpHandler, check_http_status
from api_handler import ApiHandler
from bcolors import BColors
from getpass import getpass
from config_json_handler import get_json_from_config, set_token_to_json_config,set_default_namespace_to_json_config,\
    show_namespace_token_from_config
from answer_parsers import TcpApiParser
from constructors import Constructor
import uuid
from datetime import datetime
from hashlib import sha256, md5


config_json_data = get_json_from_config()


class Client(Constructor):
    def __init__(self, version):
        self.path = os.getcwd()
        self.version = version
        self.parser = create_parser(self.version)
        uuid_v4 = str(uuid.uuid4())
        self.args = vars(self.parser.parse_args())
        self.debug = self.args.get("debug")
        self.tcp_handler = TcpHandler(uuid_v4, self.args.get("debug"))
        self.api_handler = ApiHandler(uuid_v4)
        Constructor.__init__(self)

    def go_config(self):
            if self.args.get("set_token"):
                set_token_to_json_config(self.args.get("set_token"))
            elif self.args.get("set_default_namespace"):
                if not self.test_namespace(self.args.get("set_default_namespace")):
                    return
                set_default_namespace_to_json_config(self.args.get("set_default_namespace"))
            else:
                show_namespace_token_from_config()

    @staticmethod
    def logout():
        set_token_to_json_config("")
        print("Bye!")

    def login(self):
        email_regex = r"(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)"
        is_valid = re.compile(email_regex)
        try:
            email = input('Enter your email: ')
        except KeyboardInterrupt:
            return
        if not is_valid.findall(email):
            print("Email is not valid")
            return
        try:
            pwd = getpass()
        except KeyboardInterrupt:
            return
        json_to_send = {"username": email, "password": md5((email+pwd).encode()).hexdigest()}
        if self.args.get("debug"):
            self.log_time()
        self.tcp_connect()
        api_result = self.api_handler.login(json_to_send)
        if 'ok' in api_result:
            set_token_to_json_config(api_result['token'])
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('post')
        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

    def go(self):
        self.check_file_existence()
        self.check_arguments()

        if self.args.get("kind") in ("deployments", "deploy", "deployment"):
            self.args["kind"] = "deployments"
        elif self.args.get("kind") in ("po", "pods", "pod"):
            self.args["kind"] = "pods"
        elif self.args.get("kind") in ("service", "services", "svc"):
            self.args["kind"] = "services"
        else:
            self.args["kind"] = "namespaces"

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
            self.go_config()

        elif self.args['command'] == 'expose':
            self.go_expose()

        elif self.args['command'] == 'login':
            self.login()

        elif self.args['command'] == 'logout':
            self.logout()

        elif self.args['command'] == 'set':
            self.go_set()

        elif self.args['command'] == 'restart':
            self.go_restart()

        elif self.args['command'] == 'scale':
            self.go_scale()

    def go_restart(self):
        self.log_time()
        self.tcp_connect()

        namespace = self.args['namespace']
        if not namespace:
            namespace = config_json_data.get("default_namespace")
        api_result = self.api_handler.delete("deployments", self.args["name"], namespace, True)
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('restart')
        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

    def go_scale(self):
        if self.args.get("debug"):
            self.log_time()
        self.tcp_connect()

        namespace = self.args.get('namespace')
        if not namespace:
            namespace = config_json_data.get("default_namespace")

        count = self.args.get("count")
        try:
            replicas_count = int(count)
            json_to_send = {"replicas": replicas_count}
            api_result = self.api_handler.scale(json_to_send, self.args.get("name"), namespace)
        except (ValueError, TypeError):
            print('{}{}{} {}'.format(
                BColors.FAIL,
                "Error: ",
                "Count is not integer",
                BColors.ENDC,
            ))
            return
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('scale')
        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

    def go_set(self):
        if self.args.get("debug"):
            self.log_time()
        self.tcp_connect()

        namespace = self.args.get('namespace')
        if not namespace:
            namespace = config_json_data.get("default_namespace")

        args = self.args.get("args")
        if args:
            if '=' in args:
                container_name, image = args.split('=')
                json_to_send = {"name": self.args.get("name"), "image": image}
                api_result = self.api_handler.set(json_to_send, container_name, namespace)
            else:
                try:
                    replicas_count = int(args)
                    json_to_send = {"replicas": replicas_count}
                    api_result = self.api_handler.set(json_to_send, self.args.get("name"), namespace)
                except (ValueError, TypeError):
                    print('{}{}{} {}'.format(
                        BColors.FAIL,
                        "Error: ",
                        "Count is not integer",
                        BColors.ENDC,
                    ))
                    return

            if not self.handle_api_result(api_result):
                return

            json_result = self.get_and_handle_tcp_result('set')
            self.tcp_handler.close()
            if not check_http_status(json_result, self.args.get("command")):
                return
        else:
            print('{}{}{} {}'.format(
                BColors.FAIL,
                "Error: ",
                "Empty args",
                BColors.ENDC,
            ))
            return

    def go_run(self):
        json_to_send = self.construct_run()
        if not json_to_send:
            return
        if self.debug:
            self.log_time()

        self.tcp_connect()
        namespace = self.args.get('namespace')
        if not namespace:
            namespace = config_json_data.get("default_namespace")
        api_result = self.api_handler.run(json_to_send, namespace)
        if not self.handle_api_result(api_result):
            return
        json_result = self.get_and_handle_tcp_result('run')

        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

    def go_expose(self):
        namespace = self.args.get('namespace')
        if not namespace:
            namespace = config_json_data.get("default_namespace")

        json_to_send = self.construct_expose(namespace)
        if self.debug:
            self.log_time()

        if not json_to_send:
            return
        self.tcp_connect()

        api_result = self.api_handler.expose(json_to_send, namespace)
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('expose')
        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

    def go_create(self):
        if self.args.get("debug"):
            self.log_time()
        self.tcp_connect()

        json_to_send = self.get_json_from_file()
        namespace = json_to_send["metadata"].get("namespace")
        if not namespace:
            namespace = self.args.get('namespace')
            if not namespace:
                namespace = config_json_data.get("default_namespace")

        api_result = self.api_handler.create(json_to_send, namespace)
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('create')
        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

    def test_namespace(self, namespace):
        if self.debug:
            self.log_time()
        self.tcp_connect()

        api_result = self.api_handler.get_namespaces(namespace)

        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('check namespace')
        self.tcp_handler.close()
        if not check_http_status(json_result, "check namespace"):
            return
        return True

    def go_get(self):
        kind, name = self.construct_get()
        if self.debug:
            self.log_time()
        self.tcp_connect()

        self.namespace = self.args.get('namespace')
        if not self.namespace:
            self.namespace = config_json_data.get("default_namespace")

        if kind == "namespaces":
            if self.args.get("name"):
                api_result = self.api_handler.get_namespaces(self.args.get("name"))
            else:
                api_result = self.api_handler.get_namespaces()
        else:
            api_result = self.api_handler.get(kind, name, self.namespace)
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('get')
        self.tcp_handler.close()
        if not check_http_status(json_result, "get"):
            return
        return json_result

    def get_and_handle_tcp_result(self, command_name):
        try:
            tcp_result = self.tcp_handler.receive()
            if command_name == 'get':
                if not tcp_result.get('error'):
                    if self.args.get("debug"):

                        print('{}{}{}'.format(
                            BColors.OKBLUE,
                            'get result:\n',
                            BColors.ENDC
                        ))
                    self.print_result(tcp_result)

            return tcp_result

        except RuntimeError as e:
            print('{}{}{}'.format(
                BColors.FAIL,
                e,
                BColors.ENDC
            ))
            return None

    def go_delete(self):
        kind, name = self.construct_delete()

        self.log_time()
        self.tcp_connect()

        self.args['output'] = 'yaml'
        namespace = self.args['namespace']
        if not namespace:
            namespace = config_json_data.get("default_namespace")
        if kind != 'namespaces':
            api_result = self.api_handler.delete(kind, name, namespace, self.args.get("pods"))
        else:
            api_result = self.api_handler.delete_namespaces(name)
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('delete')
        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

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
        if not self.handle_api_result(api_result):
            return

        json_result = self.get_and_handle_tcp_result('replace')

        self.tcp_handler.close()
        if not check_http_status(json_result, self.args.get("command")):
            return

    def check_file_existence(self):
        if 'file' in self.args:
            if self.args.get('file'):
                if not os.path.isfile(os.path.join(self.path, self.args.get('file'))):
                    self.parser.error('no such file: {}'.format(
                        os.path.join(self.path, self.args.get('file'))
                    ))

    def check_arguments(self):
        if not (self.args.get('version') or self.args.get('help') or self.args.get('command')):
            self.parser.print_help()

    def handle_api_result(self, api_result):
        if api_result.get('id'):
            if self.debug:
                print('{}{}...{} {}OK{}'.format(
                    BColors.OKBLUE,
                    'api connection',
                    BColors.ENDC,
                    BColors.BOLD,
                    BColors.ENDC
                ))
                print('{}{}{}{}'.format(
                    BColors.OKGREEN,
                    'Command Id: ',
                    api_result.get('id'),
                    BColors.ENDC,

                ))
            return True
        elif 'error' in api_result:
            print('{}api error: {}{}'.format(
                BColors.FAIL,
                api_result.get('error'),
                BColors.ENDC
            ))
            self.tcp_handler.close()
            return

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

    def print_result(self, result):
        if self.args.get("command") != "expose":
            if self.args.get('output') == 'yaml':
                result = result["results"]
                print(yaml.dump(result, default_flow_style=False))
            elif self.args['output'] == 'json':
                result = result["results"]
                print(json.dumps(result, indent=4))
            else:
                deploy = self.args.get("deploy")
                TcpApiParser(result, deploy=deploy)

    def log_time(self):
        if self.args["debug"]:
            print('{}{}{}'.format(
                BColors.WARNING,
                str(datetime.now())[11:19:],
                BColors.ENDC
            ))




