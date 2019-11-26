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

  ...  
