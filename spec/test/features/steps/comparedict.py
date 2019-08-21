from functools import singledispatch
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
    if is_subset(set(model), set(obj)):
        return True
    else:
        raise ValueError(f"Non-matching lists {model!r} != {obj!r}")


@is_subset.register(set)
@assert_same_type
def _(model, obj):
    return model <= obj


@is_subset.register(ANY)
def _(model, obj):
    return True
