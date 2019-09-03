Feature: Retrieve request resources from a handler in Kapow! server.
    Users can retrieve request resources by
    specifying the handler id and the
    resource path.

  Scenario Outline: Retrieve different resources for the current request.
    Get the following resources for the
    current request through the current
    handler.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern         |
      | GET    | /foo/{path} |
    When I send a request to the testing route "/foo/matchVal1" adding:
      | fieldType | name  | value      |
      | parameter | par1  | paramVal1  |
      | header    | head1 | headVal1   |
      | cookie    | cook1 | cookieVal1 |
      | body      |       | bodyVal1   |
      And I get the resource "<resourcePath>"
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response raw body:
        """
        <value>
        """

    Examples:
      | resourcePath           | value                  |
      | /request/method        | GET                    |
      | /request/path          | /foo/matchVal1 |
      | /request/host          | localhost:8080         |
      | /request/matches/path  | matchVal1              |
      | /request/params/par1   | paramVal1              |
      | /request/headers/head1 | headVal1               |
      | /request/cookies/cook1 | cookieVal1             |
      | /request/body          | bodyVal1               |
