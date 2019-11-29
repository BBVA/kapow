Quick start
===========

We'll explain you a simple example to help you understand what Kapow! really does and why it is awesome.

Scenario
--------

In this example we'll consider that our scenario is a corporate network like this:

.. _quickstart_image:
.. image:: /_static/network.png
   :align: center
   :width: 80%

Our organization has an external host as a bridge between our intranet an the public Internet.

**Our goal: our team need to check if the the host :samp:`Internal Machine` is alive.**

Limitations and constrains
--------------------------

1. We **don't want** to **grant access** to the :samp:`External Host` to anybody.
2. We **don't want** to manage VPNs or any similar solutions to access to *Internal Host* from the Internet.
3. We **want to limit the actions** that an user can perform in our intranet when while it is checking if :samp:`Internal Host` is alive.
4. We **want** the most standard way mechanism. Easy to use and that facilitates the automation.
5. We **don't have budget** to invest in a custom solution.

Study options
-------------

Alter analyze the problem and our goal we conclude that is enough **with a simple :samp:`ping` to samp:`Internal Host`.**

So, then **we need analyze how to perform the ping.**

Accessing via SSH
+++++++++++++++++

In this case we need to create a system user in samp:`External Host` for each user that needs to check if :samp:`Internal host` is alive and we also need to grant access to each user through :samp:`SSH` to the system.

Conclusion: **Not good idea**

Reasons:

    1. We need to manage users (violates our constrains)
    2. We need to access users to system (violates our constrains)
    3. We can't control the :samp:`ping` options the user choice to ping :samp:`Internal Host` (violates our constrains)

Develop custom solution
+++++++++++++++++++++++

Oks, this approach could maybe be the more customizable for our organization but:

1. We'll need to start a new project. Develop it, test it, manage it and maintain it.
2. We need time for the development.
3. We need money. Even we have developers in our organization, their time it's not free.

Conclusion: **Not good idea**

Reasons:

    1. Need to spend money (violates our constrains)
    2. Need to spend time (violates our constrains)

Using Kapow! (Spoiler: the winner!)
+++++++++++++++++++++++++++++++++++

Oks, lets analyze Kapow! and check our constrains:

1. Kapow! is Open Source. Them: **it's free**.
2. By using kapow! we don't need to program our own solution. Them: **don't waste time**.
3. By using Kapow! we can run any command in the :samp:`External Host` limiting the command parameters. Them: **it's safe**.
4. By using Kapow! we can launch any system command as HTTP API easily. Them: **we don't need to grant login access to anybody to :samp:`External Host`**

Conclusion: **Kapow! is the best choice**.

Reasons: it cover all of our requirements.

Using Kapow!
------------

Following the example of the :ref:`Scenario <quickstart_image>` we'll follow these steps:

Install Kapow!
++++++++++++++

Follow :doc:`Install Kapow! <install_and_configure>`.

Write ping.pow file
+++++++++++++++++++

Kapow! use plain text files to define the rules to expose the system command. For our example we need a file like that:

.. code-block:: console

    $ cat ping.pow
    kapow route add /ping -c 'ping -c 1 10.10.10.100 | kapow set /response/body'

Explanation:

1. :samp:`kapow route add /ping` - adds a new HTTP API end-point at :samp:`/ping`.
2. :samp:`-c` - after this parameter we write the system command that Kapow! will runs for each HTTP Request to :samp:`/ping`.
3. :samp:`ping -c 1 10.10.10.100` - sends 1 ping package to the host *10.10.10.100*, i.e. :samp:`Internal Host`.
4. :samp:`| kapow set /response/body` - sends the ping response to be the HTTP Response of HTTP End-point of :samp:`/ping`.

Launch the service
++++++++++++++++++

At this point we only need to launch kapow! with :samp:`simple.pow`:

.. code-block:: console

    $ kapow server ping.pow

Consume the service
+++++++++++++++++++

Then we can call HTTP Service as any usual tool for the web. In this example we'll use :samp:`curl`:

.. code-block:: console

    $ curl http://external.host/ping
    PING 10.10.100 (10.10.100): 56 data bytes
    64 bytes from 10.10.100: icmp_seq=0 ttl=55 time=1.425 ms

Under the hoods
++++++++++++++++

To understand what's happening in the hoods with Kapow! lets see the picture:

.. image:: /_static/sequence.png
   :align: center
   :width: 80%

As you can see, Kapow! perform the *magic* between system commands and HTTP API.
