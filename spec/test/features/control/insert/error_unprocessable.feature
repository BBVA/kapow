Feature: Kapow! server reject insert responses with semantic errors.
  Kapow! server will reject to insert routes when
  it receives a valid json but not conforming document.

  Scenario: Error because of lack of mandatory fields.
    If a request lacks of any of the mandatory fields
    the server responds with an error indicating the
    missing fields.

    Given I have a running Kapow! server
    When I insert the route:
      | entrypoint | command                    |
      | /bin/sh -c | ls -la / \| response /body |
    Then I get 422 as response code
      And I get "Missing Mandatory Field" as response reason phrase
      And I get the following entity as response body:
        | missing_mandatory_fields |
        | "url_pattern", "method" |

  Scenario: Error because of wrong route specification.
    If a request contains an invalid expression in the
    url_pattern field the server responds with an error.

    Given I have a running Kapow! server
    When I insert the route:
      | method | url_pattern  | entrypoint | command                    | index |
      | GET    | /listRootDir | /bin/sh -c | ls -la / \| response /body |     0 |
    Then I get 422 as response code
      And I get "Invalid Route Spec" as response reason phrase
      And I get an empty response body

  Scenario: Error because of wrong method value.
    If a request contains an invalid value in the
    method field the server responds with an error.

    Given I have a running Kapow! server
    When I insert the route:
      | method | url_pattern  | entrypoint | command                    | index |
      | AVECES | /listRootDir | /bin/sh -c | ls -la / \| response /body |     0 |
    Then I get 422 as response code
      And I get "Invalid Data Type" as response reason phrase
      And I get an empty response body
