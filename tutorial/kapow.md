# WELCOME <!-- .slide: data-transition="zoom-out" -->

### To the awesome

<img class="plain" style="width: 40%; background: none; border: none; box-shadow: none; margin: 0px" src="/tour/assets/logo.png" />

## GUIDED TOUR!

---

# Getting started


--

## Hello World!  <!-- .slide: data-state="clearshell" -->

- HTTP endpoints can be added on-the-fly with the `kapow` command line utility. [Check it out!](/tour/terminal/sendkeys?s=kapow+route+add+-X+GET+%27%2Fhelloworld%27+-c+%27echo+%22Hello+World%21%22+%7C+response+%2Fbody%27)

    ```bash
    $ kapow route add \
        -X GET '/helloworld' \
        -c 'echo "Hello World!" | response /body'

    ```

- New endpoints are inmediately available. [Check it out!](/tour/terminal/sendkeys?s=curl%20http%3A%2F%2Flocalhost%3A8080%2Fhelloworld)

    ```bash
    $ curl http://localhost:8080/helloworld
    ```
    
