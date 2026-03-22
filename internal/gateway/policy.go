package gateway

import "strings"

const (
	defaultLocalFastModel        = "qwen3.5:4b"
	defaultLocalBalancedModel    = "qwen3.5:9b"
	defaultRemoteVisionModel     = "qwen/qwen3.5-9b"
	defaultRemoteCheapLongModel  = "qwen/qwen3.5-flash-02-23"
	defaultRemoteStructuredModel = "google/gemini-2.5-flash"
	defaultRemotePremiumModel    = "minimax/minimax-m2.7"
	defaultAudioModel            = "whisper-large-v3-turbo"
	defaultDeepResearchModel     = "gemini-web-deep-research"
)

type DryRunRequest struct {
	Task            string `json:"task"`
	TaskClass       string `json:"task_class,omitempty"`
	OutputMode      string `json:"output_mode,omitempty"`
	RequiresTools   bool   `json:"requires_tools,omitempty"`
	RequiresVision  bool   `json:"requires_vision,omitempty"`
	LocalOnly       bool   `json:"local_only,omitempty"`
	CostSensitive   bool   `json:"cost_sensitive,omitempty"`
	PremiumRequired bool   `json:"premium_required,omitempty"`
	LatencyBudgetMS int    `json:"latency_budget_ms,omitempty"`
}

type ResponseGuards struct {
	ReasoningMode   string `json:"reasoning_mode"`
	MaxOutputTokens int    `json:"max_output_tokens"`
	SoftTimeoutMS   int    `json:"soft_timeout_ms"`
}

type RouteCandidate struct {
	Lane       string         `json:"lane"`
	Provider   string         `json:"provider"`
	Model      string         `json:"model"`
	UseRemote  bool           `json:"use_remote"`
	UseTools   bool           `json:"use_tools"`
	Reason     string         `json:"reason"`
	Guards     ResponseGuards `json:"guards"`
	BudgetLane string         `json:"budget_lane"`
}

// DryRunDecision mantém compatibilidade reversa onde testes/handlers precisarem
type DryRunDecision = RouteCandidate

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

// Decide maintains compatibility for tests/handlers that expect a single primary candidate
func (p *Planner) Decide(req DryRunRequest) DryRunDecision {
	opts := p.Plan(req)
	if len(opts) > 0 {
		return opts[0]
	}
	return DryRunDecision{} // Fallback if somehow empty
}

func (p *Planner) Plan(req DryRunRequest) []RouteCandidate {
	taskClass := classifyTask(req)
	outputMode := strings.ToLower(strings.TrimSpace(req.OutputMode))
	useTools := req.RequiresTools || taskClass == "maintenance" || taskClass == "browser_workflow"

	// Audio
	if taskClass == "audio" {
		return []RouteCandidate{{
			Lane:       "audio-stt",
			Provider:   "groq",
			Model:      defaultAudioModel,
			UseRemote:  true,
			UseTools:   false,
			Reason:     "audio/stt e um lane isolado; Groq tira o custo de STT da GPU local.",
			Guards:     ResponseGuards{ReasoningMode: "off", MaxOutputTokens: 64, SoftTimeoutMS: 20000},
			BudgetLane: "audio",
		}}
	}

	// Deep Research
	if taskClass == "deep_research" {
		return []RouteCandidate{{
			Lane:       "deep-research",
			Provider:   "gemini-web",
			Model:      defaultDeepResearchModel,
			UseRemote:  true,
			UseTools:   false,
			Reason:     "pesquisa profunda deve ficar fora do runtime critico e alimentar o RAG por curadoria.",
			Guards:     ResponseGuards{ReasoningMode: "default", MaxOutputTokens: 1200, SoftTimeoutMS: 120000},
			BudgetLane: "research",
		}}
	}

	// Vision
	if req.RequiresVision && !req.LocalOnly {
		return []RouteCandidate{{
			Lane:       "remote-cheap-vision",
			Provider:   "openrouter",
			Model:      defaultRemoteVisionModel,
			UseRemote:  true,
			UseTools:   useTools,
			Reason:     "visao fica melhor no lane remoto barato e multimodal.",
			Guards:     guardsFor(outputMode, true),
			BudgetLane: "remote_vision",
		}}
	}

	// Structured, Routing, Curation
	if outputMode == "structured_json" || taskClass == "routing" || taskClass == "curation" {
		if req.LocalOnly {
			return []RouteCandidate{{
				Lane:       "local-balanced",
				Provider:   "ollama",
				Model:      defaultLocalBalancedModel,
				UseRemote:  false,
				UseTools:   useTools,
				Reason:     "saida estruturada local_only fica no 9b local, mas com guardas para conter reasoning.",
				Guards:     guardsFor(outputMode, false),
				BudgetLane: "local",
			}}
		}
		return []RouteCandidate{
			{
				Lane:       "remote-tool-long-output",
				Provider:   "openrouter",
				Model:      defaultRemoteStructuredModel,
				UseRemote:  true,
				UseTools:   useTools,
				Reason:     "primary: JSON estruturado preferencial no gemini 2.5 flash.",
				Guards:     guardsFor(outputMode, true),
				BudgetLane: "remote_structured",
			},
			{
				Lane:       "remote-cheap-long-context",
				Provider:   "openrouter",
				Model:      defaultRemoteCheapLongModel,
				UseRemote:  true,
				UseTools:   useTools,
				Reason:     "fallback 1: modelo barato alternativo na nuvem.",
				Guards:     guardsFor(outputMode, true),
				BudgetLane: "remote_cheap",
			},
			{
				Lane:       "local-balanced",
				Provider:   "ollama",
				Model:      defaultLocalBalancedModel,
				UseRemote:  false,
				UseTools:   useTools,
				Reason:     "fallback 2: ollama local residente.",
				Guards:     guardsFor(outputMode, false),
				BudgetLane: "local",
			},
		}
	}

	// Premium / Browser
	if req.PremiumRequired || taskClass == "browser_workflow" || taskClass == "workflow_premium" {
		if req.LocalOnly {
			return []RouteCandidate{{
				Lane:       "local-balanced",
				Provider:   "ollama",
				Model:      defaultLocalBalancedModel,
				UseRemote:  false,
				UseTools:   useTools,
				Reason:     "premium req localmente forçado pro 9b.",
				Guards:     guardsFor(outputMode, false),
				BudgetLane: "local",
			}}
		}
		return []RouteCandidate{
			{
				Lane:       "remote-premium-workflow",
				Provider:   "openrouter",
				Model:      defaultRemotePremiumModel,
				UseRemote:  true,
				UseTools:   useTools,
				Reason:     "primary: workflow premium e browser humano ficam polidos no lane MiniMax.",
				Guards:     guardsFor(outputMode, true),
				BudgetLane: "remote_premium",
			},
			{
				Lane:       "remote-tool-long-output",
				Provider:   "openrouter",
				Model:      defaultRemoteStructuredModel,
				UseRemote:  true,
				UseTools:   useTools,
				Reason:     "fallback 1: gemini flash estruturado.",
				Guards:     guardsFor(outputMode, true),
				BudgetLane: "remote_structured",
			},
			{
				Lane:       "local-balanced",
				Provider:   "ollama",
				Model:      defaultLocalBalancedModel,
				UseRemote:  false,
				UseTools:   useTools,
				Reason:     "fallback 2: ollama local residente.",
				Guards:     guardsFor(outputMode, false),
				BudgetLane: "local",
			},
		}
	}

	// Latency-sensitive local
	if req.LatencyBudgetMS > 0 && req.LatencyBudgetMS <= 1500 && !useTools && !req.RequiresVision {
		return []RouteCandidate{{
			Lane:       "local-fast",
			Provider:   "ollama",
			Model:      defaultLocalFastModel,
			UseRemote:  false,
			UseTools:   false,
			Reason:     "triagem curta e sensivel a latencia cabe no lane local-fast.",
			Guards:     guardsFor(outputMode, false),
			BudgetLane: "local",
		}}
	}

	// Cost Sensitive
	if req.CostSensitive && !req.LocalOnly && !useTools {
		return []RouteCandidate{
			{
				Lane:       "remote-cheap-long-context",
				Provider:   "openrouter",
				Model:      defaultRemoteCheapLongModel,
				UseRemote:  true,
				UseTools:   false,
				Reason:     "primary: tarefas remotas sensiveis a custo.",
				Guards:     guardsFor(outputMode, true),
				BudgetLane: "remote_cheap",
			},
			{
				Lane:       "local-balanced",
				Provider:   "ollama",
				Model:      defaultLocalBalancedModel,
				UseRemote:  false,
				UseTools:   false,
				Reason:     "fallback: ollama local.",
				Guards:     guardsFor(outputMode, false),
				BudgetLane: "local",
			},
		}
	}

	// Maintenance, Local Only, General Default
	return []RouteCandidate{{
		Lane:       "local-balanced",
		Provider:   "ollama",
		Model:      defaultLocalBalancedModel,
		UseRemote:  false,
		UseTools:   useTools,
		Reason:     "manutencao, repo e tool use local priorizam cerebro local residente.",
		Guards:     guardsFor(outputMode, false),
		BudgetLane: "local",
	}}
}

// classifyRule maps keyword patterns to a task class with a weight.
type classifyRule struct {
	Class         string
	Keywords      []string
	Weight        int
	RequiresField string // "vision", or "" for none
}

var classifyRules = []classifyRule{
	{Class: "audio", Keywords: []string{"audio", "stt", "transcri", "whisper", "gravacao", "microfone"}, Weight: 10},
	{Class: "deep_research", Keywords: []string{"deep research", "pesquisa profunda", "deep search", "investigar a fundo"}, Weight: 10},
	{Class: "vision", Keywords: []string{"screenshot", "imagem", "image", "visual", "foto", "captura de tela"}, Weight: 8, RequiresField: "vision"},
	{Class: "routing", Keywords: []string{"roteamento", "classifier", "classify", "categorize", "categorizar", "route"}, Weight: 6},
	{Class: "curation", Keywords: []string{"curadoria", "facts", "tags", "resumo curto", "fatos curtos", "bullet points", "curate"}, Weight: 6},
	{Class: "browser_workflow", Keywords: []string{"browser", "navigate", "navegar", "web page", "pagina web"}, Weight: 7},
	{Class: "workflow_premium", Keywords: []string{"workflow complexo", "agentico", "multi-step", "complex workflow"}, Weight: 7},
	{Class: "maintenance", Keywords: []string{"maint", "homelab", "runbook", "reboot", "nvidia-gpu", "servidor", "deploy"}, Weight: 5},
	{Class: "general", Keywords: []string{}, Weight: 0},
}

// classifyTask uses weighted keyword scoring to determine task class.
// Explicit TaskClass in the request takes precedence.
func classifyTask(req DryRunRequest) string {
	taskClass := strings.ToLower(strings.TrimSpace(req.TaskClass))
	if taskClass != "" {
		return taskClass
	}

	task := strings.ToLower(req.Task)
	if task == "" {
		return "general"
	}

	bestClass := "general"
	bestScore := 0

	for _, rule := range classifyRules {
		if rule.RequiresField == "vision" && !req.RequiresVision {
			continue
		}
		score := 0
		for _, kw := range rule.Keywords {
			if strings.Contains(task, kw) {
				score += rule.Weight
			}
		}
		if score > bestScore {
			bestScore = score
			bestClass = rule.Class
		}
	}
	return bestClass
}

func guardsFor(outputMode string, remote bool) ResponseGuards {
	switch outputMode {
	case "structured_json", "curation":
		return ResponseGuards{
			ReasoningMode:   "minimize",
			MaxOutputTokens: 256,
			SoftTimeoutMS:   12000,
		}
	default:
		if remote {
			return ResponseGuards{
				ReasoningMode:   "default",
				MaxOutputTokens: 512,
				SoftTimeoutMS:   30000,
			}
		}
		return ResponseGuards{
			ReasoningMode:   "minimize",
			MaxOutputTokens: 384,
			SoftTimeoutMS:   15000,
		}
	}
}
