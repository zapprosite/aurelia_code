package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockLLM para simular respostas do agente chamando a tool de handoff
type MockHandoffLLM struct {
	HandoffCalled bool
}

func (m *MockHandoffLLM) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	if !m.HandoffCalled {
		m.HandoffCalled = true
		return &ModelResponse{
			ReasoningContent: "Vou passar para o coder.",
			ToolCalls: []ToolCall{
				{
					ID:   "call_1",
					Name: "handoff_to_agent",
					Arguments: map[string]interface{}{
						"target_agent":     "coder",
						"task_description": "Escreva o código solicitado.",
						"reason":           "Preciso de um especialista em implementação.",
					},
				},
			},
		}, nil
	}
	return &ModelResponse{Content: "Ok, sou o coder e terminei."}, nil
}

func TestNativeHandoffSimulation(t *testing.T) {
	ctx := context.Background()
	// Usar um path temporário único para evitar conflitos de teste
	dbPath := "/tmp/handoff_test_" + uuid.NewString() + ".db"
	store, err := NewSQLiteTaskStore(dbPath)
	require.NoError(t, err)
	defer store.Close()

	mgr, err := NewTeamManager(store)
	require.NoError(t, err)

	llm := &MockHandoffLLM{}
	registry := NewToolRegistry()
	
	service, err := NewMasterTeamService(mgr, llm, registry, 5, nil)
	require.NoError(t, err)

	// Registrar as ferramentas
	registry.Register(GetHandoffToolDefinition(), HandoffHandler(service))

	teamKey := "test_team_" + uuid.NewString()
	userID := "user_1"

	// 1. Spawn do primeiro agente (Planner)
	taskID, err := service.Spawn(ctx, teamKey, userID, "planner", "Faz o plano", "Inicie o projeto")
	require.NoError(t, err)
	assert.NotEmpty(t, taskID)

	// 2. Simular a execução do loop (o que o workerLoop faria)
	teamID, _ := mgr.GetTeamIDByKey(ctx, teamKey)
	task, _ := mgr.ClaimNextTask(ctx, teamID, "planner")
	require.NotNil(t, task)

	loop := NewLoop(llm, registry, 5)
	history := []Message{{Role: "user", Content: task.Prompt}}
	
	// Injetar contexto de task para que o handoff funcione
	taskCtx := WithTaskContext(ctx, teamID, task.ID)
	taskCtx = WithRunContext(taskCtx, task.RunID)
	taskCtx = WithAgentContext(taskCtx, "planner")

	newHistory, _, err := loop.Run(taskCtx, "Prompt sistema", history, []string{"handoff_to_agent"})
	require.NoError(t, err)

	// 3. Verificar se o handoff ocorreu
	foundHandoff := false
	for _, msg := range newHistory {
		if msg.Role == "tool" && msg.Content != "" {
			if strings.Contains(msg.Content, "handoff_to_agent") || strings.Contains(msg.Content, "Handoff iniciado") {
				foundHandoff = true
			}
		}
	}
	assert.True(t, foundHandoff, "Deveria ter encontrado a confirmação de handoff no histórico")

	// 4. Verificar se a nova task para o 'coder' foi criada no DB
	tasks, err := mgr.ListTasks(ctx, teamID)
	require.NoError(t, err)
	
	var coderTask *TeamTask
	for _, tk := range tasks {
		if tk.Title == "coder" {
			coderTask = &tk
		}
	}
	require.NotNil(t, coderTask, "Uma nova task para o 'coder' deveria ter sido criada")
	assert.Equal(t, TaskPending, coderTask.Status)
	assert.Contains(t, coderTask.Prompt, "Escreva o código solicitado.")
}
