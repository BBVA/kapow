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
import os
import tempfile

from aiohttp import web
from pyparsing import alphas, nums, White
from pyparsing import LineStart, LineEnd, SkipTo
from pyparsing import Literal, Combine, Word, Suppress
from pyparsing import OneOrMore, Optional, delimitedList
import aiofiles
import click

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
    #: Representation of the resource that can be understood by the shell
    shell_repr: str
    #: Coroutine capable of managing the resource internally
    coro: object


async def get_value(context, path):
    """Return the value of an http resource."""
    def nrd(n):
        """Return the nrd element in a path."""
        return path.split('/', n)[-1]

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
    elif path.startswith('request/form'):
        return (await context['request'].post())[nrd(2)].encode('utf-8')
    elif path == 'request/body':
        return await context['request'].read()
    else:
        raise ValueError(f'Unknown path {path!r}')


async def set_value(context, path, value):
    """
    Write to an http resource.

    File-like resources like `body` get write() calls so they have
    append semantics. Non file-like resources are just set.

    """
    def nrd(n):
        return path.split('/', n)[-1]

    if path == 'response/status':
        context['response_status'] = int(value.decode('utf-8'))
    elif path == 'response/body':
        context['response_body'].write(value)
    elif path == 'response/body':
        context['response_stream'].write(value)
    elif path.startswith('response/header/'):
        clean = value.rstrip(b'\n').decode('utf-8')
        context['response_headers'][nrd(2)] = clean
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
    view, path = resource.split(':')

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
                    shell_repr=value.decode('utf-8'),
                    coro=asyncio.sleep(0))
        elif view == 'value':
            if not is_readable(path):
                raise ValueError(f'Non-readable path "{path}".')
            else:
                value = await get_value(context, path)
                yield ResourceManager(
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
                                    chunk = await fifo.read(128)
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
                                await set_value(context, path, await fifo.read())
                        else:
                            raise RuntimeError('WTF!')
                    finally:
                        os.unlink(filename)
            yield ResourceManager(
                shell_repr=shell_quote(filename),
                coro=manage_fifo())
        elif view == 'file':
            with tempfile.NamedTemporaryFile(mode='w+b', buffering=0) as tmp:
                if is_readable(path):
                    value = await get_value(context, path)
                    tmp.write(value)
                    tmp.flush()

                yield ResourceManager(
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
                _, resource, *_ = match
                if resource not in resources:
                    manager = get_manager(resource, context)
                    resources[resource] = await stack.enter_async_context(manager())

            code = self.substitute(**{k: v.shell_repr
                                      for k, v in resources.items()})

            # print('-'*80)
            # print(code)
            # print('-'*80)

            manager_tasks = [asyncio.create_task(v.coro)
                             for v in resources.values()]
            await asyncio.sleep(0)
            shell_task = await asyncio.create_subprocess_shell(code)

            await shell_task.wait()  # Run the subshell process
            # XXX: Managers commit changes
            _, pending = await asyncio.wait(manager_tasks, timeout=1)
            if pending:
                # print(f"Warning: Resources not consumed ({len(pending)})")
                for task in pending:
                    task.cancel()
            await asyncio.sleep(0)


def create_context(request):
    """Create a request context with default values."""
    context = dict()
    context["request"] = request
    context["stream"] = None
    context["response_body"] = io.BytesIO()
    context["response_status"] = 200
    context["response_headers"] = dict()
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

        # Content-Type guessing (for demo only)
        if "Content-Type" not in headers:
            try:
                body = body.decode("utf-8")
            except UnicodeDecodeError:
                pass
            else:
                headers["Content-Type"] = "text/plain"

        return web.Response(body=body, status=status, headers=headers)


def generate_endpoint(code):
    """Return an aiohttp-endpoint coroutine to run kapow `code`."""
    async def endpoint(request):
        context = create_context(request)
        await KapowTemplate(code).run(context)  # Will change context
        return await response_from_context(context)
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
    print(f"Registering [code] methods={methods!r} pattern={pattern!r}")
    endpoint = generate_endpoint(code)
    for method in methods:  # May be '*'
        app.add_routes([web.route(method, pattern, endpoint)])


def register_path_endpoint(app, methods, pattern, path):
    """Register all needed endpoints for the defined file."""
    print(f"Registering [path] methods={methods!r} pattern={pattern!r}")
    for method in methods:  # May be '*'
        app.add_routes([web.route(method, pattern, path_server(path))])


@click.command()
@click.option('--expression', '-e')
@click.argument('program', type=click.File(), required=False)
@click.pass_context
def main(ctx, program, expression):
    """Run the kapow server with the given command-line parameters."""
    if program is None and expression is None:
        click.echo(ctx.get_help())
        ctx.exit()

    source = expression if program is None else program.read()

    app = web.Application()
    for ep, _, _ in KAPOW_PROGRAM.scanString(source):
        methods = ep.method.asList()[0].split('|')
        pattern = ''.join(ep.urlpattern)
        if ep.body:
            register_code_endpoint(app, methods, pattern, ep.body)
        else:
            register_path_endpoint(app, methods, pattern, ep.path)
    web.run_app(app)


if __name__ == '__main__':
    main()
