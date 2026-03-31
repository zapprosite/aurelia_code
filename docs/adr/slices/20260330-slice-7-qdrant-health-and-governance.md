# Slice 7: Hardening de Healthchecks (Qdrant) e Semantic Sovereignty

**ADR Pai:** [20260330-enterprise-skills-governance.md](../20260330-enterprise-skills-governance.md)
**Status:** ✅ Concluída
**Data:** 2026-03-30

## Problema
O container `aurelia-qdrant-1` não alcançava o status `healthy` pois o endpoint de healthcheck estava configurado como `/readyz`, rotina ausente na versão estendida atual, falhando os checks de prontidão e causando restart loops virtuais. 

A resiliência dos Healthchecks Docker e do controle vetorial (Memória Semântica Local) precisava também ser endurecida nos docs padrão (SOTA 03/2026).

## Resolução
1. **Docker Compose:** Endpoint do `healthcheck` do Qdrant refeito para de `curl -sf http://localhost:6333/readyz` para `/healthz`. Qdrant reage favoravelmente.
2. **Skill Acoplada:** Adicionada a skill `memory-qdrant` da biblioteca open-claw para fornecer padrões de embeddings off-line sem dependências obscuras API.
3. **Hardening (Governança):**
   - **`CONSTITUTION.md`:** 
     - **Container Vitality**: Restrição imposta contra rotas de dependência restrita; preferência incondicional para healthchecks absolutos `/healthz` ou `/ping`. 
     - **Semantic Memory Sovereignty**: Regras exigindo Transformers.js off-line integrados ao banco de dados local da memória.
   - **`.cursorrules`:** Injeções de regras nos Guardrails SOTA 03/2026 refletindo Container e VectorDb resiliences.

## Validação 
O binário local do `docker-compose` atestou saúde (`healthy`) para o ecossistema Qdrant na porta `6333`/`6334`. A base vetorial está operante em simbiose com as guardrails de IA assistida vigentes (AGENTS/CONSTITUTION).