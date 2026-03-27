# ADR 2026-03-27: Markdown Brain canônico para o `aurelia_code`

## Status
Implementada

## Objetivo

Transformar Markdown em memória operacional de primeira classe para o bot soberano `aurelia_code`, sem depender do app Obsidian instalado no host.

## Contexto

O runtime já tinha três trilhas separadas:

- `conversation_memory` no Qdrant para histórico conversacional
- sync opcional de Obsidian em pipeline lateral
- documentação `.md` do repositório fora do retrieval do agente

Esse desenho deixava o agente cego para a própria base documental do repo e criava drift entre “memória de conversa”, “vault” e “docs do código”.

## Decisão

Adotar uma coleção vetorial canônica separada, `aurelia_markdown_brain`, com estas regras:

1. indexar todo `.md` estratégico do repositório
2. indexar também o vault externo quando `obsidian_vault_path` estiver configurado
3. chunkar por seções com metadados de `section`, `repo_path` ou `vault_path`, `source_id`, `chunk_id` e `checksum`
4. expor sync por tool única: `markdown_brain_sync`
5. recuperar essa coleção no `ContextAssembler` para `aurelia` e `aurelia_code`
6. manter `conversation_memory` separado para não misturar fatos conversacionais com documentação estática

## Fora de escopo

- edição bidirecional de vault
- sincronização do app Obsidian via UI
- substituir a memória conversacional por Markdown

## Consequências

- `aurelia_code` passa a enxergar o cérebro em `.md` do repo e do vault no mesmo contrato
- dashboard `/api/brain/search` e `/api/brain/recent` passam a consultar múltiplas collections
- cron e boot sync ficam alinhados com a mesma tool e a mesma collection
- o runtime elimina o drift estrutural de “docs de um lado, cérebro de outro”

## Smoke obrigatório

```bash
go test ./internal/markdownbrain ./internal/memory ./cmd/aurelia
```

## Critério de saída

- collection `aurelia_markdown_brain` criada e sincronizada
- `aurelia_code` recupera contexto de `.md` no prompt de execução
- repo markdown e vault markdown convergem no mesmo pipeline
