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
Feature: Kapow! server rejects requests with semantic errors.
  Kapow! server will refuse to append routes when
  it receives a valid json document not conforming
  to the specification.

  Scenario: Error because lacking mandatory fields.
    If a request lacks any mandatory field the server
    responds with an error.

    Given I have a running Kapow! server
    When I append the route:
      """
      {
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | kapow set /response/body"
      }
      """
    Then I get 422 as response code
      And the response header "Content-Type" contains "application/json"
      And I get the following response body:
      """
      {
        "reason": "Invalid Route"
      }
      """

  Scenario: Error because bad route format.
    If a request contains an invalid expression in the
    field url_pattern the server responds with an error.

    Given I have a running Kapow! server
    When I append the route:
      """
      {
        "method": "GET",
        "url_pattern": "+123--",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | kapow set /response/body"
      }
      """
    Then I get 422 as response code
      And the response header "Content-Type" contains "application/json"
      And I get the following response body:
      """
      {
        "reason": "Invalid Route"
      }
      """
