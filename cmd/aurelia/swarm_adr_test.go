package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/observability"
)

// TestSwarmExecutaADR simula o humano enviando um prompt desafiador pelo Telegram
// para a Aurelia, mandando ela executar a ADR de modelos locais.
// Pela menção de API "MiniMax 3.7", assume-se que o usuário configurou Claude 3.7
// (ou equivalente rápido) para as inferências da Swarm no .env local.
func TestSwarmExecutaADR(t *testing.T) {
	if os.Getenv("RUN_SWARM_E2E") == "" {
		t.Skip("Pulando E2E real para economizar tokens. Rode com RUN_SWARM_E2E=1")
	}

	observability.Configure(observability.Options{})

	// 1. Inicia o App inteiro (carrega .env, LLM Provider, Memória, Telegram Bot, etc)
	app, err := bootstrapApp([]string{"aurelia"})
	if err != nil {
		t.Fatalf("Erro crítico ao subir o ambiente Swarm: %v", err)
	}
	defer app.close()

	app.start()

	// 2. Descobrir quem é o dono do homelab para injetar a mensagem autorizada
	cfg, _ := config.Load(app.resolver)
	if len(cfg.TelegramAllowedUserIDs) == 0 {
		t.Fatalf("Nenhum TELEGRAM_ALLOWED_USER_IDS configurado no .env")
	}
	adminID := cfg.TelegramAllowedUserIDs[0]
	chatID := adminID // Em chats privados, chatID == userID

	// 3. O Desafio (Prompt)
	prompt := "Aurelia, aqui é o humano. Eu te desafio a executar agora mesmo a alteração/aplicação descrita na 'ADR de mudança de modelos locais' (ADR-20260320-politica-modelos-hardware-vram.md). Confio na velocidade do modelo configurado (MiniMax 3.7/Claude 3.7). Trabalhe com a Swarm em tempo recorde e me dê o resultado."

	t.Logf("==== INICIANDO INJEÇÃO DO PROMPT PARA O ENXAME ====")
	t.Logf("Usuário ID: %d", adminID)
	t.Logf("Prompt: %s", prompt)
	t.Logf("===================================================")

	// 4. Injeta a mensagem no core do bot, como se tivesse vindo do Webhook/LongPolling
	err = app.bot.ProcessExternalInput(context.Background(), adminID, chatID, prompt, false)
	if err != nil {
		t.Fatalf("Falha do bot ao processar o input: %v", err)
	}

	// 5. Como o enxame roda em background criando rotinas de loop e agentes (Lead -> Planner -> Child),
	// nós seguramos o teste aberto para observar o console enquanto eles "trabalham".
	t.Log("O prompt foi aceito pelo Telegram Controller e despachado pro Loop da Swarm.")
	t.Log("Aguardando 45 segundos para ver a velocidade de resposta dos Agentes no log...")

	time.Sleep(45 * time.Second)

	// Shutdown limpo
	app.shutdown(context.Background())
	t.Log("==== TESTE E2E DA SWARM FINALIZADO ====")
}

// TestSwarmImplementaP3 simula o humano enviando o desafio do Roadmap Mestre P3
// para a Aurelia, mandando a Swarm fazer o código e provar habilidades Híbridas.
func TestSwarmImplementaP3(t *testing.T) {
	if os.Getenv("RUN_SWARM_E2E") == "" {
		t.Skip("Pulando E2E real. Rode com RUN_SWARM_E2E=1")
	}

	observability.Configure(observability.Options{})

	app, err := bootstrapApp([]string{"aurelia"})
	if err != nil {
		t.Fatalf("Erro crítico ao subir o ambiente Swarm: %v", err)
	}
	defer app.close()

	app.start()

	cfg, _ := config.Load(app.resolver)
	if len(cfg.TelegramAllowedUserIDs) == 0 {
		t.Fatalf("Nenhum TELEGRAM_ALLOWED_USER_IDS configurado no .env")
	}
	adminID := cfg.TelegramAllowedUserIDs[0]
	chatID := adminID

	prompt := "Aurelia/MiniMax, aqui é o humano te desafiando: Implemente agora a Feature [P3] listada no ADR-20260320-roadmap-mestre-slices.md. Faça o cherry-pick dos arquivos 'scripts/simulate_swarm_2026.go' e da pasta 'internal/voice' da branch agent-to-agent. Mostre que a Swarm pode usar todas as skills híbridas recém unificadas e feche isso em tempo recorde!"

	t.Logf("==== INJEÇÃO DO PROMPT DE IMPLEMENTAÇÃO P3 OMNI-SWARM ====")
	t.Logf("Prompt: %s", prompt)
	t.Logf("==========================================================")

	err = app.bot.ProcessExternalInput(context.Background(), adminID, chatID, prompt, false)
	if err != nil {
		t.Fatalf("Falha do bot ao processar o input: %v", err)
	}

	t.Log("O prompt foi injetado na mente da Swarm (Loop Agentico).")
	t.Log("Aguardando 45s de processamento MiniMax...")
	time.Sleep(45 * time.Second)

	app.shutdown(context.Background())
	t.Log("==== TESTE P3 OMNI-SWARM FINALIZADO ====")
}
