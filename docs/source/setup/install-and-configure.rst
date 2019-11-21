Installing Kapow!
=================

Kapow! has a reference implementation in Go that is under active develpment
right now.  If you want to start using Kapow! you can:
- Download a binary (linux, at this moment) from our
`releases <https://github.com/BBVA/kapow/releases>`_ section
- Install the package with the get command (you need the Go runtime installed
and `configured <https://golang.org/cmd/go/>`)

.. code-block:: bash

    go get -u github.com/BBVA/kapow

Using Kapow!
============

Kapow! binary gives you both, a way to start the server and a command line
client to interact with it.

Running the server
------------------

You start a Kapow! server by using the server command ``kapow server``.  It
automatically binds the three HTTP modules when it starts:

- The control server: Used to manage the user defined routes.  It exposes the
control API and listens by default in the loopback interface at port 8081.  You
can change this configurtation by ussing the
--control-bind <listen_address>:<listen_port> parameter.
- The data server: Allows access to the resources tree to the scripts triggered
by user's requests.  It exposes the control API and listens by default in the
loopback interface at port 8082.  You can change this configurtation by ussing
the --data-bind <listen_address>:<listen_port> parameter.
- The ``user server``: This server is the one that makes available to the
outside all the routes configured through the control server.  It listens by
default in the port 8080 of all configured interfaces (0.0.0.0). You can change this
configurtation by ussing the --bind <listen_address>:<listen_port> parameter.



********--------********--------********----------------********--------********--------********

You start a Kapow! server by using the server command ``kapow server``.  You can
configure the listen address of the different modules by using the
corresponding parameters:

- --control-bind <listen_address>:<listen_port>: Allows to manage the user
defined routes.  Defaults to 'localhost:8081' (loopback interface).
- --data-bind <listen_address>:<listen_port>: Allows access to the resources
tree to the scripts triggered by user's requests.  Defailts to 'localhost:8081'
(loopback interface).
- --bind <listen_address>:<listen_port>: Publishes the routes configured
through the control server to the outside world.  Defaults to '0.0.0.0:8080'
(all configured host's interfaces.)


Managing routes
---------------

Kapow!'s route command allows us to manage the routes that we want to publish
to the outside world.  In order to contact with the desired Kapow! server you
can use the ``--control-url`` command line parameter or the KAPOW_CONTROL_URL
environmental variable to set the correct value.

In the same way the ``--data-url`` command line parameter or the KAPOW_DATA_URL
environmental variable will allow you to set to set the server listen address
when accesing the data server, although this case is less frequent.
