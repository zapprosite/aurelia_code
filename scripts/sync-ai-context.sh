#!/bin/bash

set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
OUTPUT_DIR="$ROOT_DIR/.context"
DOCS_DIR="$OUTPUT_DIR/docs"
CODEBASE_MAP="$DOCS_DIR/codebase-map.json"
TMP_OUTPUT=$(mktemp)

cleanup() {
    rm -f "$TMP_OUTPUT"
}
trap cleanup EXIT

cd "$ROOT_DIR"

if ! command -v ai-context >/dev/null 2>&1; then
    echo "❌ ai-context não encontrado no PATH."
    exit 1
fi

mkdir -p "$DOCS_DIR"

echo "==> ai-context update --dry-run"
if ai-context update "$ROOT_DIR" --dry-run >"$TMP_OUTPUT" 2>&1; then
    cat "$TMP_OUTPUT"
else
    cat "$TMP_OUTPUT"
    echo
    echo "⚠️ ai-context update falhou. Continuando com a sincronização determinística do codebase-map."
fi

count_glob() {
    local pattern=$1
    local count
    count=$(rg --files -g "$pattern" 2>/dev/null | wc -l || true)
    printf '%s' "${count// /}"
}

count_dir() {
    local dir=$1
    find "$dir" -type f | wc -l | tr -d ' '
}

TOTAL_FILES=$(find . -path ./.git -prune -o -type f -print | wc -l | tr -d ' ')
GO_COUNT=$(count_glob '*.go')
MD_COUNT=$(count_glob '*.md')
SH_COUNT=$(count_glob '*.sh')
JSON_COUNT=$(count_glob '*.json')
YML_COUNT=$(count_glob '*.yml')
YAML_COUNT=$(count_glob '*.yaml')
PNG_COUNT=$(count_glob '*.png')

AGENTS_COUNT=$(count_dir .agents)
CONTEXT_COUNT=$(count_dir .context)
CMD_COUNT=$(count_dir cmd)
DOCS_COUNT=$(count_dir docs)
INTERNAL_COUNT=$(count_dir internal)
PKG_COUNT=$(count_dir pkg)
SCRIPTS_COUNT=$(count_dir scripts)
E2E_COUNT=$(count_dir e2e)
GITHUB_COUNT=$(count_dir .github)

GENERATED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

jq -n \
  --arg generated "$GENERATED_AT" \
  --arg total "$TOTAL_FILES" \
  --arg go_count "$GO_COUNT" \
  --arg md_count "$MD_COUNT" \
  --arg sh_count "$SH_COUNT" \
  --arg json_count "$JSON_COUNT" \
  --arg yml_count "$YML_COUNT" \
  --arg yaml_count "$YAML_COUNT" \
  --arg png_count "$PNG_COUNT" \
  --arg agents_count "$AGENTS_COUNT" \
  --arg context_count "$CONTEXT_COUNT" \
  --arg cmd_count "$CMD_COUNT" \
  --arg docs_count "$DOCS_COUNT" \
  --arg internal_count "$INTERNAL_COUNT" \
  --arg pkg_count "$PKG_COUNT" \
  --arg scripts_count "$SCRIPTS_COUNT" \
  --arg e2e_count "$E2E_COUNT" \
  --arg github_count "$GITHUB_COUNT" \
  '{
    version: "1.2.0",
    generated: $generated,
    stack: {
      primaryLanguage: "go",
      languages: ["go", "markdown", "shell", "json"],
      frameworks: ["telebot.v3", "systemd --user"],
      buildTools: ["go", "bash"],
      testFrameworks: ["go test"],
      packageManager: "go modules",
      isMonorepo: false,
      hasDocker: false,
      hasCI: true
    },
    structure: {
      totalFiles: ($total | tonumber),
      rootPath: ".",
      topDirectories: [
        {name: ".agents", fileCount: ($agents_count | tonumber), description: "Governance rules, workflows, and reusable skills"},
        {name: ".context", fileCount: ($context_count | tonumber), description: "Operational memory, generated docs, and workflow state"},
        {name: ".github", fileCount: ($github_count | tonumber), description: "CI workflows and repository automation"},
        {name: "cmd", fileCount: ($cmd_count | tonumber), description: "Application entrypoints and composition root"},
        {name: "docs", fileCount: ($docs_count | tonumber), description: "Human-facing architecture, ADRs, and reviews"},
        {name: "internal", fileCount: ($internal_count | tonumber), description: "Core runtime layers and domain logic"},
        {name: "pkg", fileCount: ($pkg_count | tonumber), description: "LLM and STT provider integrations"},
        {name: "scripts", fileCount: ($scripts_count | tonumber), description: "Build, daemon, and host operations"},
        {name: "e2e", fileCount: ($e2e_count | tonumber), description: "End-to-end and smoke tests"}
      ],
      languageDistribution: [
        {extension: ".go", count: ($go_count | tonumber)},
        {extension: ".md", count: ($md_count | tonumber)},
        {extension: ".sh", count: ($sh_count | tonumber)},
        {extension: ".json", count: ($json_count | tonumber)},
        {extension: ".yml", count: ($yml_count | tonumber)},
        {extension: ".yaml", count: ($yaml_count | tonumber)},
        {extension: ".png", count: ($png_count | tonumber)}
      ]
    },
    architecture: {
      layers: [
        {name: "cmd/aurelia", description: "Process bootstrap, CLI commands, app wiring, and runtime startup"},
        {name: "internal/runtime", description: "Instance root resolution, directory bootstrap, project-local paths, and instance locking"},
        {name: "internal/agent", description: "ReAct loop, tool execution, task graph, team orchestration, and recovery"},
        {name: "internal/telegram", description: "Telegram input/output adapters, bootstrap flow, markdown rendering, and conversation pipeline"},
        {name: "internal/memory + internal/persona", description: "Durable memory, facts, notes, archive, and prompt/context composition"},
        {name: "internal/mcp + internal/tools", description: "External MCP integrations and native tool registration"},
        {name: "internal/cron + internal/health + internal/heartbeat", description: "Scheduled execution, HTTP health, and self-monitoring"}
      ],
      patterns: [
        "modular monolith",
        "tool-driven ReAct loop",
        "master-led multi-agent orchestration",
        "SQLite-backed operational state",
        "local-first runtime",
        "single-instance daemon guard"
      ],
      entryPoints: [
        "cmd/aurelia/main.go",
        "cmd/aurelia/app.go",
        "cmd/aurelia/wiring.go",
        "scripts/build.sh",
        "scripts/install-user-daemon.sh"
      ],
      mainEntryPoints: ["cmd/aurelia/main.go"],
      moduleExports: [
        "internal/agent",
        "internal/runtime",
        "internal/telegram",
        "internal/memory",
        "internal/persona",
        "pkg/llm",
        "pkg/stt"
      ]
    },
    symbols: {
      classes: [],
      interfaces: [
        "internal/agent.LLMProvider",
        "internal/tools.ScheduleCreator",
        "pkg/stt.Transcriber"
      ],
      functions: [
        "main.run",
        "runtime.AcquireInstanceLock",
        "observability.Configure",
        "agent.(*Loop).Run"
      ],
      types: [
        "runtime.PathResolver",
        "runtime.InstanceLock",
        "telegram.BotController",
        "persona.CanonicalIdentityService",
        "mcp.Manager",
        "health.Server"
      ],
      enums: [
        "internal/agent.TaskStatus"
      ]
    },
    publicAPI: [
      "cmd/aurelia onboard",
      "cmd/aurelia auth openai",
      "scripts/build.sh",
      "scripts/install-user-daemon.sh",
      "scripts/daemon-status.sh",
      "scripts/daemon-logs.sh",
      "GET /health"
    ],
    dependencies: {
      mostImported: [
        "context",
        "fmt",
        "log/slog",
        "gopkg.in/telebot.v3",
        "modernc.org/sqlite"
      ]
    },
    stats: {
      totalSymbols: null,
      exportedSymbols: null,
      analysisTimeMs: null
    },
    keyFiles: [
      {path: "cmd/aurelia/main.go", role: "CLI entrypoint and process lifecycle"},
      {path: "cmd/aurelia/app.go", role: "Composition root and runtime bootstrap"},
      {path: "cmd/aurelia/wiring.go", role: "Tool registration and team wiring"},
      {path: "internal/runtime/instance_lock.go", role: "Single-instance enforcement"},
      {path: "internal/agent/loop.go", role: "Core ReAct execution loop"},
      {path: "internal/health/server.go", role: "Health and readiness HTTP surface"}
    ],
    navigation: {
      tests: [
        "cmd/**/*_test.go",
        "internal/**/*_test.go",
        "pkg/**/*_test.go",
        "e2e/**/*.go"
      ],
      config: [
        "~/.aurelia/config/app.json",
        "~/.aurelia/config/mcp_servers.json",
        "mcp_servers.example.json"
      ],
      types: [
        "internal/agent/provider.go",
        "internal/runtime/resolver.go",
        "internal/agent/team_types.go",
        "internal/config/config.go"
      ],
      mainLogic: [
        "cmd/aurelia",
        "internal/agent",
        "internal/runtime",
        "internal/telegram",
        "internal/tools",
        "internal/mcp",
        "pkg/llm"
      ]
    }
  }' > "$CODEBASE_MAP"

echo
echo "✅ codebase-map atualizado em $CODEBASE_MAP"
echo "ℹ️ Os .md em .context/docs continuam sendo docs curatoriais. O ai-context identifica o que revisar, mas não os materializa de forma confiável via MCP nesta workspace."
