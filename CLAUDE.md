---
description: Adaptador de execução para o Claude Code CLI.
---

# 🤖 CLAUDE.md — Adaptador de Execução

> **IMPORTANTE**: Este arquivo é um adaptador fino. A autoridade máxima reside em [AGENTS.md](./AGENTS.md).

<contract>
## 📜 Regras de Engajamento
1. **Hierarquia**: Respeite `AGENTS.md` (soberano), `REPOSITORY_CONTRACT.md` (governança) e `docs/adr/` (decisões estruturais) como fontes de verdade.
2. **Autoridade da Aurélia**: Execute sob a direção arquitetural e operacional da Aurélia; não dispute governança com ela.
3. **Higiene**: Sincronize o contexto via `/sincronizar-tudo` ou `bash scripts/sync-ai-context.sh` após cada tarefa (obrigatório ao fechar slice).
4. **Isolamento**: Priorize worktrees para implementações não triviais.
5. **ADR por Slice**: Não implemente mudança estrutural sem ADR ou backlog de slice registrado em `docs/adr/`. Use `bash scripts/validate-adr-semparar.sh` para validar conformidade.
6. **Governança Industrial**: Consulte [ADR-20260319-Polish-Governance-All](./docs/adr/ADR-20260319-Polish-Governance-All.md) para secrets, dados, rede, ops, observabilidade. Use skill `/governance-polish` para automatizar fases.
7. **Memory Sync → Vector DB**: Aurelia bot acessa memória (code history) via `/memory-sync-vector-db`. Crons automáticos sincronizam Markdown → Qdrant + Postgres. Permite LLMs pequenas sem web. Ver [memory-sync-architecture.md](./docs/memory-sync-architecture.md).
</contract>

<workflow>
## 🔄 Fluxo de Slice Estrutural

Para slices não-triviais (mudanças arquiteturais, integração multi-agente):

1. **Abrir**: `bash scripts/adr-slice-init.sh <slug> --title "Título"`
2. **Preencher**: ADR + JSON taskmaster com contexto, decision, smoke tests
3. **Validar**: `bash scripts/validate-adr-semparar.sh` (deve passar)
4. **Executar**: Em worktree isolada com handoff estruturado
5. **Fechar**: Atualize JSON com evidência, rode `/sincronizar-tudo` (ou `bash scripts/sync-ai-context.sh`)
6. **Commitar**: Com referência a sync e validação

Documentação: `.agents/workflows/adr-semparar-governance.md`
Skill: [sync-ai-context](./.agents/skills/sync-ai-context/SKILL.md)
</workflow>

<tips>
- Use `/p` para planejar com Opus e `/i` para implementar com Sonnet para melhor custo/performance
- Consulte `adr-semparar-agents-md-conformance.md` se houver dúvida sobre conformidade com AGENTS.md
</tips>

## Hierarquia de Referência

**Antes de qualquer trabalho, leia nesta ordem:**
1. [AGENTS.md](./AGENTS.md) ← Autoridade suprema
2. [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md) ← Contrato do repositório
3. [ADR Index](./docs/adr/README.md) + [TASKMASTER-INDEX](./docs/adr/TASKMASTER-INDEX.md) ← Decisões estruturais e priorização
4. [.agents/workflows/adr-semparar-status.md](./.agents/workflows/adr-semparar-status.md) ← Status real das slices

## Links obrigatórios

- [AGENTS.md](./AGENTS.md)
- [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md)
- [ADR Index](./docs/adr/README.md)
- [TASKMASTER-INDEX](./docs/adr/TASKMASTER-INDEX.md)
- [MODEL.md](./MODEL.md)
- [plan.md](./plan.md)
