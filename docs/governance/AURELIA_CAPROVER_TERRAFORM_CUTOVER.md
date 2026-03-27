# Aurelia CapRover + Terraform Cutover

## Objetivo

Deixar a Aurelia pronta para sair do runtime local `systemd + :3334` e entrar em um
deploy via CapRover, sem quebrar `https://aurelia.zappro.site`.

Este guia prepara o rollout. Ele nao faz o cutover por conta propria.

## Estado atual verificado

- `captain-definition` existe na raiz e aponta para `./Dockerfile`.
- O `Dockerfile` expoe `3334` e `8484`.
- `aurelia.zappro.site` ja existe em Terraform e no `cloudflared`, mas hoje aponta para
  `http://localhost:3334`.
- Nao existe service/app da Aurelia rodando no CapRover neste host ainda.

## Regra principal

Nao altere Terraform nem `cloudflared` para `aurelia.zappro.site` antes de a app da
Aurelia existir no CapRover e responder localmente via `captain-nginx`.

Se trocar o tunnel cedo demais, voce derruba o dominio publico.

## Configuracao alvo no CapRover

Criar uma app `aurelia` no CapRover com esta configuracao:

- Deploy source: repositório com o `captain-definition` da raiz
- Build file: `./captain-definition`
- Container HTTP Port: `3334`
- Custom domain: `aurelia.zappro.site`
- Persistent directory:
  - container path: `/home/aurelia/.aurelia`
- Environment variables recomendadas:
  - `AURELIA_HOME=/home/aurelia/.aurelia`
  - `DASHBOARD_PORT=3334`
  - `HEALTH_PORT=8484`
  - `AURELIA_MODE=sovereign`

Observacao:
- `AURELIA_HOME` precisa ser persistente. Sem isso, SQLite, config e logs morrem no redeploy.
- `HEALTH_PORT=8484` continua util para diagnostico interno, mas o trafego web do CapRover
  passa pelo container HTTP port `3334`.

## Validacao antes do cutover

Depois que a app subir no CapRover, validar localmente no host:

```bash
docker service ls | rg aurelia
curl -sS -H 'Host: aurelia.zappro.site' http://127.0.0.1/ | head
curl -sS -H 'Host: aurelia.zappro.site' http://127.0.0.1/api/homelab | jq '.services'
```

Criterio de aceite antes de mexer no tunnel:

- existe service/app da Aurelia no CapRover
- `captain-nginx` responde a `Host: aurelia.zappro.site`
- a UI abre e `/api/homelab` responde pelo caminho do CapRover

## Cutover do Terraform + Tunnel

So depois da validacao acima:

1. Alterar `/srv/ops/terraform/cloudflare/variables.tf`
   - `services.aurelia.url` de `http://localhost:3334`
   - para `http://localhost:80`

2. Alterar `/home/will/.cloudflared/config.yml`
   - `aurelia.zappro.site` de `http://localhost:3334`
   - para `http://localhost:80`

3. Aplicar Terraform

```bash
cd /srv/ops/terraform/cloudflare
terraform plan
terraform apply
```

4. Recarregar `cloudflared`

```bash
sudo systemctl reload cloudflared || sudo systemctl restart cloudflared
systemctl is-active cloudflared
```

5. Validar publico

```bash
curl -sI https://aurelia.zappro.site | head -20
curl -s https://aurelia.zappro.site/api/homelab | jq '.services'
```

## Por que `localhost:80` no cutover

No modelo CapRover, quem recebe o trafego externo no host e faz o roteamento por hostname
e o `captain-nginx` em `:80/:443`.

Por isso, depois que `aurelia.zappro.site` estiver configurado como custom domain da app,
o tunnel deve apontar para `http://localhost:80`, nao mais para `:3334`.

O `Host: aurelia.zappro.site` preservado pelo request e o que permite ao CapRover entregar
a app correta.

## Pos-cutover obrigatorio

Depois do cutover real, atualizar a governanca:

- `/srv/ops/ai-governance/NETWORK_MAP.md`
- `/srv/ops/ai-governance/SUBDOMAINS.md`

Troca esperada:

- antes: `aurelia (systemd) | :3334 :8484`
- depois: Aurelia entregue por CapRover, com tunnel apontando para `:80`

## Nao fazer

- nao apontar `aurelia.zappro.site` para `localhost:80` sem custom domain pronto no CapRover
- nao remover o runtime local antes do cutover publico passar
- nao deixar `AURELIA_HOME` sem persistencia
- nao mudar Terraform sem mudar `config.yml` no mesmo change set
