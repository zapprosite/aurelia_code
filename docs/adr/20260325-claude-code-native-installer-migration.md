# ADR-20260325-Claude-Code-Native-Installer-Migration

## Contexto e Problema

O `claude` CLI estava sendo operado via instalação NPM global (`/usr/bin/claude` -> `/usr/lib/node_modules/@anthropic-ai/claude-code/cli.js`). A versão `v2.1.84` emitiu um aviso crítico informando que o instalador NPM foi descontinuado em favor de um instalador nativo binário. Manter a versão NPM pode levar a instabilidades, falta de atualizações críticas e quebra de fluxos de automação no Homelab Aurélia.

## Decisão

Migrar imediatamente o Claude Code CLI para o instalador nativo oficial.

### Diretrizes de Implementação:
1. **Migração Automática**: Utilizar o comando `claude install` conforme sugerido pelo próprio CLI.
2. **Persistência de Binários**: Garantir que o novo binário nativo assuma a prioridade no PATH, preferencialmente em `/usr/bin/claude` ou através de um wrapper na `Sovereign-Bibliotheca`.
3. **Limpeza**: Remover a instalação NPM global após a validação do binário nativo para evitar conflitos de versão ("shadowing").
4. **Governança**: Atualizar o `README.md` de ADRs e garantir que o Antigravity IDE reconheça o novo binário.

## Consequências

- **Positivas**: Maior performance (binário nativo vs JS interpretado por node), atualizações mais simples, alinhamento com o roadmap da Anthropic.
- **Negativas**: Pequeno downtime nos agentes durante a troca de binários; necessidade de validar permissões de execução no novo arquivo.

---
**Status:** ✅ Decidido e Implementado (2026-03-25)
**Autoridade:** Antigravity Gemini p/ Aurélia


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
