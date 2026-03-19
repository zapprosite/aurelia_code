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
| Slice 10 | voz oficial da Aurelia via MiniMax Audio | audio/voice | sim | smoke TTS + voice_id |
| Slice 10 | clonagem autorizada da voz a partir de áudio local | audio/voice | sim | consentimento + smoke real |
| Slice 7 | E2E wake word -> STT -> resposta | deploy/voice | sim | wake real |
| Slice 7 | Antigravity handoff fim a fim | orchestration | sim | prompt/handoff/resposta |

## Observação profissional

As pendências mais críticas agora são:

1. wake word positivo com prova humana no deploy
2. handoff fim a fim do Antigravity
3. desktop fallback seguro
4. voz oficial da Aurelia com voice_id autorizado
5. clonagem autorizada com arquivo local e rollback pronto

## Ordem recomendada agora

1. **Slice 7 — Prova humana do voice plane**
   - testar wake word positivo até virar resposta real
   - motivo: o runtime live já está íntegro; falta a prova humana fim a fim
2. **Slice 4 + Slice 2 — Orquestração segura**
   - handoff Antigravity fim a fim
   - fluxo de login guiado seguro
   - motivo: browser/orquestração têm ROI alto e são mais seguros que desktop fallback
3. **Slice 3 — Desktop fallback seguro**
   - click seguro
   - digitação segura
   - kill-switch e limite de passos
   - motivo: desktop é o caminho mais frágil; deve entrar por último entre os blocos core

## ADRs já abertas para pendências críticas

- [20260319-voice-capture-plane.md](./20260319-voice-capture-plane.md) — cobre o próximo slice real de captura de voz
- [ADR-20260319-voice-capture-runtime.md](./ADR-20260319-voice-capture-runtime.md) — slice nonstop em execução para integrar o capture worker ao runtime
- [ADR-20260319-state-memory-runtime.md](./ADR-20260319-state-memory-runtime.md) — persistência de gateway state e transcripts locais
- [ADR-20260319-deploy-gateway-voice.md](./ADR-20260319-deploy-gateway-voice.md) — rollout contínuo em `/home/will/aurelia-24x7`
- [ADR-20260319-extensions-governance.md](./ADR-20260319-extensions-governance.md) — política final de extensões
- [ADR-20260319-aurelia-media-voice.md](./ADR-20260319-aurelia-media-voice.md) — transcript de mídia e voz oficial da Aurelia
- [ADR-20260319-aurelia-authorized-voice-clone.md](./ADR-20260319-aurelia-authorized-voice-clone.md) — execução autorizada da voz oficial
