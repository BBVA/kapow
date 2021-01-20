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
@cli
@client
Feature: Authenticate with server via token
  The control server needs to be invoked with a secret token that is
  send via the X-Kapow-Token header.

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
