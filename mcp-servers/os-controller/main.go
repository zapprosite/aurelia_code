package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/kocar/aurelia/internal/os_controller"
)

func main() {
	// Inicializar o Controller do God Mode
	ctrl := os_controller.NewController(60*time.Second, false)

	// Criar o servidor MCP
	s := server.NewMCPServer(
		"aurelia-god-mode",
		"1.0.0",
		server.WithLogging(),
	)

	// Tool: run_bash_command
	runBashTool := mcp.NewTool("run_bash_command",
		mcp.WithDescription("Executa um comando Bash no sistema host com proteção do Execution Guard"),
		mcp.WithInputSchema[string](),
	)

	s.AddTool(runBashTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		script := mcp.ParseString(request, "script", "")
		if script == "" {
			return mcp.NewToolResultError("script is required"), nil
		}

		result, err := ctrl.RunBash(ctx, script)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("God Mode Guard Block: %v", err)), nil
		}

		respMsg := fmt.Sprintf("Exit Code: %d\nStdout:\n%s\nStderr:\n%s", result.ExitCode, result.Stdout, result.Stderr)
		return mcp.NewToolResultText(respMsg), nil
	})

	// Tool: read_system_log
	readLogTool := mcp.NewTool("read_system_log",
		mcp.WithDescription("Lê as últimas N linhas de um log de sistema"),
		mcp.WithInputSchema[struct {
			Path  string `json:"path"`
			Lines int    `json:"lines"`
		}](),
	)

	s.AddTool(readLogTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := mcp.ParseString(request, "path", "")
		if path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		lines := mcp.ParseInt(request, "lines", 50)

		logLines, err := ctrl.ReadLog(ctx, path, lines)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("error reading log: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Total lines: %d\nContent:\n%s", len(logLines), strings.Join(logLines, "\n"))), nil
	})

	// Iniciar servidor em modo StdIO
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
