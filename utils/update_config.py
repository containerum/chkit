#!/usr/bin/python3
import json

from bcolors import BColors

from variables.keywords import SUCCESS_CHANGED

FILE_CONFIG = "/var/lib/containerium/CONFIG.json"
FILE_CONFIG_FROM_SRC = "/var/lib/containerium/src/CONFIG.json"


def update_config_with_token_and_namespace():
    try:
        new_json_data = open(FILE_CONFIG_FROM_SRC).read()
        json_data = open(FILE_CONFIG).read()
        data = json.loads(json_data)
        new_data = json.loads(new_json_data)
        new_data.get("tcp_handler").get("AUTH_FORM")["token"] = data.get("tcp_handler").get("AUTH_FORM").get("token")
        new_data.get("api_handler").get("headers")["Authorization"] = data.get("api_handler")\
            .get("headers").get("Authorization")
        new_data["default_namespace"] = data.get("default_namespace")
        with open(FILE_CONFIG, "w") as file:
            file.write(json.dumps(new_data, file, indent=4))
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


if __name__ == "__main__":
    update_config_with_token_and_namespace()