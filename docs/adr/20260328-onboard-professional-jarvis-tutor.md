# ADR 20260328: Onboard Profissional Estilo Jarvis Tutor

## Status
🟡 Proposto

## Contexto
Onboarding atual é um wizard CLI básico. Usuário quer onboard **profissional** estilo Jarvis/Perplexity com:
1. **Tela de boas-vindas** animada
2. **Verificação de serviços** em tempo real (Ollama, Redis, Qdrant, GPU)
3. **Configuração guiada** com preview
4. **Tutor interativo** que ensina funcionalidades

## Decisões Arquiteturais

### 1. Tela de Boas-Vindas
```tsx
// WelcomeScreen.tsx - Estilo Perplexity
<motion.div className="welcome-hero">
  <AnimatedLogo />  // Logo com glow effect
  <h1>Jarvis Tutor</h1>
  <p>Seu assistente de IA pessoal</p>
  <ServiceStatusGrid services={services} />
</motion.div>
```

### 2. Service Health Check
Verificar em paralelo:
- Ollama (localhost:11434)
- Redis (localhost:6379)
- Qdrant (localhost:6333)
- GPU (nvidia-smi)
- Telegram API

### 3. Onboard Steps
1. **Welcome** → Detecção automática de serviços
2. **Provider** → OpenRouter / Ollama / Anthropic
3. **API Keys** → Com validação inline
4. **Telegram** → Token com teste de conexão
5. **Voice** → STT/TTS config opcional
6. **Tutor Tour** → Quick tour interativo

### 4. Tutor Jarvis
Usar skill `jarvis-tutor-24-7` como guia:
- Tutorial progressivo
- Exemplos práticos
- Comandos úteis

## Dependências
- ✅ frontend/ (React + Framer Motion)
- ✅ skills/jarvis-tutor-24-7
- ⚠️ Sistema de verificação de serviços

## Referências
- [Playwright UI](https://github.com/nicepkg/playwright-ui)
- [Perplexity AI Design Patterns](https://www.saasui.design/application/perplexity-ai)
- [S-32 Multi-Bot](docs/S32_MULTI_BOT.md)

---
**Data**: 2026-03-28
**Status**: Proposto
**Autor**: Claude (Principal Engineer)
**Slice**: feature/onboard-professional
**Progress**: 0%
