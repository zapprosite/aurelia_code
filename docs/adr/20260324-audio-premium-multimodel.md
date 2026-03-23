# ADR-20260324: Arquitetura Multimodel TTS (Premium PT-BR)

## Contexto
O bot Aurelia atualmente utiliza o XTTS v2 via `openedai-speech` para respostas de voz. Embora funcional, o XTTS pode apresentar latência alta e imprecisões tonais em Português do Brasil (PT-BR). Em Março de 2026, novas tecnologias como **Kokoro TTS** e **Fish Speech** oferecem qualidade superior e sotaques mais naturais.

## Decisão
Implementaremos um sistema de **fábrica (factory)** para síntese de voz que permita:
1.  Suportar múltiplos provedores simultâneos (XTTS, Kokoro, etc.).
2.  Definir um provedor "Premium" (Kokoro) como padrão para PT-BR devido à sua melhor prosódia e baixo consumo de VRAM.
3.  Manter o XTTS como fallback ou para vozes específicas já clonadas.

### Escolha do Modelo: Kokoro TTS
- **Por quê?**: É extremamente leve (< 1s RTF em CPU), tem acentuação perfeita em PT-BR e consome pouca VRAM (< 1.5GB em modo GPU).
- **Integração**: Utilizaremos a API compatível com OpenAI do Kokoro via container Docker independente.

## Consequências
- **Consumo**: Aumentaremos o uso de VRAM em ~1GB (Kokoro) além do XTTS.
- **Latência**: Redução significativa no tempo de resposta (TTFB) para áudio.
- **Configuração**: Novas variáveis de ambiente serão necessárias para mapear os endpoints dos provedores.

## Status
- [x] Proposto (2026-03-24)
- [/] Em Implementação
