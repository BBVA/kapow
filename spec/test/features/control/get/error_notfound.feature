Feature: Fail to retrieve a route info in Kapow! server.
  When trying to get a route info in the server, if it
  does no exists the server respons with an error.

  Scenario: Try to get info for a non-existing route.
    A request of retrieving a non-existing route info
    will trigger a not found error.

    Given I have a just started Kapow! server
    When I get the info for route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
    Then I get 404 as response code
      And I get "Not Found" as response reason phrase
