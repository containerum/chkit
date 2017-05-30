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
    "image"
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
