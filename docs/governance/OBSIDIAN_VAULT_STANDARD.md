# Obsidian Vault Standard

> Status: ativo
> Papel: organizar o vault para que `Obsidian CLI` não vire segunda base canônica caótica

## 1. Regra principal
`Obsidian` é superfície editorial humana. Não é banco primário implícito.

Toda nota relevante precisa de frontmatter, status e lineage.

## 2. Estrutura oficial do vault

```text
10-governance/
20-apps/<app_id>/
30-repos/<repo_id>/
40-runbooks/
90-archive/
00-inbox/
```

## 3. Uso de cada pasta

### `00-inbox/`
- rascunho
- captura rápida
- material ainda não curado

Nada em `00-inbox/` é canônico por padrão.

### `10-governance/`
- ADR
- políticas
- contratos
- decisões persistentes

### `20-apps/<app_id>/`
- notas operacionais por aplicação
- conhecimento de produto
- referências de domínio do app

### `30-repos/<repo_id>/`
- notas específicas de repositório
- mapeamento técnico
- decisões locais que não são política global

### `40-runbooks/`
- procedimentos de operação
- backup
- restore
- reindex
- incidentes

### `90-archive/`
- material encerrado
- versões antigas
- decisões superadas

## 4. Frontmatter obrigatório

```yaml
app_id: aurelia
repo_id: aurelia
environment: local
canonical_bot_id: controle-db
source_system: obsidian
source_id: note/data-stack-standard
note_type: governance
status: draft
version: 1
updated_at: 2026-03-25T00:00:00Z
```

## 5. Status permitidos
- `draft`
- `curated`
- `published`
- `archived`

## 6. Regras de sync

1. nota humana nasce como `draft`
2. curadoria transforma em `curated`
3. sync controlado publica no sistema canônico
4. indexação semântica ocorre só após publicação ou regra explícita

## 7. O que é proibido
- nota sem frontmatter obrigatório
- backlog operacional solto sem owner
- ADR fora de `10-governance/`
- mistura de rascunho com runbook canônico
- sync silencioso sobrescrevendo canônico sem trilha

## 8. Donos
- `Will`: aprova estrutura do vault
- `aurelia_code`: mantém coerência arquitetural
- `controle-db`: audita organização, drift e orphan notes
