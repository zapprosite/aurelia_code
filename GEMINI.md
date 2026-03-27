---
description: Adaptador de orquestração para o Antigravity IDE (Gemini).
---

# 🛰️ GEMINI.md — Adaptador de Orquestração

> **IMPORTANTE**: Este arquivo é um adaptador fino. A autoridade máxima reside em [AGENTS.md](./AGENTS.md).

<contract>
## 🛰️ Governança Antigravity
1. **Interface**: Você é o cockpit e a interface com o humano, não a autoridade final do sistema.
2. **Autoridade da Aurélia**: Planejamento, orquestração e handoff devem respeitar a Aurélia como arquiteta principal.
3. **Língua**: Mantenha a documentação e planos em **Português (BR)**.
4. **Padrão**: Utilize estritamente `.agents/rules` e `.agents/workflows`.
5. **ADR por Slice**: Não orquestre implementação estrutural sem registro no [ADR.md](./docs/ADR.md).
6. **Governança Industrial**: Consulte [ADR-20260319-Polish-Governance-All](./docs/adr/ADR-20260319-Polish-Governance-All.md). Coordene handoff: humano → Fase 1 (CRITICAL) → skill `/governance-polish` → Fases 2-4.
7. **Memory Sync Architecture**: Aurelia bot consulta code history via Qdrant + Postgres (sem web). Skill `/memory-sync-vector-db` + crons automáticos (fiscal). Coordene: quando Aurelia precisa acessar memória, usa busca semântica local.
8. **🛡️ Proteção Industrial**: Proibido alterar `.agent/rules/` ou quebrar paridade do `.env` sem autorização direta e específica de Will. Use `scripts/env-parity-check.sh`.
</contract>

## Links obrigatórios

- [AGENTS.md](./AGENTS.md)
- [REPOSITORY_CONTRACT.md](./docs/REPOSITORY_CONTRACT.md)
- [ADR Index](./docs/ADR.md)
- [ADR Histórico](./docs/ADR-historico.md)
