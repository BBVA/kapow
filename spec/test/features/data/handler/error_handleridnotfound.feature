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
Feature: Fail to retrieve resources from nonexistent handler in Kapow! server.
  If trying to access a nonexistent handler then the
  server responds with a not found error.

  Scenario: Try to get a valid resource path from a nonexistent handler.
    A request to retrieve a resource from a
    nonexistent handler will trigger
    a handler ID not found error.

    Given I have a running Kapow! server
    When I get the resource "/request/path" for the handler with id "XXXXXXXXXX"
    Then I get 404 as response code
      And I get "Handler ID Not Found" as response reason phrase

  Scenario: Try to get an invalid resource from a nonexistent handler.
    A request to retrieve an invalid resource from a nonexistent
    handler will trigger a handler ID not found error
    even if the resource is invalid.

    Given I have a running Kapow! server
    When I get the resource "/invalid/path" for the handler with id "XXXXXXXXXX"
    Then I get 404 as response code
      And I get "Handler ID Not Found" as response reason phrase
