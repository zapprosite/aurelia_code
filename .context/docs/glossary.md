---
type: doc
name: glossary
description: Project terminology, type definitions, domain entities, and business rules
category: glossary
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Glossário — Aurélia Sovereign 2026.2

## Termos Principais

| Termo | Significado |
|-------|-------------|
| **SOTA 2026** | Estado da arte 2026 — práticas, padrões, modelos vigentes |
| **Porteiro** | Input guardrail (qwen2.5:0.5b) que filtra injeção de prompt antes do LLM principal |
| **Sentinel** | Nome do pacote de segurança (`pkg/porteiro/`) |
| **Cascade** | Rota determinística de fallback entre provedores LLM |
| **Jarvis Mode** | Modo de resposta conversacional (suffix de prompt + TTS em paralelo) |
| **Bootstrap** | Máquina de estados do onboarding (welcome → perfil → nome → pronto) |
| **Smart Router** | Nome do serviço LiteLLM no Docker Compose |
| **Onboarding** | Fluxo de primeira mensagem do bot (curto, direto, sem enrolação) |

## Abreviações

| Sigla | Expansão |
|-------|----------|
| STT | Speech-to-Text (áudio → texto) |
| TTS | Text-to-Speech (texto → áudio) |
| VRAM | Video RAM (memória da GPU) |
| KV Cache | Key-Value cache da GPU durante inference |
| LiteLLM | Proxy de roteamento LLM (github.com/berriai/litellm) |
| LLMOps | Operações de LLM (monitoramento, logs, deployment) |
| SOT A | Sovereign Operational Tier Architecture |

## Interfaces Exportadas

| Interface | Arquivo | Propósito |
|-----------|---------|-----------|
| `Transcriber` | `pkg/stt/stt.go` | Implementada por GroqTranscriber, LocalTranscriber |
| `Synthesizer` | `pkg/tts/tts.go` | Implementada por OpenAICompatibleSynthesizer |
| `LLMProvider` | `pkg/llm/provider.go` | Interface para todos os providers LLM |
| `ownedLLMProvider` | `internal/telegram/llm_override.go` | Provider com Close() para multi-bot pool |

## Entidades de Configuração

| Entidade | Onde | Uso |
|----------|-------|-----|
| `AppConfig` | `internal/config/config.go` | Config global do app |
| `BotConfig` | `internal/config/config.go` | Per-bot Telegram config |
| `TutorProcessor` | `internal/voice/tutor.go` | 24/7 Jarvis listener mode |
| `VoiceConfig` | `internal/voice/` | Config de captura de voz |
