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
Feature: Setting values for handler response resources in Kapow! server.
    Users can set the values in the response
    resources by specifying the handler id
    and the resource path.

  Scenario: Set status code for the current response.
    Set the status code through the current
    handler.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern  |
      | GET    | /foo         |
    When I send a request to the testing route "/foo"
      And I set the resource "/response/status" with value "418"
      And I release the testing request
    Then I get 418 as response code in the testing request

  Scenario Outline: Set different resources for the current response.
    Set the following resources for the current
    response through the current handler.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern  |
      | GET    | /foo         |
    When I send a request to the testing route "/foo"
      And I set the resource "<resourcePath>" with value "<value>"
      And I release the testing request
    Then I get 200 as response code
#      And I get "OK" as response reason phrase
      And I get the value "<value>" for the response "<fieldType>" named "<elementName>" in the testing request

    Examples:
      | resourcePath            | value      | fieldType | elementName |
      | /response/headers/head1 | headVal1   | header    | head1       |
      | /response/cookies/cook1 | cookVal1   | cookie    | cook1       |
      | /response/body          | bodyValue1 | body      | -           |
      | /response/stream        | bodyValue2 | body      | -           |

  Scenario: Overwrite a resource for the current response.
    Write twice a on a resource, such as a gzip middleware would require:
    kapow get /response/body | gzip -c | kapow set /response/body
    although for simplicity, we'll just try overwriting the status code.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern |
      | GET    | /foo        |
    When I send a request to the testing route "/foo"
      And I set the resource "/response/status" with value "418"
      And I set the resource "/response/status" with value "200"
      And I release the testing request
    Then I get 200 as response code in the testing request
