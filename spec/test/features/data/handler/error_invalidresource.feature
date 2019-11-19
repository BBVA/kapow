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
Feature: Fail to retrieve an invalid resource for a handler in Kapow! server.
  If trying to access an invalid resource for a handler
  then the server responds with an error.

  Scenario: Try to get an inexistent resource from a handler.
    A request to retrieve an invalid resource
    from a handler will trigger a invalid resource error.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern |
      | GET    | /foo        |
    When I send a request to the testing route "/foo"
      And I get the resource "/invented/path"
    Then I get 400 as response code
      And the response header "Content-Type" contains "application/json"
      And I get the following response body:
      """
      {
        "reason": "Invalid Resource Path"
      }
      """
