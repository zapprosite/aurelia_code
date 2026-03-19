#!/bin/bash
set -euo pipefail
AURELIA_HOME=${AURELIA_HOME:-"$HOME/.aurelia"}
touch "$AURELIA_HOME/logs/daemon.log" "$AURELIA_HOME/logs/daemon_error.log"
tail -f "$AURELIA_HOME/logs/daemon.log" "$AURELIA_HOME/logs/daemon_error.log"
