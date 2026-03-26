# ADR 20260326: Homelab Dashboard Nativo na Aurelia

## Status
Aprovado

## Contexto
O painel operacional do homelab estava embarcado na Aurelia apenas via `iframe` para
`/api/vrv/`, que por sua vez fazia proxy de um app Node separado em `127.0.0.1:3333`.
Esse app externo não tinha lógica de domínio própria; ele apenas consultava Docker,
health checks locais, ZFS, disco e arquivos de estado/changelog do host.

Essa topologia adicionava um processo extra, uma porta extra e um ponto de falha
desnecessário para uma visão que já pertence ao dashboard principal da Aurelia.

## Decisão
1. A Aurelia passa a coletar os dados do homelab nativamente em Go via `internal/homelab`.
2. O dashboard React deixa de usar `iframe` e renderiza a aba `Homelab` com componentes
   nativos consumindo `GET /api/homelab`.
3. A rota legada `/api/vrv/` deixa de depender de `:3333` e passa a redirecionar para a
   aba `Homelab` do dashboard.
4. O processo Node legado em `:3333` pode ser aposentado após a validação do novo painel.

## Consequências

**Positivas:**
- remove dependência operacional do app Node local
- elimina a necessidade de manter a porta `3333` viva
- reduz latência e acoplamento entre dashboard e monitor operacional
- centraliza a observabilidade do homelab no mesmo binário e no mesmo deploy da Aurelia

**Negativas / Riscos:**
- a Aurelia passa a executar comandos locais (`docker`, `zfs`, `df`) para montar o painel
- caminhos legados como `agent_state.json` continuam dependentes do ambiente do host

**Mitigação:**
- coleta com timeout curto e fallback para estado parcial, sem derrubar o dashboard
- leitura do `agent_state.json` por lista de caminhos candidatos, preservando compatibilidade
- rota legada mantida como redirecionamento enquanto o processo antigo é desligado

## Artefatos
- `internal/homelab/snapshot.go` — coleta e serialização do estado do homelab
- `cmd/aurelia/dashboard_homelab.go` — handler `GET /api/homelab`
- `cmd/aurelia/app.go` — remoção da dependência de proxy para `:3333`
- `frontend/src/components/dashboard/HomelabTab.tsx` — aba nativa do homelab
- `frontend/src/App.tsx` — integração da aba nova
- `frontend/src/components/sidebar/Sidebar.tsx` — renome de navegação
