# ADR-20260327-sovereign-streaming-pipeline

## Status
Proposto (SOTA 2026.1)

## Contexto
O sistema atual utiliza uma arquitetura baseada em `Spool` e `Polling` (1s de intervalo), resultando em uma latência perceptível (Batch Latency). Para atingir a fluidez de um assistente estilo "filme", a latência deve ser sub-vocal (início da resposta antes do fim do processamento total).

## Decisão
Implementar um pipeline de streaming reativo ponta-a-ponta:
1.  **STT Stream**: Utilizar o modo streaming do Faster-Whisper para entregar transcrições parciais.
2.  **LLM Stream**: Ativar `Stream: true` no LiteLLM/Qwen 3.5 para receber tokens em tempo real.
3.  **TTS Stream (Kodoro)**: Iniciar a síntese de áudio assim que a primeira frase completa (ou pontuação) for recebida do LLM.

## Consequências
- **Pró**: Redução da latência de resposta de ~2.5s para < 500ms.
- **Contra**: Maior complexidade na gestão de estados de interrupção (Barge-in).
- **Soberania**: 100% local, usando canais Go e WebSockets para IPC interno.
