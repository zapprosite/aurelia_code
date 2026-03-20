# Blueprint: Antigravity Gemini Operator

## Objetivo

Transformar o chat leve do Antigravity em um copiloto de microexecucao para a Aurelia, sem confundir papeis.

## Papel de Cada Camada

- `Aurelia CLI / Codex / Claude`
  - execucao tecnica principal
  - mudancas estruturais
  - verificacao real
  - testes
  - commits
- `Antigravity Chat / Gemini Flash / Minimax 2.7`
  - pesquisa curta
  - pequenas configuracoes
  - edicao localizada
  - rascunho de comandos
  - resumo de alternativas

## Regras de Roteamento

Use o chat do Antigravity quando a tarefa for:

- pequena
- reversivel
- local
- de baixo risco
- dependente de leitura contextual rapida

Nao use o chat do Antigravity como executor primario quando envolver:

- secrets
- deploy
- rede
- remocao destrutiva
- merges
- migracoes
- alteracoes multi-arquivo acopladas

## Contrato de Prompt

Toda instrucao enviada ao chat leve deve conter:

1. objetivo em uma frase
2. arquivo ou alvo exato
3. restricoes
4. prova esperada
5. formato de saida

Template:

```text
Contexto:
- Workspace: /home/will/aurelia
- Data: 2026-03-19
- Papel: voce esta no chat leve do Antigravity

Objetivo:
- <o que precisa ser feito>

Alvo:
- <arquivo, servico, comando ou area exata>

Restricoes:
- nao tocar em secrets
- nao fazer deploy
- nao inventar arquivos fora do escopo
- nao declarar sucesso sem prova

Saida esperada:
- diff proposto ou instrucoes exatas
- comandos de validacao
- risco residual em uma linha
```

## Flow Operacional

1. Aurelia classifica a tarefa
2. se for leve, delega ao chat do Antigravity com prompt estruturado
3. o chat devolve diff, comandos ou pesquisa curta
4. a Aurelia valida
5. se aprovavel, a Aurelia executa ou converte em mudanca real
6. evidencia entra no repo ou `.context/`

## Casos Bons

- ajustar uma config JSON pequena
- localizar onde um MCP e configurado
- montar um comando `curl`
- resumir diferenca entre duas abordagens
- preparar checklist de validacao

## Casos Ruins

- alterar pipeline inteiro sozinho
- operar secrets
- subir servico em producao
- fazer merge ou rebase
- executar limpeza agressiva

## Integracao com Audio/Groq

O chat leve pode ajudar a:

- revisar payload de STT
- propor filtros de silencio
- sugerir chunking de audio
- revisar persistencia em `Supabase` e `Qdrant`

Mas a execucao real continua com a Aurelia CLI.
