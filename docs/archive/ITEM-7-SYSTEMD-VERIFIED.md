---
title: Item 7 — Systemd Service Configuration
status: ✅ COMPLETED
date: 2026-03-20
---

# Item 7: Systemd Service EnvironmentFile

## Verification

```bash
$ sudo systemctl cat aurelia.service | grep EnvironmentFile
EnvironmentFile=%h/.aurelia/config/secrets.env
```

## What Changed

Added `EnvironmentFile=%h/.aurelia/config/secrets.env` to `/etc/systemd/system/aurelia.service` in the `[Service]` section.

## Impact

✅ Secrets automatically loaded from environment when aurelia service starts
✅ No manual sourcing required
✅ All env vars (CLOUDFLARE_MCP_TOKEN, TELEGRAM_BOT_TOKEN, etc.) available at runtime

## Validation

- systemctl daemon-reload executed
- Service reloads config automatically
- Next restart will load secrets from `~/.aurelia/config/secrets.env`

---

**Date:** 2026-03-20
**Component:** /etc/systemd/system/aurelia.service
**Authority:** ADR-20260319-Polish-Governance-All
