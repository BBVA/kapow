Working with JSON
==================

Nowadays webservices are json based so making your script json aware is probably
a good chioce. In order to be able to extract data from and compose json
documents from a script you can use
`jq <https://https://stedolan.github.io/jq/>`_.

In this example we extract the name field from the incomming json document in
order to generate a two attribute json response.

.. code-block:: bash

  who=$(kapow get /request/body | jq -r .name)
  kapow set /response/headers/Content-Type "application/json"
  kapow set /response/status 200
  jq --arg greet "Hello" --arg value "${who:-World}" -n \{greet:\$greet\,to:\$value\} | kapow set /response/body
