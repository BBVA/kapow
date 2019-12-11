Routes
======

A *Kapow!* route specifies the matching criteria for an incoming request on
the `User HTTP Interface`, and the details to handle it.

*Kapow!* implements a *route table* where all routes reside.

A route can be set like this:

.. code-block:: console

   $ kapow route add \
      -X POST \
      '/register/{username}' \
      -e '/bin/bash -c' \
      -c 'touch /var/lib/mydb/"$(kapow get /request/matches/username)"' \
      | jq
   {
      "id": "deadbeef-0d09-11ea-b18e-106530610c4d",
      "method": "POST",
      "url_pattern": "/register/{username}",
      "entrypoint": "/bin/bash -c",
      "command": "touch /var/lib/mydb/\"$(kapow get /request/matches/username)\""
   }

Let's use this example to discuss its elements.


Elements
--------

``id`` Route Element
~~~~~~~~~~~~~~~~~~~~

Uniquely identifies each route. It is used for instance by ``kapow route remove
<route_id>``.

.. note::

   The current implementation of *Kapow!* autogenerates a `UUID` for this field.
   In the future the user will be able to specify a custom value.


``method`` Route Element
~~~~~~~~~~~~~~~~~~~~~~~~

Specifies the HTTP method for the route to match the incoming request.

Note that the route shown above will only match a ``POST`` request.


``url_pattern`` Route Element
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

It matches the `path` component of the `URL` of the incoming request.

It can contain regex placeholders for easily capturing fragments of the path.

In the route shown above, a request with a URL ``/register/joe`` would match,
assigning `joe` to the placeholder ``username``.

*Kapow!* leverages Gorilla Mux for managing routes.  For the full story, see
https://github.com/gorilla/mux#examples


``entrypoint`` Route Element
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This sets the executable to be spawned, along with any arguments required.

In the route shown above, the entrypoint that will be run is ``/bin/bash -c``,
which is an incomplete recipe.  It is then completed by the `command` element.

.. todo::

   link to command below


.. note::

   The semantics of this element closely match `Docker`'s ``ENTRYPOINT`` directive.

.. todo::

   link to Docker docu


``command`` Route Element
~~~~~~~~~~~~~~~~~~~~~~~~~

This is an optional last argument to be passed to the ``entrypoint``.

In the route shown above, it completes the `entrypoint` to form the final
incantation to be executed:

.. todo::

   link to entrypoint above

.. code-block:: bash

   /bin/bash -c 'touch /var/lib/mydb/"$(kapow get /request/matches/username)"'

.. note::

   The semantics of this element closely match `Docker`'s ``COMMAND`` directive.

.. todo::

   link to Docker docu


Matching Algorithm
------------------

*Kapow!* leverages Gorilla Mux for this task. You can see the gory details in
their documentation.


.. todo::

   link to Gorilla Mux docu
