*Kapow!* HTTP Interfaces
========================

``kapow server`` sets up three HTTP server interfaces, each with a distinct and
clear purpose.


.. _http-user-interface:

HTTP User Interface
-------------------

The `HTTP User Interface` is used to serve final user requests.

By default it binds to address ``0.0.0.0`` and port ``8080``, but that can be
changed via the ``--bind`` flag.


.. _http-control-interface:

HTTP Control Interface
----------------------

The `HTTP Control Interface` is used by the command ``kapow route`` to
administer the list of system routes.

By default it binds to address ``127.0.0.1`` and port ``8081``, but that can be
changed via the ``--control-bind`` flag.


.. _http-data-interface:

HTTP Data Interface
-------------------

The `HTTP Data Interface` is used by the commands ``kapow get`` and ``kapow
set`` to exchange the data for a particular request.

By default it binds to address ``127.0.0.1`` and port ``8082``, but that can be
changed via the ``--data-bind`` flag.
