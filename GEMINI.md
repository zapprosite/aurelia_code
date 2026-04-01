---
description: Adaptador fino de orquestração para Gemini/Antigravity.
---

# GEMINI.md

Adaptador fino para Gemini no repositório `aurelia`.

## Ordem de leitura

1. [`AGENTS.md`](AGENTS.md)
2. [`.agent/rules/README.md`](.agent/rules/README.md)
3. [`docs/adr/0001-HISTORY.md`](docs/adr/0001-HISTORY.md)

## Contrato

- A autoridade máxima reside em [`AGENTS.md`](AGENTS.md).
- Planejamento, handoff e execução devem respeitar a governança local da Aurélia.
- Use exclusivamente [`.agent/skills/`](.agent/skills), [`.agent/workflows/`](.agent/workflows) e [`.agent/rules/`](.agent/rules) como fontes canônicas.
- Mantenha documentação e planos em português do Brasil, salvo exigência explícita em contrário.
