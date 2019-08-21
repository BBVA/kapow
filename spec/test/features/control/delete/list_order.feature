Feature: Routes auto-ordering after deleting  in a Kapow! server.

  When deleting routes the server will mantain the
  remaining routes ordered an with consecutive indexes.

  Background:
    Given I have a Kapow! server whith the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listVarDir        | /bin/sh -c | ls -la /var \| response /body                    |
      | GET    | /listEtcDir        | /bin/sh -c | ls -la /etc \| response /body                    |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |

  Scenario: Removing the first routes.

    After removing the first route the remaining ones
    will maintain their relative order and their indexes
    will be decreased by one.

    When I delete the first route inserted
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get an empty response body
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
            "url_pattern": "/listEtcDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /etc | response /body",
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

  Scenario: Removing the last routes.
    After removing the last route the remaining ones will
    maintain their relative order and indexes.

    When I delete the last route inserted
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get an empty response body
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
          "url_pattern": "/listEtcDir",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /etc | response /body",
          "index": 2,
          "id": "*"
        }
      ]
      """

  Scenario: Removing a midst route.
    After removing a midst route the remaining ones will
    maintain their relative order and the indexes of the
    following routes will be decreased by one.

    When I delete the second route inserted
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get an empty response body
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
            "url_pattern": "/listEtcDir",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /etc | response /body",
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
