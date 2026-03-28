# ADR-20260327: Poda Estrutural Soberana

## Status: Proposta

## Contexto
Após a industrialização de 27/03/2026, o monorepo foi purificado em nível de logging e dependências. No entanto, restam detritos de 2024/2025 (arquivos `.go` em `node_modules`, stubs vazios, logs antigos e diretórios órfãos) que poluem o contexto dos agentes e aumentam a latência de inferência.

## Decisão
Implementar uma "Poda Segura" (Safe Pruning) para remover artefatos não-soberanos e garantir que o repositório contenha apenas código Go funcional e documentação SOTA.

### Critérios de Exclusão (Poda):
1. **Arquivos `.go` Infiltrados**: Qualquer arquivo Go fora de `/internal`, `/pkg`, `/cmd` ou scripts de bootstrap na raiz.
2. **Diretórios Vazios**: Pastas resultantes de refatorações ou experimentos concluídos (ex: `internal/aurl`, `internal/purity` se estiverem vazios).
3. **Artefatos de Build**: Limpeza forçada de binários locais que não seguem o padrão `/bin`.
4. **Logs Efêmeros**: Remoção de `.log` e arquivos `/tmp` locais que não sejam rotacionados automaticamente.

## Consequências
- **Positivas**: Redução drástica no uso de tokens por sessão, builds mais rápidos e eliminação de "falsos positivos" em buscas globais (grep/find).
- **Negativas**: Pequeno risco de remover templates de ferramentas externas que usem extensão `.go` indevidamente ( mitigado por `dry-run`).

---
**Data**: 27/03/2026
**Autor**: Antigravity (SOTA 2026)
