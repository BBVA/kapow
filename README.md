![Kapow!](https://trello-attachments.s3.amazonaws.com/5c6edee98297dc18aa4e2b63/960x720/ff8d28fc24af11e3295afa5a9665bdc0/kapow-1601675_960_720.png)

**Kapow!** allows you to leverage the Ultimate Power™ of the UNIX® shell via HTTP.

**Kapow!** is in the process of being defined by a [specification](/spec/)

# Sample usage
## Clone the project
```bash
# clone this project
```

## Build the kapow! docker image
```bash
docker build -t bbva/kapow:0.1 /path/to/kapow/poc
```

## Build a docker image for running the nmap example
```bash
docker build -t kapow-nmap /path/to/kapow/poc/examples/nmap
```

## Run kapow
```bash
docker run \
        -it \
        -p 8080:8080 \
        kapow-nmap
```
which will output something like this:
```
======== Running on http://0.0.0.0:8080 ========
(Press CTRL+C to quit)
Route created POST /list/{ip}
ROUTE_8ed01c48_bf23_455a_8186_a1df7ab09e48
bash-4.4#
```


## Test /list endpoint
In another terminal, try running:
```bash
curl -X POST http://localhost:8080/list/github.com
```
which will respond something like:
```
Starting Nmap 7.70 ( https://nmap.org ) at 2019-05-10 14:01 UTC
Nmap scan report for github.com (140.82.118.3)
rDNS record for 140.82.118.3: lb-140-82-118-3-ams.github.com
Nmap done: 1 IP address (0 hosts up) scanned in 0.04 seconds

```
et voilà !

# License

This project is distributed under the [Apache License 2.0](/LICENSE).
