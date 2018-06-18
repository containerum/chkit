---
title: Replace Deployment
linktitle: replace deployment
description: Replaces deployment

categories: []
keywords: []

menu:
  docs:
    parent: "commands"
    weight: 5

weight: 2

draft: false
---

### replace deployment

**Description**:

Replaces deployment. Runs in one-line mode, suitable for integration with other tools, and in interactive wizard mode.

**Example**:



**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | container-name | container name, equal to image name by default |  |
|  | cpu | container CPU limit in mCPU, optional | 200 |
|  | env | container env variable in KEY0:VALUE0 KEY1:VALUE1 format |  |
|  | file | create deployment from file |  |
| -f | force | suppress confirmation | false |
|  | image | container image, optional |  |
|  | memory | container memory limit im Mb, optional | 256 |
|  | replicas | replicas, optional | 1 |



