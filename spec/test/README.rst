Kapow Spec Test Suite
=====================

This is a generic test suite to be run against ANY Kapow!
implementation.


Prerequisites
-------------

First of all you need a working ``kapow server`` installed to run the
tests against it:

Check that your installation is correct by typing this in a shell:

.. code-block:: bash

   kapow --help

You should see a help screen. If not fix it before continuing.

You also need ``gherkin-lint``. Run:

.. code-block:: bash

   npm install -g gherkin-lint


How to run the test suite?
==========================

With the above requisites in place and with your command shell in this
very directory run:

.. code-block:: bash

   make

This should test the current installed kapow implementation and output
something like:

.. code-block:: plain

   13 features passed, 0 failed, 0 skipped
   23 scenarios passed, 0 failed, 0 skipped
   99 steps passed, 0 failed, 0 skipped, 0 undefined
   Took 0m23.553s


Troubleshooting
---------------

Environment customization
~~~~~~~~~~~~~~~~~~~~~~~~~

You can customize some of the test behavior with the following
environment variables:

* ``KAPOW_SERVER_CMD``: The full command line to start a non-interactive
   listening kapow server. By default: ``kapow server``
* ``KAPOW_CONTROLAPI_URL``: URL of the Control API. By default: ``http://localhost:8081``
* ``KAPOW_DATAAPI_URL``: URL of the Data API. By default: ``http://localhost:8080``


Fixing tests one by one
~~~~~~~~~~~~~~~~~~~~~~~

If you use ``make fix`` instead of ``make`` the first failing test will stop
the execution, giving you a chance of inspecting the result. Also the
tests will be in DEBUG mode displaying all requests made to the server.


How to develop new tests?
=========================

Developing new steps
--------------------

First of all execute ``make catalog`` this will display all the existing
steps and associates features. If none of the steps fits your needs
write a new one in your feature and run ``make catalog`` again.
The detected new step will trigger the output of a new section with a
template for the step definition to be implemented (you can copy and
paste it removing the ``u'`` unicode symbol of the strings).


Developing new step definitions
-------------------------------

To make you life easier, mark the feature or scenario you are working on
with the tag ``@wip`` and use ``make wip`` to run only your
scenario/feature.

1. Paste the step definition template you just copied at the end of the
   file ``steps/steps.py``.
2. Run ``make wip`` to test that the step is triggered. You should see a
   ``NotImplementedError`` exception.
3. Implement your step.
4. Run ``make wip`` again.

When you finish implementing your step definitions remove the ``@wip`` of
your feature/scenario and run ``make`` to test everything together.
