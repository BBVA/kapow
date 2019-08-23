Feature: Kapow! server rejects insertion requests with semantic errors.
  Kapow! server will refuse to insert routes when
  it receives a valid JSON but not conforming document.

  Scenario: Error because lacking mandatory fields.
    If a request lacks any mandatory fields the server
    responds with an error.

    Given I have a running Kapow! server
    When I insert the route:
      """
      {
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body"
      }
      """
    Then I get 422 as response code
      And I get "Invalid Route" as response reason phrase

  Scenario: Error because wrong route specification.
    If a request contains an invalid expression in the
    url_pattern field the server responds with an error.

    Given I have a running Kapow! server
    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "+123--",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body",
        "index": 0
      }
      """
    Then I get 422 as response code
      And I get "Invalid Route" as response reason phrase
