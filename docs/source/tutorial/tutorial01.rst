Backup that Database!
=====================

**Junior**

  A Backup? Don't you have this kind of things already automated?

**Senior**

  Well, is not that simple. We of course have periodic backups. But, our
  project team ask us for a backup every time a project is finished.

  I've already prepared a script to do the task. Before executing it in
  production download it and test it in your own machine.

  .. todo::

     - Link backup script from Github.

**Junior**

  Ok, done! When I executed it the output says:

  .. code-block:: console

     $ ./backup_db.sh
     Backup done!
     Your log file is at /tmp/backup_db.log

**Senior**

  That's right. That script performed the backup and stored it into the
  **Backup Server** and appended some information into the backup log
  file at ``/tmp/backup_db.log``.

  Now you can SSH into the **Corporate Server** and make the real
  backup.


**Junior**

  Wait, wait... how long have you been doing this?


**Senior**

  This procedure was already here when I arrived.

**Junior**

  And why don't they do it themselves?  I mean, what do you contribute
  to the process?

**Senior**

  I am the only allowed to SSH into the **Corporate Server** for obvious
  reasons.

**Junior**

  Why do you need to SSH in the first place? Couldn't it be done
  without SSH?

**Senior**

  Actually it could be done with a promising new tool I've just found...
  Kapow!

  Is a tool that allows you to publish scripts as HTTP services.  If we
  use it here we can give them the ability to do the backup whenever
  they want.

**Junior**

  Sounds like less work for me.  I like it. 

**Senior**

  Ok then, let's try on your laptop first.

  First of all you have to follow the installation instructions XXX. 

**Junior**

  I've just installed it in my laptop, but I don't understand how all of
  this is going to work.

**Senior**

  Don't worry it is pretty easy.  Basically we will provide an HTTP
  endpoint managed by Kapow! at the **Corporate Server**; when the
  project team wants to perform a backup they only need to call the
  endpoint and Kapow! will call the backup script.

**Junior**

  It seems pretty easy.  How can I create the endpoint?

**Senior**

  First you have to start a fresh server. Please run this in your laptop:

  .. code-block:: console

     $ kapow server

  .. warning::

     It is important that you run this command in the same directory
     in which you downloaded ``backup_db.sh``.

**Junior**

  Done! But it doesn't do anything.

**Senior**

  Now you have the port 8080 open but don't have any endpoints defined.
  To define our endpoint you have to run this in another terminal:

  .. code-block:: console

     $ kapow route add -X PUT /db/backup -e ./backup_db.sh

  This will create an endpoint accessible via
  ``http://localhost:8080/db/backup``. This endpoint have to be invoked
  with the ``PUT`` method to prevent accidental calls.

**Junior**

  Cool! Do we need to do all this stuff every time we start the
  **Corporate Server**?

**Senior**

  Not at all. The have thought of everything. You can put all your route
  definitions in a special script file and pass it to the server on
  startup. They call those files `POW` files and have ``.pow``
  extension.

  It should look something like:

  .. code-block:: console

     $ cat backup.pow
     kapow route add -X PUT /db/backup -e ./backup_db.sh

  And then you can start Kapow! with it:

  .. code-block:: console

     $ kapow server backup.pow

**Junior**

  Great! Now it says:

  .. code-block:: console

     $ kapow server backup.pow
     2019/11/26 11:40:01 Running powfile: "backup.pow"
     {"id":"19bb4ac7-1039-11ea-aa00-106530610c4d","method":"PUT","url_pattern":"/db/backup","entrypoint":"./backup_db.sh","command":"","index":0}
     2019/11/26 11:40:01 Done running powfile: "backup.pow"

  I understand that this is proof that we have the endpoint available.

**Senior**

  That appears to be the case, but better we check it.

  Call it with ``curl``:

  .. code-block:: console

     $ curl -X PUT http://localhost:8080/db/backup

**Junior**

  Yay! I can see the log file at ``/tmp/backup_db.log``

**Senior**

  That's great. I am going to install all this in the *Corporate Server*
  and forget about the old procedure.

  That enough for your first day! You can go home.
