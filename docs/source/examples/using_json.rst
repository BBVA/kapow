Using JSON
==========

Modify JSON by Using Shell Commands
-----------------------------------

.. note::

    Nowadays Web services are `JSON`-based, so making your script `JSON` aware is
    probably a good choice.  In order to be able to extract data from a `JSON`
    document as well as composing `JSON` documents from a script, you can leverage
    `jq <https://stedolan.github.io/jq/>`_.


Example #1
++++++++++

In this example our *Kapow!* service will receive a `JSON` value with an incorrect
date, then our ``pow`` file will fix it and return the correct value to the user.

.. code-block:: console
   :linenos:

   $ cat fix_date.pow
   kapow route add -X POST /fix-date - <<-'EOF'
   	kapow set /response/headers/Content-Type application/json
   	kapow get /request/body | jq --arg newdate "$(date +'%Y-%m-%d_%H-%M-%S')" '.incorrectDate=$newdate' | kapow set /response/body
   EOF

Call the service with :program:`curl`:

.. code-block:: console
   :linenos:

   $ curl -X POST http://localhost:8080/fix-date -H 'Content-Type: application/json' -d '{"incorrectDate": "no way, Jose"}'
   {
      "incorrectDate": "2019-11-22_10-42-06"
   }


Example #2
++++++++++

In this example we extract the ``name`` field from the incoming `JSON` document in
order to generate a two-attribute `JSON` response.

.. code-block:: console

   $ cat echo-attribute.pow
   kapow route add -X POST /echo-attribute - <<-'EOF'
   	JSON_WHO=$(kapow get /request/body | jq -r .name)

   	kapow set /response/headers/Content-Type application/json
   	kapow set /response/status 200

   	jq --arg greet Hello --arg value "${JSON_WHO:-World}" --null-input '{ greet: $greet, to: $value }' | kapow set /response/body
   EOF

Call the service with :program:`curl`:

.. code-block:: console
   :linenos:
   :emphasize-lines: 4

   $ curl -X POST http://localhost:8080/echo-attribute -H 'Content-Type: application/json' -d '{"name": "MyName"}'
   {
     "greet": "Hello",
     "to": "MyName"
   }
