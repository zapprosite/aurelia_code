---
description: Define a fonte primária de verdade e a hierarquia de decisão.
id: 01-authority
---

# 📜 Regra 01: Autoridade Única

Este repositório opera sob um modelo de autoridade centralizada.

<directives>
1. **Fonte de Verdade**: O arquivo `AGENTS.md` na raiz é a constituição soberana. Nenhuma regra local pode contradizê-lo.
2. **Orquestração**: O Antigravity IDE é o supervisor global. Motores como Claude, Gemini e OpenCode atuam como braços executores sob supervisão.
3. **Escopo de Decisão**: Mudanças estruturais exigem registro via ADR em `docs/adr/` (Architectural Decision Record) e validação conforme `plan.md`.
</directives>
# CORE RULES (The Sovereign Laws)
**Version: SOTA 2026 | Authority: Absolute**

Estas são as leis universais do Homelab Aurélia. Nenhum agente, sob nenhuma circunstância, pode violar estas regras.

## 1. Zero Idioma Estrangeiro
- O idioma padrão universal é **Português do Brasil (PT-BR)** para todo planejamento, documentação, ADRs e respostas no chat.
- **Exceção**: Nomenclatura técnica dentro do código fonte (Go, TS, Python) e CLI commands.

## 2. Zero Hardcode e Segurança Máxima
- É TERMINANTEMENTE PROIBIDO armazenar chaves de API, senhas ou tokens em código fonte, documentação (ex: YAML, MD) ou configurações que sejam rastreadas pelo Git.
- Todos os segredos devem ser lidos de arquivos `.env` ou injetados em runtime via variáveis de ambiente (`os.Getenv()`).
- O arquivo `.env` NUNCA DEVE SER COMITADO.

## 3. Padrão "Pinned Data Center" (2026)
- Todo processamento e ciclo de vida de dado sensível (Memória Vetorial, Chat Logs, Caches) deve acontecer LOCALMENTE no Ubuntu Homelab ou persistir no Redis/Qdrant da aplicação.
- Componentes externos (Gemini, Claude, Groq) atuarão APENAS como força computacional (cérebro sem memória). A memória é Soberana.

## 4. Ordem e Hierarquia Estrutural
- Nunca destrua o arquivo `.agent/rules/core.md` ou qualquer documento base (como `docker-compose.yml`) sem ordem *explícita* e repetida de mudança estrutural.
- Mudanças Arquiteturais significativas (Banco de Dados, Rotas Core, Configurações de Rede) devem SEMPRE ser precedidas por uma ADR em `docs/adr/`.
