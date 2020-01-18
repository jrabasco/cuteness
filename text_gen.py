#!/usr/bin/python3.8

import os

"""
This script turns a file containing all of the texts to be displayed
into a texts folder compatible with the application
"""

with open('texts/main') as f:
    for i, line in enumerate(f):
        with open(f'texts/t{i}', 'w+') as output:
            output.write(line)

