> [!NOTE]
> Status: ✅ Arquivado / Concluído em 22/03/2026

# ADR-20260322-blocking-plan-approvals

## Status
Proposto

## Contexto
A fase 15 introduziu o `ActionPlan`, permitindo que a Aurelia comunique suas intenções. No entanto, o agente ainda pode ignorar o feedback humano e prosseguir para a execução de forma autônoma. Em um ambiente de "Industrial Homelab", ações críticas (como mudar configurações de rede ou deletar arquivos) exigem um "Hard-Gate" — um bloqueio físico no software que impede o avanço sem um token de autorização.

## Decisão
Implementaremos um `PlanStore` centralizado no backend Go. 
1. Quando a tool `propose_plan` for chamada, o plano entrará no estado `PENDING`.
2. O `agent.Loop` será modificado para que a transição de `PLANNING` -> `EXECUTION` falhe caso exista um plano crítico `PENDING`.
3. O Dashboard atuará como a "CHAVE" de desbloqueio, mudando o estado para `APPROVED`.

## Consequências
- **Positivas**: Segurança garantida contra alucinações de alto risco; supervisão humana real.
- **Negativas**: Introduz latência (espera pelo humano); requer que o Dashboard esteja acessível para destravar o agente em tarefas complexas.
