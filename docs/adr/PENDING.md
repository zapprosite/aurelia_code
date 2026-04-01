# PENDING — Slices Aguardando Implementação

Última atualização: 01/04/2026

## P1 — Crítico 🔴

| Slice | Descrição | Pré-requisito |
|---|---|---|
| S-43 | PostgreSQL + pgvector substituir Supabase local | — |
| S-44 | Grafana + Prometheus no docker-compose.yml | — |
| S-47 | Computer Use E2E (BUA-style) | — |
| S-48 | OS Native God Mode | S-47 |
| S-49 | Jarvis Voice + Computer | S-48 |

## P2 — Alto 🟡

| Slice | Descrição | Pré-requisito |
|---|---|---|
| S-50 | Ubuntu Desktop voice-gateway (Voice Native) | S-43 |

### Gap Analysis: S-50 (Voice Gateway)

Análise técnica de `cmd/voice-gateway/`:

- **Estado Atual**:
    - `main.py`: Em **MOCK MODE**. Simula stream de áudio injetando ruído via socket UNIX (`/tmp/aurelia-voice.sock`).
    - `wakeword.py`: Implementação funcional de captura one-shot usando `openwakeword` + `arecord`. Wake word: `"jarvis"`.
- **Gaps reais**:
    - O `main.py` precisa de integração com `PyAudio` ou similar para captura real de stream contínuo no Ubuntu Desktop.
    - O pipeline Go precisa consumir o socket de forma estável para transcrição em tempo real (Cascata Groq/Whisper).
    - Falta orquestração sistemática para rodar como serviço de background estável.
