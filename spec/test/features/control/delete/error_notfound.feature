Feature: Fail to delete a route in Kapow! server.
  When trying to delete a route in the server, if it
  does no exists the server respons with an error.

  Scenario: Delete a non-existing route.
    A request of removing a non-existing route
    will trigger a not found error.

    Given I have a just started Kapow! server
    When I delete the route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
    Then I get not found as response code
      And I get "Not Found" as response phrase
      And I get an empty response body
