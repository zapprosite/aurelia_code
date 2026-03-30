package skill

import (
	"context"

	"github.com/kocar/aurelia/internal/agent"
)

// Executor manages the runtime context injection limits
type Executor struct {
	loop *agent.Loop
}

// NewExecutor constructs an executor
func NewExecutor(loop *agent.Loop) *Executor {
	return &Executor{loop: loop}
}

// GetLoop returns the underlying agent loop
func (e *Executor) GetLoop() *agent.Loop {
	return e.loop
}

// Execute wraps the ReAct loop injecting the targeted skill file cleanly without polluting history
func (e *Executor) Execute(ctx context.Context, baseSystemPrompt string, skill *Skill, history []agent.Message, allowedTools []string) ([]agent.Message, string, error) {
	finalSystemPrompt := buildSkillPrompt(baseSystemPrompt, skill)
	return e.loop.Run(ctx, finalSystemPrompt, history, allowedTools)
}

// ExecuteStream is the streaming variant — returns a channel of tokens/tool events.
func (e *Executor) ExecuteStream(ctx context.Context, baseSystemPrompt string, skill *Skill, history []agent.Message, allowedTools []string) (<-chan agent.StreamResponse, error) {
	finalSystemPrompt := buildSkillPrompt(baseSystemPrompt, skill)

	opts := agent.LoopOptions{
		SystemPrompt:    finalSystemPrompt,
		InitialHistory:  history,
		MaxIterations:   e.loop.MaxIterations(),
		ToolDefinitions: e.loop.ResolveToolDefinitions(history, allowedTools),
	}
	return e.loop.RunWithOptionsStream(ctx, opts)
}

func buildSkillPrompt(base string, skill *Skill) string {
	if skill != nil {
		return base + "\n\n# ACTIVE SKILL CONTEXT:\n" + skill.Content
	}
	return base
}
