---
description: Visão Geral da Arquitetura do Workspace Multi-Agente.
last-updated: 2026-03-17
---

# 🏛️ Arquitetura do Workspace

## 1. Topologia de Fluxo (BMAD)
O sistema segue o modelo **BMAD (Business, Model, Architect, Developer)** orquestrado por agentes:

<flow-diagram>
[Usuário] ➡️ [/pm] ➡️ (PRD.md)
              [PM] ➡️ (Tech Spec)
          ➡️ [/architect] ➡️ (ADR/Docs)
              [Architect] ➡️ (Implementation Plan)
          ➡️ [/dev] ➡️ (Codebase)
              [Developer] ➡️ (Artifacts)
          ➡️ [/qa] ➡️ (Tests/Audit)
              [QA] ➡️ [Sincronia de Contexto]
          ➡️ [MCP ai-context] ➡️ [Sucesso/Merge]
</flow-diagram>

## 2. Higiene de Contexto (CRÍTICO)
Para evitar que o `.context/` fique desatualizado em relação ao código real:
- **Regra**: Ao finalizar qualquer feature (antes do merge), execute o comando de sincronização do MCP `ai-context`.
- **Objetivo**: Atualizar o mapa do código (`codebase map`) e garantir que o próximo agente que assumir o repositório tenha a visão real do estado atual.

## 2. Camadas de Contexto
A densidade de contexto é dividida em três níveis de persistência:

### Camada 0: Governança (Imutável por Sessão)
- **Localização**: `.agents/rules/`, `AGENTS.md`.
- **Propósito**: Define as "leis" de física do workspace.

### Camada 1: Estrutura (Estática)
- **Localização**: `docs/`, `architecture.md`.
- **Propósito**: Conhecimento técnico e decisões históricas.

### Camada 2: Execução (Efêmera)
- **Localização**: `.context/`, `.planning/`.
- **Propósito**: Estado atual, memória de curto prazo e planos ativos.

## 3. Integração de Ferramentas
- **Antigravity IDE**: Interface primária de supervisão.
- **Claude Code**: Motor de execução pesado (Subagentes).
- **ai-context (MCP)**: Indexador de símbolos e preenchimento de scaffolding.
