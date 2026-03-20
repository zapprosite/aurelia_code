package agent

import (
"fmt"
)

// MCPClient gerencia a execução de ferramentas e fluxos compostos.
type MCPClient struct {
serverURI string
}

func NewMCPClient(uri string) *MCPClient {
return &MCPClient{serverURI: uri}
}

// ComposeFlow solicita a execução de uma sequência macroscópica de ferramentas.
func (c *MCPClient) ComposeFlow(flowName string, params map[string]interface{}) string {
fmt.Printf("[MCP-CLIENT] Executando Fluxo Composto: %s no servidor %s...\n", flowName, c.serverURI)
return fmt.Sprintf("RESULTADO_FLUXO_%s_SUCCESS", flowName)
}

// CallTool executa uma ferramenta atômica via MCP.
func (c *MCPClient) CallTool(toolName string, args string) string {
fmt.Printf("[MCP-CLIENT] Chamando ferramenta %s...\n", toolName)
return "TOOL_SUCCESS"
}
