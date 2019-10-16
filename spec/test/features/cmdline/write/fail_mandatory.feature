Feature: Kapow!'s write subcommand
  As a Kapow! route developer
  In order to implement a route's command script
  I need to be able to modify current request resources for that route

  Scenario Outline: I'm warned when I fail to set mandatory flags
    In order to be able to invoke the Kapow! server some information must be
    provided such as server URL, handler id, resource path and resource value.
    When I run Kapow! "<subcommand>" sub-command with environment "<envVars>" and commandline args "<cmdlineVars>"
    Then I get an error code <errCode> with message "<errMsg>"

    Examples:
    | subcommand | envVars                                                       | cmdlineVars                                                            | errCode | errMsg                                         |
    | write      |                                                               |                                                                        | 1       | Error: Missing mandatory flag data-url         |
    | write      | KAPOW_DATA_URL=http://localhost:8080                          |                                                                        | 1       | Error: Missing mandatory flag handler-id       |
    | write      |                                                               | --data-url                                                             | 1       | Error: Missing mandatory flag data-url's arg   |
    | write      |                                                               | --data-url=http://localhost:8080                                       | 1       | Error: Missing mandatory flag handler-id       |
    | write      | KAPOW_DATA_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX |                                                                        | 1       | Error: Missing mandatory argument path         |
    | write      | KAPOW_DATA_URL=http://localhost:8080                          | --handler-id                                                           | 1       | Error: Missing mandatory flag handler-id's arg |
    | write      |                                                               | --data-url=http://localhost:8080 --handler-id                          | 1       | Error: Missing mandatory flag handler-id's arg |
    | write      |                                                               | --data-url=http://localhost:8080 --handler-id=XXXXXXX                  | 1       | Error: Missing mandatory argument path         |
    | write      | KAPOW_DATA_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX | /response/status                                                       | 1       | Error: Missing mandatory argument value        |
    | write      | KAPOW_DATA_URL=http://localhost:8080                          | --handler-id=XXXXXXX /response/status                                  | 1       | Error: Missing mandatory argument value        |
    | write      |                                                               | --data-url=http://localhost:8080 --handler-id=XXXXXXX /response/status | 1       | Error: Missing mandatory argument value        |
