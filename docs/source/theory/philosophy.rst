Philosophy
==========


We Provide a Single Static Binary
---------------------------------

- Because it makes deployments easier.


Shell Agnostic
--------------

- Kapow! knows nothing about the shell you are using.
- It only spawn binaries.
- You can use anything you want that ends interacting with the `data
  api`.
- This helps with multiplatform and with future higher level tools.


Not a Silver Bullet
-------------------

You should not use Kapow! for projects with complex business logic.

If you try to encode business logic in a shell script you will regret
**deeply**.

Kapow! is for automating simple stuff.
