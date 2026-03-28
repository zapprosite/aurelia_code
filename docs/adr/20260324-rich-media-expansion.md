# ADR 20260324: Expansão de Mídias Ricas (Rich Media)

## Status
Proposto

## Contexto
Atualmente, a Aurélia possui uma skill de voz (`aurelia-media-voice`) e um script de transcrição (`scripts/media-transcript.sh`), mas a integração com o Telegram é limitada a mensagens de voz diretas. O usuário expressou o desejo de expandir a lógica de "Voz Oficial" para mídias ricas, o que inclui vídeos do YouTube, arquivos de áudio longos e outros links de mídia.

## Decisão
Implementar um pipeline unificado de processamento de mídias ricas que:
1. **Detecta URLs de mídia** (YouTube, etc.) em mensagens de texto no Telegram.
2. **Processa arquivos de vídeo** enviados como documentos.
3. **Unifica Transcrição e Resumo**: Todo processamento de mídia rica deve resultar em uma transcrição (via Groq) seguida de um resumo executivo (via Gemma 3 local ou DeepSeek).
4. **Respeita a Identidade Visual**: Utiliza o `markdown_renderer_blocks` com alertas GitHub-style para apresentar os resultados de forma premium.

## Consequências
- **Positivas**: Melhor UX para consumo de conteúdo denso; maior utilidade da Aurélia como assistente de pesquisa.
- **Negativas**: Aumento do consumo de recursos (banda para download de vídeos, tokens para resumo).
- **Mitigação**: Uso de `yt-dlp` com extração apenas de áudio e limites de duração/tamanho conforme a política de governança de hardware.


---

## Links Obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
