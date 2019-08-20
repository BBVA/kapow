Feature: Routes auto-ordering after inserting in a Kapow! server.

  When inserting routes the server will mantain the
  whole set of routes ordered an with consecutive indexes.

  Background:
    Given I have a Kapow! server whith the following routes:
      | method | url_pattern        | entrypoint | command                                          |
      | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |
      | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |

  Scenario: Inserting before the first route.
    After inserting before the first route the previous set
    will maintain their relative order and their indexes
    will be increased by one.

    When I insert the route:
      | method | url_pattern  | entrypoint | command                       | index |
      | GET    | /listVarDir  | /bin/sh -c | ls -la /var \| response /body |     0 |
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get the following entity as response body:
        | method | url_pattern  | entrypoint | command                       | index | id |
        | GET    | /listVarDir  | /bin/sh -c | ls -la /var \| response /body |     0 |  * |
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get a list with the following elements:
        | method | url_pattern        | entrypoint | command                                          | Index | id |
        | GET    | /listVarDir        | /bin/sh -c | ls -la /var \| response /body                    |     0 |  * |
        | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |     1 |  * |
        | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |     2 |  * |

  Scenario: Inserting after the last routes.
    After inserting after the last route the previous set
    will maintain their relative order and indexes.

    When I insert the route:
      | method | url_pattern  | entrypoint | command                       | index |
      | GET    | /listVarDir  | /bin/sh -c | ls -la /var \| response /body |     2 |
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get the following entity as response body:
        | method | url_pattern  | entrypoint | command                       | index | id |
        | GET    | /listVarDir  | /bin/sh -c | ls -la /var \| response /body |     2 |  * |
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get a list with the following elements:
        | method | url_pattern        | entrypoint | command                                          | Index | id |
        | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |     0 |  * |
        | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |     1 |  * |
        | GET    | /listVarDir        | /bin/sh -c | ls -la /var \| response /body                    |     2 |  * |

  Scenario: Inserting a midst route.
    After inserting a midst route the previous route set
    will maintain their relative order and the indexes
    of thefollowing routes will be increased by one.

    When I insert the route:
      | method | url_pattern  | entrypoint | command                       | index |
      | GET    | /listVarDir  | /bin/sh -c | ls -la /var \| response /body |     1 |
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get the following entity as response body:
        | method | url_pattern  | entrypoint | command                       | index | id |
        | GET    | /listVarDir  | /bin/sh -c | ls -la /var \| response /body |     1 |  * |
    When I request a routes listing
    Then I get 200 as response code
      And I get "OK" as response phrase
      And I get a list with the following elements:
        | method | url_pattern        | entrypoint | command                                          | Index | id |
        | GET    | /listRootDir       | /bin/sh -c | ls -la / \| response /body                       |     0 |  * |
        | GET    | /listVarDir        | /bin/sh -c | ls -la /var \| response /body                    |     1 |  * |
        | GET    | /listDir/{dirname} | /bin/sh -c | ls -la /request/params/dirname \| response /body |     2 |  * |
