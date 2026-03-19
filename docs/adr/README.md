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
| [20260319-root-docs](./20260319-root-document-hygiene.md) | Higiene documental da raiz | 2026-03-19 | ✅ Aceito |
| [20260319-antigravity-light](./20260319-antigravity-copiloto-leve.md) | Antigravity como copiloto leve | 2026-03-19 | ✅ Aceito |
| [20260319-groq-stt](./20260319-groq-stt-ptbr-runtime.md) | Groq STT PT-BR no runtime | 2026-03-19 | ✅ Aceito |
| [20260319-homelab-tutor](./20260319-homelab-tutor-v2.md) | Homelab Tutor v2 | 2026-03-19 | ✅ Aceito |
| [20260319-keepassxc](./20260319-keepassxc-cofre-humano.md) | KeePassXC como cofre humano | 2026-03-19 | ✅ Aceito |
| [20260319-telegram-tts](./20260319-telegram-tts-local.md) | Telegram TTS local via voice-proxy | 2026-03-19 | ✅ Aceito |
| [20260319-aurelia-media-voice](./ADR-20260319-aurelia-media-voice.md) | Transcript de mídia e voz oficial da Aurelia | 2026-03-19 | 🔵 Em execução |
| [20260319-aurelia-authorized-voice-clone](./ADR-20260319-aurelia-authorized-voice-clone.md) | Clonagem autorizada da voz oficial da Aurelia | 2026-03-19 | 🔵 Em execução |
| [20260319-voice-capture](./20260319-voice-capture-plane.md) | Voice capture plane real | 2026-03-19 | 🟡 Proposto |
| [20260319-voice-capture-runtime](./ADR-20260319-voice-capture-runtime.md) | Voice capture runtime nonstop | 2026-03-19 | 🔵 Em execução |
| [20260319-state-memory-runtime](./ADR-20260319-state-memory-runtime.md) | State memory runtime nonstop | 2026-03-19 | 🔵 Em execução |
| [20260319-deploy-gateway-voice](./ADR-20260319-deploy-gateway-voice.md) | Deploy gateway voice nonstop | 2026-03-19 | 🔵 Em execução |
| [20260319-extensions-governance](./ADR-20260319-extensions-governance.md) | Governança de extensões | 2026-03-19 | ✅ Aceito |

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
