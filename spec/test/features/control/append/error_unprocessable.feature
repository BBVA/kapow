Feature: Kapow! server reject responses with semantic errors.
  Kapow! server will reject to append routes when
  it receives a valid json document but not conforming
  with the specification.

  Scenario: Error because of lack of mandatory fields.
    If a request lacks of any of the mandatory fields
    the server responds with an error indicating the
    missing fields.

    Given I have a running Kapow! server
    When I append the route:
      """
      {
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body"
      }
      """
    Then I get 422 as response code
      And I get "Missing Mandatory Field" as response reason phrase
      And I get the following entity as response body:
        """
        {
          "missing_mandatory_fields": [
            "url_pattern",
            "method"
          ]
        }
        """

  Scenario: Error because of wrong route specification.
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
      And I get "Invalid Route Spec" as response reason phrase
      And I get an empty response body

  Scenario: Error because of wrong method value.
    If a request contains an invalid value in the
    field method the server responds with an error.

    Given I have a running Kapow! server
    When I append the route:
      """
      {
        "method": "SOMETIMES",
        "url_pattern": "/",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body"
      }
      """
    Then I get 422 as response code
      And I get "Invalid Data Type" as response reason phrase
      And I get an empty response body
