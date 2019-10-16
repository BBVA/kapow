Feature: Kapow!'s route subcommand
  As a Kapow! administrator
  In order to manage a running Kapow! server
  I need to be able to modify the server's routes list

  Scenario Outline: I'm warned when I fail to set mandatory flags
    In order to be able to invoke the Kapow! server some information must be
    provided such as server URL, action and route info.
    When I run Kapow! "<subcommand>" sub-command with environment "<envVars>" and commandline args "<cmdlineVars>"
    Then I get an error code <errCode> with message "<errMsg>"

    Examples:
    | subcommand | envVars                                 | cmdlineVars                                  | errCode | errMsg                                          |
    | route      |                                         |                                              | 1       | Error: Missing mandatory flag control-url       |
    | route      | KAPOW_CONTROL_URL=http://localhost:8080 |                                              | 1       | Error: Missing mandatory flag action            |
    | route      |                                         | --control-url                                | 1       | Error: Missing mandatory flag control-url's arg |
    | route      |                                         | --control-url=http://localhost:8080          | 1       | Error: Missing mandatory flag action            |
    | route      | KAPOW_CONTROL_URL=http://localhost:8080 | --append                                     | 1       | Error: Missing mandatory argument path          |
    | route      |                                         | --control-url=http://localhost:8080 --append | 1       | Error: Missing mandatory argument path          |
    | route      | KAPOW_CONTROL_URL=http://localhost:8080 | -A                                           | 1       | Error: Missing mandatory argument path          |
    | route      |                                         | --control-url=http://localhost:8080 -A       | 1       | Error: Missing mandatory argument path          |
    | route      | KAPOW_CONTROL_URL=http://localhost:8080 | --delete                                     | 1       | Error: Missing mandatory argument route-id      |
    | route      |                                         | --control-url=http://localhost:8080 --delete | 1       | Error: Missing mandatory argument route-id      |
    | route      | KAPOW_CONTROL_URL=http://localhost:8080 | -D                                           | 1       | Error: Missing mandatory argument route-id      |
    | route      |                                         | --control-url=http://localhost:8080 -D       | 1       | Error: Missing mandatory argument route-id      |
