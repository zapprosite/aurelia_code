# ADR 20260319-groq-stt-ptbr-runtime

**Status**: Aceito  
**Data**: 2026-03-19

## Contexto

O homelab precisava de STT em PT-BR sem sacrificar VRAM local do runtime principal. O caminho local com Whisper concorria com o orçamento de GPU do LLM residente.

## Decisão

Adotar `Groq whisper-large-v3-turbo` como lane principal de STT, com:

- `language=pt`
- `temperature=0`
- governança própria de budget/rate limit
- fallback local de STT quando necessário

O LLM local continua responsável por instrução, tool use e decisão. `Qdrant` e `Supabase` permanecem no plano de memória/estado; `Groq` não vira fonte de verdade nem plano de controle.

## Consequências

Positivas:

- preserva VRAM do host para o modelo residente
- melhora latência de transcrição em PT-BR
- separa claramente o lane de áudio do lane de raciocínio

Trade-offs:

- introduz dependência externa para STT
- exige governor diário e fallback explícito

## Referências

- [groq_ptbr_audio_blueprint_20260319.md](../groq_ptbr_audio_blueprint_20260319.md)
- [groq_stt_simulation_20260319.md](../groq_stt_simulation_20260319.md)
- [jarvis_local_voice_blueprint_20260319.md](../jarvis_local_voice_blueprint_20260319.md)
