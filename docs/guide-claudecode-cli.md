---
description: Guia de uso do Claude Code CLI integrado ao ambiente Antigravity.
last-updated: 2026-03-17
---

# 🤖 Guia: Claude Code CLI

Instruções para extrair o máximo de performance dos subagentes Claude neste workspace.

<commands>
## ⚡ 1. Comandos de Controle
- `/bypass-on`: Ativa modo Full Auto (Dangerously Bypass).
- `/bypass-off`: Restaura Guardrails de segurança.
- `/r`, `/p`, `/i`, `/rv`: Atalhos para Researcher, Planner, Implementer e Reviewer.
</commands>

<important>
## 🔄 2. Higiene de Contexto (OBRIGATÓRIO)
Ao finalizar qualquer tarefa, você deve garantir que o `.context/` está sincronizado.
- **Ação**: Execute o comando de sincronização do MCP `ai-context`.
- **Por que?**: Garante que o próximo agente tenha a memória real do projeto.
</important>

---
*Mantenha a disciplina de Worktrees para tarefas Tier B e C.*
