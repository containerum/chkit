import requests
from config_json_handler import get_json_from_config, set_web_token_to_json_config, set_password_username_to_json_config
import json
from bcolors import BColors
import getpass
import base64
config_json_data = get_json_from_config()


def user_is_authenticated(func):
    def wrapper(self):
        url = self.server + "/api/token_status"
        r = requests.get(url, headers=config_json_data.get("webclient_api_handler").get("headers"))
        if r.status_code == 401:
            self.login()
            return func(self)
        else:
            return func(self)

    return wrapper


class WebClient:
    def __init__(self):
        self.server = config_json_data.get("webclient_api_handler").get("server")
        self.headers = config_json_data.get("webclient_api_handler").get("headers")

    def decode_password(self,password):
        return base64.b64decode(password.encode("utf8")).decode("utf8")

    def encode_password(self,password):
        return base64.b64encode(password.encode("utf8")).decode("utf8")

    def login(self):
        username = config_json_data.get("webclient_api_handler").get("username")
        password = self.decode_password(config_json_data.get("webclient_api_handler").get("password"))
        if not username or not password:
            username = input("Username:")
            password = getpass.getpass()



        url = self.server + "/api/login"
        r = requests.post(url, data={"username": username, "password": password})
        while True:
            if r.status_code == 200:
                set_password_username_to_json_config(username, self.encode_password(password))
                web_token = json.loads(r.text)["token"]
                set_web_token_to_json_config(web_token)
                self.headers["Authorization"] = web_token
                break
            else:
                print('{}{}{}{} '.format(
                    BColors.FAIL,
                    "Error: ",
                    json.loads(r.text).get("message"),
                    BColors.ENDC
                ))
                username = input("Username:")
                password = getpass.getpass()
                r = requests.post(url, data={"username": username, "password": password})

    @user_is_authenticated
    def get_namespaces(self):
        url = self.server + "/api/namespaces"
        r = requests.get(url, headers=self.headers)
        return json.loads(r.text)





