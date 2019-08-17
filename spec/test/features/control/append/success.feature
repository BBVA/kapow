Feature: Append new routes in Kapow! server.
  Append routes allow users to configure the server. New
  routes are added to the list of existing routes.

  Scenario: Append the first route.
    A fresh server, just started or with all routes removed,
    will create a new list of routes. The newly created rule
    will be at index 0.

    Given I have a just started Kapow! server
      When I append the route:
        | method | url_pattern  | entrypoint | command                    |
        | GET    | /listRootDir | /bin/sh -c | ls -la / \| response /body |
      Then I get created as response code
      And I get "Created" as response phrase
      And I get the following entity as response body:
        | method | url_pattern  | entrypoint | command                    | index | id |
        | GET    | /listRootDir | /bin/sh -c | ls -la / \| response /body |     0 |  * |

  Scenario: Append another route.
    Appending routes on a non empty list will create new routes
    at the end of the list.

    Given I have a Kapow! server whith the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |
      When I append the route:
        | method | url_pattern | entrypoint | command                       |
        | GET    | /listEtcDir | /bin/sh -c | ls -la /etc \| response /body |
      Then I get created as response code
      And I get "Created" as response phrase
      And I get the following entity as response body:
        | method | url_pattern | entrypoint | command                       | index | id |
        | GET    | /listEtcDir | /bin/sh -c | ls -la /etc \| response /body |     2 |  * |
