---
description: Contrato canônico do catálogo de skills, overlays, indexação semântica e auditoria.
status: active
owner: Will
created: 2026-03-27
---

# Skill Catalog

Este documento governa onde skills vivem, como são carregadas, como são indexadas e como o drift estrutural deve ser auditado.

## Objetivo

Eliminar catálogos paralelos, referências quebradas e duplicação de governança entre motores.

## Fonte canônica

- Catálogo versionado do repo: [`.agent/skills/`](../../.agent/skills)
- Índice humano do catálogo: [`.agent/skills/README.md`](../../.agent/skills/README.md)
- Contrato de entrada: [`AGENTS.md`](../../AGENTS.md)

## Overlays permitidos

- Global opcional: `~/.aurelia/skills`
- Overlay local opcional do projeto: `<repo>/.aurelia/skills`

Uso dos overlays:

- extensão temporária
- skill instalada por automação
- experimentação local fora do catálogo versionado

Regra de precedência:

1. `~/.aurelia/skills`
2. `<repo>/.aurelia/skills`
3. `<repo>/.agent/skills`

Como a última fonte vence em colisão, o catálogo canônico do repo sempre sobrescreve overlays com o mesmo nome lógico.

## Contrato de arquivo

Cada skill versionada deve conter:

- diretório próprio
- `SKILL.md`
- frontmatter YAML com `name` e `description`

Campos opcionais aceitos no frontmatter:

- `tags`
- `engines`
- `owner`
- `phases`

## Runtime

O daemon Go deve:

- carregar skills do catálogo canônico e dos overlays
- preferir o catálogo do repo em colisões
- auditar o catálogo na inicialização
- sincronizar embeddings da coleção `aurelia_skills` no Qdrant

## Qdrant e embeddings

Collection canônica:

- `aurelia_skills`

Modelo de embedding padrão:

- `nomic-embed-text`

Política:

- indexar por chunks derivados das seções do markdown
- preservar contexto de `name` e `description` em cada chunk
- incluir `section`, `path`, `chunk_id`, `chunk_index`, `chunk_count` e `checksum` no payload
- deduplicar skill por `name` na fase de busca

## Auditoria obrigatória

A auditoria do catálogo deve detectar:

- frontmatter inválido
- descrição vazia
- links markdown quebrados dentro de `SKILL.md`
- nomes duplicados entre diretórios-base
- referência legado `.agents` nos entrypoints de governança

## Adapters

Os adapters de motor devem ser finos e apontar para o mesmo contrato:

- [`CLAUDE.md`](../../CLAUDE.md)
- [`GEMINI.md`](../../GEMINI.md)
- [`MODEL.md`](../../MODEL.md)
- [`.github/copilot-instructions.md`](../../.github/copilot-instructions.md)

Eles não devem manter catálogo paralelo nem copiar governança longa.
