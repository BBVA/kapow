Request Life Cycle
==================

This section describes the sequence of events happening for each request
answered by the User HTTP Interface.

#. The user makes a request to the User HTTP Interface

#. The request is matched against the route table

#. ``kapow`` provides a `HANDLER_ID` to identify this request and don't mix it
   with other requests that could be running concurrently.

#. ``kapow`` spawns the binary specified as entrypoint in the matching route

   The default entrypoint is /bin/sh; we'll explain this workflow for now.

   The spawned entrypoint is run with the following variables added to its
   environment:

   - ``KAPOW_HANDLER_ID``: Containing the `HANDLER_ID`
   - ``KAPOW_DATAAPI_URL``: With the URL of the `data interface`
   - ``KAPOW_CONTROLAPI_URL``: With the URL of the `control interface`

#. During the lifetime of the shell, the request and response resources are
   available via these commands:

   - ``kapow get /request/...``
   - ``kapow set /response/...``

   These commands use the aforementioned environment variables to read
   data from the user request and to write the response.

   .. todo::

      link to resource tree

#. When the shell dies, ``kapow`` finalizes the original request
