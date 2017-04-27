import json
from bcolors import BColors
from keywords import SUCCESS_CHANGED
import os.path
FILE_CONFIG = "CONFIG.json"
FILE_CONFIG_FROM_SRC = "/var/lib/containerium/src/CONFIG.json"


def get_json_from_config():
    json_data = open(FILE_CONFIG).read()
    data = json.loads(json_data)
    return data


def set_token_to_json_config(token):
    try:
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        data.get("tcp_handler").get("AUTH_FORM")["token"] = token
        data.get("api_handler").get("headers")["Authorization"] = token
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(data, file, indent=4))
        file.close()
        print('{}{}{} '.format(
                BColors.OKBLUE,
                SUCCESS_CHANGED,
                BColors.ENDC,
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


def set_password_username_to_json_config(username,password):
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
