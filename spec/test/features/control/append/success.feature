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
Feature: Append new routes in Kapow! server.
  Appending routes allows users to configure the server. New
  routes are added at the end of the list of existing routes.

  Scenario: Append the first route.
    A just started server or one with all routes removed,
    will create a new list of routes. The newly created rule
    will be at index 0.

    Given I have a just started Kapow! server
    When I append the route:
      """
      {
        "method": "GET",
        "url_pattern": "/foo",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | kapow set /response/body"
      }
      """
    Then I get 201 as response code
      And I get "Created" as response reason phrase
      And I get the following response body:
      """
      {
        "method": "GET",
        "url_pattern": "/foo",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | kapow set /response/body",
        "index": 0,
        "id": ANY
      }
      """

  Scenario: Append another route.
    Appending routes on a non empty list will create new routes
    at the end of the list.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| kapow set /response/body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| kapow set /response/body |
    When I append the route:
      """
      {
        "method": "GET",
        "url_pattern": "/baz",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /etc | kapow set /response/body"
      }
      """
    Then I get 201 as response code
      And I get "Created" as response reason phrase
      And I get the following response body:
      """
      {
        "method": "GET",
        "url_pattern": "/baz",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /etc | kapow set /response/body",
        "index": 2,
        "id": ANY
      }
      """
