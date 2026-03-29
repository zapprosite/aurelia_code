# ADR-20260327-autonomous-visual-cortex

## Status
Proposto (SOTA 2026.1)

## Contexto
Atualmente, o Jarvis só "vê" a tela quando solicitado explicitamente ou durante tarefas de browser. Um assistente verdadeiro deve ter consciência situacional constante (Visual Awareness).

## Decisão
Implementar o "Visual Cortex" (Background Visual Worker):
1.  **Scan Analítico**: Captura de tela a cada 5-15 segundos (dinâmico baseado em atividade).
2.  **Multimodal Probing**: O Qwen 3.5 VL processa a imagem em baixa resolução para detectar "Eventos de Atenção" (erros, notificações, mudanças de contexto).
3.  **Proactive Trigger**: Se um evento relevante for detectado, o Jarvis inicia uma fala proativa ("Hey, vi que o build falhou com erro de sintaxe, quer que eu corrija?").

## Consequências
- **Pró**: Transforma o assistente de Reativo para Proativo.
- **Contra**: Consumo constante de ciclos de GPU (otimizado via quantization 4-bit).
- **Privacidade**: Processamento 100% local, sem frames saindo do Homelab.
