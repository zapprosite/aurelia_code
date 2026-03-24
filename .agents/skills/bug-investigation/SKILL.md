---
type: skill
name: Bug Investigation
description: Investigação sistemática de bugs e análise de causa raiz no Homelab Ubuntu.
skillSlug: bug-investigation
phases: [E, V]
generated: 2026-03-18
updated: 2026-03-24
status: active
scaffoldVersion: "2.0.0"
---

# 🐛 Bug Investigation: Sovereign Debugging 2026

Habilita o Antigravity a diagnosticar falhas no sistema utilizando uma abordagem científica e ferramentas de observabilidade de baixo nível.

## 🏛️ Metodologia de Investigação

### 1. Reprodução
- Consiga um passo-a-passo determinístico para disparar o bug.
- Utilize logs de sistema (`journalctl -u aurelia`) e logs de aplicação.

### 2. Isolamento (Tier Analysis)
- O bug é no roteamento (LLM)? Verifique o Tier utilizado.
- O bug é de infra (Host)? Verifique VRAM, CPU e Portas.
- O bug é de rede? Verifique o status do OpenRouter e da VPN.

### 3. Ferramentas de Diagnóstico
- `strace / gdb`: Para crashes de binários Go.
- `nvidia-smi`: Para pânicos de GPU.
- `docker inspect`: Para falhas de volume ou rede em containers.

## 📍 Quando usar
- Quando um teste falha de forma inconsistente (Flaky Tests).
- Quando o daemon `aurelia` crasha ou reinicia inesperadamente.
- Quando a resposta do bot é alucinatória ou foge dos guardrails.

## 🛡️ Guardrails
- **Não altere estado durante a investigação**: Use o modo de leitura apenas.
- **Documente a Causa Raiz**: Não apenas corrija o sintoma, explique o "Porquê".
