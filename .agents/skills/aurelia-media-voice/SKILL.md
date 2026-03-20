---
name: aurelia-media-voice
description: Extrai transcript de video/link/audio e governa a voz oficial da Aurelia em PT-BR formal via MiniMax Audio.
---

# Aurelia Media Voice

## Objetivo

Padronizar duas tarefas recorrentes:

1. extrair texto de audio, video local ou link de video
2. preparar e operar a voz oficial da Aurelia em portugues formal, feminina, doce e profissional

## Quando usar

- quando houver um link de video e for preciso transcript rapido
- quando houver audio/video local e for preciso texto
- quando for escolher, validar ou trocar a voz oficial da Aurelia
- quando for usar MiniMax Audio para benchmark, listagem de vozes ou TTS premium

## Contrato

- `Groq` continua sendo o lane de STT
- `Kokoro (Kodoro) Local GPU` é o motor principal e padrão de TTS 2026
- `MiniMax Audio` permanece como fallback premium de nuvem
- link publico de terceiro pode ser usado para transcript e estudo de estilo
- link publico de terceiro **nao** deve ser usado como base de clonagem sem autorizacao explicita
- a clonagem oficial da Aurelia exige amostra autorizada/licenciada

## Como executar

1. Para transcript de link/audio/video:
   - rode `scripts/media-transcript.sh --input <arquivo-ou-url>`
   - use `--dry-run` primeiro se faltar `yt-dlp` ou `ffmpeg`
2. Para descobrir vozes MiniMax na conta:
   - rode `scripts/minimax-voice-list.sh --type all`
3. Para validar uma voz MiniMax:
   - rode `scripts/minimax-tts-smoke.sh --voice-id <voice_id> --output /tmp/aurelia.mp3`
4. Para a voz oficial da Aurelia:
   - siga o perfil em `docs/aurelia_voice_profile_20260319.md`
   - use referencia feminina brasileira autorizada, sem girias, com fala calma e clara
5. Ao fechar a slice:
   - atualize ADR/JSON
   - rode `./scripts/sync-ai-context.sh`

## Output esperado

- transcript em texto claro
- voz/voice_id candidata documentada
- perfil de voz alinhado ao contrato do repositorio
- evidencias de smoke com scripts e/ou testes
