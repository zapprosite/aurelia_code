# ADR: 20260401-Estabilização de Áudio (Groq STT & Voxtral VRAM)

**Data**: 01/04/2026
**Status**: Implementado

## Contexto
A Aurélia Sovereign 2026 enfrentava falhas intermitentes de "Out of Memory" (OOM) na GPU RTX 4090. A execução simultânea do `gemma3:27b` (18GB VRAM) com o `whisper-stt` (faster-whisper GPU, ~4GB) e o `kokoro-tts` (~2GB) deixava pouca margem operacional, resultando em crashes no driver NVIDIA.

## Decisão
1. **Migração do STT para Groq**: Mover a transcrição de áudio para a API do Groq (`whisper-large-v3-turbo`). 
2. **Remoção do Kokoro (VRAM)**: Abandonar o `kokoro-tts` em favor do `voxtral-tts` (vLLM) ou Edge TTS (cloud).
3. **Limite de VRAM do Voxtral**: Reduzir `--gpu-memory-utilization` de 45% para 10% no container vLLM.
4. **Desativação do Whisper Local**: Comentar o serviço `whisper-stt` do Docker Compose.

## Consequências
- **Positivas**: Estabilidade total da GPU (~6-8GB de VRAM livre), respostas de transcrição quase instantâneas via Groq (<500ms), cérebro local (Gemma3) operando com performance máxima.
- **Negativas**: Dependência externa temporária para STT (pode ser revertida se novos modelos quantizados de Whisper surgirem).
- **Riscos**: Quota de API do Groq (mitigado por soft e hard caps no `internal/voice`).

## Alternativas Consideradas
- **Quantização extrema do Whisper**: Descartada por perda significativa de precisão em PT-BR.
- **Fallback para CPU**: Descartado por latência inaceitável (>10s para áudios curtos).
