from aiohttp import web
from dataclasses import dataclass
from shlex import quote as shell_quote
from string import Template
import asyncio
import contextlib
import io
import os
import re
import tempfile

from pyparsing import *
import aiofiles
import click

########################################################################
#                                Parser                                #
########################################################################

#
# Endpoint Definition
#

# Method
method = ( Literal('GET')
         | Literal('POST')
         | Literal('PUT')
         | Literal('DELETE')
         | Literal('PATCH') )
multi_method = delimitedList(method, delim="|", combine=True)
method_spec = Combine(Literal('*') | multi_method)

# Pattern
regex = Word(alphas + nums + '\\+*,[]-_')(name="regex")
symbol = Word(alphas)(name="symbol")
p_pattern = Combine('/{' + symbol + Optional(':' + regex) + '}')
p_path = Word('/', alphas + nums + '$-_.+!*\'(),')
urlpattern = Combine(OneOrMore(p_pattern | p_path))(name="urlpattern")

# Body
body = (Suppress('{') + SkipTo(Combine(LineStart() + '}' + LineEnd()))(name="body"))

# Endpoint
endpoint = (Optional(method_spec + Suppress(White()),
                     default='*')(name="method")
           + urlpattern
           + Suppress(White())
           + body)(name="endpoint")

kapow_program = OneOrMore(endpoint)


@dataclass
class ResourceManager:
    shell_repr: str
    coro: object


async def get_value(context, path):
    if path == 'request/method':
        return context['request'].method.encode('utf-8')
    elif path == 'request/path':
        return context['request'].path.encode('utf-8')
    elif path.startswith('request/match'):
        key = path.split('/', 2)[-1]
        return  context['request'].match_info[key].encode('utf-8')
    elif path.startswith('request/param'):
        key = path.split('/', 2)[-1]
        return  context['request'].rel_url.query[key].encode('utf-8')
    elif path.startswith('request/header'):
        key = path.split('/', 2)[-1]
        return  context['request'].headers[key].encode('utf-8')
    elif path.startswith('request/form'):
        key = path.split('/', 2)[-1]
        return (await context['request'].post())[key].encode('utf-8')
    elif path == 'request/body':
        return await context['request'].read()
    else:
        raise ValueError(f'Unknown path {path!r}')

async def set_value(context, path, value):
    if path == 'response/status':
        context['response_status'] = int(value.decode('utf-8'))
    elif path == 'response/body':
        context['response_body'].write(value)
    elif path.startswith('response/header/'):
        key = path.split('/', 2)[-1]
        context['response_headers'][key] = value.rstrip(b'\n').decode('utf-8')
    else:
        raise ValueError(f'Unknown path {path!r}')

def is_readable(path):
    return path.startswith('request/')


def is_writable(path):
    return path.startswith('response/')


def get_manager(resource, context):
    view, path = resource.split(':')

    @contextlib.asynccontextmanager
    async def manager():
        if view == 'value':
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
        elif view == 'source':
            raise NotImplementedError('source view not implemented')
        elif view == 'sink':
            raise NotImplementedError('sink view not implemented')
        else:
            raise ValueError(f'Unknown view type {view}')

    return manager


class KapowTemplate(Template):
    delimiter = '@'
    idpattern = r'(?a:[_a-z][_a-z0-9]*:[_a-z][-_a-z0-9/]*)'

    async def run(self, context):
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

            if False:
                print('-'*80)
                print(code)
                print('-'*80)

            shell_task = await asyncio.create_subprocess_shell(code)
            manager_tasks = [asyncio.create_task(v.coro)
                             for v in resources.values()]

            await shell_task.wait()  # Run the subshell process
            done, pending = await asyncio.wait(manager_tasks, timeout=0)  # Managers commit changes 
            if pending:
                # print(f"Warning: Resources not consumed ({len(pending)})")
                for task in pending:
                    task.cancel()


def create_runner(code):
    return KapowTemplate(code).run


def create_context(request):
    context = dict()
    context["request"] = request
    context["response_body"] = io.BytesIO()
    context["response_status"] = 200
    context["response_headers"] = dict()
    return context


def response_from_context(context):
    return web.Response(
        body=context["response_body"].getbuffer(),
        status=context["response_status"],
        headers=context["response_headers"])


def generate_endpoint(code):
    async def endpoint(request):
        context = create_context(request)
        runner = create_runner(code)
        await runner(context)  # Will change context
        return response_from_context(context)
    return endpoint


########################################################################
#                              Webserver                               #
########################################################################

def register_endpoint(app, methods, pattern, code):
    print(f"Registering methods={methods!r} pattern={pattern!r}")
    endpoint = generate_endpoint(code)
    for method in methods:  # May be '*'
        app.add_routes([web.route(method, pattern, endpoint)])


@click.command()
@click.argument('program', type=click.File())
def main(program):
    app = web.Application()
    for ep, _, _ in kapow_program.scanString(program.read()):
        register_endpoint(app,
                          ep.method.asList()[0].split('|'),
                          ''.join(ep.urlpattern),
                          ep.body)
    web.run_app(app)

if __name__ == '__main__':
    main()
