# Protocolo A2A — Agent-to-Agent

## Visão Geral

A2A (Agent-to-Agent) é o protocolo para comunicação entre agentes. Desenvolvido pelo Google, permite que agentes conversem entre si de forma padronizada.

## Conceitos Fundamentais

### Agent Card
Documento JSON que define as capacidades de um agente:

```json
{
  "name": "pesquisador",
  "description": "Agente de pesquisa web via Tavily e Context7",
  "url": "http://localhost:8081/agent",
  "version": "1.0.0",
  "capabilities": {
    "streaming": true,
    "pushNotifications": true
  },
  "skills": ["web_search", "doc_lookup", "analysis"],
  "authentication": {
    "type": "bearer",
    "credential": "..."
  }
}
```

### Task
Unidade de trabalho entre agentes:

```json
{
  "id": "task-uuid",
  "sessionId": "mission-uuid",
  "status": {
    "state": "submitted|working|completed|failed",
    "message": "Status message"
  },
  "artifacts": [...],
  "messages": [...]
}
```

### Message
Mensagem no formato role/content:

```json
{
  "role": "user|assistant|system",
  "parts": [
    {
      "type": "text",
      "text": "..."
    }
  ]
}
```

## JSON-RPC 2.0 Methods

### tasks/send
Envia uma tarefa para outro agente:

```json
{
  "jsonrpc": "2.0",
  "id": "req-1",
  "method": "tasks/send",
  "params": {
    "task": {
      "id": "task-001",
      "sessionId": "mission-001",
      "message": {
        "role": "user",
        "parts": [{"type": "text", "text": "Pesquise sobre X"}]
      }
    },
    "agentId": "pesquisador"
  }
}
```

### tasks/get
Retorna status de uma tarefa:

```json
{
  "jsonrpc": "2.0",
  "id": "req-2",
  "method": "tasks/get",
  "params": {
    "taskId": "task-001"
  }
}
```

### tasks/cancel
Cancela uma tarefa em andamento:

```json
{
  "jsonrpc": "2.0",
  "id": "req-3",
  "method": "tasks/cancel",
  "params": {
    "taskId": "task-001"
  }
}
```

### tasks/subscribe
Escuta updates de uma tarefa (streaming):

```json
{
  "jsonrpc": "2.0",
  "id": "req-4",
  "method": "tasks/subscribe",
  "params": {
    "taskId": "task-001"
  }
}
```

## Fluxo de Delegação

```
aurelia_code                    sub-agent
    │                               │
    │── tasks/send ───────────────►│
    │     (task + agentId)          │
    │                               │
    │◄── task response (working) ────│
    │                               │
    │    [processando...]          │
    │                               │
    │◄── task response (completed)──│
    │     (artifacts + messages)    │
    │                               │
    │── Qdrant.write() ───────────►│
    │     (resultado + memória)    │
```

## Implementação Go

### Server A2A

```go
type A2AServer struct {
    port    int
    agents  map[string]*Agent
    router  *Router
}

func (s *A2AServer) HandleTasksSend(w http.ResponseWriter, r *http.Request) {
    var req JSONRPCRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    task := req.Params.Task
    agentID := req.Params.AgentID
    
    agent := s.agents[agentID]
    result := agent.Execute(task)
    
    response := JSONRPCResponse{
        ID:     req.ID,
        Result: TaskResult{Task: result},
    }
    json.NewEncoder(w).Encode(response)
}
```

### Client A2A (Delegação)

```go
type A2AClient struct {
    endpoint string
    httpClient *http.Client
}

func (c *A2AClient) SendTask(ctx context.Context, agentID, prompt string) (*TaskResult, error) {
    task := Task{
        ID:        uuid.New().String(),
        Message:   Message{Role: "user", Parts: []Part{{Type: "text", Text: prompt}}},
    }
    
    req := JSONRPCRequest{
        JSONRPC: "2.0",
        ID:      uuid.New().String(),
        Method:  "tasks/send",
        Params:  SendTaskParams{Task: task, AgentID: agentID},
    }
    
    resp, err := c.httpClient.Post(c.endpoint+"/rpc", "application/json", req.Body)
    // ...
}
```

## Error Handling

| Código | Significado |
|--------|-------------|
| -32700 | Parse error |
| -32600 | Invalid request |
| -32601 | Method not found |
| -32602 | Invalid params |
| -32603 | Internal error |
| -32001 | Agent not found |
| -32002 | Task not found |
| -32003 | Agent unavailable |

## Handshake Inicial

Ao iniciar comunicação, agentes trocam Agent Cards:

```
Client: "Olá, aqui estão minhas capacidades"
Server: "Entendi, aqui estão as minhas"
Client: "Perfeito, vou enviar uma tarefa"
```

---

*Referência: a2a-protocol.org*