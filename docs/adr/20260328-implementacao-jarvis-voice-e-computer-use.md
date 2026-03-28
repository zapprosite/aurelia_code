# ADR 20260328: Implementação do Protocolo Jarvis (Voice Loop & Computer Use)

## Status
⏳ Proposta (Aguardando Aprovação de Slice)

## Contexto
Necessidade de dotar o ecossistema `aurelia` de capacidades sensoriais (escuta/fala) e motoras (navegação em interface gráfica) de forma autônoma e resiliente. O sistema deve operar sob a governança SOTA 2026.1, mantendo o `gemma3:27b` como o Juiz Soberano no LiteLLM.

## Decisões Arquiteturais

### 1. Voice Loop Isolado (CPU/RAM Focus)
Para preservar os 24GB de VRAM da RTX 4090 exclusivamente para inferência de Modelos de Linguagem (LLM) e busca semântica (RAG), os serviços de voz serão desacoplados:
- **STT (Whisper)**: Execução local via processos isolados otimizados para CPU.
- **TTS (Kokoro)**: Síntese de voz nativa rodando em RAM/CPU.
- **Integração**: O core em Go (`aurelia_code`) comunicará com estes processos via buffers de stream assíncronos, evitando bloqueios de I/O.

### 2. Computer Use via Protocolo MCP (Stagehand)
A interação com o mundo exterior (navegação web/GUI) não será hardcoded no core:
- **Ambiente**: Container isolado (Steel) para garantir segurança e persistência de sessão.
- **Controlador**: Servidor MCP escrito em Node/TypeScript utilizando o SDK do **Stagehand**.
- **Interface**: A `aurelia_code` consumirá as capacidades do Stagehand estritamente através do **Model Context Protocol (MCP)**, tratando o navegador como uma ferramenta externa dinâmica.

## Consequências
- **Positivas**: Isolamento total de falhas (um crash no navegador não derruba o agente), performance térmica otimizada (GPU livre), escalabilidade horizontal do braço motor.
- **Desafios**: Latência marginal na comunicação via buffers de áudio; necessidade de gerenciar o lifecycle do servidor MCP.

## Referências
- [internal/gateway/provider.go](file:///home/will/aurelia/internal/gateway/provider.go) (LiteLLM Gateway)
- [configs/litellm/config.yaml](file:///home/will/aurelia/configs/litellm/config.yaml) (Smart Router)
