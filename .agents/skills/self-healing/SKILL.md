---
type: skill
name: Self-Healing
description: Habilita a detecção e correção automática de falhas de sistema e binários zumbis.
skillSlug: self-healing
phases: [O]
generated: 2026-03-20
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 🩹 Self-Healing: Sovereign Watchdog 2026

Capacita a Aurélia a detectar "loops" de erro, pânicos de kernel ou travamentos de serviços e aplicar correções automáticas (Self-Correction) baseadas em heurísticas e logs.

## 🛠️ Heurísticas de Correção
1. **Docker Down**: Se um container essencial cair, a skill tenta um `docker compose up -d` após limpar volumes temporários se necessário.
2. **Git Zombie**: Invoca automaticamente a skill `/git-unblocker` se detectar locks ativos.
3. **Gateway Memory**: Se o Ollama estourar a VRAM, a skill reinicia o daemon local com parâmetros de memória otimizados.

## 🛡️ Guardrail de Segurança
- **Backoff**: Após 3 tentativas falhas de correção automática, o sistema deve parar, entrar em modo de pânico (Logs Críticos) e notificar o humano no Telegram com o status `EMERGENCY`.

## 📍 Quando usar
- Quando o daemon principal começar a reportar erros cíclicos.
- Quando ferramentas dependentes de hardware falharem em cascata.
- Durante a manutenção de madrugada para garantir que o sistema acorde saudável.
