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
from contextlib import suppress, contextmanager
from multiprocessing.pool import ThreadPool
from time import sleep
import datetime
import http.server
import ipaddress
import json
import logging
import os
import shlex
import signal
import socket
import ssl
import subprocess
import sys
import tempfile
import threading
import time

from comparedict import is_subset
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives import serialization
from cryptography import x509
from cryptography.x509.oid import NameOID, ExtensionOID
from environconfig import EnvironConfig, StringVar, IntVar, BooleanVar
from requests import exceptions as requests_exceptions
import jsonexample
import requests


WORD2POS = {"first": 0, "second": 1, "last": -1}
HERE = os.path.dirname(__file__)


class Env(EnvironConfig):
    #: How to run Kapow! server
    KAPOW_SERVER_CMD = StringVar(default="kapow server")

    #: Where the Control API is
    KAPOW_CONTROL_URL = StringVar(default="https://localhost:8081")
    KAPOW_CONTROL_PORT = IntVar(default=8081)

    #: Where the Data API is
    KAPOW_DATA_URL = StringVar(default="http://localhost:8082")

    #: Where the User Interface is
    KAPOW_USER_URL = StringVar(default="http://localhost:8080")

    KAPOW_CONTROL_TOKEN = StringVar(default="TEST-SPEC-CONTROL-TOKEN")

    KAPOW_BOOT_TIMEOUT = IntVar(default=3000)

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


def generate_ssl_cert(subject_name, alternate_name):
    # Generate our key
    key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=4096,
    )
    # Various details about who we are. For a self-signed certificate the
    # subject and issuer are always the same.
    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COMMON_NAME, subject_name),
    ])
    cert = x509.CertificateBuilder().subject_name(
        subject
    ).issuer_name(
        issuer
    ).public_key(
        key.public_key()
    ).serial_number(
        x509.random_serial_number()
    ).not_valid_before(
        datetime.datetime.utcnow()
    ).not_valid_after(
        # Our certificate will be valid for 10 days
        datetime.datetime.utcnow() + datetime.timedelta(days=10)
    ).add_extension(
        x509.SubjectAlternativeName([x509.DNSName(alternate_name)]),
        critical=True,
    ).add_extension(
        x509.ExtendedKeyUsage(
            [x509.oid.ExtendedKeyUsageOID.SERVER_AUTH
             if subject_name.endswith('_server')
             else x509.oid.ExtendedKeyUsageOID.CLIENT_AUTH]),
        critical=True,
    # Sign our certificate with our private key
    ).sign(key, hashes.SHA256())

    key_bytes = key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=serialization.NoEncryption()
        )
    crt_bytes = cert.public_bytes(serialization.Encoding.PEM)

    return (key_bytes, crt_bytes)


@contextmanager
def mtls_client(context):
    with tempfile.NamedTemporaryFile(suffix='.crt', encoding='utf-8', mode='w') as srv_cert, \
         tempfile.NamedTemporaryFile(suffix='.crt', encoding='utf-8', mode='w') as cli_cert, \
         tempfile.NamedTemporaryFile(suffix='.key', encoding='utf-8', mode='w') as cli_key:
        srv_cert.write(context.init_script_environ["KAPOW_CONTROL_SERVER_CERT"])
        srv_cert.file.flush()
        cli_cert.write(context.init_script_environ["KAPOW_CONTROL_CLIENT_CERT"])
        cli_cert.file.flush()
        cli_key.write(context.init_script_environ["KAPOW_CONTROL_CLIENT_KEY"])
        cli_key.file.flush()
        session=requests.Session()
        session.verify=srv_cert.name
        session.cert=(cli_cert.name, cli_key.name)
        yield session


def is_port_open(port):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        return sock.connect_ex(('127.0.0.1', port)) == 0


def run_kapow_server(context, extra_args=""):
    assert (not is_port_open(Env.KAPOW_CONTROL_PORT)), "Another process is already bound"

    context.server = subprocess.Popen(
        shlex.split(Env.KAPOW_SERVER_CMD) + shlex.split(extra_args) + [os.path.join(HERE, "get_environment.py")],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        env={'SPECTEST_FIFO': context.init_script_fifo_path, **os.environ},
        shell=False)

    # Check process is running with reachable APIs
    open_ports = False
    for _ in range(Env.KAPOW_BOOT_TIMEOUT):
        with suppress(requests_exceptions.ConnectionError):
            if is_port_open(Env.KAPOW_CONTROL_PORT):
                open_ports = True
                break
        sleep(.01)

    assert open_ports, "API is unreachable after KAPOW_BOOT_TIMEOUT"

    # Get init_script enviroment via fifo
    with open(context.init_script_fifo_path, 'r') as fifo:
        context.init_script_environ = json.load(fifo)


@given('I have a just started Kapow! server')
@given('I have a running Kapow! server')
def step_impl(context):
    run_kapow_server(context)


@given(u'I launch the server with the following extra arguments')
def step_impl(context):
    run_kapow_server(context, context.text)


@when('I request a route listing without providing a Control Access Token')
def step_impl(context):
    with mtls_client(context) as requests:
        context.response = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes")


@when('I request a route listing without providing an empty Control Access Token')
def step_impl(context):
    with mtls_client(context) as requests:
        context.response = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes")


@when(u'I request a route listing providing a bad Control Access Token')
def step_impl(context):
    with mtls_client(context) as requests:
        context.response = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes")


@when('I request a routes listing')
def step_impl(context):
    with mtls_client(context) as requests:
        context.response = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes")


@given('I have a Kapow! server with the following routes')
def step_impl(context):
    run_kapow_server(context)

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    with mtls_client(context) as requests:
        for row in context.table:
            response = requests.post(
                f"{Env.KAPOW_CONTROL_URL}/routes",
                json={h: row[h] for h in row.headings})
            response.raise_for_status()


@given('I have a Kapow! server with the following testing routes')
def step_impl(context):
    run_kapow_server(context)

    if not hasattr(context, 'table'):
        raise RuntimeError("A table must be set for this step.")

    with mtls_client(context) as requests:
        for row in context.table:
            response = requests.post(
                f"{Env.KAPOW_CONTROL_URL}/routes",
                json={"entrypoint": " ".join(
                          [sys.executable,
                           shlex.quote(os.path.join(HERE, "testinghandler.py")),
                           shlex.quote(context.handler_fifo_path)]),  # Created in before_scenario
                      **{h: row[h] for h in row.headings}})
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
    with mtls_client(context) as requests:
        context.response = requests.post(
            f"{Env.KAPOW_CONTROL_URL}/routes",
            data=context.text,
            headers={"Content-Type": "application/json"})

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
    with mtls_client(context) as requests:
        context.response = requests.delete(
            f"{Env.KAPOW_CONTROL_URL}/routes/{id}")


@when('I insert the route')
def step_impl(context):
    with mtls_client(context) as requests:
        context.response = requests.put(
            f"{Env.KAPOW_CONTROL_URL}/routes",
            headers={"Content-Type": "application/json"},
            data=context.text)


@when('I try to append with this malformed JSON document')
def step_impl(context):
    with mtls_client(context) as requests:
        context.response = requests.post(
            f"{Env.KAPOW_CONTROL_URL}/routes",
            headers={"Content-Type": "application/json"},
            data=context.text)


@when('I delete the {order} route')
def step_impl(context, order):
    with mtls_client(context) as requests:
        idx = WORD2POS.get(order)
        routes = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes")
        id = routes.json()[idx]["id"]
        context.response = requests.delete(
            f"{Env.KAPOW_CONTROL_URL}/routes/{id}")


@when('I try to insert with this JSON document')
def step_impl(context):
    with mtls_client(context) as requests:
        context.response = requests.put(
            f"{Env.KAPOW_CONTROL_URL}/routes",
            headers={"Content-Type": "application/json"},
            data=context.text)

@when('I get the route with id "{id}"')
def step_impl(context, id):
    with mtls_client(context) as requests:
        context.response = requests.get(
            f"{Env.KAPOW_CONTROL_URL}/routes/{id}")


@when('I get the {order} route')
def step_impl(context, order):
    with mtls_client(context) as requests:
        idx = WORD2POS.get(order)
        routes = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes")
        id = routes.json()[idx]["id"]
        context.response = requests.get(
            f"{Env.KAPOW_CONTROL_URL}/routes/{id}")


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


@given('a test HTTPS server on the {port} port')
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

    context.srv_key, context.srv_crt = generate_ssl_cert("control_server", "localhost")
    context.cli_key, context.cli_crt = generate_ssl_cert("control_client", "localhost")
    with tempfile.NamedTemporaryFile(suffix=".key") as key_file, \
         tempfile.NamedTemporaryFile(suffix=".crt") as crt_file:
        key_file.write(context.srv_key)
        key_file.flush()
        crt_file.write(context.srv_crt)
        crt_file.flush()
        context.httpserver.socket = ssl.wrap_socket(
            context.httpserver.socket,
            keyfile=key_file.name,
            certfile=crt_file.name,
            server_side=True)
    context.httpserver_thread = threading.Thread(
        target=context.httpserver.serve_forever,
        daemon=True)
    context.httpserver_thread.start()


def run_command_with_certs(context, srv_crt, cli_crt, cli_key):
    _, command = context.text.split('$')
    command = command.lstrip()

    def exec_in_thread():
        context.command = subprocess.Popen(
            command,
            shell=True,
            env={'KAPOW_CONTROL_SERVER_CERT': srv_crt,
                 'KAPOW_CONTROL_CLIENT_CERT': cli_crt,
                 'KAPOW_CONTROL_CLIENT_KEY': cli_key,
                 **os.environ})
        context.command.wait()

    context.command_thread = threading.Thread(target=exec_in_thread, daemon=True)
    context.command_thread.start()

@step('I run the following command (with invalid certs)')
def step_impl(context):
    invalid_srv_crt, _ = generate_ssl_cert("invalid_control_server",
                                           "localhost")
    run_command_with_certs(context,
                           invalid_srv_crt,
                           context.cli_crt,
                           context.cli_key)


@step('I run the following command')
def step_impl(context):
    run_command_with_certs(context,
                           context.srv_crt,
                           context.cli_crt,
                           context.cli_key)


@when('I run the following command (setting the control certs environment variables)')
def step_impl(context):
    run_command_with_certs(
        context,
        context.init_script_environ["KAPOW_CONTROL_SERVER_CERT"],
        context.init_script_environ["KAPOW_CONTROL_CLIENT_CERT"],
        context.init_script_environ["KAPOW_CONTROL_CLIENT_KEY"])


@step('the HTTPS server receives a "{method}" request to "{path}"')
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


@when('I try to connect to the control API without providing a certificate')
def step_impl(context):
    try:
        context.request_response = requests.get(f"{Env.KAPOW_CONTROL_URL}/routes", verify=False)
    except Exception as exc:
        context.request_response = exc


@then(u'I get a connection error')
def step_impl(context):
    assert issubclass(type(context.request_response), Exception), context.request_response


@when(u'I try to connect to the control API providing an invalid certificate')
def step_impl(context):
    key, cert = generate_ssl_cert("foo", "localhost")
    with tempfile.NamedTemporaryFile(suffix='.crt') as cert_file, \
         tempfile.NamedTemporaryFile(suffix='.key') as key_file:
        cert_file.write(cert)
        cert_file.flush()
        key_file.write(key)
        key_file.flush()
        with requests.Session() as session:
            session.cert = (cert_file.name, key_file.name)
            session.verify = False
            try:
                context.request_response = session.get(
                    f"{Env.KAPOW_CONTROL_URL}/routes")
            except Exception as exc:
                context.request_response = exc



@when('I inspect the automatically generated control server certificate')
def step_impl(context):
    context.control_server_cert = x509.load_pem_x509_certificate(
        context.init_script_environ["KAPOW_CONTROL_SERVER_CERT"].encode('ascii'))


@then('the extension "{extension}" contains "{value}" of type "{typename}"')
def step_impl(context, extension, value, typename):
    if extension == 'Subject Alternative Name':
        oid = ExtensionOID.SUBJECT_ALTERNATIVE_NAME
    else:
        raise NotImplementedError(f'Unknown extension {extension}')

    if typename == 'DNSName':
        type_ = x509.DNSName
        converter = lambda x: x
    elif typename == 'IPAddress':
        type_ = x509.IPAddress
        converter = ipaddress.ip_address
    else:
        raise NotImplementedError(f'Unknown type {typename}')

    ext = context.control_server_cert.extensions.get_extension_for_oid(oid)
    values = ext.value.get_values_for_type(type_)

    assert converter(value) in values, f"Value {value} not in {values}"
