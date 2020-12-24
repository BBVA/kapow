Using HTTPS and mTLS with *Kapow!*
==================================

*Kapow!* can be accesed over HTTPS and use mutual TLS for authentication. Right
now there are two possibilities to configure HTTPS/mTLS on a running server.

.. note::

   In the following sections we refer to the host running the *Kapow!* server as
   `kapow:8080`.

For testing purposes you can generate a self-signed certificate with the
following command:

.. code-block:: console

  $ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes


Using *Kapow!* built in capabilities
------------------------------------

In this section we present the option flags that the *Kapow!* server command
provides in order to set up a server with HTTPS and/or mTLS.

Enabling HTTPS
++++++++++++++

When starting the server we can provide the private key and the certificate
chain presented to the clients by giving to the server the paths to the
corresponding files using the `--keyfile` and `--certfile` option flags in the
command line:

.. code-block:: console

  $ kapow server --keyfile path/to/keyfile --certfile path/to/certfile foobar-route

Now *Kapow!* is listening on its default port (8080) accepting requests over
HTTPS. You can test it with the following command:

.. code-block:: console

    $ curl --cacert path/to/CAfile https://localhost:8080/endpoint

Where `path/to/CAfile` is the path to the file containing the CA certificate
that issued *Kapow!*'s certificate.


Enabling mTLS
++++++++++++++

Once we have *Kapow!* configured to use HTTPS we can, optionally, activate mTLS
so we can reject client connections that do not present a valid client certificate.
Currently only issuer CA matching is supported, but in a near future we'll be able
to filter by DN in order to get more fine grained authentication.

In order to activate mTLS we have to provide *Kapow!* server command with the
CA certificate issuing the client certificates we want to accept with the
`--clientcafile` option flag and toggle the `--clientauth=true` option flag:

.. code-block:: console

  $ kapow server --keyfile path/to/keyfile --certfile path/to/certfile --clientauth=true --clientcafile path/to/clientCAfile foobar-route

With this configuration *Kapow!* will reject connections that do not present a
client certificate or one certificate not issued by the specified CA. You can
test it with the following command:

.. code-block:: console

    $ curl --cacert path/to/CAfile --cert path/to/clientcredentials https://localhost:8080/endpoint

Where `path/to/clientcredentials` is the path to the file in PKCS#12 format
containing the client certificate and private key. If you have the certificate
and the private key in different files you can use the --key option to specify
it independently.


*Kapow!* Behind a Reverse Proxy
-------------------------------

Although *Kapow!* supports `HTTPS` you can use a reverse proxy to serve a
*Kapow!* service via `HTTPS`.

In this section we present a series of reverse proxy configurations that
augment the capabilities of *Kapow!*.


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

In order to enable mutual TLS authentication read the `Caddy documentation`_.


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

In order to enable mutual TLS authentication read the `HAProxy documentation`_.


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

   In order to enable mutual TLS authentication read the `Nginx documentation`_.

.. _some criteria are met: https://caddyserver.com/v1/docs/automatic-https
.. _Caddy documentation: https://caddyserver.com/docs/caddyfile/directives/tls
.. _HAProxy documentation: https://www.haproxy.com/de/documentation/aloha/12-0/traffic-management/lb-layer7/tls/
.. _Nginx documentation: https://smallstep.com/hello-mtls/doc/server/nginx
