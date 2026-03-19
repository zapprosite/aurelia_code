package main

import (
	"testing"

	"github.com/kocar/aurelia/internal/config"
)

func TestBuildPrimaryLLMHealthCheck_WarningWhenMissing(t *testing.T) {
	t.Parallel()

	result := buildPrimaryLLMHealthCheck(&config.AppConfig{})()
	if result.Status != "warning" {
		t.Fatalf("status = %q", result.Status)
	}
}

func TestBuildPrimaryLLMHealthCheck_OkWhenConfigured(t *testing.T) {
	t.Parallel()

	result := buildPrimaryLLMHealthCheck(&config.AppConfig{
		LLMProvider: "openrouter",
		LLMModel:    "minimax/minimax-m2.7",
	})()
	if result.Status != "ok" {
		t.Fatalf("status = %q", result.Status)
	}
	if result.Message != "openrouter/minimax/minimax-m2.7" {
		t.Fatalf("message = %q", result.Message)
	}
}
