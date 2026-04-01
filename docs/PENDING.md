# PENDING — Slices Aguardando Implementação
Última auditoria: 01/04/2026 17:40

> POLÍTICA: Apenas slices com código no repo são marcados ✅.
> Este arquivo é a única fonte de verdade do backlog.

## P0 — Hotfix Ativo 🔥 (impede o bot de funcionar)

| Slice | Descrição | Arquivo |
|---|---|---|
| S-51 | Fix cron expr + streaming + reporter Markdown | seed_crons, runtime, reporter |

## P1 — Crítico 🔴 (infra sem a qual o home lab não é estável)

| Slice | Descrição | Pré-requisito |
|---|---|---|
| S-52 | PostgreSQL + pgvector substituir Supabase local | — |
| S-53 | Grafana + Prometheus no docker-compose.yml | S-52 |
| S-54 | TTS voz feminina PT-BR natural (Edge TTS nativo) | — |
| S-55 | Smoke test E2E: texto → voz → cron com Markdown | S-51, S-54 |

## P2 — Alto 🟡 (qualidade da experiência)

| Slice | Descrição | Pré-requisito |
|---|---|---|
| S-56 | Ubuntu Desktop: voice-gateway integrado ao pipeline | S-51 |
| S-57 | /status handler retorna reporter.Format() em Markdown | S-51 |
| S-58 | Notificação proativa de anomalia (GPU >85°C, VRAM >90%) | S-53 |

## P3 — Médio 🟢 (expansão)

| Slice | Descrição | Pré-requisito |
|---|---|---|
| S-59 | Computer Use E2E (BUA-style) | — |
| S-60 | OS Native God Mode | S-59 |
| S-61 | Jarvis Voice + Computer | S-60 |
