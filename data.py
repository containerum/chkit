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
    'namespaces',
    'namespace'
]

output_formats = [
    'yaml',
    'json',
    'pretty',
]

deployment_json = {
    'kind': 'Deployment',
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "name": '',
        "labels": {
        }
    },
    "spec": {
        "replicas": '',
        "template": {
            "metadata": {
                "name": '',
                "labels": {
                    "run": ''
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": '',
                        "image": '',
                        "ports": [],
                        "env": [],
                        "commands": [],
                        "resources":{
                            "limits": {
                                "cpu": "",
                                "memory": ''
                            },
                            "requests":{
                                "cpu": "",
                                "memory": ""
                            }
                        }
                    }
                ]
            }
        }
    }
}
