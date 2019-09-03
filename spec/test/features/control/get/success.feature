Feature: Retrieve route details in Kapow! server.
    Users can retrieve route details from the server
    by specifying its id.

  Scenario: Retrieve route details.
    Get route details by id.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /foo       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |
    When I get the first route
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response body:
        """
        {
          "method": "GET",
          "url_pattern": "/foo",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la / | response /body",
          "index": 0,
          "id": ANY
        }
        """
