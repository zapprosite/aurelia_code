# ADR-20260330-telegram-formatter-skill

## Contexto
Modelos locais (especialmente Qwen 3.5 VL) frequentemente retornam respostas estruturadas em JSON mesmo quando solicitada a saída em Markdown, poluindo a experiência do usuário no Telegram.

## Decisão
Implementamos a skill industrial `telegram-formatter-2026` e integramos o `PorteiroMiddleware.PolishOutput` como um decorador global na camada de saída do `BotController`.

## Consequências
- Respostas técnicas são automaticamente convertidas para a interface "Master Command Gateway".
- Aumento da latência de saída de ~2s devido à inferência adicional do Qwen 0.5b (Porteiro).
- Melhoria significativa de UX e profissionalismo no atendimento.
