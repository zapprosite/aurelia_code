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
- todas as pendências abaixo agora devem manter também um JSON taskmaster em `docs/adr/taskmaster/`

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
| Slice 11 | agent bus com `PostgreSQL` | orchestration/runtime | sim | worker claim + lease |
| Slice 11 | dashboard de swarm em `Go` | orchestration/ui | sim | thread/task board vivo |
| Slice 11 | assistance queue para agentes ociosos | orchestration/runtime | sim | idle agent help flow |
| Slice 11 | memória semântica do swarm em `Qdrant` | memory/rag | sim | thread summary -> vector |

## Observação profissional

As pendências mais críticas agora são:

1. wake word positivo com prova humana no deploy
2. handoff fim a fim do Antigravity
3. voz oficial da Aurelia com voz consistente e autorizada
4. swarm hierárquico com dashboard e assistência entre agentes
5. desktop fallback seguro
6. clonagem autorizada com arquivo local e rollback pronto

## Ordem recomendada agora

1. **Onda 1 — Voz e experiência de resposta**
   - **Slice 7**
     - validar E2E de wake word -> STT -> resposta
     - validar prova humana no deploy
   - **Slice 10**
     - fechar voz oficial da Aurelia
     - manter rollback claro para TTS pronto/local
   - motivo: tudo aqui mexe no mesmo plano de voz, TTS, Telegram e experiência do operador
2. **Onda 2 — Orquestração segura de browser e Antigravity**
   - **Slice 4**
     - handoff Antigravity fim a fim
     - menos retrabalho entre chat leve e runtime
   - **Slice 2**
     - fluxo de login guiado seguro
   - motivo: browser, Antigravity e roteamento leve compartilham as mesmas fronteiras operacionais
3. **Onda 3 — Swarm hierárquico da Aurélia**
   - **Slice 11**
     - `agent_bus` em `PostgreSQL`
     - dashboard leve em `Go`
     - assistance queue entre agentes ociosos
     - memória derivada em `Qdrant`
     - copiar o contrato útil de `open-agent-supervisor` e `langgraph-supervisor`
   - motivo: depende muito mais do runtime, da orquestração e da memória do que do desktop fallback
4. **Onda 4 — Desktop fallback seguro**
   - **Slice 3**
     - click seguro
     - digitação segura
     - kill-switch e limite de passos
   - motivo: desktop é o caminho mais frágil e deve entrar por último entre os blocos core

## ADRs já abertas para pendências críticas

- [ADR-20260319-browser-safe-login.md](./ADR-20260319-browser-safe-login.md) — login guiado seguro no browser
- [ADR-20260319-antigravity-handoff-e2e.md](./ADR-20260319-antigravity-handoff-e2e.md) — handoff fim a fim com o Antigravity
- [ADR-20260319-desktop-safe-fallback.md](./ADR-20260319-desktop-safe-fallback.md) — click/digitação/kill-switch do desktop fallback
- [ADR-20260319-voice-e2e-proof-live.md](./ADR-20260319-voice-e2e-proof-live.md) — prova live do voice plane
- [ADR-20260319-hierarchical-agent-bus.md](./ADR-20260319-hierarchical-agent-bus.md) — bus do swarm hierárquico
- [20260319-voice-capture-plane.md](./20260319-voice-capture-plane.md) — cobre o próximo slice real de captura de voz
- [ADR-20260319-voice-capture-runtime.md](./ADR-20260319-voice-capture-runtime.md) — slice nonstop em execução para integrar o capture worker ao runtime
- [ADR-20260319-state-memory-runtime.md](./ADR-20260319-state-memory-runtime.md) — persistência de gateway state e transcripts locais
- [ADR-20260319-deploy-gateway-voice.md](./ADR-20260319-deploy-gateway-voice.md) — rollout contínuo em `/home/will/aurelia-24x7`
- [ADR-20260319-extensions-governance.md](./ADR-20260319-extensions-governance.md) — política final de extensões
- [ADR-20260319-aurelia-media-voice.md](./ADR-20260319-aurelia-media-voice.md) — transcript de mídia e voz oficial da Aurelia
- [ADR-20260319-aurelia-authorized-voice-clone.md](./ADR-20260319-aurelia-authorized-voice-clone.md) — execução autorizada da voz oficial

## Estado atual do modo `/adr-semparar`

Todas as pendências estruturais abertas no backlog agora possuem:

- ADR legível em `docs/adr/`
- JSON taskmaster em `docs/adr/taskmaster/`

Isso permite continuidade entre Codex, Claude e Antigravity sem perder o próximo passo operacional.
