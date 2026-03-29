# ADR 20260328: Dashboard Estilo Perplexity/Comet

## Status
🟡 Proposto

## Contexto
Dashboard atual é básico. Evoluir para **estilo Perplexity** com:
1. **Search-first UI** - Barra de busca central
2. **Feed de respostas** com markdown rico
3. **Tabs** (Timeline, Bots, Brain, Homelab, Onboard)
4. **Real-time streaming** das respostas

## Decisões Arquiteturais

### 1. Layout Search-First
```tsx
// Layout estilo Perplexity
<Layout>
  <Sidebar tabs={tabs} />
  <main>
    <SearchBar />  // Central, grande
    <ResponseFeed items={responses} />
  </main>
  <RightPanel>  // Contextual
    <AgentStatus />
    <QuickActions />
  </RightPanel>
</Layout>
```

### 2. ResponseCard
- Avatar do bot + timestamp
- Markdown renderizado (code blocks, tables)
- Citations com links
- Copy/Share buttons
- Streaming animation

### 3. Animations
- Framer Motion para transições
- Skeleton loading
- Typing indicator
- Smooth scroll

### 4. Tech Stack
- React 18 + TypeScript
- Tailwind CSS (já no projeto?)
- Framer Motion
- Lucide React icons

## Dependências
- ✅ frontend/src/App.tsx
- ✅ internal/dashboard/ (SSE)
- ⚠️ Tailwind config

## Referências
- [Perplexity AI UI](https://www.saasui.design/application/perplexity-ai)
- [UX Trends 2026](https://uxdesign.cc/the-most-popular-experience-design-trends-of-2026-3ca85c8a3e3d)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/dashboard-perplexity
**Progress**: 0%
