# Secrets — Aurelia Home Lab

## Status Atual

| Item | Status |
|------|--------|
| `secrets.env` centralizado, `chmod 600`, em `.gitignore` | ✅ |
| `scripts/secret-audit.sh` com 16+ padrões | ✅ |
| MCP e `app.json` usando env vars | ✅ |
| `systemd service` com `EnvironmentFile` | ✅ |
| KeePassXC vault (masterkey em hardware/TPM) | ⏳ deadline 2026-03-27 |
| Secret-audit no crontab semanal | ⏳ |

## Onde ficam os secrets

```
~/.aurelia/config/secrets.env   # fonte real — fora do repo
```

## Em caso de vazamento

```bash
# 1. Revogar imediatamente no provedor (Telegram, Groq, GitHub, OpenRouter, Cloudflare)
# 2. Auditar
bash scripts/secret-audit.sh
git log --all -S "password|token|key" | head -20
# 3. Criar ADR descrevendo o incidente
```

## Próximos passos

```bash
# Setup vault (humano executa):
bash scripts/setup-keepassxc-vault.sh
# Cron semanal:
echo "0 6 * * 1 bash ~/aurelia/scripts/secret-audit.sh" | crontab -
```
