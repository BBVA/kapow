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
