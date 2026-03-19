# ADR 20260319-voice-capture-plane

**Status**: Proposto  
**Data**: 2026-03-19

## Contexto

O repositório já possui a primeira fase executável do voice plane:

- spool local de áudio
- processador de fila com heartbeat
- budget diário de STT
- fallback STT por comando
- gate textual por wake phrase no transcript
- dispatch do texto aceito para o mesmo fluxo do Telegram
- mirrors opcionais para Supabase e Qdrant
- endpoints `/v1/voice/status` e `/metrics`

Esse estágio prova a integração do pipeline, mas ainda não entrega a experiência JARVIS local completa, porque o runtime não captura áudio continuamente do microfone nem faz detecção real de wake word e VAD antes do upload para o STT.

## Decisão

O próximo slice de voz será fechado como **voice capture plane real**, com estes componentes:

1. `openWakeWord` para detecção local de wake phrase
2. `Silero VAD` para corte de silêncio e ruído
3. `ring buffer` local para preservar contexto curto antes/depois do wake
4. worker dedicado de captura de microfone
5. integração com o spool já existente
6. validação posterior de `aurelia-voice.service` ou worker equivalente no deploy

## Escopo

Entra neste slice:

- captura contínua de microfone
- wake word real
- VAD real
- criação automática de jobs no spool
- smoke local do voice path sem `voice enqueue` manual

Não entra neste slice:

- TTS premium
- browser/desktop fallback
- rollout completo na worktree de deploy
- E2E de Antigravity

## Arquivos / áreas afetadas

- `internal/voice/`
- `cmd/aurelia/`
- `internal/config/`
- `scripts/`
- `docs/`

## Guardrails

- não quebrar o runtime texto/Telegram já verde
- não disparar STT sem wake word válido
- respeitar budgets diários do Groq
- manter fallback local de STT disponível
- nenhuma nova dependência entra sem impacto explícito em CPU/RAM/operabilidade

## Testes obrigatórios

- unit:
  - wake word gating
  - VAD rejeita silêncio
  - clip válido gera job no spool
- integração:
  - worker de captura gera job real
  - spool -> STT -> dispatch continua funcionando
- smoke:
  - captura local curta com wake word e resposta
- deploy posterior:
  - heartbeat do worker
  - `/v1/voice/status` coerente

## Rollout

1. implementar captura local no repositório principal
2. manter `voice enqueue` como fallback operacional
3. validar suite e smoke local
4. só então portar para `/home/will/aurelia-24x7`

## Rollback

- desabilitar captura contínua
- manter spool/processador/fallback STT já existente
- preservar `voice enqueue` como caminho manual de segurança

## Consequências

Positivas:

- aproxima a Aurelia do modo JARVIS real
- reduz dependência de gatilho manual para voz
- mantém a governança já existente do spool e budgets

Trade-offs:

- aumenta complexidade operacional local
- adiciona sensibilidade a ruído, microfone e tunagem de thresholds

## Referências

- [AGENTS.md](../../AGENTS.md)
- [plan.md](../../plan.md)
- [jarvis_local_voice_blueprint_20260319.md](../jarvis_local_voice_blueprint_20260319.md)
- [aurelia_master_blueprint_20260319.md](../aurelia_master_blueprint_20260319.md)
- [PENDING-SLICES-20260319.md](./PENDING-SLICES-20260319.md)
