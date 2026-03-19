---
description: Índice de Registros de Decisão Arquitetural (ADR).
---

# 🏛️ Architectural Decision Records (ADR)

Este diretório contém o registro histórico de todas as decisões técnicas significativas que moldaram este repositório.

## Contrato

- A autoridade primária continua em [AGENTS.md](../../AGENTS.md)
- O índice de governança está em [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- O template oficial por slice está em [TEMPLATE-SLICE.md](./TEMPLATE-SLICE.md)
- O template nonstop por slice está em [TEMPLATE-NONSTOP-SLICE.md](./TEMPLATE-NONSTOP-SLICE.md)
- O backlog oficial de pendências está em [PENDING-SLICES-20260319.md](./PENDING-SLICES-20260319.md)

## Índice de Decisões

| ID | Título | Data | Status |
| :--- | :--- | :--- | :--- |
| [20260317](./20260317-rebalanceamento-elite.md) | Rebalanceamento de Contexto Elite | 2026-03-17 | ✅ Aceito |
| [20260318-rim](./20260318-estrategia-rim.md) | Estratégia RIM | 2026-03-18 | ✅ Aceito |
| [20260318-mcp-antigravity](./20260318-integracao-mcp-antigravity.md) | Integração MCP Antigravity | 2026-03-18 | ✅ Aceito |
| [20260318-ai-context](./20260318-implementando-ai-context.md) | Implementando ai-context | 2026-03-18 | ✅ Aceito |
| [20260318-mirror-template](./20260318-estrategia-mirror-template.md) | Estratégia Mirror Template | 2026-03-18 | ✅ Aceito |
| [20260319-sync-ai-context](./20260319-sync-ai-context-como-regra-de-slice.md) | `sync-ai-context` como regra de slice | 2026-03-19 | ✅ Aceito |
| [20260319-voice-capture](./20260319-voice-capture-plane.md) | Voice capture plane real | 2026-03-19 | 🟡 Proposto |

## Como Criar um ADR

Use o template oficial:

- [TEMPLATE-SLICE.md](./TEMPLATE-SLICE.md)
- [TEMPLATE-NONSTOP-SLICE.md](./TEMPLATE-NONSTOP-SLICE.md)

Para slices de continuidade longa, use também o JSON companheiro em:

- `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`

Scaffold:

- `./scripts/adr-slice-init.sh <slug> --title "Title"`

Roteiro mínimo:

**Contexto ➜ Decisão ➜ Arquivos afetados ➜ Testes ➜ Rollout ➜ Rollback ➜ Consequências**

ADRs ajudam novos agentes e humanos a entender o "Porquê" por trás da estrutura atual.

## Quando ADR é obrigatória

- arquitetura
- slice estrutural
- provider/modelo
- storage
- runtime
- áudio/voz
- deploy
- segurança/governança
