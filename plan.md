---
title: JARVIS Desktop Agent Blueprint
status: in_progress
owner: codex
created: 2026-03-19
last_updated: 2026-03-19
feature_branch_target: 20260319-aurelia-antigravit-gemini
scope: local-desktop-control-browser-use-orchestration
---

# JARVIS Blueprint

## Status Geral

**Progresso do plano:** `58%`

```text
[###########---------] 58%
```

Estado atual:

- [x] Pedido interpretado
- [x] Skill JARVIS auditada em `~/.aurelia/skills/jarvis-desktop-agent/`
- [x] Riscos principais identificados
- [x] Blueprint inicial registrado
- [x] Arquitetura endurecida implementada
- [x] Instalação validada
- [x] Browser-use validado
- [ ] Desktop-use validado
- [ ] Integração com Aurelia validada

## Objetivo

Transformar a Aurelia em um agente estilo JARVIS local-first, capaz de:

- ver a tela
- operar navegador
- executar ações de desktop
- combinar browser-use, computer-use e automação local
- agir com guardrails fortes

Sem virar um agente cego, destrutivo ou exposto.

## Adendo Operacional

Contexto confirmado em `2026-03-19`:

- estamos operando dentro do `Antigravity`
- o ambiente menciona `subagent preview`
- podemos considerar extensões que ajudem no `Google Chrome`
- podemos considerar extensões que ajudem no próprio `Antigravity`

Regra de uso para esse adendo:

- extensão é acelerador opcional
- extensão não vira dependência central do JARVIS
- toda extensão proposta deve ter propósito claro, baixo risco e reversão simples

## Princípios Não Negociáveis

- `Playwright` primeiro para navegação web
- `Chrome DevTools` apenas em `127.0.0.1`
- `xdotool` apenas como fallback para desktop real
- screenshot antes e depois de ações relevantes
- limite de passos por tarefa
- confirmação humana para ações destrutivas
- segredos fora do fluxo visual sempre que possível
- nada de depender por padrão de `--dangerously-skip-permissions`

## Diagnóstico Cirúrgico

### O que já existe

- skill local criada em `~/.aurelia/skills/jarvis-desktop-agent/`
- documentação inicial de comportamento
- instalador de dependências
- fluxo pensado para screenshot, mouse, teclado e Chrome

### O que está fraco hoje

- o `install.sh` está logicamente incorreto no tratamento de `sudo`
- a skill ainda é descritiva, não uma integração real com o runtime da Aurelia
- faltam guardrails operacionais
- faltam testes de saúde por camada
- falta separar claramente browser-use de desktop-use
- falta política explícita para segredos, formulários e ações irreversíveis

## Arquitetura Alvo

```text
USER
  -> AURELIA ORCHESTRATOR
      -> SLICE A: planner / intent classifier
      -> SLICE B: browser-use adapter
      -> SLICE C: desktop-use adapter
      -> SLICE D: screenshot + state verification
      -> SLICE E: action safety policy
      -> SLICE F: execution log + recovery
```

### Camadas

#### Slice A. Intent Router

Decide se a tarefa vai para:

- pesquisa/documentação
- browser-use
- desktop-use
- shell/tool use
- escalonamento humano

Saída esperada:

- plano curto da ação
- ferramenta primária
- risco da ação
- necessidade ou não de confirmação

#### Slice B. Browser-Use First

Para sites e apps web:

- usar `Playwright`
- usar snapshot estruturado
- clicar por referência, não por coordenada
- usar screenshot apenas como apoio visual

Casos alvo:

- login guiado
- navegação em painel web
- pesquisa estilo Perplexity
- preenchimento de formulários não sensíveis

#### Slice C. Desktop-Use Fallback

Para o que não está acessível por DOM/web tooling:

- `xdotool`
- `wmctrl`
- screenshots
- foco de janela

Casos alvo:

- menus nativos
- janelas X11
- apps que não expõem DOM
- diálogos do sistema

#### Slice D. Verification Loop

Ciclo obrigatório:

1. observar
2. decidir
3. agir
4. verificar
5. continuar ou abortar

Sem verificação, não há sucesso declarado.

#### Slice E. Safety Policy

Regras mínimas:

- bloquear ações destrutivas sem confirmação
- bloquear digitação de segredos em tela sem opt-in explícito
- bloquear cliques repetitivos fora de contexto
- bloquear execução infinita
- registrar evidência de cada passo relevante

#### Slice F. Execution Ledger

Cada tarefa deve deixar trilha curta:

- intenção
- ferramentas usadas
- número de passos
- screenshots ou referências
- falhas
- motivo do abort

## Slices de Execução

### Slice 0. Blueprint e Guardrails

**Objetivo:** deixar o plano e o contrato de segurança definidos antes de instalar qualquer coisa.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Auditar skill existente
- [x] Identificar falhas de arquitetura
- [x] Definir camadas
- [x] Registrar blueprint em `plan.md`

### Slice 1. Hardening da Skill JARVIS

**Objetivo:** transformar a skill em algo executável e seguro.

**Progresso do slice:** `100%`

```text
[####################] 100%
```

- [x] Corrigir `install.sh`
- [x] Revisar `SKILL.md`
- [x] Revisar `JARVIS.md`
- [x] Remover instruções perigosas como padrão
- [x] Definir pré-checagens de ambiente
- [x] Definir rollback simples

Entregáveis:

- skill revisada
- instalador correto
- documentação alinhada

### Slice 2. Browser-Use Operacional

**Objetivo:** ter navegação real confiável sem depender de clique por coordenada.

**Progresso do slice:** `80%`

```text
[################----] 80%
```

- [x] Validar Playwright local
- [x] Definir fluxo padrão de navegação
- [x] Definir fluxo padrão de screenshot
- [ ] Definir fluxo de login guiado
- [x] Definir política para remote debugging local

Entregáveis:

- guia de browser-use
- smoke test web
- política de uso do Chrome DevTools

### Slice 3. Desktop-Use Operacional

**Objetivo:** controlar desktop local apenas quando browser-use não resolver.

**Progresso do slice:** `50%`

```text
[##########----------] 50%
```

- [x] Instalar `xdotool`, `wmctrl`, `scrot`, `xclip`, `x11-utils`
- [x] Validar `DISPLAY`
- [x] Validar foco de janela
- [ ] Validar click seguro
- [ ] Validar digitação segura
- [ ] Definir limite de passos

Entregáveis:

- smoke test de desktop
- fallback operacional X11

### Slice 4. Orquestração na Aurelia

**Objetivo:** conectar JARVIS ao runtime real da Aurelia.

**Progresso do slice:** `0%`

```text
[--------------------] 0%
```

- [ ] Definir roteamento por intenção
- [ ] Definir contrato entre planner e executor
- [ ] Registrar políticas de segurança
- [ ] Registrar telemetria mínima
- [ ] Definir mensagens de falha claras

Entregáveis:

- blueprint de integração
- pontos de entrada definidos

### Slice 5. Observabilidade e Recuperação

**Objetivo:** manter controle sobre o agente quando algo sair do trilho.

**Progresso do slice:** `0%`

```text
[--------------------] 0%
```

- [ ] Adicionar health checks
- [ ] Adicionar timeout por task
- [ ] Adicionar limite de ações
- [ ] Adicionar kill-switch
- [ ] Adicionar logs de execução

Entregáveis:

- plano de recuperação
- critérios de abort

### Slice 6. Rollout Seguro

**Objetivo:** ativar por fases, sem liberar poder total de uma vez.

**Progresso do slice:** `0%`

```text
[--------------------] 0%
```

- [ ] Fase 1: screenshot only
- [ ] Fase 2: browser-use only
- [ ] Fase 3: desktop click seguro
- [ ] Fase 4: formulários simples
- [ ] Fase 5: fluxos compostos

Entregáveis:

- rollout faseado
- checklist de aceite por fase

### Slice 7. Extensões e Aceleradores

**Objetivo:** avaliar extensões úteis sem acoplar o JARVIS a plugins frágeis.

**Progresso do slice:** `0%`

```text
[--------------------] 0%
```

- [ ] Mapear extensões úteis para Chrome
- [ ] Mapear extensões úteis para Antigravity
- [ ] Separar “nice to have” de “core”
- [ ] Definir política de instalação e rollback

Entregáveis:

- shortlist de extensões
- critérios de aceite
- política de remoção rápida

## Ordem de Execução Recomendada

### Agora

1. corrigir a skill JARVIS
2. endurecer o instalador
3. formalizar guardrails

### Depois

1. validar browser-use
2. validar desktop-use
3. integrar com o runtime da Aurelia

### Por último

1. automações compostas
2. fluxos mais ambiciosos estilo Perplexity/Jarvis
3. refinamento de UX

## Testes Mínimos

### Unitários

- validação de decisão de rota
- classificação de risco
- limites de passo
- política de confirmação

### Integração

- screenshot local
- Playwright abrir página
- Chrome DevTools local
- `xdotool` mover e clicar em ambiente controlado

### E2E

- abrir navegador
- navegar para página de teste
- preencher campo não sensível
- verificar resultado
- abortar corretamente em caso de falha

## Critérios de Aceite

- Aurelia consegue escolher browser-use ou desktop-use corretamente
- tarefas web simples não usam coordenadas quando DOM existe
- tarefas desktop não rodam sem verificação
- segredos não entram em fluxo visual por padrão
- existe kill-switch claro
- existe trilha curta de execução

## Activity Log

### 2026-03-19

- [x] Auditada a skill `jarvis-desktop-agent`
- [x] Identificado bug de `sudo` no `install.sh`
- [x] Definida arquitetura em slices
- [x] Registrado blueprint inicial
- [x] Reescritos `install.sh`, `jarvis-chrome.sh`, `SKILL.md`, `JARVIS.md`, `README.md` e `QUICK_REFERENCE.md`
- [x] Criado `smoke.sh` para validação segura do ambiente
- [x] Instalados `scrot`, `wmctrl`, `xclip` e `xdotool`
- [x] Validado `DISPLAY`, enumeração de janelas e screenshot local
- [x] Ativado Chrome DevTools isolado em `127.0.0.1:9222`
- [x] Validado bootstrap de Chrome isolado visual e headless
- [ ] Estabilidade de DevTools ainda insuficiente para navegação contínua
- [x] Validado Playwright headless com Chrome do sistema em `https://example.com/`
- [x] Registrado adendo operacional: Antigravity + subagent preview + extensões opcionais
- [ ] Próximo passo recomendado: click e digitação seguros

## Task Board

### Doing

- [x] Criar blueprint cirúrgico
- [x] Estruturar progresso e slices
- [x] Registrar log operacional
- [x] Hardening inicial da skill JARVIS
- [x] Validar ambiente base X11

### Next

- [ ] Validar click e digitação seguros
- [ ] Definir smoke test web de login guiado
- [ ] Validar click e digitação seguros

### Later

- [ ] Integrar ao runtime da Aurelia
- [ ] Adicionar observabilidade e replay curto
- [ ] Expandir para fluxos compostos

## Como Vou Anotar Durante a Execução

Sempre que uma etapa for executada:

- marcar o item com `[x]`
- atualizar a barra de progresso geral
- atualizar a barra do slice afetado
- registrar uma linha no `Activity Log`
- deixar o próximo passo explícito

## Próxima Ação Cirúrgica

Se você mandar executar, eu começo por este corte:

1. corrigir `install.sh`
2. endurecer `SKILL.md` e `JARVIS.md`
3. criar smoke tests mínimos
4. registrar tudo neste `plan.md`

## Evidência Atual

- `smoke.sh` passou para `curl`, `jq`, `scrot`, `wmctrl`, `xclip`, `xdotool` e `DISPLAY`
- `wmctrl -l` listou janelas reais
- `xdotool getactivewindow` respondeu corretamente
- `scrot` gerou screenshot válido em `/tmp/jarvis-smoke.png`
- `jarvis-chrome.sh start-local` ativou DevTools em `127.0.0.1:9222` com perfil isolado
- `jarvis-chrome.sh start-headless` também ativou DevTools em loopback com perfil isolado
- o processo Chrome isolado não se manteve estável o suficiente para fluxo contínuo de `navigate` via `curl`
- conclusão operacional: DevTools helper serve como bootstrap complementar, mas o baseline de browser-use deve migrar para Playwright
- Playwright headless validou `https://example.com/` com título `Example Domain`
- Playwright gerou screenshot válido em `/tmp/jarvis-playwright-smoke.png`
