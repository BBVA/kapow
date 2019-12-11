*Kapow!* HTTP Interfaces
========================

``kapow server`` sets up three HTTP server interfaces, each with a distinct and
clear purpose.


User Interface
--------------

The :any:`User HTTP Interface` is used to serve final user requests.

By default it binds to address ``0.0.0.0`` and port ``8080``, but that can be
changed via the ``--bind`` flag.


Control Interface
-----------------

The :any:`Control HTTP Interface` is used by the command ``kapow route`` to
administer the list of system routes.

By default it binds to address ``127.0.0.1`` and port ``8081``, but that can be
changed via the ``--control-bind`` flag.


Data Interface
--------------

The :any:`Data HTTP Interface` is used by the commands ``kapow get`` and ``kapow
set`` to exchange the data for a particular request.

By default it binds to address ``127.0.0.1`` and port ``8082``, but that can be
changed via the ``--data-bind`` flag.
