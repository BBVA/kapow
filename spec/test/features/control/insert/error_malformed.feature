Feature: Kapow! server reject insert requests with malformed JSON bodies.
  Kapow! server will reject to insert a route when
  it receives a malformed json document in the
  request body.

  Scenario: Error because a malformed JSON document.
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
         World | response /body",
        "index": 0,
        "id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
      }
      """
    Then I get bad request as response code
      And I get "Malformed JSON" as response phrase
      And I get an empty response body
