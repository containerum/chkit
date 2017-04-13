import json

DIR = "./CONFIG.json"


def get_json_from_config():
    json_data = open(DIR).read()
    data = json.loads(json_data)
    return data
