#!/bin/bash

{
	echo -------------------------------------------------------------------------------- 
	echo "[$(date)] Starting backup procedure"
	for i in $(seq 1 10)
	do
		echo "[$(date)] Backing up table ${i}..."
		sleep .2
	done
	echo "[$(date)] Backup finished"
} >> /tmp/backup_db.log

echo 'Backup done!'
