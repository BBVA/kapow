.. image:: https://trello-attachments.s3.amazonaws.com/5c824318411d973812cbef67/5ca1af818bc9b53e31696de3/f51eb40412bf09c8c800511d7bbe5634/kapow-1601675_480.png
    :alt: Kapow!

.. image:: https://circleci.com/gh/BBVA/kapow/tree/master.svg?style=svg
    :target: https://circleci.com/gh/BBVA/kapow/tree/master

.. image:: https://goreportcard.com/badge/github.com/bbva/kapow
    :target: https://goreportcard.com/report/github.com/bbva/kapow

**Kapow!** If you can script it, you can HTTP it.


CAVEAT EMPTOR
=============

**Warning!!! Kapow!** is under **heavy development** and `specification </spec>`_;
the provided code is a Proof of Concept and the final version will not even
share programming language.  Ye be warned.


What is it?
===========

Kapow! is an adapter between the world of Pure UNIX® Shell and a HTTP service.

Some tasks are more convenient in the shell, like cloud interactions, or some
administrative tools.  On the other hand, some tasks are more convenient as a
service, like DevSecOps tooling.

Kapow! lies between these two worlds, making your life easier.  Maybe you wonder
about how this kind of magic can happen; if you want to know the nitty-gritty
details, just read our `specification </spec>`_;.  Or, if you want to know how
Kapow! can help you first, let's start with a common situation.

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

With Kapow! you just need to create a ``.pow`` file named ``backups.pow`` that
contains:

.. code-block:: sh

    kapow route add /backups \
        -c 'cloudx storage ls /backups | grep "$(kapow get /request/params/query)" | kapow set /response/body'

and execute it in the remote host with the command:

.. code-block:: sh

    kapow server backups.pow

and that's it.  Done.  You have a web server that people can use to request
their backups every time they need only by invoking the URL
`http://remotehost/backups?query=project`

Do you like it? yes?  Then let's start learning a little more, you can access
the `documentation </doc>`_; section to find installation instructions and some
examples.







How it was born
---------------

Some awesome history is coming.


Kapow! for the impatient
========================

When you need to **share** a ``command`` but **not** a complete remote ``ssh
access``, Kapow!  will help you with the power of HTTP:

.. image:: https://trello-attachments.s3.amazonaws.com/5c824318411d973812cbef67/5ca1af818bc9b53e31696de3/784a183fba3f24872dd97ee28e765922/Kapow!.png
    :alt: Where Kapow! lives

Kapow! allows you to write a litte script that will **serve an executable as REST
service**.  This script will let you define how to connect HTTP and the  Shell
using Kapow!'s shell abstractions to the HTTP world. See it to believe:

.. image:: resources/kapow.gif?raw=true
    :alt: Kapow! in action


Superpowers
-----------

Kapow! gives you:

* A very simple way to turn any shell **executable into an API**
* A **remote administration** API
* A way to define the integration in you own terms, obligations-free!


Curses
------

Kapow! can't help when:
-----------------------

* You need high throughput: Kapow! spawns a new executable for every HTTP call
* You must perform complex logic to attend the request: never use Kapow! if
  your executables don't perform al least 90% of the hard work
* You are building a huge application


When it is your best friend:
----------------------------

* Easy command + Hard API = Kapow! to the rescue
* SSH for one command?  Kapow! allows you to share only that command
* Remote instrumentation of several machines?  Make it easy with Kapow!


The more you know
=================

If you want to know more, please follow our `documentation </doc>`_.
