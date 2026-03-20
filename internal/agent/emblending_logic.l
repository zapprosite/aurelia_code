# Lógica de Pesagem de Emblending (PicoLisp)
# Define pesos baseados na intenção da tarefa.

(de calcular-pesos (Intencao)
   (case Intencao
      ("faturamento" (list 'Vision 0.7 'MCP 0.2 'Memory 0.1))
      ("codigo"      (list 'Vision 0.1 'MCP 0.4 'Memory 0.5))
      ("seguranca"   (list 'Vision 0.2 'MCP 0.1 'Memory 0.7))
      (T             (list 'Vision 0.3 'MCP 0.3 'Memory 0.4)) ) )

(prinl (calcular-pesos "faturamento"))
(bye)
