# ADR 20260401: Industrialização OpenCode & Sincronização SOTA

## Contexto e Problema
O repositório `aurelia` utilizava múltiplas fontes de configuração para agentes (Antigravity, Claude Code, Codex), resultando em:
1. Fragmentação de segredos e tokens.
2. Dessincronização de servidores MCP entre diferentes ferramentas.
3. Dificuldade de onboarding de novos frameworks agentic.

## Decisão
Implementar o padrão **OpenCode (SOTA 2026)** para orquestração agentic no monorepo:
1. **Estrutura Centralizada**: Criada a pasta `.opencode/` na raiz do projeto como repositório único de skills, workflows e configurações.
2. **SSOT MCP**: O arquivo `.opencode/mcp_servers.json` torna-se a Fonte Única de Verdade (Single Source of Truth) para todos os servidores MCP.
3. **Refatoração Multi-Agente**: Scripts de tradução convertem o SSOT para formatos nativos (TOML para Codex, JSON para Claude/Antigravity) nas pastas de refatoração sincronizadas.
4. **Governança de Segredos**: Centralização total de tokens em `.opencode/.env`, eliminando dispersão em arquivos ocultos do SO.

## Consequências
- **Consistência Total**: Mudanças em um servidor MCP são refletidas automaticamente em todos os agentes após o sync.
- **Portabilidade**: O ecossistema de agentes pode ser movido entre máquinas apenas clonando o repositório e restaurando o `.env`.
- **Rastro de Auditoria**: Mudanças na infraestrutura de agentes agora são versionadas e documentadas via ADR.

## Status
Aprovado / Implementado (2026-04-01)
