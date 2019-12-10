We need to filter
=================

**Senior**

  Hiya!  How're you doing this morning?  I've got a new challenge from our
  grateful mates.

  As time goes on from the last log rotation, the size of the log file gets
  bigger and bigger.  Furthermore, they want to limit the output of the file to
  pick only some records, and only from the end of the file.  We need to do
  something to help them as they are wasting a lot of time reviewing the output.

**Junior**

  I have a feeling that this is going to entail some serious *bash-foo*.  What
  do you think?

**Senior**

  Sure!  But in addition to some good shell plumbing we're going to squeeze
  *Kapow!*'s superpowers a litle bit more to get a really good solution.

  Can you take a look at *Kapow!*'s documentation to see if something can be
  done?

**Junior**

  I've read in the documentation that there is a way to get access to the data
  coming in the request.  Do you think we can use this to let them choose how
  to do the filtering?

**Senior**

  Sounds great!  How have we lived without *Kapow!* all this time?

  As they requested, we can offer them a parameter to filter the registers
  they want to pick, and another parameter to limit the output size in lines.

**Junior**

  Sounds about right.  Now we have to make some modifications to our last
  endpoint definition to add this new feature.  Let's get cracking!

**Senior**

  Well, we got it again, this is exactly what they need:

  .. code-block:: sh

     kapow route add /db/backup_logs -c 'grep -- "$(kapow get /request/params/filter)" /tmp/backup_db.log \
       | tail -n "$(kapow get /request/params/lines)" \
       | kapow set /response/body'

  It looks a bit weird, but we'll have time to revise the style later.  Please
  make some tests on your laptop before we publish it on the *Corporate Server*.
  Remember to send them an example URL with the parameters they can use to
  filter and limit the amount of lines they get.

**Junior**

  OK, should look like this, doesn't it?

  .. code-block:: console

     $ curl 'http://localhost:8080/db/backup_logs?filter=rows%20inserted&lines=200'

**Senior**

  Exactly.  Another great day helping the company advance.  Let's go grab a
  beer to celebrate!
