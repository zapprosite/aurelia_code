# ADR 20260319-keepassxc-cofre-humano

**Status**: Aceito  
**Data**: 2026-03-19

## Contexto

Credenciais humanas e segredos consultados manualmente estavam sendo mantidos em arquivos texto em claro, misturando documentação operacional com material sensível.

## Decisão

Adotar `KeePassXC` como cofre humano principal para credenciais manuais.

Separação formal:

- cofre humano: `KeePassXC`
- segredos de automação: arquivos mínimos e específicos de runtime
- documentação: sem segredo real

O guia operacional detalhado deixa a raiz e passa a viver em `docs/`.

## Consequências

Positivas:

- reduz exposição de segredos em texto puro
- separa melhor humano, automação e documentação
- facilita backup cifrado do material humano

Trade-offs:

- exige disciplina de migração e manutenção do cofre
- não substitui secret handling específico de serviços automatizados

## Referências

- [keepassxc_local_vault_guide_20260319.md](../keepassxc_local_vault_guide_20260319.md)
- [SECURITY.md](../../SECURITY.md)
