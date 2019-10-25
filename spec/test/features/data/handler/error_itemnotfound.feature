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
Feature: Fail to retrieve nonexistent resource items in Kapow! server.
  If trying to access a nonexistent resource item
  then the server responds with a no content error.

  Scenario: Try to get a nonexistent resource item from a handler.
    A request to retrieve a nonexistent resource
    item from a handler will trigger a no content
    error.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern |
      | GET    | /foo        |
    When I send a request to the testing route "/foo"
      And I get the resource "/request/params/meloinvento"
    Then I get 404 as response code
#      And I get "Resource Item Not Found" as response reason phrase
