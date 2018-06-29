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
* **[create deployment-container](#create_deployment-container)** create deployment container.
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
|  | --name | service name, optional | quartz-mckay |
|  | --port | service port, optional | 0 |
|  | --port-name | service port name, optional | danzl-bistre |
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



#### <a name="create_deployment-container">create deployment-container</a>

**Description**:

Add container to deployment container set. Available methods to build deployment:
    - from flags
    - with interactive commandline wizard
    - from yaml ot json file

Use --force flag to create container without interactive wizard.
If the --container-name flag is not specified then wizard generates name RANDOM_COLOR-IMAGE.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --configmap | container configmap mount, CONFIG:MOUNT_PATH or CONFIG (then MOUNTPATH is /etc/CONFIG) |  |
|  | --container-name | container name, required on --force |  |
|  | --cpu | container CPU limit, mCPU | 0 |
|  | --deployment | deployment name, required on --force |  |
|  | --env | container environment variables, NAME:VALUE, 'NAME:$HOST_ENV' or '$HOST_ENV' (to user host env). WARNING: single quotes are required to prevent env from interpolation |  |
| -f | --force | suppress confirmation | false |
|  | --image | container image |  |
|  | --memory | container memory limit, Mb | 0 |
|  | --volume | container volume mounts, VOLUME:MOUNT_PATH or VOLUME (then MOUNT_PATH is /mnt/VOLUME) |  |


**Subcommands**:



#### <a name="create_deployment">create deployment</a>

**Description**:

Create deployment with containers and replicas.
Available methods to build deployment:
- from flags
- with interactive commandline wizard
- from yaml ot json file

Use --force flag to create container without interactive wizard.

There are several ways to specify the names of containers with flags:
- --container-name flag
- the prefix CONTAINER_NAME@ in the flags --image, --memory, --cpu, --env, --volume

If the --container-name flag is not specified and prefix is not used in any of the flags, then wizard searches for the --image flags without a prefix and generates name RANDOM_COLOR-IMAGE.

**Examples:**

---
**Single container with --container-name**

```bash
> ./ckit create depl \
        --container-name doot \
        --image nginx
```

|        LABEL        | VERSION |  STATUS  |  CONTAINERS  |    AGE    |
| ------------------- | --------| -------- | ------------ | --------- |
| akiraabe-heisenberg |  1.0.0  | inactive | doot [nginx] | undefined |

---
**Single container without --container-name**

```bash
> ./ckit create depl \
        --image nginx
```

|        LABEL        | VERSION |  STATUS  |        CONTAINERS        |    AGE    |
| ------------------- | --------| -------- | ------------------------ | --------- |
|   spiraea-kaufman   |  1.0.0  | inactive | aquamarine-nginx [nginx] | undefined |

---
**Multiple containers with --container-name**


```bash
> ./ckit create depl \
        --container-name gateway \
        --image nginx \
        --image blog@wordpress
```

|        LABEL        | VERSION |  STATUS  |        CONTAINERS        |    AGE    |
| ------------------- | --------| -------- | ------------------------ | --------- |
|   ruckers-fischer   |  1.0.0  | inactive |      gateway [nginx]     | undefined |
|                     |         |          |      blog [wordpress]    |           |

---
**Multiple containers without --container-name**
```bash
> ./ckit create depl \
        --image nginx \
        --image blog@wordpress
```

|        LABEL        | VERSION |  STATUS  |        CONTAINERS        |    AGE    |
| ------------------- | ------- | -------- | ------------------------ | --------- |
|    thisbe-neumann   |  1.0.0  | inactive |      blog [wordpress]    | undefined |
|                     |         |          |    garnet-nginx [nginx]  |           |


**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --configmap | container configmap, CONTAINER_NAME@CONFIGMAP_NAME@MOUNTPATH in case of multiple containers or CONFIGMAP_NAME@MOUNTPATH or CONFIGMAP_NAME in case of one container. If MOUNTPATH is omitted, then use /etc/CONFIGMAP_NAME as mountpath |  |
|  | --container-name | container name in case of single container |  |
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
|  | --name | configmap name | michela-zollner |


**Subcommands**:



