Working with Forms
==================

When a browser submits a form to a server all the values included in the
form are sent to the server in an HTTP call.

Kapow! handles the form decoding for you, the only thing you need to
know is the **name** of the *field* or *fields*.

In this example we respond back with the content of the form field
``myfield``:

.. code-block:: bash

   kapow get /request/form/myfield | kapow set /response/body

