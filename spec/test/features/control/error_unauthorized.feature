#
# Copyright 2021 Banco Bilbao Vizcaya Argentaria, S.A.
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
Feature: Fail to access the control API.
  When trying to access the control API without the mandatory
  Access Token, a 401 Unauthorized error should be produced.

  Scenario: Try to get routes without Access Token

    Given I have a just started Kapow! server
    When I request a route listing without providing an Access Token
    Then I get 401 as response code
      And I get "Unauthorized" as response reason phrase

  Scenario: Try to get routes with bad Access Token

    Given I have a just started Kapow! server
    When I request a route listing providing a bad Access Token
    Then I get 401 as response code
      And I get "Unauthorized" as response reason phrase