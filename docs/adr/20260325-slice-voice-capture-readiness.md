# ADR 2026-03-25: Slice 1 — Voice Capture Readiness

## Status
Implementada

## Objetivo

Eliminar a degradação real do runtime causada pelo `voice_capture` e fazer `/health` e `/api/status` refletirem um estado `up` quando a stack de voz estiver corretamente instalada.

## Contexto

O runtime principal já expõe status honesto. Hoje ele continua `degraded` por um motivo concreto: o comando de captura aponta para `/home/will/aurelia-24x7/scripts/voice-capture-openwakeword.sh`, arquivo ausente no host.

Isso não é problema de observabilidade. É problema de readiness operacional.

## Escopo

- corrigir o path/comando de `voice_capture`
- validar que o worker de captura sobe, pulsa heartbeat e não falha no boot
- alinhar health check e status snapshot com estados reais da stack de voz
- deixar smoke explícito para captura, processor e mirror

## Fora de escopo

- redesign completo da arquitetura de voz
- troca de provider STT/TTS
- features novas de wakeword

## Mudanças esperadas

1. descobrir a fonte canônica do script de captura
2. corrigir config/path para o comando efetivo
3. garantir fallback explícito se a captura estiver desabilitada por decisão operacional
4. reforçar o health para diferenciar `disabled` de `failed`
5. adicionar smoke para `voice_capture` e `voice_processor`

## Smoke obrigatório

```bash
go test ./internal/voice ./cmd/aurelia
curl -sS http://127.0.0.1:8484/health
curl -sS http://127.0.0.1:3334/api/status
```

## Critério de saída

- `voice_capture` não aparece mais como `degraded` por arquivo ausente
- `/health` continua `ok`
- `/api/status` passa para `up` ou `degraded` apenas por falha real nova

## Dependência

Nenhuma. Esta slice deve vir primeiro porque limpa o ruído operacional antes das slices estruturais seguintes.
