---
type: doc
name: security
description: Security policies, authentication, secrets management, and compliance requirements
category: security
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Segurança — Aurélia Sovereign 2026.2

## Segredos e .env

Todos os segredos vivem em `.env` na raiz do projeto. O arquivo **não é commitado** — `.gitignore` protege.

```
# .gitignore
.env
.env.local
```

Segredos carregados via `env_file: .env` nos containers Docker Compose.

## Auditoria de Segredos

Script `scripts/audit/audit-secrets.sh` roda em todo commit (hook). Valida:
1. Padrões regex de API keys, senhas, tokens no diff staged
2. Histórico git por segredos commitados
3. Arquivos `.env` no worktree
4. Logs por credenciais em texto plano

## Porteiro Sentinel (Input Guardrail)

- Modelo: `qwen2.5:0.5b` (Ollama)
- Cache: Redis com TTL 30 dias
- Bypass: owners (Telegram IDs em `telegram_allowed_user_ids`)
- Fluxo: texto → sha256 → Redis lookup → Porteiro LLM → classifica SAFE/UNSAFE

```
Cache hit SAFE  → пропускает (não consome LLM)
Cache hit UNSAFE → блокирует (não consome LLM)
Cache miss       → qwen2.5:0.5b inference → cache Redis → пропускает/блокирует
```

## Autenticação Telegram

- Token via `TELEGRAM_BOT_TOKEN` no .env
- Owners: `telegram_allowed_user_ids` — IDs numéricos do Telegram

## Modelos de Ameaças

| Ameaça | Mitigação |
|--------|-----------|
| Prompt injection via mensagem Telegram | Porteiro sentinel (qwen2.5:0.5b) + owner bypass |
| Exposição de API keys | .env em .gitignore + audit hook em commit |
| Vazamento de system prompt | Redis com TTL + não log de prompts no terminal |
| Memória contaminada | Compressão de contexto + reset de história |
| Secrets em logs | Script audit-secrets.sh detecta antes de commit | 

## Rate Limiting (Ollama)

```bash
OLLAMA_NUM_PARALLEL=2        # 2 requests simultâneos
OLLAMA_MAX_LOADED_MODELS=2   # Mantém só 2 modelos carregados
OLLAMA_KEEP_ALIVE=10m        # Descarrega após 10min idle
```
