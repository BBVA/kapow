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

    /
    │
    ├─ request
    │  ├──── method                 Used HTTP Method (GET, POST)
    │  ├──── host                   Host part of the URL
    │  ├──── path                   Complete URL path (URL-unquoted)
    │  ├──── matches
    │  │     └──── <name>           Previously matched URL path parts
    │  ├──── params
    │  │     └──── <name>           URL parameters (after the "?" symbol)
    │  ├──── headers
    │  │     └──── <name>           HTTP request headers
    │  ├──── cookies
    │  │     └──── <name>           HTTP request cookie
    │  ├──── form
    │  │     └──── <name>           Value of the form field with name <name>
    │  ├──── files
    │  │     └──── <name>
    │  │           └──── filename   Original file name of the file uploaded in the form field <name>
    │  │           └──── content    The contents of the file uploaded in the form field <name>
    │  └──── body                   HTTP request body
    │
    └─ response
      ├──── status                  HTTP status code
      ├──── headers
      │     └──── <name>            HTTP response headers
      ├──── cookies
      │     └──── <name>            HTTP request cookie
      └──── body                    Response body

Resources
---------

``/request/method`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The HTTP method of the incoming request.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl -X POST http://kapow.example:8080

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/method
   POST


``/request/host`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~

The ``Host`` header as defined in the HTTP/1.1 spec of the incoming
request.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl http://kapow.example:8080

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/host
   kapow.example


``/request/path`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the path substring of the URL.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl http://kapow.example:8080/foo/bar?qux=1

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/path
   /foo/bar

``/request/matches/<name>`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the part of the URL captured by the pattern ``name``.

**Sample Usage**

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

``/request/params/<name>`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the URL parameter ``name``

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl http://kapow.example:8080/foo?myparam=bar

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/params/myparam
   myparam


``/request/headers/<name>`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the HTTP header ``name`` of the incoming request.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl -H X-My-Header=Bar http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/headers/X-My-Header
   Bar


``/request/cookies/<name>`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the HTTP cookie ``name`` of the incoming request.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl --cookie "MYCOOKIE=Bar" http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/cookies/MYCOOKIE
   Bar

``/request/form/<name>`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the field ``name`` of the incoming request.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl -F -d myfield=foo http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/form/myfield
   foo


``/request/files/<name>/filename`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the name of the file uploaded through the incoming request.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl -F -d myfile=@filename.txt http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/files/myfile/filename
   filename.txt


``/request/files/<name>/content`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contents of the file that is being uploaded in the incoming request.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl -F -d myfile=@filename.txt http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/files/myfile/content
   ...filename.txt contents...


``/request/body`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~

Raw contents of the incoming request HTTP body.

**Sample Usage**

If the user runs:

.. code-block:: bash

   $ curl --data-raw foobar http://kapow.example:8080/

then, when handling the request:

.. code-block:: bash

   $ kapow get /request/body
   foobar


``/response/status`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the status code given in the user response.

**Sample Usage**

If during the request handling:

.. code-block:: bash

   $ kapow set /response/status 418

then the response will have the status code ``418 I am a Teapot``.


``/response/headers/<name>`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the header ``name`` in the user response.

**Sample Usage**

If during the request handling:

.. code-block:: bash

   $ kapow set /response/headers/X-My-Header Foo

then the response will contain an HTTP header named ``X-My-Header`` with
value ``Foo``.


``/response/cookies/<name>`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the cookie ``name`` that will be set to the user
response.


**Sample Usage**

If during the request handling:

.. code-block:: bash

   $ kapow set /response/cookies/MYCOOKIE Foo

then the response will set the cookie ``MYCOOKIE`` to the user in
following requests.


``/response/body`` Resource
~~~~~~~~~~~~~~~~~~~~~~~~~~~

Contains the value of the response HTTP body.

**Sample Usage**

If during the request handling:

.. code-block:: bash

   $ kapow set /response/body foobar

then the response will contain ``foobar`` in the body.
