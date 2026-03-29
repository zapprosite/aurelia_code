# ADR Slice: Actor Self-Healing (Supervision)

## Contexto
O pipeline SAP depende de atores independentes. Uma falha em um ator (ex: mpv buffer overflow) não deve travar o sistema.

## Decisão
Implementar um Supervisor pattern em Go que monitore canais liveness e realize o reinício automático dos atores sem interromper o fluxo do orquestrador.

## Consequências
- Uptime industrial 24/7.
- Maior complexidade na gestão de estado interno dos atores durante o restart.
