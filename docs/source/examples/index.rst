Using a .pow file
+++++++++++++++++

A .pow file is a plain text with shell instructions, usually, you can use Kapow!

**Starting Kapow! using .pow file**

.. code-block:: console
   :linenos:

   $ kapow server example.pow

With the example.pow:

.. code-block:: console
   :linenos:

   $ cat example.pow
   #
   # This is a simple example of a .pow file
   #
   echo "[*] Starting my script"

   # We add 2 Kapow! routes
   kapow route add '/my/route' -c 'echo "hello world" | kapow set /response/body'
   kapow route add -X POST /echo -c '"$(kapow get /request/body)" | kapow set /response/body'

.. note::

    Every manage task you with Kapow! could be done by .pow file

Load more than 1 .pow file
++++++++++++++++++++++++++

You can load more than one .pow file at time. This can help you have your .pow files ordered.

.. code-block:: console
   :linenos:

   $ ls pow-files/
   example-1.pow   example-2.pow
   $ kapow server <(cat pow-files/*.pow)

Add a new route
+++++++++++++++

.. warning::

    Be aware when you defined more than routes in same path, only first routed added will be resolved.

    For example, if you add these routes:

    1. http://localhost:8080/echo
    2. http://localhost:8080/echo/{message}

    Only first one route will be resolved.

**GET route**

Defining route:

.. code-block:: console
   :linenos:

   $ kapow route add '/my/route' -c 'echo "hello world" | kapow set /response/body'

Calling route:

.. code-block:: console
   :linenos:

   $ curl http://localhost:8080/my/route
   hello world

**POST route**

Defining route:

.. code-block:: console
   :linenos:

   $ kapow route add -X POST /echo -c '"$(kapow get /request/body)" | kapow set /response/body'

Calling route:

.. code-block:: console
   :linenos:

   $ curl -d "hello world" -X POST http://localhost:8080/echo
   hello world%

**Adding URL params**

Defining route:

.. code-block:: console
   :linenos:

   $ kapow route add '/echo/{message}' -c '"$(kapow get /request/matches/message)" | kapow set /response/body'

Calling route:

.. code-block:: console
   :linenos:

   $ curl http://localhost:8080/echo/hello%20world
   hello world%


Listing routes
++++++++++++++

You can list active route in kapow! server.

.. _examples_listing_routes:

.. code-block:: console
   :linenos:

   $ kapow route list
   [{"id":"20c98328-0b82-11ea-90a8-784f434dfbe2","method":"GET","url_pattern":"/echo/{message}","entrypoint":"/bin/sh -c","command":"echo \"$(kapow get /request/matches/message)\" | kapow set /response/body","index":0}]

Or, for pretty output, you can use :samp:`jq`:

.. code-block:: console
   :linenos:

   $ kapow route list | jq
   [
     {
       "id": "20c98328-0b82-11ea-90a8-784f434dfbe2",
       "method": "GET",
       "url_pattern": "/echo/{message}",
       "entrypoint": "/bin/sh -c",
       "command": "\"$(kapow get /request/matches/message)\" | kapow set /response/body",
       "index": 0
     }
   ]


.. note::

    Kapow! server has a administration interface, by default, listen at **localhost:8081**


Deleting routes
+++++++++++++++

If we want to delete a route you need their ID. Using de :ref:`listing routes example <examples_listing_routes>`, you can delete the route by typing:

.. code-block:: console
   :linenos:

   $ kapow route remove 20c98328-0b82-11ea-90a8-784f434dfbe2

Writing multiline .pow files
++++++++++++++++++++++++++++

Some time you need to write more complex actions. So you can write multiline commands:

.. code-block:: console
   :linenos:

   $ cat multiline.pow
   kapow route add /log_and_stuff - <<-'EOF'
       echo this is a quite long sentence and other stuff | tee log.txt | kapow set /response/body
       cat log.txt | kapow set /response/body
   EOF

.. warning::

    Be aware with the **"-"** at the end of Kapow! command. It allows to read commands from the :samp:`stdin`.

.. warning::

    Multiline depends of the shell you're using (Bash by default). If you want to learn more of multiline see: `Here Doc <https://en.wikipedia.org/wiki/Here_document>`_


Add or modify an HTTP Header
++++++++++++++++++++++++++++

Some times you want add some extra HTTP header to response.

In this example we'll adding the security header "nosniff" in a sniff.pow:

.. code-block:: console
   :linenos:

   $ cat sniff.pow
   kapow route add /sec-hello-world - <<-'EOF'
       kapow set /response/headers/X-Content-Type-Options "nosniff"

       echo "more secure hello world" | kapow set /response/body
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
   < Content-Type: text/plain; charset=utf-8
   <
   more secure hello world

.. note::

    You can read more about nosniff header `here <https://developer.mozilla.org/es/docs/Web/HTTP/Headers/X-Content-Type-Options>`_.

Modify JSON by using shell
++++++++++++++++++++++++++

.. note::

    Nowadays webservices are json based so making your script json aware is probably a good choice. In order to be able to extract data from and compose json documents from a script you can use
    `jq <https://https://stedolan.github.io/jq/>`_.

**Example 1**

In this example our Kapow! service will receive a JSON value with an incorrect date, then our .pow file will fix then and return the correct value to the user.

.. code-block:: console
   :linenos:

   $ cat fix_date.pow
   kapow route add -X POST '/fix-date' - <<-'EOF'
       kapow set /response/headers/Content-Type "application/json"
       echo "$(kapow get /request/body)" | jq --arg newdate $(date +"%Y-%m-%d_%H-%M-%S") '.incorrectDate=$newdate' | kapow set /response/body
   EOF

Call service with curl:

.. code-block:: console
   :linenos:

   $ curl -X POST http://localhost:8080/fix-date -H "Content-Type: application/json" -d '{"incorrectDate": "no way"}'
   {
      "incorrectDate": "2019-11-22_10-42-06"
   }

**Example 2**

In this example we extract the name field from the incoming json document in order to generate a two attribute json response.

.. code-block:: console

   $ cat echo-attribute.pow
   kapow route add -X POST '/echo-attribute' - <<-'EOF'
      JSON_WHO=$(kapow get /request/body | jq -r .name)

      kapow set /response/headers/Content-Type "application/json"
      kapow set /response/status 200

      jq --arg greet "Hello" --arg value "${JSON_WHO:-World}" -n \{greet:\$greet\,to:\$value\} | kapow set /response/body

   EOF

Call service with curl:

.. code-block:: console
   :linenos:
   :emphasize-lines: 4

   $ curl -X POST http://localhost:8080/echo-attribute -H "Content-Type: application/json" -d '{"name": "MyName"}'
   {
     "greet": "Hello",
     "to": "MyName"
   }

Upload files
++++++++++++

**Example 1**

Upload a file using Kapow! is very simple:

.. code-block:: console
   :linenos:

   $ cat upload.pow
   kapow route add -X POST '/upload-file' - <<-'EOF'
       echo "$(kapow get /request/files/data/content) | kapow set /response/body
   EOF

.. code-block:: console
   :linenos:

   $ cat results.json
   {"hello": "world"}
   $  curl	-X POST -H "Content-Type: multipart/form-data" -F "data=@results.json" http://localhost:8080/upload-file
   {"hello": "world"}

**Example 2**

In this example we respond back with the line count of the file received in the request:

.. code-block:: console
   :linenos:

   $ cat count-file-lines.pow
   kapow route add -X POST '/count-file-lines' - <<-'EOF'

      # Get sent file
      FNAME="$(kapow get /request/files/myfile/filename)"

      # Counting file lines
      LCOUNT="$(kapow get /request/files/myfile/content | wc -l)"

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

Protecting again Command Injection Attacks
++++++++++++++++++++++++++++++++++++++++++

When you resolve variable values be careful to *escape* by using double quotes. Otherwise you could be vulnerable to **command injection attack**.

**This examples is VULNERABLE to command injection**

In this example, an attacker can execute arbitrary command.

.. code-block:: console
   :linenos:

   $ cat command-injection.pow
   kapow route add '/vulnerable/{value}' - <<-'EOF'
        ls "$(kapow get /request/matches/value)" | kapow set /response/body
   EOF

Exploding using curl:

.. code-block:: console
   :linenos:

   $ curl "http://localhost:8080/vulnerable/;echo%20hello"

**This examples is NOT VULNERABLE to command injection**

Be aware of we add double quotes when we recover *value* data from url:

.. code-block:: console
   :linenos:

   $ cat command-injection.pow
   kapow route add '/vulnerable/{value}' - <<-'EOF'
        ls "$(kapow get /request/matches/value)" | kapow set /response/body
   EOF

.. note::

   If want to read more about command injection, you can check `OWASP site <https://www.owasp.org/index.php/Command_Injection>`_

Sending HTTP error codes
++++++++++++++++++++++++

You can specify custom status code for HTTP response:

.. code-block:: console
   :linenos:

   $ cat error.pow
   kapow route add '/error' - <<-'EOF'
       kapow set /response/status 401
       echo "401 error" | kapow set /response/body
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
   kapow route add '/redirect' - <<-'EOF'
       kapow set /response/headers/Location 'http://google.com'
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


How to execute two processes parallel
+++++++++++++++++++++++++++++++++++++

We want to :samp:`ping` two machines parallel. Kapow! get IPs from query params:

.. code-block:: console
   :linenos:

   $ cat parallel.pow
   kapow route add /parallel/{ip1}/{ip2} - <<-'EOF'
       ping -c 1 "$(kapow get /request/matches/ip1)" | kapow set /response/body &
       ping -c 1 "$(kapow get /request/matches/ip2)" | kapow set /response/body &
       wait
   EOF

Calling with curl:

.. code-block:: console
   :linenos:

    $ curl -v http://localhost:8080/parallel/10.0.0.1/10.10.10.1

Manage cookies
++++++++++++++

Sometimes you need track down some user state. Kapow! allows you manage Request/Response Cookies.

Next example we'll set a cookie:

.. code-block:: console
   :linenos:

   $ cat cookie.pow
   kapow route add /setcookie - <<-'EOF'
      CURRENT_STATUS="$(kapow get /request/cookies/kapow-status)"

      if [ -z "$CURRENT_SATUS" ]; then
         kapow set /response/cookies/Kapow-Status "Kapow Cookie Set"
      fi

      echo "Ok" | kapow set /response/body
   EOF

Calling with curl:

.. code-block:: console
   :linenos:
   :emphasize-lines: 11

   $ curl -v http://localhost:8080/set-cookie
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
