With *Kapow!* you can publish simple **shell scripts** as **HTTP
services**.

This way you can delegate in others its execution as they don't need access to the host in
which the command is ran. In summary, those repetitive tasks that everybody ask you to do
because they require administrative access to some host can be published through
a *Kapow!* server deployed in that host and the users who need the results can
invoke it directly using an easy to use interface, an HTTP request.

.. image:: https://trello-attachments.s3.amazonaws.com/5c824318411d973812cbef67/5ca1af818bc9b53e31696de3/784a183fba3f24872dd97ee28e765922/Kapow!.png
    :alt: Where Kapow! lives

Installation
============

1. Get a precompiled static binary for your system from our `releases section <https://github.com/BBVA/kapow/releases/latest>`_.

2. Put it in your ``$PATH`` (Linux example):

   .. code-block:: bash

      sudo install kapow1.0.0-rc1_linux_amd64 /usr/bin/kapow


Hello World! -- *Kapow!* style
==============================

In a shell, the traditional `Hello World!` program would be ``echo "Hello World!"``.
Let's publish it through HTTP using *Kapow!*

- First you need a *script file* with the route that will publish your command, lets's call the file ``greet.sh`` and should contain the following code:

.. code-block:: bash

    kapow route add /greet -c 'echo "Hello World!" | kapow set /response/body'

- Start the *Kapow!* server with your script as an argument

.. code-block:: bash

    kapow server greet.sh

- Finally check that all is working as intended using `curl`

.. code-block:: bash

    $ curl http://localhost:8080/greet
    Hello World!

