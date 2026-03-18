# Memory Management

## O que vale a pena persistir

### Alta prioridade
- Solução para um problema que levou mais de 30 min para resolver
- Decisão técnica com trade-offs não óbvios (e os motivos da escolha)
- Padrão de erro específico do projeto com sua causa raiz
- Configuração que foi difícil de acertar

### Média prioridade
- Preferências do usuário que se repetem
- Convenções do projeto que diferem do padrão
- Atalhos e comandos específicos do ambiente

### Baixa prioridade (não persiste)
- Soluções triviais que qualquer desenvolvedor encontraria facilmente
- Contexto de uma única sessão sem aplicação futura
- Informações que mudam frequentemente

## Formato de registro de aprendizado

```markdown
## [PROBLEMA] Título descritivo
Data: YYYY-MM-DD
Contexto: [onde/quando isso aparece]
Causa raiz: [por que acontece]
Solução: [o que funciona]
Alternativas tentadas: [o que não funcionou e por quê]
```

## Quando atualizar vs. criar novo registro
- Atualizar: quando nova informação complementa ou corrige o existente
- Criar novo: quando é um problema distinto mesmo que parecido
- Arquivar: quando a solução ficou obsoleta (mudança de versão, refatoração, etc.)

## Revisão periódica
- Após grandes mudanças na stack: revisar se registros ainda são válidos
- Após resolver problema já registrado: verificar se a solução ainda é a mesma
- Skills criadas por self-healing: revisar utilidade após 30 dias de uso
