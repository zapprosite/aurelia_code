# Glossary & Domain Concepts

This glossary captures the terms and exported concepts that recur across the runtime, governance model and operational scripts.

## Type Definitions

- [`agent.Tool`](../../internal/agent/provider.go) — schema exposed to the LLM for callable runtime capabilities
- [`agent.ToolCall`](../../internal/agent/provider.go) — a model-requested tool invocation
- [`agent.ModelResponse`](../../internal/agent/provider.go) — normalized provider response with content, reasoning and tool calls
- [`config.AppConfig`](../../internal/config/config.go) — normalized runtime config including secrets and managed paths
- [`memory.Message`](../../internal/memory/memory.go), [`memory.Fact`](../../internal/memory/memory.go), [`memory.Note`](../../internal/memory/memory.go) — persisted conversation and long-term memory entities
- [`persona.CanonicalIdentity`](../../internal/persona/loader.go) — canonical persona representation used in prompt assembly
- [`cron.CronJob`](../../internal/cron/types.go) and [`cron.CronExecution`](../../internal/cron/types.go) — persisted schedule definitions and execution records
- [`mcp.ToolSpec`](../../internal/mcp/manager.go) — MCP tool metadata projected into the runtime

## Enumerations

- [`agent.TaskStatus`](../../internal/agent/team_types.go) — exported string status for multi-agent tasks
- Most other enum-like values in the codebase are intentionally internal and unexported, such as service-control actions, monitor metrics and onboarding steps.

## Core Terms

- **AURELIA_HOME** — instance root, usually `~/.aurelia`, where config, DBs, logs and runtime state live.
- **Instance Lock** — single-instance guard persisted by [`internal/runtime/instance_lock.go`](../../internal/runtime/instance_lock.go).
- **Canonical Identity** — merged persona and memory representation used to build prompts.
- **ReAct Loop** — reasoning-and-acting loop in [`internal/agent/loop.go`](../../internal/agent/loop.go).
- **Tool Registry** — live runtime list of capabilities exposed to the LLM.
- **Master Team** — orchestration model for spawning and coordinating internal worker tasks.
- **MCP** — Model Context Protocol integration layer, loaded from config and surfaced as tools when enabled.
- **Skill** — Markdown-packaged capability loaded from global or project-specific skill directories.
- **Heartbeat** — periodic self-check execution managed by [`internal/heartbeat/service.go`](../../internal/heartbeat/service.go).

## Acronyms & Abbreviations

- **LLM** — Large Language Model
- **STT** — Speech-to-Text
- **MCP** — Model Context Protocol
- **CI** — Continuous Integration
- **ADR** — Architecture Decision Record

## Personas / Actors

The primary external actor is the Telegram user allowed by runtime config. Internally, the system also models a master agent, worker tasks, scheduled executions and optional MCP-backed external tools. The repository itself adds another layer of actors through `AGENTS.md`, where different AI executors are assigned distinct responsibilities.

## Domain Rules & Invariants

- Only one Aurelia runtime instance should hold the instance lock at a time.
- Runtime config is instance-local and file-backed, not repo-local `.env`.
- Tool availability must come from the live registry, not from assumptions embedded in prompts.
- Disabled MCP servers must stay disabled unless explicitly re-enabled through config.

## Related Resources

- [Project Overview](./project-overview.md)
- [Architecture Notes](./architecture.md)
