package tools

import (
	"context"
	"strings"
	"testing"
)

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name    string
		cpf     string
		wantOK  bool
		wantSub string
	}{
		{"valid cpf", "529.982.247-25", true, "válido ✓"},
		{"valid cpf no fmt", "52998224725", true, "válido ✓"},
		{"all same digits", "11111111111", false, "dígitos repetidos"},
		{"wrong check digit", "529.982.247-26", false, "inválido"},
		{"too short", "1234567890", false, "11 dígitos"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digits := nonDigit.ReplaceAllString(tt.cpf, "")
			result, err := validateCPF(digits)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(result, tt.wantSub) {
				t.Errorf("got %q, want substring %q", result, tt.wantSub)
			}
		})
	}
}

func TestValidateCNPJ(t *testing.T) {
	tests := []struct {
		name    string
		cnpj    string
		wantSub string
	}{
		{"valid cnpj", "11.222.333/0001-81", "válido ✓"},
		{"valid cnpj no fmt", "11222333000181", "válido ✓"},
		{"invalid cnpj", "11.222.333/0001-82", "inválido"},
		{"all same digits", "00000000000000", "dígitos repetidos"},
		{"too short", "1234567890123", "14 dígitos"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digits := nonDigit.ReplaceAllString(tt.cnpj, "")
			result, err := validateCNPJ(digits)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(result, tt.wantSub) {
				t.Errorf("got %q, want substring %q", result, tt.wantSub)
			}
		})
	}
}

func TestCPFCNPJHandler_actions(t *testing.T) {
	ctx := context.Background()

	result, err := CPFCNPJHandler(ctx, map[string]interface{}{
		"action": "validate_cpf",
		"number": "529.982.247-25",
	})
	if err != nil || !strings.Contains(result, "válido") {
		t.Errorf("validate_cpf failed: %v | %s", err, result)
	}

	result, err = CPFCNPJHandler(ctx, map[string]interface{}{
		"action": "validate_cnpj",
		"number": "11.222.333/0001-81",
	})
	if err != nil || !strings.Contains(result, "válido") {
		t.Errorf("validate_cnpj failed: %v | %s", err, result)
	}

	_, err = CPFCNPJHandler(ctx, map[string]interface{}{
		"action": "bad_action",
		"number": "123",
	})
	if err == nil {
		t.Error("expected error for bad action")
	}
}
