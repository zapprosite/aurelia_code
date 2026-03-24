---
description: Cria uma nova branch de feature com nome criativo, verificações de segurança e setup de upstream.
---

# Fluxo de Criação de Feature (v2)

## Pré-voo
1. Verificar mudanças não comitadas: `git status --short`.
   Se houver, exibir aviso mas continuar.

## Criação da Branch
2. Gerar nome no formato `[adjetivo]-[substantivo]` com personalidade técnica sênior.
   - Alta qualidade: `quantum-dispatch`, `iron-gemma`, `silent-reactor`, `stellar-pivot`,
     `neon-sentinel`, `async-oracle`, `void-prism`, `rust-signal`, `chrome-vector`
   - Evitar: nomes genéricos (`feature-1`, `test-branch`, `fix-bug`)
3. Criar e fazer checkout:
   ```bash
   git checkout -b feature/[nome-gerado]
   ```
4. Configurar upstream imediatamente:
   ```bash
   git push -u origin feature/[nome-gerado]
   ```

## Informar
5. Exibir resumo:
   - Branch criada: `feature/[nome-gerado]`
   - Remote configurado: `origin`
   - Próximos passos:
     - Implementar a feature
     - Usar `git add -A` para incluir arquivos novos (o `.gitignore` protege secrets)
     - Executar `/ship` quando pronto para abrir PR
