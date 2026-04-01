# ADR-20260401: Fix Cron Streaming + Reporter Markdown

**Status:** ✅ Implementado
**Slice:** S-51
**Data:** 01/04/2026

## Problema
- systemMonitorCronExpr = "*/180 * * * *" inválido (minutos 0-59)
- Cron disparava GenerateStream com chat-default (Nemotron)
- Todos provedores cloud falhavam em cascade → 4 erros repetidos
- Saída do reporter sem Markdown estruturado nem link Grafana

## Decisão
- Expressão corrigida para "0 */3 * * *" (a cada 3h, minuto zero)
- Cron ops usa LLMAlias: "ops-cron" (stream:false, local, timeout:120s)
- runtime.go usa Generate() não-streaming quando alias = ops-cron
- reporter.go formata saída em MarkdownV2 com link monitor.zappro.site

## Arquivos alterados
- cmd/aurelia/seed_crons.go
- internal/cron/runtime.go (se necessário)
- internal/cron/types.go (campo LLMAlias)
- internal/homelab/reporter.go
