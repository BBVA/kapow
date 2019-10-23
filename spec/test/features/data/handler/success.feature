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
Feature: Retrieve a resource from a handler in Kapow! server.
    Users can retrieve request handler resources
    from the server by specifying the handler id
    and the resource path.

  Scenario: Retrieve a resource.
    Get the "request/path" resource for the current
    request through the handler id.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern  |
      | GET    | /foo         |
    When I send a request to the testing route "/foo"
      And I get the resource "/request/path"
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response raw body:
        """
        /foo
        """

  Scenario: Retrieve a resource item.
    Get the "request/headers/Host" item resource for
    the current request through the handler id.

    Given I have a Kapow! server with the following testing routes:
      | method | url_pattern  |
      | GET    | /foo         |
    When I send a request to the testing route "/foo"
      And I get the resource "/request/headers/Host"
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response raw body:
        """
        localhost:8080
        """
