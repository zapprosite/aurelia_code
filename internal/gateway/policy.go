package gateway

import (
	"fmt"
	"strings"
)

const (
	// Tier 1 — remote cheap
	modelQwen3       = "qwen/qwen3-32b"               // $0.08/$0.24 per 1M, PT-BR quality
	modelDevstral2   = "mistralai/devstral-2512"       // $0.05/$0.22, 256K ctx, coding+PT-BR
	modelQwen35Flash = "qwen/qwen3.5-flash-02-23"      // $0.065/$0.26, 1M ctx, ultra-fast

	// Tier 2 — remote premium
	modelMiniMaxM27    = "minimax/minimax-m2.7"
	modelMiniMaxM25    = "minimax/minimax-m2.5"        // slightly cheaper than M2.7
	modelMiniMaxDirect = "MiniMax-M2"

	// Long context — Llama 4 Scout: 10M ctx at $0.08/$0.30 (-85% vs Kimi)
	modelLlama4Scout = "meta-llama/llama-4-scout"
	modelKimiK25     = "moonshotai/kimi-k2.5"          // kept as vision fallback

	// Free reasoning (rate-limited: 1000 req/day with $10+ credit)
	modelDeepSeekR1Free = "deepseek/deepseek-r1-0528:free"

	// Local
	modelGemma3 = "gemma3:12b"

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

	// LocalProbe: if true, run quality gate on the response.
	// If quality passes → return immediately (no remote cost).
	// If quality fails → continue to next candidate silently (no failure recorded).
	LocalProbe bool `json:"local_probe,omitempty"`
}

// DryRunDecision mantém compatibilidade reversa onde testes/handlers precisarem
type DryRunDecision = RouteCandidate

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

// Decide selects the best candidate for the request, honoring constraints like LocalOnly.
func (p *Planner) Decide(req DryRunRequest) DryRunDecision {
	opts := p.Plan(req)
	if req.LocalOnly {
		for _, opt := range opts {
			if !opt.UseRemote {
				return opt
			}
		}
	}
	if len(opts) > 0 {
		return opts[0]
	}
	return DryRunDecision{}
}

func (p *Planner) Plan(req DryRunRequest) []RouteCandidate {
	// 1. Initial Classification
	taskClass := req.JudgeClass
	if taskClass == "" {
		if req.TaskClass != "" {
			taskClass = req.TaskClass
		} else {
			taskClass = classifyTask(req)
		}
	}

	// 2. Override Rules
	// If input contains image: force long_context_or_multimodal
	if req.RequiresVision {
		taskClass = "vision"
	}
	// If context is large: force long_context_or_multimodal
	if req.ContextSize > 12000 {
		taskClass = "long_context_or_multimodal"
	}
	// Low confidence = the judge isn't sure → treat as simple, route local (cheap).
	// Routing uncertain tasks to premium models wastes budget without quality gain.
	if req.JudgeConfidence > 0 && req.JudgeConfidence < 0.5 {
		taskClass = "simple_short"
	}

	// 3. Routing Policy Mapping
	candidates := make([]RouteCandidate, 0)
	isStructured := req.OutputMode == "structured_json"

	switch taskClass {
	case "curation", "simple_short", "general":
		// Local-first for cost sovereignty; remote only as fallback.
		candidates = []RouteCandidate{
			{
				Lane:       "local-balanced",
				Provider:   "local",
				Model:      modelGemma3,
				UseRemote:  false,
				Reason:     fmt.Sprintf("%s: gemma3 local first (cost sovereign).", taskClass),
				Guards:     guardsFor(req.OutputMode, false),
				BudgetLane: "local",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-cheap",
				Provider:   "openrouter",
				Model:      modelQwen3,
				UseRemote:  true,
				Reason:     fmt.Sprintf("%s: qwen3-32b remote fallback.", taskClass),
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "professional":
		// Business bots: try gemma3 locally first (LocalProbe).
		// If the local response meets quality bar → zero remote cost.
		// Only escalate to DeepSeek → MiniMax if quality gate fails.
		// In practice gemma3:12b handles ~65% of professional PT-BR queries adequately.
		candidates = []RouteCandidate{
			{
				Lane:       "local-balanced",
				Provider:   "local",
				Model:      modelGemma3,
				UseRemote:  false,
				LocalProbe: true, // quality gate — escalate silently if fails
				Reason:     "professional: gemma3 local probe first (free tier).",
				Guards:     guardsForProfessional(req.OutputMode, false),
				BudgetLane: "local",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-cheap",
				Provider:   "openrouter",
				Model:      modelQwen3,
				UseRemote:  true,
				Reason:     "professional: qwen3-32b — quality PT-BR business responses ($0.08/$0.24 per 1M).",
				Guards:     guardsForProfessional(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM25,
				UseRemote:  true,
				Reason:     "professional: minimax-m2.5 free tier (1000 req/day, falls through to m2.7).",
				Guards:     guardsForProfessional(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM27,
				UseRemote:  true,
				Reason:     "professional: minimax-m2.7 paid final fallback.",
				Guards:     guardsForProfessional(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "maintenance":
		candidates = []RouteCandidate{
			{
				Lane:       "local-balanced",
				Provider:   "local",
				Model:      modelGemma3,
				UseRemote:  false,
				Reason:     "maintenance: local-balanced preferred for homelab stability.",
				Guards:     guardsFor(req.OutputMode, false),
				BudgetLane: "local",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM27,
				UseRemote:  true,
				Reason:     "maintenance fallback: minimax-m2.7.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "coding_main", "routing":
		// Prefer DeepSeek for structured output even in coding_main
		if isStructured {
			candidates = append(candidates, RouteCandidate{
				Lane:       "remote-cheap",
				Provider:   "openrouter",
				Model:      modelQwen3,
				UseRemote:  true,
				Reason:     fmt.Sprintf("%s: qwen3-32b preferred for structured output.", taskClass),
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			})
		}

		candidates = append(candidates, RouteCandidate{
			Lane:       "remote-premium",
			Provider:   "minimax",
			Model:      modelMiniMaxM27,
			UseRemote:  true,
			Reason:     fmt.Sprintf("%s: minimax-m2.7 is the primary execution model.", taskClass),
			Guards:     guardsFor(req.OutputMode, true),
			BudgetLane: "remote_premium",
			Class:      taskClass,
			Confidence: req.JudgeConfidence,
		})

		if !isStructured {
			candidates = append(candidates, RouteCandidate{
				Lane:       "remote-cheap",
				Provider:   "openrouter",
				Model:      modelQwen3,
				UseRemote:  true,
				Reason:     fmt.Sprintf("%s: qwen3-32b fallback.", taskClass),
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			})
		}

		// Local fallback
		candidates = append(candidates, RouteCandidate{
			Lane:       "local-balanced",
			Provider:   "local",
			Model:      modelGemma3,
			UseRemote:  false,
			Reason:     fmt.Sprintf("%s: gemma3 local fallback.", taskClass),
			Guards:     guardsFor(req.OutputMode, false),
			BudgetLane: "local",
			Class:      taskClass,
			Confidence: req.JudgeConfidence,
		})

	case "long_context_or_multimodal", "vision":
		candidates = []RouteCandidate{
			{
				Lane:       "local-vision",
				Provider:   "local",
				Model:      modelGemma3,
				UseRemote:  false,
				Reason:     "long_context_or_multimodal: gemma3 12b local multimodal.",
				Guards:     guardsFor(req.OutputMode, false),
				BudgetLane: "local",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-long-context",
				Provider:   "openrouter",
				Model:      modelLlama4Scout,
				UseRemote:  true,
				Reason:     "long_context_or_multimodal: llama-4-scout 10M ctx ($0.08/$0.30, -85% vs kimi).",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_vision",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "critical":
		candidates = []RouteCandidate{
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM25,
				UseRemote:  true,
				Reason:     "critical: minimax-m2.5 free tier first.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-premium",
				Provider:   "minimax",
				Model:      modelMiniMaxM27,
				UseRemote:  true,
				Reason:     "critical: minimax-m2.7 paid fallback for high accuracy.",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_premium",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "local-balanced",
				Provider:   "local",
				Model:      modelGemma3,
				UseRemote:  false,
				Reason:     "critical: gemma3 local emergency fallback.",
				Guards:     guardsFor(req.OutputMode, false),
				BudgetLane: "local",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	default:
		candidates = []RouteCandidate{
			{
				Lane:      "remote-premium",
				Provider:  "minimax",
				Model:     modelMiniMaxM27,
				UseRemote: true,
				Reason:    "default: routing to coding_main.",
				Guards:    guardsFor(req.OutputMode, true),
				Class:     "coding_main",
			},
			{
				Lane:      "local-balanced",
				Provider:  "local",
				Model:     modelGemma3,
				UseRemote: false,
				Reason:    "default: local fallback.",
				Guards:    guardsFor(req.OutputMode, false),
				Class:     "coding_main",
			},
		}
	}

	for i := range candidates {
		candidates[i].UseTools = req.RequiresTools
	}
	return candidates
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
	{Class: "workflow_premium", Keywords: []string{"workflow complexo", "agentico", "multi-step", "complex workflow", "minimax", "poderoso"}, Weight: 10},
	{Class: "maintenance", Keywords: []string{"maint", "homelab", "runbook", "reboot", "nvidia-gpu", "servidor", "deploy"}, Weight: 5},
	// Professional: business bot responses for HVAC-R commercial and construction management.
	{Class: "professional", Keywords: []string{
		"proposta", "obra", "lead", "briefing", "contrato", "orcamento", "orçamento",
		"cliente", "vendas", "hvac", "vrf", "vrv", "split", "climatizacao", "climatização",
		"cronograma", "prazo", "entrega", "fornecedor", "fornecedores", "relatorio", "relatório",
		"pmoc", "art", "rrt", "comissionamento", "evaporadora", "condensadora",
		"follow-up", "fechamento", "especificacao", "especificação", "daikin",
		"obra em andamento", "nota fiscal", "medicao", "medição",
	}, Weight: 7},
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

// guardsForProfessional returns guards for professional business responses (proposals, reports, CRM).
// Uses larger token budgets to allow complete, structured responses.
func guardsForProfessional(outputMode string, remote bool) ResponseGuards {
	_ = remote // remote/local distinction doesn't change professional guards
	return ResponseGuards{
		ReasoningMode:   "default",
		MaxOutputTokens: 1024,
		SoftTimeoutMS:   45000,
	}
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
		// S-30: Polish Premium — expandindo buffers para formatação industrial.
		return ResponseGuards{
			ReasoningMode:   "default",
			MaxOutputTokens: 1024,
			SoftTimeoutMS:   35000,
		}
	}
}
