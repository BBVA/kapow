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


def before_scenario(context, scenario):
    # Create the request_handler FIFO
    while True:
        context.handler_fifo_path = tempfile.mktemp() # Safe because using
                                                      # mkfifo
        try:
            os.mkfifo(context.handler_fifo_path)
        except OSError:
            # The file already exist
            pass
        else:
            break


def after_scenario(context, scenario):
    if hasattr(context, 'server'):
        context.server.terminate()
        context.server.wait()

    os.unlink(context.handler_fifo_path)
