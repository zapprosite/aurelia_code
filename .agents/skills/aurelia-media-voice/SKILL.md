---
name: aurelia-media-voice
description: Governa a identidade auditiva da Aurélia e o processamento de mídias ricas (Transcript, TTS, Voice Cloning).
---

# 🎙️ Aurelia Media Voice: Sovereign 2026

Esta skill habilita o Antigravity a gerenciar a voz oficial da Aurélia e o processamento de mídias (áudio/vídeo), garantindo uma identidade auditiva premium e soberana.

## 🏛️ Estratégia de Voz (2026)

### 1. Motor Principal (Sovereign Local)
- **Engine**: `Kokoro-TTS` (ou Kodoro) rodando localmente via GPU/CPU.
- **Voz**: Feminina brasileira, formal, doce e profissional.
- **Uso**: Respostas rápidas no Telegram e notificações do sistema.

### 2. Motor Premium (Cloud Fallback)
- **Engine**: `MiniMax Audio (m2.7)`.
- **Voz**: Clonagem oficial e alta fidelidade (Hi-Fi).
- **Uso**: Conteúdo de marca, mensagens de boas-vindas e síntese de textos longos.

### 3. Inteligência Auditiva (STT)
- **Primário**: `Groq` (Whisper-v3) para transcrição instantânea com baixíssima latência.
- **Auxiliar**: `MiniMax Audio` para transcrições técnicas complexas.

## 🛠️ Operações de Mídia (Industrial)

### 1. Transcrição de Vídeos e Links
- **Ferramenta**: `scripts/media-transcript.sh`.
- **Processo**: Download via `yt-dlp` -> Extração de áudio via `ffmpeg` -> STT via Groq/MiniMax.
- **Guardrail**: Evite o processamento de vídeos > 500MB no host principal se a carga de inferência estiver alta.

### 2. Gestão de Identidade (Voice Profile)
- Consulte `docs/aurelia_voice_profile_20260319.md` para diretrizes de tom e estilo.
- **Proibição**: Nunca utilize gírias ou tons informais na voz oficial, a menos que solicitado para uma persona específica de teste.

### 3. Hardware Awareness (7900x / RTX 4090)
- O processamento de áudio local deve ser pinado nos núcleos de performance do 7900x para evitar jitter.
- Monitore a VRAM se o Kokoro estiver rodando via CUDA.

## 📍 Quando usar
- Para transcrever vídeos do YouTube ou áudios do Telegram.
- Para gerar respostas em voz (TTS) para o usuário.
- Para realizar clonagem ou benchmark de novas vozes.

## 🚫 Anti-Padrões
- Clonar vozes de terceiros sem autorização.
- Usar motores de nuvem caros para mensagens triviais de sistema.
- Ignorar o tratamento de ruído de fundo em transcrições.
