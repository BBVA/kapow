Feature: Kapow! server rejects insertion requests with malformed JSON bodies.
  Kapow! server will reject to insert a route when
  it receives a malformed JSON document in the
  request body.

  Scenario: Error because of malformed JSON document.
    If a request comes with an invalid JSON document
    the server will respond with a bad request error.

    Given I have a running Kapow! server
    When I try to insert with this JSON document:
      """
      {
        "method" "GET",
        "url_pattern": /hello,
        "entrypoint": null
        "command": "echo Hello
         World | kapow set /response/body",
        "index": 0,
        "id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
      }
      """
    Then I get 400 as response code
      And I get "Malformed JSON" as response reason phrase
