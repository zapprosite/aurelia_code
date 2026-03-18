---
date: 2026-03-18
author: Codex
scope: workspace super review
skills:
  - Architect-Planner
  - Security-First
status: final
---

# Super Review — 2026-03-18

## Resumo executivo
O merge estratégico ainda **não está íntegro**. A camada de governança foi trazida para o workspace, mas a aplicação prática da estratégia RIM está incompleta: há **conflitos de merge não resolvidos**, o `.context/` está **majoritariamente scaffoldado/unfilled**, e o `codebase-map.json` **não representa a realidade** do repositório atual. Em segurança, não encontrei `app.json` versionado nem segredos óbvios expostos, mas há **vazamento de paths locais herdados**, referências quebradas e **logs verbosos** que podem registrar dados sensíveis em runtime.

## 1) Inconsistências de Merge

### 1.1 Conflitos de merge ainda presentes
Achado crítico de integridade:
- `.context/agents/README.md`
- `.context/agents/architect-specialist.md`

Ambos ainda contêm marcadores `<<<<<<<`, `=======`, `>>>>>>>`, o que prova que o merge Template/Core não foi concluído. Isso viola diretamente a governança de autoridade única e compromete o consumo automatizado da camada `.context/`.

### 1.2 `.context/docs/` não está sincronizado com `/internal`
A documentação viva prometida em `.context/docs/README.md` não corresponde ao estado do código Go:
- `.context/docs/project-overview.md` → `status: unfilled`
- `.context/docs/architecture.md` → `status: unfilled`
- `.context/docs/security.md` → `status: unfilled`
- `.context/docs/tooling.md` → `status: unfilled`
- `.context/docs/data-flow.md` → `status: unfilled`
- `.context/docs/testing-strategy.md` → `status: unfilled`

Resultado: a hierarquia definida em `AGENTS.md` manda consultar `.context/docs/`, mas essa camada ainda não é uma fonte confiável de verdade.

### 1.3 `codebase-map.json` está semanticamente incorreto
O arquivo `.context/docs/codebase-map.json` aparenta ter sido gerado, porém está desalinhado do workspace real:
- informa `totalFiles: 190`, mas o workspace atual tem **381 arquivos**;
- informa apenas **10 arquivos `.md`**, mas o repositório atual tem **191**;
- `architecture.layers`, `patterns`, `entryPoints`, `keyFiles`, `publicAPI` e `symbols` estão vazios;
- `navigation.tests` aponta para `src/**/*.test.ts` e `mainLogic` para `src`, embora o core esteja em **Go** (`cmd/`, `internal/`, `pkg/`).

Isso indica que a sincronização do `ai-context` foi executada de forma incompleta ou com parser/template incorreto após o merge.

### 1.4 Evidência de merge com resíduos de outro workspace
Há resíduos claros de um repositório anterior (`/home/will/Remote-control`) dentro do workspace atual:
- `.agents/skills/security/SKILL.md` referencia arquivos `file:///home/will/Remote-control/...`
- `.context/workflow/actions.jsonl` registra operações antigas com `repoPath` apontando para `/home/will/Remote-control`

Isso quebra a portabilidade do template e mostra que a estratégia RIM não limpou completamente artefatos herdados.

### 1.5 Divergência entre mapa contextual e arquitetura real
`docs/architecture.md` descreve corretamente um workspace multi-agente e o uso de `ai-context`, mas a camada `.context/` não reflete essa arquitetura. Hoje existe uma inconsistência entre:
- **governança declarada** (`AGENTS.md`, `docs/adr/20260318-estrategia-rim.md`, `docs/architecture.md`)
- **estado material** (`.context/docs/*`, `.context/agents/*`)

## 2) Riscos de Segurança (Tier C)

### 2.1 Logs podem capturar argumentos sensíveis e conteúdo bruto do LLM
Risco principal no core Go:
- `internal/agent/loop.go` faz log de `call.Arguments`
- `internal/skill/router.go` faz log do `Raw` retornado pelo classificador quando o parse falha

Se o usuário enviar tokens, paths privados, prompts com PII ou payloads sensíveis, esses dados podem parar em log. Isso é **Tier C condicional**: torna-se crítico em ambientes com retenção persistente de logs, coleta centralizada ou compartilhamento de stdout/stderr.

### 2.2 Exposição de caminhos absolutos locais em artefatos versionados
Foram encontrados caminhos absolutos do host do mantenedor em arquivos versionados, inclusive documentação e workflow state:
- `/home/will/Remote-control/...`
- `/home/will/.aurelia/...`
- `~/Desktop/pesquisas-gemini/`

Isso não é segredo por si só, mas expõe metadados operacionais do ambiente local e reduz a segurança/portabilidade do template. Em especial, `file://` hardcoded para outro repositório pode induzir leitura incorreta de material de segurança.

### 2.3 Não há `app.json` rastreado no repositório
Nenhum `app.json` versionado foi encontrado no workspace. Isso é positivo do ponto de vista de segredos, mas impede validar se há placeholders inseguros ou chaves expostas no arquivo citado pela aplicação (`~/.aurelia/config/app.json`). A recomendação é validar esse arquivo apenas localmente, fora do Git.

## 3) Observabilidade, logs e automação

### 3.1 Observabilidade ainda é predominantemente ad-hoc
O código usa `log.Printf/log.Println/log.Fatalf` amplamente em `cmd/`, `internal/agent/`, `internal/mcp/`, `internal/telegram/` e `pkg/stt/`, sem evidência de:
- correlação por request/run id nos logs operacionais;
- masking/redaction centralizada;
- níveis estruturados consistentes;
- política clara de retenção/sanitização.

Para um sistema multi-agente com MCP, skill routing e execução de tools, isso é um déficit relevante de observabilidade.

### 3.2 Paths do runtime Linux estão bem normalizados no core
Ponto positivo: `internal/runtime/resolver.go` e `internal/runtime/bootstrap.go` usam `os.UserHomeDir`, `filepath.Join` e `AURELIA_HOME`, o que está correto para Ubuntu Desktop 24.04 LTS e portável para outros ambientes.

### 3.3 Skills integradas: automação desigual
A camada `.agents/skills/` está presente, mas parte relevante da camada contextual correspondente continua `unfilled` ou herdada de outro repositório. Na prática, há **integração estrutural**, porém ainda falta **operacionalização confiável** para consumo por agentes sem intervenção manual.

## 4) Oportunidades de Otimização de Contexto

1. **Regenerar o `.context/docs/codebase-map.json`** com parser compatível com Go e validar `cmd/`, `internal/`, `pkg/` como entrypoints/layers reais.
2. **Preencher ou remover scaffolds `status: unfilled`**; pelo contrato de governança, placeholder persistente não deveria sobreviver ao merge.
3. **Sanear `.context/workflow/actions.jsonl`** para remover histórico herdado de outro workspace ou mover esse estado para artefato não versionado.
4. **Resolver integralmente a camada `.context/agents`** antes de qualquer novo handoff multi-agente.
5. **Padronizar referências locais** (`~`, variáveis ou paths relativos) em vez de `file:///home/will/...` hardcoded.
6. **Adicionar redaction de logs** para argumentos de tools, respostas brutas do roteador e payloads oriundos do usuário.
7. **Materializar um documento real em `.context/docs/security.md`** com política de segredos, logs, permissões e MCP.

## 5) Veredito
**Status do workspace: reprovado para considerar o merge “limpo” do ponto de vista arquitetural/governança.**

O core Go parece estruturalmente consistente, mas a camada de agentes/contexto ainda contém sinais inequívocos de merge incompleto e sincronização semântica falha. Antes de expandir features, eu trataria como prioridade:
1. resolver conflitos de merge em `.context/agents`;
2. regenerar e validar `.context/docs/*`;
3. remover resíduos de `/home/will/Remote-control`;
4. reduzir logs que possam vazar dados sensíveis.

## Limitações da auditoria
- Não houve `app.json` versionado para inspeção direta.
- Não foi possível executar `go test ./...` neste ambiente porque o binário `go` não está disponível na sessão atual.
