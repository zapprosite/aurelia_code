# ADR-20260330-her-mode-vad-monitor

## Contexto
O atual sistema de captura de voz é "one-shot", exigindo um gatilho para iniciar. O **Modo Her** exige que o Jarvis esteja ouvindo continuamente, processando áudio apenas quando há atividade de voz (VAD) para economizar GPU/CPU e garantir privacidade.

## Decisão
Implementar um **Voice Gateway** desacoplado:
1.  **Gateway em Python**: Um script residente em memória (`scripts/voice-gateway.py`) acessa o hardware via PyAudio e realiza detecção de silêncio local.
2.  **Comunicação via Unix Socket**: O gateway transmite frames PCM e eventos de VAD via `/tmp/aurelia-voice.sock`.
3.  **Monitor em Go**: O ator `VADMonitor` consome o socket e atua como o gatilho (Trigger) para os atores de `STT` (Thinker) e interrupção do `Speaker`.

## Consequências
- **Pró**: Baixa latência, desacoplamento de drivers de hardware, suporte nativo a VAD em Python (Silero/OpenWakeword).
- **Contra**: Adiciona um processo externo (`voice-gateway.py`) que deve ser gerenciado pelo systemd ou pelo Go.
- **Soberania**: 100% local, sem dependência de serviços de nuvem para gatilho de voz.
