# ADR 20260318-integracao-mcp-antigravity

**Status**: Aceito
**Data**: 2026-03-18
**Contexto**: O Aurelia precisa de ferramentas externas para análise de código, navegação web e manipulação de arquivos. O Antigravity já possui um conjunto robusto de servidores MCP configurados.
**Decisão**: Migrar e converter a configuração `mcp_config.json` do Antigravity para o `~/.aurelia/config/mcp_servers.json`.
**Consequências**:
- O Aurelia ganha capacidades de: Busca semântica (Context7/Qdrant), Navegação (Playwright), Manipulação de Arquivos (Filesystem/GitHub) e Diagnóstico (AI-Context).
- Aumento da superfície de consumo de recursos (mais processos MCP em background).
