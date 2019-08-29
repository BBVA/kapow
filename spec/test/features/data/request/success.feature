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
      | GET    | /listRootDir/{path} |
    When I send a request to the testing route "/listRootDir/otro" adding:
      | fieldType | name    | value        |
      | parameter | "par1"  | "paramVal1"  |
      | header    | "head1" | "headVal1"   |
      | cookie    | "cook1" | "cookieVal1" |
      And I get the resource <resourcePath>
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response raw body:
        """
        <value>
        """

    Examples:
      | resourcePath           | value               |
      | resource/method        | "GET"               |
      | resource/path          | "/listRootDir/otro" |
      | resource/host          | "localhost:8080"    |
      | resource/matches/path  | "otro"              |
      | resource/params/par1   | "paramVal1"         |
      | resource/headers/head1 | "headVal1"          |
      | resource/cookies/cook1 | "cookieVal1"        |
      | resource/body          | empty               |
