# Runbook: Memory Sync Fiscal — Monitoramento Automático

**Status:** 🟡 Proposto
**Responsável:** Aurelia systemd (crons)
**Logs:** `~/.aurelia/logs/memory-sync-fiscal.log`
**Dashboard:** `docs/memory-sync-architecture.md`

---

## O que é o Fiscal?

Um **sistema de crons** que executa sincronização automática em diferentes frequências:

- **5 min:** Detecção rápida (fast mode)
- **15 min:** Indexação Postgres (postgres-index)
- **6am diária:** Validação + relatório (validate)
- **2am segunda:** Compactação (compact)

Isso **mantém Qdrant + Postgres sempre sincronizados** com a memória local, sem intervenção humana.

---

## Cron Entries

### 1. Fast Sync (5 minutos)

```bash
*/5 * * * * /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode fast >> ~/.aurelia/logs/memory-sync-fiscal.log 2>&1
```

**O que faz:**
- Varrer `~/.claude/projects/-home-will-aurelia/memory/` + `docs/adr/`
- Detectar arquivos novos/modificados (usando mtime + hash SHA256)
- Gerar embeddings (bge-m3, local)
- Upsert no Qdrant

**Duração:** ~30 segundos
**Carga:** Mínima (incremental)

---

### 2. Postgres Index Sync (15 minutos)

```bash
*/15 * * * * /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode postgres-index >> ~/.aurelia/logs/memory-sync-fiscal.log 2>&1
```

**O que faz:**
- Atualizar indices em `ai_context.memory_entries`
- Atualizar indices em `ai_context.adr_registry`
- Sync metadata (owner, type, tags)
- Registrar em `sync_log`

**Duração:** ~1 minuto
**Carga:** Baixa

---

### 3. Validação Diária (6am)

```bash
0 6 * * * /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode validate >> ~/.aurelia/logs/memory-sync-fiscal.log 2>&1
```

**O que faz:**
- Comparar count: Qdrant vs Postgres
- Detectar duplicatas
- Verificar integridade de embeddings
- Gerar relatório em `sync_log`
- Alert se count mismatch

**Duração:** ~2 minutos
**Carga:** Baixa
**Nota:** Executa em off-peak (6am)

---

### 4. Compactação Semanal (Segunda 2am)

```bash
0 2 * * 1 /home/will/aurelia/scripts/memory-sync-fiscal.sh --mode compact >> ~/.aurelia/logs/memory-sync-fiscal.log 2>&1
```

**O que faz:**
- Deletar entries antigas (> 90 dias) com type=plan ou temp
- Reindex Postgres (REINDEX CONCURRENTLY)
- Compactar Qdrant (se disponível)
- Gerar estatísticas de limpeza

**Duração:** ~5 minutos
**Carga:** Média (desligar buscas paralelas)
**Agendamento:** Segunda 2am (baixo uso)

---

## Instalação no Systemd

### Adicionar ao Timer do Aurelia

**Arquivo:** `~/.config/systemd/user/aurelia-memory-sync.service`

```ini
[Unit]
Description=Aurelia Memory Sync Fiscal
Requires=aurelia.service
After=aurelia.service

[Service]
Type=oneshot
User=will
ExecStart=/home/will/aurelia/scripts/memory-sync-fiscal.sh %i
StandardOutput=journal
StandardError=journal
SyslogIdentifier=aurelia-memory-sync
```

**Arquivo:** `~/.config/systemd/user/aurelia-memory-sync.timer`

```ini
[Unit]
Description=Aurelia Memory Sync Fiscal Timers
Requires=aurelia.service

[Timer]
OnBootSec=2min
OnUnitActiveSec=5min
Unit=aurelia-memory-sync.service

[Install]
WantedBy=timers.target
```

### Ativar

```bash
systemctl --user enable aurelia-memory-sync.timer
systemctl --user start aurelia-memory-sync.timer

# Verificar status
systemctl --user list-timers aurelia-memory-sync.timer
```

---

## Monitoramento

### Logs em Tempo Real

```bash
tail -f ~/.aurelia/logs/memory-sync-fiscal.log
```

### Visualizar Status no Postgres

```bash
# Últimas sincronizações
psql -d aurelia -c "
    SELECT sync_type, status, files_processed, completed_at
    FROM ai_context.sync_log
    ORDER BY completed_at DESC
    LIMIT 20;
"

# Contar entries por tipo
psql -d aurelia -c "
    SELECT type, COUNT(*) as count
    FROM ai_context.memory_entries
    GROUP BY type
    ORDER BY count DESC;
"
```

### Verificar Qdrant

```bash
# Count de points
curl -s http://localhost:6333/collections/repository_memory/points/count | jq

# Health check
curl -s http://localhost:6333/health | jq
```

---

## Alertas

O script `memory-sync-fiscal.sh` gera alertas em 3 cenários:

### 1. Count Mismatch

```
⚠️  Validation warning: Count mismatch! Qdrant=1250, Postgres=1248
```

**Ação:** Investigar; possível dados truncados no Postgres

### 2. Sync Timeout

Se sync > 5 minutos → flag na `sync_log`

**Ação:** Revisar logs; possível Qdrant lento ou Postgres travado

### 3. Embedding Error

Se bge-m3 falhar → log de erro + skip arquivo

**Ação:** Verificar modelo; possível OOM ou model not running

---

## Troubleshooting

### Sync travado (última entrada > 1 hora atrás)

```bash
# Verificar processo
ps aux | grep memory-sync-fiscal

# Se travado: matar e restart manual
pkill -f memory-sync-fiscal
/home/will/aurelia/scripts/memory-sync-fiscal.sh --mode fast

# Verificar logs
tail -100 ~/.aurelia/logs/memory-sync-fiscal.log
```

### Qdrant indisponível

```bash
# Verificar se Qdrant está rodando
curl http://localhost:6333/health

# Se down: restart
systemctl --user restart qdrant
```

### Postgres travado

```bash
# Check connections
psql -d aurelia -c "SELECT count(*) FROM pg_stat_activity;"

# Se muitas idle: cleanup
psql -d aurelia -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE state = 'idle';"
```

---

## Metrics

### Arquivo de Métricas

Script gera metrics em `~/.aurelia/metrics/memory-sync.prom` (formato Prometheus):

```
# HELP aurelia_memory_sync_duration_seconds Sync duration by mode
# TYPE aurelia_memory_sync_duration_seconds histogram
aurelia_memory_sync_duration_seconds_bucket{mode="fast",le="0.5"} 150
aurelia_memory_sync_duration_seconds_bucket{mode="fast",le="1.0"} 198
aurelia_memory_sync_duration_seconds_bucket{mode="fast",le="+Inf"} 200

# HELP aurelia_memory_sync_files_total Files processed
# TYPE aurelia_memory_sync_files_total counter
aurelia_memory_sync_files_total{mode="fast"} 250

# HELP aurelia_memory_entries_total Total entries by type
# TYPE aurelia_memory_entries_total gauge
aurelia_memory_entries_total{type="project_memory"} 45
aurelia_memory_entries_total{type="adr"} 12
```

---

## Próximos Passos

1. ✅ Documentar runbook (este arquivo)
2. 🔧 Criar script `memory-sync-fiscal.sh`
3. ⏰ Configurar systemd timers
4. 🧪 Teste: manual + cron
5. 📊 Validar métricas Prometheus
6. 🎯 Integrar com Aurelia bot queries
