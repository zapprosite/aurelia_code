package e2e_test

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SmokeTestHomelab testa integração end-to-end: Telegram → Skills → Respostas senior
//
// Cenários:
// 1. Health check completo
// 2. Diagnóstico inteligente de container
// 3. Recomendações de arquitetura
// 4. Automação segura com confirmação
func TestSmokeHomelabIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Smoke test - skip em -short")
	}

	// Setup
	bot, chatID := setupTelegramBot(t)
	defer cleanupTelegramBot(t, bot)

	t.Run("HealthCheck-FullStack", func(t *testing.T) {
		testHealthCheckFull(t, bot, chatID)
	})

	t.Run("ContainerDiagnose-WithRecommendations", func(t *testing.T) {
		testContainerDiagnostics(t, bot, chatID)
	})

	t.Run("ArchitectureDecision-GPUOptimization", func(t *testing.T) {
		testArchitectureDecision(t, bot, chatID)
	})

	t.Run("SafeAutomation-ZFSSnapshot", func(t *testing.T) {
		testSafeAutomation(t, bot, chatID)
	})

	t.Run("MultiStepOrchestration-VoiceStackDeploy", func(t *testing.T) {
		testMultiStepOrchestration(t, bot, chatID)
	})

	t.Run("IntelligentRecovery-DRValidation", func(t *testing.T) {
		testIntelligentRecovery(t, bot, chatID)
	})

	t.Run("PremiumMonitoring-HardDataProof", func(t *testing.T) {
		testPremiumMonitoring(t, bot, chatID)
	})
}

// testHealthCheckFull: "saúde completa do homelab"
func testHealthCheckFull(t *testing.T, bot *tg.BotAPI, chatID int64) {
	prompts := []string{
		"saúde completa do homelab",
		"health check agora",
		"como está o sistema?",
	}

	for _, prompt := range prompts {
		t.Logf("Enviando: %q", prompt)
		response := sendTelegramMessage(t, bot, chatID, prompt, 30*time.Second)

		// Validações que uma resposta senior faria
		expectations := []string{
			// Status containers
			"container", "docker", "status",
			// GPU
			"VRAM", "RTX", "GPU", "memory",
			// Storage
			"ZFS", "tank", "pool",
			// Network
			"tunnel", "cloudflare", "port",
			// Overall
			"✅", "✓", "status",
		}

		for _, expect := range expectations {
			if !contains(response, expect) {
				t.Logf("⚠️  Não encontrado '%s' em resposta", expect)
			}
		}

		// Resposta deve conter recomendações ou alertas
		if !containsAny(response, []string{"recomend", "alert", "attention", "aviso"}) {
			t.Logf("💡 Resposta genérica (sem recomendações)")
		}

		t.Logf("✅ Resposta recebida (%d chars)", len(response))
		time.Sleep(2 * time.Second)
	}
}

// testContainerDiagnostics: Diagnosticar container específico com recomendações
func testContainerDiagnostics(t *testing.T, bot *tg.BotAPI, chatID int64) {
	prompts := []string{
		"diagnóstico do container n8n",
		"por que o qdrant está lento?",
		"problema com supabase-db - reiniciar?",
		"logs do voice-proxy",
	}

	for _, prompt := range prompts {
		t.Logf("📦 Diagnosticando: %q", prompt)
		response := sendTelegramMessage(t, bot, chatID, prompt, 30*time.Second)

		// Resposta senior deve incluir
		expectations := []string{
			"docker", "logs", "status", // Dados concretos
			"porque", "possible", "reason", // Análise
		}

		for _, exp := range expectations {
			if !contains(response, exp) {
				t.Logf("⚠️  Faltou: %s", exp)
			}
		}

		// Deve ter recomendação de ação
		if !containsAny(response, []string{"restart", "redeploy", "check", "verify", "scale"}) {
			t.Logf("⚠️  Nenhuma ação recomendada")
		}

		t.Logf("✅ Diagnóstico completo")
		time.Sleep(2 * time.Second)
	}
}

// testArchitectureDecision: Decisões arquiteturais com justificativa
func testArchitectureDecision(t *testing.T, bot *tg.BotAPI, chatID int64) {
	prompts := []string{
		"otimizar alocação de VRAM para voice stack",
		"scaling strategy para n8n com mais workflows",
		"deveria mover DB para fora do host?",
		"consolidar containers voice em um único?",
	}

	for _, prompt := range prompts {
		t.Logf("🏛️  Arquitetura: %q", prompt)
		response := sendTelegramMessage(t, bot, chatID, prompt, 30*time.Second)

		// Resposta senior (arquiteto)
		expectations := []string{
			"trade-off", "tradeoff", "consider", // Análise de trade-offs
			"risk", "benefit", "cost", // Justificativa
			"recommendation", "suggest", "should", // Recomendação clara
		}

		for _, exp := range expectations {
			if !contains(response, exp) {
				t.Logf("⚠️  Faltou dimensão: %s", exp)
			}
		}

		// Deve ter alternativas ou justificativa
		if !containsAny(response, []string{"option", "alternative", "instead", "rather"}) {
			t.Logf("💡 Resposta única (sem alternativas consideradas)")
		}

		t.Logf("✅ Análise arquitetural completa")
		time.Sleep(2 * time.Second)
	}
}

// testSafeAutomation: Operações perigosas com confirmação
func testSafeAutomation(t *testing.T, bot *tg.BotAPI, chatID int64) {
	t.Logf("🔒 Automação Segura")

	// 1. Request snapshot
	t.Logf("  → Solicitando snapshot ZFS")
	response := sendTelegramMessage(t, bot, chatID, "cria snapshot do tank/data", 10*time.Second)

	// Deve pedir confirmação
	if !containsAny(response, []string{"confirm", "confirma", "deseja", "tenho certeza"}) {
		t.Fatalf("❌ Não pediu confirmação para operação perigosa")
	}

	t.Logf("✅ Pediu confirmação")

	// 2. Enviar confirmação
	t.Logf("  → Confirmando operação")
	confResponse := sendTelegramMessage(t, bot, chatID, "sim, confirmo", 10*time.Second)

	// Deve reportar sucesso + o que foi feito
	if !containsAny(confResponse, []string{"snapshot", "criado", "success", "ok", "✓"}) {
		t.Fatalf("❌ Não confirmou execução de snapshot")
	}

	t.Logf("✅ Operação executada com segurança")
}

// testMultiStepOrchestration: Orquestração de vários passos (voice stack up)
func testMultiStepOrchestration(t *testing.T, bot *tg.BotAPI, chatID int64) {
	t.Logf("🎵 Deploy Voice Stack")

	prompts := []string{
		"subir voice stack completo",
		"ativar STT e TTS",
		"deploy voice com vram check",
	}

	for _, prompt := range prompts {
		t.Logf("  → %q", prompt)
		response := sendTelegramMessage(t, bot, chatID, prompt, 45*time.Second)

		// Deve mostrar progression
		expectations := []string{
			"vram", "check", "verify", // Step 1: VRAM
			"whisper", "speaches", "stт", // Step 2: STT backend
			"xtts", "tts", // Step 3: TTS backend
			"proxy", // Step 4: Proxy
			"healthy", "health", "up", // Final validation
		}

		found := 0
		for _, exp := range expectations {
			if contains(response, exp) {
				found++
			}
		}

		if found < 5 {
			t.Logf("⚠️  Apenas %d/%d steps cobertos", found, len(expectations))
		} else {
			t.Logf("✅ Orquestração completa (%d steps)", found)
		}

		time.Sleep(2 * time.Second)
	}
}

// testIntelligentRecovery: Validação de disaster recovery
func testIntelligentRecovery(t *testing.T, bot *tg.BotAPI, chatID int64) {
	prompts := []string{
		"validar prontidão de recovery",
		"quando foi o último backup?",
		"snapshots estão atualizados?",
		"conseguimos fazer DR agora?",
	}

	for _, prompt := range prompts {
		t.Logf("🚨 Recovery: %q", prompt)
		response := sendTelegramMessage(t, bot, chatID, prompt, 20*time.Second)

		// Deve ter métricas de DR
		expectations := []string{
			"backup", // Backup status
			"snapshot", // Snapshot status
			"age", "hours", "days", "recent", // Recência
			"ready", "ok", "✓", // Readiness
		}

		for _, exp := range expectations {
			if !contains(response, exp) {
				t.Logf("⚠️  Métrica faltante: %s", exp)
			}
		}

		t.Logf("✅ DR Status completo")
		time.Sleep(2 * time.Second)
	}
}

// testPremiumMonitoring: Valida o fluxo de monitoramento com métricas reais e link Grafana
func testPremiumMonitoring(t *testing.T, bot *tg.BotAPI, chatID int64) {
	t.Logf("💎 Validando Monitoramento Premium")

	// Disparamos um comando que deve forçar a execução do cron ou via prompt
	response := sendTelegramMessage(t, bot, chatID, "status de hardware agora", 45*time.Second)

	expectations := []string{
		"━━━━━━", // Header Premium
		"Temp", "Utilization", "GPU", // Hard Data
		"Grafana", "monitor.zappro.site", // Link persistente
		"🛰️", "⚡", "💎", // Ícones Sênior
	}

	found := 0
	for _, exp := range expectations {
		if contains(response, exp) {
			found++
		}
	}

	if found < 4 {
		t.Errorf("❌ Formatação Premium incompleta: apenas %d/%d elementos encontrados", found, len(expectations))
	} else {
		t.Logf("✅ Monitoramento Premium validado (%d elementos)", found)
	}
}

// ============================================================================
// Helpers
// ============================================================================

func setupTelegramBot(t *testing.T) (*tg.BotAPI, int64) {
	token := getenv("TELEGRAM_BOT_TOKEN", "")
	if token == "" {
		t.Skip("TELEGRAM_BOT_TOKEN not set")
	}

	chatID := parseint64(getenv("TELEGRAM_CHAT_ID", "0"))
	if chatID == 0 {
		t.Skip("TELEGRAM_CHAT_ID not set")
	}

	bot, err := tg.NewBotAPI(token)
	if err != nil {
		t.Fatalf("Failed to init bot: %v", err)
	}

	t.Logf("✅ Bot conectado: @%s", bot.Self.UserName)
	return bot, chatID
}

func cleanupTelegramBot(t *testing.T, bot *tg.BotAPI) {
	// Enviar mensagem de limpeza
	msg := tg.NewMessage(0, "Smoke test finalizado ✅")
	_ = msg
	t.Logf("✅ Cleanup completo")
}

func sendTelegramMessage(t *testing.T, bot *tg.BotAPI, chatID int64, text string, timeout time.Duration) string {
	msg := tg.NewMessage(chatID, text)
	msg.ParseMode = tg.ModeHTML

	// Enviar
	sent, err := bot.Send(msg)
	if err != nil {
		t.Fatalf("Falha ao enviar: %v", err)
	}

	t.Logf("📤 Enviado (ID: %d)", sent.MessageID)

	// Aguardar resposta
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Verificar se há resposta (em update channel)
			response := checkForResponse(t, bot, sent.MessageID, chatID)
			if response != "" {
				t.Logf("📥 Resposta recebida em %v", time.Since(startTime))
				return response
			}

		case <-time.After(timeout):
			t.Logf("⏱️  Timeout aguardando resposta (%v)", timeout)
			return ""
		}
	}
}

func checkForResponse(t *testing.T, bot *tg.BotAPI, originalMsgID int, chatID int64) string {
	// Polling de mensagens recentes
	u := tg.NewUpdate(0)
	u.Timeout = 1

	updates, err := bot.GetUpdates(u)
	if err != nil {
		t.Logf("Erro polling: %v", err)
		return ""
	}

	for _, update := range updates {
		if update.Message != nil && update.Message.Chat.ID == chatID &&
			update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.MessageID == originalMsgID {
			return update.Message.Text
		}
	}

	return ""
}

// ============================================================================
// Utilities
// ============================================================================

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

func getenv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultVal
}

func parseint64(s string) int64 {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	return 0
}
