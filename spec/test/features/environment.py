def after_scenario(context, scenario):
    if hasattr(context, 'server'):
        context.server.terminate()
        context.server.wait()
