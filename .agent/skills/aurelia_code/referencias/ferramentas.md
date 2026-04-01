# Ferramentas — Context7, Tavily, n8n

## Context7 MCP

### O que é
Plataforma que injeta documentação atualizada diretamente no contexto do agente AI. Resolve o problema de APIs desatualizadas nos modelos.

### Uso

```bash
# Via MCP
{
  "name": "context7_docs",
  "arguments": {
    "library": "react",
    "query": "useState hook"
  }
}
```

### Bibliotecas Suportadas
- React, Vue, Angular
- Go, Python, Rust, Node.js
- Docker, Kubernetes
- E mais 100+ libs

### Instalação

```bash
# Via npm
npx ctx7 init

# Configurar no projeto
# Adicionar ao MCP servers config
```

---

## Tavily API

### O que é
API de busca web diseñada para AI agents. Fornece resultados estruturados e respostas sintetizadas.

### Uso

```bash
# Via API REST
POST https://api.tavily.com/search
{
  "query": "Go agent framework 2026",
  "max_results": 5,
  "include_answer": true
}
```

### Response

```json
{
  "answer": "Os melhores frameworks Go para agents em 2026 são...",
  "results": [
    {
      "title": "go-agent - Protocol Lattice",
      "url": "https://github.com/...",
      "content": "...",
      "score": 0.95
    }
  ]
}
```

### Integração na Aurélia

```go
type TavilyClient struct {
    APIKey string
    BaseURL string
}

func (t *TavilyClient) Search(ctx context.Context, query string) (*SearchResult, error) {
    req, _ := http.NewRequest("POST", t.BaseURL+"/search", ...)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+t.APIKey)
    
    // ...
}
```

---

## n8n — Automação

### O que é
Plataforma de automação de workflows (self-hosted ou cloud). Permite conectar serviços, APIs, webhooks.

### Webhook na Aurélia

```bash
# Trigger workflow
curl -X POST "http://localhost:5678/webhook/my-workflow" \
  -H "Content-Type: application/json" \
  -d '{
    "trigger": "aurelia_code",
    "mission": "deploy",
    "data": {...}
  }'
```

### Workflows Úteis

| Trigger | Ação |
|---------|------|
| Missão concluída | Enviar notificação (Telegram, Slack) |
| Code review aprovado | Deploy automático |
| Task falhou | Alertar time |
| Daily standup | Resumo de missões ativas |

### Integração

```go
type N8NClient struct {
    WebhookURL string
    httpClient *http.Client
}

func (n *N8NClient) Trigger(ctx context.Context, payload map[string]interface{}) error {
    jsonData, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", n.WebhookURL, bytes.NewReader(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := n.httpClient.Do(req)
    // ...
}
```

---

## Stack Completa

```
┌─────────────────────────────────────────────────────────────┐
│                     AURELIA_CODE                           │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Context7   │  │   Tavily    │  │     n8n     │        │
│  │  (docs)     │  │  (search)   │  │ (workflow)  │        │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘        │
│         │                │                │                 │
│         ▼                ▼                ▼                 │
│    ┌─────────────────────────────────────────────┐        │
│    │              MCP Client                      │        │
│    │   (unified tool access layer)               │        │
│    └─────────────────────────────────────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

---

## Configuração

```bash
# .env

# Context7
CONTEXT7_ENABLED=true

# Tavily
TAVILY_API_KEY=tvly-...

# n8n
N8N_WEBHOOK_URL=http://localhost:5678/webhook/aurelia
N8N_API_KEY=...
```

---

*Referências oficiais: context7.com, tavily.com, docs.n8n.io*