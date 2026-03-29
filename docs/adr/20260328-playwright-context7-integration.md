# ADR 20260328: Playwright CLI + Context7 MCP Integration

## Status
🟡 Proposto

## Contexto
Playwright 1.58.2 e Context7 MCP disponíveis. Integrar para:
1. **Computer Use** via Playwright MCP
2. **Context7** para documentação atualizada
3. **E2E Testing** do dashboard

## Decisões Arquiteturais

### 1. Playwright MCP Tools
```json
// mcp_servers.json já tem playwright
"playwright": {
  "command": "npx",
  "args": ["-y", "@playwright/mcp@latest"]
}
```

Playwright MCP oferece ~22 tools para browser automation.

### 2. Context7 Integration
```json
"context7": {
  "command": "npx",
  "args": ["-y", "@upstash/context7-mcp@latest"]
}
```

Context7 fornece documentação свежая de qualquer library.

### 3. Stagehand vs Playwright
| Aspecto | Stagehand | Playwright MCP |
|---------|-----------|----------------|
| AI-powered | ✅ | ❌ |
| Tool count | 3 | 22+ |
| Headless | ✅ | ✅ |
| Screenshots | ✅ | ✅ |

**Decisão**: Usar Stagehand como fallback quando AI disponível, Playwright para automation direta.

## Dependências
- ✅ mcp-servers/stagehand/
- ✅ configs/mcp_servers.json
- ⚠️ Playwright browser install

## Referências
- [Playwright MCP](https://github.com/microsoft/playwright-mcp)
- [Context7](https://context7.com)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/mcp-integrations
**Progress**: 0%
