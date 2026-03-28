# ADR-20260325-Sovereign-Bibliotheca-v2-and-Git-Cleanup

## Contexto e Problema

O ecossistema Aurélia estava fragmentado entre Go (core), Node (Obsidiam/OpenClaw) e scripts Bash isolados. Além disso, o monorepo sofria com a "salada" de mais de 10.000 arquivos não rastreados (2.4GB) vindos de competências externas, o que degradava a performance do Git e do IDE.

## Decisão

Implementamos a **Sovereign-Bibliotheca v2**, uma camada de abstração unificada em Bash para orquestração agnóstica de agentes.

### 1. Arquitetura Unificada
- **Localização**: `homelab-bibliotheca/lib/`
- **Módulos**:
  - `config.sh`: Centralização de segredos e caminhos.
  - `memory.sh`: Sincronização SQLite ↔ Qdrant/Supabase.
  - `notes.sh`: Integração híbrida Obsidian CLI + Markdown Nativo.
  - `comms.sh`: Interface multi-bot Telegram via Gateway.
  - `skills.sh`: Indexador de competências do repositório OpenClaw.

### 2. Higiene Industrial (Git Cleanup)
- Atualização do `.gitignore` para ignorar o diretório `homelab-bibliotheca/skills/open-claw/skills/` (2.4GB).
- Manutenção do `skills-registry.json` como índice leve para busca semântica.

### 3. Governança (Regra 15)
- Documentação formal em `.agent/rules/15-sovereign-bibliotheca.md`.
- Divisão clara de responsabilidades: Go (Performance), Node (Automação/UI), Bash (Orquestração/CLI).

## Consequências

- **Positivas**: Repositório leve, busca semântica funcional, interface unificada para todos os agentes.
- **Negativas**: Dependência de scripts Bash para a cola entre tecnologias, exigindo manutenção manual do `sync.sh`.

---
**Status:** ✅ Decidido e Implementado (2026-03-25)
**Autoridade:** Antigravity Gemini p/ Aurélia


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
