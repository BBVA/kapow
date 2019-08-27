@wip
Feature: Fail to retrieve an invalid resource for a handler in Kapow! server.
  If trying to access an invalid resource for a handler
  then the server responds with an error.

  Scenario: Try to get invented/path from a existent handler.
    A request for retrieving an invalid resource for an
    existent handler will trigger a invalid resource error.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern        | entrypoint | command                    |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body |
    When I send a request to the route "/listRootDir"
      And I get the resource "request/path" for the current request handler
    Then I get 400 as response code
      And I get "Invalid Resource Path" as response reason phrase
