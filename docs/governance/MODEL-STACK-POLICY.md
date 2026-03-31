# MODEL-STACK-POLICY.md

> **Versão**: Industrial 2026

Este documento define a hierarquia e o uso de modelos (LLM, TTS, STT) no ecossistema Aurélia.

## Tier 0 (Local - Homelab)
- **Ollama**: Qwen 2.5 (0.5b / 7b / 72b), Llama 3.3.
- **Audio (STT)**: Faster-Whisper-v3 (Large).
- **Audio (TTS)**: Kokoro (Kodoro) 2026.

## Tier 1 (Cloud - Antigravity)
- **Google**: Gemini 1.5 Pro (Orquestrador) / Flash (Processamento Rápido).
- **Anthropic**: Claude 3.5 Sonnet / 3 Opus.

## Tier 2 (Public / Fallback)
- **OpenRouter**: Llama, DeepSeek, GLM-5.

---
*Assinado: Aurélia (Soberano 2026)*
*Atualizado: 2026-03-31*
