---
name: repo-health-audit
description: Varredura holística do repositório Aurélia. Interpreta cada pasta contra o stack soberano vigente (GPU/Ollama, Audio/VAD, Vision/Rod, LiteLLM Router, OpenRouter, Qdrant, Supabase, Obsidian, Telegram). Gera relatório de saúde com status ✅/⚠️/❌ por módulo.
phases: [P, E, V]
version: "1.0.0 — SOTA 2026 Q2"
---

# Repo Health Audit 🩺

Varredura holística do repositório Aurélia para validar cada módulo contra o stack soberano atual.

## Stack Soberano de Referência (2026 Q2)

| Componente | Serviço | Indicador de Saúde |
|------------|---------|-------------------|
| 🤖 Telegram | `go-telegram-bot-api` | Presença em `services/`, `TELEGRAM_BOT_TOKEN` no `.env` |
| 🎙️ Jarvis Local | `cmd/aurelia/jarvis_live.go` | `JARVIS_ALWAYS_ON` no `.env` |
| 🎤 Audio/VAD | `scripts/voice-gateway.py` + `internal/streaming/actors/vad_monitor.go` | Unix Socket `/tmp/aurelia-voice.sock` |
| 👁️ Vision | `internal/computer_use/`, `go-rod`, `mcp-servers/stagehand` | `DISPLAY` no `.env` |
| 🧠 LiteLLM Router | `docker-compose smart-router` | `LITELLM_URL`, `LITELLM_MASTER_KEY` no `.env` |
| ☁️ OpenRouter Fallback | Chamadas em `pkg/llm/openrouter.go` | `OPENROUTER_API_KEY` no `.env` |
| 🦙 GPU/Ollama | `cmd/aurelia/`, `OLLAMA_URL` | Tier 0 no smart-router |
| 🔍 Qdrant | `docker-compose qdrant` | `QDRANT_URL`, `QDRANT_API_KEY` no `.env` |
| 🗄️ Supabase | `SUPABASE_URL` | Tier Storage |
| 📔 Obsidian | `OBSIDIAN_VAULT_PATH` | Integração de conhecimento |

## Protocolo de Varredura

### Passo 1 — Listar e categorizar pastas raiz
```bash
ls -la /repo/
```
Para cada pasta/arquivo crítico, classificar como:
- ✅ Presente e íntegro
- ⚠️ Presente mas com problemas detectados
- ❌ Ausente ou corrompido

### Passo 2 — Validar módulos Go core
```bash
# Verificar compilação sem erros
cd /repo && go build ./... 2>&1
```

### Passo 3 — Verificar stack no docker-compose
```bash
docker-compose config --services
docker-compose ps
```

### Passo 4 — Verificar paridade .env / .env.example
```bash
diff <(cat .env | cut -d= -f1 | sort) <(cat .env.example | cut -d= -f1 | sort)
```

### Passo 5 — Verificar systemd services
```bash
ls configs/systemd/
```

### Passo 6 — Gerar relatório
Saída: `docs/reports/YYYYMMDD-repo-health.md`

## Critérios de Interpretação

| Status | Critério |
|--------|----------|
| ✅ | Arquivo/módulo existe, compila, chave de .env presente |
| ⚠️ | Existe mas incompleto: falta teste, documentação ou paridade .env |
| ❌ | Ausente, falha de build, ou chave crítica faltando |

## Output Obrigatório
- `docs/reports/YYYYMMDD-repo-health.md` — Relatório completo
- ADR slice atualizado com os resultados
