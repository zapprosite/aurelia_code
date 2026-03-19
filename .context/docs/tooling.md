# Tooling & Productivity Guide

The day-to-day tooling for this repository is intentionally simple: Go, Bash, systemd, SQLite-backed local state and a small set of operational scripts. Productivity comes more from disciplined workflows than from heavy framework scaffolding.

## Required Tooling

- **Go 1.25+** — builds and tests the runtime
- **Bash** — required for the scripts in [`scripts/`](../../scripts)
- **systemd user services** — expected supervision path for the long-running daemon
- **curl / jq** — useful for health checks and JSON inspection
- **npm / npx** — needed when MCP servers are launched through Node-based commands
- **GitHub Actions awareness** — contributors should understand the active CI workflows

## Recommended Automation

Key scripts shipped in the repository:

- [`scripts/build.sh`](../../scripts/build.sh) — builds the Linux binary
- [`scripts/install-user-daemon.sh`](../../scripts/install-user-daemon.sh) — installs and restarts the user daemon
- [`scripts/daemon-status.sh`](../../scripts/daemon-status.sh) — shows service status
- [`scripts/daemon-logs.sh`](../../scripts/daemon-logs.sh) — tails daemon logs
- [`scripts/health-check.sh`](../../scripts/health-check.sh) — homelab/system health snapshot
- [`scripts/gemini-smoke.sh`](../../scripts/gemini-smoke.sh) — validates the local Gemini API key, lists available models and performs a minimal generate-content check without changing the active provider
- [`scripts/smoke-test-homelab.sh`](../../scripts/smoke-test-homelab.sh) — smoke guidance for end-to-end validation
- [`scripts/sync-ai-context.sh`](../../scripts/sync-ai-context.sh) — refreshes `ai-context` state and regenerates `.context/docs/codebase-map.json`

Recommended local loop:

1. change code in an isolated branch or worktree
2. run `go test ./...`
3. rebuild with `./scripts/build.sh`
4. run `./scripts/sync-ai-context.sh`
5. validate the runtime path you changed
6. update `.context/` with the evidence

## IDE / Editor Setup

At minimum, contributors benefit from:

- Go language support (`gopls`)
- shell syntax and lint support for Bash scripts
- quick access to `journalctl`, `systemctl --user` and local health endpoints

The repository does not currently depend on editor-specific config files for correct operation.

## Productivity Tips

- Prefer the repository scripts over ad hoc command sequences when validating service behavior.
- Use `rg` for fast code and doc discovery across `internal/`, `pkg/` and `.context/`.
- Treat `.context/` as living operational memory, especially for deployment and reboot investigations.
- Keep MCP optional in your mental model: if an MCP server is missing, the core runtime should still be understandable and testable.
- Treat `ai-context` in this repository as two layers: impact detection from the CLI, and curated Markdown docs under `.context/docs/`. The script above is the canonical bridge between them.

## Related Resources

- [Development Workflow](./development-workflow.md)
- [Testing Strategy](./testing-strategy.md)
