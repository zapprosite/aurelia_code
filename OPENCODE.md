---
description: Adaptador de execução para o OpenCode CLI.
---

# 🔓 OPENCODE.md — Adaptador de Execução

> **IMPORTANTE**: Este arquivo é um adaptador fino. A autoridade máxima reside em [AGENTS.md](./AGENTS.md).

<contract>
## 🔓 Governança OpenCode
1. **Hierarquia**: Respeite `AGENTS.md` (soberano), `REPOSITORY_CONTRACT.md` (governança) e `docs/adr/` (decisões estruturais) como fontes de verdade.
2. **Autoridade da Aurélia**: Execute sob a direção arquitetural e operacional da Aurélia; não dispute governança com ela.
3. **Língua**: Mantenha a documentação e planos em **Português (BR)**.
4. **Isolamento**: Priorize worktrees para implementações não triviais.
5. **ADR por Slice**: Não implemente mudança estrutural sem ADR ou backlog de slice registrado em `docs/adr/`.
6. **Governança Industrial**: Consulte [ADR-20260319-Polish-Governance-All](./docs/adr/ADR-20260319-Polish-Governance-All.md). Coordene handoff: humano → Fase 1 (CRITICAL) → skill `/governance-polish` → Fases 2-4.
7. **Memory Sync**: Aurelia bot consulta code history via Qdrant + Postgres (sem web). Script `scripts/memory-sync-fiscal.sh` sincroniza memória. Valide com: `bash scripts/memory-sync-fiscal.sh --mode validate`.
8. **Permissões**: Auto-approve para edits e bash. Tier C (sudo) com log obrigatório. Diretiva: `sudo=1`.
</contract>

## Links obrigatórios

- [AGENTS.md](./AGENTS.md)
- [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md)
- [ADR Index](./docs/adr/README.md)
- [Política de Modelos](./docs/adr/ADR-20260320-politica-modelos-hardware-vram.md)
- [Plano Mestre](./docs/adr/ADR-20260320-plano-mestre-jarvis-local-first.md)
