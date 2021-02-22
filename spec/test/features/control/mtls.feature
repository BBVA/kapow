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
@wip
Feature: Kapow! server only allows connections from trusted clients.
  Trust is established via mTLS and certificate pinning.

  Scenario: Reject clients not providing a certificate.

    Given I have a running Kapow! server
    When I try to connect to the control API without providing a certificate
    Then I get a connection error


  Scenario: Reject clients providing an invalid certificate.

    Given I have a running Kapow! server
    When I try to connect to the control API providing an invalid certificate
    Then I get a connection error
