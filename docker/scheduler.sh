#!/bin/bash
[ -f /var/log/env ] && export $(cat /var/log/env | xargs)
/usr/bin/php /var/www/moodle/admin/cli/cron.php >> /var/log/schedule.log 2>&1
