> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

---
title: "Aurelia Integrated Planning & Verification — Phase 15"
slug: aurelia-integrated-planning-verification
status: proposed
date: 2026-03-22
decision-makers: [humano, aurelia, antigravity]
priority: P15
tags: [planning, verification, memory, dashboard]
---

# ADR-20260322: Aurelia Integrated Planning & Verification (Phase 15)

## Contexto

Embora a Aurélia já possua um loop PREV (Plan-Review-Execute-Verify), o "Plano" hoje é apenas uma instrução em texto no prompt. Não há um objeto estruturado que o humano possa ler, aprovar ou rejeitar via Dashboard. Além disso, a memória cognitiva (Knowledge Items) ainda depende muito de intervenção manual ou scripts fiscais externos.

## Decisão

Implementar a "Mão do Arquiteto": um sistema de governança onde cada tarefa complexa gera um `ActionPlan` digital, visualizável no Cockpit, e cuja conclusão dispara uma atualização automática da base de conhecimento.

### Sub-Slices Propostos

#### 1. Sub-1: ActionPlan Core (Go)
**Arquivo:** `internal/agent/planner.go`
- Definir `struct ActionPlan` com `Steps`, `RiskLevel`, `EstimatedTime` e `BackoutPlan`.
- Criar tool `propose_plan` que emite o plano como um evento SSE tipo `agent_plan`.
- Integrar no `agent.Loop` para que, em `PHASE_PLANNING`, o agente SEJA OBRIGADO a chamar esta tool antes de `set_phase(EXECUTION)`.

#### 2. Sub-2: Cockpit Plan Viewer (React)
**Arquivo:** `frontend/src/components/dashboard/PlanViewer.tsx`
- Nova aba no Dashboard para visualizar o plano ativo.
- Botões de `Approve Plan` e `Reject/Revise`.
- Endpoint em Go para segurar a execução do agente até o sinal de aprovação (Tier C).

#### 3. Sub-3: Autonomous KI Persistence (Memory)
**Skill:** `/memory-persist`
- Tool que lê o `git diff` da tarefa atual e gera um novo Knowledge Item (KI) ou atualiza um existente em `.context/docs/`.
- Garante que a "Memória" da Aurélia acompanhe o crescimento do código sem lag documental.

## Consequências

- **Segurança Sênior**: Tarefas de alto risco (sudo=1) terão um gate visual de aprovação.
- **Transparência**: O Dashboard deixa de ser um log e passa a ser uma ferramenta de co-pilotagem real.
- **Soberania Documental**: O repositório torna-se autossuficiente em sua própria explicação técnica.

## Smoke Test Planejado

1. Enviar prompt: "Crie um novo módulo de backup para o Qdrant".
2. Observar a Aurélia gerando um `ActionPlan` no Dashboard.
3. Aprovar o plano e observar a execução.
4. Validar que um novo KI foi gerado após o commit.
