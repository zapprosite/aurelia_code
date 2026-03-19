---
description: Índice único de governança, adapters e ADRs do repositório.
status: active
owner: codex
created: 2026-03-19
---

# Contrato do Repositório

Este é o índice único de governança para humanos e agentes.

## Ordem de leitura obrigatória

1. [AGENTS.md](../AGENTS.md)
2. [CLAUDE.md](../CLAUDE.md)
3. [CODEX.md](../CODEX.md)
4. [GEMINI.md](../GEMINI.md)
5. [MODEL.md](../MODEL.md)
6. [.agents/rules/](../.agents/rules/)
7. [ADR Index](./adr/README.md)
8. [plan.md](../plan.md)
9. [.context/docs/README.md](../.context/docs/README.md)

## O que cada arquivo manda

| Arquivo | Papel | Pode decidir? |
| --- | --- | --- |
| `AGENTS.md` | contrato soberano | sim |
| `CLAUDE.md` | adaptador do Claude | não, só executa sob contrato |
| `CODEX.md` | adaptador do Codex | não, só executa sob contrato |
| `GEMINI.md` | adaptador do Antigravity | não, só orquestra sob contrato |
| `MODEL.md` | política de modelos/voz | sim, mas subordinada a `AGENTS.md` e ADR |

## Regra profissional adotada

O repositório passa a operar com:

- ADR obrigatório por slice estrutural
- backlog oficial das pendências por slice
- adapters finos sempre linkados ao contrato central
- nenhuma mudança estrutural “só no chat”
- `sync-ai-context` obrigatório em slice não trivial, handoff e merge
- `sync-ai-context` dispensável em microedições triviais sem drift semântico

## Regra de contexto operacional

`sync-ai-context` é regra de higiene do repositório, mas com escopo profissional:

- **obrigatório**:
  - mudanças estruturais
  - slices não triviais
  - handoff entre agentes
  - preparação para merge/review final
- **dispensável**:
  - typo
  - comentário
  - teste muito pequeno
  - rename local sem impacto estrutural

Comando canônico:

- `./scripts/sync-ai-context.sh`

## Fontes operacionais

- [ADR Index](./adr/README.md)
- [ADR Slice Template](./adr/TEMPLATE-SLICE.md)
- [ADR Nonstop Template](./adr/TEMPLATE-NONSTOP-SLICE.md)
- [Pending Slices](./adr/PENDING-SLICES-20260319.md)
- [Aurelia Master Blueprint](./aurelia_master_blueprint_20260319.md)
- [JARVIS Master Plan](../plan.md)
- [.agents/skills/systems-engineer-homelab/SKILL.md](../.agents/skills/systems-engineer-homelab/SKILL.md)

## Regra de higiene da raiz

A raiz do repositório deve ficar reservada para:

- contratos soberanos e adaptadores (`AGENTS.md`, `CLAUDE.md`, `CODEX.md`, `GEMINI.md`, `MODEL.md`)
- docs de entrada do projeto (`README.md`, `CONTRIBUTING.md`, `SECURITY.md`)
- plano mestre ativo (`plan.md`)
- exemplos de bootstrap realmente globais (`mcp_servers.example.json`)

Não devem ficar na raiz:

- blueprints de slice
- runbooks de smoke
- guias de feature específicos
- `implementation_plan.md` e `task.md` de uma slice já encerrada

Destino correto:

- decisão arquitetural: `docs/adr/`
- blueprint/runbook/guia canônico: `docs/`
- artefatos de continuidade de slice: `.context/plans/<slice>/`

## Escopo em que ADR é obrigatório

- arquitetura
- providers/modelos
- storage/memória
- runtime/daemon/health
- áudio/voz
- deploy/rollout
- segurança/governança

## Modo Nonstop por Slice

Quando a slice precisa continuar entre sessões ou agentes, o padrão oficial é:

- `docs/adr/ADR-YYYYMMDD-slug.md`
- `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`

Scaffold canônico:

- `./scripts/adr-slice-init.sh <slug> --title "Title"`

## Escopo em que ADR pode ser dispensado

- typo
- rename local sem impacto externo
- teste pontual
- comentário/limpeza sem mudança de comportamento
