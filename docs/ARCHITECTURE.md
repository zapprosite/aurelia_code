# ARCHITECTURE

## Context

`Aurelia` is a local-first autonomous coding agent built in Go.

The system operates through Telegram, persists operational state in SQLite, and executes work through an explicit runtime made of:

- a ReAct loop
- a live tool registry
- Agent Teams orchestration
- layered memory
- controlled local execution
- optional MCP integration

This document is the architectural source of truth for the current codebase.

Use it together with:

- `AGENTS.md`
- `docs/STYLE_GUIDE.md`
- `docs/LEARNINGS.md`

## Architectural Shape

The project is a modular monolith in Go.

```text
/cmd/aurelia/
  main.go          [thin entrypoint]
  app.go           [composition root]
  wiring.go        [dependency and tool registration]

/internal/
  agent/           [loop, tool registry, team manager, task graph, recovery]
  config/          [configuration loading]
  cron/            [store, service, scheduler, runtime]
  mcp/             [manager, discovery, transport]
  memory/          [messages, facts, notes, archive]
  persona/         [canonical identity, prompt building, retrieval, file sync]
  runtime/         [instance and project path resolution]
  skill/           [skill loading and routing]
  telegram/        [handlers, bootstrap, input pipeline, output]
  tools/           [tool definitions, handlers, MCP adapter]

/pkg/
  llm/             [LLM providers and model catalogs]
  stt/             [speech-to-text]
```

Current code still contains legacy repository names and paths.
Those are implementation details under migration, not the product identity to preserve.

## Runtime Scope Separation

The runtime distinguishes three scopes:

### Repository

Contains:

- source code
- tests
- project documentation
- default assets shipped with the product

### Local Instance

Lives outside the repository in a per-user hidden directory.

Contains:

- local config
- canonical app config in `config/app.json`
- SQLite state
- logs
- instance lock
- learned notes
- runtime skills
- canonical persona files

## Process Lifecycle And Supervision

The supported desktop runtime is a user service managed by `systemd --user`.

Rules:

- the daemon must be installed from `scripts/install-user-daemon.sh`
- the service must run from the checked-out repository root
- foreground `go run` is a debug path, not the canonical production path
- duplicate runtime instances are forbidden

Single-instance enforcement is implemented through a lock file in the instance root:

- `~/.aurelia/instance.lock`

If a second process attempts to start while the lock is held, startup fails with PID and command diagnostics instead of killing the existing process.

### Target Project

The external codebase the agent is acting on.

It may define project-specific docs and working conventions.

Project-specific rules must not leak into global defaults unless explicitly promoted.

## Layer Boundaries

### Entry And Wiring

`cmd/aurelia` is responsible for:

- loading configuration
- building services
- registering tools
- starting and stopping runtimes

It must stay thin.

It is also responsible for:

- bootstrapping instance directories
- acquiring the single-instance lock
- configuring structured logging
- wiring shutdown and release of runtime resources

Runtime configuration is instance-local and file-backed.

Rule:

- app runtime config is loaded from `~/.aurelia/config/app.json`
- LLM selection is persisted as `llm_provider` plus `llm_model`
- provider-specific auth modes are explicit config fields when needed, such as `openai_auth_mode`
- repository-local `.env` files are not part of the supported runtime config model

Current auth variants:

- `google`: `api_key`
- `kilo`: `api_key`
- `zai`: `coding_plan_api_key`
- `alibaba`: `coding_plan_api_key`
- `openai`: `api_key` or experimental local `codex` CLI mode

### Interface Layer

`internal/telegram` is the interface boundary.

Responsibilities:

- receiving Telegram events
- adapting text, files, and audio into internal input
- sending output back to Telegram

It must not become the source of truth for business rules.

### Domain And Orchestration

`internal/agent` is the center of execution behavior.

Responsibilities:

- ReAct loop
- tool execution contracts
- Agent Teams orchestration
- task graph
- recovery

### Memory And Identity

`internal/memory` persists:

- recent messages
- durable facts
- compact notes
- raw archive

`internal/persona` resolves:

- canonical identity
- prompt assembly
- file synchronization
- owner and project context injection
- retrieval helpers

Important rule:

- identity and operating rules come from canonical persona and project context
- actual runtime capabilities come from the live tool registry

### Tooling

`internal/tools` owns:

- native tool definitions
- handler contracts
- schedule tools
- MCP registration adapters
- team mailbox and control tools

Tool schemas belong here, not in the composition root.

## Core Runtime Model

### ReAct Execution

The main loop is tool-driven.

The runtime injects:

- available tools
- execution guidance
- workdir guidance
- runtime capabilities

The agent must reason from real capabilities, not assumed capabilities.

### Agent Teams

The multi-agent model is master-led.

Rules:

- `master` is the only agent that answers the end user
- workers operate on explicit tasks
- tasks carry status, ownership, dependencies, result, and canonical `workdir`
- workers may coordinate through internal mailbox tools
- operational team state should persist in SQLite whenever possible

### Memory

The memory model is deterministic and layered.

Priority order:

1. canonical identity
2. stable facts
3. episodic notes
4. recent conversation window

The architecture explicitly avoids treating larger prompts or vector search as the default fix for identity and continuity problems.

### Local Execution

The runtime is expected to observe the environment, not only describe it.

That includes:

- reading files
- writing files
- listing directories
- running controlled local commands
- acting on a canonical project `workdir`

## Observability

Logging now uses `log/slog` as the canonical process logger.

Rules:

- component names must be explicit
- sensitive values must be redacted or omitted
- tool execution logs may include tool names and argument keys, but not raw sensitive payloads
- audio and file logs should use basenames instead of absolute host paths where practical

Legacy `log.Printf` call sites may still exist, but they are bridged into the structured logger at process startup.

## Architectural Rules

1. Telegram is an interface layer, not a domain layer.
2. Identity and memory rules belong in `persona` and `memory`.
3. Multi-agent orchestration belongs in `agent`.
4. Long-lived operational state should persist in SQLite when practical.
5. Tools are explicit runtime capabilities, not hidden side channels.
6. New code should preserve the modular monolith shape instead of adding ad hoc service sprawl.
7. Workers operating on external projects must preserve canonical `workdir` across local tools.
8. Architecture changes must be reflected here before the task is considered complete.

## Current Capabilities

Implemented in the current codebase:

- tool-driven ReAct loop
- Agent Teams orchestration with task graph, mailbox, recovery, and final synthesis
- SQLite-backed memory and operational persistence
- single-instance runtime guard with lock file diagnostics
- structured logging bridge with redaction-oriented call sites
- Telegram text, markdown, and audio input flow
- cron scheduling subsystem
- MCP discovery and registration
- project-aware execution through propagated `workdir`

## Current Constraints

Known constraints in the current codebase:

- repository naming still reflects the old product identity
- ai-context regeneration is still an operational dependency outside this repository
- benchmark evidence for memory and CPU footprint is not yet documented
- contribution governance and GitHub policy gates are not fully established yet

## Documentation Policy

This file captures architecture and boundaries.

Implementation conventions belong in `docs/STYLE_GUIDE.md`.

Operational mistakes, traps, and recurring lessons belong in `docs/LEARNINGS.md`.

