Redirecting a URL
=================

The HTTP protocol allows queries to a URL to be redirected to other URL.

We can do them in Kapow! with little effort:

In this example, we read the header ``User-Agent`` and feed it to the response:
.. code-block:: bash

   kapow set /response/headers/Location 'http://example.org'
   kapow set /response/status 301
