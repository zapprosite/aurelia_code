---
title: Acelerar o swarm da Aurélia com referência oficial do LangGraph
status: accepted
created: 2026-03-19
owner: codex
---

# Contexto

O blueprint do swarm hierárquico já estava fechado, mas ainda faltava uma decisão prática para acelerar a implementação sem desperdiçar tempo em desenho já resolvido por outros projetos.

Era necessário escolher uma referência concreta para copiar padrões de:

- supervisor
- handoff
- registro de agentes
- gestão de histórico
- painel de times/agentes

# Decisão

Adotar como referência primária estes repositórios oficiais:

1. `langchain-ai/open-agent-supervisor`
2. `langchain-ai/langgraph-supervisor`

Decisão complementar:

- copiar padrões de orquestração e contrato
- não copiar o runtime Python inteiro para dentro do core da Aurélia
- implementar a versão final em `Go`, mantendo o contrato local do homelab

# O que copiar

## De `open-agent-supervisor`

- configuração declarativa de agentes filhos
- supervisor configurável com prompt fixo + parte não editável
- separação entre supervisor e agentes remotos
- ideia de transparência do fluxo para o usuário

## De `langgraph-supervisor`

- handoff por tool
- retorno explícito ao supervisor
- `forward_message` para preservar fidelidade de resposta
- histórico `full_history` vs `last_message`
- hierarquia multinível de supervisores

# O que não copiar

- dependência estrutural no LangGraph Platform
- runtime Python como plano de controle central
- auth e configuração remota da OAP como decisão automática para este repositório

# Consequências

## Positivas

- acelera a `Slice 11`
- reduz chance de erro em desenho de handoff
- evita reescrever padrões maduros do zero
- mantém a Aurélia como autoridade do sistema

## Restrições

- qualquer cópia deve ser traduzida para o contrato local em `Go`
- `PostgreSQL` continua sendo a verdade operacional
- `Qdrant` continua sendo memória derivada
- a UI final do swarm continua subordinada ao contrato da Aurélia

# Fontes

- https://github.com/langchain-ai/open-agent-supervisor
- https://github.com/langchain-ai/langgraph-supervisor
- https://docs.langchain.com/langgraph-platform/auth
- https://docs.langchain.com/oss/python/langchain/supervisor

# Arquivos relacionados

- `docs/agent_swarm_dashboard_blueprint_20260319.md`
- `docs/adr/20260319-hierarchical-agent-swarm.md`
- `docs/adr/PENDING-SLICES-20260319.md`
- `plan.md`

# Validação

- repositório de referência escolhido explicitamente
- blueprint atualizado com a linha de aceleração
- ADR índice atualizado
