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
import json
from functools import partial
import re


class ANY:
    pass


class ExampleDecoder(json.JSONDecoder):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, object_hook=self.object_hook, **kwargs)

    def decode(self, s, *args, **kwargs):
        s = re.sub(r'(\W)(ANY)(\W)', r'\1{"_type": "ANY"}\3', s)
        return super().decode(s, *args, **kwargs)

    def object_hook(self, dct):
        if dct.get('_type', None) == 'ANY':
            return ANY()
        else:
            return dct

loads = partial(json.loads, cls=ExampleDecoder)
