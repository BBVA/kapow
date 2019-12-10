Using a pow File
++++++++++++++++

A ``pow`` file is just a ``bash`` script, where you make calls to the ``kapow route``
command.


**Starting** *Kapow!* **using a** ``pow`` **file**

.. code-block:: console
   :linenos:

   $ kapow server example.pow

With the example.pow:

.. code-block:: console
   :linenos:

   $ cat example.pow
   #
   # This is a simple example of a pow file
   #
   echo '[*] Starting my script'

   # We add 2 Kapow! routes
   kapow route add /my/route -c 'echo hello world | kapow set /response/body'
   kapow route add -X POST /echo -c 'kapow get /request/body | kapow set /response/body'

.. note::

   *Kapow!* can be fully configured using just ``pow`` files


Load More Than One ``pow`` File
+++++++++++++++++++++++++++++++

You can load more than one ``pow`` file at time.  This can help you keep your
``pow`` files tidy.

.. code-block:: console
   :linenos:

   $ ls pow-files/
   example-1.pow   example-2.pow
   $ kapow server <(cat pow-files/*.pow)


Add a New Route
+++++++++++++++

.. warning::

    Be aware that if you register more than one route with exactly the
    same path, only the first route added will be used.


**GET route**

Defining route:

.. code-block:: console
   :linenos:

   $ kapow route add /my/route -c 'echo hello world | kapow set /response/body'


Calling route:

.. code-block:: console
   :linenos:

   $ curl http://localhost:8080/my/route
   hello world


**POST route**

Defining route:

.. code-block:: console
   :linenos:

   $ kapow route add -X POST /echo -c 'kapow get /request/body | kapow set /response/body'


Calling route:

.. code-block:: console
   :linenos:

   $ curl -d 'hello world' -X POST http://localhost:8080/echo
   hello world


**Adding URL params**

Defining route:

.. code-block:: console
   :linenos:

   $ kapow route add '/echo/{message}' -c 'kapow get /request/matches/message | kapow set /response/body'


Calling route:

.. code-block:: console
   :linenos:

   $ curl http://localhost:8080/echo/hello%20world
   hello world


Listing Routes
++++++++++++++

You can list the active routes in the *Kapow!* server.

.. _examples_listing_routes:

.. code-block:: console
   :linenos:

   $ kapow route list
   [{"id":"20c98328-0b82-11ea-90a8-784f434dfbe2","method":"GET","url_pattern":"/echo/{message}","entrypoint":"/bin/sh -c","command":"kapow get /request/matches/message | kapow set /response/body"}]

Or, if you want human-readable output, you can use :samp:`jq`:

.. code-block:: console
   :linenos:

   $ kapow route list | jq
   [
     {
       "id": "20c98328-0b82-11ea-90a8-784f434dfbe2",
       "method": "GET",
       "url_pattern": "/echo/{message}",
       "entrypoint": "/bin/sh -c",
       "command": "kapow get /request/matches/message | kapow set /response/body",
     }
   ]


.. note::

    *Kapow!* has an `HTTP` admin interface, by default listening at **localhost:8081**


Deleting Routes
+++++++++++++++

You need the ID of a route to delete it.
Using the :ref:`listing routes example <examples_listing_routes>`, you can
obtain the ID of the route, and then delete it by typing:

.. code-block:: console
   :linenos:

   $ kapow route remove 20c98328-0b82-11ea-90a8-784f434dfbe2


Writing Multiline ``pow`` Files
+++++++++++++++++++++++++++++++

If you need to write more complex actions, you can leverage multiline commands:

.. code-block:: console
   :linenos:

   $ cat multiline.pow
   kapow route add /log_and_stuff - <<-'EOF'
   	echo this is a quite long sentence and other stuff | tee log.txt | kapow set /response/body
   	cat log.txt | kapow set /response/body
   EOF

.. warning::

    Be aware of the **"-"** at the end of the ``kapow route add`` command.
    It tells ``kapow route add`` to read commands from the :samp:`stdin`.

.. warning::

    If you want to learn more of multiline usage, see: `Here Doc
    <https://en.wikipedia.org/wiki/Here_document>`_


Add or Modify an HTTP Header
++++++++++++++++++++++++++++

You may want to add some extra HTTP header to the response.

In this example we'll be adding the header ``X-Content-Type-Options`` to the response.

.. code-block:: console
   :linenos:

   $ cat sniff.pow
   kapow route add /sec-hello-world - <<-'EOF'
   	kapow set /response/headers/X-Content-Type-Options nosniff
   	kapow set /response/headers/Content-Type text/plain

   	echo this will be interpreted as plain text | kapow set /response/body
   EOF

   $ kapow server nosniff.pow

Testing with curl:

.. code-block:: console
   :emphasize-lines: 11
   :linenos:

   $ curl -v http://localhost:8080/sec-hello-world
   *   Trying ::1...
   * TCP_NODELAY set
   * Connected to localhost (::1) port 8080 (#0)
   > GET /sec-hello-word HTTP/1.1
   > Host: localhost:8080
   > User-Agent: curl/7.54.0
   > Accept: */*
   >
   < HTTP/1.1 200 OK
   < X-Content-Type-Options: nosniff
   < Date: Wed, 20 Nov 2019 10:56:46 GMT
   < Content-Length: 24
   < Content-Type: text/plain
   <
   this will be interpreted as plain text

.. warning::

   Please be aware that if you don't explicitly specified the value of
   the ``Content-Type`` header, *Kapow!* will guess it, effectively
   negating the effect of the ``X-Content-Type-Options`` header.

.. note::

    You can read more about the ``X-Content-Type-Options: nosniff`` header `here
    <https://developer.mozilla.org/es/docs/Web/HTTP/Headers/X-Content-Type-Options>`_.


Modify JSON by Using Shell Commands
+++++++++++++++++++++++++++++++++++

.. note::

    Nowadays Web services are JSON-based so making your script JSON aware is
    probably a good choice.  In order to be able to extract data from a JSON
    document as well as composing JSON documents from a script, you can leverage
    `jq <https://https://stedolan.github.io/jq/>`_.


**Example 1**

In this example our *Kapow!* service will receive a JSON value with an incorrect
date, then our ``.pow`` file will fix it and return the correct value to the user.

.. code-block:: console
   :linenos:

   $ cat fix_date.pow
   kapow route add -X POST /fix-date - <<-'EOF'
   	kapow set /response/headers/Content-Type application/json
   	kapow get /request/body | jq --arg newdate "$(date +'%Y-%m-%d_%H-%M-%S')" '.incorrectDate=$newdate' | kapow set /response/body
   EOF

Call the service with ``curl``:

.. code-block:: console
   :linenos:

   $ curl -X POST http://localhost:8080/fix-date -H 'Content-Type: application/json' -d '{"incorrectDate": "no way, Jose"}'
   {
      "incorrectDate": "2019-11-22_10-42-06"
   }


**Example 2**

In this example we extract the name field from the incoming JSON document in
order to generate a two-attribute JSON response.

.. code-block:: console

   $ cat echo-attribute.pow
   kapow route add -X POST /echo-attribute - <<-'EOF'
   	JSON_WHO=$(kapow get /request/body | jq -r .name)

   	kapow set /response/headers/Content-Type application/json
   	kapow set /response/status 200

   	jq --arg greet Hello --arg value "${JSON_WHO:-World}" --null-input '{ greet: $greet, to: $value }' | kapow set /response/body
   EOF

Call the service with ``curl``:

.. code-block:: console
   :linenos:
   :emphasize-lines: 4

   $ curl -X POST http://localhost:8080/echo-attribute -H 'Content-Type: application/json' -d '{"name": "MyName"}'
   {
     "greet": "Hello",
     "to": "MyName"
   }


Upload Files
++++++++++++


**Example 1**

Uploading a file using *Kapow!* is very simple:

.. code-block:: console
   :linenos:

   $ cat upload.pow
   kapow route add -X POST /upload-file - <<-'EOF'
   	kapow get /request/files/data/content | kapow set /response/body
   EOF

.. code-block:: console
   :linenos:

   $ cat results.json
   {"hello": "world"}
   $ curl	-X POST -H 'Content-Type: multipart/form-data' -F data=@results.json http://localhost:8080/upload-file
   {"hello": "world"}


**Example 2**

In this example we respond back with the line count of the file received in the request:

.. code-block:: console
   :linenos:

   $ cat count-file-lines.pow
   kapow route add -X POST /count-file-lines - <<-'EOF'

   	# Get sent file
   	FNAME=$(kapow get /request/files/myfile/filename)

   	# Counting file lines
   	LCOUNT=$(kapow get /request/files/myfile/content | wc -l)

   	kapow set /response/status 200

   	echo "$FNAME has $LCOUNT lines" | kapow set /response/body
   EOF

.. code-block:: console
   :linenos:

   $ cat file.txt
   hello
   World
   $ curl -F "myfile=@file.txt" http://localhost:8080/count-file-lines
   file.txt has        2 lines


Protecting again Parameter Injection Attacks
++++++++++++++++++++++++++++++++++++++++++++

When you resolve variable values be careful to tokenize correctly by using
double quotes.  Otherwise you could be vulnerable to **parameter injection
attacks**.

**This example is VULNERABLE to parameter injection**

In this example, an attacker can inject arbitrary parameters to ``ls``.

.. code-block:: console
   :linenos:

   $ cat command-injection.pow
   kapow route add '/vulnerable/{value}' - <<-'EOF'
   	ls $(kapow get /request/matches/value) | kapow set /response/body
   EOF

Exploiting using curl:

.. code-block:: console
   :linenos:

   $ curl "http://localhost:8080/vulnerable/-lai%20hello"

**This example is NOT VULNERABLE to parameter injection**

Be aware of how we add double quotes when we recover *value* data from the
request:

.. code-block:: console
   :linenos:

   $ cat command-injection.pow
   kapow route add '/not-vulnerable/{value}' - <<-'EOF'
   	ls -- "$(kapow get /request/matches/value)" | kapow set /response/body
   EOF


.. warning::

   Quotes around parameters only protect against injection of additional
   arguments, but not against turning a non-option into option or
   vice-versa.  Note that for many commands we can leverage double-dash
   to signal the end of the options.  See the "Security Concern" section
   on the docs.


Sending HTTP error codes
++++++++++++++++++++++++

You can specify custom status code for HTTP response:

.. code-block:: console
   :linenos:

   $ cat error.pow
   kapow route add /error - <<-'EOF'
   	kapow set /response/status 401
   	echo -n '401 error' | kapow set /response/body
   EOF

Testing with curl:

.. code-block:: console
   :emphasize-lines: 10
   :linenos:

   $ curl -v http://localhost:8080/error
   *   Trying ::1...
   * TCP_NODELAY set
   * Connected to localhost (::1) port 8080 (#0)
   > GET /error HTTP/1.1
   > Host: localhost:8080
   > User-Agent: curl/7.54.0
   > Accept: */*
   >
   < HTTP/1.1 401 Unauthorized
   < Date: Wed, 20 Nov 2019 14:06:44 GMT
   < Content-Length: 10
   < Content-Type: text/plain; charset=utf-8
   <
   401 error

How to redirect using HTTP
++++++++++++++++++++++++++

In this example we'll redirect our users to Google:

.. code-block:: console
   :linenos:

   $ cat redirect.pow
   kapow route add /redirect - <<-'EOF'
   	kapow set /response/headers/Location https://google.com
   	kapow set /response/status 301
   EOF

.. code-block:: console
   :emphasize-lines: 10-11
   :linenos:

   $ curl -v http://localhost:8080/redirect
   *   Trying ::1...
   * TCP_NODELAY set
   * Connected to localhost (::1) port 8080 (#0)
   > GET /redirect HTTP/1.1
   > Host: localhost:8080
   > User-Agent: curl/7.54.0
   > Accept: */*
   >
   < HTTP/1.1 301 Moved Permanently
   < Location: http://google.com
   < Date: Wed, 20 Nov 2019 11:39:24 GMT
   < Content-Length: 0
   <
   * Connection #0 to host localhost left intact


How to Execute Two Processes in Parallel
++++++++++++++++++++++++++++++++++++++++

We want to :samp:`ping` two machines parallel. *Kapow!* gets IPs from query
params:

.. code-block:: console
   :linenos:

   $ cat parallel.pow
   kapow route add '/parallel/{ip1}/{ip2}' - <<-'EOF'
   	ping -c 1 -- "$(kapow get /request/matches/ip1)" | kapow set /response/body &
   	ping -c 1 -- "$(kapow get /request/matches/ip2)" | kapow set /response/body &
   	wait
   EOF

Calling with ``curl``:

.. code-block:: console
   :linenos:

    $ curl -v http://localhost:8080/parallel/10.0.0.1/10.10.10.1

Manage Cookies
++++++++++++++

If you track down some user state, *Kapow!* allows you manage Request/Response
Cookies.

In the next example we'll set a cookie:

.. code-block:: console
   :linenos:

   $ cat cookie.pow
   kapow route add /setcookie - <<-'EOF'
   	CURRENT_STATUS=$(kapow get /request/cookies/kapow-status)

   	if [ -z "$CURRENT_STATUS" ]; then
   	        kapow set /response/cookies/Kapow-Status 'Kapow Cookie Set'
   	fi

   	echo -n OK | kapow set /response/body
   EOF

Calling with ``curl``:

.. code-block:: console
   :linenos:
   :emphasize-lines: 11

   $ curl -v http://localhost:8080/setcookie
   *   Trying ::1...
   * TCP_NODELAY set
   * Connected to localhost (::1) port 8080 (#0)
   > GET /setcookie HTTP/1.1
   > Host: localhost:8080
   > User-Agent: curl/7.54.0
   > Accept: */*
   >
   < HTTP/1.1 200 OK
   < Set-Cookie: Kapow-Status="Kapow Cookie Set"
   < Date: Fri, 22 Nov 2019 10:44:42 GMT
   < Content-Length: 3
   < Content-Type: text/plain; charset=utf-8
   <
   Ok
   * Connection #0 to host localhost left intact
