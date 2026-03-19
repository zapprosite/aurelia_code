package gateway

import "testing"

func TestPlannerDecide_MaintenancePrefersLocalBalanced(t *testing.T) {
	t.Parallel()

	got := NewPlanner().Decide(DryRunRequest{
		TaskClass:     "maintenance",
		RequiresTools: true,
	})

	if got.Provider != "ollama" || got.Model != defaultLocalBalancedModel {
		t.Fatalf("unexpected route: %+v", got)
	}
	if !got.UseTools {
		t.Fatalf("expected tools enabled: %+v", got)
	}
}

func TestPlannerDecide_StructuredUsesDeepSeek(t *testing.T) {
	t.Parallel()

	got := NewPlanner().Decide(DryRunRequest{
		TaskClass:  "routing",
		OutputMode: "structured_json",
	})

	if got.Provider != "openrouter" || got.Model != defaultRemoteStructuredModel {
		t.Fatalf("unexpected route: %+v", got)
	}
	if got.Guards.ReasoningMode != "minimize" {
		t.Fatalf("unexpected guards: %+v", got.Guards)
	}
}

func TestPlannerDecide_VisionUsesRemoteVisionLane(t *testing.T) {
	t.Parallel()

	got := NewPlanner().Decide(DryRunRequest{
		Task:           "ler screenshot da IDE e resumir erro",
		RequiresVision: true,
	})

	if got.Lane != "remote-cheap-vision" || got.Model != defaultRemoteVisionModel {
		t.Fatalf("unexpected route: %+v", got)
	}
}

func TestPlannerDecide_AudioUsesGroq(t *testing.T) {
	t.Parallel()

	got := NewPlanner().Decide(DryRunRequest{TaskClass: "audio"})

	if got.Provider != "groq" || got.Model != defaultAudioModel {
		t.Fatalf("unexpected route: %+v", got)
	}
}

func TestPlannerDecide_LocalOnlyStructuredStaysLocal(t *testing.T) {
	t.Parallel()

	got := NewPlanner().Decide(DryRunRequest{
		TaskClass:  "curation",
		OutputMode: "structured_json",
		LocalOnly:  true,
	})

	if got.Provider != "ollama" || got.Model != defaultLocalBalancedModel {
		t.Fatalf("unexpected route: %+v", got)
	}
}
