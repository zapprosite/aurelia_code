package main

import (
	"time"

	"github.com/kocar/aurelia/internal/gateway"
)

type llmRuntimeSnapshot struct {
	RequestedProvider string                  `json:"requested_provider"`
	RequestedModel    string                  `json:"requested_model"`
	EffectiveProvider string                  `json:"effective_provider"`
	EffectiveModel    string                  `json:"effective_model"`
	ViaGateway        bool                    `json:"via_gateway"`
	CheckedAt         time.Time               `json:"checked_at"`
	Gateway           *gateway.StatusSnapshot `json:"gateway,omitempty"`
}

func buildLLMRuntimeSnapshot(a *app, checkedAt time.Time) llmRuntimeSnapshot {
	snapshot := llmRuntimeSnapshot{
		RequestedProvider: "unconfigured",
		RequestedModel:    "",
		EffectiveProvider: "unconfigured",
		EffectiveModel:    "",
		CheckedAt:         checkedAt,
	}
	if a == nil || a.cfg == nil {
		return snapshot
	}

	snapshot.RequestedProvider = a.cfg.LLMProvider
	snapshot.RequestedModel = a.cfg.LLMModel
	snapshot.EffectiveProvider = a.cfg.LLMProvider
	snapshot.EffectiveModel = a.cfg.LLMModel

	if gw, ok := a.llmProvider.(*gateway.Provider); ok && gw != nil {
		gwSnapshot := gw.StatusSnapshot()
		snapshot.EffectiveProvider = "gateway"
		snapshot.ViaGateway = true
		snapshot.Gateway = &gwSnapshot
	}

	return snapshot
}
