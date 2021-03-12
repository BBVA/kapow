#!/usr/bin/env python

import json
import os
import sys

if __name__ == '__main__':
    with open(os.environ['SPECTEST_FIFO'], 'w') as fifo:
        json.dump(dict(os.environ), fifo)
