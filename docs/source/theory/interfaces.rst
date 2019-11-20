Kapow! HTTP Interfaces
======================

User Interface
--------------

The User HTTP interface is used to serve final user requests.

By default is bind to address `0.0.0.0` and port `8080`.


Control Interface
-----------------

The Control HTTP interface is used by the command `kapow route` to
administer the list of system routes.

By default is bind to address `127.0.0.1` and port `8081`.


Data Interface
--------------

The Data HTTP interface is used by the commands `kapow get` and `kapow
set` to exchange the data for a particular request.

By default is bind to address `127.0.0.1` and port `8082`.
