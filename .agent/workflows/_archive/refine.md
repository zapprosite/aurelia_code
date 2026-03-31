---
description: Inicia o loop de refinamento de Markdown + Voz Kokoro no Telegram.
---

# Workflow: Refinamento de Comunicação (/refine)

Utilize este workflow para iterar sobre a qualidade das respostas da Aurélia.

## Passos

1. **Defina o Texto**: Escreva a mensagem que deseja testar (Markdown suportado).
2. **Execute o Script**:
   ```bash
   ./.agent/skills/communication-refinery/refine.sh "Sua mensagem aqui"
   ```
3. **Valide no Telegram**:
   - Confira a renderização das tags Markdown (Bold, Italics, Code).
   - Verifique a naturalidade da voz Kokoro no arquivo de áudio.
4. **Itere**: Ajuste o texto e rode o script novamente até o resultado ser excelente.

// turbo
5. Execute um teste de fumaça:
   ```bash
   ./.agent/skills/communication-refinery/refine.sh "Teste de fumaça da Aurélia: Markdown em *negrito* e voz Kokoro ativa."
   ```
