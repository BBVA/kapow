Feature: Append new routes in Kapow! server.
  Appending routes allows users to configure the server. New
  routes are added at the end of the list of existing routes.

  Scenario: Append the first route.
    A just started server or one with all routes removed,
    will create a new list of routes. The newly created rule
    will be at index 0.

    Given I have a just started Kapow! server
    When I append the route:
      """
      {
        "method": "GET",
        "url_pattern": "/listRootDir",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body"
      }
      """
    Then I get 201 as response code
      And I get "Created" as response reason phrase
      And I get the following response body:
      """
      {
        "method": "GET",
        "url_pattern": "/listRootDir",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la / | response /body",
        "index": 0,
        "id": ANY
      }
      """

  Scenario: Append another route.
    Appending routes on a non empty list will create new routes
    at the end of the list.

    Given I have a Kapow! server with the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |
    When I append the route:
      """
      {
        "method": "GET",
        "url_pattern": "/listEtcDir",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /etc | response /body"
      }
      """
    Then I get 201 as response code
      And I get "Created" as response reason phrase
      And I get the following response body:
      """
      {
        "method": "GET",
        "url_pattern": "/listEtcDir",
        "entrypoint": "/bin/sh -c",
        "command": "ls -la /etc | response /body",
        "index": 2,
        "id": ANY
      }
      """
