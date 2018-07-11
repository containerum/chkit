---
description: 'view pod logs. Aliases: log'
draft: false
linktitle: logs
menu:
  docs:
    parent: commands
    weight: 5
title: Logs
weight: 2

---

#### <a name="logs">logs</a>

**Description**:

view pod logs. Aliases: log

**Example**:

logs pod_label [container] [--follow] [--prev] [--tail n] [--quiet]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
| -f | --follow | follow pod logs | false |
| -q | --quiet | print only logs and errors | false |
| -t | --tail | print last <value> log lines | 100 |


**Subcommands**:



