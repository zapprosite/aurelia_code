# Testing Strategy

Quality is maintained through a broad Go test suite, explicit end-to-end coverage for key runtime paths, and CI workflows for linting, vulnerabilities and secret scanning. The test surface is strongest around core runtime modules such as agent orchestration, config normalization, lock behavior, persona assembly, cron and provider adapters.

## Test Types

- **Unit** — standard Go tests across `cmd/`, `internal/` and `pkg/` using `*_test.go`
- **Integration** — multi-component tests for cron execution, MCP manager behavior, command execution and Telegram handlers
- **E2E** — [`e2e/e2e_test.go`](../../e2e/e2e_test.go) covers persona loop, master-team recovery and cron lifecycle
- **Smoke** — [`e2e/smoke_test.go`](../../e2e/smoke_test.go) exercises homelab-oriented scenarios through Telegram
- **Static / policy** — CI runs linting, gitleaks and govulncheck

## Running Tests

- Run everything:

```bash
go test ./...
```

- Focus on end-to-end tests:

```bash
go test ./e2e -v
```

- Run smoke guidance script:

```bash
./scripts/smoke-test-homelab.sh
```

- Check shell syntax for scripts:

```bash
bash -n scripts/*.sh
```

## Quality Gates

- Runtime-critical changes should ship with targeted Go tests.
- Config and MCP behavior must preserve explicit disabled-state semantics.
- Changes in lock acquisition, startup, scheduling or tool execution should be validated beyond compile success.
- CI workflows provide baseline checks for lint, vulnerabilities and leaked secrets before merge.

## Troubleshooting

Some tests depend on environment or external services more than ordinary unit tests:

- Telegram-driven smoke tests need a working bot token and reachable chat context.
- Provider adapter tests may use stubs, but real-provider validation still depends on external credentials.
- MCP behavior can vary based on local binaries and environment variables.

When diagnosing failures, distinguish between pure unit regressions and environment-dependent operational failures.

## Related Resources

- [Development Workflow](./development-workflow.md)
- [Tooling & Productivity Guide](./tooling.md)
