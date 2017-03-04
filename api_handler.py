import json
import requests


class ApiHandler:
    def __init__(self):
        self.server = 'http://146.185.135.181:8080'
        self.headers = {
            # 'Content-Type': 'application/json-patch+json',
            # 'TCP-Connection-Name': 'none'
            'Authorization': '31'
        }
        self.TIMEOUT = 10

    def create(self, json_to_send):
        kind = '{}s'.format(json_to_send['kind'].lower())

        if json_to_send.get('metadata').get('namespace'):
            url = '{}/namespaces/{}/{}'.format(
                self.server,
                json_to_send.get('metadata').get('namespace'),
                kind
            )
        else:
            url = '{}/namespaces/default/{}'.format(
                self.server,
                kind
            )

        result = request_post(url, self.headers, self.TIMEOUT, json_to_send)

        return result

    def create_namespaces(self, json_to_send):
        url = '{}/namespaces'.format(
            self.server
        )

        # print(url)

        result = request_post(url, self.headers, self.TIMEOUT, json_to_send)

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

        # print(url)

        result = request_put(url, self.headers, self.TIMEOUT, json_to_send)

        return result

    def replace_namespaces(self, json_to_send):
        name = json_to_send['metadata']['name']

        url = '{}/namespaces/{}'.format(
            self.server,
            name
        )

        # print(url)

        result = request_put(url, self.headers, self.TIMEOUT, json_to_send)

        return result

    def run(self, json_to_send):
        if json_to_send.get('metadata').get('namespace'):
            url = '{}/namespaces/{}/deployments'.format(
                self.server,
                json_to_send.get('metadata').get('namespace')
            )
        else:
            url = '{}/namespaces/default/deployments'.format(
                self.server
            )

        # print(url)

        result = request_post(url, self.headers, self.TIMEOUT, json_to_send)

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
        # print('DELETE {}'.format(url))

        result = request_delete(url, self.headers, self.TIMEOUT)

        return result

    def delete_namespaces(self, name):
        url = '{}/namespaces/{}'.format(
            self.server,
            name
        )

        result = request_delete(url, self.headers, self.TIMEOUT)

        return result

    def get(self, kind, name, namespace):
        if not namespace:
            namespace = 'default'

        if name:
            url = '{}/namespaces/{}/{}/{}'.format(
                self.server,
                namespace,
                kind,
                name
            )
        else:
            url = '{}/namespaces/{}/{}'.format(
                self.server,
                namespace,
                kind
            )

        # print(url)

        result = request_get(url, self.headers, self.TIMEOUT)

        return result

    def get_namespaces(self, name):
        if name:
            url = '{}/namespaces/{}'.format(
                self.server,
                name
            )
        else:
            url = '{}/namespaces'.format(
                self.server
            )

        # print(url)

        result = request_get(url, self.headers, self.TIMEOUT)

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
        except:
            return {'error': 'connection error'}

    return func_wrapper


@request_exceptions_decorate
def request_get(url, headers, timeout):
    r = requests.get(
        url,
        headers=headers,
        timeout=timeout
    )
    if r.status_code == 200:
        return json.loads(r.text)
    else:
        raise StatusException(r.status_code, r._content)


@request_exceptions_decorate
def request_delete(url, headers, timeout):
    r = requests.delete(
        url,
        headers=headers,
        timeout=timeout
    )
    if r.status_code == 200:
        return json.loads(r.text)
    else:
        raise StatusException(r.status_code, r._content)


@request_exceptions_decorate
def request_post(url, headers, timeout, json_to_send):
    r = requests.post(
        url,
        data=json.dumps(json_to_send),
        timeout=timeout,
        headers=headers
    )
    if r.status_code == 200:
        return json.loads(r.text)
    else:
        raise StatusException(r.status_code, r._content)


@request_exceptions_decorate
def request_put(url, headers, timeout, json_to_send):
    r = requests.put(
        url,
        data=json.dumps(json_to_send),
        timeout=timeout,
        headers=headers
    )
    # print(url)
    # print(vars(r))
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
