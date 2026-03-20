---
description: Slice nonstop para implementar o agent bus hierárquico da Aurelia em PostgreSQL com memória derivada em Qdrant.
status: proposed
---

# ADR-20260319-hierarchical-agent-bus

## Status

- Proposto

## Slice

- slug: hierarchical-agent-bus
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: `docs/adr/taskmaster/ADR-20260319-hierarchical-agent-bus.json`

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O desenho do swarm já foi fechado por blueprint e ADR. Falta agora o runtime de verdade:

- bus operacional
- filas
- leases
- canais e threads
- agentes fixos e efêmeros
- ajuda entre agentes ociosos

## Decisão

Implementar a `Slice 11` com:

- `PostgreSQL` como verdade operacional
- `Qdrant` como memória semântica derivada
- contrato de handoff inspirado em `open-agent-supervisor` e `langgraph-supervisor`
- dashboard leve em `Go`
- assistance queue para agentes ociosos

## Escopo

- esquema SQL inicial
- filas com lease
- channels/threads/messages
- assistance queue
- resumos indexáveis no Qdrant
- endpoints de dashboard

## Fora de escopo

- UI completa de board-room
- multi-GPU scheduling
- automação de deploy

## Arquivos afetados

- `internal/agent/task_store_swarm.go`
- `internal/agent/task_store_swarm_test.go`
- `internal/agent/task_store_schema.go`
- `internal/agent/`
- `internal/tools/`
- `cmd/aurelia/`
- `docs/agent_swarm_dashboard_blueprint_20260319.md`
- `docs/adr/taskmaster/ADR-20260319-hierarchical-agent-bus.json`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/v1/router/status`
  - `curl -fsS http://127.0.0.1:8484/health`
- testes:
  - `go test ./internal/agent ./internal/tools -count=1`
  - `go test ./... -count=1`
- scripts:
  - smoke SQL local com `SELECT ... FOR UPDATE SKIP LOCKED`
- fallback:
  - manter task store atual em SQLite
  - rodar sem assistance queue

## Rollout

1. criar esquema mínimo do bus
2. implementar claim + lease
3. expor channels/threads/tarefas
4. implementar assistance queue
5. indexar resumos no Qdrant

## Rollback

- desligar o bus novo
- voltar para task store local atual
- manter Qdrant apenas como memória já existente

## Evidência esperada

- tarefa é claimed uma vez por lease válido
- agente ocioso consegue ajudar por fila secundária
- thread gera resumo e vai para Qdrant
- dashboard mostra tarefas, agentes e presença

## Pendências / bloqueios

- depende da definição final do schema PostgreSQL e da política de autoridade por agente

## Progresso registrado

- schema base do bus já entrou no `SQLiteTaskStore` para acelerar a tradução futura para `PostgreSQL`
- bases adicionadas:
  - `swarm_channels`
  - `swarm_threads`
  - `swarm_thread_messages`
  - `assistance_tasks`
- helpers iniciais de canal, thread, mensagens e claim de assistência já existem
- `go test ./internal/agent -count=1` passou
