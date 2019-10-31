# Installing Kapow!

Kapow! has a reference implementation in Go that is under active develpment right
now.  If you want to start using kapow you can:
* Download a binary (linux, at this moment) from our
[releases](https://github.com/BBVA/kapow/releases) section
* Install the package with the get command  (you need the Go runtime installed
and [configured](https://golang.org/cmd/go/))
```Shell
    go get -u github.com/BBVA/kapow
```


# Examples

Below are some examples on how to define and invoke routes in Kapow!

As you can see Kapow! binary is a server and a CLI that you can use to configure
a running server.  The server exposes an [API](/spec#http-control-api) that you
can invoke directly if you want.

In order to get information from the request that fired the scrip execution and
to help you in composing the response, the server exposes
some [resources](/spec#handlers) to interact with from the script.


## The mandatory Hello World (for WWW fans)

First you create a pow file named `greet.pow` with the following contents:

```Shell
    kapow route add /greet -c 'name=$(kapow get /request/params/name); echo Hello ${name:-World} | kapow set /response/body'
```

note you have to escape as the command will run on a shell itself. Then, you
execute:

```Shell
    kapow server greet.pow
```

to start a Kapow! server exposing your service.  Now you can check that it works
as intended with good ole' ``curl``:

```Shell
    curl localhost:8080/greet
    Hello World

    curl localhost:8080/greet?name=friend
    Hello friend
```

If you want to work with JSON you can use this version of the pow
`greet-json.pow`

```Shell
  kapow route add -X POST /greet -c 'kapow route add -X POST /greet -c 'who=$(kapow get /request/body | jq -r .name); kapow set /response/status 201; jq --arg value "${who:-World}" -n \{name:\$value\} | kapow set /response/body''
```

that uses [jq](https://stedolan.github.io/jq/) to allow you working with json
from the command line.  Check that it works with

```Shell
  curl -X POST -H "Content-Type: application/json" -d '{"name": "friend"}' localhost:8080/greet
  {"name": "friend" }

  curl -X POST -H "Content-Type: application/json" -d '' localhost:8080/greet
  {"name": "World"}
```


## The mandatory Echo (for UNIX fans)

First you create a pow file named `echo.pow` with the following contents:

```Shell
    kapow route add -X POST /echo -c 'kapow get /request/body | kapow set /response/body'
```

then, you execute:

```Shell
    kapow server echo.pow
```

and you can check that it works as intended with good ole' `curl`:

```Shell
    curl -X POST -d '1,2,3... testing' localhost:8080/echo
    1, 2, 3, 4, 5, 6, 7, 8, 9, testing
```

If you send a big file and want to see the content back as a real-time stream
you can use this version `echo-stream.pow`

```Shell
    kapow route add -X POST /echo -c 'kapow get /request/body | kapow set /response/stream'
```


## The multiline fun

Unless you're a hardcore Perl hacker, you'll probably need to write your stuff
over more than one line in order to avoid the mess we saw on our json greet
version.

Don't worry, we need to write several lines, too.  Bash, in its magnificent
UNIX® style, provides us with the
[here-documents](https://www.gnu.org/software/bash/manual/bash.html#Here-Documents)
mechanism that we can leverage precisely for this purpose.

Imagine we want to return both the standard output and a generated file from a
command execution. Let's write a `log-and-stuff.pow` file with the following content:

```Shell
  kapow route add /log_and_stuff - <<- 'EOF'
    echo this is a quite long sentence and other stuff | tee log.txt | kapow set /response/body
    cat log.txt | kapow set /response/body
  EOF
```

then we serve it with `kapow`:

```Shell
    kapow server log-and-stuff.pow
```

Yup.  As simple as that.  You can check it.

```Shell
    curl localhost:8080/log_and_stuff
    this is a quite long sentence and other stuff
    this is a quite long sentence and other stuff
```


## Interact with other systems

You can leverage all the power of the shell in your scripts and interact with
other systems by using all the available tools.  Write a
`log-and-stuff-callback.pow` file with the following content:

```Shell
  kapow route add /log_and_stuff - <<- 'EOF'
    callback_url=$(kapow get /request/params/callback)
    echo this is a quite long sentence and other stuff | tee log.txt | kapow set /response/body
    echo sending to $callback_url | kapow set /response/body
    curl -X POST --data-binary @log.txt $callback_url | kapow set /response/body
  EOF
```

serve it with `kapow`:

```Shell
    kapow server log-and-stuff-callback.pow
```

and finally check it.

```Shell
    curl localhost:8080/log_and_stuff?callback=nowhere.com
    this is a quite long sentence and other stuff
    sending to nowhere.com
    <html>
    <head><title>405 Not Allowed</title></head>
    <body>
    <center><h1>405 Not Allowed</h1></center>
    <hr><center>nginx</center>
    </body>
    </html>
```

You must be aware that you must have all the dependencies you use in your
scripts installed in the host that will run the Kapow! server.

In addition, a pow file can contain as much routes as you like so you can start
a server with several routes configured in one shoot.

# Sample Docker usage


## Clone the project

```Shell
    git clone https://github.com/BBVA/kapow.git
```


## Build the kapow! docker image

```Shell
    make docker
```

Now you have a container image with all the above pow files copied in /tmp so
you can start each example by running

```Shell
    docker run --rm -p 8080:8080 docker server example.pow
```


## Build a docker image for running the nmap example

```Shell
    cd /path/to/kapow/poc/examples/nmap; docker build -t kapow-nmap .
```

## Run kapow

```Shell
    docker run \
            -d \
            -p 8080:8080 \
            kapow-nmap
```

which will output something like this:

```Shell
   e7da20c7d9a39624b5c56157176764671e5d2d8f1bf306b3ede898d66fe3f4bf
```


## Test /list endpoint

In another terminal, try running:

```Shell
    curl http://localhost:8080/list/github.com
```

which will respond something like:

```Shell
    Starting Nmap 7.70 ( https://nmap.org ) at 2019-05-10 14:01 UTC
    Nmap scan report for github.com (140.82.118.3)
    rDNS record for 140.82.118.3: lb-140-82-118-3-ams.github.com
    Nmap done: 1 IP address (0 hosts up) scanned in 0.04 seconds
```

et voilà !


# License

This project is distributed under the [Apache License 2.0](/LICENSE).
