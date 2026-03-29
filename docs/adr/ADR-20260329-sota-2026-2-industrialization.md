# ADR-20260329-sota-2026-2-industrialization

## Status
Aceito

## Contexto
O ecossistema Aurélia precisava de industrialização para atingir os padrões SOTA 2026.2. Isso envolveu a estabilização do pipeline de streaming, implementação de capacidades de "Barge-in" (interrupção por voz) e estabelecimento de uma stack de modelos soberana.

## Decisões
1.  **Modelo L0**: O **Qwen 3.5 9B VL** foi estabelecido como o cérebro primário para visão e lógica central, purgando o Gemma3.
2.  **Pipeline de Streaming**: Implementada uma arquitetura baseada em Atores (SAP) em `internal/streaming` para orquestração multi-agente robusta.
3.  **Barge-in (Interrupção)**: Integrado monitoramento VAD e cancelamento imediato de buffer para permitir interrupção em tempo real (< 200ms).
4.  **Resiliência (Self-Healing)**: Adicionado um supervisor para reinicialização automática de serviços e atores em caso de falha.
5.  **Build Industrial**: Compilação estática (`CGO_ENABLED=0`) com injeção de metadados de versão para deploy soberano.

## Consequências
- Latência de interrupção sub-200ms.
- Binário único e portátil (`bin/aurelia`).
- Sucesso de 100% na suíte de testes (Agent Loop & E2E).
- Redução de custos operacionais via otimização de GPU (Whisper Distil).
