Writing Headers
===============

The HTTP protocol allows metadata using headers.

Kapow! allows you to set them easily:


In this example, we respond by setting the ``Content-Type`` header
to the value ``application/json``.
.. code-block:: bash

   kapow set /response/headers/Content-Type application/json

We could then return some JSON content:
.. code-block:: bash

   echo '{"data": "some relevant string"}' | kapow set /response/body
