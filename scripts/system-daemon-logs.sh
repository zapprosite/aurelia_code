#!/bin/bash
set -euo pipefail
sudo journalctl -u aurelia.service -f
