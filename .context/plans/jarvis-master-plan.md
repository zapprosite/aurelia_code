---
title: JARVIS Master Plan
status: in_progress
owner: antigravity
created: 2026-03-19
last_updated: 2026-03-19
feature_branch_target: 20260319-aurelia-antigravit-gemini
scope: local-voice-browser-antigravity-terminal
---

# JARVIS Master Plan

## Status Geral

**Progresso do plano:** `83%`

```text
[#################---] 83%
```

Estado atual:

- [x] Blueprint mestre consolidado
- [x] Skill JARVIS endurecida
- [x] Browser-use baseline validado
- [x] Skill do Antigravity instalada e versionada
- [x] Telegram gera prompt automatico para tarefa `light`
- [x] Blueprint de audio PT-BR com Groq registrado
- [x] Blueprint de voz local/JARVIS registrado
- [x] Kit local `gemma3:27b-it-q4_K_M` instalado e smokeado
- [x] Runtime ganhou suporte real a `provider=ollama`
- [x] Runtime de deploy validado sem Gemini no caminho ativo
- [ ] Mic daemon com wake word implementado
- [ ] Desktop click/digitacao seguros implementados
- [ ] Persistencia completa em Supabase + Qdrant implementada

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

- [jarvis_local_voice_blueprint_20260319.md](/home/will/aurelia/docs/jarvis_local_voice_blueprint_20260319.md)
- [local_model_kit_blueprint_20260319.md](/home/will/aurelia/docs/local_model_kit_blueprint_20260319.md)
- [groq_ptbr_audio_blueprint.md](/home/will/aurelia/groq_ptbr_audio_blueprint.md)
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

- principal pesado: `gemma3:27b-it-q4_K_M`
- alternativa pesada: `gemma3:12b:27b-q4_K_M`
- escalonamento manual: `gemma3:27b-coder:30b`
- leve/roteador/fallback: `gemma3:12b`
- embedding unico: `bge-m3`

### Regras de uso

- `gemma3:27b-it-q4_K_M` entra como padrao local para instrucao e orquestracao
- `gemma3:27b-it-q4_K_M` fica como alternativa técnica, não como default
- `OpenRouter` entra apenas para escalonamento manual, não residente
- `gemma3:12b` entra para roteamento curto, degradacao e resposta leve
- embeddings rodam fora do caminho sincrono principal

### Justificativa final

- `gemma3:27b-it-q4_K_M` encaixa melhor quando o local e orquestrador e o executor pesado fica externo
- `gemma3:12b` residente (+ browser) deixa ram para compilação local
- `gemma3:12b` continua forte e residente junto com browser e automação

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

- `Antigravity` nesta janela é o executor principal do plano
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

**Progresso do slice:** `45%`

```text
[#########-----------] 45%
```

- [x] Registrar arquitetura Groq STT PT-BR
- [x] Validar smoke de STT com `curl`
- [x] Persistir transcript no pipeline local
- [x] Registrar blueprint de voz local
- [ ] Implementar mic daemon
- [ ] Implementar wake word
- [ ] Implementar VAD + ring buffer
- [ ] Integrar spool de audio ao orquestrador

### Slice 6. Memoria, Governor e Health

**Objetivo:** deixar o sistema controlavel e auditavel.

**Progresso do slice:** `40%`

```text
[########------------] 40%
```

- [x] Definir contrato `bge-m3` para Qdrant
- [x] Definir rate limits alinhados ao host
- [x] Definir thresholds de degradacao
- [x] Provar `/health` real no deploy sem Gemini
- [ ] Ligar Supabase como source of truth do audio
- [ ] Ligar Qdrant no caminho semantico real
- [ ] Exportar metricas operacionais do loop de voz
- [ ] Fechar governor no codigo principal

### Slice 7. Rollout Seguro

**Objetivo:** ativar por fases sem quebrar o host.

**Progresso do slice:** `35%`

```text
[#######-------------] 35%
```

- [x] Validar deploy slice sem Gemini
- [x] Validar `cwd` live da worktree de deploy
- [x] Validar `primary_llm` no `/health`
- [ ] Subir `aurelia-voice.service`
- [ ] Validar E2E de wake word -> STT -> resposta
- [ ] Validar Antigravity handoff fim a fim

### Slice 8. Extensoes e Aceleradores

**Objetivo:** avaliar aceleradores sem contaminar o core.

**Progresso do slice:** `10%`

```text
[##------------------] 10%
```

- [x] Regra definida: extensao e opcional
- [ ] Mapear extensoes uteis para Chrome
- [ ] Mapear extensoes uteis para Antigravity
- [ ] Separar `nice to have` de `core`
- [ ] Definir rollback de extensoes

## Ordem de Execucao Recomendada

### Agora

1. implementar `aurelia-voice.service`
2. plugar `openWakeWord + Silero VAD + ring buffer`
3. integrar o spool de audio ao runtime principal

### Depois

1. fechar click e digitacao seguros
2. ligar Supabase + Qdrant no caminho real
3. exportar metricas do governor

### Por Ultimo

1. fechar E2E completo com Antigravity
2. refinar extensoes/aceleradores
3. revisar UX do operador

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
- [x] Revisada a escolha final de modelo local com foco no papel de orquestrador: `gemma3:27b-it-q4_K_M`
- [x] Executado o kit local do Ollama com `gemma3:27b-it-q4_K_M`, `gemma3:12b` e `bge-m3`
- [x] Validado smoke real do modelo local com `ctx=8192` retornando `OK`
- [x] Ligado `provider=ollama` no app com catalogo, onboarding e health reais

## Task Board

### Doing

- [x] Recuperar o `plan.md` como centro de verdade
- [x] Consolidar slices ja executados
- [x] Registrar a direcao sem Gemini no runtime ativo

### Next

- [ ] Implementar `aurelia-voice.service`
- [ ] Implementar wake word local
- [ ] Integrar VAD + ring buffer
- [ ] Definir contrato do spool de audio

### Later

- [ ] Fechar click e digitacao seguros
- [ ] Ligar Supabase + Qdrant no caminho real
- [ ] Fechar E2E com Antigravity e browser-use

## Proxima Acao Cirurgica

Se voce mandar executar o proximo corte, a ordem certa agora e:

1. criar o servico de voz em background
2. implementar wake word + VAD
3. integrar o spool ao orquestrador principal
4. provar o fluxo com health e evidencia

## Evidencia Atual

- [antigravity_gemini_operator_blueprint.md](/home/will/aurelia/docs/antigravity_gemini_operator_blueprint.md) existe e define o contrato do chat leve
- [groq_ptbr_audio_blueprint.md](/home/will/aurelia/groq_ptbr_audio_blueprint.md) existe e fecha a direcao de STT
- [jarvis_local_voice_blueprint_20260319.md](/home/will/aurelia/docs/jarvis_local_voice_blueprint_20260319.md) existe e fecha o desenho local
- o deploy runtime live foi provado sem `gemini_api` no `/health` na worktree de deploy
- o LLM remoto ativo esperado segue `openrouter/minimax/minimax-m2.7`
