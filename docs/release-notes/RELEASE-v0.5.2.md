Kapow! v0.5.2

## Features

* Fix handling of misbehaving http clients

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
