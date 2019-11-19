Reading Headers
===============

The HTTP protocol allows metadata using headers.

Kapow! allows you to read them easily:


In this example, we read the header ``User-Agent`` and feed it
to the response:
.. code-block:: bash

   kapow get /request/headers/User-Agent | kapow set /response/body
