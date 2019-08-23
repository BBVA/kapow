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


WORD2POS = {"first": 0, "second": 1, "last": -1}


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


@given('I have a Kapow! server with the following routes')
def step_impl(context):
    run_kapow_server(context)

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    for row in context.table:
        response = requests.post(f"{Env.KAPOW_CONTROLAPI_URL}/routes",
                                 json={h: row[h] for h in row.headings})
        response.raise_for_status()


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


@when('I delete the route with id "{id}"')
def step_impl(context, id):
    context.response = requests.delete(f"{Env.KAPOW_CONTROLAPI_URL}/routes/{id}")


@when('I insert the route')
def step_impl(context):
    context.response = requests.put(f"{Env.KAPOW_CONTROLAPI_URL}/routes",
                                    headers={"Content-Type": "application/json"},
                                    data=context.text)


@when('I try to append with this malformed JSON document')
def step_impl(context):
    context.response = requests.post(
        f"{Env.KAPOW_CONTROLAPI_URL}/routes",
        headers={"Content-Type": "application/json"},
        data=context.text)


@when('I delete the {order} route')
def step_impl(context, order):
    idx = WORD2POS.get(order)
    routes = requests.get(f"{Env.KAPOW_CONTROLAPI_URL}/routes")
    id = routes.json()[idx]["id"]
    context.response = requests.delete(f"{Env.KAPOW_CONTROLAPI_URL}/routes/{id}")


@when('I try to insert with this JSON document')
def step_impl(context):
    context.response = requests.put(
        f"{Env.KAPOW_CONTROLAPI_URL}/routes",
        headers={"Content-Type": "application/json"},
        data=context.text)

@when('I get the route with id "{id}"')
def step_impl(context, id):
    context.response = requests.get(f"{Env.KAPOW_CONTROLAPI_URL}/routes/{id}")


@when('I get the {order} route')
def step_impl(context, order):
    idx = WORD2POS.get(order)
    routes = requests.get(f"{Env.KAPOW_CONTROLAPI_URL}/routes")
    id = routes.json()[idx]["id"]
    context.response = requests.get(f"{Env.KAPOW_CONTROLAPI_URL}/routes/{id}")
