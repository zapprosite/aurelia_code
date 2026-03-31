# PLANO: Jarvis Tutor - Audio Streaming Real

## AUDITORIA COMPLETA

| Componente | Status | Detalhes |
|------------|--------|----------|
| LiteLLM | ✅ | localhost:4000 - qwen3.5 funcionando |
| Ollama | ✅ | localhost:11434 - qwen3.5 disponível |
| Kokoro TTS | ✅ | localhost:8880 - pt-br female |
| Groq STT | ✅ | whisper-large-v3 via API |
| PulseAudio | ⚠️ | Testar parec/paplay |

## ARQUITETURA STREAMING

```
MICROFONE → parec → WAV → Groq STT → texto → LiteLLM → resposta → Kokoro TTS → WAV → AUTOFALANTE
      ↓
   JarViz (loop)
```

## IMPLEMENTAÇÃO

### FASE 1: Audio Capture (parec)
```go
// Captura contínua de áudio do microfone
parec --rate=16000 --channels=1 | sox - -p audio.wav trim 0 5
```

### FASE 2: STT (Groq API)
```go
// Transcrição via Groq API
POST https://api.groq.com/openai/v1/audio/transcriptions
```

### FASE 3: LLM (LiteLLM)
```go
// qwen3.5 via LiteLLM
POST localhost:4000/v1/chat/completions
```

### FASE 4: TTS (Kokoro)
```go
// Síntese via Kokoro
POST localhost:8880/v1/audio/speech
```

### FASE 5: Playback (paplay)
```go
// Reproduz WAV no autofalante
paplay audio.wav
```

## CÓDIGO LIMPO

```go
// JarvisTutor é o orchestrator principal
type JarvisTutor struct {
    mic      *Microfone
    stt      *GroqSTT
    llm      *LiteLLM
    tts      *KokoroTTS
    speaker  *Speaker
}

func (j *JarvisTutor) Loop() {
    for {
        audio := j.mic.Capture()
        texto := j.stt.Transcribe(audio)
        resposta := j.llm.Chat(texto)
        voz := j.tts.Speak(resposta)
        j.speaker.Play(voz)
    }
}
```

## PRÓXIMOS PASSOS

1. Implementar AudioPipeline completo
2. Testar fluxo mic→STT→LLM→TTS→speaker
3. Wire no systemd service
4. Teste E2E
