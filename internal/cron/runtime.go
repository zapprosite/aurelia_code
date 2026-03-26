package cron

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/observability"
)

type AgentCronRuntime struct {
	executor      AgentExecutor
	basePrompt    string
	allowedTools  []string
	promptBuilder func(ctx context.Context, job CronJob) (string, []string, error)
}

func NewAgentCronRuntime(executor AgentExecutor, baseSystemPrompt string, allowedTools []string) *AgentCronRuntime {
	return &AgentCronRuntime{
		executor:     executor,
		basePrompt:   baseSystemPrompt,
		allowedTools: allowedTools,
	}
}

func NewAgentCronRuntimeWithPromptBuilder(executor AgentExecutor, baseSystemPrompt string, allowedTools []string, promptBuilder func(ctx context.Context, job CronJob) (string, []string, error)) *AgentCronRuntime {
	return &AgentCronRuntime{
		executor:      executor,
		basePrompt:    baseSystemPrompt,
		allowedTools:  allowedTools,
		promptBuilder: promptBuilder,
	}
}

var grafanaDashboardURL = "https://monitor.zappro.site/d/40977c27-62f1-4a7d-a455-2f4d698538d8/nvidia-gpu-metrics"

func (r *AgentCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (string, []agent.ContentPart, error) {
	logger := observability.Logger("cron.runtime").With(slog.String("job", job.ID))

	// 1. Extração de métricas reais (Hard Data) para auxílio da análise
	metrics, _ := r.extractGPUMetrics(ctx)

	// Se for o watchdog de 5 min, só reportamos se houver anomalia (Temp > 75 ou Container Down)
	if cronMarker(job.Prompt) == "[sys:sentinel-watchdog]" {
		return r.handleSentinelWatchdog(ctx, job, metrics)
	}

	ctx = agent.WithTeamContext(ctx, formatChatTeamKey(job.TargetChatID), job.OwnerUserID)

	systemPrompt := r.basePrompt
	allowedTools := r.allowedTools
	if r.promptBuilder != nil {
		prompt, tools, err := r.promptBuilder(ctx, job)
		if err == nil {
			systemPrompt = prompt
			allowedTools = tools
		}
	}

	analysisPrompt := fmt.Sprintf(
		"%s\n\n[SISTEMA: MÉTRICAS REAIS DETECTADAS]\n%s\n\nEscreva um resumo operacional curto, técnico e direto. Não se apresente como Aurélia nem como outro bot; aja apenas como um monitor automático do homelab. Se não houver anomalia relevante, prefira silêncio quando o prompt permitir. USE APENAS MARKDOWN, NADA DE JSON.",
		job.Prompt, metrics,
	)

	// Forçar execução em LLM Local (Tier 0) para evitar custos de nuvem em tarefas automáticas
	ctx = agent.WithRunOptions(ctx, agent.RunOptions{LocalOnly: true})

	history, finalAnswer, err := r.executor.Execute(ctx, systemPrompt, []agent.Message{{
		Role:    "user",
		Content: analysisPrompt,
	}}, allowedTools)
	if err != nil {
		logger.Error("visual analysis failed", slog.Any("err", err))
		return "Falha na análise visual: " + err.Error(), nil, err
	}

	// Extract media parts from the last assistant message if any
	var parts []agent.ContentPart
	if len(history) > 0 {
		lastMsg := history[len(history)-1]
		if lastMsg.Role == "assistant" && len(lastMsg.Parts) > 0 {
			parts = lastMsg.Parts
		}
	}

	// Adiciona link persistente ao final do texto
	finalMsg := fmt.Sprintf("%s\n\n📊 **Acesso Direto**: [Grafana Dashboard](%s)", finalAnswer, grafanaDashboardURL)

	return finalMsg, parts, nil
}

func (r *AgentCronRuntime) extractGPUMetrics(ctx context.Context) (string, error) {
	// Extração via nvidia-smi para maior precisão
	cmd := "nvidia-smi --query-gpu=temperature.gpu,utilization.gpu,power.draw,memory.used --format=csv,noheader,nounits"
	output, err := r.executor.RunCommand(ctx, cmd)
	if err != nil {
		return "Métricas indisponíveis", err
	}
	parts := strings.Split(strings.TrimSpace(output), ",")
	if len(parts) < 4 {
		return "Métricas parciais: " + output, nil
	}
	return fmt.Sprintf("Temp: %s°C, Util: %s%%, Power: %sW, Memory: %sMiB",
		strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]),
		strings.TrimSpace(parts[2]), strings.TrimSpace(parts[3])), nil
}

func (r *AgentCronRuntime) handleSentinelWatchdog(ctx context.Context, job CronJob, metrics string) (string, []agent.ContentPart, error) {
	// Log de monitoramento interno (silencioso)
	observability.Logger("cron.sentinel").Info("sentinel check", slog.String("metrics", metrics))

	// Lógica de emergência baseada nas métricas reais
	if strings.Contains(metrics, "Temp:") {
		// Parse simplificado para encontrar a temperatura
		idx := strings.Index(metrics, "Temp: ")
		if idx != -1 {
			tempStr := ""
			for i := idx + 6; i < len(metrics) && metrics[i] != '°'; i++ {
				tempStr += string(metrics[i])
			}
			temp, _ := strconv.Atoi(tempStr)
			if temp > 75 {
				alert := fmt.Sprintf("🚨 **ALERTA CRÍTICO DE TEMPERATURA**: GPU em %d°C. Tomando ações de emergência (Prunagem de recursos)...", temp)
				// Ação automática (exemplo: limpar cache/parar processos pesados)
				_, _, _ = r.executor.Execute(ctx, "Urgente: GPU superaquecendo. Verifique processos pesados e pare-os se necessário.", nil, []string{"run_command"})

				finalMsg := fmt.Sprintf("%s\n\n📊 **Verificar Agora**: [Dashboard](%s)", alert, grafanaDashboardURL)
				return finalMsg, []agent.ContentPart{{Type: agent.ContentPartText, Text: finalMsg}}, nil
			}
		}
	}

	return "", nil, nil
}

func cronMarker(prompt string) string {
	if !strings.HasPrefix(prompt, "[sys:") {
		return ""
	}
	end := strings.Index(prompt, "]")
	if end < 0 {
		return ""
	}
	return prompt[:end+1]
}

func formatChatTeamKey(chatID int64) string {
	return strconv.FormatInt(chatID, 10)
}

type DeliveryFunc func(ctx context.Context, job CronJob, output string, parts []agent.ContentPart, execErr error) error

type NotifyingRuntime struct {
	inner   Runtime
	deliver DeliveryFunc
}

func NewNotifyingRuntime(inner Runtime, deliver DeliveryFunc) *NotifyingRuntime {
	return &NotifyingRuntime{
		inner:   inner,
		deliver: deliver,
	}
}

func (r *NotifyingRuntime) ExecuteJob(ctx context.Context, job CronJob) (string, []agent.ContentPart, error) {
	if r.inner == nil {
		return "", nil, fmt.Errorf("inner runtime is required")
	}

	output, parts, err := r.inner.ExecuteJob(ctx, job)
	if r.deliver != nil {
		if deliverErr := r.deliver(ctx, job, output, parts, err); deliverErr != nil {
			return output, parts, deliverErr
		}
	}
	return output, parts, err
}
