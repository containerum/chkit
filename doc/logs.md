---
title: Logs
linktitle: logs
description: view pod logs

categories: []
keywords: []

menu:
  docs:
    parent: "commands"
    weight: 5

weight: 2

draft: false
---

### logs

**Description**:

view pod logs. Aliases: log

**Example**:

logs pod_label [container] [--follow] [--prev] [--tail n] [--quiet]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | follow | follow pod logs | false |
| -t | tail | print last <value> log lines | 100 |



