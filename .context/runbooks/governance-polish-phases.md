# Runbook: governance-polish Phases

**ADR:** [ADR-20260319-Polish-Governance-All](../../docs/adr/ADR-20260319-Polish-Governance-All.md)
**Skill:** `/governance-polish`
**Status:** 🟡 Proposto (Fase 1 pendente)
**Owner:** humano (Fase 1) + codex (Fases 2-4)

---

## Fase 1: CRITICAL (Requer Humano)

**Duração:** ~30 min
**Owner:** humano
**Bloqueio:** Sim — Fases 2-4 não podem começar sem conclusão

### Checklist

- [ ] Abrir KeePassXC
- [ ] Criar novo banco de dados: `/srv/data/vault/aurelia.kdbx`
- [ ] Definir senha mestra (20+ caracteres ou 5 palavras)
- [ ] Criar grupos: Cloudflare, Supabase, Infra, SSH
- [ ] Abrir `~/Desktop/rascunho-s.txt`
- [ ] Migrar cada segredo para uma entrada no KeePassXC
- [ ] Executar: `shred -u ~/Desktop/rascunho-s.txt`
- [ ] Fix Postgres password (usar `.pgpass` ou env var, não args)
- [ ] Confirmar: `ps aux | grep postgres` sem password visível
- [ ] Confirmar vault criado: `test -f /srv/data/vault/aurelia.kdbx && echo OK`

### Validação

```bash
# Após conclusão
ps aux | grep -i password | grep -v grep  # DEVE retornar vazio
test -f /srv/data/vault/aurelia.kdbx && echo "Vault OK"
```

### Próximo Passo

Quando concluir, notifique via mensagem ou comentário:
> Fase 1 concluída. Pode iniciar Fases 2-4 com `/governance-polish --phase 2 --execute`

---

## Fase 2: HIGH (Codex)

**Duração:** ~2 horas
**Owner:** codex
**Pré-requisito:** Fase 1 ✅

### Tarefas

1. **Deletar backups com credenciais**
   ```bash
   rm -v ~/app.json.bak*
   ```

2. **Refatorar MCP config** — usar env var única `$CF_API_TOKEN`
   - Arquivo: Verificar `~/.aurelia/config/app.json` ou MCP config
   - Buscar por `cloudflare-token` repetido 3x
   - Consolidar em 1 entrada com env var

3. **Escrever Schema Registry**
   - Enumerar todas as tabelas em SQLite + Postgres
   - Criar `docs/schema-registry-sqlite.md`
   - Criar `docs/schema-registry-postgres.md`
   - Registrar índices, retenção, ownership

4. **Documentar Rotation Policy**
   - Criar `docs/secrets-rotation-policy.md`
   - Definir trimestral com datas fixas
   - Registrar em calendário (cron job)

### Validação

```bash
bash scripts/secret-audit.sh  # exit 0 (sem credenciais em plaintext)
```

### Evidência

Atualizar JSON taskmaster com:
```json
{
  "phase_2_completed": {
    "backups_deleted": true,
    "mcp_config_refactored": true,
    "schema_registry_written": true,
    "rotation_policy_documented": true,
    "secret_audit_passing": true
  }
}
```

---

## Fase 3: MEDIUM (Codex)

**Duração:** ~3 horas
**Owner:** codex
**Pré-requisito:** Fase 2 ✅

### Tarefas

1. **Instalar health checks no crontab**
   ```
   */5 * * * * /home/will/aurelia/scripts/health-check.sh  # watchdog containers
   */15 * * * * /home/will/aurelia/scripts/backup-verification.sh  # backup check
   0 6 * * * /home/will/aurelia/scripts/smoke-test.sh  # daily smoke
   ```

2. **Criar backup verification script**
   - Verificar idade < 24h
   - Verificar integridade (test -f, du, md5)
   - Alert se stale

3. **Escrever incident playbook**
   - Container down → restart + alert
   - GPU OOM → kill + Telegram
   - ZFS degraded → snapshot + CRITICAL
   - Tunnel down → fallback Tailscale
   - Credential compromise → rotação imediata

4. **Documentar Qdrant schema**
   - Collection name: `conversation_memory`
   - Embedding model: `bge-m3`
   - Payload fields (timestamp, metadata, etc.)

5. **Definir data lifecycle**
   - 30d hot (active queries)
   - 90d warm (archived)
   - >90d gzipped em `/srv/backups/archive/`

6. **Definir observability contract**
   - Prometheus metrics obrigatórias
   - Alert rules em Prometheus
   - Grafana dashboards por stack

### Validação

```bash
crontab -l | grep -E '(health-check|backup|smoke)'
curl -s localhost:9090/api/v1/rules | jq '.data.groups | length'  # > 0
```

---

## Fase 4: LOW + HARDENING (Codex)

**Duração:** ~1-2 horas
**Owner:** codex + humano (UFW approval)
**Pré-requisito:** Fase 3 ✅

### Tarefas

1. **Cleanup deploy dirs**
   ```bash
   find ~/aurelia -type d -name deploy -o -name .deploy | wc -l
   # Se > 1: consolidar ou deletar duplicados
   ```

2. **Fix unclean shutdown warnings**
   - Revisar logs: `journalctl -u aurelia.service -n 50`
   - Adicionar graceful shutdown em `aurelia.service`

3. **Documentar UFW activation plan**
   - SSH first (port 22)
   - Cloudflare Tunnel (443)
   - Tailscale (não precisa UFW, é VPN)
   - Planejar com humano (requer aprovação GUARDRAILS.md)

4. **Escrever compliance matrix**
   - Tabela: Seção ↔ CONTRACT.md / GUARDRAILS.md
   - Audit schedule: semanal (secrets), mensal (health), trimestral (rotation)
   - Escalation paths e SLAs

5. **Criar `scripts/governance-audit.sh`**
   - Master script que orquestra todos os checks
   - `bash governance-audit.sh` retorna 0 (all green) ou lista falhas

### Validação

```bash
bash scripts/governance-audit.sh  # exit 0
```

---

## Status e Próximos Passos

**Atual:** Fase 1 (CRITICAL) — aguardando humano

**Próximo:**
1. Humano: Completar Fase 1 (vault, secrets, shred)
2. Codex: Executar Fases 2-4 com `/governance-polish --phase 2 --execute`
3. Todos: Validate com `bash scripts/governance-audit.sh`

**Escalation:**
- Bloqueado Fase 1 > 24h? → Telegram alert
- Fases 2-4 com erros? → Retry com `--force-rerun`

---

## Referências

- [ADR-20260319-Polish-Governance-All](../../docs/adr/ADR-20260319-Polish-Governance-All.md)
- [KeePassXC Tutorial](../../keepassxc-tutorial.html)
- [GUARDRAILS.md](/srv/ops/ai-governance/GUARDRAILS.md)
- [CONTRACT.md](/srv/ops/ai-governance/CONTRACT.md)
