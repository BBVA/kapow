*Kapow!* Behind a Reverse Proxy
===============================

In this section we present a series of reverse proxy configurations that
augment the capabilities of *Kapow!*.

.. note::

   In this section we refer to the host running the *Kapow!* server as
   `kapow:8080`.


Serving over HTTPS
------------------

*Kapow!* currently does not support `HTTPS` but you can use a
reverse proxy to serve a *Kapow!* service via `HTTPS`.

For testing purposes you can generate a self-signed certificate with the
following command:

.. code-block:: console

   $ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes


Caddy
+++++

* **Automatic Let's Encrypt Certificate**

  `Caddy` automatically enables `HTTPS` using `Let's Encrypt`
  certificates given that `some criteria are met`_.

  .. code-block:: none

     yourpublicdomain.example
     proxy / kapow:8080

* **Automatic Self-signed Certificate**

  If you want `Caddy` to automatically generate a self-signed
  certificate for testing you can use the following configuration.

  .. code-block:: none

     yourdomain.example
     proxy / kapow:8080
     tls self_signed

* **Custom Certificate**

  If you already have a valid certificate for your server use this
  configuration.

  .. code-block:: none

     yourdomain.example
     proxy / kapow:8080
     tls /path/to/cert.pem /path/to/key.pem


HAProxy
+++++++

With the following configuration you can run `HAProxy` with a custom
certificate.

.. code-block:: none

   frontend myserver.local
       bind *:443 ssl crt /path/to/myserver.local.pem
       mode http
       default_backend nodes

   backend nodes
       mode http
       server kapow1 kapow:8080


.. note::

   You can produce ``myserver.local.pem`` from the certificates in
   previous examples with this command:

   .. code-block:: console

      $ cat /path/to/cert.pem /path/to/key.pem > /path/to/myserver.local.pem


nginx
+++++

With the following configuration you can run `nginx` with a custom
certificate.

.. code-block:: none

   server {
    listen              443 ssl;
    server_name         myserver.local;
    ssl_certificate     /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://kapow:8080;
    }
   }

.. _some criteria are met: https://caddyserver.com/v1/docs/automatic-https
