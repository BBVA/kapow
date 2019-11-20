The `Kapow!` Resource Tree
==========================

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

If the user runs:

.. code-block:: bash

   $ curl -X POST http://kapow.example:8080

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/method
   POST


``/request/host``
~~~~~~~~~~~~~~~~~

The ``Host`` header as defined in the HTTP/1.1 spec of the incoming
request.

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

If the user runs:

.. code-block:: bash

   $ curl http://kapow.example:8080/foo/bar?qux=1

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/path
   /foo/bar

``/request/matches/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the part of the URL captured by the pattern ``name``.

Sample usage:

For a route defined like this:

.. code-block:: bash

   $ kapow route add /foo/{mymatch}/bar

if the user runs:

.. code-block:: bash

   $ curl http://kapow.example:8080/foo/1234/bar

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/matches/mymatch
   1234

``/request/params/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the URL parameter ``name``

Sample usage:

If the user runs:

.. code-block:: bash

   $ curl http://kapow.example:8080/foo?myparam=bar

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/params/myparam
   myparam


``/request/headers/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the HTTP header ``name`` of the incoming request.

Sample usage:

If the user runs:

.. code-block:: bash

   $ curl -H X-My-Header=Bar http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/headers/X-My-Header
   Bar


``/request/cookies/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the HTTP cookie ``name`` of the incoming request.

Sample usage:

If the user runs:

.. code-block:: bash

   $ curl --cookie "MYCOOKIE=Bar" http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/cookies/MYCOOKIE
   Bar

``/request/form/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the field ``name`` of the incoming request.

Sample usage:

If the user runs:

.. code-block:: bash

   $ curl -F -d myfield=foo http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/form/myfield
   foo


``/request/files/<name>/filename``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the name of the file uploaded through the incoming request.

Sample usage:

If the user runs:

.. code-block:: bash

   $ curl -F -d myfile=@filename.txt http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/files/myfile/filename
   filename.txt


``/request/files/<name>/content``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contents of the file that is being uploaded in the incoming request.

Sample usage:

If the user runs:

.. code-block:: bash

   $ curl -F -d myfile=@filename.txt http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/files/myfile/content
   ...filename.txt contents...


``/request/body``
~~~~~~~~~~~~~~~~~

Raw contents of the incoming request HTTP body.

Sample usage:

If the user runs:

.. code-block:: bash

   $ curl --data-raw foobar http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/body
   foobar


``/response/status``
~~~~~~~~~~~~~~~~~~~~

Contains the status code given in the user response.

Sample usage:

If during the request handling:

.. code-block:: bash

   $ kapow set /response/status 418

then the response will have the status code ``418 I am a Teapot``.


``/response/headers/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the header ``name`` in the user response.

Sample usage:

If during the request handling:

.. code-block:: bash

   $ kapow set /response/headers/X-My-Header Foo

then the response will contain an HTTP header named ``X-My-Header`` with
value ``Foo``.


``/response/cookies/<name>``
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the cookie ``name`` that will be set to the user
response.


Sample usage:

If during the request handling:

.. code-block:: bash

   $ kapow set /response/cookies/MYCOOKIE Foo

then the response will set the cookie ``MYCOOKIE`` to the user in
following requests.


``/response/body``
~~~~~~~~~~~~~~~~~~

Contains the value of the response HTTP body.

Sample usage:

.. code-block:: bash

   $ kapow set /response/body foobar

then the response will contain ``foobar`` in the body.
