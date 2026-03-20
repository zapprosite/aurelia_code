# DSL de Leis Corporativas - Swarm Immune System (2026)
# Define o que um agente PODE ou NÃO PODE fazer.

(de lei (Acao Descricao Regra)
   (put Acao 'lei Descricao)
   (put Acao 'regra Regra) )

# Leis de Segurança Zero-Trust
(lei 'delete-database "Proibido deletar bancos de dados em produção" '(not (equal Contexto "producao")))
(lei 'external-push "Dados sensíveis não podem sair da edge" '(equal Criptografia "AES-256"))

(de verificar-acao (Acao Contexto Criptografia)
   (let (Regra (get Acao 'regra))
      (if (eval Regra)
         T
         (list 'BLOQUEADO: (get Acao 'lei)) ) ) )

# Exemplo de Verificação
(prinl "Verificação Segurança: " (verificar-acao 'delete-database "producao" "NENHUMA"))
(prinl "Verificação Segurança (Edge): " (verificar-acao 'external-push "edge" "AES-256"))
(bye)
