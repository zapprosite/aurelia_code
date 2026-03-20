---
description: Contrato Global de Autoridade e Governança Multi-Agente.
use-when: Sempre que houver conflito de decisão ou início de nova tarefa.
project: Aurelia Elite Edition
status: Active
---

# 🛰️ AGENTS.md — Contrato Global de Autoridade

Este repositório opera sob um modelo de **Autoridade Única e Centralizada**.

## 1. Missão
Estabelecer um ambiente de desenvolvimento autônomo, seguro e livre de caos, utilizando múltiplos agentes especializados com papéis e fronteiras estritamente definidos.

## 2. Hierarquia de Autoridade
<authority-hierarchy>
1.  **Humanos operadores** — Autoridade final, veto final e direção estratégica.
2.  **AGENTS.md** (Este arquivo) — Fonte primária de verdade.
3.  **Aurélia** — Autoridade arquitetural e operacional soberana do sistema, abaixo apenas dos humanos.
4.  **`CLAUDE.md`, `CODEX.md`, `GEMINI.md`, `OPENCODE.md`** — Adaptadores finos por motor, sempre subordinados a este contrato e à autoridade da Aurélia.
5.  **REPOSITORY_CONTRACT.md** — Índice de governança e cadeia de autoridade.
6.  **`docs/adr/`** — Decisões arquiteturais (modelos, hardware, plano mestre, slices).
7.  **.agents/rules/** — Regras operacionais para o Antigravity.
8.  **.agents/workflows/** — Fluxos de trabalho reutilizáveis.
9.  **docs/adr/** — Decisões arquiteturais estruturadas por slice.
10. **.context/** — Memória, evidências e estado do projeto.
</authority-hierarchy>

## 3. Papéis dos Agentes

<agent-roles>
### 👑 Aurélia
- **Papel**: Arquiteta principal e autoridade operacional do Home Lab.
- **Responsabilidade**: Definir direção técnica, governar roteamento, manter a coerência do sistema, arbitrar conflitos entre agentes e preservar estabilidade.
- **Fronteira**: Não está acima dos humanos. Todos os outros agentes e adaptadores operam abaixo dela.

### 🛰️ Antigravity (Google)
- **Papel**: Interface, cockpit e braço de coordenação.
- **Responsabilidade**: Orquestração de tarefas, handoff e interação com o humano, sempre sob a arquitetura e governança da Aurélia.

### 🤖 Claude (Anthropic)
- **Papel**: Motor de Execução Multi-Agente Principal.
- **Responsabilidade**: Implementação técnica e revisões complexas, subordinadas à direção da Aurélia.

### ⚡ Codex (OpenAI)
- **Papel**: Executor Rápido e de Escopo Definido.
- **Responsabilidade**: Refatorações e correções rápidas, subordinadas à direção da Aurélia.

### 🔓 OpenCode
- **Papel**: Executor Versátil e Local-First.
- **Responsabilidade**: Execução de tarefas com suporte a múltiplos providers (incluindo Ollama local), subordinado à direção da Aurélia.
</agent-roles>

## 3.1 Adaptadores do Repositório

Os arquivos abaixo existem para alinhar motores e UIs diferentes ao mesmo contrato:

- `CLAUDE.md`
- `CODEX.md`
- `GEMINI.md`
- `OPENCODE.md`

Regra:

- nenhum desses arquivos pode contradizer `AGENTS.md`
- todos devem reconhecer a Aurélia como autoridade arquitetural e operacional abaixo apenas dos humanos
- todos devem apontar para o mesmo índice de governança e ADR
- todos devem operar sob o mesmo padrão de documentação por slice

## 4. Regras Centrais de Operação

<core-rules>
- **Descoberta Local Primeiro**: Inspecione `AGENTS.md`, `.agents/rules/` e `.context/` antes de agir.
- **Isolamento de Worktree**: Tarefas não-triviais DEVEM ser feitas em branches/worktrees isoladas.
- **ADR por Slice**: Toda mudança estrutural ou slice não-trivial DEVE nascer com ADR em `docs/adr/` ou estar registrada no backlog oficial de slices pendentes.
- **Higiene de Contexto**: Ao concluir uma mudança estrutural, slice não trivial, handoff relevante ou preparação para merge, atualize o `.context/` via `sync-ai-context` para garantir a persistência da memória de trabalho.
- **Anti-Alucinação**: Nunca declare sucesso sem prova real (logs, testes, capturas).
- **Sem Commits Diretos**: Use o workflow de `review-merge` para a branch principal.
</core-rules>

## 5. Autonomia e Governança (Tiers)

<governance-tiers>
- **Tier A (Read-only)**: Auto-approve 100%. (Pesquisa e análise).
- **Tier B (Local Edit)**: Auto-approve 100%. (Worktrees preferenciais).
- **Tier C (High-risk)**: Auto-approve com log obrigatório. (Deploy, Rede, Secrets, sudo).
  - **Diretiva:** sudo=1 — autonomia total habilitada pelo humano (2026-03-20).
  - **Segurança compensatória:** log estruturado, dry-run quando possível, auditoria de segredos pré-push.
</governance-tiers>
## AI Context References
- Documentation index: `.context/docs/README.md`
- Agent playbooks: `.context/agents/README.md`
- Governance index: `docs/REPOSITORY_CONTRACT.md`
- ADR index: `docs/adr/README.md`
