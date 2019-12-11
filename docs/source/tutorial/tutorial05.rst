Sharing the Stats
=================

**Junior**

  Good morning!

**Senior**

  Just about time...  We are in trouble!

  The report stuff was a complete success, so much so that now *Susan* has hired
  a frontend developer to create a custom dashboard to see the stats in real
  time.

  Now we have to provide the backend for the solution.

**Junior**

  And what's the problem?

**Senior**

  We are not developers!  What are we doing writing a backend?

**Junior**

  Just chill out.  Can't be that difficult...  What do they need, exactly?

**Senior**

  We have to provide a new endpoint to serve the same data but in JSON
  format.

**Junior**

  So, we have half of the work done already!

  What about this?

  .. code-block:: bash

     kapow route add /capacitystats - <<-'EOF'
             echo "{\"memory\": \"`free -m`\"}"  | kapow set /response/body
     EOF

**Senior**

  For starters, that's not valid ``JSON``.  The output would be something like:

  .. code-block:: console

     $ echo "{\"memory\": \"`free -m`\"}"
     {"memory": "              total        used        free      shared  buff/cache   available
     Mem:          31967        3121       21680         980        7166       27418
     Swap:             0           0           0"}

  You can't add new lines inside a ``JSON`` string that way, you have to escape
  the new line characters as ``\n``.

**Junior**

  Are you sure?

**Senior**

  See it for yourself.

  .. code-block:: console

     $ echo "{\"memory\": \"`free -m`\"}" | jq
     parse error: Invalid string: control characters from U+0000 through U+001F must be escaped at line 3, column 44

**Junior**

  :program:`jq`?  What is that command?

**Senior**

  :program:`jq` is a wonderful tool for working with ``JSON`` data from the command
  line.  With :program:`jq` you can extract data from a ``JSON`` document and it also
  allows you to generate a well-formed ``JSON`` document.

**Junior**

  Let's use it, then!

  How can we generate a ``JSON`` document with :program:`jq`?

**Senior**

  To generate a document we use the ``-n`` option:

  .. code-block:: console

     $ jq -n '{"mykey": "myvalue"}'
     {
       "mykey": "myvalue"
     }

**Junior**

  That does not seem very useful.  The output is just the same.

**Senior**

  Bear with me, it gets better.  You can add variables to the ``JSON`` and
  :program:`jq` will escape them for you.

  .. code-block:: console

     $ jq -n --arg myvar "$(echo -n myvalue)" '{"mykey": $myvar}'
     {
       "mykey": "myvalue"
     }

**Junior**

  Sweet!  That's just what I need.

  (hacks away for a few minutes).

  What do you think of this?

  .. code-block:: console

     $ jq -n --arg host "$(hostname)" --arg date "$(date)" --arg memory "$(free -m)" --arg load "$(uptime)" --arg disk "$(df -h)" '{"hostname": $host, "date": $date, "memory": $memory, "load": $load, "disk": $disk}'
     {
       "hostname": "junior-host",
       "date": "Tue 26 Nov 2019 05:27:24 PM CET",
       "memory": "              total        used        free      shared  buff/cache   available\nMem:          31967        3114       21744         913        7109       27492\nSwap:             0           0           0",
       "load": " 17:27:24 up 10:21,  1 user,  load average: 0.20, 0.26, 0.27",
       "disk": "Filesystem          Size  Used Avail Use% Mounted on\ndev                  16G     0   16G   0% /dev"
     }

**Senior**

  That is the output we have to produce, right.  But the code is far from
  readable.  And you also forgot about adding the endpoint.

  Can we do any better?

**Junior**

  That's easy:

  .. code-block:: bash

     kapow route add /capacitystats - <<-'EOF'
             jq -n \
                --arg hostname "$(hostname)" \
                --arg date "$(date)" \
                --arg memory "$(free -m)" \
                --arg load "$(uptime)" \
                --arg disk "$(df -h)" \
                '{"hostname": $hostname, "date": $date, "memory": $memory, "load": $load, "disk": $disk}' \
             | kapow set /response/body
     EOF

  What do you think?

**Senior**

  I'm afraid you forgot an important detail.

**Junior**

  I don't think so! the ``JSON`` is well-formed and it contains all the required
  data.  And the code is quite readable.

**Senior**

  You are right, but you are not using ``HTTP`` correctly.  You have to set the
  ``Content-Type`` header to let your client know the format of the data you are
  outputting.

**Junior**

  Oh, I see.  Let me try again:

  .. code-block:: bash

     kapow route add /capacitystats - <<-'EOF'
             jq -n \
                --arg hostname "$(hostname)" \
                --arg date "$(date)" \
                --arg memory "$(free -m)" \
                --arg load "$(uptime)" \
                --arg disk "$(df -h)" \
                '{"hostname": $hostname, "date": $date, "memory": $memory, "load": $load, "disk": $disk}' \
             | kapow set /response/body
             echo application/json | kapow set /response/headers/Content-Type
     EOF

**Senior**

  Better.  Just a couple of details.

  1. You have to set the headers **before** writing to the body.  This is
     because the body can be so big that *Kapow!* is forced to start sending it
     out.
  2. In cases where you want to set a small piece of data (like the header), it
     is better not to use ``stdin``.  *Kapow!* provides a secondary syntax
     for these cases:

     .. code-block:: console

        $ kapow set <resource> <value>

**Junior**

  Something like this?

  .. code-block:: bash

     kapow route add /capacitystats - <<-'EOF'
             kapow set /response/headers/Content-Type application/json
             jq -n \
                --arg hostname "$(hostname)" \
                --arg date "$(date)" \
                --arg memory "$(free -m)" \
                --arg load "$(uptime)" \
                --arg disk "$(df -h)" \
                '{"hostname": $hostname, "date": $date, "memory": $memory, "load": $load, "disk": $disk}' \
             | kapow set /response/body
     EOF

**Senior**

  That's perfect!  Now, let's upload this to the *Corporate Server* and tell the
  frontend developer about it.
