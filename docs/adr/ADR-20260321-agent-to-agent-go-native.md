# ADR-20260321: Comunicação Agent-to-Agent (Go Nativo)

## Status
**Proposta** (Sênior Industrial 2026)

## Contexto
Atualmente, o Aurelia depende de lógica externa (Swarm/Python) ou de implementações de ferramentas individuais para delegar tarefas entre agentes (`spawn_agent`). Isso introduz:
1.  Dependência de runtimes Python.
2.  Latência na troca de contexto.
3.  Dificuldade em manter um estado compartilhado real-time (Blackboard).

## Decisão
Implementar um **Coordinator** nativo no diretório `internal/agent` que gerencie o enxame de agentes diretamente em Go.

### Arquitetura (NG-03)
- **Coordinator**: Um orquestrador que mantém o registro de agentes ativos e suas capacidades.
- **Handoff Nativo**: Substituir o `spawn_agent` atual por uma chamada interna que transfere a `agent.Message` history de forma síncrona/assíncrona via canais Go.
- **Shared Blackboard**: Utilizar o Postgres (via Orchid ORM) para persistir o estado de curto prazo compartilhado entre agentes da mesma sessão.
- **Real-Time Integration**: Cada troca de agente deve emitir um evento `dashboard.Event` com tipo `agent_handoff`.

## Consequências
- **Positivas**: 
    - Eliminação completa de Python para orquestração core.
    - Observabilidade total no dashboard (quem está com o "token" de execução).
    - Latência zero na delegação.
- **Negativas**:
    - Maior complexidade no kernel Go inicialmente.

## Verificação (Smoke Tests)
- `go test ./internal/agent/...` deve validar handoffs entre 2 instâncias de `Loop`.
- Simulação de "Time de Resposta" com 3 agentes (Planner -> Developer -> Reviewer).
