Route Matching
==============

*Kapow!* maintains a :ref:`route <routes>` table with a list of routes as provided by the user,
and uses it to determine which handler an incoming request should be dispatched
to.

Each incoming request is matched against the routes in the route table in
strict order.  For each route in the route table, the criteria are checked.
If the request does not match, the next route in the route list is examined.
