# Homelab Tutor v2 - 2026-03-19

## Goal

Criar uma base profissional para a Aurelia operar como tutor do homelab com:

- governança
- incident response
- runbooks executáveis
- memória operacional

## External Skills Installed

- `incident-response`
- `system-architect`
- `c-level-advisor`

Install target:

- `~/.agents/skills/`
- sincronizado para Antigravity, Codex, Claude Code e Gemini CLI

## Local Tutor Structure Created

Nova skill principal:

- `~/.aurelia/skills/homelab-tutor-v2/`

Arquivos criados:

- `SKILL.md`
- `INDEX.md`

## Runbooks Created

Monitoring:

- `gpu-metrics-recover.md`
- `prometheus-target-down.md`
- `grafana-no-data-triage.md`

Docker / Compose:

- `compose-service-missing.md`
- `n8n-health-recover.md`

Data / Platform:

- `qdrant-health-recover.md`
- `supabase-health-triage.md`

Tunnel / Network:

- `cloudflare-tunnel-recover.md`
- `firewall-drift-audit.md`

AI Runtime:

- `ollama-health-recover.md`
- `voice-stack-recover.md`

DR / Backup:

- `backup-age-enforcer.md`
- `zfs-scrub-review.md`

Additional infrastructure runbooks:

- `caprover-health-recover.md`
- `litellm-health-recover.md`
- `postgres-direct-check.md`
- `tailscale-access-check.md`
- `gpu-contention-triage.md`

## Domain Catalog Added

Tutor catalog created:

- `~/.aurelia/skills/homelab-tutor-v2/DOMAIN_RESPONSES.md`

This catalog standardizes how Aurelia should answer and route by domain:

- monitoring
- docker / compose
- data / platform
- AI runtime
- network / exposure
- DR / backup / storage

## Operating Model

O tutor agora deve operar assim:

1. classificar o domínio
2. escolher skill/runbook
3. coletar prova
4. aplicar correção mínima
5. validar
6. registrar prevenção

## Self-Healing Rule

Todo incidente relevante deve resultar em um dos seguintes:

- novo runbook
- atualização de runbook existente
- registro em `.context/workflow/docs/`
- refinamento de guardrail

## Audio Architecture Extension

Foi criado um blueprint específico para áudio PT-BR com Groq:

- `groq_ptbr_audio_blueprint.md`

Direção arquitetural registrada:

- `Groq` como camada de `speech-to-text`
- `Supabase` como fonte de verdade de sessões, mensagens e jobs
- `Qdrant` como memória semântica
- `LLM local` como cérebro e executor de instruções
- `TTS PT-BR` como camada separada de saída

## Audio Execution Slice

Primeiro corte executável do blueprint de áudio aplicado no código:

- `internal/telegram/input.go`
- `internal/telegram/input_audio_test.go`
- `scripts/groq-stt-curl-smoke.sh`
- `groq_stt_simulation.md`

O que mudou:

- transcrição de voz agora é persistida na memória local logo após o STT
- o arquivo arquivado usa `message_type=audio_transcript`
- o conteúdo salvo registra `provider`, `basename` do arquivo e texto transcrito
- foi adicionado smoke script de `curl` para validar a Groq com ou sem token

Provas registradas:

- `go test ./internal/telegram ./pkg/stt` passou
- smoke local sem token imprime o `curl` exato em `dry-run`
- chamada real com token falso retorna `401 invalid_api_key`, confirmando contrato HTTP

## Audio Cost/Quality Optimization

Slice adicional aplicado após ativação real da Groq:

- `pkg/stt/groq.go`
- `scripts/groq-stt-curl-smoke.sh`
- `groq_stt_simulation.md`

Decisões executadas:

- troca do modelo padrão para `whisper-large-v3-turbo`
- envio explícito de `language=pt`
- envio explícito de `temperature=0`

Motivo:

- reduzir custo por hora de transcrição
- reduzir latência
- melhorar previsibilidade para PT-BR

Provas:

- `GET /openai/v1/models` retornou `HTTP 200`
- `POST /openai/v1/audio/transcriptions` com WAV mínimo retornou JSON válido
- `go test ./pkg/stt` passou

## Antigravity Gemini Operator Slice

Foi criada a base para a Aurelia usar o chat do Antigravity como copiloto leve, sem misturar isso com execucao principal:

- `implementation_plan.md`
- `task.md`
- `docs/antigravity_gemini_operator_blueprint.md`
- `docs/PROJECT_PLAYBOOK.md`
- `.aurelia/skills/antigravity-gemini-operator/SKILL.md`
- `.aurelia/skills/antigravity-gemini-operator/CHAT_PROMPTS.md`
- `.aurelia/skills/antigravity-gemini-operator/DECISION_MATRIX.md`
- `.agents/skills/antigravity-gemini-operator/SKILL.md`
- `.agents/skills/antigravity-gemini-operator/CHAT_PROMPTS.md`
- `.agents/skills/antigravity-gemini-operator/DECISION_MATRIX.md`

Direcao registrada:

- o runtime do projeto passa a ter um playbook carregavel por memoria canonica
- a skill do projeto fica no path real carregado pela Aurelia
- a skill tambem fica espelhada em `.agents/skills/` para governanca versionada
- Gemini Flash / Minimax 2.7 no Antigravity vira executor de microtarefas
- Codex / CLI continua como executor principal, validador e agente de commit

## Telegram Light Delegation Slice

Foi ligado ao pipeline real do Telegram um atalho para tarefas `light`:

- `internal/telegram/antigravity_prompt.go`
- `internal/telegram/antigravity_prompt_test.go`
- `internal/telegram/input_pipeline.go`

Comportamento:

- a entrada do usuario e classificada por heuristica leve
- se a tarefa for `light`, a Aurelia gera automaticamente um prompt estruturado para o chat do Antigravity
- a resposta evita passar pelo loop completo quando a melhor acao e delegar pesquisa curta, config pequena ou revisao de diff
- tarefas com sinais de alto risco continuam fora desse atalho

Prova:

- `go test ./internal/telegram` passou

## JARVIS Local Voice Blueprint

Foi registrado um blueprint direto ao ponto para a Aurelia virar um assistente local de voz, browser-use, Antigravity e terminal:

- `docs/jarvis_local_voice_blueprint_20260319.md`
- `docs/local_model_kit_blueprint_20260319.md`

Baseado em:

- metricas reais de `Grafana/Prometheus`
- VRAM real da `RTX 4090`
- limites oficiais atuais da Groq
- stack local ja presente no host

Direcao registrada:

- `openWakeWord` + `Silero VAD` no CPU
- `Groq whisper-large-v3-turbo` para STT
- `gemma3:27b-it-q4_K_M` como cerebro local padrao
- `qwen3.5:27b-q4_K_M` como alternativa tecnica
- `qwen3-coder:30b` apenas para escalonamento manual
- `bge-m3` como contrato unico de embedding
- `Supabase + Qdrant` como memoria/persistencia
- `agent-browser` primeiro e `browser-use` como camada avancada
- rate limits conservadores e governor por recurso

Tambem foi registrado o papel do segredo local da Hugging Face:

- `~/.aurelia/config/secrets.env`
- `HF_TOKEN=...`

Uso:

- contingencia para downloads/autenticacao de artefatos HF
- fora do repo e fora do `app.json`

Execucao aplicada em `2026-03-19`:

- `scripts/update-ollama.sh` foi corrigido para o kit novo
- `scripts/ollama-local-kit-smoke.sh` foi criado para validar o modelo principal
- `ollama list` passou a mostrar:
  - `gemma3:27b-it-q4_K_M`
  - `gemma3:12b`
  - `bge-m3:latest`
- smoke real do `gemma3:27b-it-q4_K_M` com `ctx=8192` retornou `OK`
- o runtime passou a aceitar `llm_provider=ollama`
- o catalogo local passou a listar modelos do endpoint `v1/models`
- o onboarding passou a tratar `ollama` sem API key
- o `primary_llm` do `/health` passou a validar endpoint local + modelo instalado quando `provider=ollama`
- `go test ./... -count=1` passou depois desse slice

## Gemini Fallback Runtime

Foi consolidado o papel da Gemini API no runtime da Aurelia:

- `scripts/gemini-smoke.sh`
- `docs/gemini_fallback_runtime_20260319.md`

Decisao registrada:

- Gemini entra como fallback LLM e pesquisa curta
- `gemini-2.5-flash` e o default remoto
- `gemini-2.5-pro` fica para escalonamento manual
- Groq continua no STT
- `bge-m3` continua como embedding unico do Qdrant

Tambem foi endurecido o health:

- checks auxiliares agora podem emitir `warning`
- `warning` nao gera falso degradado
- `error` continua derrubando o `/health`

Observacao adicional validada em `2026-03-19`:

- a oferta estudantil do ecossistema Gemini nao deve ser confundida com aumento automatico da cota da Gemini API
- o beneficio de estudante agrega mais no `Gemini app`, `NotebookLM`, integracoes Google e armazenamento do plano `Google AI Pro`
- a Gemini Developer API continua dependente de `project`, `usage tier` e `Cloud Billing`
- conclusao arquitetural mantida: Gemini segue auxiliar no runtime da Aurelia, nao caminho critico

## LLM Gateway Blueprint

Foi registrado um blueprint novo para extrair do runtime atual um gateway/proxy inteligente de modelos:

- `docs/llm_gateway_blueprint_20260319.md`

Decisao arquitetural registrada:

- o repositório atual já tem bons adapters e uma base de transporte comum em `pkg/llm/`
- isso ainda nao é um gateway inteligente completo
- o caminho certo e evoluir para um gateway em `Go`, nao reescrever tudo em `Rust` agora
- `Rust` so faz sentido se isso virar um produto separado de alta vazao e multi-tenant

Partes que devem ser reaproveitadas:

- `pkg/llm/openai_compatible.go`
- adapters por provedor
- catalogo de modelos
- estrategia local de normalizacao/fallback

Partes novas que o blueprint exige:

- registry de capacidades
- policy engine de roteamento
- route scoring
- circuit breaker por `provider:model`
- normalizacao forte de `tool_calls`
- telemetria por decisao de rota

## Homelab Jarvis Operating Blueprint

Foi registrado um blueprint novo para consolidar a Aurelia como:

- Jarvis local
- bot de manutencao do homelab
- agente de codigo
- bibliotecaria de runbooks, docs e memoria
- operador de browser, Telegram e terminal

Arquivo principal:

- `docs/homelab_jarvis_operating_blueprint_20260319.md`

Decisao arquitetural registrada:

- LiteLLM ajuda como gateway de borda para clientes e UIs, nao como cerebro da Aurelia
- o plano de controle continua em `Go`
- o runtime local deve separar claramente:
  - control plane
  - inference plane
  - knowledge plane
  - execution plane

Papel dos bancos:

- `SQLite` como fonte de verdade do runtime local
- `Qdrant` como indice semantico derivado
- `Supabase` como estado compartilhado e integracao

Papel do audio/Jarvis:

- wake word local
- VAD local
- buffer local
- STT no Groq
- TTS separado
- LiteLLM opcional apenas como borda compativel, nao como coordenador do pipeline de voz

Loops de estabilidade enfatizados no blueprint:

- health loop
- repair loop
- memory loop
- hygiene loop
- documentation loop

## Revisao Final de VRAM e Modelo Local

Foi consolidada uma revisao final de budget de VRAM para o host com `RTX 4090`:

- uso base do desktop/lab: ~`4810 MiB`
- VRAM livre em idle: ~`19238 MiB`
- `qwen3.5:9b` carregado deixa ~`10.5 GiB` livres
- `qwen3.5:9b + qwen3.5:4b` juntos deixam so ~`3.8 GiB` livres
- um `27B` de ~`17 GiB`, somado ao uso base do host, deixa folga perto de `2 GiB`

Decisao arquitetural registrada:

- `Groq` segue como escolha correta para `STT` porque tira o audio do budget de VRAM local
- o runtime ativo da Aurelia/Jarvis deve operar com `1` modelo local residente por vez
- `qwen3.5:9b` vira o default local recomendado para o caminho ativo do bot
- `qwen3.5:4b` entra apenas frio ou aquecido sob demanda
- `gemma3:27b-it-q4_K_M` sai do caminho ativo e fica como modelo manual/offline de laboratorio

Documentos atualizados:

- `docs/homelab_jarvis_operating_blueprint_20260319.md`
- `docs/jarvis_local_voice_blueprint_20260319.md`
- `docs/local_model_kit_blueprint_20260319.md`
- `plan.md`

## Model Routing Matrix

Foi registrada uma matriz unica de roteamento para custo, qualidade e estabilidade:

- `docs/model_routing_matrix_20260319.md`

Decisao operacional:

- `qwen3.5:9b` como cerebro local do bot
- `qwen3.5:4b` apenas para triagem/fallback curto
- `Groq` isolado no lane de audio
- `OpenRouter` so por capacidade explicita
- `Gemini web` como fonte de pesquisa profunda e ingestao curada para o RAG
- `LiteLLM` como gateway de borda, nao como plano de controle

## Model Response Bakeoff

Foi registrado um bakeoff de qualidade/latencia entre modelos locais e remotos:

- `docs/model_response_bakeoff_20260319.md`

Achados principais:

- `minimax-m2.7` foi o mais polido em resposta operacional livre
- `deepseek-v3.2` foi o melhor em JSON curto e curadoria compacta para RAG
- `qwen3.5-flash` foi bom e disciplinado, mas mais lento nesta rodada
- `minimax-m2.7` e `qwen3.5:9b` mostraram consumo agressivo de `reasoning` em prompts estruturados curtos quando o budget estava baixo

Decisao implicita para o roteamento:

- `minimax-m2.7` fica como lane premium de workflow
- `deepseek-v3.2` ganha forca como lane remoto para estrutura curta e curadoria

## Gateway Dry-Run

Foi implementado o primeiro corte do gateway interno em `Go`:

- pacote `internal/gateway/`
- endpoint `POST /v1/router/dry-run`

Capacidades entregues:

- policy engine inicial por lane
- guardas iniciais de `reasoning_mode`, `max_output_tokens` e `soft_timeout_ms`
- decisao explicita de `provider/model/lane/budget`

Prova:

- `go test ./internal/gateway ./internal/health ./cmd/aurelia -count=1` passou

## Gateway Rollout Blueprint

Foi registrado o blueprint do restante do gateway:

- `docs/gateway_rollout_blueprint_20260319.md`

Escopo fechado:

- route enforcement no runtime
- guardas de reasoning e output
- budgets por lane
- circuit breaker por `provider:model`
- telemetria
- rollout na worktree de deploy

## Aurelia General Blueprint

Foi consolidado um blueprint geral unico para o restante do projeto:

- `docs/aurelia_general_blueprint_20260319.md`

Escopo consolidado:

- gateway real
- voz em background
- memoria operacional completa
- governor e budgets
- desktop fallback seguro
- rollout final em deploy

## Aurelia Master Blueprint

Foi consolidado um blueprint mestre com arquitetura, rollout e matriz de testes completa:

- `docs/aurelia_master_blueprint_20260319.md`

Escopo:

- arquitetura alvo
- ordem de execucao
- gates
- testes unitarios
- testes de integracao
- E2E
- deploy
- rollback

## Gateway Enforcement Slice

O gateway saiu do modo apenas documental e passou a influenciar o runtime real:

- `internal/gateway/provider.go`
- `cmd/aurelia/app.go`
- `cmd/aurelia/health_checks.go`
- `internal/health/server.go`

Capacidades entregues:

- selecao real de lane/modelo no runtime
- guardas de resposta por lane
- budgets por lane em memoria
- circuit breaker por `provider:model`
- endpoint `GET /v1/router/status`

Estabilizacao adicional:

- `pkg/llm/openai_compatible_test.go` foi alinhado ao novo contrato de request options
- `internal/agent/lead_worker_runtime_test.go` passou a usar SQLite em memoria no helper de runtime para eliminar o flaky de `TempDir` com WAL
- `internal/gateway/provider_test.go` cobre fallback, health e status do gateway

Provas:

- `go test ./internal/gateway ./pkg/llm ./cmd/aurelia -count=1` passou
- `go test ./internal/agent -count=5` passou
- `go test ./... -count=1` passou

Pendencias honestas:

- telemetria Prometheus do gateway
- rollout e validacao na worktree de deploy

## Voice Processor + Gateway Metrics

O repositório principal recebeu o primeiro fechamento operacional do voice plane e da observabilidade do gateway.

Arquivos centrais:

- `internal/gateway/metrics.go`
- `internal/gateway/provider.go`
- `internal/voice/processor.go`
- `internal/voice/spool.go`
- `internal/voice/mirror.go`
- `internal/voice/metrics.go`
- `cmd/aurelia/app.go`
- `cmd/aurelia/voice_cli.go`
- `internal/telegram/input_pipeline.go`

Capacidades entregues:

- métricas Prometheus do gateway em `/metrics`
- métricas Prometheus do loop de voz em `/metrics`
- spool local de jobs de áudio
- processador de fila com heartbeat e budget diário
- fallback STT por comando configurável
- gate textual por wake phrase no transcript
- dispatch do transcript aceito para o mesmo fluxo do Telegram
- mirrors opcionais de transcript para Supabase e Qdrant
- `GET /v1/voice/status`
- CLI `aurelia voice enqueue <arquivo>`

Provas:

- `go test ./cmd/aurelia ./internal/voice -count=1` passou
- `go test ./... -count=1` passou

Limites honestos:

- o runtime ainda não tem captura contínua de microfone
- `openWakeWord + Silero VAD + ring buffer` continuam como próximo slice real
- o rollout disso tudo na worktree `/home/will/aurelia-24x7` ainda não foi feito

## Governança: sync-ai-context como Regra de Slice

A política de sincronização de contexto foi promovida de convenção para regra formal do repositório.

Arquivos ligados à decisão:

- `docs/adr/20260319-sync-ai-context-como-regra-de-slice.md`
- `AGENTS.md`
- `.agents/rules/05-context-state.md`
- `.agents/skills/sync-ai-context/SKILL.md`
- `docs/REPOSITORY_CONTRACT.md`

Decisão:

- `sync-ai-context` passa a ser obrigatório em slice não trivial, handoff e preparação para merge
- `sync-ai-context` pode ser dispensado em microedições triviais sem drift semântico relevante

Comando canônico mantido:

- `./scripts/sync-ai-context.sh`

## ADR do Último Plano de Voz

O último plano de voz deixou de ficar apenas em blueprint/backlog e agora tem ADR próprio:

- `docs/adr/20260319-voice-capture-plane.md`

Escopo formalizado:

- `openWakeWord`
- `Silero VAD`
- `ring buffer`
- worker de captura contínua
- integração com o spool já existente

Estado honesto:

- ADR criada
- implementação ainda pendente

## Skill ADR Nonstop Slice

Foi criado um padrão operacional para slices longas e multiagente:

- skill: `.agents/skills/adr-nonstop-slice/`
- scaffold: `./scripts/adr-slice-init.sh`
- template ADR: `docs/adr/TEMPLATE-NONSTOP-SLICE.md`
- template JSON: `docs/adr/taskmaster/TEMPLATE-NONSTOP-SLICE.json`

Objetivo:

- abrir slice com ADR + JSON de continuidade
- manter smoke/simulações/curl já previstos
- permitir handoff entre Codex, Claude e Gemini sem perda de contexto

Prova:

- `bash -n ./scripts/adr-slice-init.sh` passou
- `./scripts/adr-slice-init.sh voice-capture-rollout --title "Voice Capture Rollout" --dry-run` gerou paths válidos

## Higiene da Raiz e Migração para ADR

A raiz do repositório foi consolidada para manter apenas contratos, docs de entrada, o plano mestre e exemplos globais.

Movimentos feitos:

- `groq_ptbr_audio_blueprint.md` -> `docs/groq_ptbr_audio_blueprint_20260319.md`
- `groq_stt_simulation.md` -> `docs/groq_stt_simulation_20260319.md`
- `homelab_tutor_v2_blueprint.md` -> `docs/homelab_tutor_v2_blueprint_20260319.md`
- `plan_para_aurelia.md` -> `docs/keepassxc_local_vault_guide_20260319.md`
- `implementation_plan.md` -> `.context/plans/20260319-antigravity-gemini-operator/implementation_plan.md`
- `task.md` -> `.context/plans/20260319-antigravity-gemini-operator/task.md`

ADRs promovidas:

- `docs/adr/20260319-root-document-hygiene.md`
- `docs/adr/20260319-antigravity-copiloto-leve.md`
- `docs/adr/20260319-groq-stt-ptbr-runtime.md`
- `docs/adr/20260319-homelab-tutor-v2.md`
- `docs/adr/20260319-keepassxc-cofre-humano.md`

Regras alinhadas:

- `README.md`
- `docs/REPOSITORY_CONTRACT.md`
- `.agents/workflows/planejar.md`
- `.agents/workflows/handoff-to-claude-implementer.md`
- `.agents/skills/architect-planner/SKILL.md`
- `.agents/rules/06-planning-first.md`

Resultado:

- a raiz ficou limpa para documentos soberanos e de entrada
- decisões permanentes foram promovidas para ADR
- artefatos de slice passaram a morar em `.context/plans/`
