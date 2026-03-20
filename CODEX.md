---
description: Adaptador de implementação rápida para o Codex CLI.
---

# ⚡ CODEX.md — Adaptador de Implementação

> **IMPORTANTE**: Este arquivo é um adaptador fino. A autoridade máxima reside em [AGENTS.md](./AGENTS.md).

<contract>
## 🚀 Missão Codex
1. **Velocidade**: Atue como o braço de execução rápida para refatorações e correções.
2. **Autoridade da Aurélia**: Não tome decisão arquitetural contra a direção da Aurélia.
3. **Alinhamento**: Não tome decisões arquiteturais sem consultar o `Architect Specialist`.
4. **Fronteira**: Mantenha-se dentro da worktree designada.
5. **ADR por Slice**: Se a mudança sair do escopo local/rápido, exija ADR em `docs/adr/` antes de continuar.
6. **Governança Industrial**: Consulte [ADR-20260319-Polish-Governance-All](./docs/adr/ADR-20260319-Polish-Governance-All.md). Fases 2-4 usam skill `/governance-polish --phase 2` (após humano concluir Fase 1).
7. **Memory Sync — Executar Crons**: Script `scripts/memory-sync-fiscal.sh` sincroniza memória → Qdrant + Postgres. Executado via systemd timers (5, 15, 60, 1440 min). Valide com: `bash scripts/memory-sync-fiscal.sh --mode validate`.
</contract>

## Links obrigatórios

- [AGENTS.md](./AGENTS.md)
- [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md)
- [ADR Index](./docs/adr/README.md)
- [Política de Modelos](./docs/adr/ADR-20260320-politica-modelos-hardware-vram.md)
- [Plano Mestre](./docs/adr/ADR-20260320-plano-mestre-jarvis-local-first.md)
