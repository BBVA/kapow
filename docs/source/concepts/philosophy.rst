Philosophy
==========


Single Static Binary
--------------------

- Deployment is then as simple as it gets.

- Docker-friendly.


Shell Agnostic
--------------

- Kapow! knows nothing, and makes no assumptions, about the shell you are using.
  It only spawns executables.

- You are free to implement a client to the Data API directly if you are so
  inclined. The spec provides all the necessary details.


Not a Silver Bullet
-------------------

You should not use Kapow! if your project requires complex business logic.

If you try to encode business logic in a shell script, you will **deeply**
regret it.

Kapow! is designed for automating simple stuff.
