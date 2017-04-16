import json

FILE_CONFIG = "CONFIG.json"


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
        return True

    except Exception as e:
        print(e)
        return False
