---
title: Replace Ingress
linktitle: replace ingress
description: Replace ingress with a new one, use --force flag to write one-liner command, omitted attributes are inherited from the previous ingress

categories: []
keywords: []

menu:
  docs:
    parent: "commands"
    weight: 5

weight: 2

draft: false
---

### replace ingress

**Description**:

Replace ingress with a new one, use --force flag to write one-liner command, omitted attributes are inherited from the previous ingress.

**Example**:

chkit replace ingress $INGRESS [--force] [--service $SERVICE] [--port 80] [--tls-secret letsencrypt]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | force | replace ingress without confirmation | false |
|  | host | ingress host, optional |  |
|  | port | ingress endpoint port, optional | 8080 |
|  | service | ingress endpoint service, optional |  |
|  | tls-secret | ingress tls-secret, use 'letsencrypt' for automatic HTTPS, '-' to use HTTP, optional |  |



