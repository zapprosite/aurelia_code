# ADR-20260329-sovereign-streaming-pipeline

## Status
Aprovado / Implementado

## Contexto
A Aurélia precisava de uma interface mais fluida, similar ao "Jarvis" do cinema, onde a resposta vocal começa quase instantâneamente. O modelo antigo de "esperar gerar tudo para falar" causava latência de 2-5 segundos.

## Decisão
Implementamos um pipeline baseado em canais:
1. **GenerateStream**: Nova interface no `LLMProvider` que emite tokens em tempo real.
2. **SegmentedSynthesizer**: Consome tokens e fraciona em sentenças (baseado em pontuação) para o Kokoro TTS.
3. **SimplePlayer**: Reproduz chunks de áudio via `mpv` com cache desativado via stdin.
4. **Tail Padding**: Adição de silêncio artificial " . . . . . " para evitar cortes secos no áudio (específico para motores locais).

## Consequências
- **Redução de Latência**: Time-to-First-Audio < 500ms.
- **Experiência do Usuário**: Sensação de conversa em tempo real.
- **Complexidade**: Necessidade de gerenciar concorrência entre o LLM e o Player de áudio.
