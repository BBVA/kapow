Feature: Setting response values for handler resources in Kapow! server.
    Users can set the values in the response
    resources by specifying the handler id
    and the resource path.

  Scenario: Set status code for the current response.
    Set the status code for the current
    response through the current handler.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern  | entrypoint | command                      |
      | GET    | /listRootDir | /bin/sh -c | echo 666 \| response /status |
    When I send a request to the route "/listRootDir"
    Then I get 666 as response code
      And I get "OK" as response reason phrase

  Scenario Outline: Set all defined resources for the current response.
    Set the following resources for the current
    response through the current handler.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern  | entrypoint | command                                 |
      | GET    | /listRootDir | /bin/sh -c | echo <value> \| response <resourcePath> |
    When I send a request to the route "/listRootDir"
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the value <value> for the response <type> named <elementName>

    Examples:
      | resourcePath   | value        | type   | elementName |
      | /headers/head1 | "headVal1"   | header | "head1"     |
      | /cookies/cook1 | "cookVal1"   | cookie | "cook1"     |
      | /body          | "bodyValue1" | body   | ""          |
      | /stream        | "bodyValue2" | body   | ""          |
