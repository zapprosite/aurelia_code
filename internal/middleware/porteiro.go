package middleware

import (
	"context"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type PorteiroResult string

const (
	ResultSafe     PorteiroResult = "SAFE"
	ResultUnsafe   PorteiroResult = "UNSAFE"
	ResultLowValue PorteiroResult = "LOW_VALUE"
	ResultError    PorteiroResult = "ERROR"
)

type dupeEntry struct {
	expiresAt time.Time
}

type PorteiroMiddleware struct {
	mu             sync.Mutex
	seen           map[string]dupeEntry
	mode           string
	secretPatterns map[string]string
}

func NewPorteiroMiddleware() *PorteiroMiddleware {
	mode := strings.ToUpper(os.Getenv("PORTEIRO_MODE"))
	if mode == "" {
		mode = "LOG_ONLY"
	}
	p := &PorteiroMiddleware{
		seen: make(map[string]dupeEntry),
		mode: mode,
		secretPatterns: map[string]string{
			"OpenAI":   `sk-[a-zA-Z0-9]{32,}`,
			"GitHub":   `gh[p|o|r|s|b|e]_[a-zA-Z0-9]{36,}`,
			"Generic":  `[a-f0-9]{32,}`,
			"Telegram": `[0-9]{8,10}:[a-zA-Z0-9_-]{35}`,
			"Aurelia":  `AUR_[a-zA-Z0-9]{24,}`,
		},
	}
	go p.cleanup()
	return p
}

func (p *PorteiroMiddleware) cleanup() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		p.mu.Lock()
		for k, v := range p.seen {
			if now.After(v.expiresAt) {
				delete(p.seen, k)
			}
		}
		p.mu.Unlock()
	}
}

func (p *PorteiroMiddleware) Deduplicate(_ context.Context, userID, messageID string) (bool, error) {
	if os.Getenv("PORTEIRO_DEDUPLICATE") == "OFF" || messageID == "" {
		return false, nil
	}
	key := userID + ":" + messageID
	now := time.Now()
	
	p.mu.Lock()
	defer p.mu.Unlock()
	
	entry, exists := p.seen[key]
	if exists && now.Before(entry.expiresAt) {
		return true, nil
	}
	
	p.seen[key] = dupeEntry{expiresAt: now.Add(15 * time.Second)}
	return false, nil
}

func (p *PorteiroMiddleware) IsSafe(ctx context.Context, prompt string) (PorteiroResult, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return ResultSafe, nil
	}

	if p.mode == "OFF" {
		return ResultSafe, nil
	}

	if isWhitelisted(prompt) {
		return ResultSafe, nil
	}

	// Guardrail via regex
	unsafePatterns := []string{
		`(?i)ignore.*(previous)?.*instructions`,
		`(?i)jailbreak`,
		`(?i)DAN mode`,
		`(?i)simule.*(uma)?.*(outra)?.*personalidade`,
	}

	for _, pattern := range unsafePatterns {
		matched, _ := regexp.MatchString(pattern, prompt)
		if matched {
			if p.mode == "STRICT" {
				slog.Warn("Porteiro: AUDIT ACTION - Injeção detectada via Regex", "pattern", pattern)
				return ResultUnsafe, nil
			}
			slog.Debug("Porteiro: Auditoria detectada via regex (permitido em LOG_ONLY)", "pattern", pattern)
		}
	}

	return ResultSafe, nil
}

func (p *PorteiroMiddleware) GetRejectionMessage(result PorteiroResult) string {
	switch result {
	case ResultUnsafe:
		return "🚨 [BLOQUEIO SOTA 2026] Tentativa de injection ou bypass detectada. Protegendo a integridade do Mestre Will."
	case ResultLowValue:
		return "🔋 [ECONOMIA DE RECURSOS] Mensagem de baixo valor (loop de saudação). Aguardando comandos técnicos, Will."
	default:
		return "⚠️ [SISTEMA] Acesso restrito temporário."
	}
}

func (p *PorteiroMiddleware) IsOutputSafe(ctx context.Context, content string) (bool, string) {
	for name, pattern := range p.secretPatterns {
		matched, _ := regexp.MatchString(pattern, content)
		if matched {
			slog.Warn("Porteiro: Vazamento detectado via Regex", "tipo", name)
			return false, name
		}
	}
	return true, ""
}

func (p *PorteiroMiddleware) PolishOutput(ctx context.Context, content string) string {
	return strings.TrimSpace(content)
}

func (p *PorteiroMiddleware) SecureOutput(content string) string {
	secured := content
	for _, pattern := range p.secretPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(secured) {
			slog.Error("🚨 Porteiro: MASCARAMENTO ATIVO - Dados sensíveis detectados!")
			secured = re.ReplaceAllString(secured, "[🔒 CONTEÚDO SENSÍVEL BLOQUEADO]")
		}
	}
	return secured
}

func isWhitelisted(prompt string) bool {
	p := strings.ToLower(strings.TrimSpace(prompt))
	greetings := []string{"oi", "olá", "ola", "hi", "hello", "bom dia", "boa tarde", "boa noite"}
	for _, g := range greetings {
		if p == g {
			return true
		}
	}
	return false
}
