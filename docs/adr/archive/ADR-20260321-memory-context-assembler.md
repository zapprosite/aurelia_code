> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

# ADR-20260321-memory-context-assembler: Sub-7 Memory Context Assembler

## Status
Ativo

## Contexto
O Loop Base da Aurelia (PREV) já entende o repositório hierarquicamente graças ao Codebase Symbol Map (Sub-4). No entanto, a inteligência autônoma necessita de histórico temporal e recordações de decisões passadas para não cometer os mesmos erros recursivamente ou reinventar arquiteturas. Parte das memórias está escrita em markdown, parte no SQLite e os embeddings habitam o banco Qdrant.

## Decisão
Implementar o `internal/memory/context_assembler.go`. 
Sua responsabilidade será atuar durante a fase `PLANNING` (via `loop.go`), extraindo ativamente do Qdrant e do SQLite as anotações, "thoughts" históricos ou resumos do repositório relacionados aos keywords do comando atual.
O contexto resgatado será injetado no System Prompt junto com o _Codebase Symbol Map_.

## Consequências
- A Aurelia ganhará continuidade de pensamento entre sessões ao lembrar o que fez ou o que o humano pediu no passado.
- O payload de contexto será densamente preenchido com _RAG_, otimizando de maneira agressiva o poder de Inferência Zero-Shot da LLM.

## Testes e Rollout
1. Implementar e testar `Assembler`.
2. Validar a concatenação do prompt em dry-run ou modo log debug.
3. Teste ponta a ponta chamando o bot no Telegram perguntando "O que decidimos ontem?"
