# Aurélia — Sovereign 2026.2

Bot Telegram sênior do Will. Engenharia de software + operações + pesquisa, tudo via chat.

## Guias Essenciais

- [Visão Geral do Projeto](./project-overview.md)
- [Arquitetura](./architecture.md)
- [Fluxo de Dados](./data-flow.md)
- [Fluxo de Desenvolvimento](./development-workflow.md)
- [Estratégia de Testes](./testing-strategy.md)
- [Glossário](./glossary.md)
- [Segurança](./security.md)
- [Ferramentas e Produtividade](./tooling.md)

## ADR Slices (SOTA 2026)

| Slice | Arquivo | Status |
|-------|---------|--------|
| STT local-first | `slices/20260330-stt-local-first-cascade.md` | ✅ Implementado |
| Kokoro TTS GPU | `slices/20260330-kokoro-tts-gpu-local.md` | ✅ Implementado |
| Rate limiting | `slices/20260330-rate-limiting-smart-scheduler.md` | ✅ Implementado |
| Onboarding sênior | `slices/20260330-jarvis-onboarding-senior.md` | ✅ Implementado |
| Porteiro bypass owner | `slices/20260330-porteiro-owner-bypass.md` | ✅ Implementado |
| Voice pipeline PT-BR | `slices/20260330-voice-pipeline-ptbr-local.md` | ✅ Implementado |
| LiteLLM cascade | `slices/20260330-litellm-cascade-qwen36.md` | ✅ Implementado |

## Stack em Produção

```
LiteLLM proxy (:4000)    ✅  qwen3.5:9b local → qwen3.6-free → groq → paid
Whisper local (:8020)    ✅  faster-whisper-server → Groq fallback
Kokoro TTS (:8012)      ✅  ghcr.io/remsky/kokoro-fastapi-gpu
Redis (:6379)             ✅  Cache + Porteiro
Qdrant (:6333)           ✅  Memória vetorial
```
