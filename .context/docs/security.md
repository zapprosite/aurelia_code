# Security & Compliance Notes

Security in this repository is a mix of runtime constraints and workflow governance. The product runtime is local-first and usually user-scoped, while the repository itself adds operational rules through `AGENTS.md` and `.agents/rules/`.

## Authentication & Authorization

Authentication is provider-specific and configured through [`config.AppConfig`](../../internal/config/config.go). The runtime can authenticate against Telegram, remote LLM providers, Groq STT and optional MCP servers. The local Ollama provider is the main exception: it uses the loopback endpoint and does not require an API key. Telegram access is constrained by `telegram_allowed_user_ids`, which acts as the primary authorization boundary for chat-driven control.

OpenAI also supports a separate auth mode field, `openai_auth_mode`, which manages direct API key usage.

## Secrets & Sensitive Data

Sensitive values currently live primarily in the instance-local config file `~/.aurelia/config/app.json`. That includes provider API keys and Telegram credentials when the selected providers require them. The repository should not contain those values; `gitleaks.yml` exists in CI to help catch accidental commits.

Operationally relevant sensitive surfaces include:

- provider API keys in `app.json` for remote providers
- Telegram bot token and allowed user IDs
- MCP headers and environment variables in `mcp_servers.json`
- local keyring or desktop secret storage, when used by surrounding tooling

Logs should go through [`internal/observability/observability.go`](../../internal/observability/observability.go), which provides redaction helpers to reduce accidental leakage.

## Runtime Guardrails

- The daemon is expected to run as the user, not as `root`.
- The runtime enforces single-instance execution with a lockfile.
- MCP config is normalized so explicitly disabled servers stay disabled.
- Tool execution is centralized in `internal/tools`, which is the right place for command and service safety restrictions.

## Compliance & Policies

- `AGENTS.md` defines authority hierarchy and risk tiers.
- `.agents/rules/03-tiers-autonomy.md` marks network, secrets and deploy actions as high-risk operations.
- `.github/workflows/gitleaks.yml` and `.github/workflows/govulncheck.yml` provide baseline secret and dependency scanning in CI.

## Incident Response

There is no separate formal incident-response package in the repository today. In practice, response evidence is captured through:

- structured logs under `~/.aurelia/logs/`
- systemd service status and journal entries
- workflow notes and changelogs under `.context/workflow/`

For runtime regressions, the expected pattern is diagnose locally, collect proof, apply the minimal fix, and record the outcome in `.context/`.

## Related Resources

- [Architecture Notes](./architecture.md)
- [Development Workflow](./development-workflow.md)
- [Workflow changelog](../workflow/docs/changelog-post-reboot-validation-2026-03-19.md)
