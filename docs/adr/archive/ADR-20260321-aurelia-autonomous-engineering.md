> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

---
title: "Aurelia Autonomous Engineering — ULTRATRINK"
slug: aurelia-autonomous-engineering
status: active
date: 2026-03-21
decision-makers: [humano, aurelia, antigravity]
priority: P14
tags: [swarm, dashboard, prompt-system, skills, autonomy]
json: docs/adr/taskmaster/ADR-20260321-aurelia-autonomous-engineering.json
---

# ADR-20260321: Aurelia Autonomous Engineering — ULTRATRINK

## Contexto

A Aurélia tem ferramentas (77+), skills (~30), memória (Qdrant + SQLite) e um swarm de agentes Go nativo (v6.0-handoff). Mas falta a **cola cognitiva** que permite operar como engenheiro de software autônomo:

- ela não sabe **o que tem disponível** (77 tools chegam todas no prompt)
- ela não **planeja antes de executar** (execução linear sem PLAN→REVIEW→EXECUTE→VERIFY)
- ela não **entende o codebase** que está modificando (sem mapa de símbolos)
- ela não **combina skills** para tarefas compostas (router depende 100% da LLM)
- o dashboard é **TV** (passivo), não **cockpit** (interativo)
- a **memória não é comprimida** para caber no token budget

## Decisão

Implementar 7 capacidades em ordem de impacto decrescente, cada uma como sub-tarefa independente e deployável.

---

## Sub-Slices

### ✅ → Sub-1: Tool Introspection System
**Arquivo:** `internal/agent/tool_catalog.go`
**Impacto:** CRÍTICO — reduz contexto de 77 tools para 5-10 relevantes

- `BuildToolCatalog()` — JSON semântico de todas as tools com exemplos de uso
- `MatchToolsForTask(prompt)` — embedding match (Qdrant) para ranquear tools relevantes
- `InjectToolContext(catalog, k)` — injetar apenas top-K tools no system prompt

**Aceite:** Prompt com tarefa de debugging recebe ≤10 tools relevantes (não 77).

---

### Sub-2: Execution DNA System
**Arquivo:** `internal/persona/execution_dna.go`
**Impacto:** CRÍTICO — instrui a LLM sobre COMO agir, não apenas O QUÊ

- `ClassifyTask(prompt)` → enum: `debug|feature|refactor|research|ops|governance`
- `BuildExecutionDNA(taskType)` → workflow steps específicos por tipo
- `AssembleSystemPrompt(identity + dna + tools + memory)` → prompt assembly pipeline

**Aceite:** Task de "debug" injeta workflow de investigação (logs→config→fix→validate). Task de "feature" injeta workflow de engenharia (plan→code→test→commit).

---

### Sub-3: Planning Loop (PREV)
**Arquivo:** `internal/agent/planner.go`
**Impacto:** ALTO — transforma chatbot em engenheiro

- `Plan(task) → ActionPlan` — gerar plano estruturado antes de executar
- `SelfReview(plan) → ReviewResult` — self-review, detectar gaps
- `Execute(plan) → ExecutionLog` — executar passo a passo com checkpoints
- `Verify(log) → VerifyResult` — rodar smoke tests, backtrack se falhou

**Aceite:** Tarefa "corrige a config do bot" gera plano com 4+ passos antes de executar qualquer tool.

---

### Sub-4: Codebase Symbol Map
**Arquivo:** `internal/agent/codebase_map.go`
**Impacto:** ALTO — a Aurélia entende o código que vai modificar

- `BuildSymbolIndex(rootPath)` → parse Go AST, extrair funções/structs/interfaces/tipos
- `FindRelevantFiles(task, symbols)` → top-K arquivos por relevância semântica
- `GetDependencyGraph(pkg)` → quem importa quem (para impact analysis)
- `RefreshOnChange()` → cron de atualização automática pós-commit

**Aceite:** Prompt "modifica o token do Telegram" → lista automática de 3 arquivos afetados antes de editar.

---

### Sub-5: Semantic Skill Router
**Arquivo:** `internal/skill/semantic_router.go`
**Impacto:** MÉDIO — combina skills automaticamente

- `IndexSkills()` → embeddar descriptions no Qdrant (coleção `skill_index`)
- `MatchSkillsForTask(prompt, k)` → busca semântica top-K skills
- `ComposeSkillChain(skills)` → executar skills em sequência com shared context

**Aceite:** Prompt "investiga crash do bot" → match automático `bug-investigation` + `homelab-control` + `systems-engineer-homelab`.

---

### Sub-6: Dashboard Cockpit (P11)
**Arquivos:** `frontend/src/components/dashboard/`
**Impacto:** MÉDIO — humano pilota em real-time

- `CommandPalette.tsx` → CMD+K global para envio de prompts para a Aurélia
- `TaskBoard.tsx` → kanban de tasks ativas por agente (live via SSE)
- `PlanViewer.tsx` → visualizar ActionPlan antes de executar (com approve/reject)
- `ToolInspector.tsx` → ver em real-time qual tool está executando e com quais args

**Aceite:** Humano envia prompt via CMD+K, vê o plano gerado, aprova, e acompanha execução ferramenta por ferramenta.

---

### Sub-7: Memory Context Assembler
**Arquivo:** `internal/memory/context_assembler.go`
**Impacto:** MÉDIO-ALTO — memória usada eficientemente

- `ShouldSearch(prompt) bool` → decidir se precisa de busca semântica ou não
- `AssembleContext(prompt, budget) Context` → buscar Qdrant + codebase + skills
- `CompressContext(ctx, maxTokens)` → comprimir para caber no token budget da LLM
- `CronRefresh()` → atualizar embeddings pós-commit (já existe `memory-sync-fiscal.sh`)

**Aceite:** Token budget nunca excedido; contexto comprimido mantém informação crítica.

---

## Ordem de Execução

```
Sub-1 (Tool Introspection)   → 2-3h — maior ROI imediato
Sub-2 (Execution DNA)        → 2-3h — complemento direto do Sub-1
Sub-3 (Planning Loop)        → 4-6h — requer Sub-1 e Sub-2
Sub-4 (Codebase Symbol Map)  → 3-4h — paralelo ao Sub-3
Sub-5 (Semantic Skill Router)→ 3-4h — requer Qdrant já configurado
Sub-6 (Dashboard Cockpit)    → 4-6h — frontend independente
Sub-7 (Memory Assembler)     → 3-4h — refactor do MemoryManager existente
```

---

## Interface de Integração

Ponto de entrada único no bootstrap:

```go
// cmd/aurelia/app.go — bootstrapApp()
cognition := agent.NewCognitionEngine(
    agent.WithToolCatalog(registry),
    agent.WithExecutionDNA(cfg.DNAConfig),
    agent.WithPlanner(llmProvider),
    agent.WithCodebaseMap(cwd),
)
loop := agent.NewLoop(llmProvider, registry, cfg.MaxIterations,
    agent.WithCognition(cognition),
)
```

---

## Smoke Tests

```bash
# Sub-1: Tool Introspection
go test ./internal/agent/... -run TestToolCatalog -v

# Sub-2: Execution DNA
go test ./internal/persona/... -run TestExecutionDNA -v

# Sub-3: Planning Loop
go test ./internal/agent/... -run TestPlanningLoop -v

# Sub-4: Codebase Map
go test ./internal/agent/... -run TestCodebaseMap -v

# Sub-5: Skill Router
go test ./internal/skill/... -run TestSemanticRouter -v

# Full E2E
AURORA_COGNITION_TEST=1 go test ./cmd/aurelia/... -run TestCognitionE2E -v -timeout 120s
```

---

## Rollback

- Cada sub-slice tem flag feature (`AURELIA_COGNITION=0` desabilita)
- `WithCognition()` é opcional — se ausente, comportamento atual mantido
- Commits atômicos por sub-slice facilitam `git revert`

---

## Consequências

- A Aurélia passa de chatbot agêntico para **engenheiro de software autônomo**
- Token budget reduzido ~40% com tool filtering e context compression
- Latência de resposta aumenta ~500ms para planning (aceitável para tarefas complexas)
- Requer embeddings no Qdrant (já configurado em `conversation_memory`)
