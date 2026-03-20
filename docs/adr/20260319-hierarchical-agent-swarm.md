---
title: Swarm hierárquico da Aurélia com Postgres + Qdrant
status: accepted
created: 2026-03-19
owner: codex
---

# Contexto

O repositório já possui:

- autoridade central da Aurélia
- gateway de modelos
- voz, browser e runtime local
- ADR por slice

Faltava formalizar como um futuro enxame de agentes deve operar sem virar caos.

O objetivo é permitir:

- tarefas compartilhadas em dashboard
- ajuda entre agentes ociosos
- memória durável entre agentes
- brainstorm frequente com hierarquia

# Decisão

Adotar o seguinte desenho para o futuro swarm:

1. `Aurélia` permanece como autoridade arquitetural e operacional abaixo apenas dos humanos
2. `PostgreSQL` será a fonte de verdade operacional do swarm
3. `Qdrant` será memória semântica derivada, nunca a verdade primária
4. o swarm terá agentes fixos e efêmeros
5. agentes ociosos só podem ajudar via fila de assistência controlada
6. brainstorm multiagente será sequencial e arbitrado, nunca paralelo caótico
7. o dashboard de controle será nativo, com `Grafana` mantido para observabilidade

# Consequências

## Positivas

- reduz caos de coordenação
- dá auditabilidade completa
- facilita retomada e handoff
- permite múltiplas GPUs no futuro sem quebrar a hierarquia
- preserva a autoridade da Aurélia

## Restrições

- `Qdrant` não pode ser usado como fila ou verdade do runtime
- agentes efêmeros não podem subir de autoridade
- toda ajuda entre agentes deve respeitar lease, capability e policy
- qualquer runtime de swarm precisa nascer em slice própria e ADR própria

# Referências externas observadas em 2026-03-19

- LangGraph Platform: `assistants`, `threads`, `crons`, controle de acesso
  - https://docs.langchain.com/langgraph-platform/auth
- AutoGen Studio: dashboard de times/agentes e runs
  - https://www.microsoft.com/en-us/research/wp-content/uploads/2024/08/AutoGen_Studio-12.pdf
- CrewAI: observabilidade e disciplina de crews/flows
  - https://docs.crewai.com/

# Arquivos relacionados

- `AGENTS.md`
- `docs/REPOSITORY_CONTRACT.md`
- `docs/agent_swarm_dashboard_blueprint_20260319.md`
- `docs/adr/PENDING-SLICES-20260319.md`
- `plan.md`

# Validação

- blueprint do swarm versionado
- backlog oficial atualizado com slices do swarm
- cadeia de autoridade preservada sem contradição com `AGENTS.md`
