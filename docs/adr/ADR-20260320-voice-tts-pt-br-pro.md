# ADR-20260320-voice-tts-pt-br-pro

## Status
Proposta

## Slice
- slug: voice-tts-pt-br-pro
- owner: humano + claude
- branch/worktree: `audio-tts-voz-pt-br-pro` em `/home/will/aurelia-agent-to-agent`

## Contexto
O Swarm da Aurélia (2026) exige uma interface de voz que não soe como um robô genérico. Precisamos de:
1.  **TTS Profissional**: Suporte ao MiniMax S2V (Speech-to-Voice) para cadência e emoção humana.
2.  **Identidade Vocal**: Clonagem autorizada de voz a partir de amostras locais.
3.  **Latência Baixa**: Streaming de áudio direto para o dashboard e Telegram.

## Decisão
1.  **Motor TTS**: Utilizar **MiniMax Audio API** via `litellm` ou integração nativa em Go.
2.  **Clonagem**: Implementar um `VoiceCloningService` que processa arquivos `.wav` locais com consentimento digital.
3.  **Runtime**: Criar `internal/voice/engine.go` para orquestrar a geração e o cache de áudio.

## Plano de Rollout (Onda 1)
1.  **Fase 1**: Integração com API MiniMax (MiniMax-S2V).
2.  **Fase 2**: Implementação de `AuthorizedCloning` (Verificação de hash de áudio).
3.  **Fase 3**: E2E - Chat -> TTS -> Telegram/Web.

## Validação
- `curl` testando o endpoint do MiniMax.
- `go test` para o `VoiceCloningService`.
- Smoke test: "Aurelia, fale 'Operação 2026 Concluída' com voz de assistente sênior".
