10-SECOND PROPOSAL
===================

*kapow* is a specialized language for marrying the **web** and the **shell**.


DESCRIPTION
===========

The **web** and the **shell** are two different beasts, both packed with history.

There are some concepts in HTTP and the shell that **resemble each other**.

```
                 +------------------------+-------------------------+
                 | HTTP                   | SHELL                   |
  +--------------+------------------------+-------------------------+
  | Input        | POST form-encoding     | Command line parameters |
  | Parameters   | GET parameters         | Environment variables   |
  |              | Headers                |                         |
  |              | Serialized body (JSON) |                         |
  +--------------+------------------------+-------------------------+
  | Data Streams | Response/Request Body  | Stdin/Stdout/Stderr     |
  |              | Websocket              | Input/Output files      |
  |              | Uploaded files         |                         |
  +--------------+------------------------+-------------------------+
  | Control      | Status codes           | Signals                 |
  |              | HTTP Methods           |                         |
  +--------------+------------------------+-------------------------+
```

Any tool designed to give an HTTP interface to an existing shell command **must map concepts of boths**. For example: 

- "GET parameters" to "Command line parameters"
- "Headers" to "Environment variables"
- "Stdout" to "Response body"

*kapow* is not opinionated about the different ways you can map both worlds. Instead it provides a concise language used to express the mapping and a set of common defaults.


Why not tool...?
----------------

All the alternatives we found are **rigid** about how they match between HTTP and shell concepts.

* [shell2http](https://github.com/msoap/shell2http): HTTP-server to execute shell commands. Designed for development, prototyping or remote control. Settings through two command line arguments, path and shell command.
* [websocketd](https://github.com/joewalnes/websocketd): Turn any program that uses STDIN/STDOUT into a WebSocket server. Like inetd, but for WebSockets.
* [webhook](https://github.com/adnanh/webhook): webhook is a lightweight incoming webhook server to run shell commands.
* [gotty](https://github.com/yudai/gotty): GoTTY is a simple command line tool that turns your CLI tools into web applications. (For interactive commands only)

Tools with a rigid matching **can't evade** *[impedance mismatch](https://haacked.com/archive/2004/06/15/impedance-mismatch.aspx/)*. Resulting is an easy-to-use software, convenient in some scenarios but incapable in others.


Why not my good-old programming language...?
--------------------------------------------

* Boilerplate
* Custom code = More bugs
* Security issues (Command injection, etc)
* Dependency on developers
* "A programming language is low level when its programs require attention to the irrelevant" Alan Perlis

*kapow* aims to be halfway from one of the mentioned tools and a general programming language. A limited scripting language. Think of *awk*, 20 years later, for HTTP.


Example
-------

Imagine you want to improve your monitoring system to be able to check if a WiFi in a remote building is working. Suppose you already have a wireless capable linux server on that building.

What about exposing a wifi scan as a service?

```
$ iwlist scan | grep 'ESSID:"BBVA"'
```

The above command will list the visible WiFI networks and filter for the one we are interested in (**BBVA**); exiting with code **0** when found and **1** otherwise.

We want to signal the WiFI status with a meaningful HTTP status code: 200 for UP and 503 for DOWN.

So the only thing our service has to do is translate from the **command's exit code** to a meaningful **HTTP status code**.

Our *kapow* program (monitor.pow) should look **something like** this:

```
/monitor/wifi =
  $(iwlist scan | grep 'ESSID:"BBVA')
  | 0          = status 200  # OK
  | 1          = status 503  # Service Unavailable
```

With this code, the user can run the service with:

```
$ kapow monitor.pow
```
And then:

1. *kapow* will open a port with an HTTP server serving the URI **/monitor/wifi**.
2. When requested `curl http://<server>/monitor/wifi`. The command is executed resulting in a WiFI scan.
3. The command's **exit code** is translated to an HTTP **status code** and returned.

USE CASES OR VALUE
==================

* Reuse 40+ years of existing computer programs as microservices (nanoservices?).
* Expose command line only tools as services. Eliminating the need of interactive SSH sessions.
* Fine-grained access control to CLI based on options/parameters


REQUIRED SKILLS
===============

- Knowledge of HTTP and the command line.
- Programming.
- Low-level stuff...


ESTIMATION
==========

- 2 weeks


OUTPUT
======

- A collection of **concrete** use cases, 5 to 10.
- A DSL design (over paper) to implement the use cases.
- A working **dirty** proof of concept.
