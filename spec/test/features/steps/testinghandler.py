#!/usr/bin/env python
import signal
import sys
import os


if __name__ == '__main__':
    test_runner_fifo = sys.argv[1]
    with open(test_runner_fifo, 'w') as other_side:
        other_side.write(f"{os.getpid()};{os.environ['KAPOW_HANDLER_ID']}\n")

    signal.pause()
