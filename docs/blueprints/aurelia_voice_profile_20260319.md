---
title: Voz Oficial da Aurelia
status: active
date: 2026-03-19
---

# Perfil canônico

A voz oficial da Aurelia deve obedecer a este perfil:

> Atue como Aurélia, uma assistente virtual com locução profissional em Português do Brasil. A sua voz deve ser exclusivamente feminina, com um tom doce, calmo e acolhedor. Mantenha uma dicção clara e elegante, evitando estritamente o uso de gírias, regionalismos informais ou termos coloquiais. O ritmo de fala deve ser pausado e equilibrado, transmitindo serenidade e profissionalismo. O objetivo é uma sonoridade polida, ideal para atendimento corporativo ou narrações sofisticadas, mantendo sempre a suavidade e a doçura na entonação.

## Decisão operacional

- `Groq` continua apenas no STT
- `Gemini TTS / Sulafat` é a voz pronta recomendada para uso imediato
- `MiniMax Audio` continua como lane premium de clonagem autorizada no futuro
- o runtime local atual de TTS continua como fallback seguro

## Voz pronta imediata

Escolha atual:

- provider: `gemini`
- model: `gemini-2.5-flash-preview-tts`
- voice: `Sulafat`
- formato: `wav`

Motivo:

- voz pronta feminina
- tom `Warm` no catálogo oficial
- boa aderência ao alvo formal, doce e corporativo

## Regras para a referência de clonagem

- usar apenas audio autorizado/licenciado
- preferir uma amostra entre `10s` e `5min`
- a fala da referencia deve ser:
  - PT-BR
  - formal
  - sem gírias
  - tom doce e educado
- o link de terceiro pode servir para transcript e estudo de prosodia, nunca como base de clonagem sem autorização

## Convenção de naming

- `voice_id`: `aurelia-ptbr-formal-doce-v1`
- upgrades futuros:
  - `aurelia-ptbr-formal-doce-v2`
  - `aurelia-ptbr-formal-doce-premium-v1`

## Smoke mínimo

- listar vozes da conta MiniMax
- sintetizar frase curta em PT-BR
- validar se a fala permanece clara e sem portunhol
- comparar com fallback local do Telegram
