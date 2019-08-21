from functools import singledispatch


def assert_same_type(f):
    def wrapper(a, b):
        if type(a) != type(b):
            raise TypeError("Non-matching types")
        return f(a, b)
    return wrapper


@singledispatch
@assert_same_type
def is_subset(model, obj):
    return model == obj


@is_subset.register(dict)
@assert_same_type
def _(model, obj):
    for key, value in model.items():
        if key not in obj or not is_subset(value, obj[key]):
            return False
    return True


@is_subset.register(list)
@assert_same_type
def _(model, obj):
    if type(model) != type(obj):
        raise TypeError("Non-matching types")
    return is_subset(set(model), set(obj))


@is_subset.register(set)
@assert_same_type
def _(model, obj):
    return model <= obj
