# Chat Prompts

## Prompt base

```text
Voce esta no chat leve do Antigravity para o workspace /home/will/aurelia.

Tarefa:
- <descreva em uma frase>

Alvo exato:
- <arquivo, modulo, config, servico ou endpoint>

Restricoes:
- nao tocar em secrets
- nao fazer deploy
- nao declarar sucesso sem prova
- manter a mudanca pequena e reversivel

Entregue:
- diff proposto ou passos exatos
- comandos de validacao
- risco residual em uma linha
```

## Prompt para pequena configuracao

```text
Revise a configuracao alvo e proponha a menor mudanca possivel.
Nao reestruture o projeto.
Se houver ambiguidade, liste no maximo 2 opcoes e recomende 1.
```

## Prompt para pesquisa curta

```text
Pesquise apenas o necessario para responder ao ponto tecnico.
Nao escreva tutorial longo.
Responda com:
- conclusao
- prova
- proximo comando ou diff
```

## Prompt para revisar diff

```text
Leia o diff abaixo e diga:
- o que mudou
- risco principal
- validacao minima necessaria
```
