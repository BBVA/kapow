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

Characters
----------

- Seasoned Ops
- Junior Ops

User Journey
------------

#. Actions over the server. Launch database backup script with an HTTP call.

  - User Learns: Add a route that executes a command locally.
  - Kapow! Concepts: `kapow route add`
  - Problem/Motivation: Each time an ACME project is finished it is
    desirable to make a backup of the entire database.  Given that the
    database server is a critical machine we don't want to grant SSH
    access to lowly developers.  The script is fast because the
    database is small (for now).
  - pre-Kapow! solution: Launching the script via SSH shell.

    .. code-block:: console

       $ ssh user@server
       Password:
       (server)$ ./backup_db.sh

  - Kapow!-enabled solution: Provide an HTTP endpoint that when accessed
    triggers the run of the backup script.

    .. code-block:: console

       $ kapow route add -X PUT /db/backup -e ./backup_db.sh

    .. code-block:: console

       $ curl -X PUT http://server:8080/db/backup

#. Basic server monitoring

  - User Learns: Execute local commands and output it results to the HTTP body.
  - Kapow! Concepts: `kapow set /response/body`
  - Problem/Motivation: The backup script produces a log on /tmp/backup_db.log.
    We want to share this log over HTTP to give users feedback about the backup
    process result.
  - pre-Kapow! solution: SSH into the host + cat /tmp/backup_db.log.
  - Kapow!-enabled solution: Provide an endpoint that returns the contents of
    /tmp/backup_db.log.

    .. code-block:: console

       $ cat /tmp/backup_db.log | kapow set /response/body

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

    .. code-block:: sh

       LINES="$(kapow get /request/params/lines)"
       FILTER="$(kapow get /request/params/filter)"
       grep "$FILTER" /var/log/backup_db.log \
         | tail -n"$LINES" \
         | kapow set /response/body

#. Advanced database monitoring

  - User Learns: Compose complex HTTP responses with more than one local command.
  - Kapow! Concepts: HEREDOC and subshells
  - Problem/Motivation: The OPs manager needs to have information about
    the health status of our servers. And she is always asking to the
    team to write a report that involves calling several commands.
  - pre-Kapow! solution: SSH into the server and manually execute the
    commands, collect the output and write the report.
  - Kapow!-enabled solution:

    From this:

    .. code-block:: sh

       echo Date: | kapow set /response/body
       echo ======...==== | kapow set /response/body
       echo Memory | kapow set /response/body
       # ...


    To this:

    .. code-block:: sh

       kapow set /response/headers/Content-Type text/plain
       {
         echo Date:
         date
         echo ================================================================================
         echo Memory:
         free -m
         echo ================================================================================
         echo Load:
         uptime
         echo ================================================================================
         echo Disk:
         df -h
       } | kapow set /response/body

#. Share your achievements

  - User Learns: Format a complex HTTP response with JSON format to feed the corporate dashboard.
  - Kapow! Concepts: backtick interpolation and `kapow set /response/headers`
  - Problem/Motivation: The OPs manager wants to create a dashboard to
    see the server health information in real time. She hired a fronted
    developer to make a nice dashboard application and we need to
    provide him with the information in a format suitable for display.
  - pre-Kapow! solution: Write a php/perl/python script to serve this
  - Kapow!-enabled solution:

    Don't handwrite `JSON`

    .. code-block:: sh

       kapow set /response/body application/json
       echo "{memory: `free -m`, ...uups...}"  | kapow set /response/body

    Use ``jq``

    .. code-block:: sh

       MEMORY=$(free -m)
       LOAD=$(uptime)
       DISK=$(df -h)
       kapow set /response/body application/json
       jq -nc --arg memory "$MEMORY" '{"memory": $memory}' | kapow set /response/body

Ideas
-----

- /request/params -> Filter the results of a backup query
- Use redirects to from one Kapow! server to another. I.e: 192.168.1.1/backups/{path:.*} --> 192.168.1.2/<path>
- /request/files -> Firma el fichero que sube el usuario y te lo devuelve firmado.


.. note::

   Add this to serve the webpage that uses the implemented HTTP API
   kapow route add / -c 'kapow set /resonse/headers/Content-Type text/html ; curl --output - http:// | kapow set /response/body'


Test
----

**User**

  Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
  tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At
  vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren,
  no sea takimata sanctus est Lorem ipsum dolor sit amet.

**Admin**

  Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
  tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At
  vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren,
  no sea takimata sanctus est Lorem ipsum dolor sit amet.

**User**

  Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod
  tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. 

  .. code-block:: console

     $ cat something.txt

  Right?
