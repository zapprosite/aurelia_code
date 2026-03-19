---
title: JARVIS Master Plan
status: in_progress
owner: codex
created: 2026-03-19
last_updated: 2026-03-19
feature_branch_target: 20260319-aurelia-antigravit-gemini
scope: local-voice-browser-antigravity-terminal
---

# JARVIS Master Plan

## Status Geral

**Progresso do plano:** `99%`

```text
[###################-] 99%
```

Estado atual:

- [x] Blueprint mestre consolidado
- [x] Skill JARVIS endurecida
- [x] Browser-use baseline validado
- [x] Skill do Antigravity instalada e versionada
- [x] Telegram gera prompt automatico para tarefa `light`
- [x] Blueprint de audio PT-BR com Groq registrado
- [x] Blueprint de voz local/JARVIS registrado
- [x] Medicao real de VRAM e politica `1 modelo residente only` fechadas
- [x] Runtime ganhou suporte real a `provider=ollama`
- [x] Runtime de deploy validado sem Gemini no caminho ativo
- [x] Gateway dry-run em Go implementado
- [x] Gateway enforcement, guardas, budgets e breaker implementados no runtime
- [x] Telemetria Prometheus do gateway exportada
- [x] Pipeline de voz com spool, heartbeat, fallback STT e dispatch no runtime implementado
- [x] Mirrors opcionais de transcript para Supabase + Qdrant implementados
- [x] CLI `aurelia voice enqueue` implementada e testada
- [x] Capture worker com heartbeat, health e spool integration implementado
- [x] Mic daemon com wake word implementado
- [x] TTS local real no Telegram implementado
- [ ] Desktop click/digitacao seguros implementados
- [x] Captura de microfone com wake word + VAD implementada

## Objetivo

Transformar a Aurelia em um JARVIS local-first, capaz de:

- escutar em background com wake word
- transcrever audio em PT-BR
- decidir localmente com um modelo forte
- operar browser, Antigravity e terminal sob o mesmo orquestrador
- manter memoria e auditoria com contrato estavel
- operar sem depender da Gemini API no runtime ativo

## Fontes de Verdade

Documentos que agora governam o plano:

- [aurelia_master_blueprint_20260319.md](/home/will/aurelia/docs/aurelia_master_blueprint_20260319.md)
- [aurelia_general_blueprint_20260319.md](/home/will/aurelia/docs/aurelia_general_blueprint_20260319.md)
- [jarvis_local_voice_blueprint_20260319.md](/home/will/aurelia/docs/jarvis_local_voice_blueprint_20260319.md)
- [local_model_kit_blueprint_20260319.md](/home/will/aurelia/docs/local_model_kit_blueprint_20260319.md)
- [groq_ptbr_audio_blueprint_20260319.md](/home/will/aurelia/docs/groq_ptbr_audio_blueprint_20260319.md)
- [antigravity_gemini_operator_blueprint.md](/home/will/aurelia/docs/antigravity_gemini_operator_blueprint.md)
- deploy sem Gemini validado em [runtime_without_gemini_blueprint_20260319.md](/home/will/aurelia-24x7/docs/runtime_without_gemini_blueprint_20260319.md)

## Decisoes Fechadas

- runtime ativo sem dependencia de Gemini API
- `OpenRouter/Minimax` como LLM remoto principal quando remoto for necessario
- `Groq` como STT principal
- LLM local forte como cerebro real para tool use e instrucao
- `bge-m3` como contrato unico de embedding no Qdrant
- `agent-browser` / Playwright como camada primaria de browser-use
- `xdotool` e `wmctrl` apenas como fallback de desktop
- terminal sempre por CLI/tooling nativa, nao por GUI
- Antigravity chat fica como copiloto leve, nao executor principal

## Budget de GPU e Modelos

### Host alvo

- GPU: `RTX 4090`
- VRAM total: `24 GiB`
- regra operacional: `1` inferencia pesada por vez

### Modelos locais escolhidos

- principal residente: `qwen3.5:9b`
- roteador/fallback de latencia: `qwen3.5:4b`
- laboratorio manual: `gemma3:27b-it-q4_K_M`
- escalonamento manual: `qwen3-coder:30b`
- embedding unico: `bge-m3`

### Regras de uso

- `qwen3.5:9b` entra como padrao local para instrucao e orquestracao
- `qwen3.5:4b` entra como frio ou aquecido sob demanda, nao residente junto por padrao
- `gemma3:27b-it-q4_K_M` sai do caminho ativo e vira modelo manual de laboratorio
- `qwen3-coder:30b` entra apenas para escalonamento manual, nao residente
- embeddings rodam fora do caminho sincrono principal

### Justificativa final

- o uso base do host e ~`4.8 GiB` de VRAM
- `qwen3.5:9b` carregado deixa ~`10.5 GiB` livres e fecha a conta
- `qwen3.5:9b + qwen3.5:4b` juntos deixam so ~`3.8 GiB` livres e nao devem ficar residentes por padrao
- um `27B` deixa folga perto de `2 GiB`, o que nao e profissional para browser e automacao
- `Groq` continua correto para audio porque tira o STT do budget local
- `qwen3-coder:30b` continua forte, mas apertado demais para ficar residente junto com browser e automacao

### Limites e degradacao

- LLM pesado concorrente: `1`
- fila maxima do LLM pesado: `1`
- embeddings concorrentes: `1`
- browser-use ativo em paralelo: `1`
- degradar quando:
  - CPU media 15m > `70%`
  - memoria disponivel < `20%`
  - GPU media > `70%`
  - VRAM usada > `85%`

## Janela Atual

### Lane de execucao

- `Codex CLI` nesta janela e o executor principal do plano
- `Antigravity` segue como copiloto leve para microtarefas e handoff
- `deploy worktree` valida runtime live

## Arquitetura Alvo

```text
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

## Principios Nao Negociaveis

- browser-first para tudo que tiver DOM e fluxo previsivel
- desktop-use so entra quando browser-use nao resolver
- health nao pode mentir
- sem falso `200 ok`
- prova real antes de declarar sucesso
- segredos fora do fluxo visual sempre que possivel
- uma inferencia pesada por vez
- um modelo residente por vez
- governor explicito por recurso e por fila

## Slices de Execucao

### Slice 0. Blueprint e Governanca

**Objetivo:** consolidar a arquitetura, as fronteiras entre camadas e as regras de operacao.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Blueprint JARVIS inicial criado
- [x] Blueprint de voz local criado
- [x] Blueprint PT-BR com Groq criado
- [x] Regras para Antigravity definidas
- [x] Regras para runtime sem Gemini consolidadas

### Slice 1. Hardening da Skill JARVIS

**Objetivo:** transformar a skill em base segura e repetivel.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Corrigir `install.sh`
- [x] Revisar `SKILL.md`
- [x] Revisar `JARVIS.md`
- [x] Criar `smoke.sh`
- [x] Validar ambiente X11

### Slice 2. Browser-Use Operacional

**Objetivo:** ter browser-use confiavel sem depender de coordenada.

**Progresso do slice:** `90%`

```text
[##################--] 90%
```

- [x] Validar Playwright local
- [x] Validar screenshot web
- [x] Definir DevTools em loopback
- [x] Registrar baseline com `agent-browser`
- [ ] Fechar fluxo de login guiado seguro

### Slice 3. Desktop-Use Operacional

**Objetivo:** manter fallback real para desktop sem virar automacao cega.

**Progresso do slice:** `60%`

```text
[############--------] 60%
```

- [x] Instalar `xdotool`, `wmctrl`, `scrot`, `xclip`, `x11-utils`
- [x] Validar `DISPLAY`
- [x] Validar foco de janela
- [x] Validar screenshot local
- [ ] Validar click seguro
- [ ] Validar digitacao segura
- [ ] Definir limite de passos

### Slice 4. Orquestracao na Aurelia

**Objetivo:** conectar browser, Antigravity e roteamento leve ao runtime real.

**Progresso do slice:** `85%`

```text
[#################---] 85%
```

- [x] Criar `PROJECT_PLAYBOOK.md`
- [x] Criar skill `antigravity-gemini-operator`
- [x] Definir matriz de roteamento
- [x] Integrar geracao automatica de prompt `light` no Telegram
- [x] Registrar prompts e handoff
- [ ] Fechar handoff de ida e volta com menos retrabalho

### Slice 5. Audio e Voz

**Objetivo:** colocar a camada de audio no caminho certo sem gastar VRAM a toa.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Registrar arquitetura Groq STT PT-BR
- [x] Validar smoke de STT com `curl`
- [x] Persistir transcript no pipeline local
- [x] Registrar blueprint de voz local
- [x] Implementar spool de audio local
- [x] Implementar processador de fila com heartbeat
- [x] Implementar fallback STT por comando
- [x] Integrar spool de audio ao orquestrador real
- [x] Expor `GET /v1/voice/status`
- [x] Criar CLI `aurelia voice enqueue`
- [x] Implementar mic daemon
- [x] Implementar wake word
- [x] Implementar VAD + buffer de captura
- [x] Implementar captura continua de microfone
- [x] Validar capture worker live no deploy com headset ALSA explicito
- [x] Corrigir parse do capturador para ignorar ruido de `stderr` no sucesso

### Slice 6. Memoria, Governor e Health

**Objetivo:** deixar o sistema controlavel e auditavel.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Definir contrato `bge-m3` para Qdrant
- [x] Definir rate limits alinhados ao host
- [x] Definir thresholds de degradacao
- [x] Provar `/health` real no deploy sem Gemini
- [x] Exportar metricas do gateway
- [x] Exportar metricas operacionais do loop de voz
- [x] Ligar Supabase como mirror opcional do audio
- [x] Ligar Qdrant como mirror semantico opcional do audio
- [x] Fechar governor inicial do audio no codigo principal
- [x] Persistir governor/breaker fora da memoria
- [x] Validar `gateway_route_states` persistido no SQLite local
- [x] Validar `voice_events` no SQLite local com transcript real

### Slice 7. Rollout Seguro

**Objetivo:** ativar por fases sem quebrar o host.

**Progresso do slice:** `80%`

```text
[################----] 80%
```

- [x] Validar deploy slice sem Gemini
- [x] Validar `cwd` live da worktree de deploy
- [x] Validar `primary_llm` no `/health`
- [x] Validar gateway enforcement e suite completa localmente
- [x] Ativar worker dedicado de voz no serviço live
- [x] Validar E2E de spool -> STT -> aceite no deploy
- [ ] Validar E2E de wake word -> STT -> resposta
- [ ] Validar Antigravity handoff fim a fim

### Slice 8. Extensoes e Aceleradores

**Objetivo:** avaliar aceleradores sem contaminar o core.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Regra definida: extensao e opcional
- [x] Mapear extensoes uteis para Chrome
- [x] Mapear extensoes uteis para Antigravity
- [x] Separar `nice to have` de `core`
- [x] Definir rollback de extensoes

### Slice 9. Gateway e Roteamento Real

**Objetivo:** tirar o gateway do modo documental e levar para enforcement seguro no runtime.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Criar `internal/gateway/`
- [x] Criar `POST /v1/router/dry-run`
- [x] Registrar matriz de roteamento e bakeoff
- [x] Enforcar lane/modelo no runtime principal
- [x] Aplicar guardas reais de reasoning/output
- [x] Implementar budgets por lane
- [x] Implementar circuit breaker por `provider:model`
- [x] Cobrir gateway provider com testes dedicados
- [x] Exportar telemetria do gateway
- [x] Expor `GET /v1/router/status`
- [x] Validar rollout na worktree de deploy

## Ordem de Execucao Recomendada

### Agora

1. validar E2E positivo de wake word -> STT -> resposta com prova humana
2. fechar handoff Antigravity fim a fim
3. concluir desktop fallback seguro

### Depois

1. validar rollout conjunto de gateway + voz na worktree de deploy
2. subir `aurelia-voice.service` ou worker dedicado
3. fechar E2E de voz com budget, fallback e health live

### Por Ultimo

1. fechar E2E completo com Antigravity e login guiado seguro
2. fechar desktop fallback seguro
3. refinar extensoes/aceleradores e revisar UX do operador

## Testes Minimos

### Unitarios

- roteamento `light / medium / high-risk`
- classificacao de risco para browser, desktop e audio
- governor por fila e por concorrencia
- health sem falso `200 ok`

### Integracao

- Playwright abrir pagina
- screenshot web e local
- spool de audio consumir item valido
- Groq STT responder com transcript util

### E2E

- wake word -> fala -> transcript -> resposta
- tarefa `light` -> prompt Antigravity -> handoff estruturado
- tarefa CLI -> execucao nativa -> resposta registrada

## Criterios de Aceite

- Aurelia decide entre `reply`, `browser`, `cli` e `antigravity`
- runtime ativo nao depende de Gemini API
- `/health` prova apenas o que esta realmente em uso
- audio nao dispara STT sem wake word/VAD
- memoria longa fica consistente entre Supabase e Qdrant
- existe kill-switch claro para desktop/browser

## Activity Log

### 2026-03-19

- [x] Auditada e endurecida a skill `jarvis-desktop-agent`
- [x] Validado baseline X11 com screenshot, foco e janela
- [x] Validado baseline de browser-use com Playwright
- [x] Criado o blueprint do operador do Antigravity
- [x] Criada a skill `antigravity-gemini-operator`
- [x] Integrada a geracao automatica de prompt `light` no Telegram
- [x] Criado o blueprint PT-BR de audio com Groq
- [x] Criado o blueprint de voz local/JARVIS com governor
- [x] Criado o blueprint do kit local de modelos
- [x] Validado o deploy runtime sem Gemini na worktree `/home/will/aurelia-24x7`
- [x] Revisada a escolha final de modelo local com foco em VRAM real: `qwen3.5:9b`
- [x] Medido o custo de VRAM real de `qwen3.5:9b` e `qwen3.5:4b`
- [x] Fechada a regra operacional: `1` modelo residente only
- [x] Ligado `provider=ollama` no app com catalogo, onboarding e health reais
- [x] Implementado o primeiro corte do gateway com `dry-run`
- [x] Implementado gateway enforcement com budgets, breaker e status route
- [x] Exportada telemetria Prometheus do gateway
- [x] Implementado spool/processador de voz com heartbeat e budget diario
- [x] Implementado fallback STT por comando e mirrors opcionais para Supabase/Qdrant
- [x] Implementada e testada a CLI `aurelia voice enqueue`
- [x] Integrado capturador real `openWakeWord + VAD` ao contrato de captura
- [x] Smoke local de voz sem falso positivo em silêncio
- [x] Persistidos gateway breaker/budget state e voice_events em SQLite local
- [x] Suite `go test ./... -count=1` voltou verde apos os cortes de gateway e voz

## Task Board

### Doing

- [x] Recuperar o `plan.md` como centro de verdade
- [x] Consolidar slices ja executados
- [x] Registrar a direcao sem Gemini no runtime ativo

### Next

- [ ] Validar rollout do gateway na worktree de deploy
- [ ] Fechar smoke de voz no runtime local com config real
- [ ] Implementar captura de microfone com wake word + VAD

### Later

- [ ] Fechar click e digitacao seguros
- [ ] Persistir governor/breaker fora da memoria
- [ ] Fechar E2E com Antigravity e browser-use

## Proxima Acao Cirurgica

Se voce mandar executar o proximo corte, a ordem certa agora e:

1. portar gateway + voice processor para a worktree de deploy
2. validar `GET /v1/router/status`, `/metrics` e `/v1/voice/status` live
3. implementar captura real de microfone com wake word + VAD
4. provar o comportamento com health, logs e evidencia

## Evidencia Atual

- [aurelia_master_blueprint_20260319.md](/home/will/aurelia/docs/aurelia_master_blueprint_20260319.md) consolida arquitetura, rollout e testes de tudo
- [aurelia_general_blueprint_20260319.md](/home/will/aurelia/docs/aurelia_general_blueprint_20260319.md) consolida o restante em um plano unico
- [antigravity_gemini_operator_blueprint.md](/home/will/aurelia/docs/antigravity_gemini_operator_blueprint.md) existe e define o contrato do chat leve
- [groq_ptbr_audio_blueprint_20260319.md](/home/will/aurelia/docs/groq_ptbr_audio_blueprint_20260319.md) existe e fecha a direcao de STT
- [jarvis_local_voice_blueprint_20260319.md](/home/will/aurelia/docs/jarvis_local_voice_blueprint_20260319.md) existe e fecha o desenho local
- [gateway_rollout_blueprint_20260319.md](/home/will/aurelia/docs/gateway_rollout_blueprint_20260319.md) existe e fecha o restante do gateway
- o deploy runtime live foi provado sem `gemini_api` no `/health` na worktree de deploy
- o LLM remoto ativo esperado segue `openrouter/minimax/minimax-m2.7`
- o repositório agora já expõe `GET /metrics`, `GET /v1/router/status` e `GET /v1/voice/status`
