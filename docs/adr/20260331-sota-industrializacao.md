# ADR-20260331 — Industrialização Soberana SOTA 2026.2 (Headless Appliance)

**Status**: Aceito (Padrão 2026)
**Data**: 2026-03-31
**Autor**: Antigravity (Gemini)

## Contexto

O ecossistema Aurélia cresceu organicamente com múltiplos módulos de interface (Dashboard Web, Voice CLI, Onboarding UI), gerando dívida técnica e fragmentação. Para operar como uma infraestrutura de missão crítica (Soberania 2026), o sistema exige minimalismo, estabilidade estática e foco em interfaces de comando e automação (Telegram/CLI).

## Decisão

Implementar o padrão **SOTA 2026.2 (Headless Agentic Appliance)**:

1. **Remoção de UI Legada**: Desativar e remover todos os componentes `dashboard`, `onboard` e `voice_cli` do diretório `cmd/aurelia/`.
2. **Consolidação de Serviços**: Unificar a orquestração de squads no `MasterTeamService` e a persistência no `SQLiteTaskStore` (interface `TeamManager`).
3. **Build Industrial**: Forçar compilação estática (`CGO_ENABLED=0`) para garantir portabilidade absoluta no Homelab Ubuntu.
4. **Governança via Telegram**: Delegar toda a interação humano-máquina para o pipeline de entrada do Telegram, utilizando o middleware `Porteiro` como guardião de segurança.

## Consequências

- **Positivas**: Redução drástica da superfície de ataque, binários menores e mais rápidos, facilidade de manutenção (Single Source of Truth no Go).
- **Negativas**: Perda da interface visual do dashboard (compensada por logs estruturados e telemetria via Telegram).
- **Riscos**: Dependência da conectividade com o Telegram para controle interativo (mitigada pelo `runLiveCommand` local).

---

*Assinado: Aurélia (Soberano 2026).*
