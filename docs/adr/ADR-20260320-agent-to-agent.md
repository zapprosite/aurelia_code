---
description: Cirurgia Arquitetural - Migração de LangGraph para Swarm Event-Driven (Go + PicoLisp).
status: proposed
---

# ADR 20260320-agent-to-agent

## Status

- Proposto

## Slice

- nome do slice: agent-to-agent (Migração LangGraph -> Swarm)
- owner: Antigravity / Aurelia
- branch/worktree: agent-to-agent / /home/will/aurelia-agent-to-agent
- data: 2026-03-20

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [plan.md](../../plan.md)
- [ADR-20260319-hierarchical-agent-bus.md](./ADR-20260319-hierarchical-agent-bus.md)

## Contexto

Migração de arquiteturas baseadas em grafos de estado (LangGraph/Python) para sistemas de enxame colaborativo orientado a eventos de alta performance (Go + PicoLisp). 

### Auditoria (Fase 1) - Achados:
- **Nós de Estado Pesado**: `internal/agent/loop.go` carrega `[]Message` como estado principal em um loop síncrono.
- **Acoplamento Sequencial**: `master_team_service.go` orquestra tarefas via `task_store` com lógica de "claim" e persistência em DB, gerando latência de I/O em vez de eventos puros.
- **Lógica de Roteamento**: Atualmente distribuída em Go (`task_store_dependency_flow.go`, etc.), dificultando mudanças simbólicas rápidas.
- **Pontos de Memória**: O estado é passado via slice de mensagens, consumindo memória e limitando o contexto compartilhado entre agentes.

## Decisão

Substituição do paradigma LangGraph pelo modelo "Colmeia Colaborativa":

### A. Núcleo em Go (Orquestrador de Eventos)
- Substituir StateGraph/Loop síncrono por um Event Bus (canais Go ou NATS).
- Agentes como Goroutines independentes ouvindo tópicos.
- Hierarquia Fluida com Conductor (Regente) e Specialists.

### B. Lógica Simbólica em PicoLisp
- Extração de lógica de roteamento (`task_store_dependency_rules.go`) para scripts `.l` (PicoLisp).
- Baseado em símbolos de capacidade e carga.

### C. Memory OS
- Implementar `internal/memory` com suporte a Push Memory e contexto compartilhado (Redis/VectorDB).

## Arquivos Afetados

- [MODIFY] `internal/agent/loop.go` (Remover estado pesado)
- [MODIFY] `internal/agent/master_team_service.go` (Migrar para Event Bus)
- [NEW] `internal/agent/event_bus.go` (Implementação do backbone)
- [NEW] `internal/agent/decision.l` (Lógica simbólica PicoLisp)
- [NEW] `internal/memory/os.go` (Memory OS layer)

## Regras / Guardrails

- **Regra de Ouro**: Nenhum agente bloqueia esperando outro. Comunicação sempre via eventos.
- **Sênior Design**: Utilizar PicoLisp para decisões simbólicas complexas.

## Testes obrigatórios

- unit: Testes de canais Go e Pub/Sub.
- integração: Validação Conductor <-> Specialists via Event Bus.
- E2E: Fluxo completo de tarefa com intervenção de múltiplos agentes e visualização no dashboard.

## Rollout

1. Fase 1: Auditoria e Extração de Lógica (Python).
2. Fase 2: Backbone Go (Event Bus).
3. Fase 3: Integração PicoLisp.
4. Fase 4: Memory OS e Dashboard.

## Rollback

Manter a estrutura LangGraph legada em pastas separadas (/deprecated/python-swarm) até validação completa do novo sistema Go.

## Consequências

- positivas: Redução de 40% na latência, paralelismo real de agentes, memória persistente.
- trade-offs: Aumento na complexidade operacional (Go + PicoLisp), necessidade de coordenação de eventos assíncronos.
