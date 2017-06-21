kinds = [
    'po',
    'pods',
    'pod',
    'deployments',
    'deployment',
    'deploy',
    'service',
    'services',
    'svc',
    'ns',
    'namespaces',
    'namespace'
]

expose_kinds = [
    'deployments',
    'deployment',
    'deploy'
]

run_kinds = [
    'deployments',
    'deployment',
    'deploy'
]

fields = [
    "image",
    "replicas"
]

delete_kinds = [
    'po',
    'pods',
    'pod',
    'deployments',
    'deployment',
    'deploy',
    'service',
    'services',
    'svc'
]

output_formats = [
    'yaml',
    'json',
    'pretty',
]

deployment_json = {
    "metadata": {
        "name": "deplo5",
        "labels": {}
    },
    "kind": "Deployment",
    "spec": {
        "template": {
            "metadata": {
                "labels": {
                },
                "name": "deplo4"
            },
            "spec": {
                "containers": [
                    {
                        "name": "deplo4",
                        "resources": {
                            "requests": {
                                "memory": "128Mi",
                                "cpu": "100m"
                            }
                        },
                        "image": "ubuntu"
                    }
                ]
            }
        },
        "replicas": 1
    }
}

service_json = {
    "kind": "Service",
    "metadata": {
        "name": "",
        "labels": {}
    },
    "spec": {
        "ports": [
        ],
        "selector": {}
    }
}

config_json = {
    "api_handler": {
        "headers": {
            "Authorization": ""
        },
        "TIMEOUT": 10,
        "server": "http://sdk.containerum.io:3333"
    },
    "tcp_handler": {
        "TCP_IP": "sdk.containerum.io",
        "AUTH_FORM": {
            "token": ""
        },
        "BUFFER_SIZE": 1024,
        "TCP_PORT": 3000
    },
    "default_namespace": "default"
}
