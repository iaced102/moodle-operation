#!/bin/bash

# cron each 4 hours
0 */4 * * *  bash $HOME/gits/moodle-operator/scripts/renew_token.sh
# cron at 0h every day
0 0 * * *  bash $HOME/gits/moodle-operator/scripts/renew_kubeconfig.sh
# renew nfs tracking
*/1 * * * *  bash $HOME/gits/moodle-operator/scripts/monitor_nfs.sh
