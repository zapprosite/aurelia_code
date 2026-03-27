# MODEL.md

Adaptador genérico para motores e copilots que precisem de um ponto de entrada curto.

## Ordem de leitura

1. [`AGENTS.md`](AGENTS.md)
2. [`docs/governance/REPOSITORY_CONTRACT.md`](docs/governance/REPOSITORY_CONTRACT.md)
3. [`docs/governance/SKILL-CATALOG.md`](docs/governance/SKILL-CATALOG.md)
4. [`.agent/rules/README.md`](.agent/rules/README.md)
5. [`docs/adr/README.md`](docs/adr/README.md)

## Contrato

- `AGENTS.md` é a autoridade do repositório.
- Skills canônicas moram em [`.agent/skills/`](.agent/skills).
- Adapters não devem duplicar regras nem manter catálogos paralelos.
- Qualquer referência a `.agents` neste repo é drift estrutural.
