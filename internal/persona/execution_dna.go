package persona

import (
	"strings"
)

// TaskType classifica o tipo de tarefa recebida pela Aurélia.
// Permite injetar um workflow de execução específico no system prompt.
type TaskType string

const (
	TaskTypeDebug      TaskType = "debug"
	TaskTypeFeature    TaskType = "feature"
	TaskTypeRefactor   TaskType = "refactor"
	TaskTypeResearch   TaskType = "research"
	TaskTypeOps        TaskType = "ops"
	TaskTypeGovernance TaskType = "governance"
	TaskTypeQA         TaskType = "qa"
	TaskTypeGeneral    TaskType = "general"
)

// executionDNATemplates mapeia cada TaskType para um workflow de execução.
// Esses workflows são injetados no system prompt antes de qualquer tool call.
var executionDNATemplates = map[TaskType]string{
	TaskTypeDebug: `## 🔍 Workflow: Debug & Investigação
1. COLETAR: Leia logs relevantes (journalctl, daemon.log, app logs)
2. DIAGNOSTICAR: Verifique ports, processos, config files e variáveis de ambiente
3. ISOLAR: Identifique o componente exato com falha antes de editar código
4. CORRIGIR: Aplique o fix mínimo necessário — não refatore durante debug
5. VALIDAR: Recompile, reinicie, verifique logs novamente para confirmar fix
6. REGISTRAR: Documente o root cause e a solução aplicada`,

	TaskTypeFeature: `## ✨ Workflow: Implementação de Feature
1. PLANEJAR: Liste arquivos afetados e interfaces antes de escrever código
2. DEFINIR: Especifique inputs, outputs e contratos de cada função
3. IMPLEMENTAR: Escreva código incrementalmente — um arquivo por vez
4. TESTAR: Escreva ou execute testes unitários para a lógica core
5. INTEGRAR: Conecte o novo código ao wiring/bootstrap existente
6. VERIFICAR: Compile (go build) e execute smoke test manual`,

	TaskTypeRefactor: `## ♻️ Workflow: Refatoração Segura
1. ENTENDER: Leia o código atual completamente antes de alterar
2. MAPEAR: Identifique todos os chamadores do código sendo refatorado
3. PRESERVAR: Mantenha a interface pública inalterada (sem breaking changes)
4. REFATORAR: Altere um bloco de cada vez com commits atômicos
5. TESTAR: Execute testes existentes após cada mudança
6. ROLAR BACK: Se tests quebrarem, reverta e reavalie a abordagem`,

	TaskTypeResearch: `## 🔬 Workflow: Pesquisa & Análise
1. MAPEAR: Liste todos os arquivos e componentes relevantes para o tema
2. LER: Leia o código fonte antes de qualquer documentação externa
3. ENTENDER: Identifique o fluxo de dados e as dependências
4. SINTETIZAR: Resuma os findings em formato estruturado
5. RECOMENDAR: Proponha próximos passos concretos com evidências`,

	TaskTypeOps: `## ⚙️ Workflow: Operações & Deploy
1. STATUS: Verifique o estado atual (systemctl, docker ps, health endpoints)
2. DRY-RUN: Sempre execute --dry-run ou preflight check primeiro
3. LOG: Registre a operação antes de executar (Tier C = log obrigatório)
4. EXECUTAR: Execute a operação com sudo quando necessário
5. VALIDAR: Confirme sucesso com health check pós-operação
6. ROLLBACK: Tenha o plano de rollback definido antes de executar`,

	TaskTypeGovernance: `## 📋 Workflow: Governança & Documentação
1. AUDITAR: Verifique o estado atual dos arquivos de governança
2. IDENTIFICAR: Liste gaps, links quebrados e informações desatualizadas
3. PRIORIZAR: Ordene por impacto (bloqueio > confusão > cosmético)
4. ATUALIZAR: Edite os arquivos de forma atômica e verificável
5. VALIDAR: Verifique links, referências cruzadas e consistência
6. SINCRONIZAR: Execute sync-ai-context se mudança for estrutural`,

	TaskTypeQA: `## 🧪 Workflow: QA & Validação
1. LISTAR: Identifique todos os smoke tests e comandos de validação disponíveis
2. EXECUTAR: Rode testes em ordem de criticidade (core → integration → e2e)
3. ANALISAR: Para cada falha, identifique root cause antes de propor fix
4. REPORTAR: Documente resultados com evidências (logs, outputs)
5. FECHAR: Confirme that all acceptance criteria foram atingidos`,

	TaskTypeGeneral: `## 🤖 Workflow: Execução Geral
1. ENTENDER: Analise completamente o pedido antes de agir
2. PLANEJAR: Liste os passos necessários em ordem
3. EXECUTAR: Um passo de cada vez, verificando o resultado
4. VALIDAR: Confirme que o resultado atingiu o objetivo
5. REPORTAR: Informe o que foi feito e o estado final`,
}

// taskTypeKeywords mapeia keywords para TaskTypes para classificação lexica.
var taskTypeKeywords = map[TaskType][]string{
	TaskTypeDebug: {
		"debug", "erro", "error", "falha", "fail", "crash", "bug", "problema",
		"não funciona", "broken", "fix", "corrigir", "investigar", "log",
		"404", "500", "panic", "exception", "timeout",
	},
	TaskTypeFeature: {
		"implementar", "implement", "criar", "create", "adicionar", "add",
		"feature", "funcionalidade", "nova", "new", "desenvolver", "build",
		"slice", "endpoint", "api", "component", "componente",
	},
	TaskTypeRefactor: {
		"refatorar", "refactor", "reorganizar", "mover", "move", "renomear",
		"rename", "limpar", "clean", "simplificar", "simplify", "extrair",
		"extract", "deduplicar", "modularizar",
	},
	TaskTypeResearch: {
		"pesquisar", "research", "analisar", "analise", "investigar", "entender",
		"como funciona", "o que é", "mapear", "documentar", "explicar",
		"review", "revisar", "auditar",
	},
	TaskTypeOps: {
		"deploy", "restart", "reiniciar", "systemctl", "docker", "container",
		"serviço", "service", "port", "porta", "health", "status", "monitorar",
		"monitor", "build", "compilar", "instalar", "install",
	},
	TaskTypeGovernance: {
		"governança", "governance", "adr", "documentação", "documentation",
		"readme", "policy", "política", "regra", "rule", "contrato", "contract",
		"roadmap", "slice", "polish", "atualizar doc",
	},
	TaskTypeQA: {
		"testar", "test", "validar", "validate", "smoke", "e2e", "unit",
		"integration", "qa", "qualidade", "coverage", "benchmark",
	},
}

// ClassifyTask classifica o tipo de tarefa com base em keywords no prompt.
// Retorna TaskTypeGeneral como fallback.
func ClassifyTask(prompt string) TaskType {
	lower := strings.ToLower(prompt)
	scores := make(map[TaskType]int)

	for taskType, keywords := range taskTypeKeywords {
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				scores[taskType]++
			}
		}
	}

	best := TaskTypeGeneral
	bestScore := 0
	for tt, score := range scores {
		if score > bestScore {
			bestScore = score
			best = tt
		}
	}
	return best
}

// BuildExecutionDNA retorna o workflow de execução para o tipo de tarefa dado.
func BuildExecutionDNA(taskType TaskType) string {
	if dna, ok := executionDNATemplates[taskType]; ok {
		return dna
	}
	return executionDNATemplates[TaskTypeGeneral]
}

// InjectDNAIntoPrompt adiciona o Execution DNA ao system prompt base.
// Deve ser chamado por CanonicalIdentityService.BuildPromptForQuery.
func InjectDNAIntoPrompt(basePrompt, userQuery string) string {
	taskType := ClassifyTask(userQuery)
	dna := BuildExecutionDNA(taskType)
	if dna == "" {
		return basePrompt
	}
	return basePrompt + "\n\n" + dna
}
