Feature: Fail to retrieve an invalid resource for a handler in Kapow! server.
  If trying to access an invalid resource for a handler
  then the server responds with an error.

  Scenario: Try to get an inexistent resource from a handler.
    A request to retrieve an invalid resource
    from a handler will trigger a invalid resource error.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern        |
      | GET    | /listRootDir       |
    When I send a request to the testing route "/listRootDir"
      And I get the resource "invented/path"
    Then I get 400 as response code
      And I get "Invalid Resource Path" as response reason phrase
