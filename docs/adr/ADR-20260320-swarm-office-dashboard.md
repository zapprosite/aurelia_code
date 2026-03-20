---
description: CRM/Office Dashboard para Swarm de Agentes Colaborativos (Híbrido Go + Supabase + Qdrant + SQLite).
status: proposed
---

# ADR 20260320-swarm-office-dashboard

## Status

- Proposto

## Slice

- nome do slice: swarm-office-dashboard
- owner: Antigravity / Aurelia
- branch/worktree: agent-to-agent
- data: 2026-03-20

## Contexto

Necessidade de uma interface operacional ("Office Dashboard") que visualize os agentes como funcionários em um "Escritório Virtual". O sistema deve gerenciar memória em três camadas (quente, morna, fria) para garantir performance e profundidade de contexto.

## Decisão

### 1. Memória Híbrida (MemoryManager)
- **SQLite (Quente)**: Estado imediato e rascunhos.
- **Supabase/pgvector (Morna)**: Memória estruturada de longo prazo e busca híbrida.
- **Qdrant (Fria/Episódica)**: Experiências massivas e busca semântica pura.

### 2. Dashboard de Colaboração
- Frontend React com **ReactFlow** para visualização do grafo de ajuda mútua.
- Estética "Clean Corporate Future".
- WebSocket Hub em Go para atualizações em tempo real.

### 3. Integrações
- **Supabase**: Migrations para pgvector e tabelas de colaboração.
- **Qdrant**: Coleção `swarm_experiences` com ingestão passiva.

## Arquivos Afetados

- [NEW] `internal/memory/manager.go` (Roteador de memória)
- [NEW] `supabase/migrations/20260320_swarm_tables.sql`
- [MODIFY] `internal/agent/hub.go` (Broadcasting rico para o dashboard)
- [NEW] `web/dashboard/` (Estrutura frontend)

## Regras / Guardrails

- **Ingestão Passiva**: Toda interação bem-sucedida deve ser enviada ao Qdrant de forma assíncrona.
- **Tri-Banco**: Proibir duplicação indevida; cada camada tem seu propósito claro.

## Testes obrigatórios

- Integração: Testar busca vetorial no Supabase via query Go.
- Performance: Validar latência da busca híbrida < 200ms.
- E2E: Simular ajuda entre agentes e verificar visualização no ReactFlow.
