---
description: Set configuration variables
draft: false
linktitle: set
menu:
  docs:
    parent: commands
    weight: 5
title: Set
weight: 2

---

#### <a name="set">set</a>

**Description**:

Set configuration variables

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -n | --namespace |  |  |


**Subcommands**:

* **[set access](#set_access)** Set namespace access rights
* **[set containerum-api](#set_containerum-api)** Set Containerum API URL
* **[set default-namespace](#set_default-namespace)** Set default namespace
* **[set image](#set_image)** Set container image for specific deployment.
* **[set replicas](#set_replicas)** Set deployment replicas


#### <a name="set_replicas">set replicas</a>

**Description**:

Set deployment replicas.

**Example**:

chkit set replicas [-n namespace_label] [-d depl_label] [N_replicas]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -d | --deployment | deployment name |  |
| -r | --replicas | replicas, 1..15 | 1 |


**Subcommands**:



#### <a name="set_image">set image</a>

**Description**:

Set container image for specific deployment
If a deployment contains only one container, the command will use that container by default.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --container | container name |  |
|  | --deployment | deployment name |  |
| -f | --force | suppress confirmation | false |
|  | --image | new image |  |


**Subcommands**:



#### <a name="set_default-namespace">set default-namespace</a>

**Description**:

Set default namespace

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |


**Subcommands**:



#### <a name="set_containerum-api">set containerum-api</a>

**Description**:

Set Containerum API URL

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --allow-self-signed-certs |  | false |


**Subcommands**:



#### <a name="set_access">set access</a>

**Description**:

Set namespace access rights.
Available access levels are:
  none
  owner
  read
  read-delete
  write

**Example**:

chkit set access $USERNAME $ACCESS_LEVEL [--namespace $ID]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --force | suppress confirmation | false |


**Subcommands**:



