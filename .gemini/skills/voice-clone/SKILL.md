---
name: voice-clone
description: Skill para clonagem e injeção de vozes no motor Kokoro (Kodoro) local ou MiniMax Audio.
---

# 🎙️ Voice-Clone: Sovereign Voice Identity 2026

Habilita a criação de clones de voz digitais para a Aurélia, utilizando técnicas de Zero-Shot TTS ou Fine-tuning em hardware local.

## 🏛️ Protocolo de Clonagem (Industrial)

### 1. Motor de Clonagem
- **Local (GPU)**: `Kokoro-TTS` com técnica de Voice Styles.
- **Premium (Cloud)**: `MiniMax Voice Cloning` para fidelidade absoluta.

### 2. Identidade Oficial
- A voz padrão da Aurélia deve ser mantida consistente (Feminina, Brasileira, Doce).
- Novos clones para testes devem ser documentados em `docs/voice_clones/`.

### 3. Requisitos Técnicos
- **Amostra**: Exige áudio limpo (WAV/MP3) com duração mínima de 10 segundos sem ruído.
- **Hardware**: Processamento de inferência via CUDA na RTX 4090 para baixíssima latência (< 500ms).

## 🛡️ Guardrails Éticos
- **Permissão**: Nunca clone vozes de pessoas reais sem autorização explícita registrada.
- **Uso Malicioso**: É proibido o uso da voz clonada para impersonificação ou testes de engenharia social.

## 📍 Quando usar
- Ao atualizar a identidade sonora da Aurélia.
- Para criar vozes de assistentes especializados (ex: "Aurelia Engineer", "Aurelia Support").
- Para testar novos modelos de TTS no mercado.