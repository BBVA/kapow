<p align="center">
  <img alt="Kapow! Logo" src="https://raw.githubusercontent.com/BBVA/kapow/master/docs/source/_static/logo.png" width=300 />
  <h3 align="center">If you can script it, you can HTTP it.</h3>
  <p align="center">
    <a href="https://github.com/BBVA/kapow/actions?query=workflow%3ATest"><img alt="Test status" src="https://github.com/BBVA/kapow/workflows/Test/badge.svg" /></a>
    <a href="https://goreportcard.com/badge/github.com/bbva/kapow"><img alt="Go Report" src="https://goreportcard.com/badge/github.com/bbva/kapow" /></a>
    <a href="https://github.com/BBVA/kapow/issues/"><img alt="Open Issues" src="https://img.shields.io/github/issues/BBVA/kapow" /></a>
    <a href="https://kapow.readthedocs.io"><img alt="Documentation" src="https://img.shields.io/readthedocs/kapow" /></a>
    <a href="https://github.com/BBVA/kapow/releases"><img alt="Current Version" src="https://img.shields.io/github/v/release/BBVA/kapow?include_prereleases" /></a>
  </p>
</p>

---

## What's *Kapow!*

Say we have a nice cozy **shell command** that solves our problem. *Kapow!* lets
us easily **turn that into an HTTP API**.


### Let's see this with an example

We want to expose **log entries** for files not found on our **Apache Web
Server**, as an HTTP API.  With *Kapow!* we just need to write this file:

```bash
[apache-host]$ cat search-apache-errors.pow
kapow route add /apache-errors - <<-'EOF'
    cat /var/log/apache2/access.log | grep 'File does not exist' | kapow set /response/body
EOF
```

and then, run it using *Kapow!*

```bash
[apache-host]$ kapow server --bind 0.0.0.0:8080 search-apache-errors.pow
```

finally, we can read from the just-defined endpoint:

```bash
[another-host]$ curl http://apache-host:8080/apache-errors
[Fri Feb 01 22:07:57.154391 2019] [core:info] [pid 7:tid 140284200093440] [client 172.17.0.1:50756] AH00128: File does not exist: /usr/var/www/mysite/favicon.ico
[Fri Feb 01 22:07:57.808291 2019] [core:info] [pid 8:tid 140284216878848] [client 172.17.0.1:50758] AH00128: File does not exist: /usr/var/www/mysite/favicon.ico
[Fri Feb 01 22:07:57.878149 2019] [core:info] [pid 8:tid 140284208486144] [client 172.17.0.1:50758] AH00128: File does not exist: /usr/var/www/mysite/favicon.ico
...
```

### Why *Kapow!* shines in these cases

- We can share information **without having to grant SSH access** to anybody.
- We can share information easily **over HTTP**.
- We can effectively **limit** what gets executed.


## Documentation

You can find the complete documentation and examples [here](https://kapow.readthedocs.io).


## Authors

*Kapow!* is being developed by [BBVA-Labs Security team members](https://github.com/BBVA/kapow/blob/master/AUTHORS.rst).

*Kapow!* is Open Source Software and available under the [Apache 2
license](https://raw.githubusercontent.com/BBVA/kapow/master/LICENSE).


## Contributions

Contributions are of course welcome.  See
[CONTRIBUTING](https://raw.githubusercontent.com/BBVA/kapow/blob/master/CONTRIBUTING.rst)
or skim existing tickets to see where you could help out.
