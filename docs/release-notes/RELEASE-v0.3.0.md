Kapow! v0.3.0

## Features

* Feature parity with the original PoC (written in Python), with some exceptions
  explained in the next section.

* Built and tested against Go 1.13.5.

* (Almost) no bytes were harmed during the development of this release.


## Known issues and limitations

### User server

* [#76][i76] Only plain `http` is supported for now, since `https` support is
  not yet complete. The `kapow server` flags `--certfile` and `--keyfile` are
  present but non-functional yet.


### Data API

* [#73][i73] `/response/body` and `/response/stream` behave identically for now.

* [#92][i92] `KAPOW_DATA_URL` is always set to its default value of
  `http://localhost:8082`, and not to the actual value, as specified by the
  `kapow server --data-bind "host:port"` invocation.


### Windows®

* [#83][i83] `kapow server file.pow` will try to run `file.pow` through `bash`, not `cmd`
  or `powershell`.
* `kapow route` default entrypoint is `/bin/sh`, so in order to use `cmd` or
  `powershell`, it must be explicitly set.


[i73]: https://github.com/BBVA/kapow/issues/73
[i76]: https://github.com/BBVA/kapow/issues/76
[i83]: https://github.com/BBVA/kapow/issues/83
[i92]: https://github.com/BBVA/kapow/issues/92
