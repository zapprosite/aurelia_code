---
description: Aurélia navega no Gemini Web (gemini.google.com/app) e executa pesquisa profunda, flash ou raciocínio — salvando resultado local.
---

# /pesquisa-profunda

Aciona o Gemini Web como base de pesquisa da Aurélia via Playwright.
Conta `contatoalienmarketing@gmail.com` — login já deve estar salvo no browser.

---

## Modos disponíveis

| Modo | Quando usar | Modelo Gemini |
|------|-------------|---------------|
| `flash` | Resposta rápida, fatos simples | Gemini 2.0 Flash |
| `pro` | Análise profunda, documentos longos | Gemini 2.5 Pro |
| `deep` | Pesquisa longa com fontes na web | Deep Research (Pro) |
| `reasoning` | Raciocínio passo-a-passo, lógica | Gemini 2.0 Flash Thinking |

---

## Protocolo de execução (Playwright)

### 1. Abrir Gemini Web
```
mcp__playwright__browser_navigate({ url: "https://gemini.google.com/app" })
```
- Se pedir login: usar conta `contatoalienmarketing@gmail.com`
- Aguardar página carregar completamente (`browser_wait_for`)

### 2. Selecionar o modelo correto

**Flash (rápido):**
- Clicar no seletor de modelo (topo da interface)
- Selecionar "Gemini 2.0 Flash"

**Pro (análise profunda):**
- Selecionar "Gemini 2.5 Pro"

**Deep Research:**
- Selecionar "Gemini 2.5 Pro" ou "Gemini 1.5 Pro"
- Clicar no botão "Deep Research" (ícone de lupa com livro) ao lado do campo de texto
- Confirmar que o modo "Deep Research" está ativo (badge aparece)

**Reasoning (raciocínio):**
- Selecionar "Gemini 2.0 Flash Thinking" ou "Gemini 2.5 Pro" com "Thinking" ativo

### 3. Digitar a query
```
mcp__playwright__browser_click({ selector: "[data-testid='text-input'], textarea, [contenteditable='true']" })
mcp__playwright__browser_type({ text: "<query do usuário>" })
mcp__playwright__browser_press_key({ key: "Enter" })
```

### 4. Aguardar resposta
- Modo `flash`: aguardar 10–20s
- Modo `pro`: aguardar 20–45s
- Modo `deep`: aguardar 2–5 min (Deep Research gera relatório completo)
- Usar `browser_wait_for` com timeout adequado ao modo

### 5. Extrair resultado
```
mcp__playwright__browser_snapshot()
```
Capturar o texto completo da resposta. Para Deep Research, o relatório inclui fontes.

### 6. Screenshot opcional
```
mcp__playwright__browser_take_screenshot({ path: "~/Desktop/pesquisas-gemini/<slug>-<data>.png" })
```

### 7. Salvar resultado
Salvar em `~/Desktop/pesquisas-gemini/<slug>-<data>.md`:
```markdown
# Pesquisa: <título>
**Data:** <data>
**Modo:** <flash|pro|deep|reasoning>
**Query:** <query original>

---

<conteúdo da resposta>

---
**Fontes:** <listar se Deep Research>
```

---

## Exemplos de uso

**Flash — resposta rápida:**
> "pesquisa flash: qual a latência média do Cloudflare Tunnel em 2025?"

**Pro — análise:**
> "pesquisa pro: comparativo arquitetural entre Cloudflare Zero Trust e Tailscale para homelab"

**Deep Research — relatório completo:**
> "pesquisa deep: melhores práticas de observabilidade para sistemas Go em produção 2025"

**Reasoning — raciocínio:**
> "pesquisa reasoning: qual o melhor modelo local para agentes com tool use em hardware com 8GB VRAM?"

---

## Integração com a Aurélia (Telegram)

Quando o usuário envia mensagem com prefixo `pesquisa`, `research` ou `busca`, a Aurélia:
1. Detecta o modo (flash/pro/deep/reasoning) ou usa `pro` como padrão
2. Aciona este workflow via skill `deep-researcher`
3. Retorna resumo no Telegram + salva relatório completo no Desktop
4. Para Deep Research: envia o relatório em partes se > 4096 chars

---

## Skill associada

Skill: [`deep-researcher`](../skills/deep-researcher/SKILL.md)
