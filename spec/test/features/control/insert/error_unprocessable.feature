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
Feature: Kapow! server rejects insertion requests with semantic errors.
  Kapow! server will refuse to insert routes when
  it receives a valid JSON but not conforming document.

  Scenario: Error because lacking mandatory fields.
    If a request lacks any mandatory fields the server
    responds with an error.

    Given I have a running Kapow! server
    When I insert the route:
      """
      {
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | kapow set /response/body"
      }
      """
    Then I get 422 as response code
#      And I get "Invalid Route" as response reason phrase

  Scenario: Error because wrong route specification.
    If a request contains an invalid expression in the
    url_pattern field the server responds with an error.

    Given I have a running Kapow! server
    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "+123--",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | kapow set /response/body",
        "index": 0
      }
      """
    Then I get 422 as response code
#      And I get "Invalid Route" as response reason phrase

  Scenario: Error because negative index specified.
    If a request contains a negative number in the
    index field the server responds with an error.

    Given I have a running Kapow! server
    When I insert the route:
      """
      {
        "method": "GET",
        "url_pattern": "+123--",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | kapow set /response/body",
        "index": -1
      }
      """
    Then I get 422 as response code
#      And I get "Invalid Route" as response reason phrase
