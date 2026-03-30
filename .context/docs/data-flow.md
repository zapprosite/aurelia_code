---
type: doc
name: data-flow
description: How data moves through the system and external integrations
category: data-flow
generated: 2026-03-30
status: filled
scaffoldVersion: "2.0.0"
---

# Fluxo de Dados — Aurélia Sovereign 2026.2

## Pipeline de Entrada (Texto)

```
Telegram message (chat_id, text, user_id)
    ↓
handleText() / handleVoice()
    ↓
processInput(requiresAudio: bool)
    ↓
Porteiro Sentinel (input guardrail)
    ↓ (owner bypass se user_id in allowedUserIDs)
Memory compress (context window)
    ↓
persistIncomingContext() → Redis + archive
    ↓
dashboard.Publish("user_message")
    ↓
prepareExecution()
    ├─ router.Route() — skill routing
    ├─ buildAgentHistory() — Redis message history
    └─ resolveExecutionPrompt() — system prompt + voice suffix
    ↓
executeConversation() / executeExternalConversation()
    ↓
agent.Loop.Run() — tool calling loop
    ├─ Tools: run_command, read_file, web_search, scheduling, memory, vision
    └─ LLM: qwen3.5:9b via LiteLLM cascade
    ↓
deliverWithParallelTTS() — texto imediato + áudio em goroutine
    ↓
Telegram send (markdown) + voice note (opus/ogg)
```

## Pipeline de Entrada (Áudio)

```
Telegram voice message
    ↓
handleVoice() → bot.Download() → /tmp/<fileID>
    ↓
bc.transcribeAudioFile()
    ├─ STT: faster-whisper (:8020) → Groq (:fallback)
    └─ transcript text
    ↓
persistAudioTranscript() → Redis
    ↓
processInput(transcript, requiresAudio: true)
    ↓
voiceSystemPromptSuffix injetado → resposta conversacional
    ↓
deliverWithParallelTTS() — texto + voice note
```

## Cascade LiteLLM

```
Request "aurelia-smart"
    ↓
Priority 1: ollama/qwen3.5:9b (:11434) — 0ms RTT, free
    ↓ (se falhar)
Priority 2: openrouter/qwen/qwen-3.6-plus-preview:free — free
    ↓ (se falhar)
Priority 3: openrouter/minimax/minimax-2.5:free — free
    ↓ (se falhar)
Priority 4: groq/llama-3.3-70b-versatile — free (14.400 req/dia)
    ↓ (se falhar)
Priority 10: openrouter/minimax/minimax-m2.7 — paid
    ↓ (último recurso)
Priority 11: openrouter/moonshotai/kimi-k2.5 — paid
```

## Cache Redis

| Key Pattern | TTL | Uso |
|---|---|---|
| `aurelia:history:<user_id>` | 7 dias | Mensagens por conversa |
| `aurelia:pending:<user_id>` | 30min | Bootstrap state machine |
| `aurelia:context:<user_id>` | 30 dias | Contexto compactado |
| `porteiro:cache:<sha256(text)>` | 30 dias | Resultado do guardrail |

## Dependências Externas

| Serviço | Tipo | Auth | Rate Limit |
|---|---|---|---|
| Ollama (:11434) | LLM local | none | NUM_PARALLEL=2 |
| LiteLLM (:4000) | Router | LITELLM_MASTER_KEY | cooldown 30s |
| faster-whisper (:8020) | STT local | none | 1 concurrent |
| Groq API | STT cloud | GROQ_API_KEY | 60s timeout |
| Tavily API | Search | TAVILY_API_KEY | fallback only |
| DuckDuckGo | Search | none | fallback, 10s timeout |
| Kokoro (:8012) | TTS local | none | 2 concurrent, 30s timeout |
| Redis (:6379) | Cache | none | — |
| Qdrant (:6333) | Vetor DB | none | — |
| OpenRouter | LLM cloud | OPENROUTER_API_KEY | tier-paid |
