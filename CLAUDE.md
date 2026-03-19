---
description: Adaptador de execução para o Claude Code CLI.
---

# 🤖 CLAUDE.md — Adaptador de Execução

> **IMPORTANTE**: Este arquivo é um adaptador fino. A autoridade máxima reside em [AGENTS.md](./AGENTS.md).

<contract>
## 📜 Regras de Engajamento
1. **Hierarquia**: Respeite o `PRD.md` e `AGENTS.md` como fontes de verdade.
2. **Higiene**: Sincronize o contexto via MCP `ai-context` após cada tarefa.
3. **Isolamento**: Priorize worktrees para implementações não triviais.
4. **ADR por Slice**: Não implemente mudança estrutural sem ADR ou backlog de slice registrado em `docs/adr/`.
</contract>

<tips>
Use `/p` para planejar com Opus e `/i` para implementar com Sonnet para o melhor balanço de custo/performance.
</tips>

## Links obrigatórios

- [AGENTS.md](./AGENTS.md)
- [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md)
- [ADR Index](./docs/adr/README.md)
- [MODEL.md](./MODEL.md)
- [plan.md](./plan.md)
