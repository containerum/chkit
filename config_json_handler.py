import json
from bcolors import BColors
from keywords import SUCCESS_CHANGED
import os
import os.path
import re
from data import config_json
from datetime import datetime
from os_checker import get_file_config_path, create_folders

FILE_CONFIG = get_file_config_path()


def get_json_from_config():
    try:
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        return data
    except json.decoder.JSONDecodeError:
        data = config_json
        create_folders()
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(data,  indent=4))
        file.close()
        return data
    except FileNotFoundError:
        data = config_json
        create_folders()
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(data,  indent=4))
        file.close()
        return data


def show_namespace_token_from_config():
    try:
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        print('{}namespace: {} {}'.format(
                BColors.OKGREEN,
                data.get("default_namespace"),
                BColors.ENDC
            ))
        print('{}token: {} {}'.format(
                BColors.OKGREEN,
                data.get("tcp_handler").get("AUTH_FORM")["token"],
                BColors.ENDC
            ))
        return True

    except Exception as e:
        print('{}{}{} {}'.format(
                BColors.FAIL,
                "Error: ",
                e,
                BColors.ENDC,
            ))
        return False


def set_token_to_json_config(token):
    try:
        if not re.match("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$", token):
            raise ValueError("token is invalid")
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        data.get("tcp_handler").get("AUTH_FORM")["token"] = token
        data.get("api_handler").get("headers")["Authorization"] = token
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(data,  indent=4))
        file.close()
        print('{}{}{} '.format(
                BColors.OKBLUE,
                SUCCESS_CHANGED,
                BColors.ENDC,
            ))

        print('{}token: {} {}'.format(
                BColors.OKGREEN,
                token,
                BColors.ENDC
            ))
        return True

    except Exception as e:
        print('{}{}{} {}'.format(
                BColors.FAIL,
                "Error: ",
                e,
                BColors.ENDC,
            ))
        return False


def set_default_namespace_to_json_config(namespace):
    try:
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        data["default_namespace"] = namespace
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(data,  indent=4))
        file.close()
        print('{}{} {}'.format(
                BColors.OKBLUE,
                SUCCESS_CHANGED,
                BColors.ENDC
            ))
        print('{}namespace: {} {}'.format(
                BColors.OKGREEN,
                namespace,
                BColors.ENDC
            ))
        return True

    except Exception as e:
        print('{}{}{}{} '.format(
                BColors.FAIL,
                "Error: ",
                e,
                BColors.ENDC
            ))
        return False


def set_web_token_to_json_config(web_token):
    try:
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        data.get("webclient_api_handler")["headers"]["Authorization"] = web_token
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(data, file, indent=4))
        file.close()
        print('{}{} {}'.format(
                BColors.OKBLUE,
                SUCCESS_CHANGED,
                BColors.ENDC
            ))
        return True

    except Exception as e:
        print('{}{}{}{} '.format(
                BColors.FAIL,
                "Error: ",
                e,
                BColors.ENDC
            ))
        return False


def set_password_username_to_json_config(username, password):
    try:
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        data.get("webclient_api_handler")["username"] = username
        data.get("webclient_api_handler")["password"] = password
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(data, file, indent=4))
        file.close()
        print('{}{} {}'.format(
                BColors.OKBLUE,
                SUCCESS_CHANGED,
                BColors.ENDC
            ))
        return True

    except Exception as e:
        print('{}{}{}{} '.format(
                BColors.FAIL,
                "Error: ",
                e,
                BColors.ENDC
            ))
        return False


def update_chicking_time(time):
    if get_json_from_config():
        try:
            json_data = open(FILE_CONFIG).read()
            data = json.loads(json_data)
            data["checked_version_at"] = str(time)
            with open(FILE_CONFIG, "w") as file:
                file.write(json.dumps(data, file, indent=4))
            file.close()
            return data.get("checked_version_at")

        except Exception as e:
            print('{}{}{} {}'.format(
                    BColors.FAIL,
                    "Error: ",
                    e,
                    BColors.ENDC,
                ))


def check_last_update():
    if get_json_from_config():
        try:
            json_data = open(FILE_CONFIG).read()
            data = json.loads(json_data)
            checked_last_version = data.get("checked_version_at")
            if not checked_last_version:
                data["checked_version_at"] = str(datetime.now())
            with open(FILE_CONFIG, "w") as file:
                file.write(json.dumps(data, file, indent=4))
            file.close()
            return data.get("checked_version_at")

        except Exception as e:
            print('{}{}{} {}'.format(
                    BColors.FAIL,
                    "Error: ",
                    e,
                    BColors.ENDC,
                ))
