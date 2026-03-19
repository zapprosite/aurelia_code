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
1.  **AGENTS.md** (Este arquivo) — Fonte primária de verdade.
2.  **`CLAUDE.md`, `CODEX.md`, `GEMINI.md`, `MODEL.md`** — Adaptadores finos por motor, sempre subordinados a este contrato.
3.  **PRD.md** — Intenção do projeto e roadmap.
4.  **.agents/rules/** — Regras operacionais para o Antigravity.
5.  **.agents/workflows/** — Fluxos de trabalho reutilizáveis.
6.  **.context/** — Memória, evidências e estado do projeto.
</authority-hierarchy>

## 3. Papéis dos Agentes

<agent-roles>
### 🛰️ Antigravity (Google)
- **Papel**: Orquestrador, Supervisor e Interface de usuário.
- **Responsabilidade**: Divisão de tarefas, planejamento e coordenação.

### 🤖 Claude (Anthropic)
- **Papel**: Motor de Execução Multi-Agente Principal.
- **Responsabilidade**: Implementação técnica e revisões complexas.

### ⚡ Codex (OpenAI)
- **Papel**: Executor Rápido e de Escopo Definido.
- **Responsabilidade**: Refatorações e correções rápidas.
</agent-roles>

## 3.1 Adaptadores do Repositório

Os arquivos abaixo existem para alinhar motores e UIs diferentes ao mesmo contrato:

- `CLAUDE.md`
- `CODEX.md`
- `GEMINI.md`
- `MODEL.md`

Regra:

- nenhum desses arquivos pode contradizer `AGENTS.md`
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
- **Tier B (Local Edit)**: Auto-approve condicional (Worktrees).
- **Tier C (High-risk)**: Aprovação Humana OBRIGATÓRIA (Deploy, Rede, Secrets).
</governance-tiers>
## AI Context References
- Documentation index: `.context/docs/README.md`
- Agent playbooks: `.context/agents/README.md`
- Governance index: `docs/REPOSITORY_CONTRACT.md`
- ADR index: `docs/adr/README.md`
