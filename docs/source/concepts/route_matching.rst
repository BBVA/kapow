Route Matching
==============

*Kapow!* maintains a route table with a list of routes as provided by the user,
and uses it to determine which handler should an incoming request be dispatched
to.

.. todo::

   link routes to its section

Each incoming request is matched against the routes in the route table in
strict order.  For each route in the route table, the criteria are checked.
If the request does not match, the next route in the route list is examined.
