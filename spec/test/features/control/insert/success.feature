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
Feature: Insert new routes in Kapow! server.
  Inserting routes allows users to configure the server. New
  routes could be inserted at the beginning or before any
  existing route of the routes list.

  Background:
    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| kapow set /response/body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| kapow set /response/body |

  Scenario: Insert a route at the beginning.
    A route can be inserted at the beginning of the list
    by specifying an index 0 in the request.

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
    Then I get 201 as response code
#      And I get "Created" as response reason phrase
      And I get the following response body:
        """
        {
          "method": "GET",
          "url_pattern": "/bar",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /var | kapow set /response/body",
          "index": 0,
          "id": ANY
        }
        """

  Scenario: Insert a route in the middle.
    A route can be inserted in the middle of the list
    by specifying an index less or equal to the last
    index in the request.

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
    Then I get 201 as response code
#      And I get "Created" as response reason phrase
      And I get the following response body:
        """
        {
          "method": "GET",
          "url_pattern": "/bar",
          "entrypoint": "/bin/sh -c",
          "command": "ls -la /var | kapow set /response/body",
          "index": 1,
          "id": ANY
        }
        """
