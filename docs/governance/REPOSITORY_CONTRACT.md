---
description: Índice único de governança, adaptadores e ADRs do repositório.
status: active
owner: Antigravity
created: 2026-03-19
---

# Contrato do Repositório

Este é o índice central de governança para humanos e agentes.

## Cadeia de autoridade

1. Humanos operadores
2. `AGENTS.md`
3. Políticas vigentes em `docs/governance/`
4. Regras em `.agent/rules/`
5. ADRs em `docs/adr/`
6. Contexto operacional em `.context/`

## Ordem de leitura obrigatória

1. [AGENTS.md](../../AGENTS.md)
2. [README.md](../../README.md)
3. [CLAUDE.md](../../CLAUDE.md)
4. [GEMINI.md](../../GEMINI.md)
5. [.agent/rules/README.md](../../.agent/rules/README.md)
6. [ADR Index](../adr/README.md)
7. [.context/docs/README.md](../../.context/docs/README.md)

## Fontes permanentes de governança

- [DATA_GOVERNANCE.md](./DATA_GOVERNANCE.md)
- [DATA_STACK_STANDARD.md](./DATA_STACK_STANDARD.md)
- [SCHEMA_REGISTRY.md](./SCHEMA_REGISTRY.md)
- [OBSIDIAN_VAULT_STANDARD.md](./OBSIDIAN_VAULT_STANDARD.md)
- [MODEL-STACK-POLICY.md](./MODEL-STACK-POLICY.md)
- [SECRETS.md](./SECRETS.md)
- [S-23-cloudflare-access.md](./S-23-cloudflare-access.md)

## Regra profissional adotada

- ADR obrigatória por slice estrutural
- nenhuma mudança estrutural só no chat
- contratos permanentes em `docs/governance/`
- regras executáveis em `.agent/rules/`
- `sync-ai-context` obrigatório em slice não trivial, handoff e merge
- `sync-ai-context` dispensável em microedições sem drift semântico

Comando canônico:

- `./scripts/sync-ai-context.sh`

## Higiene da raiz

Devem ficar na raiz:

- contratos e adaptadores vigentes (`AGENTS.md`, `CLAUDE.md`, `GEMINI.md`)
- docs de entrada (`README.md`, `CONTRIBUTING.md`, `SECURITY.md`)

Não devem ficar na raiz:

- blueprints de slice
- runbooks específicos
- planos temporários de execução

Destino correto:

- ADR: `docs/adr/`
- policy e contrato durável: `docs/governance/`
- artefatos de continuidade: `.context/plans/`

## Escopo em que ADR é obrigatória

- arquitetura
- providers e modelos
- storage e memória
- runtime, daemon e health
- áudio e voz
- deploy e rollout
- segurança e governança

## Modo Nonstop por Slice

Quando a slice precisa continuar entre sessões ou agentes, o padrão oficial é:

- `docs/adr/YYYYMMDD-slug.md`
- `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`

Scaffold canônico:

- `./scripts/adr-slice-init.sh <slug> --title "Title"`
