Managing Routes
===============

Adding New Routes
-----------------

.. warning::

    Be aware that if you register more than one route with exactly the
    same path, only the first route added will be used.

GET route
+++++++++

Defining a route:

.. code-block:: console
   :linenos:

   $ kapow route add /my/route -c 'echo hello world | kapow set /response/body'


Calling route:

.. code-block:: console
   :linenos:

   $ curl http://localhost:8080/my/route
   hello world

POST route
++++++++++

Defining a route:

.. code-block:: console
   :linenos:

   $ kapow route add -X POST /echo -c 'kapow get /request/body | kapow set /response/body'


Calling a route:

.. code-block:: console
   :linenos:

   $ curl -d 'hello world' -X POST http://localhost:8080/echo
   hello world


Capturing Parts of the URL
++++++++++++++++++++++++++

Defining a route:

.. code-block:: console
   :linenos:

   $ kapow route add '/echo/{message}' -c 'kapow get /request/matches/message | kapow set /response/body'


Calling a route:

.. code-block:: console
   :linenos:

   $ curl http://localhost:8080/echo/hello%20world
   hello world


Listing Routes
--------------

You can list the active routes in the *Kapow!* server.

.. _listing-routes-example:

.. code-block:: console
   :linenos:

   $ kapow route list
   [{"id":"20c98328-0b82-11ea-90a8-784f434dfbe2","method":"GET","url_pattern":"/echo/{message}","entrypoint":"/bin/sh -c","command":"kapow get /request/matches/message | kapow set /response/body"}]

Or, if you want human-readable output, you can use :program:`jq`:

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

   *Kapow!* has a :ref:`https-control-interface`, bound by default to
   ``localhost:8081``.


Deleting Routes
---------------

You need the ID of a route to delete it.
Running the command used in the :ref:`listing routes example
<listing-routes-example>`, you can obtain the ID of the route, and then delete
it by typing:

.. code-block:: console
   :linenos:

   $ kapow route remove 20c98328-0b82-11ea-90a8-784f434dfbe2


