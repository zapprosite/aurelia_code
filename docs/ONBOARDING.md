# Aurélia Onboarding — 2026

## TL;DR (5 minutos)

```bash
git clone https://github.com/kocar/aurelia.git
cd aurelia
chmod +x iniciar.sh
./iniciar.sh              # seleciona modo: Soberano (GPU) ou Lite
go build -o aurelia .
./aurelia onboard          # TUI interativa — configura tudo
./aurelia live             # inicia Telegram + pipeline de voz
```

---

## Requisitos

| Componente | Soberano (GPU) | Lite |
|------------|----------------|------|
| **Go** | 1.23+ | 1.23+ |
| **GPU** | NVIDIA 8GB+ | — |
| **Docker** | ✅ para Qdrant, Redis | ✅ |
| **Telegram Bot** | ✅ Token obrigatório | ✅ Token obrigatório |
| **Supabase** | ✅ opcional | ✅ opcional |
| **Qdrant** | ✅ vector DB | ✅ opcional |
| **Kokoro TTS** | `localhost:8012` | `localhost:8012` |
| **Groq API** | ✅ para STT | ✅ para STT |

---

## Estrutura do Projeto

```
aurelia/
├── cmd/aurelia/          — entry points: CLI, onboard, live, tutor
├── internal/              — core packages
│   ├── agent/            — squad orchestration, task store, worker runtime
│   ├── audio/            — buffer, player, manager, segmented synth
│   ├── config/           — config loading, env overrides, JSON persistence
│   ├── mcp/              — MCP manager (connect, discovery, transport)
│   ├── memory/            — Qdrant-backed semantic memory
│   ├── observability/    — structured logging (slog)
│   ├── skill/            — skill audit, executor, router
│   ├── streaming/         — actor pipeline, VAD monitor
│   ├── tts/              — streaming TTS (Kokoro/Kodoro/Voxtral)
│   ├── voice/            — capture, spool, tutor, mirror
│   └── ...
├── pkg/                   — shared adapters
│   ├── llm/              — OpenAI, OpenRouter, Groq, MiniMax
│   ├── stt/              — Local faster-whisper, Groq Whisper
│   └── tts/              — Kokoro factory, segmented synth
├── configs/              — LiteLLM config, systemd units
├── scripts/             — operational scripts (backup, secret-audit, health)
├── services/aurelia-api/ — HTTP API (separate Go module)
├── frontend/             — dashboard React (separate)
├── docs/                — ADRs, governance, guides
└── docker-compose.yml   — Infra (Qdrant, Redis, Supabase)
```

---

## Modos de Execução

### `./aurelia onboard` — Configuração TUI

Executa o onboarding interativo com cursor de terminal. Configura:
1. LLM provider + modelo
2. STT provider (Groq Whisper)
3. Telegram bot token
4. Utilizadores autorizados (IDs Telegram)
5. Iterações máximas e janela de memória

**Não perguntado ainda** (setados via env vars):
- `OLLAMA_URL` — endpoint local do Ollama
- `AURELIA_MODE=sovereign|lite`
- `TTS_BASE_URL` — endpoint Kokoro (default: `http://127.0.0.1:8012`)
- `VOICE_ID` — voz TTS (default: `aurelia-jarvis`)

### `./aurelia live` — Modo Interativo

Inicia o bot Telegram + pipeline de voz 24/7:
- Groq Whisper para STT
- Kokoro TTS para resposta em áudio PT-BR
- Memória semântica via Qdrant
- Personas configuráveis

### `./aurelia tutor` — JARVIS Tutor

Modo de tutoring por voz com:
- Wake word "jarvis"
- VAD (Voice Activity Detection)
- Segmentação de frases
- TTS streaming Kokoro PT-BR

### `./aurelia mcp` — MCP Server Standalone

Executa como MCP server standalone via stdio:
```bash
# Claude Desktop integration
./aurelia mcp --stdio
```

---

## Variáveis de Ambiente Essenciais

```bash
# === OBRIGATÓRIOS ===
TELEGRAM_BOT_TOKEN=          # Token do bot Telegram
GROQ_API_KEY=                # Whisper STT (free tier: 15req/min)

# === SOBERANO (GPU) ===
OLLAMA_URL=http://127.0.0.1:4000/v1   # Ollama local (RTX 4090)
OLLAMA_MODEL=qwen3.5                 # ou qwen2.5-72b

# === LITE (Cloud) ===
OPENROUTER_API_KEY=          # Cascade de modelos cloud
ANTHROPIC_API_KEY=          # ou Claude direto

# === VOZ (ambos modos) ===
TTS_BASE_URL=http://127.0.0.1:8012    # Kokoro/Kodoro TTS
VOICE_ID=aurelia-jarvis                    # voz PT-BR clonada

# === OPCIONAIS (defaults funcionam) ===
AURELIA_HOME=~/.aurelia
LOG_LEVEL=info
AURELIA_MODE=sovereign
QDRANT_URL=http://localhost:6333
REDIS_URL=localhost:6379
```

---

## Comandos Essenciais

```bash
# Build
go build -o aurelia .

# Setup (primeira vez)
./aurelia onboard

# Executar
./aurelia live              # Telegram + voz
./aurelia tutor             # JARVIS tutor (voice-only)
./aurelia serve             # HTTP API server

# Operações
./aurelia mcp --stdio       # MCP server standalone
./aurelia voice enqueue <file.wav>   # enqueue TTS audio
./aurelia voice capture-once           # captura uma vez

# Health
curl http://localhost:8484/health
curl http://localhost:6333/health

# Docker (infra)
docker compose up -d         # Qdrant + Redis
docker compose down

# Audit de secrets
bash scripts/secret-audit.sh
```

---

## Troubleshooting

### "MCP server failed to connect"
```bash
# Verificar .mcp.json em ~/.aurelia/config/
cat ~/.aurelia/config/mcp_servers.json

# Rebuild MCP config
./aurelia onboard
```

### "TTS returned empty audio"
```bash
# Verificar Kokoro está a correr
curl http://localhost:8012/health

# Testar sintetização direta
curl -X POST http://localhost:8012/v1/audio/speech \
  -H "Content-Type: application/json" \
  -d '{"model":"kokoro","input":"Olá, eu sou a Aurélia.","voice":"aurelia-jarvis"}' \
  --output test.wav
```

### "STT transcription failed"
```bash
# Verificar Groq API
curl -X POST "https://api.groq.com/openai/v1/audio/transcriptions" \
  -H "Authorization: Bearer $GROQ_API_KEY" \
  -F "file=@test.wav" \
  -F "model=whisper-large-v3"
```

### "Qdrant connection refused"
```bash
# Iniciar Qdrant
docker compose up -d qdrant

# Ver logs
docker compose logs qdrant
```

---

## Próximos Passos

Após setup inicial:
1. Lê [`docs/adr/README.md`](adr/README.md) para ADRs activos
2. Lê [`docs/governance/MODEL-STACK-POLICY.md`](governance/MODEL-STACK-POLICY.md) para entender o stack de modelos
3. Lê [`docs/governance/SECRETS.md`](governance/SECRETS.md) para segurança de secrets
4. Verifica `scripts/secret-audit.sh` antes de fazer git push

---

*Actualizado: 2026-03-31 — criado como parte do cleanup MCP/A2A 2026*
