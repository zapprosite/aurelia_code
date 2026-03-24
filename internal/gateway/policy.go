package gateway

import "strings"

const (
	// Models requested by the user
	modelDeepSeekChat = "deepseek/deepseek-chat-v3.1"
	modelQwenNext     = "qwen/qwen3-coder-next"
	modelMiniMaxM27   = "minimax-m2.7-direct"
	modelKimiK25      = "moonshotai/kimi-k2.5"

	modelGemma3 = "gemma3:12b"
	modelLlama3 = "llama3.2:3b"
	modelGroq70b = "llama-3.3-70b-versatile"

	// Existing static defaults
	defaultAudioModel        = "whisper-large-v3-turbo"
	defaultDeepResearchModel = "gemini-web-deep-research"
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
	
	// New fields for the Judge
	JudgeClass      string  `json:"judge_class,omitempty"`
	JudgeConfidence float64 `json:"judge_confidence,omitempty"`
	ContextSize     int     `json:"context_size,omitempty"`
	RetryCount      int     `json:"retry_count,omitempty"`
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
	
	// Metadata from the judge
	Class      string  `json:"class,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
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
	// 1. Initial Classification
	taskClass := req.JudgeClass
	if taskClass == "" {
		taskClass = "coding_main" // Default
	}

	// 2. Override Rules
	// If input contains image / screenshot / PDF: force long_context_or_multimodal
	if req.RequiresVision {
		taskClass = "long_context_or_multimodal"
	}
	// If context is large: force long_context_or_multimodal
	if req.ContextSize > 12000 { // Threshold for "large" context
		taskClass = "long_context_or_multimodal"
	}
	// If judge confidence < 0.6: route to coding_main
	if req.JudgeConfidence > 0 && req.JudgeConfidence < 0.6 {
		taskClass = "coding_main"
	}

	// 3. Routing Policy Mapping
	switch taskClass {
	case "simple_short":
		return []RouteCandidate{
			{
				Lane:       "remote-cheap",
				Provider:   "openrouter",
				Model:      modelDeepSeekChat,
				UseRemote:  true,
				Reason:     "simple_short: deepseek-chat is cost-efficient.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-cheap-fallback",
				Provider:   "openrouter",
				Model:      modelQwenNext,
				UseRemote:  true,
				Reason:     "simple_short: qwen fallback.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "coding_main":
		return []RouteCandidate{
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM27,
				UseRemote:  true,
				Reason:     "coding_main: minimax-m2.7 is the primary execution model.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-cheap-fallback",
				Provider:   "openrouter",
				Model:      modelQwenNext,
				UseRemote:  true,
				Reason:     "coding_main: qwen fallback.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "long_context_or_multimodal":
		return []RouteCandidate{
			{
				Lane:       "remote-long-context",
				Provider:   "openrouter",
				Model:      modelKimiK25,
				UseRemote:  true,
				Reason:     "long_context_or_multimodal: kimi is specialized for long context/vision.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_vision",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM27,
				UseRemote:  true,
				Reason:     "long_context_or_multimodal: minimax fallback.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "critical":
		return []RouteCandidate{
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM27,
				UseRemote:  true,
				Reason:     "critical: minimax-m2.7 for high accuracy.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-long-context",
				Provider:   "openrouter",
				Model:      modelKimiK25,
				UseRemote:  true,
				Reason:     "critical: kimi fallback.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_vision",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	default:
		// Default to coding_main
		return []RouteCandidate{
			{
				Lane:      "remote-premium",
				Provider:  "minimax",
				Model:     modelMiniMaxM27,
				UseRemote: true,
				Reason:    "default: routing to coding_main.",
				Guards:    guardsFor(req.OutputMode, true),
				Class:     "coding_main",
			},
		}
	}
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
	{Class: "workflow_premium", Keywords: []string{"workflow complexo", "agentico", "multi-step", "complex workflow", "minimax", "poderoso", "openai", "claude"}, Weight: 10},
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
		// Para modelos locais (predominantemente gemma3:12b), 
		// permitimos o raciocínio por padrão para evitar respostas vazias.
		return ResponseGuards{
			ReasoningMode:   "default",
			MaxOutputTokens: 512,
			SoftTimeoutMS:   25000,
		}
	}
}
