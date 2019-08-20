import subprocess
from time import sleep

import requests
from environconfig import EnvironConfig, StringVar


class Env(EnvironConfig):
    #: How to run Kapow! server
    KAPOW_SERVER_CMD = StringVar(default="kapow server")

    #: Where the Control API is
    KAPOW_CONTROLAPI_URL = StringVar(default="http://localhost:8081")

    #: Where the Data API is
    KAPOW_DATAAPI_URL = StringVar(default="http://localhost:8080")


@given('I have a just started Kapow! server')
@given('I have a running Kapow! server')
def step_impl(context):
    context.server = subprocess.Popen(
        Env.KAPOW_SERVER_CMD,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        shell=True)
    is_running = context.server.poll() is None
    assert is_running, "Server is not running!"


@when('I request a routes listing')
def step_impl(context):
    context.response = requests.get(f"{Env.KAPOW_CONTROLAPI_URL}/routes")


@then('I get an empty list')
def step_impl(context):
    context.response.raise_for_status()
    assert context.response.json() == []


@given('I have a Kapow! server whith the following routes')
def step_impl(context):
    context.server = subprocess.Popen(
        Env.KAPOW_SERVER_CMD,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        shell=True)
    is_running = context.server.poll() is None
    assert is_running, "Server is not running!"

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    for row in context.table:
        response = requests.post(f"{Env.KAPOW_CONTROLAPI_URL}/routes",
                                 json={h: row[h] for h in row.headings})
        response.raise_for_status()


@then('I get a list with the following elements')
def step_impl(context):
    context.response.raise_for_status()

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    for entry, row in zip(context.response.json(), context.table):
        for header in row.headings:
            assert header in entry, f"Response does not contain the key {header}"
            if row[header] != '*':
                assert entry[header] == row[header], f"Values mismatch"

#
#
#

@when('I append the route')
def step_impl(context):
    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    for row in context.table:
        response = requests.post(f"{Env.KAPOW_CONTROLAPI_URL}/routes",
                                 json={h: row[h] for h in row.headings})
        response.raise_for_status()


@then('I get {code} as response code')
def step_impl(context, code):
    raise NotImplementedError('STEP: Then I get unprocessable entity as response code')


@then('I get "{reason}" as response reason phrase')
def step_impl(context, reason):
    raise NotImplementedError('STEP: Then I get "Missing Mandatory Field" as response phrase')


@then('I get the following entity as response body')
def step_impl(context):
    raise NotImplementedError('STEP: Then I get the following entity as response body')


@then('I get an empty response body')
def step_impl(context):
    raise NotImplementedError('STEP: Then I get an empty response body')


@when('I delete the route with id "{id}"')
def step_impl(context, id):
    raise NotImplementedError('STEP: When I delete the route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"')


@given('It has a route with id "{id}"')
def step_impl(context, id):
    raise NotImplementedError('STEP: Given It has a route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"')

@when('I insert the route')
def step_impl(context):
    raise NotImplementedError('STEP: When I insert the route')
