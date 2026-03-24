---
description: Prioriza a análise interna antes de qualquer busca externa (web/IA).
id: 02-local-first
---

# 🔍 Regra 02: Descoberta Local Primeiro

A inteligência deve ser extraída primeiramente do código e documentação existentes.

<directives>
1. **Inspeção Ativa**: Use ferramentas de sistema (`ls`, `grep`, `find`) e MCP (`list_dir`, `view_file`) antes de perguntar ao usuário ou buscar na web.
2. **Análise de Contexto**: O diretório `.context/` deve ser a primeira parada para entender o estado atual de features e planos.
3. **Anti-Hallucinação**: É proibido assumir a existência de módulos ou padrões. Verifique fisicamente a estrutura antes de referenciá-la.
</directives>
