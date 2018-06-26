---
description: Replace deployment or service
draft: false
linktitle: replace
menu:
  docs:
    parent: commands
    weight: 5
title: Replace
weight: 2

---

#### <a name="replace">replace</a>

**Description**:

Replace deployment or service

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -n | --namespace |  |  |


**Subcommands**:

* **[replace configmap](#replace_configmap)** 
* **[replace ingress](#replace_ingress)** Replace ingress with a new one.
* **[replace service](#replace_service)** Replace service.


#### <a name="replace_service">replace service</a>

**Description**:

Replace service.\nRuns in one-line mode, suitable for integration with other tools, and in interactive wizard mode.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --deployment | deployment name, optional |  |
|  | --domain | service domain, optional |  |
|  | --file | create service from file |  |
| -f | --force | suppress confirmation | false |
|  | --port | service external port, optional | 0 |
|  | --port-name | service port name |  |
|  | --protocol | service port protocol, optional | TCP |
|  | --target-port | service target port, optional | 80 |


**Subcommands**:



#### <a name="replace_ingress">replace ingress</a>

**Description**:

Replace ingress with a new one, use --force flag to write one-liner command, omitted attributes are inherited from the previous ingress.

**Example**:

chkit replace ingress $INGRESS [--force] [--service $SERVICE] [--port 80] [--tls-secret letsencrypt]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | replace ingress without confirmation | false |
|  | --host | ingress host, optional |  |
|  | --port | ingress endpoint port, optional | 8080 |
|  | --service | ingress endpoint service, optional |  |
|  | --tls-secret | ingress tls-secret, use 'letsencrypt' for automatic HTTPS, '-' to use HTTP, optional |  |


**Subcommands**:



#### <a name="replace_configmap">replace configmap</a>

**Description**:



**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --file | file with configmap data, .json, .yaml, .yml |  |
|  | --file-item | configmap file item: $KEY:$FILENAME |  |
|  | --force | suppress confirmation | false |
|  | --item | configmap item: $KEY:$VALUE |  |


**Subcommands**:



