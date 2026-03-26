# 16-secrets-governance.md — Governança de Espelhamento e Segredos

> **Autoridade**: Senior Industrial (2026-03-26)
> **Escopo**: Todos os agentes operando no monorepo Aurélia.

## 🛰️ Regra de Ouro: Paridade Total

1. **Espelhamento Obrigatório**: Os arquivos `.env` e `.env.example` DEVEM ser espelhos estruturais perfeitos (100% de paridade de chaves).
2. **Visibilidade de Estrutura**: O `.env.example` deve conter comentários e placeholders descritivos para que qualquer LLM entenda quais serviços e chaves compõem a infraestrutura, sem ver os valores reais.
3. **Opacidade de Valores**: O `.env` é de leitura exclusivamente humana e para injeção em runtime pelo sistema.
4. **Proibição de Deleção**: É terminantemente PROIBIDO para agentes deletar ou mover o arquivo `.env`. Esta é uma operação de autoridade humana exclusiva.
 NUNCA exponha o conteúdo do `.env` em logs, chats ou artefatos compartilhados.

## 🛡️ Guardrails Executáveis

- **Placeholder Padrão**: Utilize `{chave-para-env}` para sinalizar segredos faltantes ou placeholders universais em arquivos de configuração (`app.json`, etc).
- **Proteção Git**: O arquivo `.env` deve constar obrigatoriamente no `.gitignore`. Agentes devem abortar se detectarem que o `.env` está prestes a ser rastreado.
- **Auditoria de Drift**: Antes de qualquer major release ou handoff, valide a paridade entre `.env` e `.env.example`.

## 🤖 Visibilidade para IA (LLM-First)

Para que um LLM entenda o ambiente via `.env.example`:
- Use comentários `#` para agrupar categorias (LLMs, DBs, Infra).
- Use placeholders explicativos quando o nome da chave for ambíguo.
- Exemplo: `CF_TUNNEL_ID="ID_DO_TUNNEL_CLOUDFLARE_GERADO_NA_CRIACAO"`.

---
*Documento protegido sob Política Zero Hardcode 2026.*
