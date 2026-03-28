# ADR 20260325: Básico Bem Feito v2 — Plano de Implementação Concreto

**Status:** Aprovado
**Data:** 25 de Março de 2026
**Autoridade:** Will (Principal Engineer)
**Supersede:** `20260325-basico-bem-feito-aurelia-team-memory-dashboard.md`, `20260325-basico-bem-feito-swarm-memoria-dashboard.md`
**Referência:** `20260325-data-stack-contract-and-templates.md` (governança de dados)

---

## 1. Diagnóstico Honesto (25/03/2026)

Este ADR nasce de uma auditoria direta do código. Não do que está nos documentos.

### O que funciona de verdade

| Componente | Estado Real | Evidência |
|---|---|---|
| **SQLite** | Sistema de registro ativo | `modernc.org/sqlite`, powers cron/memory/mailbox/tasks |
| **Qdrant** | HTTP client real, em uso | `internal/memory/semantic_search.go`, `internal/skill/semantic_router.go` |
| **Ollama/bge-m3** | Embed pipeline real | `EmbedText()` → `/api/embed`, usado em prod |
| **Gateway cascade** | Funcional | Gemma3:12b → Groq → OpenRouter → MiniMax, com GemmaJudge |
| **Dashboard** | SSE + embedded JS | `internal/dashboard/dist/`, port 3334, fan-out real |
| **Cron** | Poll-based, SQLite | `internal/cron/scheduler.go`, executa via agent loop |
| **Telegram** | Multi-bot operacional | 5+ bots configurados |

### O que não existe no código

| Componente | Estado Real | Nota |
|---|---|---|
| **Supabase** | Zero código, zero driver | Não está no `go.mod`. Apenas em ADRs. |
| **CapRover** | Zero Dockerfile | Deploy é `sudo install` + `systemctl`. |
| **Obsidian** | Zero integração | Nenhum código de sync no runtime. |
| **Testes de integração reais** | Todos mockados | `httptest.NewServer` em 100% dos casos. |
| **Fallback Ollama** | Inexistente | Se Ollama cai, embed falha com timeout. |
| **Rate limit parsing** | Não implementado | Campos `RateLimitRemaining` existem no struct mas nunca são populados. |
| **SSE replay** | Não implementado | Eventos perdidos se nenhum cliente está conectado. |

### Drift documental identificado

- `20260325-data-stack-contract-and-templates.md` — status "Aprovado", declara Supabase como canônico. Zero código.
- `20260325-basico-bem-feito-aurelia-team-memory-dashboard.md` — status "Proposto". Visão correta, sem deliverables concretos.
- `20260325-basico-bem-feito-swarm-memoria-dashboard.md` — duplicata da anterior. Substituída por este ADR.
- `codex_fez.md` (raiz do repo) — dump de sessão de terminal. Não é um documento. Deve ser deletado.

---

## 2. Decisão

Implementar em 4 fases incrementais, na ordem de dependência. Cada fase tem critério de aceite binário.

---

## 3. Fase 1 — Limpeza (sem código novo)

**Objetivo:** Tornar o repositório honesto. Nenhum documento afirma algo que o código não entrega.

| # | Ação | Arquivo | Critério |
|---|------|---------|---------|
| 1.1 | Deletar arquivo de ruído | `codex_fez.md` | Arquivo não existe |
| 1.2 | Marcar ADR duplicada como substituída | `20260325-basico-bem-feito-swarm-memoria-dashboard.md` | Status = "Substituída" |
| 1.3 | Corrigir status da ADR data-stack | `20260325-data-stack-contract-and-templates.md` | Status = "Aprovado (Supabase/Obsidian não integrados ao runtime)" |
| 1.4 | Reescrever índice com 3 categorias | `docs/adr/README.md` | Implementada / Parcial / Proposta |

---

## 4. Fase 2 — Testes de Integração Reais

**Objetivo:** Provar que a stack existente funciona ponta-a-ponta com serviços reais.

**Build tag:** `//go:build integration`
**Comando:** `go test -tags integration ./...`
**Regra:** CI/CD normal (`go test ./...`) não roda esses testes. Homelab roda com `-tags integration`.

| # | Teste | Arquivo | Depende de |
|---|-------|---------|------------|
| 2.1 | Qdrant CRUD real: create temp collection → upsert payload canônico → search → delete | `internal/memory/semantic_search_integration_test.go` | Qdrant em `QDRANT_URL` |
| 2.2 | Ollama embed real: `EmbedText()`, verificar dimensão 1024, cosine similarity > 0.7 para textos similares | `internal/memory/embed_integration_test.go` | Ollama em `OLLAMA_URL` com bge-m3 |
| 2.3 | Gateway local lane: prompt → Gemma3:12b → resposta, circuit breaker `closed` | `internal/gateway/provider_integration_test.go` | Ollama rodando |
| 2.4 | Dashboard SSE: start server → conectar SSE → `Publish(Event)` → verificar receipt em 2s | `internal/dashboard/dashboard_integration_test.go` | Sem dependência externa |
| 2.5 | Full-stack smoke: embed → Qdrant upsert → search → dashboard event publicado | `e2e/integration_test.go` | 2.1 + 2.2 + 2.4 |

**Critério de aceite:** `go test -tags integration ./...` verde no homelab.

---

## 5. Fase 3 — Infraestrutura Mínima

### 5A. Supabase Local (pgx direto, não REST)

**Decisão arquitetural:** SQLite **não migra**. Cron, messages, mailbox, tasks ficam em SQLite. Supabase recebe apenas dados que beneficiam de queryabilidade e persistência além do runtime local.

**Tabelas iniciais no Supabase:**
- `knowledge_items` — curados, longa duração, indexáveis
- `component_status` — snapshots históricos queryáveis

| # | Ação | Arquivo |
|---|------|---------|
| 5A.1 | `go get github.com/jackc/pgx/v5` | `go.mod` |
| 5A.2 | Campos: `SupabaseURL`, `SupabaseEnabled` | `internal/config/config.go` |
| 5A.3 | Store: `Connect`, `Ping`, `Close`, schema `knowledge_items` + `component_status` | `internal/store/supabase.go` |
| 5A.4 | Health check: `supabase_db` via `Ping` | `cmd/aurelia/health_checks.go` |
| 5A.5 | Wire no startup: se Supabase down → log warning → continue com SQLite | `cmd/aurelia/app.go` |

**Critério de aceite:** `curl /health | jq '.checks.supabase_db'` retorna `healthy` ou `degraded` (não ausente, não panic).

### 5B. CapRover

**Decisão arquitetural:** `modernc.org/sqlite` é pure Go → `CGO_ENABLED=0` funciona → alpine sem deps C.

| # | Ação | Arquivo |
|---|------|---------|
| 5B.1 | Dockerfile multi-stage: `golang:1.25-alpine` → `alpine:3.20`, ports 3334 + 8484 | `Dockerfile` |
| 5B.2 | captain-definition | `captain-definition` |
| 5B.3 | docker-compose para dev local: aurelia + qdrant + postgres | `docker-compose.yml` |

**Critério de aceite:** `docker build -t aurelia . && docker run --rm aurelia --help` sem erro.

### 5C. Markdown Brain canônico + Vault Obsidian

**Status de evolução:** esta fatia foi absorvida pelo desenho implementado em [20260327-markdown-brain-aurelia-code.md](20260327-markdown-brain-aurelia-code.md).

**Decisão arquitetural:** o vault continua read-only, mas agora como uma das fontes do `Markdown Brain` canônico. O repositório `.md` e o vault externo convergem para a mesma collection vetorial.

**Contrato de payload:** `source_system: "repo_markdown"` ou `source_system: "obsidian"`, ambos com `source_id` canônico e chunking por seção.

| # | Ação | Arquivo |
|---|------|---------|
| 5C.1 | Campos: `ObsidianVaultPath`, `ObsidianSyncEnabled` como gate da fonte externa do vault | `internal/config/config.go` |
| 5C.2 | Vault reader: walk `.md`, extrair frontmatter YAML + content | `internal/obsidian/reader.go` |
| 5C.3 | Indexador canônico: repo markdown + vault markdown → embeddings + Qdrant | `internal/markdownbrain/sync.go` |
| 5C.4 | Sync state: `markdown_brain_sync_state(source_system, source_path, sha256, last_synced)` em SQLite | `internal/markdownbrain/sync.go` |
| 5C.5 | Registrar tool e cron únicos: `markdown_brain_sync` | `cmd/aurelia/wiring.go`, `cmd/aurelia/seed_crons.go` |

**Critério de aceite:** `go test ./internal/markdownbrain ./internal/memory ./cmd/aurelia` verde.

### 5D. Homelab Monitoring Cron

| # | Ação | Arquivo |
|---|------|---------|
| 5D.1 | Checks Go: Docker containers, Ollama, GPU (`nvidia-smi`), Qdrant, disco | `internal/tools/homelab_checks.go` |
| 5D.2 | Cron job seed: a cada 15 min, roda checks, publica no dashboard | `cmd/aurelia/seed_crons.go` |
| 5D.3 | Expor resultados em `/api/status` com campos `component`, `status`, `last_ok_at` | `cmd/aurelia/dashboard_status.go` |

**Critério de aceite:** `curl /api/status | jq '.components'` lista GPU, Ollama, Qdrant, Docker com status real.

---

## 6. Fase 4 — Hardening

| # | Ação | Arquivo | Critério |
|---|------|---------|---------|
| 4.1 | `ErrOllamaUnavailable`: health check antes de embed, fallback para busca lexical em SQLite | `internal/memory/semantic_search.go` | Parar Ollama → status `degraded`, não crash |
| 4.2 | Parse `X-RateLimit-*` e `Retry-After` nos responses dos providers | `internal/gateway/provider.go` | Campos `RateLimitRemaining` populados nos logs |
| 4.3 | Ring buffer de 500 eventos, replay últimos 50 no reconect SSE | `internal/dashboard/dashboard.go` | Refresh do browser não perde contexto |
| 4.4 | Todos os componentes registram status em `/api/status` | `internal/dashboard/status.go` | Dashboard, `/status` Telegram e watchdog leem a mesma fonte |

---

## 7. Ordem de Execução

```
Fase 1 (Limpeza)         → sem dependências, executar primeiro
    │
Fase 2 (Testes reais)    → prova a stack existente
    │
    ├── 5A (Supabase)    → paralelo com 5B
    ├── 5B (CapRover)    → independente, pode rodar a qualquer momento
    ├── 5C (Obsidian)    → depende de 2.1 + 2.2 (Qdrant/Ollama provados)
    └── 5D (Monitoring)  → depende de Fase 2
    │
Fase 4 (Hardening)       → depende de Fase 2 + fases 3 relevantes
```

---

## 8. Verificação Final

```bash
# Fase 1
test ! -f codex_fez.md && echo "OK: ruído removido"

# Fase 2
go test -tags integration ./internal/memory/... ./internal/gateway/... \
  ./internal/dashboard/... ./e2e/...

# Fase 3A
curl -fsS http://127.0.0.1:8484/health | jq '.checks.supabase_db'

# Fase 3B
docker build -t aurelia-test . && docker run --rm aurelia-test --help

# Fase 3C
go test -tags integration ./internal/obsidian/...

# Fase 3D
curl -fsS http://127.0.0.1:3334/api/status | jq '.components'

# Fase 4: fallback Ollama
sudo systemctl stop ollama
curl -fsS http://127.0.0.1:3334/api/status | jq '.components.ollama.status'
# Esperado: "degraded"
sudo systemctl start ollama
```

---

## 9. Regras Permanentes

Proibido no ecossistema Aurélia após este ADR:

1. Marcar ADR como "Implementada" ou "Aprovado" sem código correspondente.
2. Criar integração nova sem campo de config, health check e degradação graceful.
3. Escrever no Qdrant sem `source_system`, `source_id` e `canonical_bot_id`.
4. Expandir SQLite como banco de negócio (apenas runtime/cache/queue).
5. Adicionar dump de sessão de terminal como documento de governança.

O padrão é: **sistema que degrada explicitamente é melhor que sistema que falha silenciosamente.**


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
