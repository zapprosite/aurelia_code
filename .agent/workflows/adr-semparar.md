# /adr-semparar
Description: Abre uma slice em modo nonstop com ADR + JSON de continuidade estilo taskmaster.

---
1. Use a skill `adr-nonstop-slice`.
2. Gere o par de artefatos com `./scripts/adr-slice-init.sh <slug> --title "<Title>"`.
3. Preencha o ADR em `docs/adr/ADR-YYYYMMDD-slug.md`.
4. Preencha o JSON em `docs/adr/taskmaster/ADR-YYYYMMDD-slug.json`.
5. Registre:
   - objetivo
   - escopo
   - próximos passos
   - smoke/simulações com `curl`, `go test`, scripts e fallback
   - handoff/resume prompt
6. Ao fechar a slice, rode `./scripts/sync-ai-context.sh`.
