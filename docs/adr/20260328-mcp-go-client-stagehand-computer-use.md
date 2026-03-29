# ADR 20260328: MCP Client Go para Stagehand (Computer Use Integration)

## Status
🔴 Proposto (Bloqueado por P1)

## Contexto
O ADR `20260328-implementacao-jarvis-voice-e-computer-use` define que o `computer use` será controlado via **Model Context Protocol (MCP)**. O Stagehand MCP Server (`mcp-servers/stagehand/src/index.ts`) já está esboçado com tools `navigate`, `act` e `extract`. O que falta é o **MCP Client em Go** que expõe essas ferramentas como ferramentas nativas do `aurelia_code`.

## Decisões Arquiteturais

### 1. Arquitetura de Client MCP em Go
Implementar um wrapper nativo em Go utilizando o SDK `github.com/modelcontextprotocol/sdk-go` (ou equivalente mantenido).

```go
// internal/mcp/client.go
type StagehandClient struct {
    cmd    *exec.Cmd
    stdio  *stdio.Channel
    models *sdk.Client
}

func NewStagehandClient(ctx context.Context) (*StagehandClient, error)
func (c *StagehandClient) Navigate(ctx context.Context, url string) error
func (c *StagehandClient) Act(ctx context.Context, instruction string) error
func (c *StagehandClient) Extract(ctx context.Context, instruction string) (string, error)
func (c *StagehandClient) Close() error
```

### 2. Tool Definitions em `internal/tools/`
Integrar com o sistema de ferramentas existente:

```go
// internal/tools/stagehand_tools.go
var stagehandToolDefs = []ToolDef{
    {
        Name:        "mcp__stagehand__navigate",
        Description: "Navega para uma URL no browser isolado",
        Params:      z.Object({"url": z.String().url()}),
        Handler:     handleStagehandNavigate,
    },
    {
        Name:        "mcp__stagehand__act",
        Description: "Executa uma ação no browser (clique, digitação, etc)",
        Params:      z.Object({"instruction": z.String()}),
        Handler:     handleStagehandAct,
    },
    {
        Name:        "mcp__stagehand__extract",
        Description: "Extrai dados da página atual via instrução",
        Params:      z.Object({"instruction": z.String()}),
        Handler:     handleStagehandExtract,
    },
}
```

### 3. Process Management (Steel Container)
Gerenciar o lifecycle do servidor MCP como subprocess isolado:

- **Start**: Executar `node dist/index.js` dentro do container `steel`
- **Health check**: Ping via stdio/json-rpc
- **Restart policy**: Auto-restart em caso de crash
- **Timeout**: Kill after 10min de inatividade

### 4. Error Handling
- Server não está rodando → Tool retorna erro estruturado, não crash
- Timeout → Cancelar comando, marcar tool como `unavailable`
- Screenshot on error → Capturar estado para debug

## Consequências

### Positivas
- Computer use exposto como tool nativa do agente
- Isolamento via container previne crashes em cascata
- Reutilização do Stagehand SDK sem reimplementar browser automation

### Negativas
- Dependência de Node.js runtime no container
- Latência adicional na comunicação stdio
- Complexidade de debug quando o server morre silenciosamente

### Trade-offs
- stdio vs HTTP: stdio é mais simples mas não permite hot-reload
- Go SDK vs custom parsing: SDK oficial reduz manutenibilidade

## Dependências
- ✅ `mcp-servers/stagehand/src/index.ts` (já existe)
- ⚠️ `github.com/modelcontextprotocol/sdk-go` ou fork mantido
- ❌ `configs/litellm/config.yaml` (blocker: precisa existir primeiro)

## Referências
- [Stagehand SDK](https://github.com/browserbasehq/stagehand)
- [MCP Protocol Spec](https://modelcontextprotocol.io)
- [internal/tools/definitions.go](../../internal/tools/definitions.go)
- [ADR-20260328-jarvis-voice-computer-use.json](./taskmaster/ADR-20260328-jarvis-voice-computer-use.json)

## Links Obrigatórios
- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/neon-sentinel
**Progress**: 0%
