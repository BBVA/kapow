from contextlib import suppress
from time import sleep
import json
import os
import shlex
import signal
import socket
import subprocess
import sys
import tempfile
import threading
import time

import requests
from environconfig import EnvironConfig, StringVar, IntVar, BooleanVar
from comparedict import is_subset
import jsonexample

import logging


WORD2POS = {"first": 0, "second": 1, "last": -1}
HERE = os.path.dirname(__file__)


class Env(EnvironConfig):
    #: How to run Kapow! server
    KAPOW_SERVER_CMD = StringVar(default="kapow server")

    #: Where the Control API is
    KAPOW_CONTROLAPI_URL = StringVar(default="http://localhost:8081")

    #: Where the Data API is
    KAPOW_DATAAPI_URL = StringVar(default="http://localhost:8081")

    #: Where the User Interface is
    KAPOW_USER_URL = StringVar(default="http://localhost:8080")

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


@given('I have a Kapow! server with the following testing routes')
def step_impl(context):
    run_kapow_server(context)

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    for row in context.table:
        response = requests.post(
            f"{Env.KAPOW_CONTROLAPI_URL}/routes",
            json={"entrypoint": " ".join(
                      [sys.executable,
                       shlex.quote(os.path.join(HERE, "testinghandler.py")),
                       shlex.quote(context.handler_fifo_path)]),  # Created in before_scenario
                  **{h: row[h] for h in row.headings}})
        response.raise_for_status()

def testing_request(context, request_fn):
    # Run the request in background
    def _testing_request():
        context.testing_response = request_fn()
    context.testing_request = threading.Thread(target=_testing_request)
    context.testing_request.start()

    # Block until the handler connects and give us its pid and the
    # handler_id
    with open(context.handler_fifo_path, 'r') as fifo:
        (context.testing_handler_pid,
         context.testing_handler_id) = fifo.readline().rstrip('\n').split(';')


@when('I send a request to the testing route "{path}"')
def step_impl(context, path):
    def _request():
        try:
            return requests.get(f"{Env.KAPOW_USER_URL}{path}", stream=False)
        except:
            return None

    testing_request(context, _request)


@when('I release the testing request')
def step_impl(context):
    os.kill(int(context.testing_handler_pid), signal.SIGTERM)
    context.testing_request.join()


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


@then('I get the following response raw body')
def step_impl(context):
    assert context.text == context.response.text, f"{context.text!r} != {context.response.text!r}"


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


@when('I get the resource "{resource}"')
@when('I get the resource "{resource}" for the handler with id "{handler_id}"')
def step_impl(context, resource, handler_id=None):
    if handler_id is None:
        handler_id = context.testing_handler_id

    context.response = requests.get(
        f"{Env.KAPOW_DATAAPI_URL}/handlers/{handler_id}{resource}")


@when('I set the resource "{resource}" with value "{value}"')
def step_impl(context, resource, value):
    context.response = requests.put(
        f"{Env.KAPOW_DATAAPI_URL}/handlers/{context.testing_handler_id}{resource}",
        data=value.encode("utf-8"))


@when('I send a request to the testing route "{path}" adding')
def step_impl(context, path):
    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    params = {
        "headers": dict(),
        "cookies": dict(),
        "params": dict()}
    setters = {
        "header": params["headers"].setdefault,
        "cookie": params["cookies"].setdefault,
        "body": lambda _, v: params.setdefault("data", v.encode("utf-8")),
        "parameter": params["params"].setdefault,}

    for row in context.table:
        setters[row["fieldType"]](row["name"], row["value"])

    def _request():
        try:
            return requests.get(f"{Env.KAPOW_USER_URL}{path}",
                                stream=False,
                                **params)
        except:
            return None

    testing_request(context, _request)
