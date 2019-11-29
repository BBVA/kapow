Quick Start Guide
=================

We'll explain a simple example to help you understand what **Kapow!** can do and
why it is so awesome :-).


Scenario
--------

Let's suppose that we're in a corporate network like the following one:

.. _quickstart_image:
.. image:: /_static/network.png
   :align: center
   :width: 80%

Our organization has an external host that act as a bridge between our intranet
an the public Internet.

**Our goal: We need to check if the** :samp:`Internal Host` **is alive.**


Limitations and Constraints
---------------------------

1. We **don't want** to **grant access** to the :samp:`External Host` to anybody.
2. We **don't want** to manage VPNs or any similar solutions to access :samp:`Internal Host` from the Internet.
3. We **want to limit the actions** that a user can perform in our intranet while it is checking if :samp:`Internal Host` is alive.
4. We **want** to use the most standard mechanism.  Easy to use and automate.
5. We **don't have a budget** to invest in a custom solution.


The Desired Solution
--------------------

After analyzing the problem and our with our goal in mind, we conclude that it
is enough **to use a simple** :samp:`ping` **to** :samp:`Internal Host`.

So, the next step is to **analyze how to perform the ping.**


Accessing via SSH to :samp:`External Host`
++++++++++++++++++++++++++++++++++++++++++

If we choose this option then we need to create a user in the host and grant
them access via :samp:`SSH` to :samp:`External Host` for every person that needs
to check the :samp:status of `Internal host`.

Conclusion: **Not a good idea.**

Reasons:

  1. We need to manage users (violates a constraint.)
  2. We need to grant usesrs access to a host (violates a constraint.)
  3. We can't control what :samp:`ping` options the user can use to ping :samp:`Internal Host` (violates a constraint.)


Develop and Deploy a Custom Solution
++++++++++++++++++++++++++++++++++++

Ok, this approach could be the best choice for our organization, but:

1. We'll need to create a new project, develop, test, manage and maintain it.
2. We need to wait for for the development to be production ready.
3. We need a budget, even if we have developers in our organization.

Conclusion: **Not a good idea.**

Reasons:

1. Need to spend money (violates a constraint.)
2. Need to spend time (and time is money, see reason #1)


Using Kapow! (spoiler: it's the winner!)
++++++++++++++++++++++++++++++++++++++++

Ok, let's analyze **Kapow!** and check if it is compatible with our constraints:

1. **Kapow!** is Open Source, so **it's also free as in beer**.
2. By using kapow! we don't need to program our own solution, so we **don't have to waste time**.
3. By using Kapow! we can run any command in the :samp:`External Host` limiting the command parameters, so **it's safe**.
4. By using Kapow! we can launch any system command as an HTTP API easily, so **we don't need to grant login access to anybody to** :samp:`External Host`.

Conclusion: **Kapow! is the best choice.**

Reasons: It satisfies all of our requirements.


Using Kapow!
------------

In order to get our example :ref:`Scenario <quickstart_image>` working we need to follow the below steps.


Install Kapow!
++++++++++++++

Follow :doc:`Install Kapow! <install_and_configure>` instructions.


Write a `ping.pow` File
+++++++++++++++++++++++

Kapow! uses plain text files (called ``POW`` files) so you can define the
endpoints you want to expose the system command with.  For our example we need a
file like this:

.. code-block:: console

    $ cat ping.pow
    kapow route add /ping -c 'ping -c 1 10.10.10.100 | kapow set /response/body'

Explanation:

1. :samp:`kapow route add /ping` - adds a new HTTP API endpoint at :samp:`/ping` path in the Kapow! server.  You have to use `GET` method to invoke the endpoint.
2. :samp:`-c` - after this parameter we write the system command that **Kapow!** will run each time the endpoint is invoked.
3. :samp:`ping -c 1 10.10.10.100` - sends 1 ping package to the host *10.10.10.100*, i.e. :samp:`Internal Host`.
4. :samp:`| kapow set /response/body` - writes the output of `ping` to the body of the response, so you can see it.


Launch the Service
++++++++++++++++++

At this point we only need to launch `kapow` with our :samp:`ping.pow`:

.. code-block:: console

    $ kapow server ping.pow


Consume the Service
+++++++++++++++++++

Now we can call our newly created endpoint by using our favorite HTTP client.
In this example we're using :samp:`curl`:

.. code-block:: console

    $ curl http://external.host/ping
    PING 10.10.100 (10.10.100): 56 data bytes
    64 bytes from 10.10.100: icmp_seq=0 ttl=55 time=1.425 ms

et voilà !


Under the Hood
++++++++++++++

To understand what's happening under the hood with **Kapow!** let's see the
following diagram:

.. image:: /_static/sequence.png
   :align: center
   :width: 80%

As you can see, **Kapow!** provides the necessary *magic* to turn a **system
command** into an ``HTTP API``.
