# Regra: Sovereign 2026 (Segurança de Elite)

> **CRITICAL**: Esta regra é a autoridade máxima de segurança no monorepo.

1. **Zero Hardcode**: É estritamente proibido incluir segredos, tokens ou senhas em texto claro.
2. **Placeholders**: Utilize obrigatoriamente `{chave-para-env}` para referenciar segredos que serão injetados via ambiente.
3. **Auditoria Proativa**: O agente deve validar o conteúdo antes de qualquer commit/push.
4. **Resgate de Soberania**: Em caso de detecção de segredos expostos, o agente deve abortar a operação e alertar o usuário.
