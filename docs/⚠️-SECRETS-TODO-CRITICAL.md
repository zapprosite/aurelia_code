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
| 7 | Verificar systemd service | ✅ DONE | 5 min | 2026-03-21 |
| 1 | Setup KeePassXC vault | 🟠 HIGH | 30 min | 2026-03-27 |
| 9 | Implementar secret-audit.sh | 🟠 HIGH | 15 min | 2026-03-31 |

---

## 📞 Próximas Ações

- 2026-03-27: Setup KeePassXC vault
- 2026-03-31: Implementar secret-audit.sh

---

**Last Updated:** 2026-03-20
**Owner:** codex + humano (compartilhado)
**Authority:** [ADR-20260319-Polish-Governance-All](./adr/ADR-20260319-Polish-Governance-All.md)
