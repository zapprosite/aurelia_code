# Protocolo MCP — Model Context Protocol

## Visão Geral

MCP é o protocolo para comunicação entre agentes e ferramentas. Desenvolvido pela Anthropic, permite que agentes acessem tools de forma padronizada.

## Conceitos Fundamentais

### Resource
Dados que o agente pode acessar:

```json
{
  "uri": "file:///project/main.go",
  "name": "main.go",
  "mimeType": "text/plain",
  "size": 15234
}
```

### Tool
Função que o agente pode chamar:

```json
{
  "name": "docker_ps",
  "description": "Lista containers Docker em execução",
  "inputSchema": {
    "type": "object",
    "properties": {
      "all": {
        "type": "boolean",
        "description": "Mostrar todos os containers"
      }
    }
  }
}
```

### Prompt
Template de prompt reutilizável:

```json
{
  "name": "code_review",
  "description": "Template para code review",
  "arguments": [
    {"name": "file", "type": "string"},
    {"name": "focus", "type": "string"}
  ]
}
```

## JSON-RPC 2.0 Methods

### initialize
Inicia conexão com servidor MCP:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "clientInfo": {
      "name": "aurelia_code",
      "version": "1.0.0"
    },
    "capabilities": {}
  }
}
```

### tools/list
Lista tools disponíveis:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}
```

### tools/call
Chama uma tool:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "docker_ps",
    "arguments": {"all": true}
  }
}
```

### resources/list
Lista resources:

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "resources/list",
  "params": {}
}
```

### resources/read
Lê um resource:

```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "resources/read",
  "params": {
    "uri": "file:///project/main.go"
  }
}
```

## Ferramentas DevOps via MCP

### Docker MCP

```json
{
  "name": "docker_ps",
  "description": "Lista containers",
  "inputSchema": {"type": "object", "properties": {"all": {"type": "boolean"}}}
}
{
  "name": "docker_logs",
  "description": "Busca logs de container",
  "inputSchema": {"type": "object", "properties": {"container": {"type": "string"}, "tail": {"type": "number"}}}
}
{
  "name": "docker_exec",
  "description": "Executa comando em container",
  "inputSchema": {"type": "object", "properties": {"container": {"type": "string"}, "cmd": {"type": "string"}}}
}
```

### Kubernetes MCP

```json
{
  "name": "k8s_get_pods",
  "description": "Lista pods",
  "inputSchema": {"type": "object", "properties": {"namespace": {"type": "string"}}}
}
{
  "name": "k8s_get_deployments",
  "description": "Lista deployments",
  "inputSchema": {"type": "object", "properties": {"namespace": {"type": "string"}}}
}
{
  "name": "k8s_logs",
  "description": "Busca logs de pod",
  "inputSchema": {"type": "object", "properties": {"pod": {"type": "string"}, "namespace": {"type": "string"}}}
}
```

### GitHub MCP

```json
{
  "name": "gh_list_prs",
  "description": "Lista PRs",
  "inputSchema": {"type": "object", "properties": {"repo": {"type": "string"}, "state": {"type": "string"}}}
}
{
  "name": "gh_get_pr",
  "description": "Detalhes de PR",
  "inputSchema": {"type": "object", "properties": {"repo": {"type": "string"}, "pr": {"type": "number"}}}
}
{
  "name": "gh_create_pr",
  "description": "Cria PR",
  "inputSchema": {"type": "object", "properties": {"repo": {"type": "string"}, "title": {"type": "string"}, "body": {"type": "string"}}}
}
```

## Tavily (Search)

```json
{
  "name": "tavily_search",
  "description": "Pesquisa web via Tavily",
  "inputSchema": {
    "type": "object",
    "properties": {
      "query": {"type": "string", "description": "Query de busca"},
      "max_results": {"type": "number", "description": "Máximo de resultados"}
    },
    "required": ["query"]
  }
}
```

### Resposta Tavily

```json
{
  "results": [
    {
      "title": "Artigo",
      "url": "https://...",
      "content": "...",
      "score": 0.95
    }
  ],
  "answer": "Resumo do conteúdo..."
}
```

## Context7 (Docs)

```json
{
  "name": "context7_docs",
  "description": "Busca docs atualizadas",
  "inputSchema": {
    "type": "object",
    "properties": {
      "library": {"type": "string", "description": "Biblioteca (react, go, etc)"},
      "query": {"type": "string", "description": "O que buscar"}
    },
    "required": ["library", "query"]
  }
}
```

## Implementação Go

### MCP Client

```go
type MCPClient struct {
    endpoint string
    client   *http.Client
}

func (m *MCPClient) CallTool(ctx context.Context, name string, args map[string]interface{}) (interface{}, error) {
    req := JSONRPCRequest{
        JSONRPC: "2.0",
        ID:      uuid.New().String(),
        Method:  "tools/call",
        Params: CallToolParams{
            Name:      name,
            Arguments: args,
        },
    }
    
    resp, err := m.client.Post(m.endpoint, "application/json", req.Body)
    // parse response
}
```

## Integração com A2A

```
┌─────────────────┐      MCP       ┌─────────────────┐
│  A2A Agent      │◄──────────────►│  MCP Server     │
│  (aurelia_code) │                │  (Docker/K8s/etc)│
└─────────────────┘                └─────────────────┘
        │
        └──► tools/call ──> docker_ps
              │
              └──► docker ps
              │
              ◄── result
```

---

*Referência: modelcontextprotocol.io*