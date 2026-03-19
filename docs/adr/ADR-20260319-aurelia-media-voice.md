---
description: Template ADR para slices em modo nonstop com continuidade operacional.
status: in_progress
---

# ADR-20260319-aurelia-media-voice

## Status

- Em execução

## Slice

- slug: aurelia-media-voice
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: docs/adr/taskmaster/ADR-20260319-aurelia-media-voice.json

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)

## Contexto

O repositório já possui:

- `Groq` como lane de STT
- TTS local no Telegram via `voice-proxy`
- slices de voice capture, spool e deploy já executadas

Faltava fechar um pacote profissional para:

1. extrair transcript de vídeo/link/áudio de forma repetível
2. governar a voz oficial da Aurelia
3. preparar a integração de `MiniMax Audio` como lane premium de TTS

Também surgiu a necessidade de evitar deriva insegura: vídeo público de terceiro pode inspirar o estilo e gerar transcript, mas não deve virar base de clonagem de voz sem autorização.

## Decisão

- criar uma skill dedicada `aurelia-media-voice`
- adotar `MiniMax Audio` como lane premium de TTS para a voz oficial da Aurelia
- manter o TTS local atual como fallback operacional
- tratar links públicos apenas como fonte de transcript/estudo, não de clonagem
- formalizar o perfil vocal canônico em `docs/aurelia_voice_profile_20260319.md`

## Escopo

- skill e referência operacional
- scripts para transcript e smoke/listagem MiniMax
- suporte a `tts_provider=minimax` no runtime
- documentação e continuidade via ADR/JSON

## Fora de escopo

- ativação live da lane MiniMax sem `MINIMAX_API_KEY`
- clonagem de voz de terceiro sem autorização
- TTS em background fora do caminho atual do Telegram

## Arquivos afetados

- `pkg/tts/minimax.go`
- `pkg/tts/minimax_test.go`
- `internal/telegram/bot.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `scripts/media-transcript.sh`
- `scripts/minimax-voice-list.sh`
- `scripts/minimax-tts-smoke.sh`
- `.agents/skills/aurelia-media-voice/`
- `docs/aurelia_voice_profile_20260319.md`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/health`
- testes:
  - `go test ./pkg/tts ./internal/telegram ./internal/config -count=1`
- scripts:
  - `bash ./scripts/media-transcript.sh --input https://example.com/video --dry-run`
  - `bash ./scripts/minimax-voice-list.sh --dry-run`
  - `bash ./scripts/minimax-tts-smoke.sh --voice-id aurelia-ptbr-formal-doce-v1 --dry-run`
- fallback:
  - manter `tts_provider=openai_compatible` até existir chave e `voice_id` reais

## Rollout

1. entregar provider `minimax` no código
2. manter fallback local como caminho ativo
3. abrir slice filha para a clonagem autorizada da voz da Aurelia
4. depois da chave MiniMax e da amostra local autorizada, validar `voice_id` com smoke
5. só então considerar switch no runtime

## Rollback

- voltar para `tts_provider=openai_compatible`
- manter `voice-proxy` como lane ativo
- não remover scripts/skill; apenas marcar MiniMax como inativo

## Evidência esperada

- testes de `pkg/tts` verdes
- builder do Telegram aceitando `tts_provider=minimax`
- scripts com `--dry-run` coerentes
- skill documentada e linkada ao contrato
- slice filha aberta para a execução autorizada da voz oficial

## Pendências / bloqueios

- falta `MINIMAX_API_KEY`
- falta amostra autorizada para a voz oficial da Aurelia
- `yt-dlp` e `ffmpeg` ainda não estão instalados nesta máquina
- execução final segue em `ADR-20260319-aurelia-authorized-voice-clone`
