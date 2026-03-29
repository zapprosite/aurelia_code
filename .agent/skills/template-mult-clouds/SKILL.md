---
name: template-mult-clouds
description: Bootstrap universal do repositório multi-agente com Antigravity + Claude + OpenCode.
---

# 🚀 Template Mult-Clouds: Sovereign Infrastructure 2026

Skill estratégica para orquestrar a colaboração entre múltiplos provedores de nuvem (Cloud) e o host local (Sovereign) com Antigravity + Claude + OpenCode, garantindo redundância e performance.

## 🏛️ Estratégia de Provedores (Soberana)
- **Sovereign Host**: Ollama `qwen3.5` rodando na RTX 4090 local (Tier 0 — custo zero).
- **Cloud Rápido**: OpenRouter `google/gemini-2.5-flash` (Tier 1 — baixo custo).
- **Cloud Profundo**: OpenRouter `google/gemini-2.5-pro` (Tier 2 — raciocínio avançado).
- **Orquestração**: Antigravity + Claude + OpenCode como motores de execução.

## 🏗️ Padrões de Integração
1. **Contexto Unificado**: Uso do `ai-context` para manter a memória síncrona entre provedores.
2. **Handoff Inteligente**: Definir quando trocar de agente (ex: Claude para implementação, Antigravity para interface).
3. **Escudo de Custos**: Priorizar execução local e modelos de baixo custo conforme a política do `cost-reducer`.

## 📍 Quando usar
- No setup inicial de novos repositórios ou sub-serviços.
- Para reconfigurar o gateway em cenários de outage de provedores externos.
- Para gerenciar a identidade da Aurélia em diferentes ecossistemas.