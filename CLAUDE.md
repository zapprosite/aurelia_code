---
description: Adaptador de execução para o Claude Code CLI.
---

# CLAUDE.md — Adaptador Claude Code

> **IMPORTANTE**: Este arquivo é um adaptador fino. A autoridade máxima reside em [AGENTS.md](./AGENTS.md).

<contract>
## Governança Claude Code
1. **Hierarquia**: Respeite `AGENTS.md` (soberano), `REPOSITORY_CONTRACT.md` e `docs/adr/` como fontes de verdade.
2. **Autoridade da Aurélia**: Execute sob a direção arquitetural e operacional da Aurélia.
3. **Língua**: Mantenha documentação e planos em **Português (BR)**.
4. **Isolamento**: Prefira worktrees para implementações não triviais.
5. **ADR por Slice**: Não implemente mudança estrutural sem registro no [ADR.md](./docs/ADR.md).
6. **Planejamento**: Planos de implementação vão em `.context/plans/<slice>/`, não na raiz.
7. **Memory Sync**: Sincronize `.context/` após mudanças estruturais via `sync-ai-context`.
8. **Permissões**: Tier A/B auto-approve. Tier C (sudo) com log obrigatório. `sudo=1`.
9. **Skills e Workflows**: Use `.agents/skills/` e `.agents/workflows/` como ecossistema compartilhado.
</contract>

## Links obrigatórios

- [AGENTS.md](./AGENTS.md)
- [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md)
- [ADR Index](./docs/ADR.md)
- [ADR Histórico](./docs/ADR-historico.md)
