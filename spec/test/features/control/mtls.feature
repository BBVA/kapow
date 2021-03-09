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

Feature: Communications with the control interface are secured with mTLS.
  Trust is anchored via certificate pinning.
  The Kapow! server only allows connections from trusted clients.
  The Kapow! clients only establish connections to trusted servers.

  @server
  Scenario: Reject clients not providing a certificate.

    Given I have a running Kapow! server
    When I try to connect to the control API without providing a certificate
    Then I get a connection error

  @server
  Scenario: Reject clients providing an invalid certificate.

    Given I have a running Kapow! server
    When I try to connect to the control API providing an invalid certificate
    Then I get a connection error

  @client
  Scenario: Connect to servers providing a valid certificate.
    A valid certificate is the one provided via envvars.

    Given a test HTTPS server on the control port
    When I run the following command
    """
    $ kapow route list
    """
      And the HTTPS server receives a "GET" request to "/routes"
      And the server responds with
      | field                | value            |
      | status               | 200              |
      | headers.Content-Type | application/json |
      | body                 | []               |
    Then the command exits with "0"

  @client
  Scenario: Reject servers providing an invalid certificate.

    Given a test HTTPS server on the control port
    When I run the following command (with invalid certs)
    """
    $ kapow route list
    """
    Then the command exits immediately with "1"

  @server
  Scenario Outline: The control server is accessible through an alternative address
    The automatically generated certificated contains the Alternate Name
    provided via the `--control-reachable-addr` parameter.

    Given I launch the server with the following extra arguments
    """
    --control-reachable-addr "<reachable_addr>"
    """
    When I inspect the automatically generated control server certificate
    Then the extension "Subject Alternative Name" contains "<value>" of type "<type>"

    Examples:
      | reachable_addr | value     | type      |
      | localhost:8081 | localhost | DNSName   |
      | 127.0.0.1:8081 | 127.0.0.1 | IPAddress |
      | foo.bar:8081   | foo.bar   | DNSName   |
      | 4.2.2.4:8081   | 4.2.2.4   | IPAddress |
      | [2600::]:8081  | 2600::    | IPAddress |
