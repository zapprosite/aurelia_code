# ADR Slice: Reactive Vision Cortex

## Contexto
O scan de tela constante consome ciclos de GPU desnecessários.

## Decisão
Implementar gatilhos reativos baseados em eventos do sistema (D-Bus/X11 events) para disparar o Scan VLM proativo.

## Consequências
- Economia de 80% de GPU idle no Córtex Visual.
- Reação contextual imediata a mudanças na UI.
