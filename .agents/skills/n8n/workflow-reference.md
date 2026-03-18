# Workflow Reference

## Padrões de estrutura

### Trigger > Process > Output (padrão básico)
Webhook / Schedule / Event
  > Set / Function / Transform
  > HTTP Request / Database / Email

### Fan-out (um trigger, múltiplos destinos)
Trigger > Split In Batches > [A, B, C em paralelo]

### Aggregator (múltiplas fontes, um destino)
[Source A] + [Source B] > Merge > Output

## Nodes essenciais

| Node | Uso |
|---|---|
| Webhook | Receber dados externos via HTTP |
| Schedule Trigger | Rodar em horários fixos (cron) |
| HTTP Request | Chamar qualquer API REST |
| Set | Definir/transformar campos |
| IF | Lógica condicional |
| Split In Batches | Processar listas em lotes |
| Merge | Combinar múltiplos fluxos |
| Code | JavaScript customizado quando nodes não bastam |
| Wait | Pausar execução por tempo ou até evento |

## Boas práticas
- Nomeie cada node com o que ele faz, não o tipo
- Use sticky notes para documentar seções complexas
- Ative Continue on Error apenas quando falha parcial é aceitável
- Salve credenciais no gerenciador nativo, nunca hardcode
- Teste sempre com execução manual antes de ativar
