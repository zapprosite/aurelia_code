---
title: Secrets Migration — Item 6 (Fase 2 HIGH)
description: Refactored MCP + App config to use environment variables instead of plaintext credentials
updated: 2026-03-20
---

# Secrets Migration — Item 6 Completion

## Changes Made

### 1. Centralized Secrets File
**Location:** `~/.aurelia/config/secrets.env` (outside repo, git-ignored)

All API keys stored in a single environment file:
```bash
CLOUDFLARE_MCP_TOKEN=your-cloudflare-token-here
TELEGRAM_BOT_TOKEN=your-telegram-bot-token-here
GOOGLE_API_KEY=your-google-api-key-here
OPENROUTER_API_KEY=your-openrouter-api-key-here
GROQ_API_KEY=your-groq-api-key-here
QDRANT_API_KEY=your-qdrant-api-key-here
GITHUB_TOKEN=your-github-token-here
POSTGRES_PASSWORD=your-postgres-password-here
```

### 2. Config Files Refactored (Placeholders)
**mcp_servers.json** → Uses `${CLOUDFLARE_MCP_TOKEN}` instead of plaintext (3x duplication eliminated)
**app.json** → Uses `${TELEGRAM_BOT_TOKEN}`, `${GOOGLE_API_KEY}`, etc.

Advantages:
- ✅ No plaintext credentials in version control
- ✅ Config files safe to commit/share
- ✅ Actual secrets only in runtime environment
- ✅ Easy to rotate secrets without code changes

### 3. Helper Scripts (In repo: `/scripts/`)
**scripts/config-loader.sh** — Sources secrets.env and exports to environment
**scripts/launcher.sh** — Wrapper that loads secrets before starting app

## Security Model

| Layer | Storage | Access | Protection |
|-------|---------|--------|------------|
| **Config Files** | `~/.aurelia/config/{app,mcp}.json` | Placeholders only | Safe to commit |
| **Secrets** | `~/.aurelia/config/secrets.env` | Actual credentials | Never commit (git-ignored) |
| **Runtime Env** | Process environment | Loaded by launcher | Only available at runtime |

## Integration

### Systemd Service Setup
```ini
[Service]
EnvironmentFile=%h/.aurelia/config/secrets.env
ExecStart=/bin/bash %h/aurelia/scripts/launcher.sh
```

### Manual Startup
```bash
# Option 1: Source before running
source <(bash ~/.aurelia/scripts/config-loader.sh)
aurelia-app --config ~/.aurelia/config/app.json

# Option 2: Use launcher wrapper
bash ~/.aurelia/scripts/launcher.sh
```

## Security Checklist

- ✅ Cloudflare token duplication eliminated (1 source, 3 references)
- ✅ App.json credentials moved to env vars (4 API keys: Telegram, Google, OpenRouter, Groq)
- ✅ Postgres password in env var (was in CLI args)
- ✅ MCP credentials in env vars (Qdrant API key, GitHub token)
- ✅ secrets.env should be added to .gitignore (if repo included it)
- ⏳ Next: Verify aurelia.service sources secrets.env (Item 7)

## Next Steps

1. **Item 7** — Verify aurelia.service systemd config uses EnvironmentFile
2. **Item 8** — Define secret rotation policy (quarterly, per-key sensitivity)
3. **Item 9** — Implement secret-audit.sh to detect plaintext leaks

## Links

- [OPERATIONAL_GOVERNANCE: Incident Response](./operational-governance-incident-response.md#runbook-credential-compromise)
- [ADR-20260319 Index](./ADR-20260319-GOVERNANCE-INDEX.md)
