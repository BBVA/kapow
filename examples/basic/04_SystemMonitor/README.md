# System Monitoring as a Service

Expose several system properties and logs.

## How to run it

```
$ kapow server SystemMonitor.pow
```


## How to consume it

<details><summary>List files and directories</summary>

```
$ curl -s http://localhost:8080/file/var/log/cups
total 60
drwxr-xr-x  2 root root  4096 Dec  2 06:54 .
drwxr-xr-x 15 root root  4096 Jan 27 07:34 ..
-rw-r--r--  1 root cups     0 Dec  2 06:54 access_log
-rw-r--r--  1 root cups   125 Dec  1 19:12 access_log.1
-rw-r--r--  1 root cups   254 Nov 23 12:17 access_log.2
-rw-r--r--  1 root cups   125 Sep  8 14:41 access_log.3
-rw-r--r--  1 root cups   634 May  1  2019 access_log.4
-rw-r--r--  1 root cups     0 Sep  9  2018 error_log
-rw-r--r--  1 root cups 17312 Sep  3  2018 error_log.1
-rw-r--r--  1 root cups     0 Dec  2 06:54 page_log
-rw-r--r--  1 root cups   128 Dec  1 19:12 page_log.1
-rw-r--r--  1 root cups   188 Nov 23 12:17 page_log.2
-rw-r--r--  1 root cups   108 Sep  8 14:41 page_log.3
-rw-r--r--  1 root cups   465 May  1  2019 page_log.4
```
</details>

<details><summary>List processes</summary>
```
$ curl -s http://localhost:8080/process
nil        46717  0.0  0.0 111224  8196 pts/2    Sl   16:48   0:00 kapow server SystemMonitor.pow
root       47405  0.0  0.0      0     0 ?        I    16:50   0:00 [kworker/3:1-mm_percpu_wq]
root       47406  0.0  0.0      0     0 ?        I    16:50   0:00 [kworker/0:1]
nil        47479  0.4  0.0  60020 31124 pts/2    S+   16:50   0:01 vim README.md
root       47819  0.0  0.0      0     0 ?        I    16:52   0:00 [kworker/7:2-mm_percpu_wq]
root       47823  0.6  0.0      0     0 ?        I<   16:52   0:01 [kworker/u17:0-hci0]
nil        48097  0.0  0.1 605928 56748 ?        Sl   16:54   0:00 /usr/lib/chromium/chromium --type=renderer --disable-webrtc-apm-in-audio-service --field-trial-handle=8593857243625867040,643120771000881201,131072 --lang=en-US --disable-oor-cors --enable-auto-reload --num-raster-threads=4 --enable-main-frame-before-activation --service-request-channel-token=18057504326134146058 --renderer-client-id=112 --no-v8-untrusted-code-mitigations --shared-files=v8_context_snapshot_data:100,v8_natives_data:101
nil        48345  0.0  0.0   7124  2804 pts/2    S    16:56   0:00 /bin/sh -c ps -aux | kapow set /response/body
nil        48346  0.0  0.0   9392  3324 pts/2    R    16:56   0:00 ps -aux
nil        48347  0.0  0.0 109304  7080 pts/2    Sl   16:56   0:00 kapow set /response/body
```
</details>
