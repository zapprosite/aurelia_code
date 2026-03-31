# CONSTITUTION.md: Sovereign Enterprise 2026 🛡️

## 1. Missão
O ecossistema Aurélia é uma infraestrutura de IA soberana, local-first e resiliente, projetada para operar como uma appliance inteligente no Ubuntu 24.04 LTS.

## 2. Princípios de Engenharia
-   **Everything is a Service**: Todos os componentes críticos devem ser executados como daemons (systemd), garantindo auto-recuperação (Self-Healing).
-   **Local Sovereignty**: Funcionalidades core (VAD, STT, LLM, TTS) devem operar 100% offline ou em infraestrutura privada (Homelab).
-   **Zod-First Contract**: Esquemas de dados são a única fonte de verdade para validação I/O.
-   **ADR-Driven**: Toda mudança estrutural ou decisão arquitetural deve ser precedida por um Registro de Decisão Arquitetural (ADR).
-   **Container Vitality (03/2026)**: Healthchecks devem usar rotas universais e testadas (ex: `/healthz`, `/`, `ping`) para garantir resiliência imediata do cluster.
-   **Semantic Memory Sovereignty**: Bancos vetoriais (ex: Qdrant) devem rodar estritamente locais, integrados via embeddings off-line (Transformers.js/HuggingFace), sem dependência de chaves de API externas.
-   **Quality-Guided Routing (03/2026)**: Tarefas que exigem conformidade SOTA (Markdown 2026, nuances PT-BR, Audio/TTS) devem preferir modelos de Tier 1/Paid (`aurelia-top`) em vez do Tier 0 local, priorizando a fidelidade sobre a latência.

## 3. Padrões de Segurança
-   **Secret Isolation**: Proibido o uso de segredos hardcoded. Uso obrigatório de `EnvironmentFile` ou Vault.
-   **Auditability**: Logs estruturados em `/var/log/aurelia/` com rotação industrial.

## 4. Governança de Workflow (SDD)
Adotamos o **Spec-Driven Development (SDD)** via `spec-kit`:
1.  **Constituição** (Princípios)
2.  **Especificação** (Requisitos)
3.  **Plano** (Arquitetura)
4.  **Tarefas** (Execução)

---
*Assinado: Aurélia (SOTA 2026.03.30)*
