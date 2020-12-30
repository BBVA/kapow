Securing the server
===================

**Senior**

  Hi...  I hope you rested last night!

  Come on, I need your help here!

**Junior**

  Good morning! What's the matter? Sounds worrying

**Senior**

  We forgot to take the most basic security measures when deploying our services.
  Every body at the company can access the services and the information is
  transferred in clear text.

**Junior**

  Oh! Damn, you're right! You think we can do anything to solve this mess?

**Senior**

  Yes, I'm pretty sure that those smart guys have thought on that when building
  Kapow! Have a look at the :ref:`documentation <https_mtls>`.

**Junior**

  Got it! They did it, here're the instictions to start a server with HTTPS support.

  It's amazing! It says we can even use mTLS to control access, really promising.

**Senior**

  Ok, ok... First thigs first. We need to get a server certificate to start
  working with HTTPS. Fortunately we can ask for one to the CA we use for the
  other servers. Let's pick up one for development, they're quick to get.

**Junior**

  Yeah! I'll change the startup script to configure HTTPS:

  .. code-block:: console

     $ kapow server --keyfile /etc/kapow/tls/keyfile    \
                    --certfile /etc/kapow/tls/certfile  \
                    /etc/kapow/awesome-route

  It's easy, please copy the private key file and certificate chain to `/etc/kapow/tls` and we can restart.

**Senior**

  Great! it's working, communications are secured. Let's say everybody to change
  from http to https.

**Junior**

  Ok, did it. What are the steps to follow to limit access by using mTLS?

**Senior**

  Besides configuring the server we need to provide the users with their own
  client certificates and private keys so they can configure their browsers and
  the application server.

**Junior**

  Yes, please give me the CA certificate that will issue our client certificates
  and I'll change the startup script again

  .. code-block:: console

     $ kapow server --keyfile /etc/kapow/tls/keyfile            \
                    --certfile /etc/kapow/tls/certfile          \
                    --clientauth true                           \
                    --clientcafile /etc/kapow/tls/clientCAfile  \
                    /etc/kapow/awesome-route

  Done!

**Senior**

  Ok, let's communicate the changes to all the affected teams before we restart

**Junior**

  Oh God, After all we're starting to look like Google

  (chuckles)
