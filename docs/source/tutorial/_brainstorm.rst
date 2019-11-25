User Profile
============

Before You Start
----------------

Needed Skills:

#. Basic Linux
#. Proficient shell
#. How HTTP works

Needed Tools:

#. Linux machine
#. Installed Kapow! distribution

Motive
------

Use Kapow! to expose some Linux Box internal metrics and actions as an HTTP API for third party users.

Scenario
--------

User is a DevOps at ACME Company.

ACME Company appears to be a conglomerate which produces and sells every product type imaginable.

ACME's Infrastructure
---------------------

- 2 Linux machines

  - Corporate Server
  - Backup Server

User Journey
------------

#. Actions over the server. Launch database backup script with an HTTP call.

  - User Learns: Add a route that executes a command locally.
  - Kapow! Concepts: `kapow route add`
  - Problem/Motivation: Each time an ACME project is finished it is
    desirable to make a backup of the entire database.  Given that the
    database server is a critical machine we don't want to grant SSH
    access to lowly developers.  The script is very fast because the
    database is small (for now).
  - pre-Kapow! solution: Launching the script via SSH shell.
    ```
    ssh user@server
    $ ./backup_db.sh
    ```
  - Kapow!-enabled solution: Provide an HTTP endpoint that when accessed
    triggers the run of the backup script.
    ```
    curl -X PUT http://server:8080/db/backup
    ```
    ```
    kapow route add -X PUT /db/backup -e ./backup_db.sh
    ```

#. Basic server monitoring

  - User Learns: Execute local commands and output it results to the HTTP body.
  - Kapow! Concepts: `kapow set /response/body`
  - Problem/Motivation: The backup script produces a log on /tmp/backup_db.log.
    We want to share this log over HTTP to give users feedback about the backup
    process result.
  - pre-Kapow! solution: SSH into the host + cat /tmp/backup_db.log.
  - Kapow!-enabled solution: Provide an endpoint that returns the contents of
    /tmp/backup_db.log.
    ```
    cat /var/log/backup_db.log | kapow set /response/body
    ```

#. Filter over basic monitoring

  - User Learns: Get a parameter from the user and use it to select the
    script.
  - Kapow! Concepts: `kapow get /request/params`
  - Problem/Motivation: /tmp/backup_db.log keeps growing. It's about 100MB now.
    The users are fed up already.  We need a way to be more selective in the data
    we dump.
  - pre-Kapow! solution: SSH into the host, and then find a way to extract the
    required data from the log file. It would entitle using some combination of
    grep, tail, etc.  Or we could provide a bespoke shell script to accomplish
    this task.

  - Kapow!-enabled solution:
    ```
    LINES="$(kapow get /request/params/lines)"
    FILTER="$(kapow get /request/params/filter)"
    grep "$FILTER" /var/log/backup_db.log \
      | tail -n"$LINES" \
      | kapow set /response/body
    ```

#. Advanced database monitoring

  - User Learns: Compose complex HTTP responses with more than one local command.
  - Kapow! Concepts: HEREDOC and subshells
  - Problem/Motivation:
  - pre-Kapow! solution:
  - Kapow!-enabled solution:
    ```
    {
      echo Memory:
      free -m
      echo ================================================================================
      echo Load:
      uptime
      echo ================================================================================
      echo Disk:
      df -h
    } | kapow set /response/body
    ```

#. Share your achievements

  - User Learns: Format a complex HTTP response with JSON format to feed the corporate dashboard.
  - Kapow! Concepts: backtick interpolation and `kapow set /response/headers`
  - Problem/Motivation:
  - pre-Kapow! solution:
  - Kapow!-enabled solution:
    ``` DON'T HANDWRITE JSON
    echo "{memory: `free -m`, ...uups..}'  | kapow set /response/body
    ```

    ``` USE JQ
    MEMORY=$(free -m)
    LOAD=$(uptime)
    DISK=$(df -h)
    jq -nc --arg memory "$MEMORY" '{"memory": $memory}'
    ```

Ideas
-----

- /request/params -> Filter the results of a backup query
- Use redirects to from one Kapow! server to another. I.e: 192.168.1.1/backups/{path:.*} --> 192.168.1.2/<path>
- /request/files -> Firma el fichero que sube el usuario y te lo devuelve firmado.


.. note::

   Add this to serve the webpage that uses the implemented HTTP API
   kapow route add / -c 'kapow set /resonse/headers/Content-Type text/html ; curl --output - http:// | kapow set /response/body'
