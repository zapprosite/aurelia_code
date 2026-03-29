# ADR Slice: Barge-in VAD Interrupt

## Contexto
A fluidez "Movie-Like" exige que o Jarvis pare de falar imediatamente ao ser interrompido pelo usuário.

## Decisão
Implementar um monitor de VAD (Voice Activity Detection) como um Ator específico no pipeline reativo SAP.

## Consequências
- Interrupção instantânea de áudio e cancelamento de tokens pendentes.
- Aumento marginal de processamento de microfone (Whisper-Tiny).
