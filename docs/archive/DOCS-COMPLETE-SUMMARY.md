---
title: DOCUMENTAÇÃO COMPACTA — Aurelia Elite Edition
description: Versão enxuta do diretório `docs/`, preservando referências críticas por área e status.
owner: codex
updated: 2026-03-21
---

# Visão Compacta da Documentação (Docs/)
Este artefato centraliza em uma única página os temas essenciais de `docs/`, mantendo o contexto e apontando para cada fonte original.

## 1. Governança & ADRs
- **Governança global:** `docs/ADR-20260319-GOVERNANCE-INDEX.md` concentra os índices de dados, operação, observabilidade, segredos e rede, marcando o que está concluído/adiado/planned.
- **ADRs de suporte:** `docs/adr/ADR-20260319-Polish-Governance-All.md` (governança industrial), `docs/adr/ADR-20260319-extensions-governance.md`, `docs/adr/ADR-20260319-voice-capture-runtime.md`, `docs/adr/ADR-20260319-state-memory-runtime.md` e os demais ADRs críticos listados em `docs/adr/TASKMASTER-INDEX.md` documentam decisões sobre voz, agentes e segurança.
- **Backlog e templates:** `docs/adr/TEMPLATE-SLICE.md`, `docs/adr/TEMPLATE-NONSTOP-SLICE.md` e `docs/adr/PENDING-SLICES-20260319.md` definem como criar novas ADRs nonstop.

## 2. Blueprint operacional e arquitetura
- **General blueprints:** `docs/aurelia_master_blueprint_20260319.md`, `docs/aurelia_general_blueprint_20260319.md`, `docs/homelab_jarvis_operating_blueprint_20260319.md`, `docs/homelab_tutor_v2_blueprint_20260319.md`, `docs/agent_swarm_dashboard_blueprint_20260319.md` e `docs/gateway_rollout_blueprint_20260319.md` descrevem a visão do homelab, roteamento, dashboard e rollout.
- **LLM e roteamento:** `docs/model_routing_matrix_20260319.md`, `docs/llm_gateway_blueprint_20260319.md`, `docs/local_model_kit_blueprint_20260319.md`, `docs/groq_stt_simulation_20260319.md` e `docs/gemini_fallback_runtime_20260319.md` detalham capacidades de STT, fallback e roteamento de modelos.
- **Arquitetura e stack:** `docs/ARCHITECTURE.md`, `docs/architecture.md`, `docs/memory-sync-architecture.md` e `docs/REPOSITORY_CONTRACT.md` explicam limites modulares, fluxos de memória e contrato de governança.

## 3. Operação e confiabilidade
- **Rotinas de operação:** `docs/operational-governance-health-checks.md`, `docs/operational-governance-backup-verification.md`, `docs/operational-governance-incident-response.md` conectam cron jobs, backups e runbooks P1–P4.
- **Checklist e travas:** `docs/ITEM-7-SYSTEMD-VERIFIED.md`, `docs/SMOKE_TEST_GUIDE.md`, `docs/SMOKE_TEST_README.md` documentam passos de verificação e testes smoke, enquanto `docs/SECRETS-MIGRATION.md`, `docs/⚠️-SECRETS-TODO-CRITICAL.md` cobrem migração e riscos de segredos.
- **Monitoramento:** `docs/observability-governance-metrics.md`, `docs/data-governance-lifecycle.md`, `docs/data-governance-store-selection.md`, `docs/data-governance-compliance.md` e `docs/qdrant-collection-contract.md` explicam métricas, compliance e políticas de ciclo de vida.

## 4. Guias & instruções
- **Operação de agentes:** `docs/guide-antigravity.md` e `docs/guide-claudecode-cli.md` expõem como interagir com Antigravity e CLIs.
- **Style & playbooks:** `docs/STYLE_GUIDE.md`, `docs/PROJECT_PLAYBOOK.md`, `docs/ROADMAP_PROVIDERS.md`, `docs/BENCHMARKS.md` e `docs/LEARNINGS.md` consolidam padrões de escrita, estratégias de projeto e aprendizados.
- **Voice plane:** `docs/aurelia_voice_profile_20260319.md`, `docs/jarvis_local_voice_blueprint_20260319.md`, `docs/groq_ptbr_audio_blueprint_20260319.md`, `docs/voice_study_sources_ptbr_20260319.md` (remanejar se deletado) descrevem a identidade vocal e os estudos que a suportam.

## 5. Resumo das ações recentes
- **Sync e contexto:** `docs/adr/20260319-sync-ai-context-como-regra-de-slice.md` contextualiza quando rodar `sync-ai-context` após alterações documentais.
- **Prioridades:** `docs/PENDING-SLICES-20260319.md`, `docs/adr/TASKMASTER-INDEX.md` e `docs/ADR-20260319-GOVERNANCE-INDEX.md` são os pivot points para acompanhar Slice/Wave/Status.

## Diretrizes de uso
1. Use este documento sempre que for necessário um overview rápido de `docs/` antes de abrir os arquivos originais.
2. Para cada área citada, siga o hyperlink correspondente para acessar a versão completa e mantê-la atualizada em português.
3. Ao revisar ou condensar novos arquivos, preserve o nível de detalhe técnico (especificações de tabelas, runbooks, política de ciclo) apenas dentro do documento original — este sumário serve só de navegação mínima.
