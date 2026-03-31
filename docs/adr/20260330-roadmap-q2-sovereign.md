# ADR-20260330-roadmap-q2-sovereign

## Contexto
Após a conclusão da industrialização SOTA 2026.2, o ecossistema Aurélia possui infraestrutura estável. O objetivo do Q2/2026 é a **Autonomia Funcional Profunda** e a **Interface Humana (Jarvis-like)**.

## Roadmap Priorizado [ ] (Ordem de Execução)

### 1. 🎙️ Her-Mode: Always-On Voice Pipeline
- **Status**: [ ] Planejado
- **Decisão**: Integrar o loop VAD -> Whisper -> LLM -> Kokoro diretamente no lifecycle do `BotController`.
- **Meta**: Zero latência perceptível e interrupção por voz (Barge-in).

### 2. 🛰️ Sovereign Dashboard & Self-Healing
- **Status**: [ ] Planejado
- **Decisão**: Implementar monitoramento de telemetria em tempo real (GPU/Memória/Latência) com correção automática de processos zumbis.
- **Meta**: 99.9% de uptime do Homelab sem intervenção manual do Will.

### 3. 🕹️ Native God Mode: OS Control
- **Status**: [ ] Planejado
- **Decisão**: Ativar o `os_controller` via MCP para controle total do Desktop (X11/Display=:1).
- **Meta**: Jarvis executando tarefas complexas no Chrome e Terminal de forma autônoma.

### 4. 🧠 Ephemeral Session Summarizer
- **Status**: [ ] Planejado
- **Decisão**: Destilação automática de sessões de chat em Knowledge Items (KIs) no Qdrant.
- **Meta**: Persistência de contexto a longo prazo e redução de janelas de contexto.

### 5. 👁️ Real-Time Vision Pipeline
- **Status**: [ ] Planejado
- **Decisão**: Integração do Qwen 2-VL para análise contínua de frames do desktop.
- **Meta**: "Entendimento" visual do que o usuário está fazendo no IDE.

## Consequências
- Foco total em utilidade prática e autonomia.
- Aumento da carga computacional (GPU/NPU) exigindo orquestração fina de recursos.
