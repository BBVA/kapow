Installing Kapow!
=================

Kapow! has a reference implementation in ``Go`` that is under active development
right now.  If you want to start using **Kapow!** you can choose from several
options.


Download and Install a Binary
-----------------------------

Binaries for several platforms are available from the
`releases <https://github.com/BBVA/kapow/releases>`_ section, visit the latest
release page and download the binary corresponding to the platfom and
architecture you want to install **Kapow!** in.


Linux
^^^^^

Install the downloaded binary using the following command as a privileged user.

.. code-block:: bash

  $ install -t /usr/local/bin/kapow path_to_downloaded_binary


Windows
^^^^^^^

Copy the downloaded binary to a directory of your choice and update the system
``PATH`` variable to include that directory.


Install the package with go get
-------------------------------

If you already have `installed and configured <https://golang.org/cmd/go/>`_
the ``go`` runtime in the host where you want to run **Kapow!**, simply run:

.. code-block:: bash

  $ go get -u github.com/BBVA/kapow


Include Kapow! in your Container Image
--------------------------------------

If you want to include **Kapow!** in a ``Docker`` image you can add the binary
directly from the releases section. Following is an example ``Dockerfile`` that
includes **Kapow!**.

.. code-block:: dockerfile

  FROM debian:stretch-slim

  RUN apt-get update

  ADD https://github.com/BBVA/kapow/releases/download/<VERSION>/kapow_linux_amd64 /usr/bin/kapow

  RUN chmod 755 /usr/bin/kapow

  ENTRYPOINT ["/usr/bin/kapow"]

If the container is intended for running the server and you want to dinamically
configure it, remember to include a `--control-bind` param with an external bind
address (i.e. ``0.0.0.0``) and to map all the needed ports in order to get access
to the control interface.

After building the image you can run the container with:

.. code-block:: bash

  $ docker run --rm -i -p 8080:8080 -v $(pwd)/whatever.pow:/opt/whatever.pow kapow:latest server /opt/whatever.pow

With the `-v` parameter we map a local file into the container's filesystem so
we can use it to configure our **Kapow!** server on startup.
