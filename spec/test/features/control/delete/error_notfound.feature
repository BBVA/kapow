Feature: Fail to delete a route in Kapow! server.
  When trying to delete a route that not exists in the server
  the server respons with an error.

  Scenario: Delete a nonexistent route.
    A request of removing a nonexistent route
    will trigger a not found error.

    Given I have a just started Kapow! server
    When I delete the route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
    Then I get 404 as response code
      And I get "Not Found" as response reason phrase
