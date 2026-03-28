# Governança Sovereign-Bibliotheca SOTA 2026.1

> **Status**: Cristalizado (Março 2026)
> **Objetivo**: Gestão centralizada de ativos de inteligência via Obsidian + Qdrant.

## 1. Princípio da Centralização Soberana

- **Autoridade**: A descoberta de skills é feita pelo **Qdrant Index** (L3-Memory).
- **Interface**: A gestão humana e revisão de regras é feita via **Obsidian CLI/Vault**.
- **Sync**: O script `obsidian-sync.sh` deve ser rodado após qualquer alteração estrutural em `.agent/`.

## 2. Contrato de Skills

- **Sovereign First**: Skills em `.agent/skills/` têm precedência sobre a biblioteca em `homelab-bibliotheca`.
- **Vetorização**: Toda nova skill deve ser vetorizada para ser visível ao `Aurelia Hub`.

---
*Assinado: Aurélia (Arquiteta Líder) & Antigravity (Operador SOTA 2026.1)*
