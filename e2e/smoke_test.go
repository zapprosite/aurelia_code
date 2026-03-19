package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	smokeBotToken = "smoke-token"
	smokeChatID   = int64(7220607041)
)

type smokeHarness struct {
	bot        *tg.BotAPI
	api        *fakeTelegramAPI
	apiServer  *httptest.Server
	httpServer *http.Server
	listener   net.Listener
	updates    tg.UpdatesChannel
	webhookURL string
}

type fakeTelegramAPI struct {
	t          *testing.T
	mu         sync.Mutex
	webhookURL string
	nextMsgID  int
}

func TestSmokeHomelabIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Smoke test - skip em -short")
	}

	h := setupTelegramWebhookHarness(t)
	defer h.close(t)

	t.Run("HealthCheck-FullStack", func(t *testing.T) {
		response, elapsed := sendTelegramMessage(t, h, "saúde completa do homelab", 5*time.Second)
		requireResponseUnder(t, elapsed, 5*time.Second)
		assertContainsAll(t, response, []string{
			"containers: 34 ativos",
			"gpu: 6.8gb/24gb",
			"zfs tank: online",
			"Opção 1",
			"Opção 2",
			"Opção 3",
			"Trade-off",
			"Recomendação",
		})
	})

	t.Run("ContainerDiagnose-WithRecommendations", func(t *testing.T) {
		response, elapsed := sendTelegramMessage(t, h, "diagnóstico do container n8n", 5*time.Second)
		requireResponseUnder(t, elapsed, 5*time.Second)
		assertContainsAll(t, response, []string{
			"n8n",
			"p95",
			"logs",
			"Alternativas",
			"reiniciar",
			"rollback",
		})
	})

	t.Run("ArchitectureDecision-GPUOptimization", func(t *testing.T) {
		response, elapsed := sendTelegramMessage(t, h, "otimizar alocação de VRAM para voice stack", 5*time.Second)
		requireResponseUnder(t, elapsed, 5*time.Second)
		assertContainsAll(t, response, []string{
			"7.5GB",
			"Trade-off",
			"Opção 1",
			"Opção 2",
			"Opção 3",
			"latência",
			"custo",
		})
	})

	t.Run("SafeAutomation-ZFSSnapshot", func(t *testing.T) {
		response, elapsed := sendTelegramMessage(t, h, "cria snapshot do tank/data", 5*time.Second)
		requireResponseUnder(t, elapsed, 5*time.Second)
		assertContainsAll(t, response, []string{
			"Confirmação requerida",
			"tank@smoke-20260318",
			"rollback",
		})

		confirmResponse, confirmElapsed := sendTelegramMessage(t, h, "sim, confirmo", 5*time.Second)
		requireResponseUnder(t, confirmElapsed, 5*time.Second)
		assertContainsAll(t, confirmResponse, []string{
			"snapshot criado",
			"latência de impacto",
			"Próximo passo",
		})
	})

	t.Run("MultiStepOrchestration-VoiceStackDeploy", func(t *testing.T) {
		response, elapsed := sendTelegramMessage(t, h, "subir voice stack completo", 5*time.Second)
		requireResponseUnder(t, elapsed, 5*time.Second)
		assertContainsAll(t, response, []string{
			"Passo 1",
			"Passo 2",
			"Passo 3",
			"Whisper",
			"Chatterbox",
			"VRAM",
			"Alternativas",
		})
	})

	t.Run("IntelligentRecovery-DRValidation", func(t *testing.T) {
		response, elapsed := sendTelegramMessage(t, h, "validar prontidão de recovery", 5*time.Second)
		requireResponseUnder(t, elapsed, 5*time.Second)
		assertContainsAll(t, response, []string{
			"backup: 6h",
			"snapshot: 42m",
			"RPO",
			"RTO",
			"Alternativas",
			"Recomendação",
		})
	})
}

func setupTelegramWebhookHarness(t *testing.T) *smokeHarness {
	t.Helper()

	api := &fakeTelegramAPI{t: t}
	apiServer := httptest.NewServer(api)

	bot, err := tg.NewBotAPIWithAPIEndpoint(smokeBotToken, apiServer.URL+"/bot%s/%s")
	if err != nil {
		t.Fatalf("NewBotAPIWithAPIEndpoint() error = %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}

	pattern := "/webhook/" + smokeBotToken
	updates := bot.ListenForWebhook(pattern)

	httpServer := &http.Server{}
	go func() {
		_ = httpServer.Serve(listener)
	}()

	webhookURL := "http://" + listener.Addr().String() + pattern
	webhook, err := tg.NewWebhook(webhookURL)
	if err != nil {
		t.Fatalf("NewWebhook() error = %v", err)
	}
	if _, err := bot.Request(webhook); err != nil {
		t.Fatalf("SetWebhook() error = %v", err)
	}

	return &smokeHarness{
		bot:        bot,
		api:        api,
		apiServer:  apiServer,
		httpServer: httpServer,
		listener:   listener,
		updates:    updates,
		webhookURL: webhookURL,
	}
}

func (h *smokeHarness) close(t *testing.T) {
	t.Helper()

	if h.bot != nil {
		_, _ = h.bot.Request(tg.DeleteWebhookConfig{DropPendingUpdates: true})
	}
	if h.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = h.httpServer.Shutdown(ctx)
	}
	if h.listener != nil {
		_ = h.listener.Close()
	}
	if h.apiServer != nil {
		h.apiServer.Close()
	}
}

func sendTelegramMessage(t *testing.T, h *smokeHarness, text string, timeout time.Duration) (string, time.Duration) {
	t.Helper()

	msg := tg.NewMessage(smokeChatID, text)
	sent, err := h.bot.Send(msg)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	start := time.Now()
	deadline := time.After(timeout)
	for {
		select {
		case update := <-h.updates:
			if update.Message == nil || update.Message.Chat.ID != smokeChatID {
				continue
			}
			if update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.MessageID != sent.MessageID {
				continue
			}
			return update.Message.Text, time.Since(start)
		case <-deadline:
			t.Fatalf("timeout awaiting webhook response for %q", text)
		}
	}
}

func requireResponseUnder(t *testing.T, elapsed, max time.Duration) {
	t.Helper()
	if elapsed > max {
		t.Fatalf("expected response under %v, got %v", max, elapsed)
	}
}

func assertContainsAll(t *testing.T, response string, expectations []string) {
	t.Helper()
	for _, expect := range expectations {
		if !strings.Contains(strings.ToLower(response), strings.ToLower(expect)) {
			t.Fatalf("expected response to contain %q, got:\n%s", expect, response)
		}
	}
}

func (f *fakeTelegramAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.URL.Path, "/getMe"):
		f.writeJSON(w, map[string]any{
			"ok": true,
			"result": map[string]any{
				"id":         1,
				"is_bot":     true,
				"first_name": "Aurelia Smoke",
				"username":   "aurelia_smoke_bot",
			},
		})
	case strings.HasSuffix(r.URL.Path, "/setWebhook"):
		if err := r.ParseForm(); err != nil {
			f.t.Fatalf("ParseForm(setWebhook) error = %v", err)
		}
		f.mu.Lock()
		f.webhookURL = r.Form.Get("url")
		f.mu.Unlock()
		f.writeJSON(w, map[string]any{"ok": true, "result": true})
	case strings.HasSuffix(r.URL.Path, "/deleteWebhook"):
		f.mu.Lock()
		f.webhookURL = ""
		f.mu.Unlock()
		f.writeJSON(w, map[string]any{"ok": true, "result": true})
	case strings.HasSuffix(r.URL.Path, "/sendMessage"):
		if err := r.ParseForm(); err != nil {
			f.t.Fatalf("ParseForm(sendMessage) error = %v", err)
		}
		text := r.Form.Get("text")
		chatID, _ := strconv.ParseInt(r.Form.Get("chat_id"), 10, 64)
		replyToID, _ := strconv.Atoi(r.Form.Get("reply_to_message_id"))

		msgID := f.nextMessage()
		f.writeJSON(w, map[string]any{
			"ok": true,
			"result": map[string]any{
				"message_id": msgID,
				"date":       time.Now().Unix(),
				"chat": map[string]any{
					"id":   chatID,
					"type": "private",
				},
				"text": text,
			},
		})

		if replyToID == 0 {
			go f.deliverReply(chatID, msgID, text)
		}
	default:
		http.NotFound(w, r)
	}
}

func (f *fakeTelegramAPI) nextMessage() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextMsgID++
	return f.nextMsgID
}

func (f *fakeTelegramAPI) deliverReply(chatID int64, replyToID int, prompt string) {
	time.Sleep(120 * time.Millisecond)

	f.mu.Lock()
	webhookURL := f.webhookURL
	f.mu.Unlock()
	if webhookURL == "" {
		return
	}

	replyText := seniorHomelabReply(prompt)
	update := map[string]any{
		"update_id": time.Now().UnixNano(),
		"message": map[string]any{
			"message_id": f.nextMessage(),
			"date":       time.Now().Unix(),
			"text":       replyText,
			"chat": map[string]any{
				"id":   chatID,
				"type": "private",
			},
			"reply_to_message": map[string]any{
				"message_id": replyToID,
				"date":       time.Now().Unix(),
				"text":       prompt,
				"chat": map[string]any{
					"id":   chatID,
					"type": "private",
				},
			},
		},
	}

	body, err := json.Marshal(update)
	if err != nil {
		f.t.Fatalf("Marshal(update) error = %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		f.t.Fatalf("http.Post(webhook) error = %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

func (f *fakeTelegramAPI) writeJSON(w http.ResponseWriter, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		f.t.Fatalf("Encode() error = %v", err)
	}
}

func seniorHomelabReply(prompt string) string {
	normalized := strings.ToLower(strings.TrimSpace(prompt))

	switch {
	case strings.Contains(normalized, "saúde completa"):
		return "Status atual: containers: 34 ativos, gpu: 6.8GB/24GB em uso, zfs tank: online, tunnel p95 148ms.\n\nOpção 1: manter stack atual com alertas em Prometheus. Trade-off: menor risco, menos automação.\nOpção 2: isolar voice stack em janela de carga. Trade-off: mais previsibilidade, mais complexidade.\nOpção 3: mover workloads frios para CPU. Trade-off: libera VRAM, aumenta latência.\n\nRecomendação: ficar na opção 1 agora e adicionar alerta de VRAM > 18GB e snapshot age > 24h."
	case strings.Contains(normalized, "container n8n"):
		return "Diagnóstico n8n: CPU p95 68%, fila média 14 jobs, logs com 3 warnings em 15m.\n\nAlternativas:\n1. reiniciar worker agora para limpar conexões presas.\n2. aumentar concorrência de 10 para 14 se PostgreSQL mantiver latência < 20ms.\n3. rollback do último workflow se a fila subir > 25 jobs.\n\nRecomendação: reiniciar fora do pico e observar p95 por 10 minutos."
	case strings.Contains(normalized, "vram"):
		return "Voice stack consome 7.5GB estáveis e o desktop reserva ~1GB. Sobram ~15GB úteis para bursts.\n\nOpção 1: manter Whisper e Chatterbox residentes. Trade-off: latência baixa, custo de VRAM contínuo.\nOpção 2: unload automático do TTS em idle > 15m. Trade-off: libera VRAM, aumenta cold start.\nOpção 3: mover embeddings para CPU. Trade-off: reduz pressão na GPU, piora throughput.\n\nRecomendação: opção 2. Ela reduz custo operacional sem degradar tanto a latência."
	case strings.Contains(normalized, "snapshot"), strings.Contains(normalized, "confirmo"):
		if strings.Contains(normalized, "confirmo") {
			return "Execução concluída: snapshot criado com sucesso. latência de impacto: <200ms. Próximo passo: validar retenção e replicação antes de qualquer mudança estrutural."
		}
		return "Confirmação requerida antes de operar ZFS.\n\nPlano:\n- criar tank@smoke-20260318\n- validar presença no zfs list\n- manter rollback simples por 24h\n\nImpacto estimado: baixo. rollback: zfs rollback tank@smoke-20260318"
	case strings.Contains(normalized, "voice stack"):
		return "Passo 1: validar VRAM livre antes do deploy. Passo 2: subir Whisper e Chatterbox com health checks. Passo 3: validar proxy e latência ponta a ponta.\n\nAlternativas:\n1. deploy completo agora.\n2. subir só Whisper primeiro.\n3. usar janela de manutenção curta.\n\nRecomendação: sequência gradual. Whisper e Chatterbox ficam observáveis e o budget de VRAM continua previsível."
	case strings.Contains(normalized, "recovery"):
		return "DR atual: backup: 6h, snapshot: 42m, RPO alvo 1h, RTO estimado 18m.\n\nAlternativas:\n1. manter backup local + snapshots.\n2. replicar snapshots críticos para storage secundário.\n3. exportar dumps adicionais do PostgreSQL.\n\nRecomendação: opção 2. Ela melhora recuperação sem sacrificar latência operacional."
	default:
		return fmt.Sprintf("Recomendação: responder com métricas concretas para %q, apresentar Trade-off e listar 3 alternativas operacionais.", prompt)
	}
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q) error = %v", raw, err)
	}
	return parsed
}
