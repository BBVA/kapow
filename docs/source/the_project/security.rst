Security Concerns
=================

Special care has to be taken when using parameters provided by the user when
composing command line invocations.

Sanitizing user input is not a new problem, but in the case of *Kapow!*, we
have to take into account also the way that the shell parses its arguments,
as well as the way the command itself interprets them, in order to get it right.

.. warning::

   It is **imperative** that the user input is sanitized properly if we are
   going feed it as a parameter to a command line program.


Parameter Injection Attacks
---------------------------

When you resolve variable values be careful to tokenize correctly by using
double quotes.  Otherwise you could be vulnerable to **parameter injection
attacks**.

**This example is VULNERABLE to parameter injection**

In this example, an attacker can inject arbitrary parameters to :command:`ls`.

.. code-block:: console
   :linenos:

   $ cat command-injection.pow
   kapow route add '/vulnerable/{value}' - <<-'EOF'
   	ls $(kapow get /request/matches/value) | kapow set /response/body
   EOF

Exploiting using :program:`curl`:

.. code-block:: console
   :linenos:

   $ curl http://localhost:8080/vulnerable/-lai%20hello

**This example is NOT VULNERABLE to parameter injection**

Note how we add double quotes when we recover *value* data from the
request:

.. code-block:: console
   :linenos:

   $ cat command-injection.pow
   kapow route add '/not-vulnerable/{value}' - <<-'EOF'
   	ls -- "$(kapow get /request/matches/value)" | kapow set /response/body
   EOF


.. warning::

   Quotes around parameters only protect against the injection of additional
   arguments, but not against turning a non-option into option or vice-versa.
   Note that for many commands we can leverage double-dash to signal the end of
   the options.  See the "Security Concern" section on the docs.


Parameter Mangling Attacks
--------------------------

Let's consider the following route:

.. code-block:: bash

   #!/bin/sh
   kapow route add /find -c <<-'EOF'
          BASEPATH=$(kapow get /request/params/path)
          find "$BASEPATH" | kapow set /response/body
   EOF


The expected use for this endpoint is something like this:

.. code-block:: console

   $ curl http://kapow-host/find?path=/tmp
   /tmp
   /tmp/.X0-lock
   /tmp/.Test-unix
   /tmp/.font-unix
   /tmp/.XIM-unix
   /tmp/.ICE-unix
   /tmp/.X11-unix
   /tmp/.X11-unix/X0


.. todo:: Meanwhile, in Russia:

Let's suppose that a malicious attacker gets access to this service and
makes this request:

.. code-block:: console

   $ curl http://kapow-host/find?path=-delete


Let's see what happens:

The command that will eventually be executed by :command:`bash` is:

.. code-block:: bash

   find -delete | kapow set /response/body

This will *silently delete all the files below the current directory*, no
questions asked.  Probably not what you expected.

This happens because :command:`find` has the last word on how to interpret its
arguments.  For :command:`find`, the argument :nref-option:`-delete` is not a path.

Let's see how we can handle this particular case:

.. code-block:: bash

   #!/bin/sh
   kapow route add /find -c <<-'EOF'
           USERINPUT=$(kapow get /request/params/path)
           BASEPATH=$(dirname -- "$USERINPUT")/$(basename -- "$USERINPUT")
           find "$BASEPATH" | kapow set /response/body
   EOF

.. note::

   Since this is critical for keeping your *Kapow!* services secure, we are working
   on a way to make this more transparent and safe, while at the same time keeping
   it *Kapowy*.


