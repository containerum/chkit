---
title: Get Deployment-Versions
linktitle: get deployment-versions
description: Get deployment versions

categories: []
keywords: []

menu:
  docs:
    parent: "commands"
    weight: 5

weight: 2

draft: false
---

### get deployment-versions

**Description**:

Get deployment versions. You can filter versions by specifying version query (--version): Valid queries are:    - "<1.0.0"   - "<=1.0.0"   - ">1.0.0"   - ">=1.0.0"   - "1.0.0", "=1.0.0", "==1.0.0"   - "!1.0.0", "!=1.0.0" A query can consist of multiple querys separated by space: queries can be linked by logical AND:   - ">1.0.0 <2.0.0" would match between both querys, so "1.1.1" and "1.8.7" but not "1.0.0" or "2.0.0"   - ">1.0.0 <3.0.0 !2.0.3-beta.2" would match every version between 1.0.0 and 3.0.0 except 2.0.3-beta.2 Queries can also be linked by logical OR:   - "<2.0.0 || >=3.0.0" would match "1.x.x" and "3.x.x" but not "2.x.x" AND has a higher precedence than OR. It's not possible to use brackets. Queries can be combined by both AND and OR  - `>1.0.0 <2.0.0 || >3.0.0 !4.2.1` would match `1.2.3`, `1.9.9`, `3.1.1`, but not `4.2.1`, `2.1.1`

**Example**:

chkit get deployment-versions MY_DEPLOYMENT [--last-n 4] [--version >=1.0.0] [--output yaml] [--file versions.yaml]

**Flags**:

| Short | Name | Usage | Default value |
| ----- | ---- | ----- | ------------- |
|  | file | output file, optional, default is STDOUT |  |
|  | last-n | limit n versions to show | 0 |
|  | version | version query, examples: <1.0.0, <=1.0.0, !1.0.0 |  |



