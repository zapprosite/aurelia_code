# Pipeline Soberano — Como usar

## Executar a fila completa
Cole no Antigravity: //pipeline

## Adicionar novo slice
Edite .agent/tasks.json, adicione objeto na queue com status="pending".

## Checar status
Cole no Antigravity: //pipeline status

## Estrutura de um slice
{
  "id": "S-XX",
  "title": "descrição curta",
  "priority": "P0|P1|P2|P3",
  "status": "pending|running|done|failed",
  "depends_on": "S-XX (opcional)",
  "smoke": { "cmd": "comando que termina com echo SMOKE_OK" },
  "instructions": "o que o agente deve fazer"
}

## Regra do smoke
Se o smoke não imprimir SMOKE_OK, o slice falha.
Smokes válidos: grep, go build, go test, curl /health, ls -la arquivo.
