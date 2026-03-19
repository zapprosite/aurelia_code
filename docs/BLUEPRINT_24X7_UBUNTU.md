# Blueprint: Aurelia 24/7 no Ubuntu Desktop sem login

Data de referência: 2026-03-19
Commit estável base: `322dcd3`
Tag estável: `20260319-main-stable-homelab-v1`

## Objetivo

Deixar a Aurelia:

- com `main` limpo
- Telegram-first
- focada em homelab
- iniciando sozinha no boot
- rodando sem sessão gráfica e sem login interativo
- estável 24/7 no Ubuntu Desktop

## Decisão principal

### Git

Não usar `git filter-repo` como ferramenta principal.

Usar `git filter-repo` só se houver:

- segredo vazado no histórico
- artefato binário indevido no histórico
- arquivo sensível que precise sair de todos os commits

Para “main sem lixo”, o caminho certo é:

1. congelar a `main` estável com tag
2. abrir branch/worktree de feature
3. fazer os ajustes 24/7 fora da `main`
4. validar
5. só então mergear ou reescrever `main` se realmente necessário

### Runtime

Não depender de `systemctl --user`.

Para modo “igual datacenter”, usar:

- serviço system-level
- `User=will`
- `WantedBy=multi-user.target`
- `After=network-online.target`

Isso remove a dependência de login no desktop.

## Estado atual

Hoje a base já tem:

- instance lock
- supressão de duplicate launch
- heartbeat
- health server
- smoke de homelab
- CI Linux
- ferramentas de homelab no core

O ponto ainda errado para produção 24/7 é o modelo de daemonização:

- hoje ainda existe foco em unit de usuário
- isso é ruim para boot sem login

## Shape final desejado

### Código

`main` deve conter só:

- baseline estável do bot
- runtime Ubuntu
- Telegram
- homelab tools
- smoke e testes

### Serviço

Trocar de:

- `~/.config/systemd/user/aurelia.service`

Para:

- `/etc/systemd/system/aurelia.service`

### Inicialização

Boot do Ubuntu:

- sobe rede
- sobe `aurelia.service`
- serviço roda como `will`
- usa `AURELIA_HOME=/home/will/.aurelia`
- não precisa abrir sessão gráfica

## Unit alvo

```ini
[Unit]
Description=Aurelia Homelab Bot
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=will
Group=will
WorkingDirectory=/opt/aurelia/current
Environment=AURELIA_HOME=/home/will/.aurelia
ExecStart=/usr/local/bin/aurelia
Restart=always
RestartSec=2
StartLimitIntervalSec=0
TimeoutStopSec=30
KillMode=process
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

## Scripts desta branch

- `scripts/aurelia.system.service`
- `scripts/install-system-daemon.sh`
- `scripts/system-daemon-status.sh`
- `scripts/system-daemon-logs.sh`

Instalação prevista:

```bash
bash ./scripts/install-system-daemon.sh
```

Comportamento da migração:

- para e desabilita `systemctl --user aurelia.service`
- move a unit antiga de usuário para backup
- instala `/etc/systemd/system/aurelia.service`
- instala `/usr/local/bin/aurelia`
- reinicia o serviço de sistema
- mantém logs em `~/.aurelia/logs/`

## Validação

```bash
go build ./cmd/aurelia
go test ./...
bash ./scripts/smoke-test-homelab.sh
systemctl status aurelia
curl 127.0.0.1:8484/health
```

## Definição de pronto

Pronto significa:

1. máquina liga e a Aurelia sobe sem ninguém logar
2. `systemctl status aurelia` mostra `active (running)`
3. `curl 127.0.0.1:8484/health` responde
4. Telegram responde em poucos segundos
5. `go test ./...` passa
6. `smoke-test-homelab.sh` passa

## Ordem correta

1. congelar `main`
2. abrir branch/worktree de feature
3. migrar para system service
4. validar boot sem login
5. mergear de volta para `main`
