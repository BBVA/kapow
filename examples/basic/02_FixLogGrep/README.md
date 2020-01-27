# Fix Log Grep as a Service

A simple service that exposes log entries for files not found on our Apache Web Server.

## How to run it

```
$ kapow server FixLogGrep.pow
```


## How to consume it

```
$ curl http://localhost:8080/apache-errors
[Fri Feb 01 22:07:57.154391 2019] [core:info] [pid 7:tid 140284200093440] [client 172.17.0.1:50756] AH00128: File does not exist: /usr/var/www/mysite/favicon.ico
[Fri Feb 01 22:07:57.808291 2019] [core:info] [pid 8:tid 140284216878848] [client 172.17.0.1:50758] AH00128: File does not exist: /usr/var/www/mysite/favicon.ico
[Fri Feb 01 22:07:57.878149 2019] [core:info] [pid 8:tid 140284208486144] [client 172.17.0.1:50758] AH00128: File does not exist: /usr/var/www/mysite/favicon.ico
```
