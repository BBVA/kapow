I Need My Report
================

**Junior**

  Good morning!

  You look very busy, what's the matter?

**Senior**

  I am finishing the capacity planning report.  Let me just
  mail it... Done!

  Today I am going to teach you how to do this report so we can split
  the workload.

**Junior**

  Oh. That sounds... fun.  Ok, tell me about this report.

**Senior**

  Here at ACME Inc. we take capacity planning seriously.  It is
  important that our employees always have the best resources to
  accomplish their job.

  We prepare a report with some statistics about the load of our
  servers.  This way we know when we have to buy another one.
  
**Junior**

  I see this company scales just like Google.

**Senior**

  We have a procedure:

  1. SSH into the machine.
  2. Execute the following commands copying its output for later fill in
     the report:

     - ``hostname`` and ``date``:  To include in the report.
     - ``free -m``:  To know if we have to buy more RAM.
     - ``uptime``:  To see the load of the system.
     - ``df -h``:  Just in case we need another hard disk.

  3. Copy all this in a mail and send it to *Susan*, the operations
     manager.

**Junior**

  And why *Susan* don't enter the server herself to see all of this?

**Senior**

  She doesn't have time for this. She is a manager, she is very busy!

**Junior**

  Well, I guess we can make a Kapow! endpoint to let her see all this
  information from the browser.  This way she doesn't need to waste any
  time asking us.

  I started to write it already:

  .. code-block:: console

     kapow route add /capacityreport -c 'hostname | kapow set /response/body; date | kapow set /response/body; free -m | kapow set /response/body; uptime | kapow set /response/body; df -h | kapow set /response/body'

**Senior**

  That is preposterous!

  First of all that code is not readable.  And the output would be
  something like.

  .. code-block:: text

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

  What *Susan* is used to see is more like:

  .. code-block:: text

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

  .. code-block:: console

     kapow route add /capacityreport -c 'hostname | kapow set /response/body; echo ================================================================================ | kapow set /response/body; ...'

**Senior**

  That fix the issue for *Susan* but make it worst for us.

  What about a HEREDOC to help us make the code more readable.

**Junior**

  A *HEREwhat*?

**Senior**

  A HEREDOC or **here document** is the method Unix shells use to
  express multi-line literals.

  They look like this:

  .. code-block:: console

     $ cat <<HERE
        you can put
        more than one line
        here
       HERE
  
  The shell will put the data between the first ``HERE`` and the second
  ``HERE`` as the ``stdin`` of the ``cat`` process.

**Junior**

  If I want to use this with Kapow! I have to make it read the script
  from ``stdin``.  To do this I know that I have to put a ``-`` at the
  end.

  Let me try:

  .. code-block:: bash

     kapow route add /capacityreport - <<-HERE
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
     HERE

**Senior**

  That would work. Nevertheless I am not satisfied.

  What about all the repeated ``kapow set /response/body`` statements?
  Could we do any better?

**Junior**

  Maybe we can redirect all to a file and use the file as the input of
  ``kapow set /response/body``.

**Senior**

  There is a better way. You can make use of another neat ``bash``
  feature: **group commands**.

  Group commands allows you to execute several commands treating the
  group as one single command.

  You can use this way:

  .. code-block:: bash

     { command1; command2; } | command3

**Junior**

  What about this:

  .. code-block:: bash

     kapow route add /capacityreport - <<-HERE
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
     HERE

**Senior**

  I am not worried about maintaining that script. Good job!

**Junior**

  You know me. Whatever it takes to avoid writing reports ;)
