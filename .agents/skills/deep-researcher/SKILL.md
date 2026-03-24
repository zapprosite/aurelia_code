---
type: skill
name: deep-researcher
description: Pesquisa profunda utilizando Gemini Web e ferramentas externas para descoberta de novos conceitos.
skillSlug: deep-researcher
phases: [P]
generated: 2026-03-19
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 🔬 Deep Researcher: Gemini Web — Base de Pesquisa da Aurélia

Skill que navega no **https://gemini.google.com/app** via Playwright para executar pesquisas em 4 modos: flash, pro, deep e reasoning.

## 🏗️ Modos de pesquisa

| Modo | Modelo | Tempo | Uso |
|------|--------|-------|-----|
| `flash` | Gemini 2.0 Flash | ~15s | Fatos rápidos, perguntas simples |
| `pro` | Gemini 2.5 Pro | ~30s | Análise profunda, documentos longos |
| `deep` | Gemini 2.5 Pro + Deep Research | 2–5min | Relatório com fontes web |
| `reasoning` | Gemini 2.0 Flash Thinking | ~45s | Raciocínio passo-a-passo |

## 🔧 Protocolo de execução

Seguir workflow: [`.agents/workflows/pesquisa-profunda.md`](../../workflows/pesquisa-profunda.md)

1. `browser_navigate` → `https://gemini.google.com/app`
2. Selecionar modelo conforme modo
3. Ativar "Deep Research" se modo=`deep`
4. Digitar query + Enter
5. Aguardar resposta (timeout por modo)
6. `browser_snapshot` → extrair texto
7. Salvar em `~/Desktop/pesquisas-gemini/<slug>-<data>.md`
8. Retornar resumo ao Telegram

## 📍 Quando usar
- Pesquisa técnica com fontes verificadas (modo `deep`)
- Análise rápida de conceitos (modo `flash`)
- Relatórios longos e documentos (modo `pro`)
- Raciocínio encadeado e lógica complexa (modo `reasoning`)
- Feed direto para `architect-planner` → criação de ADRs
