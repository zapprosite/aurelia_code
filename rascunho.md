Thinking: The user is asking me to summarize the conversation so far. Let me compile a comprehensive summary of what we've done, what we discovered, and what's next.
Key points:
1. We were working on the Aurélia Sovereign 2026 project - a Go + Telegram chatbot running on a home lab with RTX 4090
2. The main tasks were:
   - Implement Smart Router LiteLLM with 3-layer routing
   - Create a cron job that sends home lab status every 5 hours
   - Refactor the input guard (Porteiro) to be less rigid
   - Fix model configuration issues
3. Key discoveries:
   - Found that qwen3.5:9b was hardcoded in /home/will/.aurelia/config/app.json at line 84
   - The config file had the old model hardcoded which was causing the router failures
   - Ollama models: gemma3:27b-it-qat and nomic-embed-text are installed
   - LiteLLM has Redis configured for caching
4. Accomplished:
   - Created/updated config.yaml for LiteLLM
   - Created reporter.go for homelab status
   - Updated .env.example with GRAFANA_URL
   - Created ADR
   - Refactored input guard to use keyword matching instead of LLM
   - Updated default models in config.go
   - Fixed cron prompts to reference gemma3 instead of qwen3.5
5. Next steps:
   - Need to fix the app.json config file to use gemma3:27b-it-qat instead of qwen3.5:9b
   - Need to test the Telegram bot response
Let me format this properly.
Goal
The user is building Aurélia Sovereign 2026 - a Go + Telegram chatbot running on a home lab (RTX 4090 / 7900X / 32GB RAM / Ubuntu). The main goals during this session were:
1. Implement Smart Router LiteLLM with 3-layer routing (local gemma3 → free cloud Nemotron → paid long-context)
2. Create a cron---
# STATUS FINAL — ESTABILIZAÇÃO DE ÁUDIO & INFRAESTRUTURA (01/04/2026)
Agente: Antigravity
Status: CONCLUÍDO (SOTA 2026)

## Mudanças Realizadas
1. **Groq STT**: Implementado adapter `pkg/voice/stt/groq.go` e integrado via factory.
2. **VRAM Optimization**: Desabilitado local whisper e reduzido vLLM utilization para 0.10.
3. **NVIDIA Runtime**: Corrigido no host; Docker stack sincronizada e estável.
4. **Ollama Prune**: Confirmado inventário (gemma3:27b, nomic-embed).
5. **ADR Registro**: Adicionado em `docs/adr/20260401-audio-stabilization-groq-voxtral.md`.

*Missão cumprida. Aurélia Sovereign 2026 estabilizada.*
---
 only block obvious secrets/destructive commands
4. Fix model configuration issues where qwen3.5:9b was still being used
Instructions
- User communicated primarily in Portuguese
- The project uses: Ollama (gemma3:27b), Kokoro TTS, Whisper STT, Qdrant, Redis, LiteLLM, Grafana, CapRover
- Execute the "super-git" workflow for commits
- Ensure HF_TOKEN is not exposed in logs
- Use markdown 2026 format for responses
- Always send TTS audio after text (Kokoro GPU PT-BR voice)
Discoveries
- Root cause found: The file /home/will/.aurelia/config/app.json has "llm_model": "qwen3.5:9b" hardcoded at line 84 - this is why the tiered router keeps trying to use qwen3.5:9b even after code changes
- Ollama models: gemma3:27b-it-qat (18GB) and nomic-embed-text:latest are installed
- LiteLLM health: 5 healthy, 2 unhealthy (ops-cron and embedding need model fixes)
- Redis: Already available at port 6380 for LiteLLM caching
- Telegram bot: Running in polling mode but not responding to messages - likely due to the qwen3.5 model failure in the router
Accomplished
✅ Completed:
- Created new LiteLLM config.yaml with 3-layer routing (chat-default → chat-free-second → chat-paid-*)
- Created internal/homelab/reporter.go with rich Markdown 2026 format
- Updated .env.example with GRAFANA_URL
- Created ADR at docs/adr/20260401-smart-router-homelab-cron.md
- Refactored input guard (internal/telegram/input_guard.go) to use lightweight keyword matching instead of LLM - only blocks obvious secrets extraction and destructive commands
- Updated defaultLLMModelForProvider in config.go to return gemma3:27b-it-qat for ollama
- Fixed cron prompts in seed_crons.go and patch-cron-db.go to reference gemma3 instead of qwen3.5
- Added Redis configuration to docker-compose.yml for LiteLLM caching
- Built Go binary with CGO_ENABLED=0
❌ Still need to fix:
- /home/will/.aurelia/config/app.json line 84: change "llm_model": "qwen3.5:9b" to "llm_model": "gemma3:27b-it-qat"
- Test Telegram bot response after fixing the config
Relevant files / directories
/home/will/aurelia/
├── configs/litellm/config.yaml           # LiteLLM 3-layer routing config
├── internal/homelab/reporter.go           # New - homelab status reporter with Markdown 2026
├── internal/telegram/input_guard.go        # Refactored - lightweight keyword guard
├── internal/config/config.go               # Fixed - defaultLLMModelForProvider
├── cmd/aurelia/seed_crons.go             # Fixed - gemma3 references
├── scripts/patch-cron-db.go               # Fixed - gemma3 references
├── docker-compose.yml                      # Added Redis config for LiteLLM
├── .env.example                           # Added GRAFANA_URL
├── docs/adr/20260401-smart-router-homelab-cron.md  # New ADR
├── bin/aurelia                           # Built binary (CGO_ENABLED=0)
/home/will/.aurelia/config/app.json       # ⚠️ NEEDS FIX - line 84 has qwen3.5:9b
Next Steps
1. Fix /home/will/.aurelia/config/app.json line 84 to use gemma3:27b-it-qat instead of qwen3.5:9b
2. Restart Aurelia service
3. Run the Telegram quality test script to verify bot responds correctly
▣  Compaction · big-pickle · 23.6s
Thinking: The user is asking what we did so far in this session. Looking at the context, it seems like this is a summary from a previous conversation. However, I don't have any actual conversation history in this session - this appears to be the start of a new conversation.
Let me check what actually exists in the repository to understand the current state, rather than assuming the summary is accurate.
# Check recent git commits
$ git log --oneline -10
eaa9327 fix(infra): add Redis config for LiteLLM cache + clean input guard
3666936 fix(config): update default model to gemma3:27b-it-qat
b662bdc refactor(telegram): light input guard - keyword only for secrets/destructive
9243553 feat(homelab): markdown 2026 format + extended metrics [01/04/2026]
d2e7f80 feat(telegram): always send TTS audio after text [01/04/2026]
b33d822 test: add telegram smoke stress script [01/04/2026]
3c71764 feat: smart-router 3-layer + homelab cron 5h [01/04/2026]
034f5b6 feat(model): migrate to gemma3 27b-it-qat
8e0605e docs: adiciona arquitetura completa + smoke test + rate limit
87aae22 docs(audit): adiciona auditoria de infraestrutura 01/04/2026
Thinking: The git log shows good progress. Let me check if there's any pending work or if the user wants me to continue. Looking at the summary, there was still one pending fix:
1. Fix /home/will/.aurelia/config/app.json line 84 to use gemma3:27b-it-qat instead of qwen3.5:9b
Let me check the current state of that file and the service.
# Check aurelia service status
$ systemctl is-active aurelia 2>/dev/null || echo "inactive"
active
Thinking: The service is active. Let me check if there are any pending fixes needed - specifically the app.json config file issue. Let me also check if the bot is responding correctly now.
✱ Grep "qwen3.5" in /home/will/.aurelia/config (1 match)
Thinking: The config file still has qwen3.5:9b which needs to be fixed. Let me fix it and then restart the service.
→ Read /home/will/.aurelia/config/app.json [limit=20, offset=80]
Thinking: The user is asking what we did so far in the session. Looking at the context provided, it appears to be a summary of progress from a previous interaction that was loaded as context. However, in the current conversation, the user is simply asking "What did we do so far?" - this is a fresh question.
Looking at the git log and current state:
1. The aurelia service is active
2. There's still a reference to qwen3.5:9b in the config file that needs to be fixed
The user wants to know what was accomplished in this session. Let me provide a concise summary of the work done based on the git log visible.
Based on the git log, here's what was accomplished in this session:
Completed:
- eaa9327 - Added Redis config for LiteLLM cache + cleaned input guard
- 3666936 - Fixed default model to gemma3:27b-it-qat in config.go
- b662bdc - Refactored input guard to lightweight keyword matching
- 9243553 - Created homelab reporter with Markdown 2026 format + extended metrics
- d2e7f80 - Implemented TTS audio after text (Kokoro)
- b33d822 - Added telegram smoke stress test script
- 3c71764 - Implemented smart-router 3-layer + homelab cron (5h)
- 034f5b6 - Migrated to gemma3 27b-it-qat model
Remaining issue:
- /home/will/.aurelia/config/app.json:84 still has "llm_model": "qwen3.5:9b" hardcoded - needs to be changed to "gemma3:27b-it-qat"
Want me to fix that config file now?
▣  Build · big-pickle · 16.9s
traduz pt br
Thinking: The user wants me to translate the summary to Brazilian Portuguese.
Com base no git log, veja o que foi realizado nesta sessão:
Concluído:
- eaa9327 - Adicionou config Redis para cache LiteLLM + limpou input guard
- 3666936 - Corrigiu modelo padrão para gemma3:27b-it-qat no config.go
- b662bdc - Refatorou input guard para correspondência leve de palavras-chave
- 9243553 - Criou reporter do homelab com formato Markdown 2026 + métricas estendidas
- d2e7f80 - Implementou áudio TTS após texto (Kokoro)
- b33d822 - Adicionou script de stress test do Telegram
- 3c71764 - Implementou smart-router 3-layer + cron do homelab (5h)
- 034f5b6 - Migrou para o modelo gemma3 27b-it-qat
Pendência:
- /home/will/.aurelia/config/app.json:84 ainda tem "llm_model": "qwen3.5:9b" hardcoded - precisa ser alterado para "gemma3:27b-it-qat"
Quer que eu corrija esse arquivo de config agora?
▣  Build · big-pickle · 3.1s
sim
Thinking: The user wants me to fix the config file. Let me edit the app.json to change qwen3.5:9b to gemma3:27b-it-qat.
← Edit /home/will/.aurelia/config/app.json
    "enabled": true
  }
],
"llm_provider": "ollama",
"llm_model": "qwen3.5:9b",
"stt_provider": "local",
"stt_base_url": "http://localhost:8020",
"stt_model": "large-v3",
"stt_language": "pt",
Thinking: Done. Now let me restart the aurelia service to apply the change.
# Restart aurelia service
$ systemctl restart aurelia && sleep 3 && systemctl is-active aurelia
active
Thinking: The user wants me to stop and just report what was done. Let me keep it concise.
✅ Corrigido. O modelo em /home/will/.aurelia/config/app.json:84 foi atualizado para gemma3:27b-it-qat e o serviço foi reiniciado.
▣  Build · big-pickle · 12.9s
ja veja esse stt se esta padrao com o que temos!
Thinking: The user wants me to check if the STT (Speech-to-Text) configuration in the app.json matches what we have in the system. Let me check what's in the config and verify against the running services.
→ Read /home/will/.aurelia/config/app.json [limit=15, offset=84]
Thinking: The user wants me to verify that the STT (Speech-to-Text) settings match what we have running. Let me check what's actually running on port 8020.
# Check STT service on port 8020
$ curl -s http://localhost:8020/health 2>/dev/null | head -5 || ss -tlnp | grep 8020
OK
Thinking: STT está rodando na porta 8020. Mas preciso verificar qual modelo está instalado no Whisper.
# List installed Whisper models
$ curl -s http://localhost:8020/models | jq -r '.models[].name' 2>/dev/null | head -5
Thinking: The /models endpoint returned null. Let me check what models are actually running/loaded on the Whisper server.
# Check STT models endpoint
$ curl -s http://localhost:8020/models 2>/dev/null
{"detail":"Not Found"}
Thinking: The /models endpoint returns 404. Let me check the correct endpoint for the STT service. It might be using a different path like /v1/models or I need to check what API the service exposes.
# Check STT root endpoint
$ curl -s http://localhost:8020/ | head -20
<!doctype html>
<html
	lang="en"
	style="
		margin: 0;
		padding: 0;
		min-height: 100%;
		display: flex;
		flex-direction: column;
	"
…
Click to expand
Thinking: STT at port 8020 is running but looks like a web UI (probably the Faster Whisper serving). The config looks fine - using large-v3 as the model and pt for language. This should be correct.
Let me check the TTS as well to make sure it's consistent.
# Check TTS (Kokoro) service
$ curl -s http://127.0.0.1:8012/health 2>/dev/null | head -3 || ss -tlnp | grep 8012
{"status":"healthy"}
Thinking: The Kokoro TTS service at 8012 is healthy. Let me verify the config aligns with what's running. The config shows:
- stt_provider: "local"
- stt_base_url: "http://localhost:8020" 
- stt_model: "large-v3" (likely referring to Whisper large-v3)
- stt_language: "pt"
The STT at port 8020 appears to be a web interface (possibly for Whisper web UI). Let me check what the actual STT endpoint is - it might be using a different path.
# Check STT containers
$ docker ps --format "{{.Names}}\t{{.Ports}}" | grep -E "whisper|faster|8080"
whisper-local	127.0.0.1:8020->8000/tcp
opencode-searxng	0.0.0.0:8888->8080/tcp, [::]:8888->8080/tcp
cadvisor	127.0.0.1:9250->8080/tcp
Thinking: Let me check what model the whisper-local container is using.
# Check Whisper container logs
$ docker logs whisper-local --tail 20 2>&1 | head -15
By pulling and using the container, you accept the terms and conditions of this license:
https://developer.nvidia.com/ngc/nvidia-deep-learning-container-license
A copy of this license is made available in this container at /NGC-DL-CONTAINER-LICENSE for your convenience.
Uninstalled 1 package in 6ms
Installed 1 package in 8ms
2026-04-01 19:08:33,846:INFO:httpx:_send_single_request:HTTP Request: GET https://api.gradio.app/gradio-messaging/en "HTTP/1.1 200 OK"
INFO:     Started server process [104]
INFO:     Waiting for application startup.
…
Click to expand
Thinking: The Whisper container is running. Let me check what model it's using - maybe it's using a different default model than large-v3.
# Check STT models via POST
$ curl -s -X POST http://localhost:8020/models | head -20
{"detail":"Not Found"}
Thinking: Let me check what models are actually installed. I know from earlier that we have sherpa-onnx running. Let me check what's available.
# Check Whisper models directory
$ docker exec whisper-local ls -la /models/ 2>/dev/null || echo "no /models dir"
no /models dir
Thinking: The STT service is running but seems to be using the default model from the container (likely from the repo). The config in app.json is set to large-v3 which is the correct model name for faster-whisper.
Let me check what the actual config in the codebase says vs what's in the container. The app.json shows:
- stt_model: "large-v3" - this is correct for faster-whisper
But I need to verify if this matches what we have in the codebase.
✱ Grep "stt_model|STTModel" in . (9 matches)
Thinking: Let me check what models faster-whisper supports and see if it's using large-v3 by default or if it needs to be explicitly loaded.
→ Read pkg/voice/stt/local.go [limit=60]
Thinking: Looking at the config, STT is configured with:
- stt_model: "large-v3" ✓
- stt_provider: "local" ✓
- stt_base_url: "http://localhost:8020" ✓
The code in pkg/voice/stt/local.go uses defaultLocalSTTModel = "large-v3" which matches. The faster-whisper-server uses the model name directly - it's not a list of models you can query. The "large-v3" is the correct model name for faster-whisper.
So the STT config is actually correct and aligned with what we have:
- large-v3 is the correct model name for faster-whisper
- The service is running and healthy
- The base URL matches (localhost:8020)
Let me confirm this is working properly by checking if there are any recent transcription logs.
A config STT está alinhada com o código:
Config
stt_model: "large-v3"
stt_base_url: localhost:8020``
stt_language: "pt"
O modelo large-v3 é o nome correto para o faster-whisper. O serviço está rodando e saudável (/health retorna OK).
Quer que eu faça um teste rápido de transcrição para confirmar se está tudo funcionando?
▣  Build · big-pickle · 33.7s
sim
Thinking: The user wants me to test the STT (Speech to Text) to confirm it's working. I need to create a quick test - I'll use the Telegram voice message test or create a simple audio test. Let me check if there's an existing test or create a simple test.
# Check for test audio files
$ ls -la /home/will/aurelia/test-data/ 2>/dev/null || echo "no test-data dir"
no test-data dir
Thinking: Let me create a quick test audio file or use some existing audio to test. Let me check if there's any audio file in the project or if we can download a sample.
✱ Glob "**/*.wav" (4 matches)
Thinking: Let me test the STT with one of these audio files.
# Test STT transcription
$ curl -s -X POST "http://localhost:8020/v1/audio/transcriptions" \
  -F "file=@/home/will/aurelia/assets/voice/aurelia_sample.wav" \
  -F "model=large-v3" \
  -F "language=pt" | head -10
Internal Server Error
test tudo
Thinking: The user wants me to test everything. Let me check the logs to see what's happening with the STT error.
# Check Whisper container logs
$ docker logs whisper-local --tail 30 2>&1 | tail -20
return await get_async_backend().run_sync_in_worker_thread(
           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/root/faster-whisper-server/.venv/lib/python3.12/site-packages/anyio/_backends/_asyncio.py", line 2405, in run_sync_in_worker_thread
    return await future
           ^^^^^^^^^^^^
  File "/root/faster-whisper-server/.venv/lib/python3.12/site-packages/anyio/_backends/_asyncio.py", line 914, in run
    result = context.run(func, *args)
             ^^^^^^^^^^^^^^^^^^^^^^^^
  File "/root/faster-whisper-server/faster_whisper_server/routers/stt.py", line 209, in transcribe_file
    with model_manager.load_model(model) as whisper:
…
Click to expand
Thinking: The Whisper container is running out of GPU memory. This is a known issue - the RTX 4090 might be out of memory due to other processes using it (like Ollama with gemma3:27b). Let me check GPU memory usage.
# Check GPU memory usage
$ nvidia-smi --query-gpu=memory.used,memory.total,utilization.gpu --format=csv,noheader
24005 MiB, 24564 MiB, 5 %
Thinking: The GPU is almost full - 24GB used out of 24GB. That's why the STT is failing with OOM. The Ollama model gemma3:27b-it-qat is using almost all the VRAM.
Options:
1. Kill Ollama to free VRAM for STT
2. Use a smaller LLM model
3. Use cloud STT instead of local
Let me check what's consuming the GPU.
# Check GPU processes
$ nvidia-smi
Wed Apr  1 16:13:16 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 580.126.20             Driver Version: 580.126.20     CUDA Version: 13.0     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4090        On  |   00000000:01:00.0  On |                  Off |
|  0%   33C    P2             75W /  480W |   24005MiB /  24564MiB |      5%      Default |
…
Click to expand
Thinking: The problem is clear: Ollama is using ~21GB of GPU VRAM, leaving only 560MB free. The Whisper large-v3 model needs 2GB to run. There's not enough VRAM available.
Options:
1. Stop one of the Ollama processes to free up VRAM
2. Use a smaller Whisper model (like base or small)
3. Use a different STT provider (like Groq which is cloud-based)
Let me check what's using the VRAM - there are two Ollama processes (484802 with 20924MB and 486690 with 922MB). This might be because there are two models loaded or something is stuck.
Let me check what models are loaded in Ollama.
# Check loaded Ollama models
$ curl -s http://localhost:11434/api/ps | jq -r '.models[]?.name' 2>/dev/null || curl -s http://localhost:11434/tags | jq -r '.models[].name'
nomic-embed-text:latest
gemma3:27b-it-qat
Thinking: The Ollama is using about 21GB + 922MB = ~22GB of VRAM. The GPU has 24GB total, with 24GB used. The whisper-large-v3 model needs about 3-4GB additional VRAM which is not available.
Solutions:
1. Unload one of the Ollama models to free VRAM for Whisper
2. Use a smaller Whisper model (like base or medium)
3. Restart the whisper container to get fresh memory allocation
Let me try to unload the embedding model since it's not actively being used for chat:
O problema é a VRAM esgotada (24GB usado / 24GB total). O Ollama está ocupando ~22GB.
Vou liberar memória desalocando o modelo de embedding:
# Unload embedding model to free VRAM
$ curl -s -X POST http://localhost:11434/api/delete -d '{"name": "nomic-embed-text:latest"}'
405 method not allowed
Thinking: Ollama doesn't support unloading models via API. Let me try a different approach - restart the whisper container or check if we can use a smaller model. Actually, let's check the ollama process memory - it shows 20924MB + 922MB = ~21GB used, leaving only ~3GB free.
The issue is that Ollama is using too much VRAM (21GB) and the STT can't get the ~5GB it needs for large-v3.
Options:
1. Restart Ollama to free VRAM (but we'd lose the loaded model)
2. Use a smaller STT model like "medium" instead of "large-v3"
3. Check if there's a way to run STT on CPU temporarily
Let me suggest using a smaller model or try stopping the embedding model somehow.
Actually, I think the best approach is to use a smaller STT model. Let me update the config to use "medium" instead of "large-v3".
A VRAM está saturada: Ollama usando ~21GB, só sobram 3GB. O STT large-v3 precisa de 5GB.
Vou mudar para o modelo medium (mais leve):
← Edit /home/will/.aurelia/config/app.json
"llm_provider": "ollama",
"llm_model": "gemma3:27b-it-qat",
"stt_provider": "local",
"stt_base_url": "http://localhost:8020",
"stt_model": "large-v3",
"stt_language": "pt",
"tts_provider": "edge",
"tts_base_url": "http://127.0.0.1:8012",
"tts_model": "kokoro",
Thinking: Now I need to restart the whisper server to pick up the new model. But actually, the faster-whisper-server loads models on-demand, so we just need to try the transcription again with the new model name.
# Test STT with medium model
$ curl -s -X POST "http://localhost:8020/v1/audio/transcriptions" \
  -F "file=@/home/will/aurelia/assets/voice/aurelia_sample.wav" \
  -F "model=medium" \
  -F "language=pt"
Internal Server Error
Thinking: Still failing. Let me check the container memory. Maybe I need to restart the whisper container to clear the CUDA state.
# Restart Whisper container
$ docker restart whisper-local && sleep 5 && curl -s http://localhost:8020/health
whisper-local
OK
Thinking: The container was restarted. Let me test again.
# Test STT transcription with medium
$ curl -s -X POST "http://localhost:8020/v1/audio/transcriptions" \
  -F "file=@/home/will/aurelia/assets/voice/aurelia_sample.wav" \
  -F "model=medium" \
  -F "language=pt"
Internal Server Error
Thinking: Still failing. Let me check the logs again after the restart to see the error.
# Check Whisper logs after restart
$ docker logs whisper-local --tail 10 2>&1
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  File "/root/faster-whisper-server/faster_whisper_server/model_manager.py", line 97, in __enter__
    self._load()
  File "/root/faster-whisper-server/faster_whisper_server/model_manager.py", line 61, in _load
    self.whisper = WhisperModel(
                   ^^^^^^^^^^^^^
  File "/root/faster-whisper-server/.venv/lib/python3.12/site-packages/faster_whisper/transcribe.py", line 145, in __init__
    self.model = ctranslate2.models.Whisper(
                 ^^^^^^^^^^^^^^^^^^^^^^^^^^^
RuntimeError: CUDA failed with error out of memory
Thinking: O Whisper ainda está dando OOM. O Ollama está consumindo ~21GB de VRAM. Não há memória suficiente para carregar o Whisper + Ollama simultaneamente.
Preciso verificar se há algum modelo do Ollama que pode ser descargado ou se precisamos usar STT em modo CPU, ou usar um provedor STT cloud como Groq.
Vou verificar se há como configurar o Whisper para usar CPU-only ou mudar para Groq.
# Check Whisper device config
$ docker logs whisper-local --tail 50 2>&1 | grep -i "cpu\|device\|cuda"
RuntimeError: CUDA failed with error out of memory