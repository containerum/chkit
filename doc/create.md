---
description: Create resource (deployment, service...)
draft: false
linktitle: create
menu:
  docs:
    parent: commands
    weight: 5
title: Create
weight: 2

---

#### <a name="create">create</a>

**Description**:

Create resource (deployment, service...)

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -n | --namespace |  |  |


**Subcommands**:

* **[create configmap](#create_configmap)** 
* **[create deployment](#create_deployment)** create deployment
* **[create ingress](#create_ingress)** create ingress
* **[create service](#create_service)** create service


#### <a name="create_service">create service</a>

**Description**:

Create service for the specified pod in the specified namespace.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --deploy | service deployment, required |  |
| -f | --file | file with service data | - |
|  | --force | create service without confirmation | false |
|  | --name | service name, optional | flavescent-anaximander |
|  | --port | service port, optional | 0 |
|  | --port-name | service port name, optional | massalia-penny |
|  | --proto | service protocol, optional | TCP |
|  | --target-port | service target port, optional | 80 |


**Subcommands**:



#### <a name="create_ingress">create ingress</a>

**Description**:

Create ingress. Available options: TLS with LetsEncrypt and custom certs.

**Example**:

chkit create ingress [--force] [--filename ingress.json] [-n prettyNamespace]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | create ingress without confirmation | false |
|  | --host | ingress host (example: prettyblog.io), required |  |
|  | --path | path to endpoint (example: /content/pages), optional |  |
|  | --port | ingress endpoint port (example: 80, 443), optional | 8080 |
|  | --service | ingress endpoint service, required |  |
|  | --tls-cert | TLS cert file, optional |  |
|  | --tls-secret | TLS secret string, optional |  |


**Subcommands**:



#### <a name="create_deployment">create deployment</a>

**Description**:

Create a new deployment. Runs in one-line mode, suitable for integration with other tools, and in interactive wizard mode.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --configmap | container configmap, CONTAINER_NAME@CONFIGMAP_NAME@MOUNTPATH in case of multiple containers or CONFIGMAP_NAME@MOUNTPATH or CONFIGMAP_NAME in case of one container. If MOUNTPATH is omitted, then use /etc/CONFIGMAP_NAME as mountpath |  |
|  | --cpu | container memory limit, mCPU, CONTAINER_NAME@CPU in case of multiple containers or CPU in case of one container |  |
|  | --env | container environment variable, CONTAINER_NAME@KEY:VALUE in case of multiple containers or KEY:VALUE in case of one container |  |
|  | --file | file with configmap data, .json, .yaml, .yml, optional |  |
| -f | --force | suppress confirmation, optional | false |
|  | --image | container image, CONTAINER_NAME@IMAGE in case of multiple containers or IMAGE in case of one container |  |
|  | --memory | container memory limit, Mb, CONTAINER_NAME@MEMORY in case of multiple containers or MEMORY in case of one container |  |
|  | --name | deployment name, optional |  |
|  | --replicas | deployment replicas, optional | 0 |
|  | --volume | container volume, CONTAINER_NAME@VOLUME_NAME@MOUNTPATH in case of multiple containers or VOLUME_NAME@MOUNTPATH or VOLUME_NAME in case of one container. If MOUNTPATH is omitted, then use /mnt/VOLUME_NAME as mountpath |  |


**Subcommands**:



#### <a name="create_configmap">create configmap</a>

**Description**:



**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --file | file with configmap data |  |
| -f | --force | suppress confirmation | false |
|  | --item-file | configmap file, KEY:FILE_PATH or FILE_PATH |  |
|  | --item-string | configmap item, KEY:VALUE string pair |  |
|  | --name | configmap name | eunomia-knoll |


**Subcommands**:



