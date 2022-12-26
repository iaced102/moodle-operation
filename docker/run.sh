#!/bin/bash
env > /var/log/env
cron > /var/log/cron.log
/usr/sbin/apache2ctl -D FOREGROUND
