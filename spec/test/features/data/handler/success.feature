Feature: Retrieve a resource from a handler in Kapow! server.
    Users can retrieve request handler resources
    from the server by specifying the handler id
    and the resource path.

  @wip
  Scenario: Retrieve a resource for the current request.
    Get the "request/path" resource for the current
    request through the handler id.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern        | entrypoint | command                    |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body |
    When I send a request to the route "/listRootDir"
      And I the get the resource "request/path" for the current request handler
    Then I get 200 as response code
      And I get "OK" as response reason phrase
      And I get the following response body:
        """
        /listRootDir
        """
