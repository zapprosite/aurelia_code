# 20260328 - Arquitetura Multimodal Profissional e Otimização de GPU

## Status: Proposto

## Contexto
O sistema atual sofre de contenção de recursos na RTX 4090 (24GB). A carga combinada de **Gemma3-27B** (Ollama), **Whisper-Large-v3** (STT) e **Kodoro** (TTS) atinge o limite de VRAM, resultando em `CUDA Out of Memory` e falhas na transcrição ("Erro de Transcrição").

## Decisões
1. **Dimensionamento Sensato de STT**: Substituir o modelo `large-v3` pelo `medium` (multilingual) para transcrição local. O ganho de performance e economia de ~2GB de VRAM compensam a perda marginal de precisão para áudios de bot.
2. **Estratégia Hybrid-Cloud (Soberania Pragmática)**: Habilitar o uso da Groq API para STT como provedor prioritário (ou fallback instantâneo). Isso reduz o stress na GPU local para tarefas críticas de inferência (LLM).
3. **UX Silencioso**: Remodelar o tratamento de erros em `internal/telegram/messages.go` para que o sistema tente automaticamente um segundo provedor antes de reportar erro ao usuário.
4. **Governança de Wake-Word**: Manter o `openwakeword` em CPU (via factory default) para não onerar desnecessariamente os núcleos CUDA.

## Consequências
- **Estabilidade Multimodal**: Fim das quedas do serviço de transcrição por falta de memória.
- **Latência Reduzida**: Groq entrega transcrições em milissegundos, liberando a GPU local para o Kodoro/Gemma3.
- **Soberania Preservada**: O processamento local continua sendo o primário (ou fallback robusto), mantendo a funcionalidade offline.

## Referências
- `app.json` configuration.
- `internal/telegram/messages.go` error handling.
- `whisper-local` docker-compose configuration.
