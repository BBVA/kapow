Shell Tricks
============

How to Execute Two Processes in Parallel
----------------------------------------

We want to :command:`ping` two machines parallel.  *Kapow!* can get IP addresses
from query params:

.. code-block:: console
   :linenos:

   $ cat parallel.pow
   kapow route add '/parallel/{ip1}/{ip2}' - <<-'EOF'
   	ping -c 1 -- "$(kapow get /request/matches/ip1)" | kapow set /response/body &
   	ping -c 1 -- "$(kapow get /request/matches/ip2)" | kapow set /response/body &
   	wait
   EOF

Calling with :program:`curl`:

.. code-block:: console
   :linenos:

    $ curl -v http://localhost:8080/parallel/10.0.0.1/10.10.10.1

