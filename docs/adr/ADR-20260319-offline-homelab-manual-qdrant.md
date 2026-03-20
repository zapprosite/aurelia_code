---
description: Slice nonstop para consolidar o manual operacional offline do homelab consumido pelo modelo local via Qdrant.
status: proposed
---

# ADR-20260319-offline-homelab-manual-qdrant

## Status

- Proposto

## Slice

- slug: offline-homelab-manual-qdrant
- owner: codex
- branch/worktree: `20260319-aurelia-antigravit-gemini` em `/home/will/aurelia`
- json de continuidade: `docs/adr/taskmaster/ADR-20260319-offline-homelab-manual-qdrant.json`

## Links obrigatórios

- [AGENTS.md](../../AGENTS.md)
- [REPOSITORY_CONTRACT.md](../REPOSITORY_CONTRACT.md)
- [ADR Index](./README.md)
- [plan.md](../../plan.md)
- [20260319-groq-stt-ptbr-runtime.md](./20260319-groq-stt-ptbr-runtime.md)
- [20260319-homelab-tutor-v2.md](./20260319-homelab-tutor-v2.md)
- [20260319-hierarchical-agent-swarm.md](./20260319-hierarchical-agent-swarm.md)

## Contexto

O runtime já fechou um desenho local-first:

- `qwen3.5:9b` como cérebro residente
- `SQLite` como verdade operacional local
- `Qdrant` como memória semântica derivada
- `Groq` restrito ao STT
- TTS e browser fora da decisão principal de manutenção do homelab

O gap restante é documental e operacional: o modelo local ainda não possui um manual offline único, consistente e governado com tudo o que existe no homelab. Sem isso, ele depende demais de contexto disperso em código, docs soltas e memória recente.

## Problema

Precisamos de um manual operacional canônico do homelab, legível por humanos e indexável para um `9B`, cobrindo:

- inventário de serviços, portas, containers e systemd units
- runbooks de incidentes e recuperação
- topologia de storage, GPU, modelos locais e budgets
- políticas de governança, health, backups e rollout
- mapas de bancos (`SQLite`, `PostgreSQL`, `Qdrant`, `Supabase` quando houver)
- procedimentos sem dependência de web externa

Esse manual deve ser suficientemente compacto e estruturado para ser chunkado, versionado, resumido e indexado no `Qdrant` sem virar entulho semântico.

## Decisão

- criar um `Manual Offline do Homelab` como fonte documental única para o runtime local
- tratar o manual como artefato operacional primário, versionado no repositório
- derivar embeddings para `Qdrant` a partir desse manual e de seus anexos/runbooks
- manter `SQLite` como verdade do estado vivo e `Qdrant` como mecanismo de recuperação semântica
- impedir ingestão arbitrária: só entram no índice documentos aprovados, versionados e com dono claro

## Arquitetura proposta

### Fonte de verdade

O manual será organizado em um conjunto pequeno de documentos estáveis:

- `docs/homelab/manual/00-overview.md`
- `docs/homelab/manual/10-inventory.md`
- `docs/homelab/manual/20-services.md`
- `docs/homelab/manual/30-models-and-gpu.md`
- `docs/homelab/manual/40-data-and-memory.md`
- `docs/homelab/manual/50-runbooks.md`
- `docs/homelab/manual/60-operations-and-rollout.md`

### Contrato de ingestão

Cada documento elegível para ingestão semântica deve ter:

- título
- status
- owner
- data de revisão
- tags
- nível de criticidade
- origem versionada no repositório

### Contrato de recuperação

O `9B` não decide só pelo vetor. O fluxo correto é:

1. recuperar contexto do `Qdrant`
2. cruzar com estado vivo do `SQLite`
3. aplicar política/guardrails do runtime
4. só então responder ou agir

### Regra de atualização

- toda mudança estrutural no homelab deve atualizar o manual
- toda feature que mude operação deve atualizar runbook correspondente
- `sync-ai-context` continua obrigatório no fechamento de slices não triviais
- indexação do `Qdrant` acontece só após docs consistentes e aprovados

## Escopo

- definir a estrutura do manual offline
- definir o contrato de documentos indexáveis
- definir a fronteira entre `SQLite` e `Qdrant`
- definir a rotina de ingestão, revisão e compactação
- preparar a slice futura de implementação do manual e do indexador

## Fora de escopo

- implementar agora o indexador completo
- migrar tudo automaticamente para `PostgreSQL`
- inventar memória infinita literal
- deixar o vetor mandar sem consultar o estado vivo

## Arquivos afetados

- `docs/adr/ADR-20260319-offline-homelab-manual-qdrant.md`
- `docs/adr/README.md`
- opcionalmente `docs/adr/PENDING-SLICES-20260319.md`
- futuros documentos em `docs/homelab/manual/`

## Simulações e smoke previstos

- estrutura:
  - `find docs/homelab/manual -maxdepth 1 -type f | sort`
- consistência:
  - validar frontmatter mínimo dos documentos do manual
- ingestão:
  - smoke de chunking e payload para `Qdrant`
- recuperação:
  - consulta semântica seguida de verificação em `SQLite`

## Rollout

1. abrir esta ADR
2. criar a árvore `docs/homelab/manual/`
3. consolidar inventário e runbooks existentes
4. definir indexador assíncrono para `Qdrant`
5. validar consultas do `9B` só com contexto offline

## Rollback

- manter o manual como documentação estática se a ingestão semântica ainda não estiver pronta
- não permitir que o runtime dependa exclusivamente do `Qdrant`

## Evidência esperada

- manual offline versionado e navegável
- política clara de ingestão semântica
- recuperação híbrida `Qdrant + SQLite`
- runbooks suficientes para operação sem web

## Consequências

### Positivas

- o `qwen3.5:9b` ganha contexto operacional estável sem depender da web
- o homelab fica transferível entre agentes e humanos
- o `Qdrant` recebe contexto curado, não ruído

### Negativas

- exige disciplina documental contínua
- aumenta custo inicial de curadoria
- documentos velhos e não revisados podem contaminar a recuperação

## Próximos passos

1. abrir a slice de implementação do manual offline
2. criar os documentos-base do manual
3. ligar a ingestão semântica ao contrato desta ADR
