Feature: Kapow!'s read subcommand
  As a kapow route developer
  In order to implement a route's command script
  I need to be able to access current request resources for that route

  Scenario Outline: Be warned when I fail to set mandatory flags
    In order to be able to invoke the Kapow! server some information should be
    provided such as server URL, handler id and resource path.
    When I run kapow "<subcommand>" sub-command with environment "<envVars>" and commandline args "<cmdlineVars>"
    Then I get an error code <errCode> with message "<errMsg>"

    Examples:
    | subcommand | envVars                                                       | cmdlineVars                                           | errCode | errMsg                                   |
    | read       |                                                               |                                                       | 1       | Error: Missing mandatory flag url        |
    | read       | KAPOW_DATA_URL=http://localhost:8080                          |                                                       | 1       | Error: Missing mandatory flag handler-id |
    | read       |                                                               | --data-url=http://localhost:8080                      | 1       | Error: Missing mandatory flag handler-id |
    | read       | KAPOW_DATA_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX |                                                       | 1       | Error: Missing mandatory argument path   |
    | read       | KAPOW_DATA_URL=http://localhost:8080                          | --handler-id=XXXXXXX                                  | 1       | Error: Missing mandatory argument path   |
    | read       |                                                               | --data-url=http://localhost:8080 --handler-id=XXXXXXX | 1       | Error: Missing mandatory argument path   |
