Feature: Fail to retrieve route details in Kapow! server.
  When trying to get route details for a route that
  does no exist the server responds with an error.

  Scenario: Try to get details for a nonexistent route.
    A request for retrieving details for a nonexistent
    route will trigger a not found error.

    Given I have a just started Kapow! server
    When I get the route with id "xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx"
    Then I get 404 as response code
      And I get "Not Found" as response reason phrase
