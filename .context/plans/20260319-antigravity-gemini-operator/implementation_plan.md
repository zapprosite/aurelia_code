---
title: Antigravity Gemini Operator Integration
status: in_progress
owner: antigravity
created: 2026-03-19
---

# Implementation Plan

## Goal

Configurar a Aurelia para operar com um modo tutor/orquestrador que saiba quando e como delegar microtarefas ao chat do Antigravity, usando Gemini Flash e o contexto do workspace sem quebrar a governanca.

## Architecture

1. Runtime memory
   - registrar um `PROJECT_PLAYBOOK.md` com a regra operacional
2. Project skill
   - criar uma skill em `.aurelia/skills/antigravity-gemini-operator/`
3. Human/operator blueprint
   - criar um blueprint curto e acionavel para o workspace
4. Execution discipline
   - definir matrix de roteamento, prompts, evidencias e limites

## Verification

1. skill criada no path que o runtime carrega
2. frontmatter valido no `SKILL.md`
3. playbook do projeto presente em `docs/PROJECT_PLAYBOOK.md`
4. task board atualizado

## Notes

- a delegacao para o chat do Antigravity sera restrita a pesquisa leve, pequenas configuracoes e preparacao de diffs
- rede, deploy, secrets e acoes irreversiveis continuam fora do escopo do chat leve
