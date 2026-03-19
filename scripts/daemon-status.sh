#!/bin/bash
set -euo pipefail
systemctl --user status aurelia.service --no-pager -l
