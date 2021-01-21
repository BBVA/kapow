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
from multiprocessing.pool import ThreadPool
import time
import http.server

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
    KAPOW_CONTROL_URL = StringVar(default="http://localhost:8081")

    #: Where the Data API is
    KAPOW_DATA_URL = StringVar(default="http://localhost:8082")

    #: Where the User Interface is
    KAPOW_USER_URL = StringVar(default="http://localhost:8080")

    KAPOW_CONTROL_TOKEN = StringVar(default="TEST-SPEC-CONTROL-TOKEN")

    KAPOW_BOOT_TIMEOUT = IntVar(default=1000)

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

def run_kapow_server(context, control_token=True):
    with suppress(requests.exceptions.ConnectionError):
        open_ports = (
            requests.head(Env.KAPOW_CONTROL_URL, timeout=1).status_code
            and requests.head(Env.KAPOW_DATA_URL, timeout=1).status_code)
        assert (not open_ports), "Another process is already bound"

    control_token = {'KAPOW_CONTROL_TOKEN': Env.KAPOW_CONTROL_TOKEN} if control_token else {}

    context.server = subprocess.Popen(
        shlex.split(Env.KAPOW_SERVER_CMD),
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        env={**control_token, **os.environ},
        shell=False)

    # Check process is running with reachable APIs
    open_ports = False
    for _ in range(Env.KAPOW_BOOT_TIMEOUT):
        is_running = context.server.poll() is None
        assert is_running, "Server is not running!"
        with suppress(requests.exceptions.ConnectionError):
            open_ports = (
                requests.head(Env.KAPOW_CONTROL_URL, timeout=1).status_code
                and requests.head(Env.KAPOW_DATA_URL, timeout=1).status_code)
            if open_ports:
                break
        sleep(.01)

    assert open_ports, "API is unreachable after KAPOW_BOOT_TIMEOUT"


@given('I have a just started Kapow! server with {config}')
@given('I have a just started Kapow! server')
@given('I have a running Kapow! server')
def step_impl(context, config=None):
    control_token = config != 'no control token'
    run_kapow_server(context, control_token)


@when('I request a route listing without providing an Access Token')
def step_impl(context):
    context.response = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes")


@when('I request a route listing without providing an empty Access Token')
def step_impl(context):
    context.response = requests.get(
        f"{Env.KAPOW_CONTROL_URL}/routes",
        headers={"X-Kapow-Token": ""})


@when(u'I request a route listing providing a bad Access Token')
def step_impl(context):
    context.response = requests.get(
        f"{Env.KAPOW_CONTROL_URL}/routes",
        headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN + "42"})


@when('I request a routes listing')
def step_impl(context):
    context.response = requests.get(
        f"{Env.KAPOW_CONTROL_URL}/routes",
        headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})


@given('I have a Kapow! server with the following routes')
def step_impl(context):
    run_kapow_server(context)

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    for row in context.table:
        response = requests.post(
            f"{Env.KAPOW_CONTROL_URL}/routes",
            json={h: row[h] for h in row.headings},
            headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})
        response.raise_for_status()


@given('I have a Kapow! server with the following testing routes')
def step_impl(context):
    run_kapow_server(context)

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    for row in context.table:
        response = requests.post(
            f"{Env.KAPOW_CONTROL_URL}/routes",
            json={"entrypoint": " ".join(
                      [sys.executable,
                       shlex.quote(os.path.join(HERE, "testinghandler.py")),
                       shlex.quote(context.handler_fifo_path)]),  # Created in before_scenario
                  **{h: row[h] for h in row.headings}},
            headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})
        response.raise_for_status()

def testing_request(context, request_fn):
    # Run the request in background
    context.testing_request = ThreadPool(processes=1).apply_async(request_fn)

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
    context.testing_handler_pid = None
    context.testing_response = context.testing_request.get()


@when('I append the route')
def step_impl(context):
    context.response = requests.post(
        f"{Env.KAPOW_CONTROL_URL}/routes",
        data=context.text,
        headers={"Content-Type": "application/json",
                 "X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})

@then('I get {code} as response code')
def step_impl(context, code):
    assert context.response.status_code == int(code), f"Got {context.response.status_code} instead"


@then('I get {code} as response code in the testing request')
def step_impl(context, code):
    assert context.testing_response.status_code == int(code), f"Got {context.testing_response.status_code} instead"


@then('the response header "{header_name}" contains "{value}"')
def step_impl(context, header_name, value):
    assert context.response.headers.get(header_name, "").split(';')[0] == value, f"Got {context.response.headers.get(header_name)} instead"


@then('the testing response header {header_name} contains {value}')
def step_impl(context, header_name, value):
    assert context.testing_response.headers.get(header_name) == value, f"Got {context.testing_response.headers.get(header_name)} instead"


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
    context.response = requests.delete(
        f"{Env.KAPOW_CONTROL_URL}/routes/{id}",
        headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})


@when('I insert the route')
def step_impl(context):
    context.response = requests.put(
        f"{Env.KAPOW_CONTROL_URL}/routes",
        headers={"Content-Type": "application/json",
                 "X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN},
        data=context.text)


@when('I try to append with this malformed JSON document')
def step_impl(context):
    context.response = requests.post(
        f"{Env.KAPOW_CONTROL_URL}/routes",
        headers={"Content-Type": "application/json",
                 "X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN},
        data=context.text)


@when('I delete the {order} route')
def step_impl(context, order):
    idx = WORD2POS.get(order)
    routes = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes",
                          headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})
    id = routes.json()[idx]["id"]
    context.response = requests.delete(
        f"{Env.KAPOW_CONTROL_URL}/routes/{id}",
        headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})


@when('I try to insert with this JSON document')
def step_impl(context):
    context.response = requests.put(
        f"{Env.KAPOW_CONTROL_URL}/routes",
        headers={"Content-Type": "application/json",
                 "X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN},
        data=context.text)

@when('I get the route with id "{id}"')
def step_impl(context, id):
    context.response = requests.get(
        f"{Env.KAPOW_CONTROL_URL}/routes/{id}",
        headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})


@when('I get the {order} route')
def step_impl(context, order):
    idx = WORD2POS.get(order)
    routes = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes",
                          headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})
    id = routes.json()[idx]["id"]
    context.response = requests.get(
        f"{Env.KAPOW_CONTROL_URL}/routes/{id}",
        headers={"X-Kapow-Token": Env.KAPOW_CONTROL_TOKEN})


@when('I get the resource "{resource}"')
@when('I get the resource "{resource}" for the handler with id "{handler_id}"')
def step_impl(context, resource, handler_id=None):
    if handler_id is None:
        handler_id = context.testing_handler_id

    context.response = requests.get(
        f"{Env.KAPOW_DATA_URL}/handlers/{handler_id}{resource}")


@when('I set the resource "{resource}" with value "{value}"')
def step_impl(context, resource, value):
    context.response = requests.put(
        f"{Env.KAPOW_DATA_URL}/handlers/{context.testing_handler_id}{resource}",
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


@then('I get the value "{value}" for the response "{fieldType}" named "{elementName}" in the testing request')
def step_impl(context, value, fieldType, elementName):
    if fieldType == "header":
        actual = context.testing_response.headers.get(elementName)
    elif fieldType == "cookie":
        actual = context.testing_response.cookies.get(elementName)
    elif fieldType == "body":
        actual = context.testing_response.text
    else:
        raise ValueError("Unknown fieldtype {fieldType!r}")

    assert actual == value, f"Expecting {fieldType} {elementName!r} to be {value!r}, got {actual!r} insted"


@given('a test HTTP server on the {port} port')
def step_impl(context, port):
    context.request_ready = threading.Event()
    context.request_ready.clear()
    context.response_ready = threading.Event()
    context.response_ready.clear()

    class SaveResponseHandler(http.server.BaseHTTPRequestHandler):
        def do_verb(self):
            context.request_response = self
            context.request_ready.set()
            context.response_ready.wait()
        do_GET=do_verb
        do_POST=do_verb
        do_PUT=do_verb
        do_DELETE=do_verb
        do_HEAD=do_verb

    if port == "control":
        port = 8081
    elif port == "data":
        port = 8082
    else:
        raise ValueError(f"Unknown port {port}")

    context.httpserver = http.server.HTTPServer(('127.0.0.1', port),
                                                SaveResponseHandler)
    context.httpserver_thread = threading.Thread(
        target=context.httpserver.serve_forever,
        daemon=True)
    context.httpserver_thread.start()


@step('I run the following command')
def step_impl(context):
    _, command = context.text.split('$')
    command = command.lstrip()

    def exec_in_thread():
        context.command = subprocess.Popen(command, shell=True)
        context.command.wait()

    context.command_thread = threading.Thread(target=exec_in_thread, daemon=True)
    context.command_thread.start()


@then('the HTTP server received a "{method}" request to "{path}"')
def step_impl(context, method, path):
    context.request_ready.wait()
    assert context.request_response.command == method, f"Method {context.request_response.command} is not {method}"
    assert context.request_response.path == path, f"Method {context.request_response.path} is not {path}"



@then('the received request has the header "{name}" set to "{value}"')
def step_impl(context, name, value):
    context.request_ready.wait()
    matching = context.request_response.headers[name]
    assert matching, f"Header {name} not found"
    assert matching == value, f"Value of header doesn't match. {matching} != {value}"


@when('the server responds with')
def step_impl(context):
    # TODO: set the fields given in the table
    has_body = False
    for row in context.table:
        if row['field'] == 'status':
            context.request_response.send_response(int(row['value']))
        elif row['field'].startswith('headers.'):
            _, header = row['field'].split('.')
            context.request_response.send_header(header, row['value'])
        elif row['field'] == 'body':
            has_body = True
            payload = row['value'].encode('utf-8')
            context.request_response.send_header('Content-Length', str(len(payload)))
            context.request_response.end_headers()
            context.request_response.wfile.write(payload)

    if not has_body:
        context.request_response.send_header('Content-Length', '0')
        context.request_response.end_headers()

    context.response_ready.set()

@then('the command exits {immediately} with "{returncode}"')
@then('the command exits with "{returncode}"')
def step_impl(context, returncode, immediately=False):
    context.command_thread.join(timeout=3.0 if immediately else None)
    if context.command_thread.is_alive():
        try:
            print("killing in the name of")
            context.command.kill()
        finally:
            assert False, "The command is still alive"

    else:
        context.command.wait()
        assert context.command.returncode == int(returncode), f"Command returned {context.command.returncode} instead of {returncode}"


@then('the received request doesn\'t have the header "{name}" set')
def step_impl(context, name):
    context.request_ready.wait()
    assert name not in context.request_response.headers, f"Header {name} found"
