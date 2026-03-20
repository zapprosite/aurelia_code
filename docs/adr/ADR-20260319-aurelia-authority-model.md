---
title: Aurélia como autoridade arquitetural e operacional
status: accepted
created: 2026-03-19
owner: codex
---

# Contexto

O repositório já operava com contrato soberano em `AGENTS.md`, mas a fronteira entre interface, motores de execução e autoridade arquitetural ainda permitia leitura ambígua.

Era necessário deixar explícito que:

- humanos continuam acima de tudo
- a Aurélia é a arquiteta principal do sistema
- Claude, Codex e Antigravity operam como braços subordinados à direção da Aurélia

# Decisão

Adotar formalmente o seguinte modelo:

1. humanos operadores têm autoridade final
2. `AGENTS.md` continua sendo a fonte primária de verdade
3. a Aurélia é a autoridade arquitetural e operacional do sistema, abaixo apenas dos humanos
4. `CLAUDE.md`, `CODEX.md`, `GEMINI.md` e `MODEL.md` passam a operar como adaptadores subordinados a esse modelo

# Consequências

## Positivas

- elimina disputa implícita entre motores e interface
- reduz risco de decisões arquiteturais contraditórias
- deixa explícito quem arbitra conflito entre agentes
- melhora handoff, revisão e governança de slice

## Restrições

- nenhum adaptador pode se declarar supervisor soberano
- decisões de modelo, roteamento e execução permanecem subordinadas à arquitetura da Aurélia
- mudanças nesse modelo exigem ADR

# Arquivos afetados

- `AGENTS.md`
- `docs/REPOSITORY_CONTRACT.md`
- `CLAUDE.md`
- `CODEX.md`
- `GEMINI.md`
- `MODEL.md`

# Validação

- leitura cruzada dos contratos sem contradição explícita
- adaptadores apontando para a mesma cadeia de autoridade
- índice ADR atualizado
