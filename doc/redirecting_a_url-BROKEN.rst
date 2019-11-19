Redirecting a URL
=================

The HTTP protocol allows queries to a URL to be redirected to other URL.

We can do them in Kapow! with little effort:

In this example, we read the header ``User-Agent`` and feed it to the response:
.. code-block:: bash

   echo -n 302 | kapow set /response/status
   echo -n http://foobar-url.example | kapow set /response/headers/Location
