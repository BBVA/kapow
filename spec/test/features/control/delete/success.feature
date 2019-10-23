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
Feature: Delete routes in Kapow! server.
  Deleting routes allows users to remove undesired
  routes from the server.

  Scenario: Delete a route.
    Routes are removed from the sever by specifying their id.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern    | entrypoint | command                                          |
      | GET    | /foo           | /bin/sh -c | ls -la / \| kapow set /response/body                       |
      | GET    | /qux/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| kapow set /response/body |
    When I delete the first route
    Then I get 204 as response code
      And I get "No Content" as response reason phrase
