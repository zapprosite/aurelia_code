# ADR Index

Índice canônico de ADRs vigentes. Dividido em três categorias honestas.

## Naming
- Padrão: `YYYYMMDD-slug.md`
- ADR estrutural obrigatória para: storage, memória, runtime, governança, segurança, rede, modelos, deploy
- **Conformidade Semparar**: ADRs estruturais têm par `.md + taskmaster/*.json`

---

## Implementada (código existe e foi validado)

- [20260327-markdown-brain-aurelia-code.md](20260327-markdown-brain-aurelia-code.md) — cérebro Markdown canônico para o `aurelia_code`
- [20260326-homelab-dashboard-native.md](20260326-homelab-dashboard-native.md) — absorção nativa do monitor do homelab no dashboard da Aurelia
- [20260325-slice-voice-capture-readiness.md](20260325-slice-voice-capture-readiness.md) — voice capture pipeline e config normalizado
- [20260325-slice-runtime-governance-enforcement.md](20260325-slice-runtime-governance-enforcement.md) — enforcement de payload canônico no Qdrant
- [20260325-slice-team-orchestration-honesty.md](20260325-slice-team-orchestration-honesty.md) — renomear swarm→team, contrato de coordenação
- [20260325-claude-code-native-installer-migration.md](20260325-claude-code-native-installer-migration.md) — migração do CLI Claude Code
- [20260325-openclaw-skill-vault-isolation.md](20260325-openclaw-skill-vault-isolation.md) — isolamento do vault de skills OpenClaw
- [20260324-multi-bot-dashboard.md](20260324-multi-bot-dashboard.md) — dashboard multi-bot com SSE e embedded JS
- [20260324-install-tavily-web-search.md](20260324-install-tavily-web-search.md) — integração Tavily web search
- [20260325-caixa-bot-persona.md](20260325-caixa-bot-persona.md) — persona do bot Caixa PF/PJ
- [20260325-controle-db-bot-governance.md](20260325-controle-db-bot-governance.md) — governança do bot controle-db

---

## Parcial (governança documentada, código incompleto)

- [20260325-data-stack-contract-and-templates.md](20260325-data-stack-contract-and-templates.md) — contrato de camadas SQLite/Qdrant/Supabase/Obsidian; **Supabase e Obsidian não integrados ao runtime**
- [20260325-sovereign-bibliotheca-v2-and-git-cleanup.md](20260325-sovereign-bibliotheca-v2-and-git-cleanup.md) — bibliotheca v2 e limpeza de git

---

## Proposta / Plano Ativo

- [20260326-implementacao-master-skill-global.md](20260326-implementacao-master-skill-global.md) — implementação do Orquestrador Global Master Skill (Sovereign 2026)
- [20260325-basico-bem-feito-v2-implementation.md](20260325-basico-bem-feito-v2-implementation.md) — **plano de implementação concreto**, 4 fases, critérios binários de aceite — **leitura obrigatória**
- [20260324-rich-media-expansion.md](20260324-rich-media-expansion.md) — expansão de mídias ricas

---

## Slices em Execução (ADR Semparar)

| ADR | Status | Progress | taskmaster |
|---|---|---|---|
| [20260328-adr-semparar-docs-adr-resolve.md](20260328-adr-semparar-docs-adr-resolve.md) | Em Execução | 15% | [JSON](./taskmaster/ADR-20260328-adr-semparar-docs-adr-resolve.json) |
| [20260328-implementacao-jarvis-voice-e-computer-use.md](20260328-implementacao-jarvis-voice-e-computer-use.md) | Proposta | 0% | [JSON](./taskmaster/ADR-20260328-jarvis-voice-computer-use.json) |
| [20260328-implementacao-linux-god-mode.md](20260328-implementacao-linux-god-mode.md) | Proposta | 0% | [JSON](./taskmaster/ADR-20260328-linux-god-mode.json) |

---

## Substituída

- [20260325-basico-bem-feito-swarm-memoria-dashboard.md](20260325-basico-bem-feito-swarm-memoria-dashboard.md) — duplicata; substituída por `basico-bem-feito-v2-implementation.md`

---

## Referências Obrigatórias

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)

---

## Validação

Execute para verificar conformidade Semparar:
```bash
grep -l 'AGENTS.md' docs/adr/*.md | wc -l  # deve ser 26 (todos os ADRs)
find docs/adr/taskmaster/ -name '*.json' | wc -l  # deve ser 3+
```

---

**Última atualização**: 2026-03-28
**Conformidade Semparar**: 3/3 taskmaster JSONs criados
