---
description: Matriz de extensões opcionais para Chrome/Antigravity com foco em rollback simples.
status: active
---

# Extension Matrix 2026-03-19

## Regra principal

Nenhuma extensão faz parte do core da Aurelia.

O core continua sendo:

- `Go`
- `CLI`
- `agent-browser`
- `Playwright`
- `Chrome DevTools`
- `systemd`

## Core

- nenhuma extensão permitida

## Nice-to-have

| Área | Classe | Valor | Risco | Status |
| --- | --- | --- | --- | --- |
| Chrome | bloqueio leve de ruído visual/anúncios | melhora foco do operador humano | pode alterar DOM de páginas públicas | opcional |
| Chrome | visualizador JSON | leitura mais rápida de payloads | baixo | opcional |
| Chrome | salvamento de página única | arquivamento manual | baixo | opcional |
| Antigravity | nenhuma extensão externa por padrão | evitar drift de UX/segurança | n/a | recomendado |

## Risky

| Classe | Motivo do risco |
| --- | --- |
| manipulador de headers/cookies | pode mascarar bugs e alterar sessões |
| automação por userscript | cria runtime paralelo difícil de governar |
| extensões de IA com leitura global de página | mistura contexto e secrets fora do contrato do repo |

## Política de instalação

- sempre em perfil isolado
- nunca no perfil principal de operação
- nunca como pré-requisito do runtime
- sempre com rollback simples

## Rollback

1. remover a extensão do perfil isolado
2. reiniciar o Chrome desse perfil
3. validar que o fluxo base com DevTools/Playwright continua verde
