# Runbook: Systemd Timers — Memory Sync Fiscal

**Status:** ✅ Ativo
**Responsável:** antigravity
**Data:** 2026-03-19

---

## O que é?

4 systemd user timers que executam `memory-sync-fiscal.sh` em diferentes frequências para manter memória sincronizada entre repositório, embeddings (Qdrant) e metadados (Postgres).

---

## Timers Ativos

| Timer | Frequência | Próxima Execução | Propósito |
|-------|-----------|------------------|----------|
| `aurelia-memory-sync-fast` | 5 min | ~5min | Scan incremental de arquivos MD + embedding |
| `aurelia-memory-sync-postgres-index` | 15 min | ~15min | Reindex tabelas Postgres + sync metadata |
| `aurelia-memory-sync-validate` | 6:00am | Amanhã 6am | Validar consistência Qdrant + Postgres |
| `aurelia-memory-sync-compact` | 2am segunda | Seg 2am | Compactar storage + limpeza de métricas |

---

## Unit Files

**Location:** `~/.config/systemd/user/aurelia-memory-sync-*`

```
aurelia-memory-sync-fast.service
aurelia-memory-sync-fast.timer

aurelia-memory-sync-postgres-index.service
aurelia-memory-sync-postgres-index.timer

aurelia-memory-sync-validate.service
aurelia-memory-sync-validate.timer

aurelia-memory-sync-compact.service
aurelia-memory-sync-compact.timer
```

---

## Script

**Location:** `/home/will/aurelia/scripts/memory-sync-fiscal.sh`

**Funcionalidade:**
- Modo: fast | postgres-index | validate | compact
- Logging: `~/.aurelia/logs/memory-sync-fiscal.log`
- Métricas: `~/.aurelia/metrics/memory-sync.prom` (Prometheus format)
- Fallback: funciona mesmo sem Qdrant/Postgres rodando (validation-only)

**Invocação Manual:**
```bash
bash scripts/memory-sync-fiscal.sh --mode fast
bash scripts/memory-sync-fiscal.sh --mode postgres-index
bash scripts/memory-sync-fiscal.sh --mode validate
bash scripts/memory-sync-fiscal.sh --mode compact
```

---

## Status

### Verificar Timers
```bash
systemctl --user list-timers aurelia-memory-sync-* --all
```

### Verificar Logs
```bash
tail -f ~/.aurelia/logs/memory-sync-fiscal.log
```

### Verificar Métricas
```bash
cat ~/.aurelia/metrics/memory-sync.prom
```

### Verificar Status de Um Timer
```bash
systemctl --user status aurelia-memory-sync-fast.timer
```

---

## Troubleshooting

### Timer não está rodando
```bash
systemctl --user daemon-reload
systemctl --user enable aurelia-memory-sync-*.timer
systemctl --user start aurelia-memory-sync-*.timer
```

### Ver última execução
```bash
journalctl --user -u aurelia-memory-sync-fast.service -n 20
```

### Script falha
Verificar:
1. `~/.aurelia/logs/memory-sync-fiscal.log` para erros
2. `~/.claude/projects/-home-will-aurelia/memory/` existe?
3. `nc -z localhost 6333` (Qdrant)
4. `nc -z localhost 5432` (Postgres)

Se ambos unavailable, script roda em validation-only mode (sem erro).

---

## Próximos Passos

- [ ] Implementar real embedding generation (bge-m3)
- [ ] Conectar Qdrant real com collection "repository_memory"
- [ ] Conectar Postgres com schema ai_context.*
- [ ] Criar Grafana dashboard para métricas
- [ ] Adicionar alertas (Telegram) para falhas

---

## Arquivos Relacionados

- [memory-sync-architecture.md](../memory-sync-architecture.md)
- [memory-sync-fiscal-cron.md](./memory-sync-fiscal-cron.md)
- [Skill: /memory-sync-vector-db](../../.context/skills/memory-sync-vector-db/SKILL.md)
