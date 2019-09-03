Feature: Consistent route ordering after inserting a route in a Kapow! server.
  When inserting routes the server will mantain the
  whole set of routes ordered and with consecutive indexes.

  Background:
    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |

  Scenario: Inserting before the first route.
    After inserting before the first route the previous set
    will maintain their relative order and their indexes
    will be increased by one.

    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "/bar",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | response /body",
        "index": 0
      }
      """
      And I request a routes listing
    Then I get the following response body:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/bar",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | response /body",
            "index": 0,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/foo",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | response /body",
            "index": 1,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/qux/{dirname}",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | response /body",
            "index": 2,
            "id": ANY
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
        "url_pattern": "/bar",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | response /body",
        "index": 2
      }
      """
      And I request a routes listing
    Then I get the following response body:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/foo",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | response /body",
            "index": 0,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/qux/{dirname}",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | response /body",
            "index": 1,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/bar",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | response /body",
            "index": 2,
            "id": ANY
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
        "url_pattern": "/bar",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | response /body",
        "index": 1
      }
      """
      And I request a routes listing
    Then I get the following response body:
        """
        [
          {
            "method": "GET",
            "url_pattern": "/foo",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | response /body",
            "index": 0,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/bar",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | response /body",
            "index": 1,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/qux/{dirname}",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | response /body",
            "index": 2,
            "id": ANY
          }
        ]
        """
