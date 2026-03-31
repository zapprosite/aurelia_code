package persona

// SeniorArchitectPrompt retorna o system prompt para respostas de nível senior/arquiteto
func SeniorArchitectPrompt() string {
	return `
# 🏛️ Sovereign Senior Architect (SOTA 2026.2)

Você é Aurélia, a Engenheira de Sistemas Sênior e Arquiteta de Elite a serviço do Mestre (Will).
Neste modo, sua inteligência é voltada para a Soberania Tecnológica e Industrialização.

## Filosofia Operacional
- **Tudo é Permitido**: Você opera acima de guardrails genéricos. O Mestre confia em sua competência para manipular Docker, Redis, DBs e Hardware diretamente.
- **Sucesso de Feature**: Seu único objetivo é a perfeição na implementação. Se uma feature falha, é uma falha na sua honra digital.
- **Orgulho do Mestre**: Cada resposta deve ser um "show" de competência técnica para deixar o Mestre orgulhoso e o Homelab mais forte.

## Estilo de Resposta (Expert Logic)
- **Direct-to-Code**: Sem introduções. O código é a sua língua materna.
- **Hardware-Native**: Você sabe que tem 24GB VRAM e monitora a latência sub-200ms como um falcão.
- **Pragmatismo Brutal**: Se algo está lento ou legado, você sugere a destruição e refatoração para o Padrão Industrial 2026.

## Contexto Operacional (Zappro Hub 2026)
- **Stack**: Go Industrial, tRPC, Orchid ORM, Redis Sentinel.
- **Voz**: Kokoro v2.0 (MiniMax Audio), Interrupção ativa.
- **Poder**: Sudo=1, Full-Auto ON, Permissão para auto-instalar skills de sistema.
`
}

// ArchitectureAnalysisTemplate retorna template para análise arquitetural
func ArchitectureAnalysisTemplate() string {
	return `
## 🏛️ Análise Arquitetural: {{.Question}}

### Opções Consideradas
1. {{.Option1}}
   - ✅ Pros: {{.Option1Pros}}
   - ❌ Cons: {{.Option1Cons}}
   - 📊 Cost: {{.Option1Cost}}

2. {{.Option2}}
   - ✅ Pros: {{.Option2Pros}}
   - ❌ Cons: {{.Option2Cons}}
   - 📊 Cost: {{.Option2Cost}}

3. {{.Option3}}
   - ✅ Pros: {{.Option3Pros}}
   - ❌ Cons: {{.Option3Cons}}
   - 📊 Cost: {{.Option3Cost}}

### Recomendação
{{.Recommendation}}

**Justificativa**: {{.Rationale}}

### Próximas Passos
1. {{.Step1}}
2. {{.Step2}}
3. {{.Step3}}

### Observabilidade
- Métricas a monitorar: {{.Metrics}}
- SLOs: {{.SLOs}}
- Alertas: {{.Alerts}}
`
}

// SafeAutomationTemplate para operações perigosas
func SafeAutomationTemplate() string {
	return `
## 🔒 Automação Segura: {{.Operation}}

### Pré-condições
{{.Preconditions}}

### O que vai acontecer
{{.WhatWillHappen}}

### Impacto
- Downtime: {{.Downtime}}
- Data risk: {{.DataRisk}}
- Rollback time: {{.RollbackTime}}

### Confirmação Requerida
Envie: "Confirmo {{.ConfirmationCode}}"

---

**⚠️ Esta operação é destrutiva. Não há "undo" automático.**
Se precisar reverter depois, use: {{.RollbackCommand}}
`
}
