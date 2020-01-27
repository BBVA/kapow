# Network Scanner (nmap) as a Service

Run a long network scan in background with support for webhook on completion.

* The user can define the destination IP and port(s).
* The service answers immediately with a `jobid`.
* If a webhook url is defined it will be called on completion with the result and the jobid.
* At any moment the user can request the status of the scan at /scan/{jobid}

## How to run it

```
$ kapow server NetworkScanner.pow
```


## How to consume it

* Scan your own host
```
$ curl --data 'ports=1-65535&ip=127.0.0.1' http://localhost:8080/scan
{
  "job": "dba2edbc-527d-453a-9c25-0608bb8f06da"
}
```

* Grab the result

```
$ curl -v http://localhost:8080/scan/dba2edbc-527d-453a-9c25-0608bb8f06da
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE nmaprun>
<?xml-stylesheet href="file:///usr/bin/../share/nmap/nmap.xsl" type="text/xsl"?>
<!-- Nmap 7.80 scan initiated Mon Jan 27 20:01:42 2020 as: nmap -Pn -n -p 1-65535 -oX dba2edbc-527d-453a-9c25-0608bb8f06da.running.xml -&#45; 127.0.0.1 -->
<nmaprun scanner="nmap" args="nmap -Pn -n -p 1-65535 -oX dba2edbc-527d-453a-9c25-0608bb8f06da.running.xml -&#45; 127.0.0.1" start="1580151702" startstr="Mon Jan 27 20:01:42 2020" version="7.80" xmloutputversion="1.04">
<scaninfo type="connect" protocol="tcp" numservices="65535" services="1-65535"/>
<verbose level="0"/>
<debugging level="0"/>
<host starttime="1580151702" endtime="1580151704"><status state="up" reason="user-set" reason_ttl="0"/>
<address addr="127.0.0.1" addrtype="ipv4"/>
<hostnames>
</hostnames>
<ports><extraports state="closed" count="65531">
<extrareasons reason="conn-refused" count="65531"/>
</extraports>
<port protocol="tcp" portid="631"><state state="open" reason="syn-ack" reason_ttl="0"/><service name="ipp" method="table" conf="3"/></port>
<port protocol="tcp" portid="8080"><state state="open" reason="syn-ack" reason_ttl="0"/><service name="http-proxy" method="table" conf="3"/></port>
<port protocol="tcp" portid="8081"><state state="open" reason="syn-ack" reason_ttl="0"/><service name="blackice-icecap" method="table" conf="3"/></port>
<port protocol="tcp" portid="8082"><state state="open" reason="syn-ack" reason_ttl="0"/><service name="blackice-alerts" method="table" conf="3"/></port>
</ports>
<times srtt="78" rttvar="70" to="100000"/>
</host>
<runstats><finished time="1580151704" timestr="Mon Jan 27 20:01:44 2020" elapsed="1.91" summary="Nmap done at Mon Jan 27 20:01:44 2020; 1 IP address (1 host up) scanned in 1.91 seconds" exit="success"/><hosts up="1" down="0" total="1"/>
</runstats>
</nmaprun>

```

*If you receive a 202 it means the scan is still on progress*
