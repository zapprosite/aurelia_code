package telegram

import (
	"strings"
	"testing"
)

func TestMaybeBuildAntigravityDelegationPrompt_Config(t *testing.T) {
	got := maybeBuildAntigravityDelegationPrompt("localize a config do Groq e monte um curl curto")
	if got == nil {
		t.Fatal("expected delegation prompt")
	}
	if got.Kind != antigravityTaskConfig {
		t.Fatalf("unexpected kind %q", got.Kind)
	}
	if got.Request.Kind != antigravityTaskConfig {
		t.Fatalf("unexpected structured kind %q", got.Request.Kind)
	}
	if !strings.Contains(got.Prompt, "menor mudanca possivel") {
		t.Fatalf("expected config guidance, got %q", got.Prompt)
	}
	if !strings.Contains(got.Prompt, "\"status\": \"approved|revise|blocked\"") {
		t.Fatalf("expected structured response contract, got %q", got.Prompt)
	}
}

func TestMaybeBuildAntigravityDelegationPrompt_Research(t *testing.T) {
	got := maybeBuildAntigravityDelegationPrompt("pesquise as melhores praticas de chunking de audio")
	if got == nil {
		t.Fatal("expected delegation prompt")
	}
	if got.Kind != antigravityTaskResearch {
		t.Fatalf("unexpected kind %q", got.Kind)
	}
	if !strings.Contains(got.Prompt, "Nao escreva tutorial longo") {
		t.Fatalf("expected research guidance, got %q", got.Prompt)
	}
}

func TestMaybeBuildAntigravityDelegationPrompt_HighRiskReturnsNil(t *testing.T) {
	if got := maybeBuildAntigravityDelegationPrompt("configure o token e faca deploy em producao"); got != nil {
		t.Fatalf("expected nil for high-risk task, got %+v", got)
	}
}

func TestParseAntigravityHandoffResult(t *testing.T) {
	got, err := parseAntigravityHandoffResult(`{"status":"approved","summary":"ok","commands":["go test ./..."],"validation":["smoke"],"residual_risk":"baixo"}`)
	if err != nil {
		t.Fatalf("parseAntigravityHandoffResult() error = %v", err)
	}
	if got.Status != "approved" {
		t.Fatalf("unexpected status %q", got.Status)
	}
}
