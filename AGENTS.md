# Project Rules and Guidelines

> Entry point de documentação e governança do repositório.

## Leitura obrigatória

1. [README.md](./README.md)
2. [docs/governance/REPOSITORY_CONTRACT.md](./docs/governance/REPOSITORY_CONTRACT.md)
3. [CLAUDE.md](./CLAUDE.md)
4. [GEMINI.md](./GEMINI.md)
5. [.agent/rules/README.md](./.agent/rules/README.md)
6. [docs/adr/README.md](./docs/adr/README.md)
7. [.context/docs/README.md](./.context/docs/README.md)

## Core Guides

- [Project Overview](./.context/docs/project-overview.md)
- [Architecture Notes](./.context/docs/architecture.md)
- [Development Workflow](./.context/docs/development-workflow.md)
- [Testing Strategy](./.context/docs/testing-strategy.md)
- [Glossary & Domain Concepts](./.context/docs/glossary.md)
- [Data Flow & Integrations](./.context/docs/data-flow.md)
- [Security & Compliance Notes](./.context/docs/security.md)
- [Tooling & Productivity Guide](./.context/docs/tooling.md)

## Governance Highlights

- contrato central: [docs/governance/REPOSITORY_CONTRACT.md](./docs/governance/REPOSITORY_CONTRACT.md)
- stack de modelos: [docs/governance/MODEL-STACK-POLICY.md](./docs/governance/MODEL-STACK-POLICY.md) e [.agent/rules/13-model-stack-policy.md](./.agent/rules/13-model-stack-policy.md)
- stack de dados: [docs/governance/DATA_POLICY.md](./docs/governance/DATA_POLICY.md), [docs/governance/OBSIDIAN_VAULT_STANDARD.md](./docs/governance/OBSIDIAN_VAULT_STANDARD.md)
- rede e subdomínios: [docs/governance/S-23-cloudflare-access.md](./docs/governance/S-23-cloudflare-access.md) e [.agent/rules/12-network-governance.md](./.agent/rules/12-network-governance.md)
- ADR por slice estrutural: [docs/adr/README.md](./docs/adr/README.md)
- continuidade de slice: [.agent/workflows/adr-semparar.md](./.agent/workflows/adr-semparar.md) e [scripts/adr-slice-init.sh](./scripts/adr-slice-init.sh)
- bibliotheca e integração Go/Node/Bash: [.agent/rules/15-sovereign-bibliotheca.md](./.agent/rules/15-sovereign-bibliotheca.md)
- voz e mídia: [docs/jarvis_local_voice_blueprint_20260319.md](./docs/jarvis_local_voice_blueprint_20260319.md) e [.agent/skills/aurelia-media-voice/SKILL.md](./.agent/skills/aurelia-media-voice/SKILL.md)

## Related Resources

- [Agent Handbook](./.context/agents/README.md)
- [Workflow changelog](./.context/workflow/docs/changelog-post-reboot-validation-2026-03-19.md)
- [Codebase map](./.context/docs/codebase-map.json)
