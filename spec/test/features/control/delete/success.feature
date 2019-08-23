Feature: Delete routes in Kapow! server.
  Delete routes allow users to remove non-desired
  routes from the server.

  Scenario: Delete a route.
    Routes are removed from the sever by specifying their id.

    Given I have a Kapow! server whith the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |
    When I delete the first route
    Then I get 200 as response code
      And I get "OK" as response reason phrase
