Feature: Kapow! server rejects requests with semantic errors.
  Kapow! server will refuse to append routes when
  it receives a valid json document not conforming
  to the specification.

  Scenario: Error because lacking mandatory fields.
    If a request lacks any mandatory field the server
    responds with an error.

    Given I have a running Kapow! server
    When I append the route:
      """
      {
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body"
      }
      """
    Then I get 422 as response code
      And I get "Invalid Route" as response reason phrase

  Scenario: Error because bad route format.
    If a request contains an invalid expression in the
    field url_pattern the server responds with an error.

    Given I have a running Kapow! server
    When I append the route:
      """
      {
        "method": "GET",
        "url_pattern": "+123--",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body"
      }
      """
    Then I get 422 as response code
      And I get "Invalid Route" as response reason phrase
