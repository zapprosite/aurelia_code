# Auditoria Estrutural Soberana (27/03/2026)

Este documento registra a crítica arquitetural de elite do monorepo `aurelia` após a industrialização SOTA 2026.

## Crítica por Pasta (Raiz)

| Pasta | Status | Crítica (Sênior/SOTA 2026) |
| :--- | :--- | :--- |
| `bin/` | ✅ Estável | **Ponto Positivo**: Binários agora são isolados por arquitetura. **Crítica**: Falta um `Makefile` ou `Taskfile` unificado que gerencie os cross-builds de forma soberana (atualmente depende de scripts bash). |
| `cmd/` | 🟠 Regular | **Ponto Positivo**: Bootstrap limpo. **Crítica**: O `app.go` está se tornando um "God Object". Muita lógica de fiação manual (wiring) que deveria ser extraída para um `internal/wiring` ou similar. |
| `docs/` | ✅ SOTA | **Ponto Positivo**: ADRs e Governança são a lei. **Crítica**: Muita documentação legada misturada com SOTA 2026. Recomenda-se um expurgo de `.md` de 2024/2025. |
| `frontend/` | 🟠 Regular | **Ponto Positivo**: Next.js é sólido. **Crítica**: Presença de detritos Go e arquivos de teste órfãos. Precisa de uma limpeza de `dist/` e `node_modules`. |
| `internal/` | 🚀 Elite | **Ponto Positivo**: Desacoplamento de `alog` e `observability` concluído. **Crítica**: A pasta `internal/agent` está sobrecarregada com Heartbeat, Loop e Planner. Deveria ser segmentada. |
| `pkg/` | ✅ Estável | **Ponto Positivo**: STT e LLM isolados. **Crítica**: O pacote `stt` carrega dependências de logging. Deveria ser 100% puro (usando interfaces de logger injetadas). |
| `scripts/` | 🟠 Legado | **Ponto Positivo**: `aureliactl` é a interface de comando. **Crítica**: Muitos scripts fragmentados. Necessário consolidar em uma CLI única (Go) para evitar dependência de bash em Homelab. |
| `vendor/` | 🟢 Seguro | **Ponto Positivo**: Garante soberania offline. **Crítica**: Ocupa muito espaço. Avaliar uso de `proxy.golang.org` se a soberania total de rede não for mais o único objetivo (mas para Homelab SOTA, manter é correto). |

## Ações de Refinamento Pendentes (Abril 2026)
1. **[X] Purificação de Log**: Concluído via `internal/purity/alog`.
2. **[ ] Extração de Wiring**: Mover lógica de bootstrap de `cmd/aurelia/app.go` para pacotes específicos.
3. **[ ] Unificação de CLI**: Transformar `aureliactl.sh` em um comando interno nativo do binário `aurelia`.

---
**Assinado**: Antigravity (SOTA 2026)
**Data**: 27/03/2026
**Soberania**: Total
