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
    "kind": "Deployment",
    "metadata": {
        "name": ""
    },
    "spec": {
        "replicas": 1,
        "template": {
            "metadata": {
                "name": "",
                "labels": {
                    "test": "app"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "",
                        "image": "",
                        "resources": {
                            "requests": {
                                "cpu": "100m",
                                "memory": "128Mi"
                            }
                        }
                    }
                ]
            }
        }
    }
}

service_json = {
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "name": "",
    },
    "spec": {
        "ports": [
        ],
        "selector": {}
    }
}
