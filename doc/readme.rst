What is Kapow!
==============

Kapow! is an adapter between the world of Pure UNIX® Shell and an HTTP service.

Some tasks are more convenient in the shell, like cloud interactions, or some
administrative tools.  On the other hand, some tasks are more convenient as a
service, like DevSecOps tooling.

Kapow! lies between these two worlds, making your life easier.  Maybe you wonder
about how this kind of magic can happen; if you want to know the nitty-gritty
details, just read our [spec](/spec/).  Or, if you want to know how Kapow! can
help you first, let's start with a common situation.

Think about that awesome command that you use every day, something very
familiar, like ``cloudx storage ls /backups``.  Then someone asks you for an
specific backup, so you ``ssh`` into the host, execute your command, possibly
``grepping`` through its output, copy the result and send it back to him. 
And that's fine... for the 100 first times.

Then you decide, let's use an API for this and generate an awesome web server
with it.  So, you create a project, manage its dependencies, code the server,
parse the request, learn how to use the API, call the API and deploy it
somewhere.  And that's fine... until you find yourself again in the same
situation with another awesome command.

The awesomeness of UNIX® commands is infinite, so you'll be in this situation
an infinite number of times!  Instead, let's put Kapow! into action.

With Kapow!, when someone asks you for an specific backup (remember your
familiar command?) you just need to create a ``.pow`` file named ``backups.pow``
that contains:

.. code-block:: bash

    kapow route add /backups \
        -c 'cloudx storage ls /backups | grep $(request /params/query) | response /body'

and execute it in the host with the command:

.. code-block:: bash

    kapow server backups.pow

and that's it.  Done.  Do you like it? yes?  Then let's start learning a little
more.


The mandatory Hello World (for WWW fans)
----------------------------------------

First you must create a pow file named ``hello.pow`` with the following contents:

.. code-block:: bash

    kapow route add /greet -c "echo 'hello world' | response /body"

then, you must execute:

.. code-block:: bash

    kapow server hello.pow

and you can check that it works as intended with good ole' ``curl``:

.. code-block:: bash

    curl localhost:8080/greet


The mandatory Echo (for UNIX fans)
----------------------------------

First you must create a pow file named ``echo.pow`` with the following contents:

.. code-block:: bash

    kapow route add -X POST /echo -c 'request /body | response /body'

then, you must execute:

.. code-block:: bash

    kapow server echo.pow

and you can check that it works as intended with good ole ``curl``:

.. code-block:: bash

    curl -X POST -d '1,2,3... testing' localhost:8080/echo


The multiline fun
-----------------

Unless you're a hardcore Perl hacker, you'll probably need to write your stuff
over more than one line.

Don't worry, we need to write several lines, too. Bash, in its magnificent
UNIX® style, provides us with the
`here-documents`_ mechanism that we can leverage precisely for this purpose.

.. _here-documents: https://www.gnu.org/software/bash/manual/bash.html#Here-Documents

Let's write a ``multiline.pow`` file with the following content:

.. code-block:: bash

    kapow route add /log_and_love - <<- 'EOF'
        echo "[$(date)] and stuff" >> stuff.log
        echo love | response /body
    EOF

and then we serve it with ``kapow``:

.. code-block:: bash

    kapow server multiline.pow

Yup.  As simple as that.


Sample Docker usage
===================

Clone the project
-----------------

.. code-block:: bash

    # clone this project


Build the kapow! docker image
-----------------------------

.. code-block:: bash

    cd /path/to/kapow/poc; docker build -t bbva/kapow:0.1 .

Build a docker image for running the nmap example

.. code-block:: bash

    cd /path/to/kapow/poc/examples/nmap; docker build -t kapow-nmap .

Run kapow
---------
.. code-block:: bash

    docker run \
            -it \
            -p 8080:8080 \
            kapow-nmap

which will output something like this:

.. code-block:: bash

    ======== Running on http://0.0.0.0:8080 ========
    (Press CTRL+C to quit)
    Route created POST /list/{ip}
    ROUTE_8ed01c48_bf23_455a_8186_a1df7ab09e48
    bash-4.4#


Test /list endpoint
-------------------
In another terminal, try running:

.. code-block:: bash

    curl http://localhost:8080/list/github.com

which will respond something like:

.. code-block:: bash

    Starting Nmap 7.70 ( https://nmap.org ) at 2019-05-10 14:01 UTC
    Nmap scan report for github.com (140.82.118.3)
    rDNS record for 140.82.118.3: lb-140-82-118-3-ams.github.com
    Nmap done: 1 IP address (0 hosts up) scanned in 0.04 seconds

et voilà !

License
=======

This project is distributed under the [Apache License 2.0](/LICENSE).
