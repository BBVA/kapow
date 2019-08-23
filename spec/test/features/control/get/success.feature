Feature: Retrieve route info in Kapow! server.
  Users can retrieve route info from the server by
  specifying its id.

  Scenario: Retrieve route info.
    Get route info by spscifying its id.

    Given I have a Kapow! server whith the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |
    When I get the first route info
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response body:
        """
        {
          "method": "GET",
          "url_pattern": "/listRootDir",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la / | response /body",
          "index": 0,
          "id": ANY
        }
        """
