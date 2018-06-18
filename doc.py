#!/usr/bin/env python
"""Generates doc files for chkit CLI util"""

import subprocess

__author__ = "ninedraft"
__copyright__ = "Copyright 2018, Exon Lab"
__license__ = "MIT"
__version__ = "1.0.0"
__maintainer__ = __author__
__status__ = "Production"

result = subprocess.run(['./chkit', 'doc', '--list'], stdout=subprocess.PIPE)
commands = list(cmd.decode("utf-8")
                for cmd in result.stdout.split(b'\n')
                if len(cmd.strip()) > 0)

for command in commands:
    result = subprocess.run(['./chkit', 'doc', '--md',
                             '--command', command,
                             '--output', './doc/' + command.replace(' ', '_') + '.md'],
                            stdout=subprocess.PIPE)
