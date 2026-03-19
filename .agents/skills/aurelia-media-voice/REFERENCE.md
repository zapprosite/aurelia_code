# Referência — Aurelia Media Voice

## Persona vocal alvo

- feminina
- brasileira
- tom doce, calmo e acolhedor
- dicção clara e elegante
- sem gírias
- sem regionalismo informal
- sem portunhol
- ritmo pausado e equilibrado
- sonoridade polida para atendimento corporativo

## Pré-requisitos

- `curl`
- `jq`
- `ffmpeg`
- `ffprobe`
- `yt-dlp` para links
- `GROQ_API_KEY` para transcript
- `MINIMAX_API_KEY` para voz MiniMax

## Regras de segurança

- nao clonar voz de terceiro sem autorizacao
- usar links publicos de terceiros apenas para transcript, resumo e estudo de estilo
- a voz oficial da Aurelia deve nascer de referencia autorizada

## Convenção sugerida de voz

- `voice_id`: `aurelia-ptbr-formal-doce-v1`
- `modelo`: `speech-2.8-hd`
- `language_boost`: `Portuguese`
- `formato`: `mp3`

## Referências operacionais

- `scripts/media-transcript.sh`
- `scripts/minimax-voice-list.sh`
- `scripts/minimax-tts-smoke.sh`
- `docs/aurelia_voice_profile_20260319.md`
- `docs/adr/ADR-20260319-aurelia-media-voice.md`
