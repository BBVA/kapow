@wip
Feature: Fail to retrieve nonexistent resource items in Kapow! server.
  If trying to access a nonexistent resource item
  then the server responds with a no content error.

  Scenario: Try to get a nonexistent resource item from a handler.
    A request to retrieve a nonexistent resource
    item from a handler will trigger a no content
    error.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern        |
      | GET    | /listRootDir       |
    When I send a request to the testing route "/listRootDir"
      And I get the resource "/request/params/meloinvento"
    Then I get 204 as response code
      And I get "Resource Item Not Found" as response reason phrase
