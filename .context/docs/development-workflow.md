# Development Workflow

The repository uses an explicit governance-first workflow. `AGENTS.md` defines the authority hierarchy, `.agents/rules/` defines operational guardrails, and `.context/` stores short-lived project memory and workflow evidence.

## Branching & Releases

- Use isolated branches or worktrees for non-trivial implementation.
- Branch prefixes in the local ruleset are `feat/`, `fix/` and `research/`.
- Avoid direct commits to `main`; the repository governance expects review-driven integration.
- Deployment-style validation currently exists both as user-service scripts and system-service worktree flows in related worktrees.

Examples of active worktree conventions in this repository family:

- `/home/will/aurelia` → main feature workspace
- `/home/will/aurelia-24x7` → deploy/runtime validation
- `/home/will/aurelia-main-merge` → integration/merge workspace

## Local Development

- Onboard instance-local config:

```bash
go run ./cmd/aurelia onboard
```

- Build the binary:

```bash
./scripts/build.sh
```

- Run all tests:

```bash
go test ./...
```

- Run the daemon in foreground for debugging:

```bash
go run ./cmd/aurelia
```

- Install or refresh the user daemon:

```bash
./scripts/install-user-daemon.sh
```

- Refresh `.context` documentation state:
  - required for structural slices, handoff and merge prep
  - optional for trivial typo/comment-only edits

```bash
./scripts/sync-ai-context.sh
```

- Open a long-running slice with ADR + continuity JSON:

```bash
./scripts/adr-slice-init.sh <slug> --title "Title"
```

- Store slice-local plans outside the repository root:
  - `.context/plans/<slice>/implementation_plan.md`
  - `.context/plans/<slice>/task.md`

- Inspect runtime health:

```bash
curl -fsS http://127.0.0.1:8484/health
```

## Code Review Expectations

Review should focus on runtime safety, regressions and evidence. The local rules emphasize:

- local discovery before assumptions
- worktree isolation for non-trivial changes
- anti-hallucination: never claim success without logs, tests or direct validation
- context hygiene: update `.context/` after completing meaningful work
- ADR discipline: structural slices should be linked to `docs/adr/` and can use `/adr-semparar` when continuity matters
- root hygiene: root-level docs should stay limited to contracts, project entry docs, global examples and the active `plan.md`

For code changes, reviewers should check startup behavior, lock safety, MCP fallback behavior, tool boundary correctness and tests for new runtime logic.

## CI & Automation

The repository currently ships GitHub Actions workflows for:

- `ci.yml`
- `golangci-lint.yml`
- `govulncheck.yml`
- `gitleaks.yml`

These workflows complement, but do not replace, local validation of daemon startup and health behavior.

For repository memory, the supported local refresh path is now [`scripts/sync-ai-context.sh`](../../scripts/sync-ai-context.sh). It runs `ai-context update --dry-run` for impact detection and regenerates `.context/docs/codebase-map.json` deterministically from the checkout.

## Onboarding Tasks

New contributors should read [`AGENTS.md`](../../AGENTS.md), inspect `.agents/rules/`, review [`docs/ARCHITECTURE.md`](../../docs/ARCHITECTURE.md), and only then start editing runtime code. For operational changes, prefer validating both foreground execution and the systemd-backed path.

## Related Resources

- [Testing Strategy](./testing-strategy.md)
- [Tooling & Productivity Guide](./tooling.md)
- [Workflow changelog](../workflow/docs/changelog-post-reboot-validation-2026-03-19.md)
