Handling HTTP Requests
======================

Add or Modify an HTTP Header
----------------------------

You may want to add some extra HTTP header to the response.

In this example we'll be adding the header ``X-Content-Type-Options`` to the response.

.. code-block:: console
   :linenos:

   $ cat sniff-route
   #!/usr/bin/env sh
   kapow route add /sec-hello-world - <<-'EOF'
   	kapow set /response/headers/X-Content-Type-Options nosniff
   	kapow set /response/headers/Content-Type text/plain

   	echo this will be interpreted as plain text | kapow set /response/body
   EOF

   $ kapow server nosniff-route

Testing with :program:`curl`:

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

   Please be aware that if you don't explicitly specify the value of
   the ``Content-Type`` header, *Kapow!* will guess it, effectively
   negating the effect of the ``X-Content-Type-Options`` header.

.. note::

    You can read more about the ``X-Content-Type-Options: nosniff`` header `here
    <https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options>`_.


Upload Files
------------

Example #1
++++++++++

Uploading a file using *Kapow!* is very simple:

.. code-block:: console
   :linenos:

   $ cat upload-route
   #!/usr/bin/env sh
   kapow route add -X POST /upload-file - <<-'EOF'
   	kapow get /request/files/data/content | kapow set /response/body
   EOF

.. code-block:: console
   :linenos:

   $ cat results.json
   {"hello": "world"}
   $ curl	-X POST -H 'Content-Type: multipart/form-data' -F data=@results.json http://localhost:8080/upload-file
   {"hello": "world"}


Example #2
++++++++++

In this example we reply the line count of the file received in the request:

.. code-block:: console
   :linenos:

   $ cat count-file-lines
   #!/usr/bin/env sh
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
   $ curl -F myfile=@file.txt http://localhost:8080/count-file-lines
   file.txt has        2 lines


Sending HTTP error codes
------------------------

You can specify custom status code for `HTTP` response:

.. code-block:: console
   :linenos:

   $ cat error-route
   #!/usr/bin/env sh
   kapow route add /error - <<-'EOF'
   	kapow set /response/status 401
   	echo -n '401 error' | kapow set /response/body
   EOF

Testing with :program:`curl`:

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
--------------------------

In this example we'll redirect our users to `Google`:

.. code-block:: console
   :linenos:

   $ cat redirect
   #!/usr/bin/env sh
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


Manage Cookies
--------------

If you track down some user state, *Kapow!* allows you manage Request/Response
Cookies.

In the next example we'll set a cookie:

.. code-block:: console
   :linenos:

   $ cat cookie
   #!/usr/bin/env sh
   kapow route add /setcookie - <<-'EOF'
   	CURRENT_STATUS=$(kapow get /request/cookies/kapow-status)

   	if [ -z "$CURRENT_STATUS" ]; then
   		kapow set /response/cookies/Kapow-Status 'Kapow Cookie Set'
   	fi

   	echo -n OK | kapow set /response/body
   EOF

Calling with :program:`curl`:

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
   OK
   * Connection #0 to host localhost left intact
