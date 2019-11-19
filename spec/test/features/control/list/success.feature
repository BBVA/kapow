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
Feature: Listing routes in a Kapow! server.
  Listing routes allows users to know what URLs are
  available on a Kapow! server. The List endpoint returns
  a list of the routes the server has configured.

  Scenario: List routes on a fresh started server.
    A just started or with all routes removed Kapow! server,
    will show an empty list of routes.

    Given I have a just started Kapow! server
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response body:
        """
        []
        """

  Scenario: List routes on a server with routes loaded.
    After some route creation/insertion operations the server
    must return an ordered list of routes stored.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| kapow set /response/body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| kapow set /response/body |
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response body:
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
          }
        ]
        """
