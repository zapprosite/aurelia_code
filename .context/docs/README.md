# Documentation Index

Este diretório é a memória técnica curta do repositório. Use-o como ponto de entrada para entender a forma atual do código, o fluxo de desenvolvimento e as decisões operacionais que aparecem em `AGENTS.md`, `.agents/rules/` e nos artefatos de workflow em `.context/workflow/`.

## Core Guides

- [Project Overview](./project-overview.md)
- [Architecture Notes](./architecture.md)
- [Development Workflow](./development-workflow.md)
- [Testing Strategy](./testing-strategy.md)
- [Glossary & Domain Concepts](./glossary.md)
- [Data Flow & Integrations](./data-flow.md)
- [Security & Compliance Notes](./security.md)
- [Tooling & Productivity Guide](./tooling.md)

## Current Repository Snapshot

- Root: `/home/will/aurelia`
- Module: `github.com/kocar/aurelia`
- Primary languages: Go (`198` files), Shell (`13` files), Markdown (`38` files), JSON (`2` files)
- Main runtime entrypoint: [`cmd/aurelia/main.go`](../../cmd/aurelia/main.go)
- Composition root: [`cmd/aurelia/app.go`](../../cmd/aurelia/app.go)
- Main architectural source: [`docs/ARCHITECTURE.md`](../../docs/ARCHITECTURE.md)
- Governance source: [`AGENTS.md`](../../AGENTS.md)
- Governance contract: [`docs/REPOSITORY_CONTRACT.md`](../../docs/REPOSITORY_CONTRACT.md)

## Repository Map

- `.agents/` — autoridade operacional local, regras, workflows e skills do workspace
- `.context/` — memória operacional, docs sintéticos e estado de workflow
- `.context/plans/` — implementation plans e task boards por slice, fora da raiz
- `cmd/` — entrypoints do binário e onboarding
- `internal/` — domínio principal, runtime, gateway, voz, ferramentas, Telegram, MCP, cron e memória
- `pkg/` — provedores LLM, STT e TTS reutilizáveis
- `scripts/` — build, instalação do daemon, health-check e smoke scripts
- `docs/` — documentação arquitetural e ADRs do produto
- `e2e/` — testes end-to-end e smoke integration
- `.github/workflows/` — CI, lint, gitleaks e govulncheck

## Document Map

| Guide | Focus | Key Inputs |
| --- | --- | --- |
| `project-overview.md` | posicionamento do projeto, entrypoints e stack | `README.md`, `go.mod`, `cmd/aurelia/*` |
| `architecture.md` | shape do sistema e limites entre módulos | `docs/ARCHITECTURE.md`, `cmd/aurelia/app.go`, `internal/*` |
| `development-workflow.md` | ciclo de trabalho, branching e revisão | `AGENTS.md`, `.agents/rules/`, scripts, CI |
| `testing-strategy.md` | estratégia de testes e gates | `*_test.go`, `e2e/`, `.github/workflows/` |
| `glossary.md` | termos, tipos e conceitos recorrentes | `internal/agent`, `internal/persona`, `internal/runtime` |
| `data-flow.md` | entrada, raciocínio, ferramentas e persistência | `internal/telegram`, `internal/agent`, `internal/memory`, `internal/health` |
| `security.md` | auth, segredos e guardrails | `internal/config`, `internal/tools`, `AGENTS.md`, CI |
| `tooling.md` | CLIs, scripts e automação diária | `scripts/`, Go toolchain, systemd, npm/npx MCP |

## Governance Highlights

- adapters finos por motor: [`CLAUDE.md`](../../CLAUDE.md), [`CODEX.md`](../../CODEX.md), [`GEMINI.md`](../../GEMINI.md), [`MODEL.md`](../../MODEL.md)
- ADR obrigatória por slice estrutural: [`docs/adr/README.md`](../../docs/adr/README.md)
- modo nonstop para slices longas: [`/adr-semparar`](../../.agents/workflows/adr-semparar.md)
- comando canônico de continuidade: [`scripts/adr-slice-init.sh`](../../scripts/adr-slice-init.sh)
- higiene da raiz: contratos e docs de entrada ficam na raiz; blueprints/ADRs/planos de slice não

## Related Resources

- [Agent Handbook](../agents/README.md)
- [Workflow changelogs](../workflow/docs/changelog-post-reboot-validation-2026-03-19.md)
- [Codebase map](./codebase-map.json)
