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
Feature: Authenticate server via fingerprint
  The client gets to verify the authenticity of the server by matching
  its SSL certificate with the provided fingerprint.

  @server
  @wip
  Scenario: The served certificate matches the provided fingerprint
    The fingerprint of the authentic SSL cert is provided via
    the KAPOW_CONTROL_FINGERPRINT variable to the init script.

    Given I have a just started Kapow! server
    When I connect to the control server
    Then the fingerprint of the certificate matches the value in "KAPOW_CONTROL_FINGERPRINT"
