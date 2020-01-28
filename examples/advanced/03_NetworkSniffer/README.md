# Network Sniffer (tcpdump) as a Service

Provides an HTTP service that allows the user to sniff the network in real time.  The packet capture data is served as an HTTP stream that can be injected to a packet analysis tool on the fly.


## How to run it

For the sake of simplicity, run:

```
$ sudo -E kapow server NetworkSniffer.pow
```

In a production environment, tcpdump should be run with the appropiate
permissions, but kapow can (and should) run as an unprivileged user.


## How to consume it

```
$ curl http://localhost:8080/sniff | sudo -E wireshark -k -i -
```
