# Architecture Notes

The current Aurelia codebase is a modular Go monolith that runs as a user-scoped autonomous agent. Its architecture centers on a single process that wires together configuration, memory, persona construction, LLM access, tool execution, Telegram I/O, scheduling, gateway routing, a queue-driven voice processor and optional MCP integrations.

## System Architecture Overview

At runtime, the process starts from [`cmd/aurelia/main.go`](../../cmd/aurelia/main.go), bootstraps instance paths and config in [`cmd/aurelia/app.go`](../../cmd/aurelia/app.go), registers tools in [`cmd/aurelia/wiring.go`](../../cmd/aurelia/wiring.go), and then starts three main surfaces:

- the Telegram bot controller
- the cron scheduler
- the HTTP health server
- the optional voice spool processor

The deployment model is local-first. State lives under `~/.aurelia/`, the runtime is expected to be supervised by systemd, and the repository itself carries governance and workflow metadata alongside product code.

## Architectural Layers

- **Entrypoint & wiring**: `cmd/aurelia/` composes the runtime and handles CLI subcommands.
- **Execution core**: `internal/agent/` owns the ReAct loop, tool contracts, team orchestration and recovery.
- **Interface layer**: `internal/telegram/` translates Telegram updates into internal messages and renders responses back.
- **Runtime services**: `internal/runtime/`, `internal/config/`, `internal/observability/`, `internal/health/` handle bootstrap, config, lockfile, logging and health.
- **Gateway & voice services**: `internal/gateway/` and `internal/voice/` handle model routing, budgets, circuit breaking, audio spool processing and transcript mirrors.
- **State & identity**: `internal/memory/` and `internal/persona/` store durable context and build prompts.
- **Automation layer**: `internal/cron/`, `internal/skill/`, `internal/tools/`, `internal/mcp/` implement scheduled work, skills, native tools and external MCP tools.
- **Provider layer**: `pkg/llm/`, `pkg/stt/` and `pkg/tts/` adapt third-party model, transcription and speech APIs, including the local Telegram TTS path and the optional MiniMax premium voice lane.

## Detected Design Patterns

| Pattern | Confidence | Locations | Description |
| --- | --- | --- | --- |
| Composition Root | High | [`cmd/aurelia/app.go`](../../cmd/aurelia/app.go) | All major runtime dependencies are instantiated in one bootstrap function. |
| Registry Pattern | High | [`internal/agent/provider.go`](../../internal/agent/provider.go), [`internal/tools/definitions.go`](../../internal/tools/definitions.go) | Tools are registered dynamically and exposed as LLM-callable capabilities. |
| Adapter Pattern | High | [`internal/telegram/`](../../internal/telegram), [`pkg/llm/`](../../pkg/llm), [`pkg/stt/`](../../pkg/stt) | External APIs are wrapped behind internal interfaces. |
| Repository / Store Pattern | High | [`internal/memory/`](../../internal/memory), [`internal/cron/store.go`](../../internal/cron/store.go), [`internal/agent/task_store.go`](../../internal/agent/task_store.go) | SQLite-backed stores isolate persistence concerns. |
| Supervisor / Recovery Pattern | Medium | [`internal/agent/master_team_service_recovery.go`](../../internal/agent/master_team_service_recovery.go), [`internal/runtime/instance_lock.go`](../../internal/runtime/instance_lock.go) | The runtime enforces single-instance execution and supports team recovery. |

## Entry Points

- [`cmd/aurelia/main.go`](../../cmd/aurelia/main.go)
- [`cmd/aurelia/app.go`](../../cmd/aurelia/app.go)
- [`cmd/aurelia/wiring.go`](../../cmd/aurelia/wiring.go)
- [`cmd/aurelia/onboard.go`](../../cmd/aurelia/onboard.go)
- [`scripts/build.sh`](../../scripts/build.sh)
- [`scripts/install-user-daemon.sh`](../../scripts/install-user-daemon.sh)

## Public API

| Symbol | Type | Location |
| --- | --- | --- |
| `Loop` | struct | [`internal/agent/loop.go`](../../internal/agent/loop.go) |
| `ToolRegistry` | struct | [`internal/agent/provider.go`](../../internal/agent/provider.go) |
| `AppConfig` | struct | [`internal/config/config.go`](../../internal/config/config.go) |
| `EditableConfig` | struct | [`internal/config/config.go`](../../internal/config/config.go) |
| `PathResolver` | struct | [`internal/runtime/resolver.go`](../../internal/runtime/resolver.go) |
| `MemoryManager` | struct | [`internal/memory/memory.go`](../../internal/memory/memory.go) |
| `CanonicalIdentityService` | struct | [`internal/persona/canonical_service.go`](../../internal/persona/canonical_service.go) |
| `Manager` | struct | [`internal/mcp/manager.go`](../../internal/mcp/manager.go) |
| `Service` | struct | [`internal/cron/service.go`](../../internal/cron/service.go) |
| `Server` | struct | [`internal/health/server.go`](../../internal/health/server.go) |

## Internal System Boundaries

`internal/telegram` should stay as an adapter layer, not become the source of truth for orchestration rules. `internal/agent` owns the reasoning loop and team execution. `internal/tools` defines the capabilities surfaced to the LLM, while `internal/mcp` is responsible only for optional external tool transport. `internal/persona` and `internal/memory` own identity and long-term context, not the interface layer.

## External Service Dependencies

- **Telegram Bot API** — primary interaction channel
- **LLM providers** — Anthropic, Google, Kimi, Kilo, Ollama, OpenAI, OpenRouter, ZAI, Alibaba
- **Groq** — speech-to-text backend
- **voice-proxy / Chatterbox TTS** — local text-to-speech backend and safe fallback for Telegram audio replies
- **MiniMax Audio** — optional premium TTS lane for the official Aurelia voice when an authorized `voice_id` and API key are configured
- **Qdrant / Supabase** — optional transcript mirrors and semantic memory targets
- **MCP servers** — optional local or remote tools loaded from JSON config
- **systemd** — supported long-running supervision model

## Key Decisions & Trade-offs

The project favors a single deployable process over split microservices. That keeps local operation simple, but it concentrates many responsibilities in the composition root and raises the importance of clear internal boundaries. Optional MCP integration is intentionally non-fatal: if MCP config is missing or disabled, the runtime still boots.

## Risks & Constraints

- The composition root is broad and accumulates boot responsibilities quickly.
- Many integrations are environment-driven, so config drift can surface at startup.
- Telegram remains the dominant UX surface; local CLI flows are mostly onboarding and auth support.
- `.context/docs/codebase-map.json` is currently sparse, so prose docs remain more useful than the machine summary.

## Top Directories Snapshot

- `cmd/` — `10` files
- `internal/` — `158` files
- `pkg/` — `27` files
- `scripts/` — `16` files
- `docs/` — `23` files
- `e2e/` — `2` files

## Related Resources

- [Project Overview](./project-overview.md)
- [Data Flow & Integrations](./data-flow.md)
- [`docs/ARCHITECTURE.md`](../../docs/ARCHITECTURE.md)
