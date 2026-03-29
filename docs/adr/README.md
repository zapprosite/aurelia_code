# ADR Index

Índice canônico de ADRs vigentes. Dividido em três categorias honestas.

## Naming
- Padrão: `YYYYMMDD-slug.md`
- ADR estrutural obrigatória para: storage, memória, runtime, governança, segurança, rede, modelos, deploy
- **Conformidade Semparar**: ADRs estruturais têm par `.md + taskmaster/*.json`

---

## Implementada (código existe e foi validado)

- [20260328-adr-semparar-docs-adr-resolve.md](ADR-20260328-adr-semparar-docs-adr-resolve.md) — **conformidade Semparar** — 26 ADRs com links, 6 JSONs taskmaster ✅
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
| [20260328-adr-semparar-docs-adr-resolve.md](ADR-20260328-adr-semparar-docs-adr-resolve.md) | ✅ Aceito | 100% | [JSON](./taskmaster/ADR-20260328-adr-semparar-docs-adr-resolve.json) |
| [20260328-claude-folder-reorganize.md](ADR-20260328-claude-folder-reorganize.md) | ✅ Aceito | 100% | [JSON](./taskmaster/ADR-20260328-claude-folder-reorganize.json) |
| [20260328-implementacao-jarvis-voice-e-computer-use.md](20260328-implementacao-jarvis-voice-e-computer-use.md) | Proposta | 0% | [JSON](./taskmaster/ADR-20260328-jarvis-voice-computer-use.json) |
| [20260328-implementacao-linux-god-mode.md](20260328-implementacao-linux-god-mode.md) | Proposta | 0% | [JSON](./taskmaster/ADR-20260328-linux-god-mode.json) |

## Jarvis/Computer Use (P1-P3) — gemma3 27b + Kokoro + LiteLLM

### P1 - Crítico (dependências base)
| ADR | Descrição | Progress | taskmaster |
|---|---|---|---|
| [20260328-smart-router-litellm-gemma3-27b.md](20260328-smart-router-litellm-gemma3-27b.md) | LiteLLM config com tiers | 0% | [JSON](./taskmaster/ADR-20260328-smart-router-litellm-gemma3-27b.json) |
| [20260328-mcp-go-client-stagehand-computer-use.md](20260328-mcp-go-client-stagehand-computer-use.md) | Go client para Stagehand | 0% | [JSON](./taskmaster/ADR-20260328-mcp-go-client-stagehand-computer-use.json) |
| [20260328-container-steel-browser-isolation.md](20260328-container-steel-browser-isolation.md) | Container Docker isolado | 0% | [JSON](./taskmaster/ADR-20260328-container-steel-browser-isolation.json) |

### P1 - SOTA Open Source (BUA, vnc-use, computer-use-mcp)
| ADR | Descrição | Progress | taskmaster |
|---|---|---|---|
| [20260328-bua-browser-use-agent-go.md](20260328-bua-browser-use-agent-go.md) | **BUA-style agent Go** (loop observe→act, 20+ tools) | 0% | [JSON](./taskmaster/ADR-20260328-bua-browser-use-agent-go.json) |
| [20260328-go-rod-browser-layer.md](20260328-go-rod-browser-layer.md) | **go-rod browser layer** (CDP, stealth mode) | 0% | [JSON](./taskmaster/ADR-20260328-go-rod-browser-layer.json) |
| [20260328-mcp-tool-schema-computer-use.md](20260328-mcp-tool-schema-computer-use.md) | **MCP tool schema** (normalized coords, 0-999) | 0% | [JSON](./taskmaster/ADR-20260328-mcp-tool-schema-computer-use.json) |
| [20260328-normalized-coordinates-hitl.md](20260328-normalized-coordinates-hitl.md) | **HitL + Safety** (dangerous patterns, confirm gate) | 0% | [JSON](./taskmaster/ADR-20260328-normalized-coordinates-hitl.json) |

### P2 - Qualidade (features principais)
| ADR | Descrição | Progress | taskmaster |
|---|---|---|---|
| [20260328-vision-pipeline-computer-use.md](20260328-vision-pipeline-computer-use.md) | Screenshot → LLM pipeline | 0% | [JSON](./taskmaster/ADR-20260328-vision-pipeline-computer-use.json) |
| [20260328-whisper-groq-gpu-budget.md](20260328-whisper-groq-gpu-budget.md) | STT optimization | 0% | [JSON](./taskmaster/ADR-20260328-whisper-groq-gpu-budget.json) |
| [20260328-computer-use-e2e-autonomous-gui.md](20260328-computer-use-e2e-autonomous-gui.md) | Agent loop autônomo | 0% | [JSON](./taskmaster/ADR-20260328-computer-use-e2e-autonomous-gui.json) |
| [20260328-anthropic-sdk-go-integration.md](20260328-anthropic-sdk-go-integration.md) | **SDK Go + LiteLLM** (betatoolrunner) | 0% | [JSON](./taskmaster/ADR-20260328-anthropic-sdk-go-integration.json) |

### P3 - Polish (E2E + UX)
| ADR | Descrição | Progress | taskmaster |
|---|---|---|---|
| [20260328-e2e-jarvis-loop-wake-tts.md](20260328-e2e-jarvis-loop-wake-tts.md) | Loop completo wake→TTS | 0% | [JSON](./taskmaster/ADR-20260328-e2e-jarvis-loop-wake-tts.json) |

### Meta-ADR de Tracking
| ADR | Descrição | taskmaster |
|---|---|---|
| [20260328-computer-use-dependency-map.md](20260328-computer-use-dependency-map.md) | Mapa de dependências completo | [JSON](./taskmaster/ADR-20260328-computer-use-dependency-map.json) |

---

## Substituída

- [20260325-basico-bem-feito-swarm-memoria-dashboard.md](20260325-basico-bem-feito-swarm-memoria-dashboard.md) — duplicata; substituída por `basico-bem-feito-v2-implementation.md`

---

## Documentação Relacionada

### PR Reviews
- [docs/pr-review/README.md](../pr-review/) — Índice de revisões de PRs
- [docs/pr-review/PR-0008-INDUSTRIALIZE-SOTA-2026.md](../pr-review/PR-0008-INDUSTRIALIZE-SOTA-2026.md) — Revisão completa do PR #8 (854 arquivos, 55k+ linhas)

### Workflows
- [.agent/workflows/super-git.md](../../.agent/workflows/super-git.md) — Combo soberano de build + delivery com gate secret scanning

---

## Referências Obrigatórias

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../governance/REPOSITORY_CONTRACT.md)

---

## Validação

Execute para verificar conformidade Semparar:
```bash
grep -l 'AGENTS.md' docs/adr/*.md | wc -l  # deve ser 33 (todos os ADRs)
find docs/adr/taskmaster/ -name '*.json' | wc -l  # deve ser 14+
```

---

**Última atualização**: 2026-03-28
**Conformidade Semparar**: ✅ 44/44 ADRs com links | ✅ 20/20 JSONs taskmaster
**ADRs Jarvis/Computer Use**: 13 novos (P1-P3 + SOTA Open Source)
**Slice completada**: 20260328-adr-semparar-docs-adr-resolve ✅
