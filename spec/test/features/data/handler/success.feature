Feature: Retrieve a resource from a handler in Kapow! server.
    Users can retrieve request handler resources
    from the server by specifying the handler id
    and the resource path.

  Scenario: Retrieve a resource.
    Get the "request/path" resource for the current
    request through the handler id.

<<<<<<< HEAD
    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern        |
      | GET    | /listRootDir       |
    When I send a request to the testing route "/listRootDir"
      And I get the resource "request/path"
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response raw body:
        """
        /listRootDir
        """

  Scenario: Retrieve a resource item.
    Get the "request/headers/Host" item resource for
    the current request through the handler id.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern        |
      | GET    | /listRootDir       |
    When I send a request to the testing route "/listRootDir"
      And I get the resource "request/headers/Host"
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response raw body:
        """
        localhost:8080
        """
