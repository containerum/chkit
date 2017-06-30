import re
import os
from bcolors import BColors
from data import deployment_json, service_json
from keywords import LOWER_CASE_ERROR, NO_IMAGE_AND_CONFIGURE_ERROR, JSON_TEMPLATES_EXPOSE_FILE, JSON_TEMPLATES_RUN_FILE
from run_image import RunImage
from run_configure import RunConfigure
import json
import yaml
from parser import ONE_REQUIRED_ARGUMENT_ERROR, KIND_OR_FILE_BOTH_ERROR, NAME_OR_FILE_BOTH_ERROR, NAME_WITH_KIND_ERROR
from hashlib import sha256, md5
from datetime import datetime


class Constructor:

    def construct_run(self):
        json_to_send = deployment_json

        if not self.args["name"].islower():
            e = LOWER_CASE_ERROR
            print('{}{}{} {}'.format(
                BColors.FAIL,
                "Error: ",
                e,
                BColors.ENDC,
            ))
            return
        name_check = r"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
        is_valid = re.compile(name_check)
        if not is_valid.findall(self.args['name']):
            print('{}{}{} {}'.format(
                BColors.FAIL,
                "Error: ",
                ValueError("Deploy name must consist of lower case alphanumeric characters or '-', and must start and"
                           " end with an alphanumeric character"),
                BColors.ENDC,
            ))
            return
        json_to_send['metadata']['name'] = self.args['name']

        if self.args["configure"] and not self.args["image"]:
            runconfigure = RunConfigure()
            param_dict = runconfigure.get_data_from_console()

        elif self.args["image"] and not self.args["configure"]:
            runimage = RunImage()
            param_dict = runimage.parse_data(self.args)

        if not param_dict:
            return
        image = param_dict["image"]
        ports = param_dict["ports"]
        labels = param_dict["labels"]
        env = param_dict["env"]
        cpu = param_dict["cpu"]
        memory = param_dict["memory"]
        replicas = param_dict["replicas"]
        commands = param_dict["commands"]

        if not self.args["configure"] and not self.args["image"]:
            self.parser.error(NO_IMAGE_AND_CONFIGURE_ERROR)
            return

        json_to_send['spec']['replicas'] = replicas
        json_to_send['spec']['template']['spec']['containers'][0]['name'] = self.args['name']
        json_to_send['spec']['template']['spec']['containers'][0]['image'] = image
        if commands:
            json_to_send['spec']['template']['spec']['containers'][0]['command'] = commands
        if ports:
            json_to_send['spec']['template']['spec']['containers'][0]['ports'] = []
            for port in ports:
                json_to_send['spec']['template']['spec']['containers'][0]['ports'].append({
                    'containerPort': port
                })
        if labels:
            for key, value in labels.items():
                json_to_send['metadata']['labels'].update({key: value})
                json_to_send['spec']['template']['metadata']['labels'].update({key: value})
        if env:
            json_to_send['spec']['template']['spec']['containers'][0]['env'] = [
                {
                    "name": key,
                    "value": value
                }
                for key, value in env.items()]
        json_to_send['spec']['template']['spec']['containers'][0]['resources']["requests"]['cpu'] = cpu
        json_to_send['spec']['template']['spec']['containers'][0]['resources']["requests"]['memory'] = memory
        with open(os.path.join(os.getenv("HOME") + "/.containerum/src/", JSON_TEMPLATES_RUN_FILE), 'w',
                  encoding='utf-8') as w:
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
            pass

        try:
            with open(file_name, 'r', encoding='utf-8') as f:
                body = yaml.load(f)
                return body
        except FileNotFoundError:
            self.parser.error('no such file: {}'.format(
                file_name
            ))
        except yaml.YAMLError as e:
            self.parser.error('bad json or yaml: {}'.format(
                e
            ))

    def construct_delete(self):
        if self.args['file'] and self.args['kind'] == "namespaces" and not self.args['name']:
            body = self.get_json_from_file()
            name = body['metadata']['name']
            kind = '{}s'.format(body['kind'].lower())
            return kind, name

        elif not self.args['file'] and self.args['kind'] != "namespaces" and self.args['name']:
            kind = self.args['kind']
            name = self.args['name']
            return kind, name

        elif not self.args['file'] and self.args['kind'] == "namespaces":
            self.parser.error(ONE_REQUIRED_ARGUMENT_ERROR)
        elif self.args['file'] and self.args['kind'] == "namespaces":
            self.parser.error(KIND_OR_FILE_BOTH_ERROR)
        elif self.args['file'] and self.args['name']:
            self.parser.error(NAME_OR_FILE_BOTH_ERROR)
        elif self.args['kind'] != "namespaces" and not self.args['name']:
            self.parser.error(NAME_WITH_KIND_ERROR)

    def construct_get(self):
        if self.args.get('file') and not self.args.get('kind') and not self.args.get('name'):
            body = self.get_json_from_file()
            name = body['metadata']['name']
            kind = '{}s'.format(body['kind'].lower())
            return kind, name

        elif not self.args.get('file') and self.args.get('kind'):
            kind = self.args['kind']
            name = self.args.get('name')
            return kind, name

        elif not self.args['file'] and not self.args['kind']:
            self.parser.error(ONE_REQUIRED_ARGUMENT_ERROR)
        elif self.args['file'] and self.args['kind']:
            self.parser.error(KIND_OR_FILE_BOTH_ERROR)
        elif self.args['file'] and self.args['name']:
            self.parser.error(NAME_OR_FILE_BOTH_ERROR)

    def construct_expose(self, namespace):
        labels = {}
        is_external = {"external": "true"}
        json_to_send = service_json
        ports = self.args.get("ports")
        self.args["kind"] = "deployments"
        if ports:
            for p in ports:
                p = p.split(":")
                if len(p) == 3:
                    if p[2].upper() == "TCP" or p[2].upper() == "UDP":
                        json_to_send["spec"]["ports"].append(
                            {"name": p[0], "protocol": p[2].upper(), "targetPort": int(p[1])})
                    else:
                        json_to_send["spec"]["ports"].append(
                            {"name": p[0], "protocol": "TCP", "port": int(p[2]), "targetPort": int(p[1])})
                        is_external["external"] = "false"
                if len(p) == 4:
                    json_to_send["spec"]["ports"].append({"name": p[0], "port": int(p[2]), "protocol": p[3].upper(),
                                                          "targetPort": int(p[1])})
                    is_external["external"] = "false"
                elif len(p) == 2:
                    json_to_send["spec"]["ports"].append({"name": p[0], "protocol": "TCP", "targetPort": int(p[1])})
        result = self.go_get()
        if not result:
            return
        namespace_hash = sha256(namespace.encode('utf-8')).hexdigest()[:32]
        labels.update({namespace_hash: self.args.get("name")})
        json_to_send["metadata"]["labels"].update(labels)

        json_to_send["metadata"]["labels"].update(is_external)
        json_to_send["metadata"]["name"] = self.args["name"] + "-" + \
                                           md5((self.args.get("name") + str(datetime.now()))
                                               .encode("utf-8")).hexdigest()[:4]
        json_to_send["spec"]["selector"].update(labels)
        with open(os.path.join(os.getenv("HOME") + "/.containerum/src/", JSON_TEMPLATES_EXPOSE_FILE), 'w',
                  encoding='utf-8') as w:
            json.dump(json_to_send, w, indent=4)
        return json_to_send
