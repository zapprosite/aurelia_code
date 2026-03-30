---
type: doc
name: architecture
description: System architecture, layers, patterns, and design decisions
category: architecture
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Arquitetura — Aurélia Sovereign 2026.2

## Visão Geral da Arquitetura

Aurélia é um bot Telegram Go com pipeline de áudio/texto nativo, roteamento inteligente de LLMs via LiteLLM, e sovereignidade máxima (50-80% local). O sistema é monolítico Go com fronteira clara entre entrada (Telegram), processamento (skill routing), execução (agent loop), e saída (TTS + Telegram).

## Camadas Arquiteturais

- **Telegram Bot Handler** (`internal/telegram/`) — Interface com usuário via Telegram (voz, texto, callbacks)
- **Skill Router** (`internal/skill/`) — Classifica intent e seleciona skill ativa
- **Agent Loop** (`internal/agent/`) — Loop de execução de ferramentas com tool calling
- **LLM Providers** (`pkg/llm/`) — thin wrappers sobre OpenAI-compatible APIs
- **Gateway** (`internal/gateway/`) — Roteamento LiteLLM com cascade determinística
- **Tools** (`internal/tools/`) — Ferramentas de execução (run_command, read_file, web_search, scheduling)
- **STT/TTS** (`pkg/stt/`, `pkg/tts/`) — Pipeline de mídia (faster-whisper + Kokoro)
- **Segurança** (`internal/sentinel/`, `pkg/porteiro/`) — Guardrails de entrada e portaria de execução

## Padrões Detectados

| Padrão | Confiança | Localização | Descrição |
|--------|-----------|-------------|-----------|
| Factory | 90% | `pkg/stt/factory.go`, `pkg/tts/factory.go` | Cria transcriber/synthesizer baseado em config |
| Circuit Breaker | 80% | `internal/gateway/provider.go` | Retry com cooldown entre tiers |
| Pipeline | 95% | `internal/telegram/input_pipeline.go` | Sequência entrada → guard → skill → LLM → output |
| Singleton de Config | 85% | `internal/config/config.go` | AppConfig carregado uma vez no startup |
| Decorator | 80% | `pkg/tts/segmented.go` | Wraps synthesizer com chunking para textos longos |

## Entry Points

- `cmd/aurelia/main.go` — entry point do binário
- `cmd/aurelia/app.go` — inicialização do app (buildLLMProvider, wiring)
- `internal/telegram/input.go` — handleVoice, handleText (handlers Telegram)
- `/v1/telegram/impersonate` — rota HTTP para testes de stress

## API Pública (pacotes exportados)

| Símbolo | Tipo | Localização |
|---------|------|-------------|
| `BotController` | struct | `internal/telegram/` |
| `ProcessExternalInput` | func | `internal/telegram/input_pipeline.go` |
| `NewTranscriber` | func | `pkg/stt/factory.go` |
| `NewSynthesizer` | func | `pkg/tts/factory.go` |
| `NewProvider` (gateway) | func | `internal/gateway/provider.go` |
| `WebSearchHandler` | func | `internal/tools/web_search.go` |
| `RunCommandHandler` | func | `internal/tools/run_command.go` |

## Limites de Sistema

- **VRAM RTX 4090**: 24.5GB total — ~21GB consumidos por whisper + qwen3.5 + Kokoro
- **RAM**: 30GB — ~9GB disponível para burst
- **Concorrência Ollama**: 2 requests simultâneos (rate limit aplicado via systemd)
- **Timeout LLM**: 90s global, 45s para local
- **Timeout STT**: 120s (whisper), 60s (Groq)
- **Cache Redis**: 30 dias TTL para porteiro

## ADR Index

Ver `docs/adr/README.md` — ADRs em `docs/adr/slices/` cobrem todas as decisões arquiteturais desde 2026-03-28.
