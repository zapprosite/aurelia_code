// Package tools provides tool definitions and handlers for the Aurelia agent
// Computer Use handlers via MCP Manager integration
// ADR: 20260328-mcp-go-client-stagehand-computer-use

package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/kocar/aurelia/internal/mcp"
)

// stagehandServerName is the name of the stagehand server in MCP config
const stagehandServerName = "stagehand"

// globalMCPManager stores the MCP manager for computer use tools
// Set during app initialization via SetGlobalMCPManager()
var (
	globalMCPManager mcp.Caller
	globalMCPMu     sync.RWMutex
)

// SetGlobalMCPManager sets the global MCP manager for computer use tools
func SetGlobalMCPManager(manager mcp.Caller) {
	globalMCPMu.Lock()
	defer globalMCPMu.Unlock()
	globalMCPManager = manager
}

// GetMCPManager returns the global MCP manager instance
// ADR: 20260328-mcp-go-client-stagehand-computer-use
func ComputerNavigateHandler(ctx context.Context, params map[string]any) (string, error) {
	manager := GetMCPManager()
	if manager == nil {
		return "", fmt.Errorf("MCP manager not initialized")
	}

	url, ok := params["url"].(string)
	if !ok || url == "" {
		return "", fmt.Errorf("url parameter is required")
	}

	result, err := manager.CallTool(ctx, stagehandServerName, "navigate", map[string]interface{}{
		"url": url,
	})
	if err != nil {
		return "", fmt.Errorf("navigate failed: %w", err)
	}

	if result.IsError {
		return "", fmt.Errorf("navigate error: %s", result.Content)
	}

	return fmt.Sprintf("Navegado com sucesso para %s", url), nil
}

// ComputerActHandler handles the mcp__stagehand__act tool
// ADR: 20260328-mcp-go-client-stagehand-computer-use
func ComputerActHandler(ctx context.Context, params map[string]any) (string, error) {
	manager := GetMCPManager()
	if manager == nil {
		return "", fmt.Errorf("MCP manager not initialized")
	}

	instruction, ok := params["instruction"].(string)
	if !ok || instruction == "" {
		return "", fmt.Errorf("instruction parameter is required")
	}

	result, err := manager.CallTool(ctx, stagehandServerName, "act", map[string]interface{}{
		"instruction": instruction,
	})
	if err != nil {
		return "", fmt.Errorf("act failed: %w", err)
	}

	if result.IsError {
		return "", fmt.Errorf("act error: %s", result.Content)
	}

	return fmt.Sprintf("Acao executada: %s", instruction), nil
}

// ComputerExtractHandler handles the mcp__stagehand__extract tool
// ADR: 20260328-mcp-go-client-stagehand-computer-use
func ComputerExtractHandler(ctx context.Context, params map[string]any) (string, error) {
	manager := GetMCPManager()
	if manager == nil {
		return "", fmt.Errorf("MCP manager not initialized")
	}

	instruction, ok := params["instruction"].(string)
	if !ok || instruction == "" {
		return "", fmt.Errorf("instruction parameter is required")
	}

	result, err := manager.CallTool(ctx, stagehandServerName, "extract", map[string]interface{}{
		"instruction": instruction,
	})
	if err != nil {
		return "", fmt.Errorf("extract failed: %w", err)
	}

	if result.IsError {
		return "", fmt.Errorf("extract error: %s", result.Content)
	}

	return result.Content, nil
}

// ComputerScreenshotHandler handles the mcp__stagehand__screenshot tool
// ADR: 20260328-vision-pipeline-computer-use
func ComputerScreenshotHandler(ctx context.Context, params map[string]any) (string, error) {
	manager := GetMCPManager()
	if manager == nil {
		return "", fmt.Errorf("MCP manager not initialized")
	}

	result, err := manager.CallTool(ctx, stagehandServerName, "screenshot", map[string]interface{}{})
	if err != nil {
		return "", fmt.Errorf("screenshot failed: %w", err)
	}

	if result.IsError {
		return "", fmt.Errorf("screenshot error: %s", result.Content)
	}

	return fmt.Sprintf("Screenshot capturado: %s", result.Content), nil
}

// GetMCPManager returns the global MCP manager instance
func GetMCPManager() mcp.Caller {
	globalMCPMu.RLock()
	defer globalMCPMu.RUnlock()
	return globalMCPManager
}
