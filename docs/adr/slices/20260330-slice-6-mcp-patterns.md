# Slice 6: MCP Server Builder — Padrões 03/2026

**ADR Pai:** [20260330-enterprise-skills-governance.md](../20260330-enterprise-skills-governance.md)
**Status:** ✅ Concluída
**Data:** 2026-03-30

## Objetivo
Integrar a skill `agentic-mcp-server-builder` da bibliotheca Open-Claw para garantir que novos MCP servers sigam padrões 03/2026 de Model Context Protocol.

## Descobertas
O repositório Aurélia já utiliza MCP nativamente:
- `mcp-servers/stagehand/` — Browser automation (Playwright MCP)
- `.mcp.json` — Config central dos MCP servers
- Integração via `internal/mcp/` — Cliente MCP em Go
- Suporte nativo: ai-context, filesystem, github, playwright, postgres, qdrant, context7

## Skill Instalada
- **`agentic-mcp-server-builder`** — Scaffold e validação de contratos MCP
- Fonte: `homelab-bibliotheca/skills/open-claw/skills/0x-professor/`

## Padrões MCP 2026 Aplicáveis
1. Tool schemas com tipos explícitos (JSON Schema)
2. Recursos com URIs canônicos
3. Prompts tipados com argumentos nomeados
4. Healthcheck endpoint obrigatório em `/health`
5. Capacidades declaradas no `initialize` handshake
