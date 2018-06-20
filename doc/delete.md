---
description: Delete resource
draft: false
linktitle: delete
menu:
  docs:
    parent: commands
    weight: 5
title: Delete
weight: 2

---

#### <a name="delete">delete</a>

**Description**:

Delete resource

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -n | --namespace |  |  |


**Subcommands**:

* **[delete configmap](#delete_configmap)** delete configmap
* **[delete deployment](#delete_deployment)** delete deployment in specific namespace
* **[delete ingress](#delete_ingress)** delete ingress
* **[delete namespace](#delete_namespace)** delete namespace
* **[delete pod](#delete_pod)** delete pod in specific namespace
* **[delete service](#delete_service)** delete service in specific namespace
* **[delete volume](#delete_volume)** delete volume


#### <a name="delete_volume">delete volume</a>

**Description**:

delete volume

**Example**:

chkit delete volume [--force]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | suppress confirmation | false |


**Subcommands**:



#### <a name="delete_service">delete service</a>

**Description**:

Delete service in namespace.

**Example**:

chkit delete service service_label [-n namespace]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | force delete without confirmation | false |


**Subcommands**:



#### <a name="delete_pod">delete pod</a>

**Description**:

Delete pods.

**Example**:

chkit delete pod pod_name [-n namespace]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | delete pod without confirmation | false |


**Subcommands**:



#### <a name="delete_namespace">delete namespace</a>

**Description**:

Delete namespace provided in the first arg.

**Example**:

chkit delete namespace $ID

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | force delete without confirmation | false |


**Subcommands**:



#### <a name="delete_ingress">delete ingress</a>

**Description**:

Delete ingress.

**Example**:

chkit delete ingress $INGRESS [-n $NAMESPACE] [--force]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | delete ingress without confirmation | false |


**Subcommands**:



#### <a name="delete_deployment">delete deployment</a>

**Description**:

Delete deployment in specific namespace. Use --force flag to suppress confirmation.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | delete without confirmation | false |


**Subcommands**:



#### <a name="delete_configmap">delete configmap</a>

**Description**:

delete configmap

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | suppress confirmation | false |


**Subcommands**:



