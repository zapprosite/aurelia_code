package main

import (
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/health"
	"github.com/kocar/aurelia/internal/runtime"
)

func registerAuxiliaryHealthChecks(healthSrv *health.Server, cfg *config.AppConfig, resolver *runtime.PathResolver) {
	if healthSrv == nil || cfg == nil {
		return
	}
	_ = resolver

	healthSrv.RegisterCheck("primary_llm", buildPrimaryLLMHealthCheck(cfg))
}

func buildPrimaryLLMHealthCheck(cfg *config.AppConfig) func() health.CheckResult {
	return func() health.CheckResult {
		if cfg == nil || cfg.LLMProvider == "" {
			return health.CheckResult{Status: "warning", Message: "llm provider not configured"}
		}
		return health.CheckResult{Status: "ok", Message: cfg.LLMProvider + "/" + cfg.LLMModel}
	}
}
