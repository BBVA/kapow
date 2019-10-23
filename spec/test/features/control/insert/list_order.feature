#
# Copyright 2019 Banco Bilbao Vizcaya Argentaria, S.A.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
Feature: Consistent route ordering after inserting a route in a Kapow! server.
  When inserting routes the server will mantain the
  whole set of routes ordered and with consecutive indexes.

  Background:
    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| kapow set /response/body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| kapow set /response/body |

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
        "command": "ls -la /var | kapow set /response/body",
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
            "command": "ls -la /var | kapow set /response/body",
            "index": 0,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/foo",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la / | kapow set /response/body",
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

  Scenario: Inserting after the last routes.
    After inserting after the last route the previous set
    will maintain their relative order and indexes.

    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "/bar",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | kapow set /response/body",
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
            "command": "ls -la / | kapow set /response/body",
            "index": 0,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/qux/{dirname}",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | kapow set /response/body",
            "index": 1,
            "id": ANY
          },
          {
            "method": "GET",
            "url_pattern": "/bar",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /var | kapow set /response/body",
            "index": 2,
            "id": ANY
          }
        ]
        """

  Scenario: Inserting a route in the middle.
    After inserting a route in the middle, the previous
    route set will maintain their relative order and the indexes
    of the following routes will be increased by one.

    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "/bar",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /var | kapow set /response/body",
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
            "url_pattern": "/qux/{dirname}",
            "entrypoint": "/bin/sh -c",
            "command": "ls -la /request/params/dirname | kapow set /response/body",
            "index": 2,
            "id": ANY
          }
        ]
        """
