# ADR 20260327: Master Skill Global (Bootstrap 2026) 🛠️🌍

## Status
Proposto

## Contexto
O ecossistema Aurélia cresce em múltiplos repositórios e workspaces. Atualmente, o setup de um novo projeto com as regras de governança e kits de IA (BMad, Spec-Kit) é manual e propenso a erros.

## Decisão
Implementar a **Master Skill Global**, uma ferramenta de CLI e skill de agente que automatiza:
1.  **Init**: Criação de `.context/`, `.agent/rules/` e arquivos de governança base.
2.  **Kit Install**: Instalação do `Antigravity Kit` e `BMad Method` via link simbólico ou cópia controlada.
3.  **Skill Sync**: Importação de skills de um diretório central (`~/.aurelia/skills`) para o projeto local.

## Arquitetura
- **Comando**: `/master-skill init` / `/master-skill install {kit}`.
- **Engine**: Script Python/Go resiliente que valida a estrutura do monorepo antes de agir.
- **Soberania**: Permite que qualquer pasta se torne um nó do cluster Aurélia em segundos.

## Consequências
- **Positivas**: Padronização absoluta, velocidade de escala, facilidade de manutenção de regras globais.
- **Negativas**: Necessidade de manter retrocompatibilidade com versões antigas dos kits.
