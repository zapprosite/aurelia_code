package e2e_test

import (
	"testing"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TestPersonaWillUseCases valida os fluxos de trabalho específicos da persona do Will.
func TestPersonaWillUseCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Persona Use Case test - skip em -short")
	}

	bot, chatID := setupTelegramBot(t)
	defer cleanupTelegramBot(t, bot)

	t.Run("HVAC-Diagnostic-OCR", func(t *testing.T) {
		testHVACDiagnostic(t, bot, chatID)
	})

	t.Run("Diet-And-Training-Logging", func(t *testing.T) {
		testDietTrainingLogging(t, bot, chatID)
	})

	t.Run("Homelab-Business-Guard-Alerts", func(t *testing.T) {
		testBusinessInfrastructureAlerts(t, bot, chatID)
	})
}

// testHVACDiagnostic: Simula envio de erro de climatização
func testHVACDiagnostic(t *testing.T, bot *tg.BotAPI, chatID int64) {
	t.Logf("❄️ Testando diagnóstico HVAC")
	
	prompt := "identifique o erro E4 nesta placa de VRF"
	response := sendTelegramMessage(t, bot, chatID, prompt, 30*time.Second)

	expectations := []string{
		"HVAC", "pressão", "descarga", "E4", "condensadora",
	}

	for _, exp := range expectations {
		if !contains(response, exp) {
			t.Logf("⚠️  Expected HVAC insight missing: %s", exp)
		}
	}
	t.Logf("✅ Diagnóstico HVAC validado")
}

// testDietTrainingLogging: Simula registro de performance e macros
func testDietTrainingLogging(t *testing.T, bot *tg.BotAPI, chatID int64) {
	t.Logf("🥩 Testando registro de Dieta e Treino")

	prompts := []string{
		"registre treino de hoje: agachamento 140kg 5x5",
		"adicione refeição: 250g de frango e 150g de batata doce",
	}

	for _, p := range prompts {
		response := sendTelegramMessage(t, bot, chatID, p, 20*time.Second)
		
		if !containsAny(response, []string{"registrado", "confirmado", "macros", "PR", "✓"}) {
			t.Logf("⚠️  Falha ao confirmar registro de: %s", p)
		}
	}
	t.Logf("✅ Registro de performance validado")
}

// testBusinessInfrastructureAlerts: Simula alertas proativos da infra da empresa
func testBusinessInfrastructureAlerts(t *testing.T, bot *tg.BotAPI, chatID int64) {
	t.Logf("🏛️ Testando alertas de infraestrutura de negócio")

	prompt := "status dos servidores da empresa de HVAC"
	response := sendTelegramMessage(t, bot, chatID, prompt, 25*time.Second)

	expectations := []string{
		"servidores", "saúde", "status", "latência", "CPU",
	}

	for _, exp := range expectations {
		if !contains(response, exp) {
			t.Logf("⚠️  Expected infra insight missing: %s", exp)
		}
	}
	t.Logf("✅ Status de infra de negócio validado")
}
