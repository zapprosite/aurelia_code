package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// BrasilAPIBaseURL is exported so tests can override the endpoint.
var BrasilAPIBaseURL = "https://brasilapi.com.br/api"

var nonDigit = regexp.MustCompile(`\D`)

// CPFCNPJHandler handles cpf_cnpj tool calls.
// - action "validate_cpf": valida CPF via algoritmo (sem API externa)
// - action "validate_cnpj": valida CNPJ via algoritmo
// - action "lookup_cnpj": consulta dados da empresa na BrasilAPI (gratuito)
func CPFCNPJHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	action, err := requireStringArg(args, "action")
	if err != nil {
		return "", err
	}
	number := nonDigit.ReplaceAllString(optionalStringArg(args, "number"), "")

	switch action {
	case "validate_cpf":
		return validateCPF(number)
	case "validate_cnpj":
		return validateCNPJ(number)
	case "lookup_cnpj":
		return lookupCNPJ(ctx, number)
	default:
		return "", fmt.Errorf("ação inválida: %q — use validate_cpf, validate_cnpj ou lookup_cnpj", action)
	}
}

// ── CPF ──────────────────────────────────────────────────────────────────────

func validateCPF(digits string) (string, error) {
	if len(digits) != 11 {
		return fmt.Sprintf("CPF inválido: deve ter 11 dígitos, recebido %d (%s)", len(digits), digits), nil
	}
	// Reject all-same digits (111.111.111-11 etc.)
	if strings.Count(digits, string(digits[0])) == 11 {
		return "CPF inválido: sequência de dígitos repetidos", nil
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(digits[i]-'0') * (10 - i)
	}
	r1 := (sum * 10) % 11
	if r1 == 10 || r1 == 11 {
		r1 = 0
	}
	if r1 != int(digits[9]-'0') {
		return fmt.Sprintf("CPF %s — inválido (dígito verificador incorreto)", formatCPF(digits)), nil
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(digits[i]-'0') * (11 - i)
	}
	r2 := (sum * 10) % 11
	if r2 == 10 || r2 == 11 {
		r2 = 0
	}
	if r2 != int(digits[10]-'0') {
		return fmt.Sprintf("CPF %s — inválido (dígito verificador incorreto)", formatCPF(digits)), nil
	}

	return fmt.Sprintf("CPF %s — válido ✓", formatCPF(digits)), nil
}

func formatCPF(d string) string {
	if len(d) != 11 {
		return d
	}
	return d[:3] + "." + d[3:6] + "." + d[6:9] + "-" + d[9:]
}

// ── CNPJ ─────────────────────────────────────────────────────────────────────

func validateCNPJ(digits string) (string, error) {
	if len(digits) != 14 {
		return fmt.Sprintf("CNPJ inválido: deve ter 14 dígitos, recebido %d (%s)", len(digits), digits), nil
	}
	if strings.Count(digits, string(digits[0])) == 14 {
		return "CNPJ inválido: sequência de dígitos repetidos", nil
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i, w := range weights1 {
		sum += int(digits[i]-'0') * w
	}
	r1 := sum % 11
	if r1 < 2 {
		r1 = 0
	} else {
		r1 = 11 - r1
	}
	if r1 != int(digits[12]-'0') {
		return fmt.Sprintf("CNPJ %s — inválido (dígito verificador incorreto)", formatCNPJ(digits)), nil
	}

	sum = 0
	for i, w := range weights2 {
		sum += int(digits[i]-'0') * w
	}
	r2 := sum % 11
	if r2 < 2 {
		r2 = 0
	} else {
		r2 = 11 - r2
	}
	if r2 != int(digits[13]-'0') {
		return fmt.Sprintf("CNPJ %s — inválido (dígito verificador incorreto)", formatCNPJ(digits)), nil
	}

	return fmt.Sprintf("CNPJ %s — válido ✓", formatCNPJ(digits)), nil
}

func formatCNPJ(d string) string {
	if len(d) != 14 {
		return d
	}
	return d[:2] + "." + d[2:5] + "." + d[5:8] + "/" + d[8:12] + "-" + d[12:]
}

// ── BrasilAPI CNPJ Lookup ─────────────────────────────────────────────────────

func lookupCNPJ(ctx context.Context, digits string) (string, error) {
	if len(digits) != 14 {
		return fmt.Sprintf("CNPJ deve ter 14 dígitos para consulta, recebido %d", len(digits)), nil
	}

	url := fmt.Sprintf("%s/cnpj/v1/%s", BrasilAPIBaseURL, digits)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "aurelia-bot/1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar BrasilAPI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Sprintf("CNPJ %s não encontrado na Receita Federal.", formatCNPJ(digits)), nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("BrasilAPI retornou status %d: %s", resp.StatusCode, string(body)), nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("erro ao parsear resposta JSON: %w", err)
	}

	return formatCNPJResponse(digits, data), nil
}

func formatCNPJResponse(digits string, data map[string]interface{}) string {
	str := func(key string) string {
		v, _ := data[key].(string)
		return strings.TrimSpace(v)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**CNPJ %s**\n", formatCNPJ(digits)))
	sb.WriteString(fmt.Sprintf("**Razão Social:** %s\n", str("razao_social")))
	if nome := str("nome_fantasia"); nome != "" {
		sb.WriteString(fmt.Sprintf("**Nome Fantasia:** %s\n", nome))
	}
	sb.WriteString(fmt.Sprintf("**Situação:** %s\n", str("descricao_situacao_cadastral")))
	sb.WriteString(fmt.Sprintf("**Porte:** %s\n", str("descricao_porte")))
	sb.WriteString(fmt.Sprintf("**Natureza Jurídica:** %s\n", str("natureza_juridica")))

	// Address
	logradouro := str("logradouro")
	numero := str("numero")
	municipio := str("municipio")
	uf := str("uf")
	cep := str("cep")
	if logradouro != "" {
		sb.WriteString(fmt.Sprintf("**Endereço:** %s, %s — %s/%s CEP %s\n", logradouro, numero, municipio, uf, cep))
	}

	// Activity
	if atividade, ok := data["cnae_fiscal_descricao"].(string); ok && atividade != "" {
		sb.WriteString(fmt.Sprintf("**Atividade Principal:** %s\n", atividade))
	}

	// Dates
	if abertura := str("data_inicio_atividade"); abertura != "" {
		sb.WriteString(fmt.Sprintf("**Abertura:** %s\n", abertura))
	}
	if capital, ok := data["capital_social"].(float64); ok && capital > 0 {
		sb.WriteString(fmt.Sprintf("**Capital Social:** R$ %.2f\n", capital))
	}

	return sb.String()
}
