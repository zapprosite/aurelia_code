package gateway

import (
	"fmt"
	"strings"
)

const (
	// Tier 0.5 — free remote (1000 req/day with $10+ OpenRouter credit)
	modelMiniMaxM25Free = "minimax/minimax-m2.5:free"

	// Tier 1 — remote cheap (paid)
	modelDeepSeekV31 = "deepseek/deepseek-chat-v3.1" // Tier 1 cheap remote standard
	modelDevstral2   = "mistralai/devstral-2512"     // optional cheap coding fallback

	// Tier 2 — remote premium (paid)
	modelMiniMaxM27    = "minimax/minimax-m2.7"
	modelMiniMaxM25    = "minimax/minimax-m2.5" // slightly cheaper than M2.7
	modelMiniMaxDirect = "MiniMax-M2"

	// Long context — Llama 4 Scout: 10M ctx at $0.08/$0.30 (-85% vs Kimi)
	modelLlama4Scout = "meta-llama/llama-4-scout"
	modelKimiK25     = "moonshotai/kimi-k2.5" // kept as vision fallback

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
				Model:      modelDeepSeekV31,
				UseRemote:  true,
				Reason:     fmt.Sprintf("%s: deepseek-chat-v3.1 remote fallback.", taskClass),
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_cheap",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
		}

	case "professional":
		// Hard technical queries (ABNT, COP, ROI, multi-step calculations) bypass local probe
		// and go directly to Groq (free, 0.8s) for accuracy. Soft queries still probe locally.
		//
		// Cascade:
		//   hard  → Groq (free, 0.8s) → DeepSeek V3.1 → M2.7 ($0.30/$1.20)
		//   soft  → gemma3 probe (free, 2s) → Groq (free) → DeepSeek V3.1 → M2.7
		if isHardProfessional(req.Task) {
			candidates = []RouteCandidate{
				{
					Lane:       "remote-free",
					Provider:   "groq",
					Model:      "llama-3.3-70b-versatile",
					UseRemote:  true,
					Reason:     "professional[hard]: groq llama-3.3 free tier — fast, accurate for technical.",
					Guards:     guardsForProfessional(req.OutputMode, true),
					BudgetLane: "remote_free",
					Class:      taskClass,
					Confidence: req.JudgeConfidence,
				},
				{
					Lane:       "remote-cheap",
					Provider:   "openrouter",
					Model:      modelDeepSeekV31,
					UseRemote:  true,
					Reason:     "professional[hard]: deepseek-chat-v3.1 paid fallback.",
					Guards:     guardsForProfessional(req.OutputMode, true),
					BudgetLane: "remote_cheap",
					Class:      taskClass,
					Confidence: req.JudgeConfidence,
				},
				{
					Lane:       "remote-premium",
					Provider:   "openrouter",
					Model:      modelMiniMaxM27,
					UseRemote:  true,
					Reason:     "professional[hard]: minimax-m2.7 premium ceiling ($0.30/$1.20/1M).",
					Guards:     guardsForProfessional(req.OutputMode, true),
					BudgetLane: "remote_premium",
					Class:      taskClass,
					Confidence: req.JudgeConfidence,
				},
			}
		} else {
			candidates = []RouteCandidate{
				{
					Lane:       "local-balanced",
					Provider:   "local",
					Model:      modelGemma3,
					UseRemote:  false,
					LocalProbe: true,
					Reason:     "professional[soft]: gemma3 probe — free, zero cost if ≥80w confident.",
					Guards:     guardsForProfessional(req.OutputMode, false),
					BudgetLane: "local",
					Class:      taskClass,
					Confidence: req.JudgeConfidence,
				},
				{
					Lane:       "remote-free",
					Provider:   "groq",
					Model:      "llama-3.3-70b-versatile",
					UseRemote:  true,
					Reason:     "professional[soft]: groq llama-3.3 free tier fallback (0.8s).",
					Guards:     guardsForProfessional(req.OutputMode, true),
					BudgetLane: "remote_free",
					Class:      taskClass,
					Confidence: req.JudgeConfidence,
				},
				{
					Lane:       "remote-cheap",
					Provider:   "openrouter",
					Model:      modelDeepSeekV31,
					UseRemote:  true,
					Reason:     "professional[soft]: deepseek-chat-v3.1 paid fallback.",
					Guards:     guardsForProfessional(req.OutputMode, true),
					BudgetLane: "remote_cheap",
					Class:      taskClass,
					Confidence: req.JudgeConfidence,
				},
				{
					Lane:       "remote-premium",
					Provider:   "openrouter",
					Model:      modelMiniMaxM27,
					UseRemote:  true,
					Reason:     "professional[soft]: minimax-m2.7 premium ceiling ($0.30/$1.20/1M).",
					Guards:     guardsForProfessional(req.OutputMode, true),
					BudgetLane: "remote_premium",
					Class:      taskClass,
					Confidence: req.JudgeConfidence,
				},
			}
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
				Model:      modelDeepSeekV31,
				UseRemote:  true,
				Reason:     fmt.Sprintf("%s: deepseek-chat-v3.1 preferred for structured output.", taskClass),
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
				Model:      modelDeepSeekV31,
				UseRemote:  true,
				Reason:     fmt.Sprintf("%s: deepseek-chat-v3.1 fallback.", taskClass),
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
				Lane:       "remote-free",
				Provider:   "groq",
				Model:      "llama-3.3-70b-versatile",
				UseRemote:  true,
				Reason:     "critical: groq llama-3.3 free tier (fast, accurate).",
				Guards:     guardsFor(req.OutputMode, true),
				BudgetLane: "remote_free",
				Class:      taskClass,
				Confidence: req.JudgeConfidence,
			},
			{
				Lane:       "remote-premium",
				Provider:   "openrouter",
				Model:      modelMiniMaxM27,
				UseRemote:  true,
				Reason:     "critical: minimax-m2.7 premium ceiling ($0.30/$1.20/1M).",
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

// isHardProfessional returns true when a professional query contains signals that indicate
// high technical complexity: specific norms, calculations, multi-criteria analysis, or
// planning tasks. These bypass the local probe and go directly to Groq for accuracy.
func isHardProfessional(task string) bool {
	t := strings.ToLower(task)
	hardSignals := []string{
		// Technical standards and certifications
		"abnt", "nbr", "norma", "nr-", "iso ", "pcld",
		// Engineering metrics
		"cop ", " cop\t", "kw", "btu", "eer", "capacidade frigorifica", "carga termica", "carga térmica",
		// Financial analysis
		"roi", "payback", "retorno", "tir ", " tir\t", "vpl", "fluxo de caixa",
		// Complex planning
		"cronograma", "matriz de", "plano de manutencao", "plano de manutenção",
		"comissionamento", "comissionar",
		// Multi-factor comparison
		"comparar", "comparacao", "comparação", "vs ", " versus ", "analise", "análise",
		// Specific technical documents
		"script de vendas", "proposta tecnica", "proposta técnica", "memorial descritivo",
		// Multi-unit / scale indicators
		"unidades internas", "unidades externas",
	}
	hits := 0
	for _, sig := range hardSignals {
		if strings.Contains(t, sig) {
			hits++
			if hits >= 2 {
				return true
			}
		}
	}
	// Single high-weight signal is enough
	highWeight := []string{"abnt", "nbr", "comissionamento", "payback", "roi", "script de vendas"}
	for _, sig := range highWeight {
		if strings.Contains(t, sig) {
			return true
		}
	}
	return false
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
