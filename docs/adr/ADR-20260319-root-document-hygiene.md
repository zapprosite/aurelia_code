# ADR 20260319-root-document-hygiene

**Status**: Aceito  
**Data**: 2026-03-19

## Contexto

A raiz do repositório acumulou blueprints, guias, smoke docs e artefatos de slice que não são contratos soberanos nem pontos de entrada do projeto. Isso reduz legibilidade, dificulta onboarding e espalha a governança.

## Decisão

A raiz fica reservada a:

- `AGENTS.md`
- `CLAUDE.md`
- `CODEX.md`
- `GEMINI.md`
- `MODEL.md`
- `README.md`
- `CONTRIBUTING.md`
- `SECURITY.md`
- `plan.md`
- exemplos globais como `mcp_servers.example.json`

Artefatos fora desse conjunto devem ser movidos para:

- `docs/adr/` quando representarem decisão arquitetural
- `docs/` quando forem blueprints, guias ou runbooks canônicos
- `.context/plans/<slice>/` quando forem `implementation_plan.md` e `task.md` de slice

## Consequências

Positivas:

- raiz mais limpa e previsível
- menor ambiguidade sobre o que é contrato versus histórico de slice
- onboarding mais rápido para humano e agente

Trade-offs:

- exige disciplina de atualização de links
- requer manutenção explícita do índice ADR e dos documentos canônicos

## Referências

- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)
