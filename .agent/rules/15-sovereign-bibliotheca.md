# Governança Sovereign-Bibliotheca v2

> **Status**: Ativo (Março 2026)
> **Objetivo**: Manter a ordem arquitetural entre o motor Go (Aurélia), o ecossistema Node (Obsidian/OpenClaw) e a orquestração Bash.

## 1. Princípio da Soberania da Camada (Layer Sovereignty)

Para evitar a "salada" técnica, seguimos a regra dos três pilares:

- **Core (Go)**: Responsável por performance, concorrência, drivers de hardware e lógica de estado persistente (`internal/`).
- **Automação (Node/JS/Python)**: Responsável pela interface de usuário (Obsidian), scripts de terceiros e skills dinâmicas (`homelab-bibliotheca/skills/`).
- **Orquestração (Bash)**: A "cola" universal que unifica os dois mundos via CLI e contratos JSON (`homelab-bibliotheca/lib/`).

## 2. Regras de Interoperabilidade

### ⚖️ Contrato de Linguagem
NUNCA importe logicamente uma linguagem na outra (ex: `cgo` ou `node-gyp` para conectar os mundos).
- Use **CLI** (Scripts Bash em `lib/`) para tarefas síncronas.
- Use **REST API** (Aurélia no `8484`) para tarefas assíncronas ou de rede.

### 📝 Contrato de Dados
Toda troca de dados entre Go e Node DEVE ser via:
1. **JSON**: Validado via schemas Go (internal/) ou Zod quando o pacote TypeScript for implementado.
2. **Markdown**: Para persistência de conhecimento humano e notas.

## 3. Manutenção da Biblioteca

### 🔄 Sincronização (Master Sync)
O script `homelab-bibliotheca/sync.sh` é a autoridade máxima de integridade. Ele DEVE ser executado após:
- Importação de novas skills.
- Mudanças estruturais no SQLite ou Qdrant.
- Migrações do Supabase.

### 📂 Higiene de Diretórios
- **Proibido** criar arquivos de script soltos na raiz. Use `homelab-bibliotheca/lib/`.
- **Proibido** duplicar segredos. Use o arquivo central `~/.aurelia/config/secrets.env` carregado pelo `config.sh`.

## 4. Gestão de Memória e Skills
- **Skills** são ativos vivos. O `skills-registry.json` deve ser mantido atualizado para que os agentes possam "descobrir" novas competências autonomamente via `skills.sh manifest`.

---
*Assinado: Aurélia (Arquiteta Líder) & Antigravity (Operador IDE)*
