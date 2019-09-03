Feature: Setting values for handler response resources in Kapow! server.
    Users can set the values in the response
    resources by specifying the handler id
    and the resource path.

  Scenario: Set status code for the current response.
    Set the status code through the current
    handler.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern  |
      | GET    | /foo |
    When I send a request to the testing route "/foo"
      And I set the resource "/response/status" with value "418"
      And I release the testing request
    Then I get 418 as response code in the testing request

  Scenario Outline: Set different resources for the current response.
    Set the following resources for the current
    response through the current handler.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern  |
      | GET    | /foo |
    When I send a request to the testing route "/foo"
      And I set the resource "<resourcePath>" with value "<value>"
      And I release the testing request
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the value "<value>" for the response "<fieldType>" named "<elementName>" in the testing request

    Examples:
      | resourcePath            | value      | fieldType | elementName |
      | /response/headers/head1 | headVal1   | header    | head1       |
      | /response/cookies/cook1 | cookVal1   | cookie    | cook1       |
      | /response/body          | bodyValue1 | body      | -           |
      | /response/stream        | bodyValue2 | body      | -           |
