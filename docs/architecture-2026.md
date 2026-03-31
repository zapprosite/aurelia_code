# Arquitetura Soberana (Aurelia SOTA 2026)

> **Autoridade**: Supremacia Arquitetural
> **Filosofia**: Pinned Data Center & Zero Drift

Este documento mapeia a arquitetura definitiva do ecossistema Aurélia. Se você é um agente operando neste repositório, este é o único mapa da mina que você precisa ler para entender como a máquina funciona no nível macro.

## 1. O Padrão "Pinned Data Center"
A fundação do Homelab (Ubuntu 24.04). Todo dado sensível, estado e memória semântica está "ancorado" localmente. O Cloud é usado estritamente como poder computacional volátil.

- **Redis (State & Sync)**: O coração transacional.
  - *Deduplicação de Eventos*: `telegram:seen:{chat_id}:{msg_id}` com `SetNX` impede que retries disparem execuções paralelas.
  - *Caches de LLM*: Retornos pesados do Porteiro (Polimento) são cacheados em hashes SHA256 por 24h para zerar latência em re-execuções.
- **Qdrant (Vector Memory)**: A memória de longo prazo da Aurélia. Conversas são comprimidas e vetorizadas localmente.
- **SQLite / Postgres**: Persistência estruturada do negócio.

## 2. A Camada "Porteiro" (Sentinela Sênior)
O middleware que blinda o sistema.

- **Input Guardrail**: Analisa a intenção da mensagem (via Qwen) antes de tocar no roteador principal. Bloqueia jailbreaks e injeções.
- **Output Guardrail (Secrets)**: Interceptação determinística baseada em Regex/Heurística (`sk-`, `ghp_`, `AUR_`). Nenhuma string sensível cruza o túnel do Telegram.
- **Polisher (SOTA 2026)**: A saída crua (JSON, XML de LLMs) é interceptada e convertida para Markdown amigável no padrão da assistente, com cache em Redis.

## 3. Roteamento de Modelos (Tiered Fallback)
A Aurélia não confia em uma única API. O roteamento é resiliente e desce a escada (Fallback) automaticamente:

1. **Tier 1 (Cloud Orchestrator)**: Gemini 1.5 Pro / Claude 3.5 Sonnet. Puxa a carga pesada de raciocínio lógico e orquestração.
2. **Tier 2 (OpenRouter/Fallback)**: Rota secundária para quando as APIs primárias falham ou atingem Rate Limits.
3. **Tier 0 (Soberano/Local)**: Ollama (Qwen 0.5b / 7b) rodando na GPU. Se a internet cair, o Sentinel e o Polisher continuam vivos e bloqueando acesso graças aos modelos leves rodando no bare-metal.

## 4. Governança Multi-Agente (Zero Noise)
Para evitar que múltiplos LLMs (Gemini, Claude, OpenCode) entrem em colapso devido a vazamento de prompt:

- **Adaptadores na Raiz**: `AGENTS.md`, `GEMINI.md` e `CLAUDE.md` são vazios de lógica. Eles agem como routers apontando para `.agent/rules/`.
- **Regras Injetáveis**: A pasta `.agent/rules/` contém `.md` hiper-densos que são lidos on-demand dependendo de quem pega a task (Ex: Gemini lê `gemini.md` focado em arquitetura; Claude lê `claude.md` focado em bash e containers).

## 5. O Fluxo de Execução (O Pipeline Perfeito)
1. **Telegram Retry** -> **Redis Deduplication** (Bloqueio).
2. **Payload Autorizado** -> **Porteiro IsSafe** (Qwen valida intenção).
3. **Módulo de Agente** -> **Tier 1 LLM** (Raciocínio & Ferramentas).
4. **Resposta Gerada** -> **Porteiro Polisher** (Markdown + Emoji).
5. **Saída Bruta** -> **Porteiro SecretMask** (Censura de chaves).
6. **Delivery** -> Usuário via Telegram API.

---
*Assinado: Arquiteto de Sistemas (Antigravity SOTA 2026)*
*Resiliência é a única métrica que importa.*
