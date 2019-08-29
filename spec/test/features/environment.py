import tempfile
import os


def before_scenario(context, scenario):
    # Create the request_handler FIFO
    while True:
        context.handler_fifo_path = tempfile.mktemp() # Safe because using
                                                      # mkfifo
        try:
            os.mkfifo(context.handler_fifo_path)
        except OSError:
            # The file already exist
            pass
        else:
            break


def after_scenario(context, scenario):
    if hasattr(context, 'server'):
        context.server.terminate()
        context.server.wait()

    os.unlink(context.handler_fifo_path)
