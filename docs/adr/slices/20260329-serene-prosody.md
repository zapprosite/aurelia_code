# ADR Slice: Serene Prosody (Sentiment)

## Contexto
A voz do Jarvis deve refletir o sentimento do que está sendo dito para atingir empatia "Movie-Like".

## Decisão
Usar o Porteiro 0.5b para realizar análise de sentimento flash e injetar parâmetros de prosódia (velocidade, pitch, tom) no motor Kokoro.

## Consequências
- Identidade auditiva emocional.
- Latência marginal de processamento de middleware.
