I Need my Report!
=================

**Junior**

  Good morning!

  You look very busy, what's going on?

**Senior**

  I am finishing the capacity planning report.  Let me just mail it... Done!

  Today I am going to teach you how to write this report so we can split the
  workload.

**Junior**

  Oh.  That sounds... fun.  OK, tell me about this report.

**Senior**

  Here at ACME Inc. we take capacity planning seriously.  It is
  important that our employees always have the best resources to
  accomplish their job.

  We prepare a report with some statistics about the load of our
  servers.  This way we know when we have to buy another one.

**Junior**

  I see this company scales up just like Google...

**Senior**

  Smartass...

**Junior**

  (chuckles)

**Senior**

  We have a procedure:

  1. ``ssh`` into the machine.
  2. Execute the following commands copying its output for later filling in the
     report:

     - ``hostname`` and ``date``:  To include in the report.
     - ``free -m``:  To know if we have to buy more RAM.
     - ``uptime``:  To see the load of the system.
     - ``df -h``:  Just in case we need another hard disk drive.

  3. Copy all this in an email and send it to *Susan*, the operations manager.

**Junior**

  And why *Susan* can't ``ssh`` into the server herself to see all of this?

**Senior**

  She doesn't have time for this.  She is a manager, and she is very busy!

**Junior**

  Well, I guess we can make a *Kapow!* endpoint to let her see all this
  information from the browser.  This way she doesn't need to waste any time
  asking us.

  I started to write it already:

  .. code-block:: bash

     kapow route add /capacityreport -c 'hostname | kapow set /response/body; date | kapow set /response/body; free -m | kapow set /response/body; uptime | kapow set /response/body; df -h | kapow set /response/body'

**Senior**

  Not good enough!

  First of all, that code is not readable.  And the output would be something
  like:

  .. code-block:: none

     corporate-server
     Tue 26 Nov 2019 01:03:44 PM CET
                   total        used        free      shared  buff/cache   available
     Mem:          31967        2286       23473         729        6207       28505
     Swap:             0           0           0
      13:03:44 up  5:57,  1 user,  load average: 0.76, 0.63, 0.45
     Filesystem          Size  Used Avail Use% Mounted on
     dev                  16G     0   16G   0% /dev
     run                  16G  1.7M   16G   1% /run

  Which is also very difficult to read!

  What *Susan* is used to see is more like this:

  .. code-block:: none

     Hostname:
     ... the output of `hostname` ...
     ================================================================================
     Date:
     ... the output of `date` ...
     ================================================================================
     Memory:
     ... the output of `free -m` ...
     ================================================================================
     ... and so on ...

**Junior**

  All right, what about this?

  .. code-block:: bash

     kapow route add /capacityreport -c 'hostname | kapow set /response/body; echo ================================================================================ | kapow set /response/body; ...'

**Senior**

  That fixes the issue for *Susan*, but makes it worse for us.

  What about a HEREDOC to help us make the code more readable?

**Junior**

  A *HEREwhat*?

**Senior**

  A HEREDOC or **here document** is the method Unix shells use to
  express multi-line literals.

  They look like this:

  .. code-block:: console

     $ cat <<-'EOF'
             you can put
             more than one line
             here
     EOF

  The shell will put the data between the first ``EOF`` and the second
  ``EOF`` as the `stdin` of the :command:`cat` process.

**Junior**

  OK, I understand. That's cool, by the way.

  So, if I want to use this with *Kapow!*, I have to make it read the script
  from `stdin`.  To do this I know that I have to put a :nref-option:`-` at the end.

  Let me try:

  .. code-block:: bash

     kapow route add /capacityreport - <<-'EOF'
             hostname | kapow set /response/body
             echo ================================================================================ | kapow set /response/body
             date | kapow set /response/body
             echo ================================================================================ | kapow set /response/body
             free -m | kapow set /response/body
             echo ================================================================================ | kapow set /response/body
             uptime | kapow set /response/body
             echo ================================================================================ | kapow set /response/body
             df -h | kapow set /response/body
             echo ================================================================================ | kapow set /response/body
     EOF

**Senior**

  That would work.  Nevertheless I am not yet satisfied.

  What about all the repeated ``kapow set /response/body`` statements?
  Do you think we could do any better?

**Junior**

  Maybe we can redirect all output to a file and use the file as the input of
  ``kapow set /response/body``.

**Senior**

  There is a better way.  You can make use of another neat :command:`bash`
  feature:  `command grouping`_.

  Command grouping allows you to execute several commands treating the group as
  one single command.

  You can use this way:

  .. code-block:: bash

     { command1; command2; } | command3

**Junior**

  What about this:

  .. code-block:: bash

     kapow route add /capacityreport - <<-'EOF'
             {
                      hostname
                      echo ================================================================================
                      date
                      echo ================================================================================
                      free -m
                      echo ================================================================================
                      uptime
                      echo ================================================================================
                      df -h
                      echo ================================================================================
             } | kapow set /response/body
     EOF

**Senior**

  Nice!  Now I am not worried about maintaining that script.  Good job!

**Junior**

  You know me.  Whatever it takes to avoid writing reports ;-)

  (both chuckle).

.. _command grouping: https://www.gnu.org/software/bash/manual/html_node/Command-Grouping.html
