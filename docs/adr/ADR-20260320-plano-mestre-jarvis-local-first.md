---
title: Plano Mestre JARVIS — Aurelia Local-First
status: accepted
date: 2026-03-20
decision-makers: [humano, codex, aurelia]
supersedes: plan.md (raiz)
---

# ADR-20260320: Plano Mestre JARVIS — Aurelia Local-First

## Contexto

O projeto operava com um `plan.md` monolítico na raiz que misturava roadmap, decisões, activity log e task board. Com a padronização do repositório como template multi-agente, o plano mestre precisa ser registrado como ADR para governança adequada.

## Objetivo

Transformar a Aurelia em um JARVIS local-first, capaz de:

- Escutar em background com wake word
- Transcrever áudio em PT-BR
- Decidir localmente com um modelo forte (instruction, não code)
- Operar browser, Antigravity e terminal sob o mesmo orquestrador
- Manter memória e auditoria com contrato estável
- Operar sem depender da Gemini API no runtime ativo

## Arquitetura Alvo

```
MIC
  -> wake word local
  -> VAD + ring buffer
  -> Groq STT
  -> intent router
      -> reply
      -> Antigravity handoff
      -> browser-use
      -> CLI tools
  -> local LLM
  -> memory/audit
      -> Supabase
      -> Qdrant
```

## Decisões Fechadas

- Runtime ativo sem dependência de Gemini API.
- `OpenRouter/Minimax` como LLM remoto principal quando remoto for necessário.
- `Groq` como STT principal.
- LLM local instruction como cérebro real para tool use e orquestração de enxame.
- `bge-m3` como contrato único de embedding no Qdrant.
- `agent-browser` / Playwright como camada primária de browser-use.
- `xdotool` e `wmctrl` apenas como fallback de desktop.
- Terminal sempre por CLI/tooling nativa, não por GUI.
- Antigravity chat fica como copiloto leve, não executor principal.

## Princípios Não Negociáveis

- Browser-first para tudo que tiver DOM e fluxo previsível.
- Desktop-use só entra quando browser-use não resolver.
- Health não pode mentir — sem falso `200 ok`.
- Prova real antes de declarar sucesso.
- Segredos fora do fluxo visual sempre que possível.
- Uma inferência pesada por vez, um modelo residente por vez.
- Governor explícito por recurso e por fila.

## Slices de Execução

### ✅ Concluídos (100%)

| Slice | Descrição |
|-------|-----------|
| 0 | Blueprint e Governança |
| 1 | Hardening da Skill JARVIS |
| 5 | Áudio e Voz (spool, STT, wake word, VAD, mic daemon) |
| 6 | Memória, Governor e Health (bge-m3, rate limits, SQLite) |
| 8 | Extensões e Aceleradores |
| 9 | Gateway e Roteamento Real (dry-run, enforcement, breaker, Prometheus) |

### 🔄 Em Progresso

| Slice | Descrição | Progresso |
|-------|-----------|-----------|
| 2 | Browser-Use Operacional (falta login guiado seguro) | 90% |
| 4 | Orquestração (falta handoff Antigravity ida-volta) | 85% |
| 7 | Rollout Seguro (falta E2E wake word→resposta) | 80% |
| 3 | Desktop-Use (falta click/digitação seguros) | 60% |

### Cobertura por ADR Nonstop

- Slice 2: `ADR-20260319-browser-safe-login`
- Slice 3: `ADR-20260319-desktop-safe-fallback`
- Slice 4: `ADR-20260319-antigravity-handoff-e2e`
- Slice 7: `ADR-20260319-voice-e2e-proof-live`
- Slice 10: `ADR-20260319-aurelia-media-voice` + `ADR-20260319-aurelia-authorized-voice-clone`
- Slice 11: `ADR-20260319-hierarchical-agent-bus`

## Ondas de Execução

### Onda 1 — Voz e Experiência de Resposta
- [ ] Validar E2E: wake word → STT → resposta com prova humana
- [ ] Fechar voz oficial da Aurelia com rollback claro

### Onda 2 — Orquestração Segura
- [ ] Handoff Antigravity fim a fim
- [ ] Login guiado seguro no browser

### Onda 3 — Swarm Hierárquico
- [ ] `agent_bus` em PostgreSQL
- [ ] Dashboard leve em Go
- [ ] Assistance queue
- [ ] Memória derivada em Qdrant

### Onda 4 — Desktop Fallback Seguro
- [ ] Click e digitação seguros com limit de passos
- [ ] Refinar UX do operador

## Testes Mínimos

- **Unitários:** roteamento, classificação de risco, governor, health sem falso 200
- **Integração:** Playwright, screenshot, spool de áudio, Groq STT
- **E2E:** wake word→transcript→resposta, tarefa light→Antigravity→handoff, CLI→execução→registro

## Critérios de Aceite

- Aurelia decide entre `reply`, `browser`, `cli` e `antigravity`.
- Runtime ativo não depende de Gemini API.
- `/health` prova apenas o que está realmente em uso.
- Áudio não dispara STT sem wake word/VAD.
- Memória longa consistente entre Supabase e Qdrant.
- Kill-switch claro para desktop/browser.

## Consequências

- `plan.md` da raiz removido — este ADR é a fonte de verdade do roadmap.
- Activity log e task board ficam no `.context/` e nos ADRs de slice.
- Hardware e modelos ficam em ADR separado: `ADR-20260320-politica-modelos-hardware-vram.md`.
