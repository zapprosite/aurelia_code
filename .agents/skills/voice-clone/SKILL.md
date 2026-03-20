---
name: voice-clone
description: Skill para clonagem e injeção de vozes no motor Kokoro (Kodoro) local.
---

# Voice Clone (Kokoro Local)

## Objetivo
Clonar e injetar identidades vocais no ecossistema da Aurélia, garantindo que o motor de síntese local (RTX 4090) utilize vozes personalizadas com sotaque correto em Português-Brasil.

## Quando usar
- Para atualizar a voz oficial da Aurélia (`Aurelia.wav`).
- Para criar novas vozes de agentes secundários ou bots de terceiros.
- Quando o motor de voz regredir para fonemas estrangeiros (Espanhol).

## Processo de Clonagem

### 1. Amostragem (Input)
- Obtenha um áudio de 10 a 30 segundos (claro, sem ruído de fundo).
- Formato preferencial: `.wav` ou `.mp3`.

### 2. Sanitização e Conversão
A Aurélia exige amostras em 24kHz Mono para performance ideal no Kokoro. Use `ffmpeg`:
```bash
ffmpeg -i input.mp3 -ar 24000 -ac 1 Aurelia.wav
```

### 3. Injeção no Container
O motor de voz (Chatterbox/Kodoro) lê as vozes de um volume Docker. Injete o arquivo:
```bash
docker cp Aurelia.wav chatterbox-tts:/app/voices/
```

### 4. Configuração do Código
Ao solicitar o áudio via API (OpenAI Compatible), use os parâmetros:
- `voice`: "Aurelia.wav" (nome do arquivo injetado).
- `language`: "pt" (OBRIGATÓRIO para PT-BR).

## Governança e Ética
- **Autorização**: Apenas use vozes sob licença ou autorização explícita (como a voz da Aurélia).
- **Local-First**: Priorize sempre a injeção local em vez de APIs de nuvem para garantir a soberania do áudio.

## Troubleshooting
- **Sotaque Espanhol**: `internal/config` já impõe `tts_language: "pt"` e `pkg/tts/openai_compatible.go` só remove essa chave se ficar vazia, portanto o valor oficial nunca pode ser alterado durante o runtime. Use `scripts/verify-tts-language.sh` para validar o JSON de configuração local e (opcionalmente) pingar `http://127.0.0.1:8484/v1/voice/synthesize`.
- **Voz não encontrada**: Verifique se o arquivo está dentro de `/app/voices/` no container `chatterbox-tts`.
