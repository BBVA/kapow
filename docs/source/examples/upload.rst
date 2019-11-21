Upload a file
=============

HTTP request allows us to send and receive files by using the Multipart standard.

Kapow! allow us to handle files received in the request. In this example we
respond back with the line count of the file received in the request.

.. code-block:: bash

  fname=$(kapow get /request/files/myfile/filename)
  lcount=$(kapow get /request/files/myfile/content | wc -l)
  kapow set /response/status 200
  echo "$fname has $lcount lines" | kapow set /response/body

You can try this by using the following curl:

.. code-block:: bash

  curl -F "myfile=@README.rst" http://localhost:8080/linecount
