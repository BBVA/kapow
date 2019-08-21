Feature: Routes auto-ordering after inserting in a Kapow! server.

  When inserting routes the server will mantain the
  whole set of routes ordered an with consecutive indexes.

  Background:
    Given I have a Kapow! server whith the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |

  Scenario: Inserting before the first route.
    After inserting before the first route the previous set
    will maintain their relative order and their indexes
    will be increased by one.

    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "/listVarDir",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | response /body",
        "index": 0
      }
      """
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get the following entity as response body:
        """
        {
          "method": "GET",
          "url_pattern": "/listVarDir",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /var | response /body",
          "index": 0,
          "id": "*"
        }
        """
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get a list with the following elements:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/listVarDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | response /body",
            "index": 0,
            "id": "*"
          },
          {
            "method": "GET",
            "url_pattern": "/listRootDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | response /body",
            "index": 1,
            "id": "*"
          },
          {
            "method": "GET",
            "url_pattern": "/listDir/:dirname",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | response /body",
            "index": 2,
            "id": "*"
          }
        ]
        """

  Scenario: Inserting after the last routes.
    After inserting after the last route the previous set
    will maintain their relative order and indexes.

    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "/listVarDir",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | response /body",
        "index": 2
      }
      """
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get the following entity as response body:
        """
        {
          "method": "GET",
          "url_pattern": "/listVarDir",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /var | response /body",
          "index": 2,
          "id": "*"
        }
        """
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get a list with the following elements:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/listRootDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | response /body",
            "index": 0,
            "id": "*"
          },
          {
            "method": "GET",
            "url_pattern": "/listDir/:dirname",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | response /body",
            "index": 1,
            "id": "*"
          },
          {
            "method": "GET",
            "url_pattern": "/listVarDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | response /body",
            "index": 2,
            "id": "*"
          }
        ]
        """

  Scenario: Inserting a midst route.
    After inserting a midst route the previous route set
    will maintain their relative order and the indexes
    of thefollowing routes will be increased by one.

    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "/listVarDir",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | response /body",
        "index": 1
      }
      """
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get the following entity as response body:
        """
        {
          "method": "GET",
          "url_pattern": "/listVarDir",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /var | response /body",
          "index": 1,
          "id": "*"
        }
        """
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get a list with the following elements:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/listRootDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | response /body",
            "index": 0,
            "id": "*"
          },
          {
            "method": "GET",
            "url_pattern": "/listVarDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | response /body",
            "index": 1,
            "id": "*"
          },
          {
            "method": "GET",
            "url_pattern": "/listDir/:dirname",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | response /body",
            "index": 2,
            "id": "*"
          }
        ]
        """
