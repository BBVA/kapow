Working with pow Files
======================

Starting *Kapow!* using a pow file
----------------------------------

A :file:`pow` file is just a :command:`bash` script, where you make calls to the
``kapow route`` command.

.. code-block:: console
   :linenos:

   $ kapow server example.pow

With the :file:`example.pow`:

.. code-block:: console
   :linenos:

   $ cat example.pow
   #
   # This is a simple example of a pow file
   #
   echo '[*] Starting my script'

   # We add 2 Kapow! routes
   kapow route add /my/route -c 'echo hello world | kapow set /response/body'
   kapow route add -X POST /echo -c 'kapow get /request/body | kapow set /response/body'

.. note::

   *Kapow!* can be fully configured using just :file:`pow` files


Load More Than One pow File
---------------------------

You can load more than one :file:`pow` file at time.  This can help you keep
your :file:`pow` files tidy.

.. code-block:: console
   :linenos:

   $ ls pow-files/
   example-1.pow   example-2.pow
   $ kapow server <(cat pow-files/*.pow)


Writing Multiline pow Files
---------------------------

If you need to write more complex actions, you can leverage multiline commands:

.. code-block:: console
   :linenos:

   $ cat multiline.pow
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


Keep things tidy
----------------

Sometimes things grow, and keep things tidy it's the only way to mantain the
hole thing.

You can distribute your endpoints among several pow files. And you can keep the
hole thing documented in one html server with kapow.

.. code-block:: console
    :linenos:

    $ cat index.pow
    #!/bin/bash

    kapow route add / - <<-'EOF'
        cat howto.html | kapow set /response/body
    EOF

    source ./info_stuff.pow
    source ./other_endpoints.pow

As you can see the pow files can be imported into another pow files using
source. In fat a pow file it's a regular shell file.

