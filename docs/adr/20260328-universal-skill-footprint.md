# 20260328 - Pegada Universal de Skills (Universal Skill Footprint)

## Status: Proposto

## Contexto
O ecossistema Aurélia evoluiu para uma arquitetura multi-agente e multi-camada. Atualmente, as "skills" (habilidades agentivas) residem primariamente em `.agent/skills/`, mas não são automaticamente propagadas para:
1. **Claude Code** (.claude/skills): Impedindo o Claude de usar o mesmo catálogo.
2. **Obsidian** (Bibliotheca): Dificultando a gestão humana e a visualização do catálogo.
3. **Qdrant/DB**: Impedindo a busca semântica cross-agente e a governança centralizada.

## Decisões
1. **Sovereign Mirroring**: Adotar o diretório `.agent/skills/` como a Fonte de Verdade (Source of Truth). Todas as adições ou deleções devem ser espelhadas para `.claude/skills/` e `homelab-bibliotheca/skills/aurelia/`.
2. **Indexing Semântico**: Cada skill (`SKILL.md`) deve ser fragmentada e indexada na collection `aurelia_skills` no Qdrant, utilizando o modelo de embedding local.
3. **Persistência de Metadados**: Utilizar o SQLite (`aurelia.db`) para rastrear o estado de sincronização e o "mapa de calor" de uso de cada skill.
4. **Orquestração Master-Skill**: O comando `/master-skill sync` será o gatilho único para disparar este fluxo de industrialização de habilidades.

## Consequências
- **Consistência Cross-Agente**: Claude, Antigravity e Codex compartilham o mesmo cérebro operacional.
- **Busca Semântica de Alta Performance**: Localização instantânea da ferramenta certa para cada tarefa.
- **Soberania de Contexto**: Toda a documentação de habilidades está disponível offline no Obsidian para auditoria humana.

## Referências
- Master Skill: `/home/will/.gemini/antigravity/skills/master-skill/SKILL.md`
- Governance: `AGENTS.md`
