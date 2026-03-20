---
title: Voice Plane Unificado (Arquitetura, Captura e Identidade)
status: accepted
date: 2026-03-19
owner: codex
---

# ADR-20260319-voice-plane-unified

## Contexto
O projeto Aurelia evoluiu de um agente puramente textual para um sistema JARVIS local-first. A fragmentação de decisões sobre áudio (STT, TTS, captura contínua e clonagem) exigia uma visão unificada para garantir consistência arquitetural e ética.

## Decisão
Implementar um **Voice Plane** soberano e local-first, composto por quatro camadas principais:

### 1. Camada de Escuta (Capture Plane)
- **Detecção**: `openWakeWord` local para frase de ativação.
- **Qualidade**: `Silero VAD` para supressão de silêncio/ruído.
- **Contexto**: `ring buffer` local para preservar o áudio imediatamente anterior e posterior ao gatilho.
- **Runtime**: Worker dedicado de captura integrado ao spool de jobs local.

### 2. Camada de Transcrição (STT)
- **Primary**: `Groq whisper-large-v3-turbo` (lane externo de alta performance).
- **Motivo**: Preservar a VRAM do host para o LLM residente.
- **Guardrails**: Budget diário rigoroso e fallback local disponível.
- **Identidade Contextual**: Hard override para PT-BR (`language=pt`) e temperatura zero.

### 3. Camada de Síntese (TTS e Identidade)
- **Motor Primário**: `chatterbox-tts` (Kokoro) rodando localmente em Docker.
- **Identidade Canônica**: Voz "Aurelia" (Doce, educada, PT-BR formal) definida em `docs/aurelia_voice_profile_20260319.md`.
- **Clonagem**: Permitida apenas mediante amostra autorizada (`Aurelia.wav`) e consentimento registrado, com trilha de auditoria.
- **Fallback**: Degradação automática para texto se o pipeline de áudio falhar.

### 4. Governança e Operação
- **Spool & Spooler**: Gerenciamento de fila com heartbeat e estados persistidos em SQLite.
- **Observabilidade**: Métricas exportadas via `/metrics` e status detalhado em `/v1/voice/status`.

## Consequências
- **Positivas**: Redução drástica da latência de resposta, identidade vocal protegida e consistente, e uso eficiente de hardware.
- **Riscos**: Dependência de tokens de API externos para STT e sensibilidade ambiental do microfone (ruído).

## Referências
- `docs/blueprints/jarvis_local_voice_blueprint_20260319.md`
- `docs/aurelia_voice_profile_20260319.md`
- `internal/voice/`, `internal/tts/`, `pkg/stt/`
