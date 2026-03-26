# ==============================================================================
# 🛰️ AURELIA — SOBERANO MAKEFILE (2026)
# ==============================================================================

.PHONY: help setup build run test audit sync clean

help: ## Mostra esta ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

setup: ## Inicia o setup interativo e configuração de ambiente
	@chmod +x iniciar.sh
	@./iniciar.sh

build: ## Compila o binário da Aurélia
	@chmod +x scripts/build.sh
	@./scripts/build.sh

run: ## Inicia o daemon via Docker Compose
	@docker-compose up -d

test: ## Executa todos os testes unitários
	@go test -v ./...

audit: ## Executa a auditoria de segredos (Secrets)
	@chmod +x scripts/secret-audit.sh
	@./scripts/secret-audit.sh

sync: ## Sincroniza o contexto de IA (ai-context)
	@chmod +x scripts/sync-ai-context.sh
	@./scripts/sync-ai-context.sh

clean: ## Limpa artefatos de build e arquivos temporários
	@rm -f aurelia
	@find . -name "*.log" -delete
	@find . -name "aurelia-tts-*" -delete
