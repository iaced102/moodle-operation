#!/bin/bash

# cron each 4 hours
0 */4 * * *  bash $HOME/gits/moodle-operator/moodle-operator/renew_token.sh
# cron at 0h every day
0 0 * * *  bash $HOME/gits/moodle-operator/moodle-operator/renew_kubeconfig.sh
