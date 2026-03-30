---
type: doc
name: project-overview
description: High-level overview of the project, its purpose, and key components
category: overview
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Visão Geral — Aurélia Sovereign 2026.2

## O que é

Assistente de engenharia sênior do Will via Telegram. Responde comandos de software, infraestrutura, pesquisa e operações de Homelab. Executa localmente (RTX 4090 + AMD 7900X) com fallback inteligente para cloud.

## Stack de Produção

- **RTX 4090 24GB VRAM** — qwen3.5:9b local + faster-whisper + Kokoro
- **AMD Ryzen 9 7900X 12c/24t** — 30GB RAM
- **NVMe Gen5 4TB** — armazenamento
- **Ubuntu Desktop** — controle total via DISPLAY=:1

## Stack de Software

- **Go** — Telegram bot, pipeline, agents, tools
- **Docker** — LiteLLM proxy, Redis, Qdrant, Whisper, Kokoro
- **LiteLLM** — proxy local com cascade determinística
- **Ollama** — inference engine local
- **Redis** — cache + porteiro
- **Qdrant** — vetor DB para memória
- **faster-whisper-server** — STT local (CUDA)
- **Kokoro TTS** — síntese vocal PT-BR local

## Repositório

- `/home/will/aurelia/` — código-fonte Go + configs
- `docs/adr/slices/` — ADRs de cada decisão desde 2026-03-29
- `pkg/` — libs reutilizáveis (llm, stt, tts, porteiro)
- `internal/` — telegram, gateway, skill, tools, agent, config
- `cmd/aurelia/` — entry point
- `scripts/audit/audit-secrets.sh` — auditoria soberana de segredos

## Rotas HTTP de Debug

- `GET /health` — status do bot
- `POST /v1/telegram/impersonate` — stress test (mock Telegram user)
- `GET localhost:4000/health` — LiteLLM proxy
- `GET localhost:6333/healthz` — Qdrant
- `GET localhost:8020/health` — Whisper local
- `GET localhost:8012/health` — Kokoro TTS

## Checklist de Início Rápido

1. `docker compose -f /home/will/aurelia/docker-compose.yml up -d`
2. `sudo systemctl restart aurelia`
3. `curl http://localhost:8585/health`
4. Teste: envie áudio pelo Telegram
