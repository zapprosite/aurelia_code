# 20260328 - Industrialização de TTS e Fix de Prosódia Brasileira

## Status: Aprovado e Implementado

## Contexto
O ecossistema Aurélia utiliza o motor Kodoro (Kokoro-82M) como sintetizador padrão via contêiner Docker. Recentemente, foram reportados problemas críticos de qualidade:
1. **Sotaque Não-Natural**: Percepção de sotaque espanhol ou europeu devido à falta de fonemização brasileira nativa.
2. **Truncamento**: O motor parava de ler frases longas abruptamente.
3. **Decaimento Abrupto**: A última palavra de cada áudio era cortada antes de terminar.

## Decisões
1. **Ativação SOTA de Fonemização BR**: Alteramos o payload da API de `language` para `lang_code: "pt-br"`. Isso força o `KPipeline` a utilizar o perfil `pt-br` do `espeak-ng` em vez de deduzir `p` (Portugal) do nome da voz (`pf_dora`).
2. **Enhanced Tail Padding**: Implementamos um sufixo prosódico de `" . . . . ."` (5 pontos espaçados) no `pkg/tts/openai_compatible.go`. Isso provê o "silêncio generativo" necessário para que o modelo finalize a síntese da última sílaba de forma natural.
3. **Controle de Janela de Tokens**: Reduzimos o limite de segmentação (`maxChars`) de 4000 para **1200 caracteres** em `pkg/tts/factory.go`. Considerando a densidade de tokens do Português, 1200 caracteres garantem margem de segurança dentro do limite rígido de 510 tokens da arquitetura Kokoro-82M.
4. **Unificação de Configuração**: O `app.json` foi atualizado para `tts_language: "pt-br"` para alinhar o frontend e o backend na mesma prosódia.

## Consequências
- **Prosódia Brasileira Autêntica**: Fim do "efeito espanhol" relatado pelo usuário.
- **Estabilidade em Produção**: Suporte a textos de qualquer tamanho via chunking seguro.
- **Fidelidade Auditiva**: Áudio completo e com encerramento natural (sem cortes).

## Referências
- Motor: `ghcr.io/remsky/kokoro-fastapi-cpu` (8880/8012)
- Modelo: Kokoro-v1_0.pth
- Voz em uso: `pf_dora` (Mapeada via `pt-br` alias)
