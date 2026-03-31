---
name: latenca
description: Monitora e otimiza a latência do Homelab (Ollama) e Cloud (LiteLLM/Groq/OpenRouter).
---

# Latência Monitor (SOTA 2026)

Este skill permite auditar a performance em tempo real do ecossistema Aurélia, comparando a inferência local (Soberania) com a nuvem (Escalabilidade).

## Comandos

- `/latenca check`: Executa um benchmark completo em todos os endpoints configurados.
- `/latenca ollama`: Verifica especificamente o status e latência do Ollama local.
- `/latenca cloud`: Verifica latência dos provedores externos via LiteLLM.

## Thresholds de Referência (2026)

| Provedor | SOTA Latency (TTFT) | Status |
|----------|----------------------|--------|
| Groq     | < 150ms              | Elite  |
| Ollama   | < 500ms              | Native |
| OpenRouter| < 800ms             | Hybrid |

## Troubleshooting
Se a latência do Ollama estiver > 2000ms, verifique a contenção de VRAM via `nvidia-smi`.
Se o LiteLLM estiver lento, verifique a configuração do Smart Router na porta 4000.

---
**Criado via /master-skill**
