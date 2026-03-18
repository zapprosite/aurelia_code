# Skill Creation Guide (Self-Healing)

## Processo de criação automática de skill

### 1. Capturar o aprendizado da sessão
Ao final de uma sessão onde um problema complexo foi resolvido, extraia:
- Qual era o problema central
- O que foi tentado e não funcionou
- O que funcionou e por quê
- Como identificar esse problema mais rapidamente da próxima vez

### 2. Estruturar como skill

```
nome-descritivo/
├── SKILL.md          # processo e quando usar
└── reference.md      # detalhes técnicos da solução
```

### 3. Nomear corretamente
- Nome descreve o problema, não a tecnologia: "debug-memory-leak" não "node-js-tools"
- Kebab-case sempre
- Específico o suficiente para ser encontrado, genérico o suficiente para ser reusado

### 4. Validar antes de salvar
Responda:
- Outro desenvolvedor entenderia essa skill sem contexto adicional?
- A skill é específica o suficiente para ser útil?
- A skill é genérica o suficiente para aparecer de novo?

### 5. Localização
- Problema específico do projeto: .agents/skills/ na raiz do repo
- Problema geral de desenvolvimento: ~/.config/agents/skills/

## Melhoria de skills existentes

Quando resolver um problema e já existe skill sobre ele:
1. Compare a solução atual com o que está documentado
2. Se a solução foi diferente: atualize a skill com o novo aprendizado
3. Se a solução foi a mesma: a skill está funcionando, nenhuma ação necessária
4. Se a skill estava errada: corrija e documente por que estava errada

## Critério de qualidade de uma skill auto-gerada
- Resolve o problema em menos tempo que da primeira vez
- Não requer conhecimento prévio da sessão original para entender
- Tem pelo menos um exemplo concreto
- Documenta o que não funciona além do que funciona
