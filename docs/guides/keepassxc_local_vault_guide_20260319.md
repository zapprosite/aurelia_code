# Plano de Feature: Cofre Local de Credenciais com KeePassXC

## Objetivo

Substituir o uso de arquivos de texto com segredos em claro por um cofre local criptografado, operado por humano, preservando acesso simples às credenciais sem perder autonomia operacional da máquina.

Escolha registrada para esta feature:

- Solução escolhida: `KeePassXC`

## Problema Atual

Hoje existe pelo menos um padrão inseguro de armazenamento manual de credenciais:

- arquivo `.txt` com tokens, senhas e chaves em texto puro
- material sensível salvo no Desktop
- mistura entre documentação operacional e segredo real

Esse modelo é ruim por quatro motivos:

1. um único vazamento expõe toda a operação
2. backups e syncs acidentais podem copiar segredos em claro
3. documentação vira vetor de exfiltração
4. rotação e auditoria ficam confusas

## Decisão

Adotar `KeePassXC` como cofre humano principal para credenciais manuais.

O `KeePassXC` foi escolhido porque oferece:

- banco criptografado local em arquivo `.kdbx`
- operação offline
- interface humana simples
- organização por grupos e tags
- suporte a anexos, notas e histórico
- possibilidade de usar senha mestra forte e arquivo-chave
- backup fácil do cofre criptografado

## O Que Entra e O Que Não Entra

### Entra no KeePassXC

- senhas de painéis administrativos
- tokens de API usados por humano
- chaves de serviço consultadas manualmente
- credenciais de Cloudflare, Grafana, CapRover, Supabase e afins
- notas de recuperação
- links de login

### Não entra como fonte primária

- segredos que precisam ser lidos automaticamente por serviços no boot
- credenciais de runtime injetadas por `systemd`, `.env`, arquivos de config ou secret stores já usados pela automação
- valores temporários gerados por scripts efêmeros

## Princípio Operacional

Separar claramente dois mundos:

- Cofre humano: `KeePassXC`
- Segredos de automação: arquivos mínimos, específicos e com permissão restrita

Ou seja:

- humano consulta e mantém credenciais no cofre
- serviços usam apenas o segredo estritamente necessário no local estritamente necessário
- documentação não contém segredo real

## Estrutura Recomendada

### Local do Cofre

Salvar o banco em um caminho privado, fora do Desktop:

```bash
~/.local/share/secrets/will-zappro.kdbx
```

Diretório recomendado:

```bash
~/.local/share/secrets/
```

Permissões esperadas:

- diretório: `700`
- arquivo `.kdbx`: `600`

### Estrutura Interna do Cofre

Criar grupos no `KeePassXC`:

- `Cloudflare`
- `Tunnels`
- `Supabase`
- `Grafana`
- `CapRover`
- `n8n`
- `Qdrant`
- `Infra`
- `SSH`
- `Recovery`

Para cada entrada, preencher:

- Título
- Usuário
- Senha ou token
- URL
- Notas
- Data de rotação
- Dono operacional

## Tutorial Humano: Como Implantar

### 1. Instalar o KeePassXC

No Ubuntu:

```bash
sudo apt update
sudo apt install -y keepassxc
```

### 2. Criar o diretório do cofre

```bash
mkdir -p ~/.local/share/secrets
chmod 700 ~/.local/share/secrets
```

### 3. Criar o banco

Abrir o `KeePassXC` e criar um novo banco com:

- nome: `will-zappro`
- caminho: `~/.local/share/secrets/will-zappro.kdbx`

### 4. Definir proteção do banco

Usar no mínimo:

- uma senha mestra forte, longa e única

Opcional e recomendado para endurecimento:

- um arquivo-chave adicional

Modelo recomendado:

- senha mestra com 5 a 7 palavras aleatórias ou 20+ caracteres
- arquivo-chave salvo fora do Desktop
- se possível, arquivo-chave em mídia separada

### 5. Criar os grupos

Dentro do banco, criar os grupos definidos nesta seção:

- `Cloudflare`
- `Tunnels`
- `Supabase`
- `Grafana`
- `CapRover`
- `n8n`
- `Qdrant`
- `Infra`
- `SSH`
- `Recovery`

### 6. Migrar os segredos do rascunho

Para cada segredo atualmente em `.txt`:

1. criar uma entrada no grupo correto
2. copiar o valor para o `KeePassXC`
3. preencher URL e notas de uso
4. registrar data de rotação, se existir

Exemplos de entradas:

- `Cloudflare API Token`
- `Cloudflare Tunnel Secret`
- `CapRover Admin`
- `n8n PostgreSQL`
- `Qdrant API Key`
- `Supabase Studio`
- `Supabase ANON_KEY`
- `Supabase SERVICE_ROLE_KEY`
- `Grafana Admin`

### 7. Sanitizar a documentação

Depois da migração, os documentos humanos devem ficar assim:

- podem citar nome da credencial
- podem citar onde ela é usada
- podem citar onde ela está armazenada
- não podem conter o valor real

Exemplo correto:

```text
Cloudflare API Token
Local de armazenamento: KeePassXC > Cloudflare > Cloudflare API Token
Uso: Terraform e manutenção de DNS/Tunnel
```

Exemplo incorreto:

```text
Cloudflare API Token: abc123...
```

### 8. Tratar automações separadamente

Se algum serviço precisa subir sozinho após reboot e consumir credenciais:

- manter apenas o segredo necessário no arquivo/config específico do serviço
- restringir permissões
- evitar duplicar tudo no Desktop ou em documentos de arquitetura

Exemplos aceitáveis:

- `~/.cloudflared/config.yml`
- `~/.cloudflared/<tunnel-id>.json`
- `~/.codex/secrets/...`
- `/etc/systemd/system/<servico>.service.d/*.conf`

## Política de Uso Diário

### Permitido

- consultar credenciais no `KeePassXC`
- copiar senha/token apenas no momento do uso
- manter notas de recuperação no próprio cofre
- exportar backup criptografado do `.kdbx`

### Proibido

- salvar segredo real em `~/Desktop/*.txt`
- repetir segredo em documentos de arquitetura
- colar credenciais em issue, chat, commit ou markdown de projeto
- usar um arquivo “rascunho central” em texto puro como fonte de verdade

## Modelo de Documentação Segura

Os documentos do sistema devem separar:

- arquitetura
- operação
- localização do segredo

Mas nunca o segredo em si.

Modelo recomendado para documentação:

```text
Serviço: Grafana
URL: https://monitor.zappro.site
Credencial: Grafana Admin
Armazenamento: KeePassXC > Grafana > Grafana Admin
Rotação: manual
Observação: alterar após incidente ou troca de operador
```

## Backup e Recuperação

### Backup mínimo

Fazer cópia do arquivo:

```bash
~/.local/share/secrets/will-zappro.kdbx
```

Guardar cópias em pelo menos dois lugares:

- storage local seguro
- mídia externa ou local offline

### Regras

- backup sempre criptografado
- nunca exportar CSV ou TXT em claro como rotina
- testar abertura do backup periodicamente

### Se usar arquivo-chave

Guardar o arquivo-chave separado do `.kdbx`.

Não adianta criptografar o banco e deixar o arquivo-chave no mesmo lugar público e óbvio.

## Migração do Estado Atual

### Fase 1

- criar o banco `KeePassXC`
- criar grupos
- migrar todas as entradas do rascunho manual

### Fase 2

- revisar documentos do Desktop
- remover valores reais dos documentos
- manter apenas referências para o cofre

### Fase 3

- revisar quais segredos precisam continuar em runtime para automação
- reduzir duplicação
- confirmar permissões de arquivos sensíveis

## Critérios de Aceite

Esta feature estará concluída quando:

1. existir um banco `KeePassXC` criado em caminho privado
2. todos os segredos manuais relevantes estiverem no cofre
3. documentos humanos não tiverem mais valores reais
4. a automação continuar funcionando sem depender de `.txt` no Desktop
5. o operador humano souber localizar e atualizar credenciais sem improviso

## Tutorial Rápido Para o Humano

Resumo prático:

1. abrir `KeePassXC`
2. desbloquear `will-zappro.kdbx`
3. navegar até o grupo do serviço
4. copiar a credencial necessária
5. usar
6. se rotacionar a credencial, atualizar imediatamente a entrada correspondente
7. nunca registrar o novo valor em `.md`, `.txt` ou chat

## Observação Importante

`KeePassXC` resolve o armazenamento humano de credenciais.

Ele não substitui:

- configuração segura de serviços
- controle de acesso de rede
- hardening de painéis públicos
- rotação de segredos
- segregação entre acesso humano e acesso de automação

## Resultado Esperado

Ao final desta mudança:

- o humano continua tendo acesso fácil às senhas
- as credenciais deixam de ficar expostas em texto puro no Desktop
- a documentação permanece útil sem virar vazamento
- a operação autônoma da máquina não é quebrada
