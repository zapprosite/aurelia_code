---
description: Índice único de governança, adaptadores e ADRs do repositório.
status: active
owner: Antigravity
created: 2026-03-19
---

# Contrato do Repositório

> **Autoridade**: [AGENTS.md](../../AGENTS.md) | **Soberania 2026**

Este é o índice central de governança para humanos e agentes.

## 1. Cadeia de Autoridade

1. Humanos operadores
2. `AGENTS.md`
3. Políticas vigentes em `docs/governance/`
4. Regras em `.agent/rules/`
5. ADRs em `docs/adr/`
6. Contexto operacional em `.context/`

## 2. Ordem de Leitura Obrigatória

1. [AGENTS.md](../../AGENTS.md)
2. [README.md](../../README.md)
3. [Skill Catalog](./SKILL-CATALOG.md)
4. [CLAUDE.md](../../CLAUDE.md)
5. [GEMINI.md](../../GEMINI.md)
6. [MODEL.md](../../MODEL.md)
7. [.agent/rules/README.md](../../.agent/rules/README.md)
8. [ADR Index](../adr/README.md)
9. [.context/docs/README.md](../../.context/docs/README.md)

## 3. Fontes Permanentes de Governança

- [DATA_POLICY.md](./DATA_POLICY.md)
- [OBSIDIAN_VAULT_STANDARD.md](./OBSIDIAN_VAULT_STANDARD.md)
- [MODEL-STACK-POLICY.md](./MODEL-STACK-POLICY.md)
- [SKILL-CATALOG.md](./SKILL-CATALOG.md)
- [SECRETS.md](./SECRETS.md)
- [S-23-cloudflare-access.md](./S-23-cloudflare-access.md)

## 4. Regras Profissionais (2026)

- **ADR Obrigatória**: Necessária para cada slice estrutural.
- **Mudança Estrutural**: Proibida via chat; deve ser documentada em ADR.
- **Contratos**: Permanentes em `docs/governance/`.
- **Regras Executáveis**: Localizadas em `.agent/rules/`.
- **Sincronização**: `sync-ai-context` obrigatório em slices não triviais, handoffs e merges.
- **Paridade .env**: `.env` e `.env.example` devem ser espelhos estruturais (Zero Drift).
- **Zero Hardcode**: Placeholders `{chave-para-env}` obrigatórios para segredos.
- **Segredos Locais**: `app.json` pode persistir segredos; mascaramento obrigatório em UI, logs e docs.
- **Persistência**: Proibida a deleção do arquivo `.env` por agentes (Permissão Humana Exclusiva).

Comando canônico: `./scripts/sync-ai-context.sh`

## 5. Higiene da Raiz

**Devem ficar na raiz:**
- Contratos e adaptadores vigentes (`AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, `MODEL.md`)
- Docs de entrada (`README.md`, `CONTRIBUTING.md`, `SECURITY.md`)

**Não devem ficar na raiz:**
- Blueprints de slice, runbooks específicos e planos temporários.

**Destino correto:**
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

### 1.4 Persistent Configuration Deletion
```bash
# FORBIDDEN: Delete the environment configuration (ONLY HUMAN CAN DELETE)
rm .env
rm -f .env
rm /home/will/aurelia/.env
```

**Why:** Deleting the `.env` file wipes out the identity and infrastructure secrets of the Sovereign system.

### 1.5 Docker Catastrophe
## Modo Nonstop por Slice

Quando a slice precisa continuar entre sessões ou agentes, o padrão oficial é:

- `docs/adr/YYYYMMDD-slug.md`
- `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`

Scaffold canônico:

- `./scripts/adr-slice-init.sh <slug> --title "Title"`
