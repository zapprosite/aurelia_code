# 📚 Aurélia Project Playbook

> **Auto-sync**: Este documento é gerado continuamente pelo sistema de memory-sync. Reflete CLAUDE.md + AGENTS.md + Skills + ADRs.

## 🤖 Hierarquia de Autoridade

1. **AGENTS.md** ← Soberano (agent roles, responsibilities)
2. **CLAUDE.md** ← Governance rules (execução, slices, ADRs)
3. **REPOSITORY_CONTRACT.md** ← Contrato de repositório
4. **docs/adr/** ← Decisões arquiteturais

---

## 📋 Contexto de Execução

### Regras de Engajamento

- ✅ **Autoridade**: Respeite AGENTS.md como fonte suprema
- ✅ **Sincronização**: Use `/sincronizar-tudo` após tarefas estruturais
- ✅ **ADR-First**: Não implemente mudanças estruturais sem ADR ou backlog
- ✅ **Isolamento**: Use worktrees para mudanças não-triviais
- ✅ **Memory-Sync**: Crons automáticos sincronizam Markdown → Qdrant + Postgres

### Skills Disponíveis (Telegram-Facing)

As seguintes skills estão SEMPRE disponíveis para o telegram bot:

- **governance-polish**: Automatizar fases de governança
- **sync-ai-context**: Sincronizar contexto AI via `/sincronizar-tudo`
- **adr-slice**: Workflow de ADR + slice estruturado
- **memory-sync-vector-db**: Integrar histórico em Qdrant (LLM pequenas)

---

## 🔄 Fluxo de Slice Estrutural

Para tarefas não-triviais (mudanças arquiteturais, multi-agente):

1. **Abrir Slice**: `bash scripts/adr-slice-init.sh <slug> --title "Título"`
2. **Preencher ADR**: Decisões + smoke tests no JSON
3. **Validar**: `bash scripts/validate-adr-semparar.sh`
4. **Executar**: Em worktree isolada
5. **Fechar**: Rode `/sincronizar-tudo` + commit com ref de sync

---

## 🎯 Responsabilidades de Agentes

### Claude (AI Assistant)
- Implementar sob direção de Aurelia
- Respeitando AGENTS.md como soberano
- Sincronizar contexto com `/sincronizar-tudo`
- Não disputar governança com Aurelia

### Aurelia (Orchestrator)
- Direção arquitetural e operacional
- Validar conformidade com AGENTS.md
- Autoridade na governança de projeto

---

## 💾 Memory & Vector DB

- **Arquivos**: Markdown em `docs/adr/`, `docs/governance/`, `.agents/skills/`
- **Sync**: Cron automático Markdown → Qdrant + Postgres
- **Benefício**: LLMs pequenas acesso a histórico sem web
- **Documentação**: `docs/memory-sync-architecture.md`

---

## 🔗 Referências Obrigatórias

| Documento | Propósito | Autoridade |
|-----------|----------|-----------|
| AGENTS.md | Role mapping e responsibilities | ⭐ Suprema |
| CLAUDE.md | Execution rules e governance | ⭐ Alta |
| REPOSITORY_CONTRACT.md | Contrato de repo | ⭐ Alta |
| ADR Index (docs/adr/) | Decisões arquiteturais | ⭐ Alta |
| TASKMASTER-INDEX | Priorização de slices | ⭐ Alta |

---

## 📌 Última Sincronização

- **Data**: 2026-03-23
- **Script**: `bash scripts/sync-ai-context.sh`
- **Status**: ✅ Auto-sync ativo (Qdrant + Postgres)

---

**NOTA**: Este playbook é regenerado automaticamente. Não edite manualmente — altere CLAUDE.md, AGENTS.md ou scripts de sync.
