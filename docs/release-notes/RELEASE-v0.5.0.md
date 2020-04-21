Kapow! v0.5.0

## Features

* [#89][i89] Wrong environment variables exported to process

* [#98][i98] Implement a proper logging system (Partially solved. Added a logging feature for scripts in User-server for aid in debugging)

* [#92][i92] Spawn package uses static default value for KAPOW_DATA_URL

* [#102][i102] Handle race condition between Control API and pow files

* [#105][i105] Handle unexpected exit of servers


## Known issues and limitations

### Data API

* [#73][i73] `/response/body` and `/response/stream` behave identically for now.


### Windows®

* [#83][i83] `kapow server file.pow` will try to run `file.pow` through `bash`, not `cmd`
  or `powershell`.
* `kapow route` default entrypoint is `/bin/sh`, so in order to use `cmd` or
  `powershell`, it must be explicitly set.


[i73]: https://github.com/BBVA/kapow/issues/73
[i83]: https://github.com/BBVA/kapow/issues/83
[i89]: https://github.com/BBVA/kapow/issues/89
[i98]: https://github.com/BBVA/kapow/issues/98
[i92]: https://github.com/BBVA/kapow/issues/92
[i102]: https://github.com/BBVA/kapow/issues/102
[i105]: https://github.com/BBVA/kapow/issues/105
