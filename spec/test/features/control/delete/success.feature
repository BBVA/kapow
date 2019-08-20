Feature: Delete routes in Kapow! server.
  Delete routes allow users to remove non-desired
  routes from the server.

  Scenario: Delete a route.
    Routes are removed from the sever by specifying their id.

    Given I have a running Kapow! server
      And It has a route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
    When I delete the route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
    Then I get ok as response code
      And I get "OK" as response phrase
      And I get an empty response body
