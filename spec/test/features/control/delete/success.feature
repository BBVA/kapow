Feature: Delete routes in Kapow! server.
  Deleting routes allows users to remove undesired
  routes from the server.

  Scenario: Delete a route.
    Routes are removed from the sever by specifying their id.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |
    When I delete the first route
    Then I get 204 as response code
      And I get "No Content" as response reason phrase
