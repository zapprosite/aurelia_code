# ADR 2026-03-26: Zero Hardcode Policy

## Status
Implementada

## Contexto

O problema real não era "qualquer segredo em qualquer lugar". O problema era drift entre:

- segredos hardcoded em código, docs ou exemplos
- exibição insegura de tokens em UI, logs e sumários
- confusão entre artefato operacional local e documentação versionada

O runtime local da Aurélia usa `app.json` como configuração efetiva da instância. Tratar esse arquivo como se fosse documentação pública gerava um contrato ruim e quebrava o onboarding.

## Decisão

Adotar **Zero Hardcode** com escopo explícito:

1. **Código, docs e templates**: segredos reais são proibidos; usar `{chave-para-env}` quando for preciso representar credenciais.
2. **Saída humana**: onboarding, dashboard, logs e qualquer superfície de observabilidade devem mascarar segredos.
3. **Configuração local da instância**: o arquivo operacional `app.json` pode persistir segredos informados pelo operador, porque ele é estado privado do runtime, não artefato de documentação.
4. **Overrides por ambiente**: variáveis de ambiente continuam soberanas quando definidas.
5. **Paridade `.env`**: `.env.example` deve continuar espelhando estruturalmente `.env`.

## Consequências

- onboarding e `SaveEditable` preservam credenciais reais da instância local
- o contrato de segurança passa a ser coerente com o runtime soberano
- a proteção obrigatória fica concentrada onde faz sentido: código, docs, exemplos, logs e UI
