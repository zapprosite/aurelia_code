# ADR-20260329-autonomous-visual-cortex

## Status
Proposto / Em Implementação

## Contexto
O Jarvis precisa ser capaz de "ver" o que está acontecendo no ambiente do usuário (terminal, IDE, browser) sem que o usuário precise descrever manualmente. Isso permite proatividade (sugerir correções, notar bugs, antecipar comandos).

## Decisão
Implementamos um worker de monitoramento visual:
1. **Periodic Scan**: Um worker em background captura screenshots ou extrai o DOM/Buffer do terminal a cada 10-30 segundos.
2. **Semantic Analysis (Tier 0)**: O modelo Qwen 3.5 VL analisa a imagem em busca de "Anomalias" ou "Oportunidades de Ajuda".
3. **Proactive Trigger**: Se algo crítico for detectado, o Jarvis interrompe o loop de espera e inicia uma interação sugerindo a solução.
4. **Context Awareness**: O buffer visual é mantido como parte da memória de curto prazo (Short-term Memory).

## Consequências
- **Aumento de Carga**: O scan contínuo consome GPU. Devemos otimizar a frequência com base na atividade do sistema.
- **Privacidade**: O usuário deve ter controle claro sobre o que o Jarvis "vê". Implementamos o comando `/eyes-off` para pausar o córtex visual.
