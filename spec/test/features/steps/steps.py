from contextlib import suppress
from time import sleep
import json
import shlex
import socket
import subprocess

import requests
from environconfig import EnvironConfig, StringVar, IntVar, BooleanVar
from comparedict import is_subset
import jsonexample

import logging

class Env(EnvironConfig):
    #: How to run Kapow! server
    KAPOW_SERVER_CMD = StringVar(default="kapow server")

    #: Where the Control API is
    KAPOW_CONTROLAPI_URL = StringVar(default="http://localhost:8081")

    #: Where the Data API is
    KAPOW_DATAAPI_URL = StringVar(default="http://localhost:8080")

    KAPOW_BOOT_TIMEOUT = IntVar(default=10)

    KAPOW_DEBUG_TESTS = BooleanVar(default=False)


if Env.KAPOW_DEBUG_TESTS:
    # These two lines enable debugging at httplib level
    # (requests->urllib3->http.client) You will see the REQUEST,
    # including HEADERS and DATA, and RESPONSE with HEADERS but without
    # DATA.  The only thing missing will be the response.body which is
    # not logged.
    try:
        import http.client as http_client
    except ImportError:
        # Python 2
        import httplib as http_client
    http_client.HTTPConnection.debuglevel = 1

    # You must initialize logging, otherwise you'll not see debug output.
    logging.basicConfig()
    logging.getLogger().setLevel(logging.DEBUG)
    requests_log = logging.getLogger("requests.packages.urllib3")
    requests_log.setLevel(logging.DEBUG)
    requests_log.propagate = True

def run_kapow_server(context):
    context.server = subprocess.Popen(
        shlex.split(Env.KAPOW_SERVER_CMD),
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        shell=False)

    # Check process is running with reachable APIs
    open_ports = False
    for _ in range(Env.KAPOW_BOOT_TIMEOUT):
        is_running = context.server.poll() is None
        assert is_running, "Server is not running!"
        with suppress(requests.exceptions.ConnectionError):
            open_ports = (
                requests.head(Env.KAPOW_CONTROLAPI_URL, timeout=1).status_code
                and requests.head(Env.KAPOW_DATAAPI_URL, timeout=1).status_code)
            if open_ports:
                break
        sleep(1)

    assert open_ports, "API is unreachable after KAPOW_BOOT_TIMEOUT"

@given('I have a just started Kapow! server')
@given('I have a running Kapow! server')
def step_impl(context):
    run_kapow_server(context)


@when('I request a routes listing')
def step_impl(context):
    context.response = requests.get(f"{Env.KAPOW_CONTROLAPI_URL}/routes")


@then('I get an empty list')
def step_impl(context):
    context.response.raise_for_status()
    assert context.response.json() == []


@given('I have a Kapow! server whith the following routes')
def step_impl(context):
    run_kapow_server(context)

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
    context.response = requests.post(f"{Env.KAPOW_CONTROLAPI_URL}/routes",
                                     data=context.text,
                                     headers={"Content-Type": "application/json"})


@then('I get {code} as response code')
def step_impl(context, code):
    assert context.response.status_code == int(code), f"Got {context.response.status_code} instead"


@then('I get "{reason}" as response reason phrase')
def step_impl(context, reason):
    assert context.response.reason == reason, f"Got {context.response.reason} instead"


@then('I get the following response body')
def step_impl(context):
    assert is_subset(jsonexample.loads(context.text), context.response.json())


@then('I get an empty response body')
def step_impl(context):
    assert context.response.content == b'', f"Response body is not empty. Got {context.response.content} instead."


@when('I delete the route with id "{id}"')
def step_impl(context, id):
    raise NotImplementedError('STEP: When I delete the route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"')


@given('It has a route with id "{id}"')
def step_impl(context, id):
    raise NotImplementedError('STEP: Given It has a route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"')


@when('I insert the route')
def step_impl(context):
    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    row = context.table[0]
    context.response = requests.put(f"{Env.KAPOW_CONTROLAPI_URL}/routes",
                                    json={h: row[h] for h in row.headings})


@when('I try to append with this malformed JSON document')
@when('I try to append with this JSON document')
def step_impl(context):
    context.response = requests.post(
        f"{Env.KAPOW_CONTROLAPI_URL}/routes",
        headers={"Content-Type": "application/json"},
        data=context.text)


@when('I delete the first route inserted')
def step_impl(context):
    raise NotImplementedError('STEP: When I delete the first route inserted')


@when('I delete the last route inserted')
def step_impl(context):
    raise NotImplementedError('STEP: When I delete the last route inserted')


@when('I delete the second route inserted')
def step_impl(context):
    raise NotImplementedError('STEP: When I delete the second route inserted')


@when('I try to insert with this JSON document')
def step_impl(context):
    raise NotImplementedError('STEP: When I try to insert with this JSON document')
