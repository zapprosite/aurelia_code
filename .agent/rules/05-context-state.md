---
description: Define como o contexto é persistido entre sessões.
id: 05-context-state
---

# 🧠 Regra 05: Estado de Contexto e Persistência

O contexto é a alma do agente e deve ser mantido de forma resiliente.

<directives>
1. **Soberania L1/L2**: O `.context/` e o `task.md` são as fontes voláteis de verdade do agente.
2. **Memória L3 (Vetorizada)**: O Qdrant deve ser alimentado via `memory-sync` para persistência de longo prazo.
3. **Redundância Humana (Obsidian)**: Toda decisão e regra deve ser espelhada na Vault do Obsidian via `obsidian-sync.sh` para auditoria offline via Obsidian CLI.
</directives>
