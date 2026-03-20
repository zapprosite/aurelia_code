🚨 **SEGREDOS — TAREFAS CRÍTICAS PENDENTES** 🚨

**STATUS ATUAL:** ✅ secrets.env implementado (mínimo configurado)

---

## ⚠️ PRÓXIMOS 30 DIAS (NÃO ESQUECER!)

### 1️⃣ KeePassXC Vault Setup (Fase 2 — Item 1, CRITICAL)
**Por quê:** secrets.env é plaintext no disco. Qualquer acesso ao sistema = todas credenciais comprometidas.

```bash
# Quando pronto, executar:
cd ~/aurelia
bash scripts/setup-keepassxc-vault.sh
# Resultado: ~/.aurelia/config/aurelia.kdbx + masterkey em hardware token ou TPM
```

**Responsável:** humano (requer decisão de como armazenar masterkey)
**Deadline:** 2026-03-27 (7 dias) — antes de qualquer mudança de rede

---

### 3️⃣ Secret Audit Script (Fase 2 — Item 9, HIGH)
**Status:** ✅ COMPLETED (2026-03-20)

**Implementado:** `scripts/secret-audit.sh` com cobertura completa:

```bash
# Executar auditoria:
bash scripts/secret-audit.sh

# Detecta:
✅ Plaintext passwords em logs
✅ API keys em codigo fonte
✅ Tokens em comments
✅ Credenciais em .env non-centralizadas
✅ Bearer tokens
✅ AWS keys (AKIA prefix)
✅ GitHub tokens (ghp_)
✅ Telegram tokens
✅ Private keys (RSA/DSA/EC/PGP/OpenSSH)
```

**Features:**
- Padrão-matching regex com 16+ patterns
- Git history scanning (`git log -p -S`)
- Log file scanning (`~/.aurelia/logs/`)
- Source code scanning (*.go, *.py, *.js, *.ts)
- .env file detection
- Color-coded output (RED=critical, YELLOW=warnings, GREEN=pass)
- Audit logging com timestamp
- Exit codes: 0=pass, 1=findings

**Próximo:** Adicionar ao crontab semanal
```bash
0 6 * * 1 bash ~/aurelia/scripts/secret-audit.sh
```

---

### 4️⃣ Systemd Service EnvironmentFile (Fase 2 — Item 7, HIGH)
**Por quê:** aurelia.service PRECISA sourcar secrets.env ao iniciar.

```bash
# Validar em:
sudo systemctl cat aurelia.service | grep EnvironmentFile

# Deve ter:
EnvironmentFile=%h/.aurelia/config/secrets.env
```

**Status:** ⏳ NÃO VERIFICADO
**Responsável:** codex (verificar + adicionar se falta)

---

## 🔐 Security Checklist

- ✅ secrets.env protegido (chmod 600)
- ✅ secrets.env em .gitignore
- ✅ MCP config usa env vars
- ✅ app.json usa env vars
- ⏳ KeePassXC vault criado
- ⏳ Rotation policy agendada
- ⏳ Secret audit cron rodando
- ⏳ systemd service sourca secrets

---

## 📋 Se houver Breach/Vazamento

1. **Imediatamente (< 5 min):**
   ```bash
   # Revogar credenciais:
   # - Cloudflare: dashboard → API tokens → revoke
   # - Telegram: @BotFather → revoke token
   # - Google: console.cloud.google.com → revoke
   # - OpenRouter: openrouter.ai → revoke
   # - Groq: console.groq.com → revoke
   # - GitHub: github.com/settings/tokens → delete
   # - Qdrant: http://localhost:6333 → change API key
   # - Postgres: ALTER USER root WITH PASSWORD 'new';
   ```

2. **Dentro de 30 min:**
   ```bash
   bash scripts/secret-audit.sh  # Verificar se há outras leaks
   git log --all -S "password\|token\|key" | head -20  # Audit history
   ```

3. **Dentro de 24h:**
   - Criar ADR para documentar como isso aconteceu
   - Implementar preventivo (e.g., pre-commit hook)

---

## 🎯 Ordem de Execução Recomendada

| # | Task | Prioridade | Esforço | Deadline | Status |
|---|------|-----------|---------|----------|--------|
| 7 | Verificar systemd service | ✅ DONE | 5 min | 2026-03-20 | ✅ Completo |
| 9 | Implementar secret-audit.sh | ✅ DONE | 15 min | 2026-03-20 | ✅ Completo |
| 1 | Setup KeePassXC vault | 🟠 HIGH | 30 min | 2026-03-27 | ⏳ Próxima semana |

---

## 📞 Próximas Ações

- ✅ 2026-03-20: Item 7 + Item 9 completados
- ⏳ 2026-03-27: Setup KeePassXC vault (Item 1)
- 📅 2026-03-27: Agendar secret-audit no cron semanal (`0 6 * * 1`)

---

**Last Updated:** 2026-03-20 (Item 7 + Item 9 completos)
**Owner:** codex + humano (compartilhado)
**Authority:** [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
