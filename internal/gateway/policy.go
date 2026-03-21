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

type DryRunDecision struct {
	Lane       string         `json:"lane"`
	Provider   string         `json:"provider"`
	Model      string         `json:"model"`
	UseRemote  bool           `json:"use_remote"`
	UseTools   bool           `json:"use_tools"`
	Reason     string         `json:"reason"`
	Guards     ResponseGuards `json:"guards"`
	BudgetLane string         `json:"budget_lane"`
}

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

func (p *Planner) Decide(req DryRunRequest) DryRunDecision {
	taskClass := normalizeTaskClass(req)
	outputMode := strings.ToLower(strings.TrimSpace(req.OutputMode))
	useTools := req.RequiresTools || taskClass == "maintenance" || taskClass == "browser_workflow"

	switch taskClass {
	case "audio":
		return DryRunDecision{
			Lane:       "audio-stt",
			Provider:   "groq",
			Model:      defaultAudioModel,
			UseRemote:  true,
			UseTools:   false,
			Reason:     "audio/stt e um lane isolado; Groq tira o custo de STT da GPU local.",
			Guards:     ResponseGuards{ReasoningMode: "off", MaxOutputTokens: 64, SoftTimeoutMS: 20000},
			BudgetLane: "audio",
		}
	case "deep_research":
		return DryRunDecision{
			Lane:       "deep-research",
			Provider:   "gemini-web",
			Model:      defaultDeepResearchModel,
			UseRemote:  true,
			UseTools:   false,
			Reason:     "pesquisa profunda deve ficar fora do runtime critico e alimentar o RAG por curadoria.",
			Guards:     ResponseGuards{ReasoningMode: "default", MaxOutputTokens: 1200, SoftTimeoutMS: 120000},
			BudgetLane: "research",
		}
	}

	if req.RequiresVision && !req.LocalOnly {
		return DryRunDecision{
			Lane:       "remote-cheap-vision",
			Provider:   "openrouter",
			Model:      defaultRemoteVisionModel,
			UseRemote:  true,
			UseTools:   useTools,
			Reason:     "visao fica melhor no lane remoto barato e multimodal, sem puxar o lane premium.",
			Guards:     guardsFor(outputMode, true),
			BudgetLane: "remote_vision",
		}
	}

	if outputMode == "structured_json" || taskClass == "routing" || taskClass == "curation" {
		if req.LocalOnly {
			return DryRunDecision{
				Lane:       "local-balanced",
				Provider:   "ollama",
				Model:      defaultLocalBalancedModel,
				UseRemote:  false,
				UseTools:   useTools,
				Reason:     "saida estruturada local_only fica no 9b local, mas com guardas para conter reasoning.",
				Guards:     guardsFor(outputMode, false),
				BudgetLane: "local",
			}
		}
		return DryRunDecision{
			Lane:       "remote-tool-long-output",
			Provider:   "openrouter",
			Model:      defaultRemoteStructuredModel,
			UseRemote:  true,
			UseTools:   useTools,
			Reason:     "JSON curto e curadoria compacta ficaram mais previsiveis no DeepSeek do que em MiniMax/local com budget curto.",
			Guards:     guardsFor(outputMode, true),
			BudgetLane: "remote_structured",
		}
	}

	if req.PremiumRequired || taskClass == "browser_workflow" || taskClass == "workflow_premium" {
		return DryRunDecision{
			Lane:       "remote-premium-workflow",
			Provider:   "openrouter",
			Model:      defaultRemotePremiumModel,
			UseRemote:  true,
			UseTools:   useTools,
			Reason:     "workflow premium e browser humano ficam mais polidos no lane MiniMax.",
			Guards:     guardsFor(outputMode, true),
			BudgetLane: "remote_premium",
		}
	}

	if req.LatencyBudgetMS > 0 && req.LatencyBudgetMS <= 1500 && !useTools && !req.RequiresVision {
		return DryRunDecision{
			Lane:       "local-fast",
			Provider:   "ollama",
			Model:      defaultLocalFastModel,
			UseRemote:  false,
			UseTools:   false,
			Reason:     "triagem curta e sensivel a latencia cabe no lane local-fast.",
			Guards:     guardsFor(outputMode, false),
			BudgetLane: "local",
		}
	}

	if req.CostSensitive && !req.LocalOnly && !useTools {
		return DryRunDecision{
			Lane:       "remote-cheap-long-context",
			Provider:   "openrouter",
			Model:      defaultRemoteCheapLongModel,
			UseRemote:  true,
			UseTools:   false,
			Reason:     "tarefas remotas sensiveis a custo, sem tool use, cabem melhor no lane barato de contexto longo.",
			Guards:     guardsFor(outputMode, true),
			BudgetLane: "remote_cheap",
		}
	}

	return DryRunDecision{
		Lane:       "local-balanced",
		Provider:   "ollama",
		Model:      defaultLocalBalancedModel,
		UseRemote:  false,
		UseTools:   useTools,
		Reason:     "manutencao, repo e tool use local devem priorizar o cerebro local residente.",
		Guards:     guardsFor(outputMode, false),
		BudgetLane: "local",
	}
}

func normalizeTaskClass(req DryRunRequest) string {
	taskClass := strings.ToLower(strings.TrimSpace(req.TaskClass))
	if taskClass != "" {
		return taskClass
	}
	task := strings.ToLower(req.Task)
	switch {
	case strings.Contains(task, "audio") || strings.Contains(task, "stt") || strings.Contains(task, "transcri"):
		return "audio"
	case strings.Contains(task, "research") || strings.Contains(task, "pesquisa") || strings.Contains(task, "deep search"):
		return "deep_research"
	case strings.Contains(task, "screenshot") || strings.Contains(task, "imagem") || strings.Contains(task, "vision"):
		return "vision"
	case strings.Contains(task, "json") || strings.Contains(task, "roteamento"):
		return "routing"
	case strings.Contains(task, "curadoria") || strings.Contains(task, "facts") || strings.Contains(task, "tags"):
		return "curation"
	case strings.Contains(task, "browser") || strings.Contains(task, "ide"):
		return "browser_workflow"
	case strings.Contains(task, "maint") || strings.Contains(task, "homelab") || strings.Contains(task, "runbook"):
		return "maintenance"
	default:
		return "general"
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
		return ResponseGuards{
			ReasoningMode:   "minimize",
			MaxOutputTokens: 384,
			SoftTimeoutMS:   15000,
		}
	}
}
