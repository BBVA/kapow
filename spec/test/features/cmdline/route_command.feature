Feature: Kapow!'s route subcommand
  As a Kapow! administrator
  In order to manage a running Kapow! server
  I need to be able to manage the server's routes list

  Scenario Outline: Be warned when I fail to set mandatory flags
    In order to be able to invoke the Kapow! server some information should be
    provided such as server URL, handler id and resource.
    When I run Kapow! "<subcommand>" sub-command with environment "<envVars>" and commandline args "<cmdlineVars>"
    Then I get an error code <errCode> with message "<errMsg>"

    Examples:
    | subcommand | envVars                                                          | cmdlineVars                                                       | errCode | errMsg                                       |
    | get        |                                                                  |                                                                   | 1       | Error: Missing mandatory flag url            |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          |                                                                   | 1       | Error: Missing mandatory flag handler-id     |
    | get        |                                                                  | --control-url=http://localhost:8080                               | 1       | Error: Missing mandatory flag handler-id     |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX |                                                                   | 1       | Error: Missing mandatory flag action         |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          | --handler-id=XXXXXXX                                              | 1       | Error: Missing mandatory flag action         |
    | get        |                                                                  | --control-url=http://localhost:8080 --handler-id=XXXXXXX          | 1       | Error: Missing mandatory flag action         |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX | --append                                                          | 1       | Error: Missing mandatory argument path       |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          | --handler-id=XXXXXXX --append                                     | 1       | Error: Missing mandatory argument path       |
    | get        |                                                                  | --control-url=http://localhost:8080 --handler-id=XXXXXXX --append | 1       | Error: Missing mandatory argument path       |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX | -A                                                                | 1       | Error: Missing mandatory argument path       |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          | --handler-id=XXXXXXX -A                                           | 1       | Error: Missing mandatory argument path       |
    | get        |                                                                  | --control-url=http://localhost:8080 --handler-id=XXXXXXX -A       | 1       | Error: Missing mandatory argument path       |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX | --delete                                                          | 1       | Error: Missing mandatory argument route-id   |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          | --handler-id=XXXXXXX --delete                                     | 1       | Error: Missing mandatory argument route-id   |
    | get        |                                                                  | --control-url=http://localhost:8080 --handler-id=XXXXXXX --delete | 1       | Error: Missing mandatory argument route-id   |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX | -D                                                                | 1       | Error: Missing mandatory argument route-id   |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          | --handler-id=XXXXXXX -D                                           | 1       | Error: Missing mandatory argument route-id   |
    | get        |                                                                  | --control-url=http://localhost:8080 --handler-id=XXXXXXX -D       | 1       | Error: Missing mandatory argument route-id   |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX | --list                                                            | 0       |                                              |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          | --handler-id=XXXXXXX --list                                       | 0       |                                              |
    | get        |                                                                  | --control-url=http://localhost:8080 --handler-id=XXXXXXX --list   | 0       |                                              |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080 KAPOW_HANDLER_ID=XXXXXXX | -L                                                                | 0       |                                              |
    | get        | KAPOW_CONTROL_URL=http://localhost:8080                          | --handler-id=XXXXXXX -L                                           | 0       |                                              |
    | get        |                                                                  | --control-url=http://localhost:8080 --handler-id=XXXXXXX -L       | 0       |                                              |
