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

#. Actions over the database. Launch database backup script with an HTTP call.
  - User Learns: Add a route that executes a command locally.
  - Kapow! Concepts: `kapow route add` 
#. Basic database monitorization
  - User Learns: Execute local commands and output it results to the HTTP body.
  - Kapow! Concepts: `kapow set /response/body`
#. Advanced database monitorization
  - User Learns: Compose complex HTTP responses with more than one local command.
  - Kapow! Concepts: HEREDOC and subshells
#. Share your achievements
  - User Learns: Format a complex HTTP response with JSON format to feed the corporate dashboard.
  - Kapow! Concepts: backtick interpolation and `kapow set /response/headers`
#. Unifing the interface (???)
  - User Learns: Add logic to the handler. React to a specific request. 
  - Kapow! Concepts: `kapow get /request/headers` and IF (bash)
  
  
 Ideas
------

/request/params -> Filter the results of a backup query
Use redirects to from one Kapow! server to another. I.e: 192.168.1.1/backups/{path:.*} --> 192.168.1.2/<path>
/request/files -> Firma el fichero que sube el usuario y te lo devuelve firmado.
  
  
.. note::

   Add this to serve the webpage that uses the implemented HTTP API
   kapow route add / -c 'kapow set /resonse/headers/Content-Type text/html ; curl --output - http:// | kapow set /response/body'
 
