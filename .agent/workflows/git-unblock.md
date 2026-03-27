---
description: Workflow rápido para invocar a skill de destravamento de Git zumbi.
---

# /git-unblock

Description: Workflow de emergência para destravar o repositório quando comandos do Git param de responder (timeouts, locks presos). Invoca a skill `git-unblocker`.

---

1. O agente identifica que houve timeout no terminal ou erro de `index.lock`.
2. O agente chama imediatamente a skill `git-unblocker` para conhecimento e contexto.
3. Executa o one-liner de limpeza agressiva (SIGKILL em git/gh + remoção de locks).
4. Emite um `git status` para validar o restabelecimento.
5. Retoma o fluxo original interrompido, mas agora em passos isolados de terminal.
