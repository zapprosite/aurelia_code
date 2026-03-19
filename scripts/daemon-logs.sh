#!/bin/bash
set -euo pipefail
journalctl --user -u aurelia.service -f
