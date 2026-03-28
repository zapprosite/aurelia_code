package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProposePlanTool_RequiresBackoutPlan(t *testing.T) {
	t.Parallel()
	SetupTest()

	reg := NewToolRegistry()
	reg.RegisterPlannerTools()

	var planner *Tool
	for _, d := range reg.GetDefinitions() {
		if d.Name == "propose_plan" {
			planner = &d
			break
		}
	}

	require.NotNil(t, planner, "propose_plan tool should be registered")

	// Verificar se 'backout_plan' está na lista de required no JSONSchema
	requiredRaw, ok := planner.JSONSchema["required"]
	require.True(t, ok, "JSONSchema required field is missing")

	var required []string
	if interf, ok := requiredRaw.([]interface{}); ok {
		for _, v := range interf {
			if s, ok := v.(string); ok {
				required = append(required, s)
			}
		}
	} else if ss, ok := requiredRaw.([]string); ok {
		required = ss
	} else {
		require.Fail(t, "JSONSchema required field is not a string slice or interface slice")
	}

	assert.Contains(t, required, "backout_plan", "backout_plan should be a required field in propose_plan tool")
}
