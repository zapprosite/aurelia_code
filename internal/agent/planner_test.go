package agent

import (
	"testing"
)

func TestProposePlanTool_RequiresBackoutPlan(t *testing.T) {
	reg := NewToolRegistry()
	reg.RegisterPlannerTools()

	var planner *Tool
	for _, d := range reg.GetDefinitions() {
		if d.Name == "propose_plan" {
			planner = &d
			break
		}
	}

	if planner == nil {
		t.Fatal("propose_plan tool not registered")
	}

	// Verificar se 'backout_plan' está na lista de required no JSONSchema
	requiredRaw, ok := planner.JSONSchema["required"]
	if !ok {
		t.Fatal("JSONSchema required field is missing")
	}

	required, ok := requiredRaw.([]string)
	if !ok {
		// Tentar []interface{} se []string falhar (comum em Go maps literais)
		if interf, ok := requiredRaw.([]interface{}); ok {
			for _, v := range interf {
				if s, ok := v.(string); ok {
					required = append(required, s)
				}
			}
		} else {
			t.Fatal("JSONSchema required field is not a string slice or interface slice")
		}
	}

	found := false
	for _, field := range required {
		if field == "backout_plan" {
			found = true
			break
		}
	}

	if !found {
		t.Error("backout_plan should be a required field in propose_plan tool")
	}
}
