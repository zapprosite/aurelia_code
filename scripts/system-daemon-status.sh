#!/bin/bash
set -euo pipefail
sudo systemctl status aurelia.service --no-pager -l
