# Security Skill (Elite Level)

## Objetivo
Identificar, prevenir e mitigar riscos de segurança em todo o monorepo, garantindo conformidade com os padrões de "Elite Repo" e as Regras Globais de Governança.

## Princípios de Elite (Obrigatórios)
- **Zod-First**: Toda validação de input e esquema de dados deve residir em `packages/zod-schemas/`. Proibido duplicar lógica.
- **Segurança Proativa**: Auditoria de segredos e `dry-run` para mudanças de infraestrutura são mandatórios.
- **Observabilidade**: Logging estruturado e tratamento de erro padronizado em todos os módulos core e endpoints tRPC.

## Quando usar
- Revisão de segurança em PRs (Tier B/C).
- Auditoria de segredos antes de `git push`.
- Implementação de novos recursos de infraestrutura ou banco de dados.
- Configuração de autenticação/autorização no monorepo.

## Como executar
1. **Contexto**: Inspecione `AGENTS.md` e `RULE[user_global]` antes de qualquer ação de segurança.
2. **Superfície de Ataque**: Identifique se o foco é: Auth, DB/Deps, Desktop ou Web.
3. **Checklist**: Aplique as diretrizes dos sub-arquivos técnicos.
4. **Remediação**: Documente vulnerabilidades com severidade, impacto e plano de correção seguindo o padrão ADR se necessário.

## Arquivos de Referência (Manual Técnico)
- [auth-and-secrets.md](file:///home/will/Remote-control/.agents/skills/security/auth-and-secrets.md): Gestão de credenciais, JWT e auditoria de segredos.
- [database-and-deps.md](file:///home/will/Remote-control/.agents/skills/security/database-and-deps.md): Zod schemas, ORM security e supply chain.
- [desktop-security.md](file:///home/will/Remote-control/.agents/skills/security/desktop-security.md): Segurança Electron e isolamento de contexto.
- [web-security.md](file:///home/will/Remote-control/.agents/skills/security/web-security.md): OWASP Top 10, headers e observabilidade tRPC.

## Severidade e Governança
- **Crítica (Tier C)**: Vazamento de segredos, acesso root ou PII exposta. Exige correção imediata e ADR.
- **Alta (Tier B)**: Bypass de auth, Injeção (SQL/Command) ou falhas de isolamento.
- **Média/Baixa**: Headers ausentes, logs verbosos ou dependências desatualizadas.
