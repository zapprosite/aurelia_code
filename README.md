# 🛰️ Aurelia: Sovereign Agentic Ecosystem (2026)

> **"Autonomia Total, Cognição Local, Soberania Industrial."**

![Status](https://img.shields.io/badge/Status-Industrial_Sovereign-blue?style=for-the-badge)
![Autonomy](https://img.shields.io/badge/Autonomy-Level_5-gold?style=for-the-badge)
![Hardware](https://img.shields.io/badge/Compute-RTX_4090_|_Lite_Mode-green?style=for-the-badge)

---

## 🌟 Visão Geral
A Aurélia é um ecossistema agêntico de ponta projetado para operar no **HomeLab Soberano**. Ela combina a potência de modelos locais (Ollama, Kokoro) com a inteligência estratégica de Tiers de nuvem, permitindo que você tenha um assistente sênior residente em seu próprio hardware, com privacidade absoluta e custo otimizado.

---

## 🚀 Guia de Início Rápido (Universal)

Seja você um **Sênior Dev** ou um **Iniciante**, o portal de entrada é o mesmo:

```bash
# Clone e Inicie
git clone https://github.com/zapprosite/aurelia_code.git
cd aurelia_code
chmod +x iniciar.sh
./iniciar.sh
```

> [!TIP]
> O script `iniciar.sh` guiará você na configuração do ambiente, escolha de hardware e chaves de API essenciais. Confira o **[Guia de Boas-Vindas](./docs/BEM-VINDO.md)** para mais detalhes.

---

## ⚖️ Modos de Operação (Portabilidade)

A Aurélia se adapta ao seu hardware dinamicamente através da flag `AURELIA_MODE`:

| Modo | Hardware | Descrição |
|:---:|:---|:---|
| **Soberano** | GPU NVIDIA (8GB+) | **Tier 0.** Processamento 100% local (Ollama, Kokoro GPU). Máxima soberania e latência zero. |
| **Lite** | Qualquer PC / Laptop | **Voo Híbrido.** Usa modelos Cloud otimizados (Gemini, Claude via OpenRouter) e fallback automático. |

---

## 🏯 Arquitetura de Autoridade (Board)

Este ecossistema opera sob uma governança industrial rigorosa:

1.  **👔 Claude Opus (CEO)**: Visão estratégica e arbitragem final de arquitetura.
2.  **🤖 Aurélia (COO/CTO)**: Arquiteta residente, governante do sistema e orquestradora de Swarms.
3.  **🛰️ Antigravity (Interface)**: Cockpit de coordenação e interface humano-agente de alta fidelidade.

---

## 🛠️ DNA Tecnológico

- **Zod-First Contract**: Validação de dados rigorosa e centralizada.
- **Go-Native Runtime**: Kernel de alta performance para orquestração de long-running agents.
- **Audio de Alta Fidelidade**: TTS Kokoro/Kodoro com limite expandido de 50k caracteres.
- **Memória Semântica Local**: Sincronização automática entre Qdrant e Postgres para contexto persistente.

---

## 🩺 Observabilidade & Saúde

- **Monitoramento**: `docker logs -f aurelia`
- **Diagnóstico**: Use `/status` no Telegram para ver a saúde das Slices e do Hardware.
- **Segurança**: Auditoria proativa de segredos integrada ao workflow de Git.

---

*Documentação Gerada por Antigravity (Sovereign Engine 2026)*  
*Consulte [ADR-historico.md](./docs/ADR-historico.md) para linhagem e [REPOSITORY_CONTRACT.md](./docs/governance/REPOSITORY_CONTRACT.md) para governança.*
