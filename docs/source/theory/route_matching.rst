Route Matching
==============

Route table
-----------

    Kapow! maintains a route table as provided by the user, and uses it to determine
    which handler should an incoming request be dispatched to.

    Each incoming request is matched against the routes in the route table in
    strict order, for each route in the route table, the criteria are checked.
    If the request does not match, the next route in the route list is examined.
