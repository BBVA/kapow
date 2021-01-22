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
Feature: Authenticate with server via token
  The control server needs to be invoked with a secret token that is
  send via the X-Kapow-Token header.

  @server
  Scenario: Try to get routes without Access Token
    The Access Token is mandatory.

    Given I have a just started Kapow! server
    When I request a route listing without providing an Access Token
    Then I get 401 as response code
      And I get "Unauthorized" as response reason phrase

  @server
  Scenario: Try to get routes with bad Access Token
    If the provided Access Token doesn't match with the one on the
    server side, the request must be denied.

    Given I have a just started Kapow! server
    When I request a route listing providing a bad Access Token
    Then I get 401 as response code
      And I get "Unauthorized" as response reason phrase

  @server
  Scenario: Auto-generate Access Token if KAPOW_CONTROL_TOKEN is undefined
    At startup and if undefined, a new random control token must be
    generated.  Any communication attempt from a client with an empty
    Control Token must be denied.

    Given I have a just started Kapow! server with no control token
    When I request a route listing without providing an empty Access Token
    Then I get 401 as response code
      And I get "Unauthorized" as response reason phrase

  @cli
  @server
  Scenario: Fail to start the server if KAPOW_CONTROL_TOKEN is empty
    At startup and if the provided token is an empty string the server
    will fail to start.

    Given I run the following command
       """
       $ KAPOW_CONTROL_TOKEN="" kapow server

       """
    Then the command exits immediately with "1"

  @cli
  @client
  Scenario Outline: Use cli to communicate with control server
    The provided kapow subcommand must send the X-Kapow-Token
    header.

    Given a test HTTP server on the control port
    When I run the following command
       """
       $ KAPOW_CONTROL_TOKEN=testing kapow <parameters>

       """
    Then the HTTP server received a "<method>" request to "<path>"
      And the received request has the header "X-Kapow-Token" set to "testing"
    When the server responds with:
      | field                | value            |
      | status               | <status>         |
      | headers.Content-Type | application/json |
      | body                 | {}               |
    Then the command exits with "0"

    Examples:
      | parameters           | method | path        | status |
      | route list           | GET    | /routes     | 200    |
      | route add / -c 'foo' | POST   | /routes     | 201    |
      | route remove bar     | DELETE | /routes/bar | 204    |

  @cli
  @client
  Scenario: Fail to start the client if KAPOW_CONTROL_TOKEN is empty
    At startup and if the provided token is an empty string the client
    will fail to start.

    Given a test HTTP server on the control port
    When I run the following command
       """
       $ KAPOW_CONTROL_TOKEN="" kapow route list

       """
    Then the command exits immediately with "1"

  @cli
  @client
  Scenario: Fail to start the client if KAPOW_CONTROL_TOKEN is missing
    At startup and if the variable KAPOW_CONTROL_TOKEN is not set, the
    client will fail to start.

    Given a test HTTP server on the control port
    When I run the following command
       """
       $ kapow route list

       """
    Then the command exits immediately with "1"
