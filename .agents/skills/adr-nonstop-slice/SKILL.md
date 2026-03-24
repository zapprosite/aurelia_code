---
type: skill
name: ADR Nonstop Slice
description: Fluxo de desenvolvimento ininterrupto baseado em ADRs dinâmicos e sessões de trabalho contínuas.
skillSlug: adr-nonstop-slice
phases: [P, E, V, C]
generated: 2026-03-20
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 🍕 ADR Nonstop Slice: Sovereign Workflow 2026

Esta skill otimiza o desenvolvimento de grandes funcionalidades através do fatiamento em "slices" independentes, mantendo o contexto vivo através de ADRs (Architectural Decision Records) incrementais.

## 🚀 Protocolo "Nonstop"
1. **Arquitetura Dinâmica**: O ADR não é estático. Ele evolui conforme descobrimos novas restrições técnicas (MiniMax Feedback).
2. **Ciclo de Trabalho**: `Plano -> ADR -> Implementação -> Teste -> Sincronismo -> Próxima Slice`.
3. **Autoridade do ADR**: O arquivo `docs/adr/README.md` é a bússola oficial. Se não está lá, não é oficial.

## 📍 Quando usar
- No desenvolvimento de grandes módulos (ex: Sistema de Visão, OCR, Dashboards novos).
- Quando a tarefa exigir mais de 4h de foco contínuo.
- Para manter o rastreio de decisões complexas entre sessões.
