#!/usr/bin/env python
"""
A Kapow! interpreter written in Python.

"""

from dataclasses import dataclass
from shlex import quote as shell_quote
from string import Template
import asyncio
import contextlib
import io
import logging
import os
import sys
import tempfile
import threading
import traceback

from aiohttp import web
from pyparsing import alphas, nums, White
from pyparsing import LineStart, LineEnd, SkipTo
from pyparsing import Literal, Combine, Word, Suppress
from pyparsing import OneOrMore, Optional, delimitedList
import aiofiles
import click

log = logging.getLogger('kapow')

########################################################################
#                                Parser                                #
########################################################################

# Method
METHOD = (Literal('GET')
          | Literal('POST')
          | Literal('PUT')
          | Literal('DELETE')
          | Literal('PATCH'))
MULTI_METHOD = delimitedList(METHOD, delim="|", combine=True)
METHOD_SPEC = Combine(Literal('*') | MULTI_METHOD)

# Pattern
REGEX = Word(alphas + nums + '\\+*,.[]-_')(name="regex")
SYMBOL = Word(alphas)(name="symbol")
P_PATTERN = Combine('/{' + SYMBOL + Optional(':' + REGEX) + '}')
P_PATH = Word('/', alphas + nums + '$-_.+!*\'(),')
URLPATTERN = Combine(OneOrMore(P_PATTERN | P_PATH))(name="urlpattern")

# Body
BODY = (Suppress('{')
        + SkipTo(Combine(LineStart() + '}' + LineEnd()))(name="body"))

# Endpoint head
ENDPOINT = (Optional(METHOD_SPEC + Suppress(White()),
                     default='*')(name="method")
            + URLPATTERN
            + Suppress(White()))

# Endpoint bodies
CODE_EP = (ENDPOINT + BODY)(name="code_ep")
PATH_EP = (ENDPOINT + '=' + SkipTo(LineEnd())(name="path"))(name="path_ep")

KAPOW_PROGRAM = OneOrMore(CODE_EP | PATH_EP)


########################################################################
#                              Resources                               #
########################################################################

@dataclass
class ResourceManager:
    """A resource exposed to the subshell."""
    #: Kapow resource representation
    kapow_repr: str
    #: Representation of the resource that can be understood by the shell
    shell_repr: str
    #: Coroutine capable of managing the resource internally
    coro: object
    #: Path to readed fifo. Needs to be written for coro to release.
    #: XXX: Use proper fifo async instead
    fifo_path: str = None
    #: Fifo direction 'read'/'write'
    fifo_direction: str = None



async def get_value(context, path):
    """Return the value of an http resource."""
    def nrd(n):
        """Return the nrd element in a path."""
        return path.split('/')[n]

    try:
        if path == 'request/method':
            return context['request'].method.encode('utf-8')
        elif path == 'request/path':
            return context['request'].path.encode('utf-8')
        elif path.startswith('request/match'):
            return context['request'].match_info[nrd(2)].encode('utf-8')
        elif path.startswith('request/param'):
            return context['request'].rel_url.query[nrd(2)].encode('utf-8')
        elif path.startswith('request/header'):
            return context['request'].headers[nrd(2)].encode('utf-8')
        elif path.startswith('request/cookie'):
            return context['request'].cookies[nrd(2)].encode('utf-8')
        elif path.startswith('request/form'):
            return (await context['request'].post())[nrd(2)].encode('utf-8')
        elif path.startswith('request/file'):
            name = nrd(2)
            content = nrd(3)  # filename / content
            field = (await context['request'].post())[name]
            if content == 'filename':
                try:
                    return field.filename.encode('utf-8')
                except:
                    return b''
            elif content == 'content':
                try:
                    return field.file.read()
                except:
                    return b''
            else:
                raise ValueError(f'Unknown content type {content!r}')
        elif path == 'request/body':
            return await context['request'].read()
        else:
            raise ValueError(f'Unknown path {path!r}')
    except KeyError:
        return b''


async def set_value(context, path, value):
    """
    Write to an http resource.

    File-like resources like `body` get write() calls so they have
    append semantics. Non file-like resources are just set.

    """
    if not value:
        return

    def nrd(n):
        return path.split('/')[n]

    if path == 'response/status':
        context['response_status'] = int(value.decode('utf-8'))
    elif path == 'response/body':
        context['response_body'].write(value)
    elif path == 'response/body':
        context['response_stream'].write(value)
    elif path.startswith('response/header/'):
        clean = value.rstrip(b'\n').decode('utf-8')
        context['response_headers'][nrd(2)] = clean
    elif path.startswith('response/cookie/'):
        clean = value.rstrip(b'\n').decode('utf-8')
        context['response_cookies'][nrd(2)] = clean
    else:
        raise ValueError(f'Unknown path {path!r}')


def is_readable(path):
    return path.startswith('request/')


def is_writable(path):
    return path.startswith('response/')


def get_manager(resource, context):
    """
    Return an async context manager capable of manage the given
    resource.
    """
    try:
        view, path = resource.split(':')
    except:
        log.error(f"Invalid resource %r", resource)
        raise

    @contextlib.asynccontextmanager
    async def manager():
        """
        Manage the given `resource` as an async context manager.

        This context manager has different behavior depending on the
        `view` and/or `path` of the resource.

        As a context manager it has three sections:
        - Before `yield`: Prepare, if needed, the physical resource on
        disk.
        - `yield`: Return a `ResourceManager` object containing the
        shell representation of the object and the coroutine
        consuming/generating the resource data.
        - After `yield`: Cleanup any disk resource.
        """
        if view == 'raw':
            if not is_readable(path):
                raise ValueError(f'Non-readable path "{path}".')
            else:
                value = await get_value(context, path)
                yield ResourceManager(
                    kapow_repr=resource,
                    shell_repr=value.decode('utf-8'),
                    coro=asyncio.sleep(0))
        elif view == 'value':
            if not is_readable(path):
                raise ValueError(f'Non-readable path "{path}".')
            else:
                value = await get_value(context, path)
                yield ResourceManager(
                    kapow_repr=resource,
                    shell_repr=shell_quote(value.decode('utf-8')),
                    coro=asyncio.sleep(0))
        elif view == 'fifo':
            # No race condition here? Shut your ass!!
            # https://stackoverflow.com/a/1430566
            filename = tempfile.mktemp()
            os.mkfifo(filename)
            if path.startswith('response/stream'):
                async def manage_fifo():
                    initialized = False
                    try:
                        async with aiofiles.open(filename, 'rb') as fifo:
                            while True:
                                if path.endswith('/lines'):
                                    chunk = await fifo.readline()
                                else:
                                    chunk = await fifo.read(1024*10)
                                if chunk:
                                    if not initialized:
                                        # Give a chance to other coroutines
                                        # to write changes to response
                                        # (headers, etc)
                                        await asyncio.sleep(0)
                                        response = web.StreamResponse(
                                            status=200,
                                            headers=context["response_headers"],
                                            reason="OK")
                                        for name, value in context["response_cookies"]:
                                            response.set_cookie(name, value)
                                        context["stream"] = response
                                        await response.prepare(context["request"])
                                        initialized = True
                                    await response.write(chunk)
                                else:
                                    break
                    finally:
                        os.unlink(filename)
            else:
                async def manage_fifo():
                    try:
                        if is_readable(path):
                            async with aiofiles.open(filename, 'wb') as fifo:
                                await fifo.write(await get_value(context, path))
                        elif is_writable(path):
                            async with aiofiles.open(filename, 'rb') as fifo:
                                buf = io.BytesIO()
                                while True:
                                    chunk = await fifo.read(128)
                                    if not chunk:
                                        break
                                    buf.write(chunk)
                                await set_value(context, path, buf.getvalue())
                        else:
                            raise RuntimeError('WTF!')
                    finally:
                        os.unlink(filename)
            yield ResourceManager(
                kapow_repr=resource,
                shell_repr=shell_quote(filename),
                coro=manage_fifo(),
                fifo_path=filename,
                fifo_direction='read' if is_readable(path) else 'write')
        elif view == 'file':
            with tempfile.NamedTemporaryFile(mode='w+b', buffering=0) as tmp:
                if is_readable(path):
                    value = await get_value(context, path)
                    tmp.write(value)
                    tmp.flush()

                yield ResourceManager(
                    kapow_repr=resource,
                    shell_repr=shell_quote(tmp.name),
                    coro=asyncio.sleep(0))

                if is_writable(path):
                    tmp.seek(0)
                    await set_value(context, path, tmp.read())
        else:
            raise ValueError(f'Unknown view type {view}')

    return manager


class KapowTemplate(Template):
    """Shell-code templating for @view:path variables substitution"""

    delimiter = '@'
    idpattern = r'(?a:[_a-z][_a-z0-9]*:[_a-z][-_a-z0-9/]*)'

    async def run(self, context):
        """Run this template allocating and deallocating resources."""
        async with contextlib.AsyncExitStack() as stack:
            # Initialize all resources creating a mapping
            resources = dict()  # resource: (shell_repr, manager)
            for match in self.pattern.findall(self.template):
                delim, resource, *rest = match
                if not resource:  # When is braced
                    resource = rest[0]
                if delim and not resource and rest == ['', '']:
                    # Escaped
                    continue
                if resource not in resources:
                    try:
                        manager = get_manager(resource, context)
                    except:
                        log.error(f"Invalid match %r, %r, %r", delim, resource, rest)
                        raise
                    resources[resource] = await stack.enter_async_context(manager())

            code = self.substitute(**{k: v.shell_repr
                                      for k, v in resources.items()})

            log.debug("Creating tasks")
            manager_tasks = {asyncio.create_task(v.coro): v
                             for k, v in resources.items()}

            await asyncio.sleep(0)
            log.debug("Creating subprocess")
            shell_task = await asyncio.create_subprocess_shell(
                code,
                executable=os.environ.get('SHELL', '/bin/sh'))

            log.debug("Waiting for subprocess")
            await shell_task.wait()  # Run the subshell process

            done, pending = await asyncio.wait(manager_tasks.keys(), timeout=0.1)

            if pending:
                for task in pending:
                    resource = manager_tasks[task]
                    if resource.fifo_path is not None:
                        log.debug(f"Trying to stop %s", resource.kapow_repr)
                        if resource.fifo_direction == 'write':
                            os.system(f"echo -n > {resource.fifo_path} &")
                        elif resource.fifo_direction == 'read':
                            os.system(f"cat {resource.fifo_path} > /dev/null &")
                        else:
                            raise ValueError("Unknown direction")
                    else:
                        log.debug(f"Non-fifo resource pending!! %s", resource.kapow_repr)

                log.debug("Waiting for pending resources...")
                await asyncio.wait(pending)

            await asyncio.sleep(0)


def create_context(request):
    """Create a request context with default values."""
    context = dict()
    context["request"] = request
    context["stream"] = None
    context["response_body"] = io.BytesIO()
    context["response_status"] = 200
    context["response_headers"] = dict()
    context["response_cookies"] = dict()
    return context


async def response_from_context(context):
    """Return the appropia aiohttp response for a given context."""
    if context["stream"] is not None:
        await context["stream"].write_eof()
        return context["stream"]
    else:
        body = context["response_body"].getvalue()
        status = context["response_status"]
        headers = context["response_headers"]
        cookies = context["response_cookies"]

        # Content-Type guessing (for demo only)
        if "Content-Type" not in headers:
            try:
                body = body.decode("utf-8")
            except UnicodeDecodeError:
                pass
            else:
                headers["Content-Type"] = "text/plain"

        response = web.Response(body=body, status=status, headers=headers)
        for name, value in cookies.items():
            response.set_cookie(name, value)

        return response


def generate_endpoint(code, path=None):
    """Return an aiohttp-endpoint coroutine to run kapow `code`."""
    async def endpoint(request):
        context = create_context(request)
        log.debug("Running endpoint %r", path)
        try:
            await KapowTemplate(code).run(context)  # Will change context
        except:
            log.exception("Template crashed!")
        log.debug("Endpoint finished, creating response %r", path)
        response = await response_from_context(context)
        log.debug("Responding %r", path)
        return response
    return endpoint


def path_server(path):
    """Return an aiohttp-endpoint coroutine to serve the file in `path`."""
    # At initialization check
    if not os.path.isfile(path):
        raise NotImplementedError("Only files can be served.")

    async def serve_path(request):
        # Per request check
        if os.path.isdir(path):
            raise NotImplementedError("Cannot serve whole directories yet.")
        return web.FileResponse(path)
    return serve_path


########################################################################
#                              Webserver                               #
########################################################################

def register_code_endpoint(app, methods, pattern, code):
    """Register all needed endpoints for the defined endpoint code."""
    endpoint = generate_endpoint(code, pattern)
    for method in methods:  # May be '*'
        app.add_routes([web.route(method, pattern, endpoint)])


def register_path_endpoint(app, methods, pattern, path):
    """Register all needed endpoints for the defined file."""
    for method in methods:
        if method != 'GET':
            raise ValueError("Invalid method for serving files.")
        else:
            app.add_routes([web.static(pattern, path)])


async def debug_tasks():
    while True:
        await asyncio.sleep(1)
        log.debug("Tasks: %s | Threads: %s",
                  len(asyncio.Task.all_tasks()),
                  threading.active_count())


async def start_background_tasks(app):
    app["debug_tasks"] = app.loop.create_task(debug_tasks())


@click.command()
@click.option('--expression', '-e')
@click.option('--verbose', '-v', count=True)
@click.argument('program', type=click.File(), required=False)
@click.pass_context
def main(ctx, program, verbose, expression):
    """Run the kapow server with the given command-line parameters."""
    if program is None and expression is None:
        program = sys.stdin

    app = web.Application(client_max_size=1024*1024*1024)

    if verbose == 0:
        _print = lambda _: None
        logging.basicConfig(stream=sys.stderr, level=logging.ERROR)
    else:
        _print = lambda s: print(s, file=sys.stderr)
        if verbose == 1:
            logging.basicConfig(stream=sys.stderr, level=logging.INFO)
        else:
            logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)
            if verbose > 2:
                app.on_startup.append(start_background_tasks)

    source = expression if program is None else program.read()

    for ep, _, _ in KAPOW_PROGRAM.scanString(source):
        methods = ep.method.asList()[0].split('|')
        pattern = ''.join(ep.urlpattern)
        if ep.body:
            log.info(f"Registering [code] methods=%r pattern=%r", methods, pattern)
            register_code_endpoint(app, methods, pattern, ep.body)
        else:
            log.info(f"Registering [path] methods=%r pattern=%r", methods, pattern)
            register_path_endpoint(app, methods, pattern, ep.path)
    web.run_app(app, print=_print)

if __name__ == '__main__':
    main()
