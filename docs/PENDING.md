# PENDING — Slices Aguardando Implementação
Última auditoria: 01/04/2026

> POLÍTICA: Apenas slices com código no repo são marcados ✅.
> Este arquivo é a única fonte de verdade do backlog.

## P0 — Hotfix Ativo 🔥 (impede o bot de funcionar)

Nenhum. Sistema SOTA 2026 Sovereign operando de forma estável.

## P1 — Crítico 🔴 (infra sem a qual o home lab não é estável)

Nenhum limitador conhecido no momento.

## P2 — Alto 🟡 (qualidade da experiência)

Nenhum pendente.

## P3 — Médio 🟢 (expansão da capacidade nativa)

Nenhum pendente.

---

## ✅ Slices Concluídos Recentes (Sovereign 2026)
- **S-51 a S-55**: Estabilização do Homelab, TTS via Edge, Groq STT, remoção do Supabase.
- **S-56**: Extração do `AURELIA.md` (System Prompt injetado em runtime).
- **S-57**: Setup de Feature Flags via env (FEAT_VOICE, FEAT_DREAM, FEAT_KAIROS).
- **S-58**: Implementação do `DreamConsolidator` (SQLite -> LLM -> Qdrant).
- **S-59**: Refatoração do `input_pipeline.go` (unidades lógicas < 100 linhas).
- **S-60**: `SharedMemory` Redis Pub/Sub e infraestrutura de Swarm Local.
- **S-61**: Sincronização definitiva do `PENDING.md`.
- **S-62**: Correção de cron expression para 6 campos (repo_guardian).
- **S-63**: Persistência de `IDENTITY.md` via volume Docker RO.
- **S-64**: Persistência de `HEARTBEAT.md` via volume Docker RO.
- **S-65**: Remoção total de referências residuais ao Supabase.
- **S-66**: Fix `OllamaURL` default e estabilização do roteamento LiteLLM local.
- **S-67**: Ollama models → tank/models (OLLAMA_MODELS=/srv/models)
- **S-68**: SearXNG como tool search_web_local na Aurélia
