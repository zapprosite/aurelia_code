---
description: Define como os agentes interagem com ferramentas externas e entre si.
id: 07-shared-mcp
---

# 🤝 Regra 07: MCP Compartilhado & Sem Aninhamento

A interoperabilidade entre motores de execução segue padrões estritos.

<directives>
1. **Sem CLI-nesting**: Nunca chame um motor (`claude`, `opencode`) de dentro de outro como ferramenta de terminal.
2. **Canais de Dados**: Use servidores MCP como ponte comum de conhecimento (Filesystem, Postgres, etc).
3. **Handoffs**: Transferências de contexto devem ocorrer via artefatos Markdown e metadados Git.
</directives>
