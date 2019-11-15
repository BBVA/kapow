Welcome to *Kapow!*
===================

With *Kapow!* you can publish simple **shell scripts** as **HTTP services** easily.

*Kapow!* with an example
------------------------

**Goal**

We want users on the Internet to be able to ``ping`` an *Internal Host*
which is inside a private network.

.. image:: _static/network.png

**Limitations**

- We can't allow users to log into any host. 

- We need to have full control over the precise command is run as
  well as the parameters used.

**Solution**

With a *Kapow!* one-liner you can allow your users to run a command inside
*External Host* through an HTTP call.

.. image:: _static/sequence.png

This is the only line you'll need:

.. code-block:: bash

   $ kapow route add /ping -c 'ping -c1 10.10.10.100 | kapow set /response/body'


.. todo::

   Mention license and contributing



