---
description: ADR mestra de governança industrial do homelab — secrets, dados, rede, operações, observabilidade e compliance.
status: proposed
---

# ADR-20260319-Polish-Governance-All

## Status

- Proposto

## Slice

- slug: polish-governance-all
- owner: humano + codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: docs/adr/taskmaster/ADR-20260319-Polish-Governance-All.json

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)
- [CONTRACT.md](/srv/ops/ai-governance/CONTRACT.md)
- [GUARDRAILS.md](/srv/ops/ai-governance/GUARDRAILS.md)

## Contexto

Auditoria completa do homelab will-zappro revelou **17 findings** (4 CRITICAL, 5 HIGH, 6 MEDIUM, 2 LOW). A infraestrutura tem governança forte em nível de autoridade (AGENTS.md, CONTRACT.md, GUARDRAILS.md) mas gaps críticos em segredos, dados, rede e operações.

### Estado Real do Sistema

- 32 containers Docker (6 stacks: platform, supabase, caprover, monitoring, voice, litellm)
- 4 PostgreSQL, 1 Qdrant, múltiplos SQLite, Supabase como mirror opcional
- RTX 4090 (14.3GB/24GB VRAM), ZFS tank 3.62TB (0.7% usado)
- Cloudflare Tunnel (6 subdomains públicos), Tailscale VPN, **UFW INATIVO**
- aurelia.service (systemd) com 8 MCP servers

## Decisão

Criar uma ADR unificada com 6 seções de governança que funciona como "constituição operacional" do homelab, subordinada a AGENTS.md e CONTRACT.md. Cada seção segue: Estado Atual → Gap → Meta Industrial → Steps → Validação → Owner.

## Escopo

6 seções de governança cobrindo todos os 17 findings da auditoria.

## Fora de escopo

- Execução automatizada dos remédios (requer humano para CRITICAL)
- Mudanças em AGENTS.md ou CONTRACT.md
- Criação de novos serviços ou containers

---

## Section 1: DATA_GOVERNANCE

**Findings cobertos:** #4 (sem contrato de Data Governance), #8 (sem Schema Registry), #13, #14, #15

### Estado Atual

| Componente | Papel | Source of Truth? | Governança |
|---|---|---|---|
| SQLite (`aurelia.db`) | gateway_route_states, voice_events, cron, tasks, memory | SIM — primário | Nenhuma |
| Postgres (4 instâncias) | n8n, supabase-db, litellm-db, dev | SIM por serviço | Nenhuma |
| Qdrant | conversation_memory (bge-m3) | NÃO — índice derivado | Mencionado em 8 ADRs |
| Supabase | sessions, messages (mirror opcional) | NÃO — réplica | Blueprints apenas |

### Entregas

1. **Store Selection Matrix** — árvore de decisão para onde cada tipo de dado vive
2. **Domain Ownership Table** — owner, SLA, criticidade por tabela
3. **Schema Registry** — todas as tabelas, colunas, índices, política de retenção
4. **Sync Flow** — SQLite (truth) → Supabase (mirror) → Qdrant (index)
5. **Data Lifecycle** — 30d hot, 90d warm, archive gzipped
6. **Qdrant Contract** — collection schema, embedding model (bge-m3), payload fields

### Validação

```bash
# Tabelas reais vs documentadas
sqlite3 ~/.aurelia/data/aurelia.db ".tables" | diff - docs/schema-registry-sqlite.md
# Postgres
docker exec supabase-db psql -U postgres -c "\dt" | diff - docs/schema-registry-postgres.md
```

### Owner

codex (documentação) + humano (validação de ownership)

---

## Section 2: SECRETS_GOVERNANCE

**Findings cobertos:** #1 (secrets plaintext Desktop), #2 (Postgres password em `ps aux`), #3 (credential leak em logs), #5 (API keys em `app.json` + 6 backups), #6 (MCP config repete token 3x), #7 (KeePassXC vault nunca criado), #9 (sem rotação de secrets)

### Estado Atual

- Credenciais em `~/Desktop/rascunho-s.txt` (plaintext)
- Postgres password visível via `ps aux | grep postgres`
- 6 cópias de `app.json.bak*` com API keys
- MCP config repete mesmo Cloudflare token 3 vezes
- KeePassXC guia existe mas vault nunca criado
- Nenhuma política de rotação

### Entregas

1. **KeePassXC vault** em `/srv/data/vault/aurelia.kdbx` (ZFS-backed)
2. **Migrar credenciais** do Desktop para vault, depois `shred -u ~/Desktop/rascunho-s.txt`
3. **Fix Postgres CLI** — usar `.pgpass` ou env var, nunca args de processo
4. **Deletar `app.json.bak*`** (6 cópias com credenciais)
5. **Refatorar MCP config** — env var único `$CF_API_TOKEN` para Cloudflare
6. **Log redaction rules** — journald + Docker (filtrar patterns de API key)
7. **`scripts/secret-audit.sh`** — scan de padrões (API keys, tokens, passwords)
8. **Rotação trimestral** — calendário com datas fixas

### Validação

```bash
# Nenhuma senha visível em processos
ps aux | grep -i password | grep -v grep  # deve retornar vazio

# Scan de segredos
bash scripts/secret-audit.sh  # exit 0

# Vault existe
test -f /srv/data/vault/aurelia.kdbx && echo "OK"
```

### Owner

**humano** (CRITICAL — requer manuseio manual de credenciais)

---

## Section 3: NETWORK_GOVERNANCE

**Gap foundacional — não coberto explicitamente no audit mas crítico para postura industrial**

### Estado Atual

- UFW INATIVO
- Cloudflare Tunnel com 6 subdomains públicos
- Tailscale VPN ativo mas sem ACL documentada
- Docker networks sem isolamento por stack

### Entregas

1. **UFW activation plan** — SSH first, per GUARDRAILS.md (requer aprovação humana)
2. **Port exposure matrix** — quais portas acessíveis de onde (localhost, Tailscale, público)
3. **Cloudflare Tunnel hardening** — revisar access policies por subdomain
4. **Docker network isolation** — redes separadas por stack
5. **Tailscale ACL policy** — documentar e aplicar regras

### Validação

```bash
# UFW ativo
sudo ufw status verbose  # deve mostrar regras ativas

# Port scan externo confirma lockdown
# (executar de máquina externa ou usar Cloudflare Radar)
```

### Owner

**humano** (requer aprovação per GUARDRAILS.md para mudanças de rede)

---

## Section 4: OPERATIONAL_GOVERNANCE

**Findings cobertos:** #11, #12, #16, #17

### Estado Atual

- Cron jobs existentes mas sem watchdog abrangente
- Sem verificação automática de backups
- Sem playbook de incidentes
- Deploy dirs duplicados
- Unclean shutdown warnings

### Entregas

1. **Cron schedule:**
   - `*/5 * * * *` — watchdog de containers
   - `*/15 * * * *` — health-check geral
   - `0 6 * * *` — smoke-test diário
2. **Backup verification script** — idade < 24h + integridade
3. **Incident Response Playbook:**
   - Service down → restart + alert
   - GPU OOM → kill processo + Telegram
   - ZFS degraded → snapshot + alert CRITICAL
   - Tunnel down → fallback Tailscale
   - Credential compromise → rotação imediata
4. **Service Dependency Map**
5. **Cleanup** de deploy dirs duplicados
6. **Fix** unclean shutdown warnings

### Validação

```bash
# Cron jobs instalados
crontab -l | grep -E '(watchdog|health-check|smoke)'

# Simular container unhealthy → alert Telegram
docker stop test-container && sleep 310 && docker start test-container
```

### Owner

codex (scripts) + humano (cron install, Telegram bot)

---

## Section 5: OBSERVABILITY_GOVERNANCE

**Finding coberto:** #10

### Estado Atual

- Prometheus e Grafana em stack monitoring
- Alertas ad-hoc, sem contrato de métricas
- Sem log rotation em Docker

### Entregas

1. **Prometheus Metrics Contract** — métricas obrigatórias, scrape intervals
2. **Alert Rules:**
   - GPU VRAM > 90%
   - Disk usage > 80%
   - Container unhealthy > 5min
   - Backup stale > 24h
3. **Alert delivery** via Telegram bot
4. **Grafana dashboard requirements** — 1 dashboard por stack
5. **Docker log rotation** — max-size/max-file em `daemon.json`
6. **Database metrics** — postgres_exporter + SQLite file size trends

### Validação

```bash
# Prometheus rules carregadas
curl -s localhost:9090/api/v1/rules | jq '.data.groups | length'

# Log rotation configurado
docker info --format '{{.LoggingDriver}}'
cat /etc/docker/daemon.json | jq '.["log-opts"]'
```

### Owner

codex (config) + humano (Telegram bot setup)

---

## Section 6: COMPLIANCE_MATRIX

**Transversal — mapeia todas as seções a CONTRACT.md e GUARDRAILS.md**

### Entregas

1. **Mapping table:** cada seção ↔ cláusula de CONTRACT.md/GUARDRAILS.md
2. **Audit schedule:**
   - Semanal: secrets scan
   - Mensal: health + capacity review
   - Trimestral: rotação de secrets + capacity planning
3. **Escalation paths:**
   - Alert → Telegram → human review
   - 24h SLA para CRITICAL
   - 72h SLA para HIGH
4. **Exception/waiver process** com data de expiração
5. **`scripts/governance-audit.sh`** — master script que orquestra todos os checks

### Validação

```bash
# Todas as seções referenciadas
grep -c "CONTRACT.md" docs/adr/ADR-20260319-Polish-Governance-All.md  # >= 3

# Audit cron ativo
crontab -l | grep governance-audit
```

### Owner

codex (scripts + documentação) + humano (aprovação de waivers)

---

## Execution Order (by severity)

### Fase 1 — CRITICAL (Dia 1, requer humano)

1. Criar vault KeePassXC per guia existente
2. Migrar credenciais do Desktop para vault
3. `shred -u ~/Desktop/rascunho-s.txt`
4. Fix Postgres CLI password (`.pgpass` ou env)

### Fase 2 — HIGH (Dias 2-3)

5. Deletar `app.json.bak*` com credenciais
6. Refatorar MCP config (env var para Cloudflare token)
7. Escrever Schema Registry (enumerar DBs primeiro)
8. Definir e documentar rotation policy

### Fase 3 — MEDIUM (Dias 4-7)

9. Instalar health checks no crontab
10. Criar backup verification script
11. Escrever incident response playbook
12. Documentar Qdrant collection schema
13. Definir data lifecycle policy
14. Definir observability contract + alert rules

### Fase 4 — LOW + HARDENING (Semana 2)

15. Cleanup deploy dirs
16. Fix unclean shutdown warnings
17. Planejar e documentar UFW activation
18. Escrever compliance matrix + audit schedule
19. Criar `governance-audit.sh`

---

## Arquivos afetados

| Arquivo | Ação |
|---|---|
| `docs/adr/ADR-20260319-Polish-Governance-All.md` | CREATE — esta ADR |
| `docs/adr/taskmaster/ADR-20260319-Polish-Governance-All.json` | CREATE — JSON taskmaster |
| `docs/adr/README.md` | UPDATE — registrar ADR |
| `/srv/data/vault/aurelia.kdbx` | CREATE — vault KeePassXC |
| `scripts/secret-audit.sh` | CREATE — scan de segredos |
| `scripts/governance-audit.sh` | CREATE — audit master |
| `/etc/docker/daemon.json` | UPDATE — log rotation |
| Múltiplos docs em `docs/` | CREATE — schema registry, playbooks |

## Simulações e smoke previstos

```bash
# Secrets
ps aux | grep -i password | grep -v grep  # vazio
bash scripts/secret-audit.sh  # exit 0

# Data
sqlite3 ~/.aurelia/data/aurelia.db ".tables"  # lista conhecida

# Network
sudo ufw status  # active (após Fase 4)

# Observability
curl -s localhost:9090/api/v1/rules | jq '.status'  # success

# Compliance
bash scripts/governance-audit.sh  # all green
```

## Rollback

- Cada seção é independente; rollback é por seção
- Scripts novos: deletar arquivo
- Cron jobs: `crontab -r` (com backup prévio via `crontab -l > /tmp/cron.bak`)
- UFW: `sudo ufw disable` (reversível)
- Vault: não destrutivo (adiciona, não remove)

## Consequências

### Positivas

- Homelab tratado como infraestrutura industrial
- Todos os 17 findings da auditoria endereçados
- Política de secrets com vault real e rotação
- Observabilidade contratual com alertas
- Compliance auditável com schedule automático

### Negativas

- Overhead operacional aumenta (cron jobs, audits)
- Requer comprometimento do humano para Fase 1 (CRITICAL)
- Manutenção de Schema Registry exige disciplina
