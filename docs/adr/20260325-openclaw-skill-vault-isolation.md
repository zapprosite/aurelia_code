# ADR 2026-03-25: Isolar OpenClaw Skill Vault do Módulo Raiz

## Status
Aceita

## Contexto

A árvore `homelab-bibliotheca/skills/open-claw/skills/` contém milhares de artefatos de terceiros, incluindo código Go inválido, dependências ausentes e subprojetos que não fazem parte do runtime soberano da Aurélia.

Isso fazia `go test ./...` na raiz falhar por lixo externo não governado.

## Decisão

Isolar `homelab-bibliotheca/skills/open-claw/skills/` como módulo Go separado via `go.mod` próprio.

## Consequências

- `go test ./...` na raiz deixa de varrer esse acervo
- o runtime principal continua validável sem filtros ad hoc
- o acervo permanece preservado, sem deleção destrutiva
- qualquer execução Go dentro desse vault passa a ser decisão explícita e local



---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
