Request Life Cycle
==================

This section describes the sequence of events happening for each request
answered by the User HTTP Interface.

#. The user makes a request to the User HTTP Interface

#. The request is matched against the route table

#. Kapow! provides a HANDLER_ID to identify this request

#. Kapow! spawns the binary specified as entrypoint in the matching route

   The default entrypoint is /bin/sh; we'll explain this workflow for now.

   The spawned entrypoint is run with the following variables added to its environment:

   - KAPOW_HANDLER_ID
   - KAPOW_DATAAPI_URL
   - KAPOW_CONTROLAPI_URL

#.  During the lifetime of the shell, the request and response resources are available via these commands:

   - kapow get /request/...
   - kapow set /response/...

TODO: link to resource tree

#. When the shell dies, Kapow! finalizes the original request.
