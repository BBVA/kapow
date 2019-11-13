# What is it?

Kapow! is an adapter between the world of Pure UNIX® Shell and a HTTP service.
It lies between these two worlds, making your life easier.

Kapow! allows yo to publish a simple shell script as a REST HTTP service so you
can delegate in others its execution as they don't need access to the host in
which the command is ran.  Those repetitive tasks that everybody ask you to do
because they require administrative access to some host can be published through
a Kapow! server deployed in that host and the users who need the results can
invoke it directly using a easy to use interfaz, a HTTP request.

In the tradicional way you needed to create a project, manage its dependencies,
code the server (probably including only a command execution) and deploy it
somewhere.  And that's fine... until you find yourself again in the same
situation with another awesome command.

### From now on you can put Kapow! into action

- Create a pow file containing an api call to Kapow! for creating the route
that will publish your command, lets's call it `greet.pow`
  ```sh
    kapow route add /greet -c 'name=$(kapow get /request/params/name); echo Hello ${name:-World} | kapow set /response/body'
  ```
- Start the Kapow! server providing your pow file to configure the route
  ```sh
    kapow server greet.pow
  ```
- check that all it is working as intended using `curl`
  ```sh
    $ curl localhost:8080/greet
    Hello World

    $ curl localhost:8080/greet?name=friend
    Hello friend
  ```

# Installing Kapow!

Kapow! has a reference implementation in Go that is under active development right
now.  If you want to start using Kapow! you can:
* Download a binary (only linux is available , at this moment) from our
[releases](https://github.com/BBVA/kapow/releases/latest) section
* Install the package with Go's `get` command (you need the Go runtime installed
and [configured](https://golang.org/cmd/go/))
```sh
go get -u github.com/BBVA/kapow
```
