import json
import requests
from config_json_handler import get_json_from_config

config_json_data = get_json_from_config()


class ApiHandler:
    def __init__(self, uuid_v4):
        self.server = config_json_data.get("api_handler").get("server")
        self.headers = config_json_data.get("api_handler").get("headers")
        self.headers.update({"Channel": uuid_v4})
        self.TIMEOUT = config_json_data.get("api_handler").get("TIMEOUT")

    def create(self, json_to_send, namespace=None):
        kind = '{}s'.format(json_to_send['kind'].lower())

        if namespace:
            url = '{}/namespaces/{}/{}'.format(
                self.server,
                namespace,
                kind
            )
        else:
            url = '{}/namespaces/default/{}'.format(
                self.server,
                kind
            )

        result = make_request(url, self.headers, self.TIMEOUT, "POST", json_to_send)
        return result

    def create_namespaces(self, json_to_send):
        url = '{}/namespaces'.format(
            self.server
        )

        result = make_request(url+"?cpu=500M&memory=500M&user=00000000-0000-0000-0000-000000000007", self.headers,
                              self.TIMEOUT, "POST", json_to_send)

        return result

    def replace(self, json_to_send, namespace):
        kind = '{}s'.format(json_to_send['kind'].lower())
        name = json_to_send['metadata']['name']

        if not namespace:
            if json_to_send.get('namespace'):
                namespace = json_to_send['namespace']
            else:
                namespace = 'default'

        url = '{}/namespaces/{}/{}/{}'.format(
            self.server,
            namespace,
            kind,
            name
        )

        result = make_request(url, self.headers, self.TIMEOUT, json_to_send, method="PUT")

        return result

    def replace_namespaces(self, json_to_send):
        name = json_to_send['metadata']['name']

        url = '{}/namespaces/{}'.format(
            self.server,
            name
        )

        result = make_request(url, self.headers, self.TIMEOUT, "PUT", json_to_send)

        return result

    def run(self, json_to_send, namespace=None):
        if namespace:
            url = '{}/namespaces/{}/deployments'.format(
                self.server,
                namespace
            )
        else:
            url = '{}/namespaces/default/deployments'.format(
                self.server
            )

        result = make_request(url, self.headers, self.TIMEOUT, "POST", json_to_send)

        return result

    def expose(self, json_to_send, namespace):
        if namespace:
            url = '{}/namespaces/{}/services'.format(
                self.server,
                namespace
            )
        else:
            url = '{}/namespaces/default/services'.format(
                self.server
            )

        result = make_request(url, self.headers, self.TIMEOUT, "POST", json_to_send)

        return result

    def delete(self, kind, name, namespace):
        if not namespace:
            namespace = 'default'

        url = '{}/namespaces/{}/{}/{}'.format(
            self.server,
            namespace,
            kind,
            name
        )

        result = make_request(url, self.headers, self.TIMEOUT, "DELETE")

        return result

    def delete_namespaces(self, name):
        url = '{}/namespaces/{}'.format(
            self.server,
            name
        )

        result = make_request(url, self.headers, self.TIMEOUT, "DELETE")

        return result

    def get(self, kind, name, namespace):

        if name:
            url = '{}/namespaces/{}/{}/{}'.format(
                self.server,
                namespace,
                kind,
                name[0]
            )
        else:
            url = '{}/namespaces/{}/{}'.format(
                self.server,
                namespace,
                kind
            )
        result = make_request(url, self.headers, self.TIMEOUT, "GET")
        return result

    def get_namespaces(self, name=None):
        if name:
            url = '{}/namespaces/{}'.format(
                self.server,
                name
            )
        else:
            url = '{}/namespaces'.format(
                self.server
            )

        result = make_request(url, self.headers, self.TIMEOUT, "GET")

        return result


def request_exceptions_decorate(func):
    def func_wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except requests.exceptions.Timeout:
            return {'error': 'timeout'}
        except json.decoder.JSONDecodeError as e:
            return {'error': str(e)}
        except StatusException as e:
            return {'error': str(e)}
        except Exception as e:
            print(e)
            return {'error': 'connection error'}

    return func_wrapper


@request_exceptions_decorate
def make_request(url, headers, timeout, method, json_to_send=None):
    if method == "DELETE":
        r = requests.delete(
            url,
            headers=headers,
            timeout=timeout
        )
    elif method == "POST":
        #print(url,headers, json_to_send, timeout)
        r = requests.post(
            url,
            data=json.dumps(json_to_send),
            timeout=timeout,
            headers=headers
        )
    elif method == "PUT":
        r = requests.put(
            url,
            data=json.dumps(json_to_send),
            timeout=timeout,
            headers=headers
        )
    else:
        r = requests.get(
            url,
            headers=headers,
            timeout=timeout
        )
    if r.status_code == 200:
        return json.loads(r.text)
    else:
        raise StatusException(r.status_code, r._content)


class StatusException(Exception):
    def __init__(self, status_code, content):
        self.status_code = status_code
        self.content = json.loads(content.decode('utf-8'))

    def __str__(self):
        return 'status error: {}\n{}'.format(self.status_code, self.content.get('error'))
