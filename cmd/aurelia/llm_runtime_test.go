package main

import (
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/gateway"
)

func TestBuildLLMRuntimeSnapshot_DirectProvider(t *testing.T) {
	checkedAt := time.Unix(10, 0).UTC()
	a := &app{
		cfg: &config.AppConfig{
			LLMProvider: "ollama",
			LLMModel:    "gemma3:27b",
		},
	}

	snapshot := buildLLMRuntimeSnapshot(a, checkedAt)
	if snapshot.RequestedProvider != "ollama" || snapshot.EffectiveProvider != "ollama" {
		t.Fatalf("unexpected provider snapshot: %#v", snapshot)
	}
	if snapshot.ViaGateway {
		t.Fatalf("expected direct provider, got gateway snapshot %#v", snapshot)
	}
}

func TestBuildLLMRuntimeSnapshot_GatewayProvider(t *testing.T) {
	checkedAt := time.Unix(20, 0).UTC()
	cfg := &config.AppConfig{
		LLMProvider:      "openrouter",
		LLMModel:         "deepseek/deepseek-chat-v3.1",
		OpenRouterAPIKey: "secret",
	}
	provider, err := gateway.NewProvider(cfg)
	if err != nil {
		t.Fatalf("gateway.NewProvider() error = %v", err)
	}
	defer provider.Close()

	a := &app{cfg: cfg, llmProvider: provider}
	snapshot := buildLLMRuntimeSnapshot(a, checkedAt)
	if snapshot.EffectiveProvider != "gateway" {
		t.Fatalf("expected effective provider gateway, got %#v", snapshot)
	}
	if !snapshot.ViaGateway || snapshot.Gateway == nil {
		t.Fatalf("expected gateway metadata, got %#v", snapshot)
	}
}
