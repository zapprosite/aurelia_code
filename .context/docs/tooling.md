# Tooling

A infraestrutura de desenvolvimento do Aurelia é minimalista e eficiente.

## Core Stack
- **Go 1.25+**: Linguagem principal do runtime.
- **SQLite**: Persistência local (modernc.org/sqlite).
- **slog**: Logging estruturado nativo de Go.
- **Systemd**: Gestão de processos (daemon de usuário).

## Ferramentas de Scripting
- **Bash**: Scripts de build, deploy e monitoramento (`scripts/`).
- **Sed**: Transformação de templates de unit services.

## Integrações (Externas)
- **Telegram Bot API**: Interface de usuário.
- **LLM Providers**: Google Gemini, Anthropic Claude, OpenAI, etc.
- **Groq API**: Transcrição de áudio via Whisper.
- **MCP**: Model Context Protocol para suporte a ferramentas externas.

## MCP Servers Status (Real Estado em 2026-03-18)

**HTTP MCPs — Funcionam sem deps:**
- `cloudflare-api` ✅ — Requere token CLOUDFLARE_API_TOKEN
- `cloudflare-observability` ✅ — Requere conta autenticada
- `cloudflare-radar` ✅ — Public API, sem token obrigatório

**Stdio MCPs — Dependem de npx/binários:**
- `ai-context` — Requere `@anthropic-ai/ai-context` via npm ou instalado
- `playwright` — Requere binários Chromium (150MB+)
- `filesystem` — Requere MCP.js rodando
- `context7` — Requere biblioteca privada
- `postgres` — Requere `psql` ou driver native
- `qdrant` — Requere cliente Qdrant ou curl para :6333
- `github` — Requere token GITHUB_TOKEN

**Bootstrap Policy:**
- Desabilitar MCPs se faltar binário, token ou permissão → daemon continua estável
- Reabilitar apenas após validar dependência real
- Logs do daemon em `~/.cache/aurelia/daemon.log`
