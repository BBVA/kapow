# Controlling Access with Mutual TLS (mTLS)

## CAVEAT EMPTOR

This notes demonstrate the mTLS functionality of *Kapow!*.  For production
environments, observe proper security best practices regarding the creation,
storage and securing of PKI certificates.

## Prerequisites

1. This should be run in an UNIX®-like OS.
2. Install latest release of [easyRSA](https://github.com/OpenVPN/easy-rsa/releases).
3. `mykapowserver` must resolve to `127.0.0.1`;
you can easily accomplish this by adding:
```
127.0.0.1   mykapowserver
```
to your `/etc/hosts`.

4. Copy `kapow/tools/validsslclient` to somewhere in your `$PATH`.

## Prepare a CA, Server and Client Certificates

Go to an empty directory, and run the following commands.

### Init the PKI system

``` console
$ easyrsa init-pki

init-pki complete; you may now create a CA or requests.
Your newly created PKI dir is: /home/user/mTLS/pki
```

### Build a CA

``` console
$ easyrsa build-ca nopass
[...]
Generating a RSA private key
........................+++++
.................+++++
writing new private key to '/home/user/mTLS/pki/private/ca.key.tUuvKT5D4h'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Common Name (eg: your user, host, or server name) [Easy-RSA CA]:

CA creation complete and you may now import and sign cert requests.
Your new CA certificate file for publishing is at:
/home/user/mTLS/pki/ca.crt
```

### Server Certificates

#### Generate a Request for the Server
``` console
$ easyrsa gen-req mykapowserver nopass
Generating a RSA private key
..........................................................................+++++
.+++++
writing new private key to '/home/user/mTLS/pki/private/mykapowserver.key.0FKwkzgE4X'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Common Name (eg: your user, host, or server name) [mykapowserver]:

Keypair and certificate request completed. Your files are:
req: /home/user/mTLS/pki/reqs/mykapowserver.req
key: /home/user/mTLS/pki/private/mykapowserver.key

```

#### Sign the Server Request
``` console
$ easyrsa sign-req server mykapowserver


You are about to sign the following certificate.
Please check over the details shown below for accuracy. Note that this request
has not been cryptographically verified. Please be sure it came from a trusted
source or that you have verified the request checksum with the sender.

Request subject, to be signed as a server certificate for 3650 days:

subject=
    commonName                = mykapowserver


Type the word 'yes' to continue, or any other input to abort.
  Confirm request details: yes
Using configuration from .../easyrsa-3.0.0/share/easyrsa/openssl-1.0.cnf
Check that the request matches the signature
Signature ok
The Subject's Distinguished Name is as follows
commonName            :ASN.1 12:'mykapowserver'
Certificate is to be certified until Dec 20 10:59:29 2030 GMT (3650 days)

Write out database with 1 new entries
Data Base Updated

Certificate created at: /home/user/mTLS/pki/issued/mykapowserver.crt


```

Optionally Validate that the Server Certificate has the Correct Data:
``` console
$ easyrsa show-req mykapowserver full

Showing req details for 'mykapowserver'.
This file is stored at:
/home/user/mTLS/pki/reqs/mykapowserver.req

Certificate Request:
    Data:
        Version: 1 (0x0)
        Subject:
            commonName                = mykapowserver
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:d1:55:72:19:52:51:8a:1e:6e:29:a3:d6:da:7f:
                    e3:e4:e7:5b:14:b9:59:7f:d9:8c:6a:78:b4:e6:71:
                    20:d5:aa:e9:8c:54:e1:14:09:4e:11:0d:13:6b:0c:
                    04:df:b4:d2:f9:27:b3:17:f6:bc:a0:45:3c:e1:2e:
                    57:32:6a:7b:4e:84:4b:6e:e2:cb:1f:91:b1:e5:67:
                    31:17:79:db:d7:54:d2:72:32:12:2e:a6:52:c7:49:
                    98:fa:73:8e:7c:a4:62:c9:1d:bd:0b:a0:8a:98:2a:
                    9f:19:bf:2c:f7:4a:06:a9:92:f5:99:64:db:6a:21:
                    05:09:c4:04:de:1c:e6:14:98:10:0d:b8:1a:6e:71:
                    ca:e1:85:e6:c5:46:34:09:ff:9f:e3:05:b7:3d:35:
                    22:93:a2:84:eb:e0:cb:42:0c:ef:c4:8d:8f:28:4a:
                    c3:4b:d5:e1:ad:c2:a3:6b:6c:03:a2:1c:9f:7e:70:
                    84:8c:b9:24:99:5e:43:bf:cd:1b:ed:40:20:70:ec:
                    55:46:00:9d:16:9e:d5:c5:e2:d7:40:0a:60:bb:ac:
                    17:ed:c6:4e:f1:c4:62:d1:f7:14:20:21:12:57:c1:
                    c5:ca:3b:58:88:f6:47:93:52:62:30:b1:4e:21:e4:
                    21:6d:a1:c1:a5:0a:6c:da:62:a3:d2:15:18:d7:f8:
                    bb:45
                Exponent: 65537 (0x10001)
        Attributes:
            a0:00
    Signature Algorithm: sha256WithRSAEncryption
         ba:cf:0e:77:a5:2f:01:4c:ac:a6:9e:5a:92:df:0c:fd:e3:37:
         0d:e0:b8:41:3d:44:36:85:31:fa:7e:ac:0f:f9:b2:ec:89:e5:
         7e:cf:92:6f:02:e3:2b:71:d5:b3:ce:97:d8:4c:38:10:a0:ac:
         b4:f6:87:d2:d6:77:24:4d:5f:95:ce:6b:19:3b:09:3e:0b:bc:
         83:a4:f0:d8:2c:6e:b6:aa:53:c7:a7:a6:29:eb:f2:a9:e8:8d:
         18:bd:d1:8f:15:de:fc:01:94:30:df:e6:cd:10:f6:a8:2b:8b:
         42:16:b1:02:e7:1b:b3:0d:81:33:73:94:bc:20:f0:9a:3e:e8:
         26:2e:46:50:ca:ae:1a:ad:30:90:2b:5a:b9:de:6d:f1:bd:53:
         7a:4e:cf:d4:56:7c:74:e0:33:8e:40:b0:72:1d:e4:bc:ac:91:
         ad:7c:3d:6d:8f:09:04:2f:04:16:bc:9a:b6:15:ba:1f:0a:d9:
         6f:2f:e3:a9:c5:34:86:f1:40:b7:a7:04:47:3b:47:9e:f4:a4:
         73:72:1a:df:50:d6:b4:e9:bb:7a:23:94:c6:c8:6b:d6:75:ab:
         f3:46:55:24:7d:a8:bc:7b:08:35:9d:09:0d:75:07:b2:14:3f:
         63:85:31:c1:38:9e:41:0d:fb:b6:dc:48:cb:6c:2e:8c:ab:a7:
         19:cf:37:25
```

### Client Certificate

#### Generate a Request for the Client
``` console
$ easyrsa gen-req authorizedclient nopass
Generating a RSA private key
........+++++
................................................................................+++++
writing new private key to '/home/user/mTLS/pki/private/authorizedclient.key.74zBlwL81Z'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Common Name (eg: your user, host, or server name) [authorizedclient]:

Keypair and certificate request completed. Your files are:
req: /home/user/mTLS/pki/reqs/authorizedclient.req
key: /home/user/mTLS/pki/private/authorizedclient.key


```

#### Sign the Client Request
``` console
$ easyrsa sign-req client authorizedclient


You are about to sign the following certificate.
Please check over the details shown below for accuracy. Note that this request
has not been cryptographically verified. Please be sure it came from a trusted
source or that you have verified the request checksum with the sender.

Request subject, to be signed as a client certificate for 3650 days:

subject=
    commonName                = authorizedclient


Type the word 'yes' to continue, or any other input to abort.
  Confirm request details: yes
Using configuration from .../easyrsa-3.0.0/share/easyrsa/openssl-1.0.cnf
Check that the request matches the signature
Signature ok
The Subject's Distinguished Name is as follows
commonName            :ASN.1 12:'authorizedclient'
Certificate is to be certified until Dec 20 11:04:18 2030 GMT (3650 days)

Write out database with 1 new entries
Data Base Updated

Certificate created at: /home/user/mTLS/pki/issued/authorizedclient.crt



```

Optionally Validate that the Client Certificate has the Correct Data:
``` console
$ easyrsa  show-req authorizedclient full

Showing req details for 'authorizedclient'.
This file is stored at:
/home/user/mTLS/pki/reqs/authorizedclient.req

Certificate Request:
    Data:
        Version: 1 (0x0)
        Subject:
            commonName                = authorizedclient
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:cf:6b:49:ea:f5:f7:6a:9e:00:4c:50:a9:9d:f6:
                    70:32:3c:bc:07:ce:dc:c9:53:01:79:a7:28:63:7b:
                    60:ff:31:e7:b0:03:f1:67:93:af:f8:e3:80:51:d9:
                    d4:9e:b2:01:31:a5:19:bc:e8:f7:92:8d:32:e8:6f:
                    8c:7d:b5:38:43:17:a8:9e:f0:f4:c6:fc:90:c5:b2:
                    1f:87:39:70:d3:03:bb:45:8f:f6:a3:c5:8e:4d:0a:
                    c2:24:a6:23:40:e9:f4:0e:20:7d:c2:34:49:48:92:
                    5a:dc:9c:fa:43:c4:8f:35:c4:77:c3:4c:c5:e7:b5:
                    a8:53:8f:89:51:09:29:ba:82:93:0a:39:02:79:83:
                    19:4b:60:03:d3:fd:26:65:25:1a:5d:80:4f:7f:84:
                    4a:77:13:81:c8:c8:37:ad:bd:0f:bf:9b:62:48:57:
                    ee:1a:2f:e4:00:35:d2:82:23:73:0b:8b:f7:56:3e:
                    58:4d:ed:e7:87:a1:1c:a0:db:0a:3a:bc:a5:d2:a4:
                    93:92:88:2f:24:29:3c:86:2d:ce:72:64:1e:ed:bd:
                    c5:3d:06:da:29:4f:16:36:f8:14:d9:0e:9a:1d:fc:
                    0f:5a:e6:a9:65:06:fd:f1:59:f3:fb:6f:3a:ac:7c:
                    70:59:c8:a2:60:a9:04:c3:2c:57:e8:95:11:ef:e0:
                    99:59
                Exponent: 65537 (0x10001)
        Attributes:
            a0:00
    Signature Algorithm: sha256WithRSAEncryption
         40:53:17:1b:42:97:b3:c9:5d:e1:a9:b5:c1:4a:bd:69:12:1a:
         8e:9d:99:38:ce:91:c8:20:82:e6:dc:b1:5b:b3:7a:15:ae:6f:
         ad:90:1c:35:c9:b3:9d:dd:96:d6:4f:31:f0:aa:fd:ea:f1:76:
         04:73:c0:57:e8:a8:80:20:82:17:e2:d7:1c:f7:5c:4b:39:6d:
         c8:15:43:81:8f:87:fd:eb:4b:ff:77:7b:3f:56:94:42:2d:ac:
         fa:6c:7e:3e:1a:3d:9a:dc:5b:2d:07:8d:7a:da:c9:ea:55:56:
         0e:cc:7d:c9:ec:30:a7:d4:24:94:2e:85:de:11:ba:34:ea:01:
         d0:79:43:42:0a:c6:0e:07:1d:10:b9:53:a9:c1:ad:65:55:d5:
         73:bc:1d:8b:65:bb:d1:36:61:5a:fe:4d:a7:4e:d9:9d:41:27:
         5b:97:fc:f0:5e:ff:30:f9:b6:10:92:61:cd:30:ca:c6:d8:bb:
         8c:df:fe:0a:31:e3:29:90:62:6c:3d:4a:b9:e5:ad:7a:42:9a:
         32:25:f3:01:65:49:af:25:9e:f9:30:f7:ea:23:49:15:1e:57:
         9c:f8:62:77:2a:36:dc:a6:d5:02:13:3e:d1:ba:91:88:7f:a1:
         e3:bc:81:2f:8d:98:0a:b2:51:21:8c:56:56:57:bc:f6:a8:2b:
         48:fa:e6:13
```

## Launch a Kapow! Server Using the Newly Created Certs

``` console
$ kapow server  --certfile pki/issued/mykapowserver.crt     \
                --clientcafile pki/ca.crt                   \
                --clientauth                                \
                --keyfile pki/private/mykapowserver.key     \
                --debug &
2020/12/22 11:39:18.169479 UserServer using CA certs from pki/ca.crt
2020/12/22 11:39:18.169487 UserServer listening at 0.0.0.0:8080
2020/12/22 11:39:18.169657 ControlServer listening at localhost:8081
2020/12/22 11:39:18.169686 DataServer listening at localhost:8082
```

## Restrict Access to an Route Depending on the Client's DN

Create `./test_mTLS` with this content:
``` sh
#!/usr/bin/env sh
# /!\ DON'T REMOVE THIS! OTHERWISE UNAUTHORIZED CLIENTS WILL GAIN ACCESS!
set -e
echo CN=authorizedclient | validsslclient

# Put your logic beyond this point.
kapow set /response/body "You have been granted access.  Use it wisely."
```
and set the executable bits:

``` sh
chmod +x ./test_mTLS
```

Add a route to *Kapow!* to be handled by `./test_mTLS`:
``` console
$ kapow route add /restricted -e ./test_mTLS
{"id":"e40c13b5-444a-11eb-9c32-002b671b12f9","method":"GET","url_pattern":"/restricted","entrypoint":"./test_mTLS","command":"","index":0,"debug":true} 
```

Make an authenticated request:
``` console
$ curl --cacert pki/ca.crt                        \
       --cert pki/issued/authorizedclient.crt     \
       --key pki/private/authorizedclient.key     \
       https://mykapowserver:8080/restricted
2020/12/22 11:46:49.501552 605dc278-444b-11eb-9c32-002b671b12f9 validsslclient: Found valid user: 'CN=authorizedclient'
127.0.0.1:35120 605dc278-444b-11eb-9c32-002b671b12f9 - [22/Dec/2020:11:46:49 +0000] "GET /restricted HTTP/2.0" 200 45 "-" "curl/7.74.0"
You have been granted access.  Use it wisely. 
```
