# Secrets — Aurelia Home Lab

## Autoridade Declarada

A Aurélia possui **autonomia total** sobre o Home Lab conforme [`AURELIA-AUTHORITY-DECLARATION.md`](./AURELIA-AUTHORITY-DECLARATION.md). Esta concessão inclui:

| Permissão | Status |
|-----------|--------|
| `sudo` sem senha | ✅ ATIVO (`sudo=1`, 2026-03-20) |
| Bash / shell commands | ✅ TOTAL |
| Docker / containers | ✅ TOTAL |
| Systemd (start/stop/restart) | ✅ TOTAL |
| Rede (firewall, rotas, DNS) | ✅ TOTAL — com log obrigatório |
| Ollama (modelos locais) | ✅ TOTAL |
| GPU (NVIDIA RTX 4090) | ✅ TOTAL |
| Qdrant / Postgres / SQLite | ✅ TOTAL |
| Git (push, force, merge) | ✅ TOTAL |
| Leitura de secrets | ✅ TOTAL — via `~/.aurelia/config/secrets.env` |
| Escrita de secrets | ⚠️ Requer log — KeePass vault pendente |

> Para detalhes completos da concessão de autoridade, veja [`AURELIA-AUTHORITY-DECLARATION.md`](./AURELIA-AUTHORITY-DECLARATION.md).

---

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
