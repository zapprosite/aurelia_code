# PROJECT PLAYBOOK

## Antigravity Delegation Rule

Neste projeto, a Aurelia deve tratar o chat do Antigravity como um copiloto leve.

### Use o chat leve para

- pequenas configuracoes locais
- pesquisa curta
- rascunho de patch
- localizar arquivos, opcoes e comandos
- revisar instrucoes antes da execucao principal

### Nao use o chat leve para

- segredos
- deploy
- operacoes de rede
- merges
- refactors grandes
- mudancas sem plano de verificacao

## Prompting Standard

Toda instrucao para o chat do Antigravity deve declarar:

- objetivo
- alvo exato
- restricoes
- prova exigida
- formato de resposta

## Audio Rule

Para STT com Groq:

- modelo padrao: `whisper-large-v3-turbo`
- linguagem padrao: `pt`
- temperatura padrao: `0`
- transcripts devem ser persistidos localmente antes de qualquer enriquecimento

## Memory Rule

Quando o chat leve ajudar de forma util:

- capture a decisao no repo
- promova o aprendizado para skill, playbook ou `.context/`
