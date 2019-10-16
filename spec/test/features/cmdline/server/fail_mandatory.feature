Feature: Kapow!'s server subcommand
  As a Kapow! administrator
  In order to allow remote access to commands in a host
  I need to be able to start a Kapow! server on that host

  Scenario Outline: I'm warned when I fail to set mandatory flags
    In order to be able to invoke the Kapow! server some information must be
    provided such as server URL, handler id and resource path.
    When I run Kapow! "<subcommand>" sub-command with environment "<envVars>" and commandline args "<cmdlineVars>"
    Then I get an error code <errCode> with message "<errMsg>"

    Examples:
    | subcommand | envVars | cmdlineVars | errCode | errMsg                                   |
    | server     |         |             | 1       | Error: Missing mandatory flag bind       |
    | server     |         | --bind      | 1       | Error: Missing mandatory flag bind's arg |
