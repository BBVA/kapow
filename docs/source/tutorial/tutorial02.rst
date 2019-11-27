What we've done?
================

**Senior**

  Ey, I come from seeing our project team mates.  They're delighted with their
  new toy, but they miss something.

  I forgot to tell you that after the backup run they need to review the log
  file to check that all went ok.

**Junior**

  Makes sense.  Do you think that Kapow! can help with this? I feel that this is
  the way to do it.

**Senior**

  Sure!  Lets take a look at the documentation to see how can we tweak with the
  logic of the request.

**Junior**

  Got it!  There're a
  `lot of resources to work with </theory/resource_tree.rst>`_, I see that we
  can write to the response, do you think this will work for us?

**Senior**

  Yeah, they used to cat the log file contents to see what happened in the
  last execution:

  .. code-block:: console

     $ cat /tmp/backup_db.log

  I've make it easy for you, do you dare to do it?

**Junior**

  Let me try add this to our ``POW`` file:

  .. code-block:: console

     kapow route add /db/backup_logs -e cat /tmp/backup_db.log | kapow set /response/body

**Senior**

  Looks good to me, clean and simple, and it is a very good idea to use ``GET``
  here as it wont change anything in the server. Restart Kapow! and try it.

**Junior**

  Wooow! I get back the content of the file. If they liked the first one they're
  going to love this.

**Senior**

  Agree. we've done the day with this...
