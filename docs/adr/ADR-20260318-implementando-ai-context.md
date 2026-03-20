# ADR 20260318-implementando-ai-context

**Status**: Proposto
**Data**: 2026-03-18
**Contexto**: O repositório Aurelia é um monorepo modular em Go que se beneficia de um contexto semântico denso para agentes de IA. Atualmente, faltam metadados estruturados para orquestração de agentes.
**Decisão**: Adotar o framework `ai-context` para gerenciar metadados de agentes, documentação e fluxos de trabalho (PREVC).
**Consequências**:
- Criação do diretório `.context/`.
- Padronização de documentação em `.context/docs/`.
- Uso de `workflows` para tarefas complexas.
