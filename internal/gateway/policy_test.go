package gateway

import (
	"testing"
)

func TestPlannerDecide_MaintenancePrefersOpenRouterWithOllamaFallback(t *testing.T) {
	t.Parallel()

	opts := NewPlanner().Plan(DryRunRequest{
		TaskClass:     "maintenance",
		RequiresTools: true,
	})

	// Primary should be local-balanced for maintenance
	if len(opts) == 0 {
		t.Fatalf("expected at least one candidate")
	}
	if opts[0].Provider != "local" || opts[0].Model != modelGemma3 {
		t.Fatalf("unexpected primary route for maintenance: %+v", opts[0])
	}
	if !opts[0].UseTools {
		t.Fatalf("expected tools enabled: %+v", opts[0])
	}

	// Fallback should include remote-premium
	if len(opts) < 2 || opts[1].Provider != "minimax" || opts[1].Model != modelMiniMaxM27 {
		t.Fatalf("expected minimax fallback for maintenance, got: %+v", opts)
	}
}

func TestPlannerDecide_StructuredUsesDeepSeek(t *testing.T) {
	t.Parallel()

	got := NewPlanner().Decide(DryRunRequest{
		TaskClass:  "routing",
		OutputMode: "structured_json",
	})

	if got.Provider != "openrouter" || got.Model != modelQwen3 {
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

	if got.Lane != "local-vision" || got.Model != modelGemma3 {
		t.Fatalf("unexpected route: %+v", got)
	}
}

func TestPlannerDecide_AudioUsesGroq(t *testing.T) {
	t.Parallel()

	got := NewPlanner().Decide(DryRunRequest{TaskClass: "audio"})

	if got.Provider != "minimax" || (got.Model != modelMiniMaxM27 && got.Model != modelMiniMaxDirect) {
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

	if got.Provider != "local" || got.Model != modelGemma3 {
		t.Fatalf("unexpected route: %+v", got)
	}
}

func TestClassifyTask_TabularCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       DryRunRequest
		wantClass string
	}{
		{name: "explicit class wins", req: DryRunRequest{TaskClass: "audio", Task: "research this"}, wantClass: "audio"},
		{name: "audio by keyword stt", req: DryRunRequest{Task: "transcribe this audio via stt"}, wantClass: "audio"},
		{name: "audio by keyword transcri", req: DryRunRequest{Task: "transcrição do áudio"}, wantClass: "audio"},
		{name: "audio by keyword whisper", req: DryRunRequest{Task: "use whisper model"}, wantClass: "audio"},
		{name: "deep research explicit", req: DryRunRequest{Task: "do deep research on this topic"}, wantClass: "deep_research"},
		{name: "deep research pt", req: DryRunRequest{Task: "pesquisa profunda sobre o assunto"}, wantClass: "deep_research"},
		{name: "vision with flag", req: DryRunRequest{Task: "analyze this screenshot", RequiresVision: true}, wantClass: "vision"},
		{name: "vision without flag stays general", req: DryRunRequest{Task: "analyze this screenshot"}, wantClass: "general"},
		{name: "routing keyword", req: DryRunRequest{Task: "roteamento da mensagem"}, wantClass: "routing"},
		{name: "routing classifier", req: DryRunRequest{Task: "run the classifier on this input"}, wantClass: "routing"},
		{name: "curation facts", req: DryRunRequest{Task: "extract the key facts from this"}, wantClass: "curation"},
		{name: "curation tags", req: DryRunRequest{Task: "generate tags for this content"}, wantClass: "curation"},
		{name: "browser workflow", req: DryRunRequest{Task: "open the browser and navigate to this page"}, wantClass: "browser_workflow"},
		{name: "maintenance reboot", req: DryRunRequest{Task: "reboot the server"}, wantClass: "maintenance"},
		{name: "maintenance homelab", req: DryRunRequest{Task: "homelab setup needs attention"}, wantClass: "maintenance"},
		{name: "maintenance nvidia", req: DryRunRequest{Task: "check nvidia-gpu status"}, wantClass: "maintenance"},
		{name: "general fallback", req: DryRunRequest{Task: "tell me about the weather"}, wantClass: "general"},
		{name: "empty task", req: DryRunRequest{}, wantClass: "general"},
		{name: "ambiguous prefers higher weight", req: DryRunRequest{Task: "research this route classifier"}, wantClass: "routing"},
		{name: "workflow premium", req: DryRunRequest{Task: "run a complex workflow with multi-step processing"}, wantClass: "workflow_premium"},
		{name: "deploy as maintenance", req: DryRunRequest{Task: "deploy the new version to the servidor"}, wantClass: "maintenance"},
		{name: "case insensitive", req: DryRunRequest{Task: "TRANSCRIBE THIS AUDIO"}, wantClass: "audio"},
		{name: "explicit class overrides keywords", req: DryRunRequest{TaskClass: "maintenance", Task: "transcribe audio"}, wantClass: "maintenance"},
		{name: "gravacao as audio", req: DryRunRequest{Task: "gravacao de voz"}, wantClass: "audio"},
		{name: "curate keyword", req: DryRunRequest{Task: "curate content for the feed"}, wantClass: "curation"},
		{name: "navegar as browser", req: DryRunRequest{Task: "navegar ate a pagina principal"}, wantClass: "browser_workflow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := classifyTask(tt.req)
			if got != tt.wantClass {
				t.Errorf("classifyTask(%+v) = %q, want %q", tt.req, got, tt.wantClass)
			}
		})
	}
}

func TestIsHardProfessional(t *testing.T) {
	t.Parallel()

	soft := []string{
		"O que significa HVAC?",
		"Explique diferenca entre VRF e split para cliente leigo",
		"Quais pontos de negociacao numa proposta VRV Daikin?",
		"Como funciona um sistema de climatizacao?",
	}
	hard := []string{
		"Script de vendas com ROI e payback para VRV Daikin incorporadora",
		"Pontos por norma ABNT para comissionamento VRV edificio comercial",
		"Comparacao VRV vs chiller+fancoil 2000m2 analise COP custo operacional",
		"Plano manutencao cronograma anual 20 unidades internas Daikin",
		"Calcule carga termica e COP para edificio 500m2",
		"Proposta tecnica completa com payback e ROI",
	}

	for _, q := range soft {
		if isHardProfessional(q) {
			t.Errorf("expected soft for %q, got hard", q)
		}
	}
	for _, q := range hard {
		if !isHardProfessional(q) {
			t.Errorf("expected hard for %q, got soft", q)
		}
	}
}

func TestProfessionalHardRoutesToGroq(t *testing.T) {
	t.Parallel()

	hardQueries := []string{
		"Script de vendas com ROI e payback para VRV Daikin",
		"Pontos ABNT para comissionamento VRV edificio comercial",
		"Analise VRV vs chiller+fancoil analise COP custo",
	}

	p := NewPlanner()
	for _, q := range hardQueries {
		candidates := p.Plan(DryRunRequest{
			Task: q, JudgeClass: "professional", JudgeConfidence: 0.85,
		})
		if len(candidates) == 0 {
			t.Fatalf("no candidates for %q", q)
		}
		first := candidates[0]
		if first.Lane != "remote-free" || first.Provider != "groq" {
			t.Errorf("hard query should route to groq remote-free first, got lane=%s provider=%s for %q",
				first.Lane, first.Provider, q)
		}
	}
}

func TestProfessionalSoftRoutesToLocalProbe(t *testing.T) {
	t.Parallel()

	softQueries := []string{
		"Explique diferenca entre VRF e split para cliente leigo",
		"O que e VRV Daikin?",
	}

	p := NewPlanner()
	for _, q := range softQueries {
		candidates := p.Plan(DryRunRequest{
			Task: q, JudgeClass: "professional", JudgeConfidence: 0.85,
		})
		if len(candidates) == 0 {
			t.Fatalf("no candidates for %q", q)
		}
		first := candidates[0]
		if first.Lane != "local-balanced" || !first.LocalProbe {
			t.Errorf("soft query should route to local probe first, got lane=%s probe=%v for %q",
				first.Lane, first.LocalProbe, q)
		}
	}
}
