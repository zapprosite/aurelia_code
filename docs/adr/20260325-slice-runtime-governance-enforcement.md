# ADR 2026-03-25: Slice 2 — Runtime Governance Enforcement

## Status
Implementada

## Objetivo

Transformar governança de dados e runtime em enforcement real de código, não só documentação.

## Contexto

O repositório já tem contrato bom para `Supabase`, `SQLite`, `Qdrant`, `Obsidian` e status operacional. O gap restante é clássico: os docs mandam, mas o runtime ainda aceita drift em pontos importantes.

Se isso ficar assim, futuros LLMs voltam a inventar schema, payload e wiring.

## Escopo

- validar payloads canônicos antes de escrita em Qdrant
- bloquear escrita sem `canonical_bot_id`, `source_system`, `source_id` e `version` quando aplicável
- explicitar no runtime o provider/model efetivo por bot e por app
- adicionar testes de contrato para governança soberana
- endurecer rotas de observabilidade para não mascararem desvio de configuração

## Fora de escopo

- migração completa para Supabase como fonte de verdade em todo módulo
- refactor grande de MCPs
- mudança do modelo stack vigente

## Mudanças esperadas

1. criar validadores de payload para memória/indexação
2. aplicar validadores nos writers/mirrors relevantes
3. adicionar testes de contrato no runtime
4. reforçar exposição de provider/model efetivo nos endpoints
5. falhar cedo quando config e comportamento divergirem de forma ilegítima

## Smoke obrigatório

```bash
go test ./cmd/... ./internal/... ./pkg/... ./e2e
curl -sS http://127.0.0.1:3334/api/runtime/llm
curl -sS http://127.0.0.1:3334/api/status
```

## Critério de saída

- payload sem contrato canônico falha cedo
- provider efetivo fica explícito e auditável
- docs soberanos continuam validados por teste

## Dependência

Pode começar depois da Slice 1. Não depende da Slice 3.


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
