Feature: Listing routes in a Kapow! server

  Listing routes allow users to know what commands are
  available on a Kapow! server. The List endpoint returns
  a list of the routes the server has configured.

  Scenario: Listing routes on a fresh started server
    A fresh server, just started or with all routes removed,
    will show an empty list of routes.

    Given I have a just started Kapow! server
      When I request a routes listing
      Then I get an empty list

  Scenario: Listing routes on a server with routes loaded.
    After some route creation/insertion operations the server
    must return an ordered list of routes stored.

    Given I have a Kapow! server whith the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |
      When I request a routes listing
      Then I get a list with the following elements:
        | method | url_pattern        | entrypoint | command                                          | index | id |
        | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |  0    | *  |
        | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |  1    | *  |
