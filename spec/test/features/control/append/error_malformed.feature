Feature: Kapow! server reject append requests with malformed JSON bodies.
  Kapow! server will reject to append a route when
  it receives a malformed json document in the
  request body.

  Scenario: Error because a malformed JSON document.
    If a request comes with an invalid JSON document
    the server will respond with a bad request error.

    Given I have a running Kapow! server
    When I try to append with this JSON document:
      """
      {
        "method" "GET",
        "url_pattern": /hello,
        "entrypoint": null
        "command": "echo Hello
         World | response /body",
        "id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
      }
      """
    Then I get 400 as response code
      And I get "Malformed JSON" as response phrase
      And I get an empty response body
