.. image:: https://trello-attachments.s3.amazonaws.com/5c824318411d973812cbef67/5ca1af818bc9b53e31696de3/f51eb40412bf09c8c800511d7bbe5634/kapow-1601675_480.png
    :alt: Kapow!

**Kapow!** If you can script it, you can HTTP it.

+-----------------+------------------------------------------------+
| Project site    | https://github.com/BBVA/kapow                  |
+-----------------+------------------------------------------------+
| Issues          | https://github.com/BBVA/kapow/issues           |
+-----------------+------------------------------------------------+
| Specification   | https://github.com/BBVA/kapow/tree/master/spec |
+-----------------+------------------------------------------------+
| Documentation   | https://github.com/BBVA/kapow/tree/master/doc  |
+-----------------+------------------------------------------------+
| Python versions | 3.7 or above                                   |
+-----------------+------------------------------------------------+


CAVEAT EMPTOR
=============

**Warning!!! Kapow!** is under **heavy development** and `specification </spec>`_;
the provided code is a Proof of Concept and the final version will not even
share programming language.  Ye be warned.


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

.. image:: https://github.com/BBVA/kapow/blob/master/resources/kapow.gif?raw=true
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
--------------------------

* Easy command + Hard API = Kapow! to the rescue
* SSH for one command?  Kapow! allows you to share only that command
* Remote instrumentation of several machines?  Make it easy with Kapow!


The more you know
=================

If you want to know more, please follow our `documentation </doc>`_.
