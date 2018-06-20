---
description: Get resource data
draft: false
linktitle: get
menu:
  docs:
    parent: commands
    weight: 5
title: Get
weight: 2

---

#### <a name="get">get</a>

**Description**:

Get resource data

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -n | --namespace |  |  |


**Subcommands**:

* **[get access](#get_access)** print namespace access data
* **[get configmap](#get_configmap)** show configmap data
* **[get containerum-api](#get_containerum-api)** print Containerum API URL
* **[get default-namespace](#get_default-namespace)** print default namespace
* **[get deployment](#get_deployment)** show deployment data
* **[get deployment-versions](#get_deployment-versions)** get deployment versions
* **[get ingress](#get_ingress)** show ingress data
* **[get namespace](#get_namespace)** show namespace data or namespace list
* **[get pod](#get_pod)** show pod info
* **[get profile](#get_profile)** show profile info
* **[get service](#get_service)** show service info
* **[get solution](#get_solution)** get solutions


#### <a name="get_solution">get solution</a>

**Description**:

Show list of available solutions templates. To search solution by name add arg.

**Example**:

chkit get solution [name]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |


**Subcommands**:



#### <a name="get_service">get service</a>

**Description**:

Show service info.

**Example**:

chkit get service service_label [-o yaml/json] [-f output_file]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --file | output file | - |
| -o | --output | output format [yaml/json] |  |


**Subcommands**:



#### <a name="get_profile">get profile</a>

**Description**:

Shows profile info.

**Example**:

chkit get profile

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |


**Subcommands**:



#### <a name="get_pod">get pod</a>

**Description**:

Show pod info.

**Example**:

chkit get pod pod_label [-o yaml/json] [-f output_file]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --file | output file | - |
| -o | --output | output format (json/yaml) |  |


**Subcommands**:



#### <a name="get_namespace">get namespace</a>

**Description**:

show namespace data or namespace list.

**Example**:

chkit get $ID... [-o yaml/json] [-f output_file]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --file | output file |  |
| -o | --output | output format (json/yaml) |  |


**Subcommands**:



#### <a name="get_ingress">get ingress</a>

**Description**:

Print ingress data.

**Example**:

chkit get ingress ingress_names... [-n namespace_label] [-o yaml/json]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --file | output file |  |
| -o | --output | output format (yaml/json) |  |


**Subcommands**:



#### <a name="get_deployment-versions">get deployment-versions</a>

**Description**:

Get deployment versions. You can filter versions by specifying version query (--version): Valid queries are:    - "<1.0.0"   - "<=1.0.0"   - ">1.0.0"   - ">=1.0.0"   - "1.0.0", "=1.0.0", "==1.0.0"   - "!1.0.0", "!=1.0.0" A query can consist of multiple querys separated by space: queries can be linked by logical AND:   - ">1.0.0 <2.0.0" would match between both querys, so "1.1.1" and "1.8.7" but not "1.0.0" or "2.0.0"   - ">1.0.0 <3.0.0 !2.0.3-beta.2" would match every version between 1.0.0 and 3.0.0 except 2.0.3-beta.2 Queries can also be linked by logical OR:   - "<2.0.0 || >=3.0.0" would match "1.x.x" and "3.x.x" but not "2.x.x" AND has a higher precedence than OR. It's not possible to use brackets. Queries can be combined by both AND and OR  - `>1.0.0 <2.0.0 || >3.0.0 !4.2.1` would match `1.2.3`, `1.9.9`, `3.1.1`, but not `4.2.1`, `2.1.1`

**Example**:

chkit get deployment-versions MY_DEPLOYMENT [--last-n 4] [--version >=1.0.0] [--output yaml] [--file versions.yaml]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --file | output file, optional, default is STDOUT |  |
|  | --last-n | limit n versions to show | 0 |
|  | --version | version query, examples: <1.0.0, <=1.0.0, !1.0.0 |  |


**Subcommands**:



#### <a name="get_deployment">get deployment</a>

**Description**:

Print deployment data.

**Example**:

namespace deployment_names... [-n namespace_label]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --file | output file |  |
| -o | --output | output format (yaml/json) |  |


**Subcommands**:



#### <a name="get_default-namespace">get default-namespace</a>

**Description**:

Print default namespace.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |


**Subcommands**:



#### <a name="get_containerum-api">get containerum-api</a>

**Description**:

print Containerum API URL

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |


**Subcommands**:



#### <a name="get_configmap">get configmap</a>

**Description**:

show configmap data

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | --file | output file | - |
| -o | --output | output format yaml/json |  |


**Subcommands**:



#### <a name="get_access">get access</a>

**Description**:

Print namespace access data.

**Example**:

chkit get ns-access $ID

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |


**Subcommands**:



