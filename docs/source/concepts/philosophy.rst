Philosophy
==========


Single Static Binary
--------------------

- Deployment is then as simple as it gets.

- `Docker`-friendly.


Shell Agnostic
--------------

- *Kapow!*, like John Snow, knows nothing, and makes no assumptions about the
  shell you are using.  It only spawns executables.

- You are free to implement a client to the Data API directly if you are so
  inclined.  The spec provides all the necessary details.


Not a Silver Bullet
-------------------

You should not use *Kapow!* if your project requires complex business logic.

If you try to encode business logic in a shell script, you will **deeply**
regret it soon enough.

*Kapow!* is designed for automating simple stuff.


Interoperability over Performance
---------------------------------

We want *Kapow!* to be as performant as possible, but not at the cost of
flexibility.  This is the reason why our :ref:`Data API
<http-data-interface>` leverages HTTP
instead of a lighter protocol for example.

When we have to choose between making things faster or more
interoperable the latter usually wins.
