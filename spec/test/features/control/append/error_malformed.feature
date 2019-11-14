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
Feature: Kapow! server reject append requests with malformed JSON bodies.
  Kapow! server will reject to append a route when
  it receives a malformed json document in the
  request body.

  Scenario: Error because a malformed JSON document.
    If a request comes with an invalid JSON document
    the server will respond with a bad request error.

    Given I have a running Kapow! server
    When I try to append with this malformed JSON document:
      """
      Hi! I am an invalid JSON document.
      """
    Then I get 400 as response code
      And the response header "Content-Type" contains "application/json"
      And I get the following response body:
      """
      {
        "reason": "Malformed JSON"
      }
      """
