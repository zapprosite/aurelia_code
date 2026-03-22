---
title: "ADR-20260321-planning-loop-prev"
status: proposed
date: 2026-03-21
---

# ADR: Implementação do Planning Loop (PREV)

## Contexto e Problema
Atualmente, os agentes na `agent/loop.go` operam em um ciclo reativo e sequencial (Thought ➔ Action). Em arquiteturas sêniores de AI Engineering, agentes que têm acesso direto ao SO (sudo=1) precisam de previsibilidade metodológica. Sem um planejamento delimitado, o bot pode disparar N ferramentas de uma só vez numa direção equivocada. Módulos paralelos como o ULTRATRINK (Dashboard UI) e `squad.go` já esperam visualmente por eventos mais explícitos sobre *em que fase* o agente se encontra, para feedback na UI.

## Decisão (PREV Architecture)
Vamos injetar formalmente uma Máquina de Estados (State Machine) ou delimitadores de fase no Loop da Aurelia:
1. **Plan (P):** O agente faz RAG ou pesquisas de contexto, não disparando comandos diretos que mutem estado, e emite um "Plano de Ação" (Execution Plan / ADR local).
2. **Review (R):** (Opcional) A Swarm pausa ou avalia as regras para ver se o plano tem consentimento, pedindo handoff ou aprovação para seguir mediante `RunOptions` ou policy guardrails.
3. **Execute (E):** A implementação real do plano via ferramentas destrutivas (write/run).
4. **Verify (V):** A validação cruzada final e os testes para provar ao humano (Walkthrough/Smoke).

## Rollout e Design
- Adicionar ou adequar literais de Enum de `Phase` no `agent.Loop` (`squad.go` ou `event_bus`).
- Atualizar a emissão de SSE (`EmitPhaseChanged`) no `loop.go` para alimentar o React.
- Adequar os prompts mestres da Identidade para incentivar as ferramentas de planejamento.

## Rollback
Se o loop engasgar infinitamente no RAG, reverter as checagens e voltar para o Thought ➔ Action cru de único turno.
