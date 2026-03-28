package main

import (
	"context"
	"encoding/json"
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
	)
	// Adicionando schema via AddTool manualmente se o helper estiver ausente
	runBashTool.InputSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"script": map[string]interface{}{
				"type":        "string",
				"description": "O script bash ou comando a ser executado",
			},
		},
		"required": []string{"script"},
	}

	s.AddTool(runBashTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Type assertion para []byte
		argsRaw, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments format"), nil
		}
		
		script, _ := argsRaw["script"].(string)
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
	)
	readLogTool.InputSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Caminho absoluto do arquivo de log",
			},
			"lines": map[string]interface{}{
				"type":        "integer",
				"description": "Número de linhas a serem lidas (default: 50)",
			},
		},
		"required": []string{"path"},
	}

	s.AddTool(readLogTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		argsRaw, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid arguments format"), nil
		}

		path, _ := argsRaw["path"].(string)
		linesFloat, _ := argsRaw["lines"].(float64) // JSON numbers are float64 in interface{}
		lines := int(linesFloat)

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
