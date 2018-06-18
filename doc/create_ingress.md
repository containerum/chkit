---
title: Create Ingress
linktitle: create ingress
description: Create ingress

categories: []
keywords: []

menu:
  docs:
    parent: "commands"
    weight: 5

weight: 2

draft: false
---

### create ingress

**Description**:

Create ingress. Available options: TLS with LetsEncrypt and custom certs.

**Example**:

chkit create ingress [--force] [--filename ingress.json] [-n prettyNamespace]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | force | create ingress without confirmation | false |
|  | host | ingress host (example: prettyblog.io), required |  |
|  | path | path to endpoint (example: /content/pages), optional |  |
|  | port | ingress endpoint port (example: 80, 443), optional | 8080 |
|  | service | ingress endpoint service, required |  |
|  | tls-cert | TLS cert file, optional |  |
|  | tls-secret | TLS secret string, optional |  |



