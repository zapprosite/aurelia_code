# Custom Nodes Reference

## Estrutura de um custom node

```
my-node/
├── package.json
├── nodes/
│   └── MyNode/
│       ├── MyNode.node.ts
│       └── mynode.svg
└── credentials/
    └── MyNodeApi.credentials.ts
```

## Campos obrigatórios no node
- displayName: nome visível no editor
- name: identificador interno (camelCase)
- group: categoria ('transform', 'trigger', 'output')
- version: número da versão
- inputs/outputs: array com 'main'
- properties: array de campos configuráveis pelo usuário

## Instalação no n8n self-hosted
1. Coloque o pacote em ~/.n8n/custom/
2. npm install no diretório
3. Reinicie o n8n

## Tipos de propriedades disponíveis
- string, number, boolean
- options (dropdown)
- collection (grupo de campos opcionais)
- fixedCollection (grupo repetível)
- json (textarea com validação JSON)

## Boas práticas
- Sempre valide inputs antes de processar
- Use this.helpers.request() para chamadas HTTP dentro do node
- Implemente paginação quando a API retorna listas grandes
- Trate erros com mensagens claras para o usuário
