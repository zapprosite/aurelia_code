---
description: Automação de pesquisa profunda via Gemini Web para suporte técnico.
id: deep-researcher
---

# 🔬 Skill: Deep-Researcher

Automatiza a descoberta de conhecimento externo para subsidiar decisões arquiteturais e debugging.

<directives>
1. **Verificação de Ambiente**:
   - Navegue para `https://myaccount.google.com/email`.
   - Confirme se a conta ativa é a conta de desenvolvimento autorizada para este workspace.
   - Se a conta for pessoal ou incorreta, solicite via `notify_user` que o usuário troque de perfil no navegador.
2. **Navegação de Elite**:
   - Use `https://gemini.google.com` com o modelo mais avançado disponível.
   - Ative o modo **Deep Research** para tópicos complexos.
3. **Parâmetros de Pesquisa**:
   - Priorize documentação oficial, benchmarks de performance e artigos técnicos (`arXiv`, `Google Scholar`).
   - Extraia exemplos de código (snippets) prontos para integração.
4. **Persistência de Conhecimento**:
   - Salve o relatório em `~/Desktop/pesquisas-gemini/` com o formato `YYYY-MM-DD_topico.md`.
</directives>

## Fluxo de Trabalho
1. Definir o objetivo da pesquisa.
2. Executar a navegação via ferramenta `browser`.
3. Consolidar os resultados em um artefato no repositório ou no Desktop.
