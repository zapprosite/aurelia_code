# ADR-20260324: Instalar Tavily Web Search

**Status:** ✅ Decidido
**Data:** 24 de Março de 2026
**Autoridade:** Aurélia (Arquiteta Principal) / Antigravity (Coordenação)
**Contexto:** Substituição do DuckDuckGo pelo Tavily Web Search API para maior confiabilidade.

---

## 1. Contexto e Problema

A Aurélia usa atualmente **DuckDuckGo HTML scraping** como motor de `web_search`. Este método é frágil, sujeito a bloqueios por rate-limiting e não possui uma API oficial garantida no formato atual. Além disso, o **Claude Code CLI** não possui nenhum MCP de busca configurado no ambiente local de Will, limitando sua capacidade de pesquisa profunda.

A disponibilidade de uma API Key Tavily (`tvly-dev-...`) oferece uma oportunidade de modernização industrial, substituindo o scraping por chamadas de API oficiais, mais rápidas e com metadados mais ricos (título, URL e conteúdo limpo).

## 2. Decisão

Implementar o **Tavily Web Search** como provedor canônico de busca em dois destinos principais:

1.  **Claude Code CLI**: Configuração via MCP Server.
2.  **Aurélia Daemon (Go)**: Integração nativa na ferramenta `web_search`.

## 3. Abordagem Técnica

### Parte 1 — Claude Code CLI (MCP Server)

O arquivo `~/.claude/settings.json` será atualizado para incluir o servidor MCP do Tavily:

```json
{
  "skipDangerousModePermissionPrompt": true,
  "defaultMode": "bypassPermissions",
  "mcpServers": {
    "tavily": {
      "command": "npx",
      "args": ["-y", "tavily-mcp@0.1.4"],
      "env": {
        "TAVILY_API_KEY": "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"
      }
    }
  }
}
```

### Parte 2 — Aurélia Daemon (Go Native Tool)

1.  **`internal/tools/web_search.go`**: Reescrever a lógica para consumir `https://api.tavily.com/search`.
    - Payload: `{"api_key": "...", "query": "...", "max_results": N}`.
    - Resposta integrada ao fluxo atual.
2.  **`internal/tools/definitions.go`**: Atualizar a descrição da ferramenta para refletir o uso da API Tavily.
3.  **`internal/config/config.go`**: Adicionar o campo `TavilyAPIKey` mapeado para a env `TAVILY_API_KEY`.
4.  **`~/.aurelia/config/secrets.env`**: Persistir a chave de API para o runtime.

## 4. Consequências

- **Positivas**: 
    - Fim da fragilidade do scraping do DuckDuckGo.
    - Resultados de busca mais precisos e estruturados.
    - Claude Code agora possui capacidade de pesquisa na web via MCP.
- **Negativas**: 
    - Dependência de um serviço externo (Tavily) e consumo de créditos de API.
    - Necessidade de fallback para DuckDuckGo caso a chave expire ou atinja o limite. (Fallback será mantido como opcional).

## 5. Verificação

- [ ] `go build ./...` (Validar compilação).
- [ ] `go test ./internal/tools/... -run WebSearch -v` (Teste unitário).
- [ ] `/mcp` no Claude Code para validar o status do servidor Tavily.
- [ ] Smoke test do daemon com a nova env.


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
