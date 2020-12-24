Working with Init Scripts
=========================

Starting *Kapow!* using an init script
--------------------------------------

An init program, which can be just a shell script, allows you to make calls to
the ``kapow route`` command.

.. code-block:: console
   :linenos:

   $ kapow server example-init-program

With the :file:`example-init-program`:

.. code-block:: console
   :linenos:

   $ cat example-init-program
   #!/usr/bin/env sh
   #
   # This is a simple example of an init program
   #
   echo '[*] Starting my init program'

   # We add 2 Kapow! routes
   kapow route add /my/route -c 'echo hello world | kapow set /response/body'
   kapow route add -X POST /echo -c 'kapow get /request/body | kapow set /response/body'

.. note::

   *Kapow!* can be fully configured using just init scripts


Writing Multiline Routes
------------------------

If you need to write more complex actions, you can leverage multiline routes:

.. code-block:: console
   :linenos:

   $ cat multiline-route
   #!/usr/bin/env sh
   kapow route add /log_and_stuff - <<-'EOF'
   	echo this is a quite long sentence and other stuff | tee log.txt | kapow set /response/body
   	cat log.txt | kapow set /response/body
   EOF

.. warning::

    Be aware of the **"-"** at the end of the ``kapow route add`` command.
    It tells ``kapow route add`` to read commands from `stdin`.

.. warning::

    If you want to learn more about multiline usage, see: `Here Doc
    <https://en.wikipedia.org/wiki/Here_document>`_


Keeping Things Tidy
-------------------

Sometimes things grow, and keeping things tidy is the only way to mantain the
whole thing.

You can distribute your endpoints in several init programs.  And you can keep
the whole thing documented in one html file, served with *Kapow!*.

.. code-block:: console
    :linenos:

    $ cat index-route
    #!/usr/bin/env sh
    kapow route add / - <<-'EOF'
    	cat howto.html | kapow set /response/body
    EOF

    source ./info_stuff
    source ./other_endpoints

You can import other shell script libraries with `source`.


Debugging Init Programs/Scripts
-------------------------------

Since *Kapow!* redirects the standard output and the standard error of the init
program given on server startup to its own, you can leverage ``set -x`` to see
the commands that are being executed, and use that for debugging.

To support debugging user request executions, the server subcommand has a
``--debug`` option flag that prompts *Kapow!* to redirect both the script's
standard output and standard error to *Kapow!*'s standard output, so you can
leverage ``set -x`` the same way as with init programs.


.. code-block:: console

    $ cat withdebug-route
    #!/usr/bin/env sh
    kapow route add / - <<-'EOF'
        set -x
        echo "This will be seen in the log"
    	echo "Hi HTTP" | kapow set /response/body
    EOF

    $ kapow server --debug withdebug-route
