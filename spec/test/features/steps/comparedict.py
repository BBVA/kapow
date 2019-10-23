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
from functools import singledispatch
from itertools import zip_longest
from jsonexample import ANY


def assert_same_type(f):
    def wrapper(a, b):
        if type(a) != type(b):
            raise TypeError(f"Non-matching types {a!r} != {b!r}")
        return f(a, b)
    return wrapper


@singledispatch
@assert_same_type
def is_subset(model, obj):
    if model == obj:
        return True
    else:
        raise ValueError(f"Non-matching values {model!r} != {obj!r}")


@is_subset.register(dict)
@assert_same_type
def _(model, obj):
    for key, value in model.items():
        if key not in obj or not is_subset(value, obj[key]):
            raise ValueError(f"Non-matching dicts {model!r} != {obj!r}")
    return True


@is_subset.register(list)
@assert_same_type
def _(model, obj):
    for a, b in zip_longest(model, obj):
        if not is_subset(a, b):
            raise ValueError(f"Non-matching list member {a!r} in {b!r}")
    return True


@is_subset.register(ANY)
def _(model, obj):
    return True
