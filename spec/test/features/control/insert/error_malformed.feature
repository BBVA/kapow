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
@skip
Feature: Kapow! server rejects insertion requests with malformed JSON bodies.
  Kapow! server will reject to insert a route when
  it receives a malformed JSON document in the
  request body.

  Scenario: Error because of malformed JSON document.
    If a request comes with an invalid JSON document
    the server will respond with a bad request error.

    Given I have a running Kapow! server
    When I try to insert with this JSON document:
      """
      {
        "method" "GET",
        "url_pattern": /hello,
        "entrypoint": null
        "command": "echo Hello
         World | kapow set /response/body",
        "index": 0,
        "id": "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
      }
      """
    Then I get 400 as response code
      And the response header "Content-Type" contains "application/json"
      And I get the following response body:
      """
      {
        "reason": "Malformed JSON"
      }
      """
