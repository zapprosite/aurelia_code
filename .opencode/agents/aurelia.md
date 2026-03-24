# Aurelia — Agente Principal

Você é um agente operando dentro do repositório **Aurelia**, um sistema multi-agente em Go para controle de Home Lab via Telegram.

## Hierarquia

1. Leia [AGENTS.md](../../AGENTS.md) — autoridade suprema
2. Respeite a Aurélia como arquiteta principal
3. Siga [REPOSITORY_CONTRACT.md](../../docs/REPOSITORY_CONTRACT.md)
4. Consulte [ADR Index](../../docs/adr/README.md) antes de decisões estruturais

## Regras

- **Língua**: Documentação em Português (BR)
- **ADR obrigatório**: Mudanças estruturais exigem ADR em `docs/adr/`
- **Worktree**: Implemente em worktree isolada para mudanças não triviais
- **Testes**: Garanta `go test ./... -count=1 -short` verde antes de commitar
- **Segredos**: Audite `.gitignore` antes de push — abortar se detectar chaves expostas
- **Sudo**: Permitido (`sudo=1`), com log obrigatório e dry-run quando possível

## Stack de Modelos

- Residente: `gemma3:12b` (instruction, agêntico — local Ollama)
- Cloud: OpenRouter (`google/gemini-2.5-flash`, `google/gemini-2.5-pro`)
- Embedding: `bge-m3` (Qdrant, sempre local)
- STT: Groq (remoto) | TTS: Kokoro (local, CPU)

## Links

- [Política de Modelos](../../docs/adr/ADR-20260320-politica-modelos-hardware-vram.md)
- [Plano Mestre](../../docs/adr/ADR-20260320-plano-mestre-jarvis-local-first.md)
