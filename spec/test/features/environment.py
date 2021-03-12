#
# Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import tempfile
import os
import signal
from contextlib import suppress

def tmpfifo():
    while True:
        fifo_path = tempfile.mktemp() # The usage mkfifo make this safe
        try:
            os.mkfifo(fifo_path)
        except OSError:
            # The file already exist
            pass
        else:
            break

    return fifo_path


def before_scenario(context, scenario):
    context.handler_fifo_path = tmpfifo()
    context.init_script_fifo_path = tmpfifo()


def after_scenario(context, scenario):
    # Real Kapow! server being tested
    if hasattr(context, 'server'):
        context.server.terminate()
        context.server.wait()

    os.unlink(context.handler_fifo_path)
    os.unlink(context.init_script_fifo_path)

    # Mock HTTP server for testing
    if hasattr(context, 'httpserver'):
        context.response_ready.set()
        context.httpserver.shutdown()
        context.httpserver_thread.join()

    if getattr(context, 'testing_handler_pid', None) is not None:
        with suppress(ProcessLookupError):
            os.kill(int(context.testing_handler_pid), signal.SIGTERM)
