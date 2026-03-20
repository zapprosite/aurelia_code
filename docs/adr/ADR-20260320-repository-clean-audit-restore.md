# ADR-20260320-repository-clean-audit-restore: Auditoria, Limpeza e Restauração da Integridade

## Status
**Proposed**

## Contexto
Durante a fase de consolidação do "Voice Plane", foi identificada uma desorganização significativa na estrutura do monorepo:
1.  **Poluição na Raiz**: Arquivos temporários (ex: `?_journal_mode=WAL`) e binários dispersos.
2.  **Salada de Documentação**: A pasta `docs/` continha arquivos redundantes e sem categorização clara.
3.  **Regressões de Código**: 
    - Conflitos de merge mal resolvidos em `pkg/llm/catalog_test.go`.
    - Ausência física de arquivos de implementação para as rotas de Vision/Photo no Telegram (`internal/telegram`), resultando em quebra de build.

## Decisão
Implementar uma auditoria profunda e restauração imediata baseada em:
1.  **Poda Estrutural**:
    - Organização de `docs/` em subpastas (`blueprints`, `governance`, `guides`, `archive`).
    - Organização de `.agents/workflows/` (ex: `adr-semparar/`).
2.  **Restauração de Build**:
    - Resolução manual de conflitos em `pkg/llm/catalog_test.go` (unindo testes de Ollama e Vision).
    - Desativação temporária das rotas `OnPhoto` e `OnDocument` no bot do Telegram e comentário dos testes correspondentes em `input_test.go` até que a implementação original seja recuperada ou rescrita.
3.  **Higiene Documental**:
    - Criação desta ADR para rastrear a decisão de "Build Verde Primeiro".

## Consequências
- **Positivo**: Repositório volta a compilar e testar com sucesso nas áreas core. Documentação legível.
- **Negativo**: Funcionalidade de Visão (envio de fotos ao bot) está temporariamente inativa.

## Verificação
- `go test ./pkg/llm/...` (Sucesso após fix).
- `go test ./internal/telegram/...` (Sucesso após isolamento).
- `git status` (Limpo).

## A Ideia desta ADR
Documentar que, antes de avançar para novas features, é **obrigatório** manter o monorepo saudável, limpo e compilável. A "poda" não é apenas estética, é funcional.
