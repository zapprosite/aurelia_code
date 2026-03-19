---
description: Backlog oficial das pendências abertas por slice em 2026-03-19.
status: active
owner: codex
---

# Pending Slices Backlog

Este é o backlog oficial das pendências abertas do plano JARVIS/Aurelia.

## Regras

- toda pendência estrutural daqui exige ADR ao iniciar execução
- itens menores podem ser fechados direto se não alterarem arquitetura/runtime
- a fonte primária de status continua em [plan.md](../../plan.md)

## Pendências abertas

| Slice | Pendência | Tipo | ADR obrigatória | Teste mínimo |
| --- | --- | --- | --- | --- |
| Slice 2 | fluxo de login guiado seguro | browser/runtime | sim | smoke browser |
| Slice 3 | click seguro | desktop fallback | sim | screenshot antes/depois |
| Slice 3 | digitação segura | desktop fallback | sim | ação reversível validada |
| Slice 3 | limite de passos | segurança operacional | sim | teste de abort |
| Slice 4 | handoff de ida e volta com menos retrabalho | orchestration | sim | E2E Antigravity |
| Slice 5 | mic daemon | voz/runtime | sim | worker local vivo |
| Slice 5 | wake word real | voz/runtime | sim | ruído não dispara |
| Slice 5 | VAD + ring buffer | voz/runtime | sim | silêncio não chama STT |
| Slice 5 | captura contínua de microfone | voz/runtime | sim | job real criado sem enqueue manual |
| Slice 6 | persistir governor/breaker fora da memória | runtime/state | sim | restart preserva estado |
| Slice 6 | source of truth compartilhada fim a fim | storage/memory | sim | SQLite/Supabase/Qdrant coerentes |
| Slice 7 | rollout na worktree de deploy | deploy | sim | suite + health live |
| Slice 7 | `aurelia-voice.service` ou worker dedicado | runtime/deploy | sim | systemd worker ativo |
| Slice 7 | E2E spool -> STT -> resposta no deploy | deploy/voice | sim | resposta real |
| Slice 7 | E2E wake word -> STT -> resposta | deploy/voice | sim | wake real |
| Slice 7 | Antigravity handoff fim a fim | orchestration | sim | prompt/handoff/resposta |
| Slice 8 | mapear extensões úteis para Chrome | optional tooling | não | doc curado |
| Slice 8 | mapear extensões úteis para Antigravity | optional tooling | não | doc curado |
| Slice 8 | separar core de nice-to-have | governance | sim | matriz atualizada |
| Slice 8 | rollback de extensões | governance | sim | procedimento escrito |
| Slice 9 | rollout do gateway na worktree de deploy | deploy/gateway | sim | `/v1/router/status` live |

## Observação profissional

As pendências mais críticas hoje são:

1. voice capture real
2. persistência real de estado/governor
3. rollout de gateway/voice no deploy

## Ordem recomendada agora

1. **Slice 5 — Voice plane real**
   - fechar `mic daemon`, `wake word`, `VAD + ring buffer`, `captura contínua`
   - motivo: o capture worker já entrou; agora vale completar o lane de voz antes de expandir superfície
2. **Slice 6 — Estado e memória reais**
   - persistir `governor/breaker`
   - fechar `SQLite/Supabase/Qdrant` como truth flow coerente
   - motivo: evita levar para deploy um runtime que ainda perde estado importante ao reiniciar
3. **Slice 7 + Slice 9 — Deploy gateway/voice**
   - rollout em `/home/will/aurelia-24x7`
   - `aurelia-voice.service` ou worker dedicado
   - E2E `spool -> STT -> resposta`
   - `GET /v1/router/status` live
   - motivo: gateway e voz já devem subir juntos no ambiente real
4. **Slice 4 + Slice 2 — Orquestração segura**
   - handoff Antigravity fim a fim
   - fluxo de login guiado seguro
   - motivo: browser/orquestração têm ROI alto e são mais seguros que desktop fallback
5. **Slice 3 — Desktop fallback seguro**
   - click seguro
   - digitação segura
   - kill-switch e limite de passos
   - motivo: desktop é o caminho mais frágil; deve entrar por último entre os blocos core
6. **Slice 8 — Extensões**
   - mapear, separar core de nice-to-have e definir rollback
   - motivo: opcional, não deve contaminar o core antes do fechamento do runtime

## ADRs já abertas para pendências críticas

- [20260319-voice-capture-plane.md](./20260319-voice-capture-plane.md) — cobre o próximo slice real de captura de voz
- [ADR-20260319-voice-capture-runtime.md](./ADR-20260319-voice-capture-runtime.md) — slice nonstop em execução para integrar o capture worker ao runtime
