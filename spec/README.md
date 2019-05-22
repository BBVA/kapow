# Kapow!

## Why?

Because we think that:

- UNIX is great and we love it
- The UNIX shell is great
- HTTP interfaces are convenient and everywhere
- CGI is not a good way to mix them


## How?

So, how we can mix the **web** and the **shell**? Let's see...

The **web** and the **shell** are two different beasts, both packed with
history.

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
  |              | HTTP Methods           | Exit Codes              |
  +--------------+------------------------+-------------------------+
```

Any tool designed to give an HTTP interface to an existing shell command
**must map concepts of boths**.  For example:

- "GET parameters" to "Command line parameters"
- "Headers" to "Environment variables"
- "Stdout" to "Response body"

Kapow! is not opinionated about the different ways you can map both worlds.
Instead it provides a concise set of tools used to express the mapping and a
set of common defaults.


### Why not tool "X"?

All the alternatives we found are **rigid** about how they match between HTTP
and shell concepts.

* [shell2http](https://github.com/msoap/shell2http): HTTP-server to execute
  shell commands.  Designed for development, prototyping or remote control.
Settings through two command line arguments, path and shell command.
* [websocketd](https://github.com/joewalnes/websocketd): Turn any program that
  uses STDIN/STDOUT into a WebSocket server.  Like inetd, but for WebSockets.
* [webhook](https://github.com/adnanh/webhook): webhook is a lightweight
  incoming webhook server to run shell commands.
* [gotty](https://github.com/yudai/gotty): GoTTY is a simple command line tool
  that turns your CLI tools into web applications.  (For interactive commands
only)
* [shell-microservice-exposer](https://github.com/jaimevalero/shell-microservice-exposer):
  Expose your own scripts as a cool microservice API dockerizing it.

Tools with a rigid matching **can't evade** *[impedance
mismatch](https://haacked.com/archive/2004/06/15/impedance-mismatch.aspx/)*.
Resulting in an easy-to-use software, convenient in some scenarios but
incapable in others.


### Why not my good-old programming language "X"?

* Boilerplate
* Custom code = More bugs
* Security issues (Command injection, etc)
* Dependency on developers
* **"A programming language is low level when its programs require attention to
  the irrelevant"** *Alan Perlis*
* **There is more Unix-nature in one line of shell script than there is in ten
  thousand lines of C** *Master Foo*


### Why not CGI?

TODO: Small explanation and example.


## What?

We named it Kapow!.  It is pronounceable, short and meaningless...  like every
good UNIX command ;-)

TODO: Definition
TODO: Intro to Architecture


# HTTP API

Kapow! server interacts with the outside world only through its HTTP API.  Any
program making the correct HTTP request to a Kapow! server, can change its
behavior.

## Design Principles

* All requests and responses will leverage JSON as the data encoding method.

* The API calls responses will have two distinct parts:
  * The HTTP status code (e.g., 400, Bad Request).  The target audience of this
    information is the client code.  The client can thus use this information to
    control the program flow.
  * The JSON-encoded message.  The target audience in this case is the human
    operating the client.  The human can use this information to make a
    decision on how to proceed.

Let's illustrate these ideas with an example: TODO

  * All successful API calls will return a representation the *final* state
    attained by the objects which have been addressed (requested, set or
    deleted).

    FIXME: consider what to do when deleting objects.  Isn't it too much to
    return the list of all deleted objects in such a request?

## API elements

### Servers

TODO: Define servers' API

### Routes

Routes are the mechanism that allows Kapow! to find the correct program to
respond to an external event (e.g.  an incomming HTTP request).

#### List routes

Returns JSON data about the current routes.

* **URL**

  `/routes`

* **Method**

  `GET`

* **Success Response**

  * **Code**: `200 OK`<br />
    **Content**: TODO

* **Sample Call**

  TODO

* **Notes**

  Currently all routes are returned; in the future, a filter may be accepted.

#### Append route

  Accepts JSON data that defines a new route to be appended to the current routes.

* **URL**

  `/routes`

* **Method**

  `POST`

* **Header**: `Content-Type: application/json`

* **Data Params**

  * **Content**:
  ```
  {
    "method": "GET",
    "url_pattern": "/hello",
    "entrypoint": null,
    "command": "echo Hello World | response /body"
  }
  ```

* **Success Response**

  * **Code**: `200 OK`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**:
  ```
  {
    "method": "GET",
    "url_pattern": "/hello",
    "entrypoint": null,
    "command": "echo Hello World | response /body",
    "index": 0
  }
  ```

* **Error Response**

  * **Code**: `400 Bad Request`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**: `{ "error": "Malformed JSON." }`

  * **Code**: `400 Bad Request`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**: `{ "error": "Mandatory field(s) not provided." }`

* **Sample Call**
TODO

* **Notes**

  * A successful request will yield a response containing all the effective
    parameters that were applied.

#### Insert a route

  Accepts JSON data that defines a new route to be inserted at the specified
  index to the current routes.

* **URL**

  `/routes`

* **Method**

  `PUT`

* **Header**: `Content-Type: application/json`

* **Data Params**

  * **Content**:
  ```
  {
    "method": "GET",
    "url_pattern": "/hello",
    "entrypoint": null,
    "command": "echo Hello World | response /body",
  }
  ```

* **Success Response**

  * **Code**: `200 OK`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**:
  ```
  {
    "method": "GET",
    "url_pattern": "/hello",
    "entrypoint": null,
    "command": "echo Hello World | response /body",
    "index": 0
  }
  ```

* **Error Response**

  * **Code**: `400 Bad Request`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**: `{ "error": "Malformed JSON." }`

  * **Code**: `400 Bad Request`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**: `{ "error": "Mandatory field(s) not provided." }`

* **Sample Call**
TODO

* **Notes**

  * Route numbering starts at zero.
  * When `index` is not provided or is less than 0 the route will be inserted
    first, effectively making it index 0.
  * Conversely when `index` is greater than the number of entries on the route
    table it will be inserted last.
  * A successful request will yield a response containing all the effective
    parameters that were applied.

#### Delete a route

Removes the route identified by `:id`.

* **URL**

  `/routes/:id`

* **Method**

  `DELETE`

* **Success Response**

  * **Code**: `200 OK`<br />
    **Content**: TODO

* **Error Response**

  * **Code**: `404 Not Found`<br />
    **Header**: `Content-Type: application/json`<br />
    **Content**: `{ "error": "Unknown route", "route_id": "{{ :id }}" }`


* **Sample Call**
TODO

* **Notes**


### Handlers

Handlers are in-memory data structures exposing the data of the current request
and response.

Each handler is identified by a `handler_id` and provide access to the
following keys:

```
/                        The root of the keys tree
│
├─ request               All information related to the HTTP request.  Read-Only
│  ├──── method          Used HTTP Method (GET, POST)
│  ├──── path            Complete URL path (URL-unquoted)
│  ├──── matches         Previously matched URL path parts
│  │     └──── <key>
│  ├──── params          URL parameters (post ? symbol)
│  │     └──── <key>
│  ├──── headers         HTTP request headers
│  │     └──── <key>
│  ├──── cookies         HTTP request cookie
│  │     └──── <key>
│  ├──── form            form-urlencoded form fields
│  │     └──── <key>
│  └──── body            HTTP request body
│  
└─ response              All information related to the HTTP request.  Write-Only
   ├──── status          HTTP status code
   ├──── body            Response body.  Mutually exclusive with response/stream
   ├──── stream          Chunk-encoded body.  Streamed response.  Mutually exclusive with response/body
   └──── headers         HTTP response headers
         └──── <key>
```

#### Example Keys

1.
  Scenario: Request URL is `http://localhost:8080/example?q=foo&r=bar`
  Key: `/request/path`
  Access: Read-Only
  Returned Value: `/example?q=foo&r=bar`
  Comment: That would provide read-only access to the request URL path.

2. 
  Scenario: Request URL is `http://localhost:8080/example?q=foo&r=bar`
  Key: `/request/params/q`
  Access: Read-Only
  Returned Value: `foo`
  Comment: That would provide read-only access to the request URL parameter `q`.

3.
  Scenario: A POST request with a JSON body and the header `Content-Type` set to `application/json`.
  Key: `/request/headers/Content-Type`
  Access: Read-Only
  Returned Value: `application/json`
  Comment: That would provide read-only access to the value of the request header `Content-Type`.

4.
  Scenario: A request generated by submitting this form:
    ```
    <form method="post">
      First name:<br>
      <input type="text" name="firstname" value="Jane"><br>
      Last name:<br>
      <input type="text" name="lastname" value="Doe">
      <input type="submit" value="Submit">
    </form>
    ```
  Key: `/request/form/firstname`
  Access: Read-Only
  Returned Value: `Jane`
  Comment: That would provide read-only access to the value of the field `firstname` of the form.

5.
  Scenario: A request is being attended.
  Key: `/response/status`
  Access: Write-Only
  Acceptable Value: A 3-digit integer.  Must match `[0-9]{3}`.
  Default Value: 200
  Comment: It is customary to use the HTTP status code as defined at [https://www.w3.org/Protocols/rfc2616/rfc2616-sec6.html#sec6.1.1](RFC2616).

6.
  Scenario: A request is being attended.
  Key: `/response/body`
  Access: Write-Only
  Acceptable Value: Any string of bytes.
  Default Value: N/A
  Comment: For media types other than `application/octet-stream` you should
  specify the appropiate `Content-Type` header.


**Note**: Parameters under `request` are read-only and, conversely, parameters under
`response` are write-only.

#### Get handler key

* **URL**
* **Method**
* **URL Params**
* **Data Params**
* **Success Response**
* **Error Response**
* **Sample Call**
* **Notes**

#### Overwrite the value for a handler key
* **URL**
* **Method**
POST
* **URL Params**
* **Data Params**
* **Success Response**
* **Error Response**
* **Sample Call**
* **Notes**

## Usage Example

## Test Suite Notes

The test suite is located on [blebleble] directory.
You can run it by ...


# Framework
## Commands
Any compliant implementation of Kapow! must provide these commands:

### `kapow`
This implements the server, yaddayadda

#### Example

### `kroute`
TODISCUSS: maybe consider using `kapow route` instead

#### Example

### `request`

#### Example

### `response`

#### Example

## Full-fledged example  (TODO: express it more simply)


## Test Suite Notes

# Server


## Test Suite Notes
