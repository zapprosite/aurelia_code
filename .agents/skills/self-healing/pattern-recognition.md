# Pattern Recognition

## Como identificar padrões que valem uma skill

### Sinais de que vale criar uma skill
- Mesmo tipo de problema apareceu 2+ vezes em sessões diferentes
- A solução envolve mais de 3 passos não óbvios
- Outras pessoas no time provavelmente vão enfrentar o mesmo problema
- O problema tem variações mas o approach central é sempre o mesmo

### Sinais de que não vale criar uma skill
- Problema muito específico de um contexto único
- Solução muda frequentemente
- Já existe documentação oficial clara sobre o tema

## Categorias comuns de padrões

### Padrões de debug
- Erro específico com causa não óbvia
- Sequência de passos de diagnóstico que sempre funciona
- Combinação de ferramentas para investigar problema específico

### Padrões de implementação
- Estrutura de código que funciona bem para um caso de uso recorrente
- Integração com API/serviço externo específico
- Configuração de ferramenta com opções não óbvias

### Padrões de processo
- Workflow de deploy que reduz erros
- Checklist antes de operação arriscada
- Sequência de revisão para tipo específico de PR

## Template de análise de padrão

```
PROBLEMA: [descrição em uma frase]
FREQUÊNCIA: [quantas vezes apareceu]
VARIAÇÕES: [como o problema se manifesta de formas diferentes]
SOLUÇÃO CENTRAL: [o que resolve em todos os casos]
EXCEÇÕES: [quando a solução padrão não se aplica]
VALE SKILL?: [sim/não e por quê]
```
