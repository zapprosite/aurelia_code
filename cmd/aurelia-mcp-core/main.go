package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/kocar/aurelia/internal/tools"
)

func main() {
	// Criar o servidor MCP Soberano 2026
	s := server.NewMCPServer(
		"aurelia-mcp-core",
		"1.0.0",
		server.WithLogging(),
	)

	// 1. Filesystem Tools
	s.AddTool(mcp.NewTool("read_file",
		"Le o conteudo de um arquivo local.",
		mcp.NewToolSchema(tools.ReadFileDefinition().JSONSchema),
	), wrapHandler(tools.ReadFileHandler))

	s.AddTool(mcp.NewTool("write_file",
		"Escreve conteudo integral em um arquivo.",
		mcp.NewToolSchema(tools.WriteFileDefinition().JSONSchema),
	), wrapHandler(tools.WriteFileHandler))

	s.AddTool(mcp.NewTool("list_dir",
		"Lista os arquivos dentro de um diretorio.",
		mcp.NewToolSchema(tools.ListDirDefinition().JSONSchema),
	), wrapHandler(tools.ListDirHandler))

	// 2. System Control Tools
	s.AddTool(mcp.NewTool("run_command",
		"Executa um comando local (Bash) de forma controlada no Ubuntu 24.04.",
		mcp.NewToolSchema(tools.RunCommandDefinition().JSONSchema),
	), wrapHandler(tools.RunCommandHandler))

	s.AddTool(mcp.NewTool("docker_control",
		"Controla containers Docker e Docker Compose.",
		mcp.NewToolSchema(tools.DockerControlDefinition().JSONSchema),
	), wrapHandler(tools.DockerControlHandler))

	s.AddTool(mcp.NewTool("service_control",
		"Controla serviços systemd (status, restart, logs).",
		mcp.NewToolSchema(tools.ServiceControlDefinition().JSONSchema),
	), wrapHandler(tools.ServiceControlHandler))

	s.AddTool(mcp.NewTool("system_monitor",
		"Monitora recursos do sistema (CPU, RAM, GPU, Processos).",
		mcp.NewToolSchema(tools.SystemMonitorDefinition().JSONSchema),
	), wrapHandler(tools.SystemMonitorHandler))

	// 3. AI & Search Tools
	s.AddTool(mcp.NewTool("ollama_control",
		"Controla o Ollama local (list, pull, run).",
		mcp.NewToolSchema(tools.OllamaControlDefinition().JSONSchema),
	), wrapHandler(tools.OllamaControlHandler))

	s.AddTool(mcp.NewTool("web_search",
		"Pesquisa na internet via DuckDuckGo.",
		mcp.NewToolSchema(tools.WebSearchDefinition().JSONSchema),
	), wrapHandler(tools.WebSearchHandler))

	// Iniciar servidor via Stdio
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "MCP Server Error: %v\n", err)
		os.Exit(1)
	}
}

// wrapHandler adapta o handler interno da Aurelia para o formato mcp-go
func wrapHandler(h func(context.Context, map[string]interface{}) (string, error)) server.Handler {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := h(ctx, request.Arguments)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(res), nil
	}
}
