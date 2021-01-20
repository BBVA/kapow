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
Feature: Don't leak the control token to the data server
  The control server needs to be invoked with a secret token that shall
  not be used when calling the data server.

  Scenario Outline: Use cli to communicate with data server
    The provided kapow subcommand must not send the X-Kapow-Token
    header.

    Given a test HTTP server on the data port
    When I run the following command
       """
       $ KAPOW_HANDLER_ID=myhandlerid KAPOW_CONTROL_TOKEN=testing kapow <subcommand> /

       """
    Then the HTTP server received a "<method>" request to "/handlers/myhandlerid/"
      And the received request doesn't have the header "X-Kapow-Token" set
    When the server responds with:
      | field                | value            |
      | status               | 200              |
      | headers.Content-Type | application/json |
      | body                 | null             |
    Then the command exits with "0"

    Examples:
      | subcommand | method |
      | get        | GET    |
      | set        | PUT    |
