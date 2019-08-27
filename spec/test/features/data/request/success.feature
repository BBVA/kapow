Feature: Retrieve request resources from a handler in Kapow! server.
    Users can retrieve request handler resources
    from the server by specifying the handler id
    and the resource path.

  Scenario Outline: Retrieve all defined resources for the current request.
    Get the following resources for the current
    request through the current handler.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern         | entrypoint | command                    |
      | GET    | /listRootDir/{path} | /bin/sh -c | ls -la / \| response /body |

    When I send a request to the route "/listRootDir/otro" setting this values:
      | type      | name    | value        |
      | parameter | "par1"  | "paramVal1"  |
      | header    | "head1" | "headVal1"   |
      | cookie    | "cook1" | "cookieVal1" |

      And I get the resource <resourcePath> for the current request handler
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response body <value>

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
