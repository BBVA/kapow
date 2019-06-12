![Kapow!](https://trello-attachments.s3.amazonaws.com/5c6edee98297dc18aa4e2b63/960x720/ff8d28fc24af11e3295afa5a9665bdc0/kapow-1601675_960_720.png)

**Kapow!** allows you to leverage the Ultimate Power™ of the UNIX® shell via HTTP.

# Warning!!!

**Warning!! Kapow!** is in the process of being defined by a [specification](/spec/) the provided code it's an unstable Proof of Concept. 

## How Kapow! born

Some awesome history it's comming.

# What is Kapow!
Kapow! is an adapter between the world of Pure Unix Shell and a HTTP service.

Some tasks are more convenient in the shell. Like cloud interactions, or some adminstration tools. In the other side some tasks are more convenient as a service. Like DevSecOps tooling.

Kapow! lies between this two worlds, making you life easier. Maybe you wonder about how kind of magic happen, if you want to know the very last detail just read our [spec](/spec/). If you need to know how Kapow! can you help first, let start with a common situation.

Think about that awesome command that you use every day, something very confortable, like `cloudx storage ls /backups`. Then someone ask you for an specific backup, so you go into the machine throught ssh, execute your command (maybe you `grep` it), copy the result and send it. And everything it's ok... for the 100 first times.

Then you decide, let's use the API for this and generate an awesome web server with it. So, create a proyect, manage their dependencies, code the server, parse the request, learn how to use the API, call the API and deploy somewere. And everything it's ok... until you find again in the same situation with another awesome command.

The awensomeness of unix commands it's infinite, so you'll be in this situation infinite times!!. Let's put Kapow! in this equation.

With Kapow! when someone ask you for an specific backup (remember your confortable command?) you create a pow file named "backups.pow" that contains:
```bash
kapow route add '/backups' \
    -c 'cloudx storage ls /backups | grep $(request /params/query) | response /body'
```

And execute in the machine with the command:
```bash
kapow server backups.pow
```

And that's it. Done. Do you like it? yes? let's start learning a litte more.

## The obligatory Hello World (for www boys)

First you must create a pow file named "hello.pow" with the following contents:
```bash
kapow route add "/greet" -c "echo 'hello world' | response /body"
```

Then you must execute:
```bash
kapow server hello.pow
```
And you can check the works as intented with our good old curl:
```bash
curl localhost:8080/greet
```

## The obligatory Echo (for UNIX boys)
First you must create a pow file named "echo.pow" with the following contents:
```bash
kapow route add -X POST "/echo" -c "request /body | response /body"
```

Then you must execute:
```bash
kapow server echo.pow
```
And you can check the works as intented with our good old curl:
```bash
curl -X POST -d "1,2,3... testing" localhost:8080/echo
```

## The multiline fun
Unless you're a hardcore Perl boy, you need write your stuff in more than one line.

Don't worry, we need write several lines too. Bash, in their magnificent UNIX flavour bring us HERE doc.

Let's write a "multiline.pow" file with the following content:
```bash
kapow route add "/log_and_love" - <<-'EOF'
echo "[$(date)] and stuff" >> stuff.log
echo "love" | response /body
EOF
```

And then we serve it with kapow:
```bash
kapow server multiline.pow
```

As simple as that.

# Sample Docker usage
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
curl http://localhost:8080/list/github.com
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
