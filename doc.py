#!/usr/bin/env python
"""Generates doc files for chkit CLI util"""

import subprocess

__author__ = "ninedraft"
__copyright__ = "Copyright 2018, Exon Lab"
__license__ = "MIT"
__version__ = "1.0.0"
__maintainer__ = __author__
__status__ = "Production"

result = subprocess.run(['./chkit', 'doc',
                         '--list',
                         '--format', '{{.Path}};{{printf \"%q\" .Description}}'], stdout=subprocess.PIPE)
if result.returncode != 0:
    print(result)
    exit(result.returncode)
commands = list(cmd.decode("utf-8")
                for cmd in result.stdout.split(b'\n')
                if len(cmd.strip()) > 0)

for command in commands:
    name, description = (command.split(";") + ['', ''])[:2]
    yaml = f"""---
title: {name.title()}
linktitle: {name}
description: {(description.split('.') + [''])[0].replace('"', '')}

categories: []
keywords: []

menu:
  docs:
    parent: "commands"
    weight: 5

weight: 2

draft: false
---

"""
    result = subprocess.run(['./chkit', 'doc',
                             '--md',
                             '--command', name],
                            stdout=subprocess.PIPE)
    if result.returncode != 0:
        print(result)
        exit(result.returncode)
    md = yaml + bytes(result.stdout).decode("utf-8")
    fname = name.replace(' ', '_')
    with open(f'./doc/{fname}.md', mode='w') as docFile:
        docFile.write(md)
