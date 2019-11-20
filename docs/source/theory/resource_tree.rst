The Resource Tree
=================

This is the model that Kapow! uses to expose the internals of the user request
being serviced.


We use this tree to get access to any data that comes in the request,
as well as to compose the response.

We access the resource tree easily with the ``kapow set`` and ``kapow get``
subcommands.


Overview
--------

.. code-block:: plain

    /                               The root of the resource paths tree
    │
    ├─ request                      All information related to the HTTP request.  Read-Only
    │  ├──── method                 Used HTTP Method (GET, POST)
    │  ├──── host                   Host part of the URL
    │  ├──── path                   Complete URL path (URL-unquoted)
    │  ├──── matches                Previously matched URL path parts
    │  │     └──── <name>
    │  ├──── params                 URL parameters (after the "?" symbol)
    │  │     └──── <name>
    │  ├──── headers                HTTP request headers
    │  │     └──── <name>
    │  ├──── cookies                HTTP request cookie
    │  │     └──── <name>
    │  ├──── form                   Form-urlencoded form fields (names only)
    │  │     └──── <name>           Value of the form field with name <name>
    │  ├──── files                  Files uploaded via multi-part form fields (names only)
    │  │     └──── <name>
    │  │           └──── filename   Original file name
    │  │           └──── content    The file content
    │  └──── body                   HTTP request body
    │
    └─ response                     All information related to the HTTP request.  Write-Only
      ├──── status                 HTTP status code
      ├──── headers                HTTP response headers
      │     └──── <name>
      ├──── cookies                HTTP request cookie
      │     └──── <name>
      ├──── body                   Response body.  Mutually exclusive with response/stream
      └──── stream                 Alias for /response/body

Resources
---------

``/request/method``
~~~~~~~~~~~~~~~~~~~
The HTTP method of the incoming request.

Sample usage:

.. code-block:: bash
   $ kapow get /request/method
   GET

``/request/host``
~~~~~~~~~~~~~~~~~

The ``Host`` header as defined in the HTTP/1.1 spec.

Sample usage:

If the user runs:

.. code-block:: bash
   $ curl http://kapow.example:8080


then, when handling the request:

.. code-block:: bash
   $ kapow get /request/host
   kapow.example


``/request/path``
~~~~~~~~~~~~~~~~~
Contains the path substring of the URL.


Sample usage:

.. code-block:: bash
   # GET http://url.example/foo/bar?q=1
   $ kapow get /request/path
   /foo/bar

``/request/matches/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/params/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/headers/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/cookies/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/form/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/files/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/files/<name>/filename``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/files/<name>/content``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/request/body``
~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/response/status``
~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/response/headers/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/response/cookies/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get

``/response/body``
~~~~~~~~~~~~~~~~~~
Sample usage:

.. code-block:: bash
   $ kapow get
