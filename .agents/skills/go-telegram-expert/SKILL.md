---
name: go-telegram-expert
description: Especialista em arquitetura síncrona/assíncrona de bots Telegram em Go (Padrões 2026).
---

# 🤖 Go Telegram Expert Skill

Esta skill aplica padrões arquiteturais de alta performance e manutenibilidade para bots Telegram escritos em Go, atualizados para as práticas de Março/2026.

## 🏛️ Padrões Arquiteturais
1. **Layered Design**:
   - `internal/telegram/handlers`: Filtros de entrada e roteamento.
   - `internal/services`: Lógica de negócio agnóstica de interface.
   - `internal/repository`: Persistência (SQLite/PostgreSQL).
2. **Finite State Machines (FSM)**:
   - Use estados explícitos para conversas complexas (ex: onboarding, checklists).
3. **Concurrency Control**:
   - Utilize Goroutines com `context.Context` para evitar vazamento de memória em processos de longa duração.

## 🚀 Performance & Observabilidade
- **Rate Limiting**: Implemente janelas de execução para respeitar os limites da API do Telegram.
- **Structured Logging**: Use `slog` ou similar para logs rotulados com `user_id` e `update_id`.
- **Healthchecks**: Mantenha monitoramento do status da conexão (Long Polling vs Webhooks).

## 🛡️ Segurança
- **Input Sanitization**: Valide todo input vindo do Telegram antes de passar para o LLM.
- **Secret Management**: Nunca hardcode tokens; utilize o `AURELIA_HOME` para resolver caminhos de configuração.

## Quando usar
- Durante reviews de código Go focado em Telegram.
- Para refatorar handlers de mensagem.
- Para implementar novos fluxos interativos.
