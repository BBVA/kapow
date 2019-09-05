Feature: Consistent route order after a route deletion in Kapow! server.
  When deleting routes the server will mantain the
  remaining routes ordered and with consecutive indexes.

  Background:
    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| kapow set /response/body                       |
      | GET    | /bar           | /bin/sh -c | ls -la /var \| kapow set /response/body                    |
      | GET    | /baz           | /bin/sh -c | ls -la /etc \| kapow set /response/body                    |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| kapow set /response/body |

  Scenario: Removing the first route.
    After removing the first route the remaining ones
    will maintain their relative order and their indexes
    will be decreased by one.

    When I delete the first route
      And I request a routes listing
    Then I get the following response body:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/bar",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | kapow set /response/body",
            "index": 0,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/baz",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /etc | kapow set /response/body",
            "index": 1,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/qux/{dirname}",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | kapow set /response/body",
            "index": 2,
            "id": ANY
          }
        ]
        """

  Scenario: Removing the last route.
    After removing the last route the remaining ones will
    maintain their relative order and indexes.

    When I delete the last route
      And I request a routes listing
    Then I get the following response body:
      """
      [
        {
          "method": "GET",
          "url_pattern": "/foo",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la / | kapow set /response/body",
          "index": 0,
          "id": ANY
        },
        {
          "method": "GET",
          "url_pattern": "/bar",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /var | kapow set /response/body",
          "index": 1,
          "id": ANY
        },
        {
          "method": "GET",
          "url_pattern": "/baz",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /etc | kapow set /response/body",
          "index": 2,
          "id": ANY
        }
      ]
      """

  Scenario: Removing a midst route.
    After removing a midst route the remaining ones will
    maintain their relative order and the indexes of the
    following routes will be decreased by one.

    When I delete the second route
      And I request a routes listing
    Then I get the following response body:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/foo",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | kapow set /response/body",
            "index": 0,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/baz",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /etc | kapow set /response/body",
            "index": 1,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/qux/{dirname}",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | kapow set /response/body",
            "index": 2,
            "id": ANY
          }
        ]
        """
