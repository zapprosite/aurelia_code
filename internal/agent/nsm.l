# Neuro-Symbolic Memory (NSM) - Symbolic Graph Core (2026)
# Mapeia relações causais entre agentes e tarefas em RAM.

(de relacao (A Rel B)
   (put A Rel (union (get A Rel) (list B))) )

(de caminho-causal (Inicio Fim)
   (let (Visto NIL)
      (de busca (Atual Destino)
         (if (== Atual Destino)
            (list Destino)
            (unless (memq Atual Visto)
               (push 'Visto Atual)
               (mapcan
                  '((B) (let (R (busca B Destino)) (and R (cons Atual R))))
                  (get Atual 'ajudou) ) ) ) )
      (busca Inicio Fim) ) )

# Exemplo de Ingestão de Símbolo
(relacao 'regente 'ajudou 'especialista-dev)
(relacao 'especialista-dev 'ajudou 'validador-sec)

# Busca Causal
(prinl "Causal Path: " (caminho-causal 'regente 'validador-sec))
(bye)
