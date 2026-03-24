---
name: template-mult-clouds
description: Bootstrap universal do repositório multi-agente com Antigravity + Claude + Codex + OpenCode.
---

# 🚀 Template Mult-Clouds: Sovereign Infrastructure 2026

Skill estratégica para orquestrar a colaboração entre múltiplos provedores de nuvem (Cloud) e o host local (Sovereign), garantindo redundância e performance.

## 🏛️ Estratégia de Provedores (Soberana)
- **Primary Inference**: OpenRouter (MiniMax/DeepSeek/Kimi).
- **Secondary / Backup**: Claude (Anthropic) / Gemini (Google).
- **Sovereign Host**: Ollama (Gemma 3) rodando na RTX 4090 local.

## 🏗️ Padrões de Integração
1. **Contexto Unificado**: Uso do `ai-context` para manter a memória síncrona entre provedores.
2. **Handoff Inteligente**: Definir quando trocar de agente (ex: Claude para implementação, Antigravity para interface).
3. **Escudo de Custos**: Priorizar execução local e modelos de baixo custo conforme a política do `cost-reducer`.

## 📍 Quando usar
- No setup inicial de novos repositórios ou sub-serviços.
- Para reconfigurar o gateway em cenários de outage de provedores externos.
- Para gerenciar a identidade da Aurélia em diferentes ecossistemas.
