---
type: doc
name: tooling
description: Scripts, IDE settings, automation, and developer productivity tips
category: tooling
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Ferramentas e Produtividade — Aurélia Sovereign 2026.2

## Scripts Principais

| Script | Propósito |
|--------|----------|
| `scripts/audit/audit-secrets.sh` | Auditoria soberana de segredos (roda em cada commit) |
| `scripts/audit/audit-secrets.sh --check-all` | Scan completo (git history + logs + arquivos) |

## Comandos Úteis

```bash
# Build
CGO_ENABLED=0 go build -trimpath -o bin/aurelia ./cmd/aurelia

# Deploy
sudo systemctl restart aurelia

# Logs
journalctl -u aurelia -f

# VRAM
nvidia-smi

# Ollama models
ollama list

# Flush cache Redis (Porteiro)
docker exec aurelia-redis-1 redis-cli KEYS "porteiro:cache:*" | xargs docker exec aurelia-redis-1 redis-cli DEL
```

## Rate Limit Ollama (systemd)

```bash
sudo systemctl edit ollama --add-edit Section=Service
# Adicionar em /etc/systemd/system/ollama.service.d/rate-limit.conf

[Service]
Environment="OLLAMA_NUM_PARALLEL=2"
Environment="OLLAMA_MAX_LOADED_MODELS=2"
Environment="OLLAMA_KEEP_ALIVE=10m"

sudo systemctl daemon-reload && sudo systemctl restart ollama
```

## Docker Compose

```bash
# Up comdependências
docker-compose -f /home/will/aurelia/docker-compose.yml up -d

# Up serviço específico
docker-compose -f /home/will/aurelia/docker-compose.yml up -d kokoro-tts

# Logs
docker logs -f aurelia-kokoro
docker logs -f aurelia-smart-router
```

## Variáveis de Ambiente Críticas

| Variável | Serviço | Uso |
|----------|---------|-----|
| `TELEGRAM_BOT_TOKEN` | Bot | Autenticação Telegram |
| `TAVILY_API_KEY` | web_search | Pesquisa primária |
| `OPENROUTER_API_KEY` | LiteLLM | Modelos pagos e free tier |
| `GROQ_API_KEY` | LiteLLM | STT cloud + LLM free tier |
| `OLLAMA_HOST` | Ollama | Inference local |
| `LITELLM_MASTER_KEY` | LiteLLM proxy | Autenticação da API LiteLLM |
