# Aurelia → JARVIS

Bot Telegram em Go, local-first, agêntico. Evolui para JARVIS com voz + swarm de agentes.

**Stack:** Go 1.25+ · Ollama · Qdrant · Postgres · SQLite · Telegram · Groq STT · Kokoro TTS

---

## Setup rápido

```bash
go run ./cmd/aurelia onboard      # configura ~/.aurelia/config/app.json
./scripts/build.sh                # binário ~23 MB
./scripts/install-user-daemon.sh  # systemd --user
./scripts/daemon-status.sh        # health check
```

Runtime: ~25 MB RAM idle, <20 ms startup, instância única via `flock`.

---

## Autoridade

Leia nesta ordem antes de qualquer mudança:

1. [`AGENTS.md`](./AGENTS.md) — hierarquia e papéis dos agentes
2. [`docs/governance/AURELIA-AUTHORITY-DECLARATION.md`](./docs/governance/AURELIA-AUTHORITY-DECLARATION.md) — permissões totais (sudo=1) e missão JARVIS
3. [`docs/governance/REPOSITORY_CONTRACT.md`](./docs/governance/REPOSITORY_CONTRACT.md) — contrato do repositório
4. [`docs/adr/ADR-20260320-roadmap-mestre-slices.md`](./docs/adr/ADR-20260320-roadmap-mestre-slices.md) — slices por prioridade (P1→P8)

---

## Modelos (AMD RX 7900 XTX / ROCm)

| Papel | Modelo | VRAM |
|-------|--------|------|
| Residente | `gemma3:12b` | ~8.1 GB |
| Fallback 262K ctx | `qwen3.5:9b` | ~6.7 GB |
| Escalonamento | `gemma3:27b-it-q4_K_M` | ~17 GB |
| Embedding | `bge-m3` | ~1.2 GB |
| STT | Groq (remoto) | 0 |
| TTS | Kokoro (CPU) | 0 |

Detalhes: [`docs/adr/ADR-20260320-politica-modelos-hardware-vram.md`](./docs/adr/ADR-20260320-politica-modelos-hardware-vram.md)

---

## Agentes e Workflows

| Engine | Adaptador | Papel |
|--------|-----------|-------|
| Antigravity (Gemini) | `GEMINI.md` | Cockpit e orquestração |
| Claude | `CLAUDE.md` | Execução complexa |
| Codex | `CODEX.md` | Execução rápida |
| OpenCode | `OPENCODE.md` + `.opencode/` | Local-first, multi-provider |

Comandos rápidos (no IDE):
```
/git-ship      → commit semântico + push + PR
/git-feature   → nova branch com verificações de segurança
/sincronizar-tudo → commit sênior completo
/adr-semparar  → slice longa com ADR + JSON de continuidade
/sincronizar-ai-context → sync .context/ + codebase-map
```

Workflows completos: `.agents/workflows/` | Skills: `.agents/skills/`

---

## Estrutura

```
cmd/           → entrypoints Go
internal/      → pacotes internos (telegram, voice, health, tools)
pkg/           → pacotes públicos (llm, config)
scripts/       → build, daemon, memory-sync, secret-audit
docs/
  adr/         → 3 ADRs ativos (roadmap, jarvis, modelos)
  governance/  → authority, contract, secrets
.agents/
  workflows/   → /git-ship, /dev, /qa, /pm, etc.
  skills/      → voice-clone, sync-ai-context, adr-nonstop-slice, etc.
.context/      → memória de trabalho (ai-context MCP)
```

---

## Memória e Observabilidade

```bash
# Memory sync (automático via systemd timers)
bash scripts/memory-sync-fiscal.sh --mode fast      # 5 min
bash scripts/memory-sync-fiscal.sh --mode validate  # diário 6h

# Auditoria de secrets (rodar antes de push)
bash scripts/secret-audit.sh

# Logs do daemon
journalctl --user -u aurelia.service -f
```

---

## Segurança

- `~/.aurelia/config/secrets.env` — fora do repo, `chmod 600`
- `.gitignore` protege `.env*` e `secrets.env`
- Secrets: [`docs/governance/SECRETS.md`](./docs/governance/SECRETS.md)
- KeePassXC vault: deadline 2026-03-27

---

## Desenvolvimento

```bash
go test ./... -count=1 -short   # testes rápidos
go build ./...                  # verifica compilação
./scripts/adr-slice-init.sh <slug> --title "Título"  # nova slice estrutural
```

Branch: `feat/`, `fix/`, `research/` — worktrees isoladas para slices não-triviais.
Commits de agente: sempre `--no-verify` (hooks lentos; CI valida via GitHub Actions).
