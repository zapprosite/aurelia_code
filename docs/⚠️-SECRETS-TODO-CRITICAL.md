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

### 2️⃣ Secret Rotation Policy (Fase 2 — Item 8, HIGH)
**Por quê:** API keys nunca rotacionadas = risco crescente.

```
CLOUDFLARE_MCP_TOKEN     → Trimestral (90 dias)
TELEGRAM_BOT_TOKEN       → Anual (365 dias)
GOOGLE_API_KEY           → Trimestral (90 dias)
OPENROUTER_API_KEY       → Mensal (30 dias) — exposto em logs facilmente
GROQ_API_KEY             → Trimestral (90 dias)
GITHUB_TOKEN             → Semestral (180 dias)
QDRANT_API_KEY           → Trimestral (90 dias)
POSTGRES_PASSWORD        → Anual (365 dias)
```

**Próximo agendamento:**
- 🔴 Cloudflare: 2026-06-20 (92 dias)
- 🔴 OpenRouter: 2026-04-20 (31 dias) ⚠️ PRIMEIRO!
- 🟡 Google: 2026-06-20
- 🟡 Groq: 2026-06-20
- 🟡 Qdrant: 2026-06-20

**Responsável:** codex (automático via cron) OU humano (manual)

---

### 3️⃣ Secret Audit Script (Fase 2 — Item 9, HIGH)
**Por quê:** Prevenir novo vazamento de credenciais.

```bash
# Implementar:
bash scripts/secret-audit.sh

# Deve detectar:
✅ Plaintext passwords em logs
✅ API keys em codigo fonte
✅ Tokens em comments
✅ Credenciais em .env non-centralizadas
```

**Status:** Scripts existem (verify-backups.sh, health-check.sh), mas secret-audit.sh NÃO.
**Responsável:** codex (implementar + cron semanal)

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

| # | Task | Prioridade | Esforço | Deadline |
|---|------|-----------|---------|----------|
| 1 | Rotar OpenRouter key | 🔴 CRÍTICA | 5 min | **2026-04-20** |
| 2 | Verificar systemd service | 🔴 CRÍTICA | 5 min | 2026-03-21 |
| 3 | Setup KeePassXC vault | 🟠 HIGH | 30 min | 2026-03-27 |
| 4 | Implementar secret-audit.sh | 🟠 HIGH | 15 min | 2026-03-31 |
| 5 | Agendar rotations | 🟡 MEDIUM | 10 min | 2026-04-01 |

---

## 📞 Lembretes Automáticos

Adicione ao seu calendario/alarme:
```bash
# 2026-04-20: Rotar OPENROUTER_API_KEY (30 dias)
# 2026-06-20: Rotar CLOUDFLARE, GOOGLE, GROQ, QDRANT (90 dias)
# 2026-09-20: Proximo ciclo trimestral
```

---

**Last Updated:** 2026-03-20
**Owner:** codex + humano (compartilhado)
**Authority:** [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
