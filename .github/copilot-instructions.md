# Copilot Instructions

Use este arquivo apenas como adapter fino.

## Mandatory Read Order

1. [`../AGENTS.md`](../AGENTS.md)
2. [`../docs/governance/REPOSITORY_CONTRACT.md`](../docs/governance/REPOSITORY_CONTRACT.md)
3. [`../docs/governance/SKILL-CATALOG.md`](../docs/governance/SKILL-CATALOG.md)
4. [`../.agent/rules/README.md`](../.agent/rules/README.md)
5. [`../docs/adr/README.md`](../docs/adr/README.md)

## Contract

- `AGENTS.md` is the repository authority.
- The canonical skill catalog lives in [`../.agent/skills/`](../.agent/skills).
- Do not invent a parallel governance layer in Copilot-specific files.
- Treat `.agents` references as legacy drift and prefer `.agent`.
