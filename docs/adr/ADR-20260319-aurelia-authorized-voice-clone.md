---
description: Template ADR para slices em modo nonstop com continuidade operacional.
status: in_progress
---

# ADR-20260319-aurelia-authorized-voice-clone

## Status

- Em execução

## Slice

- slug: aurelia-authorized-voice-clone
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: docs/adr/taskmaster/ADR-20260319-aurelia-authorized-voice-clone.json

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)
- [ADR-20260319-aurelia-media-voice.md](./ADR-20260319-aurelia-media-voice.md)
- [Voz Oficial da Aurelia](../aurelia_voice_profile_20260319.md)

## Contexto

A slice `ADR-20260319-aurelia-media-voice` já fechou a infraestrutura para:

- transcript de áudio/vídeo/link
- suporte a `tts_provider=minimax`
- scripts de smoke/listagem MiniMax
- perfil vocal canônico da Aurelia

O que falta agora é a parte sensível e final:

1. receber um áudio **local e autorizado**
2. gerar ou validar um `voice_id` da Aurelia na MiniMax
3. só então trocar o Telegram TTS do fallback local para a voz oficial

O vídeo público enviado pelo usuário pode servir para análise de estilo e transcript, mas não será usado como base de clonagem.

Áudio local já entregue pelo usuário:

- `/home/will/aurelia/clone-voz/aurelia.mp3`

## Decisão

- a voz oficial da Aurelia será criada somente a partir de amostra local autorizada
- a persona vocal alvo continua:
  - feminina
  - brasileira
  - formal
  - doce e acolhedora
  - sem gírias
  - sem portunhol
- o `voice_id` canônico esperado é `aurelia-ptbr-formal-doce-v1`
- a troca do runtime só ocorrerá após smoke real e rollback definido

## Escopo

- checklist de consentimento e origem do áudio
- smoke real da MiniMax com `voice_id`
- switch controlado do Telegram TTS
- rollback imediato para `voice-proxy`

## Fora de escopo

- usar vídeo público de terceiro como base de clonagem
- mudar o runtime sem smoke real
- remover o fallback local atual

## Arquivos afetados

- `docs/adr/ADR-20260319-aurelia-authorized-voice-clone.md`
- `docs/adr/taskmaster/ADR-20260319-aurelia-authorized-voice-clone.json`
- `docs/adr/README.md`
- `docs/adr/PENDING-SLICES-20260319.md`
- opcionalmente `docs/aurelia_voice_profile_20260319.md`
- opcionalmente `~/.aurelia/config/app.json` quando houver `MINIMAX_API_KEY`

## Simulações e smoke previstos

- curl:
  - `curl -fsS http://127.0.0.1:8484/health`
- testes:
  - `go test ./pkg/tts ./internal/telegram ./internal/config -count=1`
- scripts:
  - `bash ./scripts/minimax-voice-list.sh --dry-run`
  - `bash ./scripts/minimax-tts-smoke.sh --voice-id aurelia-ptbr-formal-doce-v1 --dry-run`
  - `bash ./scripts/media-transcript.sh --input /caminho/para/amostra-autorizada.wav --dry-run`
- fallback:
  - manter `tts_provider=openai_compatible`
  - manter `voice-proxy` ativo até a validação final

## Rollout

1. receber arquivo local autorizado
2. configurar `MINIMAX_API_KEY`
3. listar vozes ou criar/ativar a voz clonada
4. validar com `minimax-tts-smoke.sh`
5. trocar config ativa do Telegram TTS
6. reiniciar o serviço e validar resposta de voz
7. manter rollback imediato pronto

## Rollback

- voltar para `tts_provider=openai_compatible`
- restaurar `tts_base_url=http://127.0.0.1:8011`
- manter a voz local atual como fallback
- preservar a voz MiniMax apenas como lane opcional até nova validação

## Evidência esperada

- prova de que a origem do áudio é local e autorizada
- arquivo local recebido: `/home/will/aurelia/clone-voz/aurelia.mp3`
- `voice_id` válido com smoke real MiniMax
- resposta de voz no Telegram com a nova voz
- `/health` saudável após a troca
- rollback funcional

## Pendências / bloqueios

- falta `MINIMAX_API_KEY`
- falta o smoke real do `voice_id` final
