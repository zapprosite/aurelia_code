# Data Flow & Integrations

The dominant data path in Aurelia starts with a Telegram update, passes through preprocessing and prompt assembly, enters the ReAct loop, optionally executes tools, persists memory, and then returns a formatted response to Telegram. Side channels include cron-triggered executions, health checks, optional MCP tool calls, and a queue-driven voice path that can inject accepted transcripts into the same Telegram execution flow.

## Module Dependencies

- `cmd/aurelia` → `internal/runtime`, `internal/config`, `internal/memory`, `internal/persona`, `internal/telegram`, `internal/cron`, `internal/mcp`, `internal/tools`, `pkg/llm`, `pkg/stt`, `pkg/tts`
- `internal/telegram` → `internal/memory`, `internal/persona`, `internal/skill`, `gopkg.in/telebot.v3`
- `internal/agent` → tool registry contracts and LLM provider interfaces
- `internal/cron` → `internal/agent`, `internal/persona`, persistent SQLite store
- `internal/mcp` → normalized MCP config and external transport clients
- `pkg/llm` / `pkg/stt` / `pkg/tts` → third-party provider SDKs and HTTP APIs

## Service Layer

- [`memory.MemoryManager`](../../internal/memory/memory.go) — persists conversational and long-term context
- [`persona.CanonicalIdentityService`](../../internal/persona/canonical_service.go) — assembles canonical persona and retrieval context
- [`agent.Loop`](../../internal/agent/loop.go) — orchestrates model calls and tool execution
- [`agent.MasterTeamService`](../../internal/agent/master_team_service.go) — handles multi-agent team execution
- [`cron.Service`](../../internal/cron/service.go) and [`cron.Scheduler`](../../internal/cron/scheduler.go) — manage scheduled jobs
- [`mcp.Manager`](../../internal/mcp/manager.go) — connects and exposes MCP-backed tools
- [`telegram.BotController`](../../internal/telegram/bot.go) — receives, dispatches and sends Telegram messages
- [`health.Server`](../../internal/health/server.go) — exposes `/health` and `/ready`
- [`gateway.Provider`](../../internal/gateway/provider.go) — applies the routing matrix, budgets and breakers across local and remote models
- [`voice.Processor`](../../internal/voice/processor.go) — consumes voice spool jobs, applies STT budget/fallback, gates on wake phrase and dispatches accepted commands

## High-level Flow

1. Telegram delivers a user message to [`internal/telegram/input.go`](../../internal/telegram/input.go) and related handlers.
2. Audio inputs are optionally transcribed through [`pkg/stt/groq.go`](../../pkg/stt/groq.go).
3. Persona and memory layers assemble the system prompt and contextual history.
4. [`agent.Loop`](../../internal/agent/loop.go) calls the selected LLM provider with the live tool definitions.
5. If the model requests tools, the runtime executes handlers from [`internal/tools/`](../../internal/tools) or MCP-provided adapters.
6. Outputs are formatted through Telegram renderers; when `requiresAudio=true`, the controller synthesizes speech through the configured TTS provider, using `voice-proxy` as the current local fallback and MiniMax as an optional premium lane.
7. Messages, notes, facts and archives are stored in SQLite for later retrieval.

Voice path in the current checkout:

1. An audio file is enqueued into [`internal/voice/spool.go`](../../internal/voice/spool.go), currently through `aurelia voice enqueue`.
2. [`internal/voice/processor.go`](../../internal/voice/processor.go) claims the oldest job, applies Groq STT and falls back to a local command when budget or cooldown blocks the primary path.
3. The transcript is sanitized, checked for the configured wake phrase, mirrored optionally to Supabase/Qdrant, and then injected into [`telegram.BotController.ProcessExternalInput`](../../internal/telegram/input_pipeline.go).

## Internal Movement

Data moves mostly through direct function calls rather than queues. The main persistent boundaries are SQLite databases for memory, cron state and team tasks. Team orchestration creates its own task graph and mailbox model inside `internal/agent`, while cron schedules replay prompts through the same execution core used for interactive chat.

## External Integrations

- **Telegram** — inbound updates and outbound messages
- **LLM APIs / local Ollama** — reasoning and tool-call generation through the internal gateway
- **Groq STT** — audio transcription
- **voice-proxy / Chatterbox TTS** — local OpenAI-compatible speech synthesis for Telegram output and fallback
- **MiniMax Audio** — optional premium TTS lane for the official Aurelia voice profile
- **Supabase / Qdrant** — optional transcript mirrors for shared audit and semantic retrieval
- **MCP servers** — optional remote or local tool capabilities
- **HTTP health consumers** — local monitoring via `/health` and `/ready`

## Observability & Failure Modes

Structured logging is provided by [`internal/observability/observability.go`](../../internal/observability/observability.go). Startup can fail on config load, missing auth material, lock acquisition or provider initialization. MCP is intentionally soft-fail: missing or disabled MCP config should not stop the core runtime. Health checks surface uptime and registered dependency checks, but only checks registered explicitly by the runtime appear in `/health`; when `provider=ollama`, the primary LLM check now validates both the local endpoint and the configured model presence.

## Related Resources

- [Architecture Notes](./architecture.md)
- [Security & Compliance Notes](./security.md)
- [Testing Strategy](./testing-strategy.md)
