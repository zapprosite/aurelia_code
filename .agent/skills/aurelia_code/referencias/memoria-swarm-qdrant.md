# Memória Swarm — Qdrant

## Visão Geral

A memória compartilhada é o que transforma um grupo de agentes em um **swarm verdadeiro**. Sem memória compartilhada, cada agente é isolado; com ela, o swarm age como um organismo coletivo.

## Arquitetura de Memória

```
┌─────────────────────────────────────────────────────────────┐
│                     Qdrant Cluster                          │
├─────────────────────────────────────────────────────────────┤
│  Collections:                                               │
│  ├── aurelia_swarm_missions   (missões)                    │
│  ├── aurelia_swarm_context    (contexto compartilhado)      │
│  ├── aurelia_swarm_decisions  (decisões do líder)          │
│  ├── aurelia_code_memory      (memória persistente)        │
│  └── conversation_memory      (conversas)                  │
└─────────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────┐
│  agentes:                                                   │
│  ├── aurelia_code (líder)                                   │
│  ├── Pesquisador                                            │
│  ├── Coder                                                  │
│  └── Revisor                                               │
└─────────────────────────────────────────────────────────────┘
```

## Collections

### aurelia_swarm_missions

Armazena missões completas:

```json
{
  "mission_id": "uuid",
  "created_at": "2026-04-01T10:00:00Z",
  "status": "active|completed|cancelled|failed",
  "leader": "aurelia_code",
  "description": "Implementar sistema de login",
  "tasks": [
    {
      "task_id": "uuid",
      "role": "pesquisador",
      "agent_id": "sub-001",
      "description": "Pesquisar bibliotecas de auth",
      "status": "completed",
      "result": "...",
      "created_at": "...",
      "completed_at": "..."
    }
  ],
  "artifacts": [...],
  "errors": [...],
  "metrics": {
    "total_tasks": 5,
    "completed": 4,
    "failed": 1,
    "duration_minutes": 45
  }
}
```

### aurelia_swarm_context

Contexto compartilhado entre agentes:

```json
{
  "context_id": "uuid",
  "mission_id": "uuid",
  "shared_data": {
    "project_name": "auth-system",
    "tech_stack": ["go", "postgres", "redis"],
    "constraints": ["sem libs externas", "100% coverage"],
    "decisions": [
      {"decision": "usar bcrypt", "reason": "segurança", "approved_by": "aurelia_code"}
    ]
  },
  "agents_context": {
    "pesquisador": {"last_search": "...", "docs_loaded": [...]},
    "coder": {"files_modified": [...], "tests_added": [...]},
    "revisor": {"comments": [...], "approved_files": [...]}
  }
}
```

### aurelia_swarm_decisions

Decisões do líder documentadas:

```json
{
  "decision_id": "uuid",
  "mission_id": "uuid",
  "type": "delegation|architecture|priority|tradeoff",
  "description": "Delegar pesquisa para agente",
  "rationale": "Pesquisador tem acesso a Tavily e Context7",
  "alternatives_considered": [
    "Fazer eu mesmo - tomaria muito tempo",
    "Usar Coder - não tem ferramentas de pesquisa"
  ],
  "outcome": "positive|negative|pending",
  "timestamp": "2026-04-01T10:05:00Z"
}
```

### aurelia_code_memory

Memória de longo prazo do líder:

```json
{
  "memory_id": "uuid",
  "type": "pattern|lesson|preference|skill",
  "content": "Quando delegar tarefas de pesquisa, sempre especificar prazo",
  "source_mission": "uuid",
  "confidence": 0.85,
  "tags": ["delegation", "research", "timing"],
  "created_at": "2026-04-01T11:00:00Z"
}
```

## Schema Qdrant

### Pontos (Vectors)

Cada memória é um ponto vetorial:

```go
type MemoryPoint struct {
    ID        string `json:"id"`
    Vector    []float32
    Payload   map[string]interface{}
}
```

### Vector Configuration

```json
{
  "size": 768,
  "distance": "Cosine"
}
```

## Escrita (Write)

### Após cada tarefa completada:

```go
func (s *SwarmMemory) WriteResult(ctx context.Context, task Task, result string) error {
    point := &MemoryPoint{
        ID:      task.ID,
        Vector:  embed(result), // embedding da resposta
        Payload: map[string]interface{}{
            "mission_id": task.MissionID,
            "role":       task.Role,
            "result":     result,
            "timestamp":  time.Now().UTC(),
            "status":     "completed",
        },
    }
    return s.client.Upsert("aurelia_swarm_missions", point)
}
```

### Após decisão do líder:

```go
func (s *SwarmMemory) WriteDecision(ctx context.Context, d Decision) error {
    point := &MemoryPoint{
        ID:      d.ID,
        Vector:  embed(d.Description),
        Payload: map[string]interface{}{
            "mission_id":   d.MissionID,
            "type":         d.Type,
            "description":  d.Description,
            "rationale":    d.Rationale,
            "outcome":      d.Outcome,
            "timestamp":    time.Now().UTC(),
        },
    }
    return s.client.Upsert("aurelia_swarm_decisions", point)
}
```

## Leitura (Read)

### Buscar contexto relevante:

```go
func (s *SwarmMemory) GetContext(ctx context.Context, missionID string, query string) ([]MemoryContext, error) {
    queryVector := embed(query)
    
    results, err := s.client.Search(
        "aurelia_swarm_context",
        queryVector,
        filters.MustEqual("mission_id", missionID),
        limit=5,
    )
    
    var contexts []MemoryContext
    for _, r := range results {
        contexts = append(contextes, r.Payload)
    }
    return contexts, nil
}
```

### Buscar decisões passadas similares:

```go
func (s *SwarmMemory) GetSimilarDecisions(ctx context.Context, query string, limit int) []Decision {
    queryVector := embed(query)
    
    results, _ := s.client.Search(
        "aurelia_swarm_decisions",
        queryVector,
        filters.Should(filters.Eq("type", "delegation")),
        limit=limit,
    )
    
    // Retornar decisões similares
}
```

## Padrões de Uso

### 1. Início de Missão
```bash
# Buscar类似 missões anteriores
curl -X POST "http://localhost:6333/collections/aurelia_swarm_missions/points/search" \
  -d '{"query": embed("similar mission"), "limit": 3}'
```

### 2. Durante Delegação
```bash
# Buscar contexto existente
curl -X POST "http://localhost:6333/collections/aurelia_swarm_context/points/search" \
  -d '{"query": embed("project auth"), "limit": 5}'
```

### 3. Após Conclusão
```bash
# Armazenar resultado
curl -X POST "http://localhost:6333/collections/aurelia_swarm_missions/points" \
  -d '{"points": [{"id": "task-001", "vector": [...], "payload": {...}}]}'
```

## Integração com Agentes

### Sub-agent recebe contexto:

```go
func (s *A2AAgent) ExecuteTask(ctx context.Context, task Task) Result {
    // 1. Buscar contexto relevante
    context, _ := s.memory.GetContext(ctx, task.MissionID, task.Description)
    
    // 2. Adicionar ao prompt
    prompt := task.Description + "\n\nContexto relevante:\n" + context
    
    // 3. Executar
    result := s.llm.Execute(prompt)
    
    // 4. Armazenar resultado
    s.memory.WriteResult(ctx, task, result)
    
    return result
}
```

---

*Qdrant é a memória do swarm — sem ela, não há aprendizado coletivo.*