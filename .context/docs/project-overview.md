# Project Overview

Aurelia is a local-first autonomous coding agent written in Go and designed to run as a long-lived desktop or homelab runtime. The current codebase combines a Telegram interface, a ReAct execution loop, layered memory, optional MCP integrations, and multi-agent orchestration in a single modular binary.

The repository is also an elite multi-agent workspace: `AGENTS.md`, `.agents/rules/`, `.agents/workflows/` and `.context/` are part of the operating model, not just auxiliary docs.

## Codebase Reference

> **Detailed analysis**: use [`codebase-map.json`](./codebase-map.json) as the machine-readable index, but treat the prose docs in this directory as the authoritative human summary of the current checkout.

## Quick Facts

- Root: `/home/will/aurelia`
- Module: `github.com/kocar/aurelia`
- Languages: Go (`197` files), Shell (`13` files), Markdown (`37` files), JSON (`2` files)
- Primary binary entry: [`cmd/aurelia/main.go`](../../cmd/aurelia/main.go)
- Main wiring path: [`cmd/aurelia/app.go`](../../cmd/aurelia/app.go)
- Runtime health endpoint: `GET /health` from [`internal/health/server.go`](../../internal/health/server.go)
- Canonical runtime config: `~/.aurelia/config/app.json`

## Entry Points

- [`cmd/aurelia/main.go`](../../cmd/aurelia/main.go) — CLI runtime entrypoint for daemon mode, `onboard`, and `auth openai`
- [`cmd/aurelia/app.go`](../../cmd/aurelia/app.go) — composition root for runtime bootstrap, providers, memory, Telegram, cron, MCP and health server
- [`cmd/aurelia/onboard.go`](../../cmd/aurelia/onboard.go) — interactive onboarding flow that writes the instance-local config
- [`scripts/build.sh`](../../scripts/build.sh) — production-oriented build wrapper
- [`scripts/install-user-daemon.sh`](../../scripts/install-user-daemon.sh) — user-service installation flow for `systemd --user`

## Key Exports

- [`agent.Loop`](../../internal/agent/loop.go)
- [`agent.ToolRegistry`](../../internal/agent/provider.go)
- [`config.AppConfig`](../../internal/config/config.go)
- [`memory.MemoryManager`](../../internal/memory/memory.go)
- [`persona.CanonicalIdentityService`](../../internal/persona/canonical_service.go)
- [`mcp.Manager`](../../internal/mcp/manager.go)
- [`health.Server`](../../internal/health/server.go)

## File Structure & Code Organization

- `cmd/` — program entrypoints, onboarding flow and dependency wiring
- `internal/agent/` — ReAct loop, task store, team manager, worker orchestration and recovery flows
- `internal/config/` — runtime config loading and MCP config normalization
- `internal/cron/` — persistent schedules, scheduler runtime and job delivery
- `internal/health/` — lightweight HTTP health and readiness server
- `internal/mcp/` — MCP manager, connection bootstrap, discovery and transport wrappers
- `internal/memory/` — SQLite-backed messages, facts, notes and archive
- `internal/observability/` — structured logging and redaction helpers
- `internal/persona/` — canonical identity, prompt assembly, retrieval and file synchronization
- `internal/runtime/` — instance paths, bootstrap and single-instance locking
- `internal/skill/` — skill loader, installer, router and executor
- `internal/telegram/` — input/output adapters and chat command handling
- `internal/tools/` — tool definitions and handlers exposed to the LLM
- `pkg/llm/` — concrete LLM providers and model catalog logic
- `pkg/stt/` — speech-to-text provider abstraction and Groq implementation

## Technology Stack Summary

The runtime is a Go `1.25` modular monolith. It persists local state with SQLite, uses `log/slog` for observability, exposes a small HTTP health surface, and relies on Telegram as the main user-facing transport. Optional integrations include multiple remote LLM providers, the local Ollama provider, Groq STT, Gemini fallback smoke/config support, and MCP servers loaded from JSON config.

## Development Tools Overview

Daily work happens with the Go toolchain, Bash scripts in `scripts/`, systemd user services for daemonized execution, and GitHub Actions for CI enforcement. The repository also carries its own governance layer in `.agents/` and `.context/`.

## Getting Started Checklist

1. Review [`AGENTS.md`](../../AGENTS.md) and `.agents/rules/` before changing anything non-trivial.
2. Run onboarding with `go run ./cmd/aurelia onboard`.
3. Build the binary with `./scripts/build.sh`.
4. Run the test suite with `go test ./...`.
5. Install the user daemon with `./scripts/install-user-daemon.sh` when validating the long-running runtime path.
6. Confirm operational health with `curl -fsS http://127.0.0.1:8484/health`.

## Related Resources

- [Architecture Notes](./architecture.md)
- [Development Workflow](./development-workflow.md)
- [Tooling & Productivity Guide](./tooling.md)
- [Testing Strategy](./testing-strategy.md)
