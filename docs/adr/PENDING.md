# PENDING — Slices Aguardando Implementação

Última atualização: 01/04/2026

## P1 — Crítico 🔴

| Slice | Descrição | Dependência |
|---|---|---|
| S-40 | gemma3:27b pull + Ollama register (HF token) | gemma3 no Ollama |
| S-41 | configs/litellm/config.yaml — router 3 camadas | S-40 |
| S-42 | cron 5h homelab — seed_crons.go ajuste | S-41 |

## P2 — Alto 🟡

| Slice | Descrição | Dependência |
|---|---|---|
| S-43 | PostgreSQL + pgvector substituir Supabase local | — |
| S-44 | Grafana dashboard rota ativa LiteLLM | S-41 |
| S-45 | .aurelia/onboard.md commit no repo | — |
| S-46 | .opencode/agents/ aurelia + fiscal | — |

## P3 — Médio 🟢

| Slice | Descrição | Dependência |
|---|---|---|
| S-47 | Computer Use E2E (BUA-style) | — |
| S-48 | OS Native God Mode | S-47 |
| S-49 | Jarvis Voice + Computer | S-48 |
