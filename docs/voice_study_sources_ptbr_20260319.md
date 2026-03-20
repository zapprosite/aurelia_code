---
title: Fontes de Estudo de Voz PT-BR
status: active
date: 2026-03-19
---

# Objetivo

Registrar fontes abertas e licenciadas úteis para estudar prosódia, dicção e estilo de voz em português do Brasil, sem depender de vozes clonadas de terceiros.

# Fontes recomendadas

## Mozilla Common Voice PT

- uso: estudo de sotaque, dicção e variação
- observação: não usar para tentar identificar pessoas
- link: https://datacollective.mozillafoundation.org/datasets/cmflnuzw6rbhzx9lai7ha3emb

## TTS-Portuguese-Corpus

- uso: estudo de TTS em português brasileiro
- link: https://github.com/Edresson/TTS-Portuguese-Corpus

## CORAA

- uso: estudo e benchmarking de fala em PT-BR
- observação: licença mais restritiva para derivados
- link: https://github.com/nilc-nlp/CORAA
- paper: https://arxiv.org/abs/2110.15731

## FalaBrasil Speech Datasets

- uso: índice curado de bases de áudio transcrito em PT-BR
- link: https://github.com/falabrasil/speech-datasets

# Regra da Aurelia

- essas fontes servem para **estudo**
- a voz oficial da Aurelia deve vir de **amostra local autorizada**
- não usar “vozes clonadas prontas” de terceiros como base da voz oficial

# Decisão prática

A voz `aurelia-ptbr-formal-doce-v1` deve ser construída a partir do arquivo local autorizado entregue pelo usuário:

- `/home/will/aurelia/clone-voz/aurelia.mp3`

Quando houver `MINIMAX_API_KEY`, o caminho será:

1. validar a amostra
2. criar ou ativar a voz autorizada na MiniMax
3. rodar smoke do `voice_id`
4. só então trocar o TTS do Telegram
