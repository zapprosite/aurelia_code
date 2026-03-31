// Package computer_use provides autonomous GUI navigation capabilities
// via Stagehand MCP for browser automation.
// ADR: 20260328-computer-use-e2e-autonomous-gui

package computer_use

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/mcp"
)

// stagehandServerName is the name of the stagehand MCP server
const stagehandServerName = "stagehand"

// Agent is the main computer use agent that orchestrates browser automation
// ADR: 20260328-computer-use-e2e-autonomous-gui
type Agent struct {
	llm      LLMClient
	mcpMgr   *mcp.Manager
	maxSteps int
}

// LLMClient interface for LLM communication
type LLMClient interface {
	Chat(ctx context.Context, messages []Message, system string) (string, error)
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AgentConfig configures the computer use agent
type AgentConfig struct {
	LLM             LLMClient
	MCPManager      *mcp.Manager
	MaxSteps        int
	ScreenshotBudget int
}

// NewAgent creates a new computer use agent
// ADR: 20260328-computer-use-e2e-autonomous-gui
func NewAgent(cfg AgentConfig) *Agent {
	if cfg.MaxSteps <= 0 {
		cfg.MaxSteps = 10
	}

	return &Agent{
		llm:      cfg.LLM,
		mcpMgr:   cfg.MCPManager,
		maxSteps: cfg.MaxSteps,
	}
}

// AgentState represents the current state of the agent
type AgentState struct {
	Intent string      `json:"intent"`
	Steps  []ActionStep `json:"steps"`
	Done   bool        `json:"done"`
	Summary string      `json:"summary,omitempty"`
}

// ActionStep represents a single action taken by the agent
type ActionStep struct {
	Action     string         `json:"action"`
	Params     map[string]any `json:"params"`
	Result     string         `json:"result"`
	Error      string         `json:"error,omitempty"`
	Screenshot string         `json:"screenshot,omitempty"`
}

// Run executes the computer use agent to accomplish the given intent
// ADR: 20260328-computer-use-e2e-autonomous-gui
func (a *Agent) Run(ctx context.Context, intent string) (*AgentState, error) {
	state := &AgentState{
		Intent: intent,
		Steps:  make([]ActionStep, 0),
		Done:   false,
	}

	for step := 0; step < a.maxSteps; step++ {
		// 1. Observe - capture screenshot
		screenshot, err := a.captureScreenshot(ctx)
		if err != nil {
			state.Steps = append(state.Steps, ActionStep{
				Action: "observe",
				Result: "",
				Error:  err.Error(),
			})
			continue
		}

		// 2. Decide - ask LLM what to do next

		// 2. Decide - ask LLM what to do next
		decision, err := a.decideAction(ctx, state.Intent, screenshot, state.Steps)
		if err != nil {
			return nil, fmt.Errorf("decision failed: %w", err)
		}

		if decision.Type == "done" {
			state.Done = true
			state.Summary = decision.Result
			break
		}

		// 3. Validate action (safety guardrails)
		if err := a.validateAction(decision); err != nil {
			state.Steps = append(state.Steps, ActionStep{
				Action: decision.Type,
				Params: decision.Params,
				Error:  fmt.Sprintf("blocked: %s", err.Error()),
			})
			continue
		}

		// 4. Act - execute the action
		result, err := a.executeAction(ctx, decision)
		stepResult := ActionStep{
			Action:     decision.Type,
			Params:     decision.Params,
			Result:     result,
			Screenshot: screenshot,
		}
		if err != nil {
			stepResult.Error = err.Error()
		}

		state.Steps = append(state.Steps, stepResult)

		// Small delay between actions
		time.Sleep(500 * time.Millisecond)
	}

	if !state.Done && len(state.Steps) >= a.maxSteps {
		state.Summary = fmt.Sprintf("Agent reached max steps (%d) without completing", a.maxSteps)
	}

	return state, nil
}

// captureScreenshot captures a screenshot via MCP
func (a *Agent) captureScreenshot(ctx context.Context) (string, error) {
	if a.mcpMgr == nil {
		return "", fmt.Errorf("MCP manager not initialized")
	}

	result, err := a.mcpMgr.CallTool(ctx, stagehandServerName, "screenshot", nil)
	if err != nil {
		return "", fmt.Errorf("screenshot failed: %w", err)
	}

	if result.IsError {
		return "", fmt.Errorf("screenshot error: %s", result.Content)
	}

	return result.Content, nil
}

// Decision represents an LLM decision
type Decision struct {
	Type      string         `json:"type"` // "navigate", "act", "extract", "done"
	Params   map[string]any `json:"params,omitempty"`
	Result   string         `json:"result,omitempty"`
	Reasoning string         `json:"reasoning,omitempty"`
}

// decideAction asks the LLM to decide the next action
// ADR: 20260328-computer-use-e2e-autonomous-gui
func (a *Agent) decideAction(ctx context.Context, intent string, screenshot string, steps []ActionStep) (*Decision, error) {
	if a.llm == nil {
		return nil, fmt.Errorf("no LLM client configured")
	}

	// Build prompt
	prompt := a.buildDecisionPrompt(intent, screenshot, steps)

	// Call LLM
	response, err := a.llm.Chat(ctx, []Message{
		{Role: "user", Content: prompt},
	}, "")
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// Parse response
	return parseDecision(response)
}

// buildDecisionPrompt creates the prompt for the LLM
func (a *Agent) buildDecisionPrompt(intent string, screenshot string, steps []ActionStep) string {
	var sb strings.Builder

	sb.WriteString("Voce e um agente de computer use. Analise o estado atual e decida a proxima acao.\n\n")
	sb.WriteString("Intent original: " + intent + "\n\n")

	if len(steps) > 0 {
		sb.WriteString("Historico de acoes:\n")
		for i, step := range steps {
			if step.Error != "" {
				sb.WriteString(fmt.Sprintf("%d. %s(%v) → ERRO: %s\n", i+1, step.Action, step.Params, step.Error))
			} else {
				sb.WriteString(fmt.Sprintf("%d. %s(%v) → %s\n", i+1, step.Action, step.Params, step.Result))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("Screenshot atual: [presente]\n\n")
	sb.WriteString("Decida a proxima acao (responda em JSON):\n")
	sb.WriteString(`{
    "type": "navigate|act|extract|done",
    "params": {"url": "...", "instruction": "...", "query": "..."},
    "reasoning": "por que estou tomando esta decisao"
}`)

	return sb.String()
}

// executeAction performs the decided action
// ADR: 20260328-computer-use-e2e-autonomous-gui
func (a *Agent) executeAction(ctx context.Context, decision *Decision) (string, error) {
	if a.mcpMgr == nil {
		return "", fmt.Errorf("MCP manager not initialized")
	}

	switch decision.Type {
	case "navigate":
		url, _ := decision.Params["url"].(string)
		if url == "" {
			return "", fmt.Errorf("url is required for navigate")
		}
		result, err := a.mcpMgr.CallTool(ctx, stagehandServerName, "navigate", map[string]interface{}{
			"url": url,
		})
		if err != nil {
			return "", fmt.Errorf("navigate failed: %w", err)
		}
		if result.IsError {
			return "", fmt.Errorf("navigate error: %s", result.Content)
		}
		return fmt.Sprintf("Navegado para %s", url), nil

	case "act":
		instruction, _ := decision.Params["instruction"].(string)
		if instruction == "" {
			return "", fmt.Errorf("instruction is required for act")
		}
		result, err := a.mcpMgr.CallTool(ctx, stagehandServerName, "act", map[string]interface{}{
			"instruction": instruction,
		})
		if err != nil {
			return "", fmt.Errorf("act failed: %w", err)
		}
		if result.IsError {
			return "", fmt.Errorf("act error: %s", result.Content)
		}
		return fmt.Sprintf("Acao executada: %s", instruction), nil

	case "extract":
		query, _ := decision.Params["query"].(string)
		if query == "" {
			return "", fmt.Errorf("query is required for extract")
		}
		result, err := a.mcpMgr.CallTool(ctx, stagehandServerName, "extract", map[string]interface{}{
			"instruction": query,
		})
		if err != nil {
			return "", fmt.Errorf("extract failed: %w", err)
		}
		if result.IsError {
			return "", fmt.Errorf("extract error: %s", result.Content)
		}
		return fmt.Sprintf("Extraido: %s", result.Content), nil

	case "done":
		return decision.Result, nil

	default:
		return "", fmt.Errorf("unknown action type: %s", decision.Type)
	}
}

// Dangerous patterns for safety guardrails
var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(rm\s+-rf|dd\s+|mkfs|wipefs|shred)`),
	regexp.MustCompile(`(?i)(sudo\s+rm|chmod\s+777|ssh\s+.*@)`),
	regexp.MustCompile(`(?i)(curl\s+.*\|\s*sh|wget\s+.*\|\s*bash)`),
	regexp.MustCompile(`(?i)(drop\s+database|delete\s+from\s+.*where)`),
	regexp.MustCompile(`(?i)(format\s+c:|#!/bin/bash.*rm\s+)`),
}

// validateAction checks if an action is safe to execute
// ADR: 20260328-computer-use-e2e-autonomous-gui
func (a *Agent) validateAction(decision *Decision) error {
	if decision.Type == "done" {
		return nil
	}

	// Check for dangerous patterns in all string params
	for _, v := range decision.Params {
		if str, ok := v.(string); ok {
			for _, pattern := range dangerousPatterns {
				if pattern.MatchString(str) {
					return fmt.Errorf("blocked dangerous pattern: %s", pattern.String())
				}
			}
		}
	}

	return nil
}

// parseDecision parses the LLM response into a Decision
func parseDecision(response string) (*Decision, error) {
	// Try to extract JSON from response
	response = strings.TrimSpace(response)

	// Remove markdown code blocks if present
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var decision Decision
	if err := json.Unmarshal([]byte(response), &decision); err != nil {
		return nil, fmt.Errorf("failed to parse decision JSON: %w", err)
	}

	if decision.Type == "" {
		return nil, fmt.Errorf("decision type is required")
	}

	return &decision, nil
}
