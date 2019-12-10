We need to filter
=================

**Senior**

  Hi guy!  How're you doing this morning?  I've got a new challenge from our
  grateful mates.

  As time goes on from the last rotation the size of the log file gets bigger.
  Furthermore, they want to limit the output of the file to pick only some
  records and only from the end of the file.  We need to do something to help
  them as they waste a lot of time reviewing the output.

**Junior**

  My two cents is that this time is going to be more of a deep bash knowledge.
  Do you agree?

**Senior**

  By sure, but in addition to some good shell plumbing we're going to squeeze
  Kapow!'s superpowers a litle bit more to get a really good solution.

  Can you take a look at Kapow!'s documentation to see if something can be done?

**Junior**

  I've seen in the documentation and there is a way to get access to the data
  comming in the request.  Do you think we can use this to let them choose how
  to do the filtering?

**Senior**

  Sounds great!  How we have lived without Kapow! all this time?

  As they requested, we can offer them with a parameter to filter the registers
  they want to pick and another parameter to limit the output size in lines.

**Junior**

  Sounds that will be enough.  Now we have to make some modifications to our
  last endpoint definition to add this new feature.  Let's start working...

**Senior**

  Well, we got it again, this is exactly what they need:

  .. code-block:: sh

     kapow route add /db/backup_logs -c 'grep -- "$(kapow get /request/params/filter)" /var/log/backup_db.log \
       | tail -n "$(kapow get /request/params/lines)" \
       | kapow set /response/body'

  It looks a bit weird but we'll have time to re-styling later.  Please make
  some tests on your laptop before to publish on the *Corporate Server*.
  Remember to send them an example URL with the parameters the can use to
  filter.

**Junior**

  Ok, Should look like this, isn't it?

  .. code-block:: console

     $ curl http://localhost:8080/db/backup_logs?filter=rows%20inserted&lines=200

**Senior**

  Exactly.  Another great day helping the company to advance.  Let's go for a
  beer for celebrating!
