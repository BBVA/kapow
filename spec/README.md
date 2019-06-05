# Kapow!


## Why?

Because we think that:

- UNIX® is great and we love it
- The UNIX® shell is great
- HTTP interfaces are convenient and everywhere
- CGI is not a good way to mix them


## How?

So, how we can mix the **web** and the **shell**?  Let's see...

The **web** and the **shell** are two different beasts, both packed with
history.

There are some concepts in HTTP and the shell that **resemble each other**.

  |                        | HTTP                                                                           | Shell                                              |
  |------------------------|--------------------------------------------------------------------------------|----------------------------------------------------|
  | Input<br /> Parameters | POST form-encoding<br >Get parameters<br />Headers<br />Serialized body (JSON) | Command line parameters<br />Environment variables |
  | Data Streams           | Response/Request Body<br />Websocket<br />Uploaded files                       | stdin/stdout/stderr<br />Input/Output files        |
  | Control                | Status codes<br />HTTP Methods                                                 | Signals<br />Exit Codes                            |

Any tool designed to give an HTTP interface to an existing shell command
**must map concepts from both domains**.  For example:

- "GET parameters" to "Command line parameters"
- "Headers" to "Environment variables"
- "stdout" to "Response body"

Kapow! is not opinionated about the different ways you can map both worlds.
Instead, it provides a concise set of tools, with a set of sensible defaults,
allowing the user to express the desired mapping in an explicit way.


### Why not tool "X"?

All the alternatives we found are **rigid** about the way they match between
HTTP and shell concepts.

* [shell2http](https://github.com/msoap/shell2http): HTTP-server to execute
  shell commands.  Designed for development, prototyping or remote control.
Settings through two command line arguments, path and shell command.
* [websocketd](https://github.com/joewalnes/websocketd): Turn any program that
  uses STDIN/STDOUT into a WebSocket server.  Like inetd, but for WebSockets.
* [webhook](https://github.com/adnanh/webhook): webhook is a lightweight
  incoming webhook server to run shell commands.
* [gotty](https://github.com/yudai/gotty): GoTTY is a simple command line tool
  that turns your CLI tools into web applications.  Note that this tool works
  only with interactive commands.
* [shell-microservice-exposer](https://github.com/jaimevalero/shell-microservice-exposer):
  Expose your own scripts as a cool microservice API dockerizing it.

Tools with a rigid matching **can't evade** *[impedance
mismatch](https://haacked.com/archive/2004/06/15/impedance-mismatch.aspx/)*.
Resulting in an easy-to-use software, convenient in some scenarios but
incapable in others.


### Why not my good-old programming language "X"?

* Boilerplate
* Custom code = More bugs
* Security issues (command injection, etc)
* Dependency on developers
* *"A programming language is low level when its programs require attention to
  the irrelevant."<br />&mdash;Alan Perlis*
* *"There is more Unix-nature in one line of shell script than there is in ten
  thousand lines of C."<br />&mdash;Master Foo*


### Why not CGI?

* CGI is also **rigid** about how it matches between HTTP and UNIX® process
  concepts.  Notably CGI *meta-variables* are injected into the script's
  environment; this behavior can and has been exploited by nasty attacks such as
  [Shellshock](https://en.wikipedia.org/wiki/Shellshock_(software_bug)).
* Trying to leverage CGI from a shell script could be less cumbersome in some
  cases but possibly being more error-prone.  For instance, since in CGI
  everything written to the standard output becomes the body of the response,
  any leaked command output would corrupt the HTTP response.


## What?

We named it Kapow!.  It is pronounceable, short and meaningless...  like every
good UNIX® command ;-)

TODO: Definition

TODO: Intro to Architecture


### API

Kapow! server interacts with the outside world only through its HTTP API.  Any
program making the correct HTTP request to a Kapow! server can change its
behavior.

Kapow! exposes two distinct APIs, a control API and a data API, described
below.


# HTTP Control API

It allows you to configure the Kapow! service. This API is available during the
whole lifetime of the server.


## Design Principles

* All requests and responses will leverage JSON as the data encoding method.
* The API calls responses will have two distinct parts:
  * The HTTP status code (e.g., `400`, which is a bad request).  The target
    audience of this information is the client code.  The client can thus use
    this information to control the program flow.
  * The JSON-encoded message.  The target audience in this case is the human
    operating the client.  The human can use this information to make a
    decision on how to proceed.
* All successful API calls will return a representation of the *final* state
  attained by the objects which have been addressed (either requested, set or
  deleted).

For instance, given this request:
```http
HTTP/1.1 GET /routes
```

an appropiate reponse may look like this:
```http
200 OK
Content-Type: application/json
Content-Length: 189

[
  {
    "method": "GET",
    "url_pattern": "/hello",
    "entrypoint": null,
    "command": "echo Hello World | response /body",
    "index": 0
  }
]
```


## API Elements

Kapow! provides a way to control its internal state through these elements.


### Routes

Routes are the mechanism that allows Kapow! to find the correct program to
respond to an external event (e.g.  an incoming HTTP request).


#### List routes

Returns JSON data about the current routes.

* **URL**: `/routes`
* **Method**: `GET`
* **Success Responses**:
  * **Code**: `200 OK`<br />
    **Content**:<br />
    ```json
    [
      {
        "method": "GET",
        "url_pattern": "/hello",
        "entrypoint": null,
        "command": "echo Hello World | response /body"
      },
      {
        "method": "POST",
        "url_pattern": "/bye",
        "entrypoint": null,
        "command": "echo Bye World | response /body"
      }
    ]
    ```
* **Sample Call**: `$ curl $KAPOW_URL/routes`
* **Notes**: Currently all routes are returned; in the future, a filter may be accepted.


#### Append route

Accepts JSON data that defines a new route to be appended to the current routes.

* **URL**: `/routes`
* **Method**: `POST`
* **Header**: `Content-Type: application/json`
* **Data Params**:<br />
  ```json
  {
    "method": "GET",
    "url_pattern": "/hello",
    "entrypoint": null,
    "command": "echo Hello World | response /body"
  }
  ```
* **Success Responses**:
  * **Code**: `201 Created`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**:<br />
    ```json
    {
      "method": "GET",
      "url_pattern": "/hello",
      "entrypoint": null,
      "command": "echo Hello World | response /body",
      "index": 0
    }
    ```
* **Error Responses**:
  * **Code**: `400 Malformed JSON`
  * **Code**: `400 Invalid Data Type`
  * **Code**: `400 Invalid Route Spec`
  * **Code**: `400 Missing Mandatory Field`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**:<br />
    ```json
    {
      "missing_mandatory_fields": [
        "url_pattern",
        "command"
      ]
    }
    ```
* **Sample Call**:<br />
    ```sh
    $ curl -X POST --data-binary @- $KAPOW_URL/routes <<EOF
    {
      "method": "GET",
      "url_pattern": "/hello",
      "entrypoint": null,
      "command": "echo Hello World | response /body",
      "index": 0
    }
    EOF
    ```
* **Notes**:
  * A successful request will yield a response containing all the effective
    parameters that were applied.
  * Kapow! won't try to validate the submitted command.  Any errors will happen
    at runtime, and trigger a `500` status code.


#### Insert a route

  Accepts JSON data that defines a new route to be inserted at the specified
  index to the current routes.

* **URL**: `/routes`
* **Method**: `PUT`
* **Header**: `Content-Type: application/json`
* **Data Params**:<br />
  ```json
  {
    "method": "GET",
    "url_pattern": "/hello",
    "entrypoint": null,
    "command": "echo Hello World | response /body",
  }
  ```
* **Success Responses**:
  * **Code**: `200 OK`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**:<br />
    ```json
    {
      "method": "GET",
      "url_pattern": "/hello",
      "entrypoint": null,
      "command": "echo Hello World | response /body",
      "index": 0
    }
    ```
* **Error Responses**:
  * **Code**: `400 Malformed JSON`
  * **Code**: `400 Invalid Data Type`
  * **Code**: `400 Invalid Route Spec`
  * **Code**: `400 Missing Mandatory Field`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**:<br />
    ```json
    {
      "missing_mandatory_fields": [
        "url_pattern",
        "command"
      ]
    }
    ```
  * **Code**: `400 Invalid Index Type`
  * **Code**: `400 Index Already in Use`
  * **Code**: `404 Invalid Index`
  * **Code**: `404 Invalid Route Spec`
* **Sample Call**:<br />
    ```sh
    $ curl -X PUT --data-binary @- $KAPOW_URL/routes <<EOF`
    {
      "method": "GET",
      "url_pattern": "/hello",
      "entrypoint": null,
      "command": "echo Hello World | response /body",
      "index": 0
    }
    EOF
    ```
* **Notes**:
  * Route numbering starts at zero.
  * When `index` is not provided or is less than `0` the route will be inserted
    first, effectively making it index `0`.
  * Conversely, when `index` is greater than the number of entries on the route
    table, it will be inserted last.
  * A successful request will yield a response containing all the effective
    parameters that were applied.


#### Delete a route

Removes the route identified by `:id`.

* **URL**: `/routes/:id`
* **Method**: `DELETE`
* **Success Responses**:
  * **Code**: `200 OK`<br />
    **Content**:<br />
    ```json
    {
      "method": "GET",
      "url_pattern": "/hello",
      "entrypoint": null,
      "command": "echo Hello World | response /body",
      "index": 0
    }
    ```
* **Error Responses**:
  * **Code**: `404 Not Found`
* **Sample Call**:<br />
  ```sh
  $ curl -X DELETE $KAPOW_URL/routes/ROUTE_1f186c92_f906_4506_9788_a1f541b11d0f
  ```
* **Notes**:


# HTTP Data API

It is the channel through which the actual HTTP data flows during the
request/response cycle, both reading from the request as well as writing to the
response.


## Design Principles

* According to established best practices we use the HTTP methods as follows:
  * `GET`: Read data without any side-effects.
  * `PUT`: Overwrite existing data.
* The API calls responses will have two distinct parts:
  * The HTTP status code (e.g., `400`, which is a bad request).  The target
    audience of this information is the client code.  The client can thus use
    this information to control the program flow.
  * The HTTP reason phrase.  The target audience in this case is the human
    operating the client.  The human can use this information to make a
    decision on how to proceed.
* Regarding HTTP request and response bodies:
  * The response body will be empty in case of error.
  * It will transport binary data in other case.


## API Elements

The data API consists of a single element, the handler.


### Handlers

Handlers are in-memory data structures exposing the data of the current request
and response.

Each handler is identified by a `handler_id` and provide access to the
following resource paths:

```
/                               The root of the resource paths tree
│
├─ request                      All information related to the HTTP request.  Read-Only
│  ├──── method                 Used HTTP Method (GET, POST)
│  ├──── path                   Complete URL path (URL-unquoted)
│  ├──── matches                Previously matched URL path parts
│  │     └──── <name>
│  ├──── params                 URL parameters (post ? symbol)
│  │     └──── <name>
│  ├──── headers                HTTP request headers
│  │     └──── <name>
│  ├──── cookies                HTTP request cookie
│  │     └──── <name>
│  ├──── form                   Form-urlencoded form fields
│  │     └──── <name>
│  ├──── files                  Files uploaded via multi-part form fields
│  │     └──── <name>
│  │           └──── filename   Original file name
│  │           └──── content    The file content
│  └──── body                   HTTP request body
│
└─ response                     All information related to the HTTP request.  Write-Only
   ├──── status                 HTTP status code
   ├──── headers                HTTP response headers
   │     └──── <name>
   ├──── cookies                HTTP request cookie
   │     └──── <name>
   ├──── body                   Response body.  Mutually exclusive with response/stream
   └──── stream                 Chunk-encoded body.  Streamed response.  Mutually exclusive with response/body
```


#### Example Keys

- Read the request URL path.
  - Scenario: Request URL is `http://localhost:8080/example?q=foo&r=bar`
  - Key: `/request/path`
  - Access: Read-Only
  - Returned Value: `/example?q=foo&r=bar`
  - Comment: That would provide read-only access to the request URL path.
- Read an specific URL parameter.
  - Scenario: Request URL is `http://localhost:8080/example?q=foo&r=bar`
  - Key: `/request/params/q`
  - Access: Read-Only
  - Returned Value: `foo`
  - Comment: That would provide read-only access to the request URL parameter `q`.
- Obtain the `Content-Type` header of the request.
  - Scenario: A POST request with a JSON body and the header `Content-Type` set to `application/json`.
  - Key: `/request/headers/Content-Type`
  - Access: Read-Only
  - Returned Value: `application/json`
  - Comment: That would provide read-only access to the value of the request header `Content-Type`.
- Read a field from a form.
  - Scenario: A request generated by submitting this form:<br />
    ```html
    <form method="post">
      First name:<br>
      <input type="text" name="firstname" value="Jane"><br>
      Last name:<br>
      <input type="text" name="lastname" value="Doe">
      <input type="submit" value="Submit">
    </form>
    ```
  - Key: `/request/form/firstname`
  - Access: Read-Only
  - Returned Value: `Jane`
  - Comment: That would provide read-only access to the value of the field `firstname` of the form.
- Set the response status code.
  - Scenario: A request is being attended.
  - Key: `/response/status`
  - Access: Write-Only
  - Acceptable Value: A 3-digit integer.  Must match `[0-9]{3}`.
  - Default Value: `200`
  - Comment: It is customary to use the HTTP status code as defined at [https://www.w3.org/Protocols/rfc2616/rfc2616-sec6.html#sec6.1.1](RFC2616).
- Set the response body.
  - Scenario: A request is being attended.
  - Key: `/response/body`
  - Access: Write-Only
  - Acceptable Value: Any string of bytes.
  - Default Value: N/A
  - Comment: For media types other than `application/octet-stream` you should specify the appropiate `Content-Type` header.

**Note**: Parameters under `request` are read-only and, conversely, parameters under
`response` are write-only.


#### Get handler resource

Returns the value of the requested resource path, or an error if the resource path doesn't exist or is invalid.

* **URL**: `/handlers/{:handler_id}{:resource_path}`
* **Method**: `GET`
* **URL Params**: FIXME: We think that here should be options to cook the value in some way, or get it raw.
* **Success Responses**:
  * **Code**: `200 OK`<br />
    **Header**: `Content-Type: application/octet-stream`<br />
    **Content**: The value of the resource.  Note that it may be empty.
* **Error Responses**:
    **Code**: `400 Invalid Resource Path`<br />
    **Notes**: Check the list of valid resource paths at the top of this section.
  * **Code**: `404 Not Found`
* **Sample Call**:<br />
  ```sh
  $ curl /handlers/$KAPOW_HANDLER_ID/request/body
  ```
* **Notes**: The content may be empty.


#### Overwrite the value of a resource

* **URL**: `/handlers/{:handler_id}{:resource_path}`
* **Method**: `PUT`
* **URL Params**: FIXME: We think that here should be options to cook the value in some way, or pass it raw.
* **Data Params**: Binary payload.
* **Success Responses**:
  * **Code**: `200 OK`
* **Error Responses**:
  * **Code**: `400 Invalid Payload`
  * **Code**: `400 Invalid Resource Path`<br />
    **Notes**: Check the list of valid resource paths at the top of this section.
  * **Code**: `404 Handler Not Found`
  * **Code**: `404 Name Not Found`<br />
    **Notes**: Although the resource path is correct, no such name is present in the request.  For instance, `/request/headers/Foo`, when no `Foo` header is not present in the request.
* **Sample Call**:<br />
  ```sh
  $ curl -X --data-binary '<h1>Hello!</h1>' PUT /handlers/$KAPOW_HANDLER_ID/response/body
  ```
* **Notes**:


## Usage Example

TODO: End-to-end example of the data API.


## Test Suite Notes

The test suite is located on [blebleble] directory.
You can run it by ...


# Framework


## Commands

Any compliant implementation of Kapow! must provide these commands:


### `kapow server`

This is the master command, that shows the help if invoked without args, and
runs the sub-commands when provided to it.


#### Example
```sh
$ kapow
Usage: kapow [OPTIONS] COMMAND [ARGS]...

Options:
  TBD

Commands:
  start    starts a Kapow! server
  route    operates on routes
...
```


### `kapow start`

This command runs the Kapow! server, which is the core of Kapow!.  If
run without parameters, it will run an unconfigured server.  It can accept a path
to a `.pow` file, which is a shell script that contains commands to configure
the Kapow! server.

The `.pow` can leverage the `kapow route` command, which is used to define a route.
The `kapow route` command needs a way to reach the Kapow! server, and for that,
`kapow` provides the `KAPOW_URL` variable in the environment of the
aforementioned shell script.

Every time the kapow server receives a request, it will spawn a process to
handle it, according to the specified entrypoint, `/bin/sh -c` by default, and then
execute the specified command.  This command is tasked with processing the
incoming request, and can leverage the `request` and `response` commands to
easily access the `HTTP Request` and `HTTP Response`, respectively.

In order for `request` and `response` to do their job, they require a way to
reach the Kapow! server, as well as a way to identify the current request being
served.  Thus, the Kapow! server adds the `KAPOW_URL` and `KAPOW_HANDLER_ID` to the
process' environment.


#### Example
```sh
$ kapow start /path/to/service.pow
```


### `kapow route`

To serve an endpoint, you must first register it.

`kapow route` registers/deregisters a route, which maps an
`HTTP` method and a URL pattern to the code that will handle the request.

When registering, you can specify an *entrypoint*, which defaults to `/bin/sh -c`,
and an argument to it, the *command*.

To deregister a route you must provide a *route_id*.

**Notes**:
 * The entrypoint definition matches *Docker*'s.
 * The index matches the way *netfilter*'s `iptables` handles rule numbering.


#### **Environment**
- `KAPOW_URL`


#### **Help**
```sh
$ kapow route --help
Usage: kapow route [OPTIONS] COMMAND [ARGS]...

Options:
  --help  Show this message and exit.

Commands:
  add
  remove
```
```sh
$ kapow route add --help
Usage: kapow route add [OPTIONS] URL_PATTERN [COMMAND_FILE]

Options:
  -c, --command TEXT
  -e, --entrypoint TEXT
  -X, --method TEXT
  --url TEXT
  --help                 Show this message and exit.
```
```sh
$ kapow route remove --help
Usage: kapow route remove [OPTIONS] ROUTE_ID

Options:
  --url TEXT
  --help      Show this message and exit.
```


#### Example
```sh
kroute add -X GET '/list/{ip}' -c 'nmap -sL $(request /matches/ip) | response /body'
```

### `request`

Exposes the requests' resources.


#### **Environment**
- `KAPOW_URL`
- `KAPOW_HANDLER_ID`


#### Example
```sh
# Access the body of the request
request /body
```


### `response`

Exposes the response's resources.


#### **Environment**
- `KAPOW_URL`
- `KAPOW_HANDLER_ID`


#### Example
```sh
# Write to the body of the response
echo 'Hello, World!' | response /body
```


## An End-to-End Example
```sh
$ cat nmap.kpow
kroute add -X GET '/list/{ip}' -c 'nmap -sL $(request /matches/ip) | response /body'
```
```sh
$ kapow ./nmap.kapow
```
```sh
$ curl $KAPOW_URL/list/127.0.0.1
Starting Nmap 7.70 ( https://nmap.org ) at 2019-05-30 14:45 CEST
Nmap scan report for localhost (127.0.0.1)
Host is up (0.00011s latency).
Not shown: 999 closed ports
PORT   STATE SERVICE
22/tcp open  ssh

Nmap done: 1 IP address (1 host up) scanned in 0.06 seconds
```


## Test Suite Notes


# Server


## Test Suite Notes
