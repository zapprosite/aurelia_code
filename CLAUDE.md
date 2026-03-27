# CLAUDE.md

Adaptador fino para Claude no repositório `aurelia`.

## Ordem de leitura

1. [`AGENTS.md`](AGENTS.md)
2. [`docs/governance/REPOSITORY_CONTRACT.md`](docs/governance/REPOSITORY_CONTRACT.md)
3. [`docs/governance/SKILL-CATALOG.md`](docs/governance/SKILL-CATALOG.md)
4. [`.agent/rules/README.md`](.agent/rules/README.md)
5. [`docs/adr/README.md`](docs/adr/README.md)

## Contrato

- Este arquivo não é autoridade. A autoridade está em [`AGENTS.md`](AGENTS.md).
- Use [`.agent/skills/`](.agent/skills), [`.agent/workflows/`](.agent/workflows) e [`.agent/rules/`](.agent/rules) como caminhos canônicos.
- Trate `.agents` como legado e corrija referências quando encontradas.
- Skills, roteamento semântico e auditoria do catálogo são definidos em [`docs/governance/SKILL-CATALOG.md`](docs/governance/SKILL-CATALOG.md).
