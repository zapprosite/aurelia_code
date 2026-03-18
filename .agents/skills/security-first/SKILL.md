---
name: security-first
description: Elite skill for security-first.
---



# 🛡️ Skill: Security-First

Esta habilidade garante que o workspace permaneça seguro e livre de vazamentos de credenciais.

<directives>
1. **Secrets Scanning**: Antes de qualquer commit ou push, verifique se há chaves de API, tokens JWT ou senhas em texto puro.
2. **Permissões de Arquivo**: Valide se arquivos sensíveis (como `.env`, `.key`) não estão com permissões excessivas (ex: deve ser 600 ou 400).
3. **Auditoria de Código**: Procure padrões de risco como Injeção de SQL, XSS ou falhas de lógica em permissões.
4. **Deny List**: Respeite rigorosamente os arquivos ignorados no `.gitignore` e as regras de segurança do sistema.
</directives>

## Checklist de Pré-Commit
- [ ] Nenhum segredo detectado no `git diff --cached`.
- [ ] `.env` e chaves privadas estão no `.gitignore`.
- [ ] Novas dependências foram verificadas contra vulnerabilidades conhecidas.
