@wip
Feature: Fail to retrieve resources from nonexistent handler in Kapow! server.
  If trying to access a nonexistent handler then the
  server responds with a noptfound error.

  Scenario: Try to get a valid resource from a nonexistent handler.
    A request for retrieving a resource from a nonexistent
    handler will trigger a not found error.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern        | entrypoint | command                    |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body |
    When I get the resource "request/path" for the handler with id XXXXXXXXXX
    Then I get 404 as response code
      And I get "Not Found" as response reason phrase

  Scenario: Fail to get an invalid resource from a nonexistent handler.
    A request for retrieving a resource from a nonexistent
    handler will trigger a not found error even if the
    resource is invalid.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern        | entrypoint | command                    |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body |
    When I get the resource "invented/path" for the handler with id XXXXXXXXXX
    Then I get 404 as response code
      And I get "Not Found" as response reason phrase
