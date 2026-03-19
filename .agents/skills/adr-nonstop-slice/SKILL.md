---
name: adr-nonstop-slice
description: Abre e mantém um slice em modo nonstop com ADR + JSON de continuidade estilo taskmaster.
---

# ADR Nonstop Slice

## Objetivo

Criar e manter slices que não morrem no meio do caminho, com:

- ADR do slice
- JSON de continuidade estilo taskmaster
- comandos de simulação/smoke já previstos
- handoff claro para Codex, Claude ou Gemini

## Quando usar

- quando uma slice estrutural precisa continuar entre sessões/agentes
- quando o trabalho depende de smoke com `curl`, `go test`, scripts ou simuladores
- quando o usuário quer "ativar" um modo que empurre a slice até o fim com menos perda de contexto

## Como executar

1. Criar o par de artefatos do slice com `./scripts/adr-slice-init.sh`.
2. Preencher o ADR `.md` com contexto, decisão, testes, rollout e rollback.
3. Preencher o `.json` com:
   - estado atual
   - próximos comandos
   - smoke commands
   - simulações
   - evidências
   - bloqueios
4. Durante a execução, atualizar o `.json` sempre que o próximo passo mudar.
5. Ao fechar a slice:
   - atualizar o ADR
   - atualizar o `.json`
   - rodar `./scripts/sync-ai-context.sh`

## Contrato

- o `.md` é a decisão humana legível
- o `.json` é a memória operacional para continuidade
- nenhum slice estrutural em modo nonstop existe sem os dois
- o nome canônico de novos slices é `ADR-YYYYMMDD-slug`

## Output esperado

- `docs/adr/ADR-YYYYMMDD-slug.md`
- `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`
- comandos de smoke/simulação preenchidos
- continuação clara para qualquer agente
