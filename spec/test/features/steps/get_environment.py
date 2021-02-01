#!/usr/bin/env python

import json
import os
import sys

if __name__ == '__main__':
    with open(sys.argv[1], 'w') as fifo:
        json.dump(os.environ, fifo)
