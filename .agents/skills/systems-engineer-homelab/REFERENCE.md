# Systems Engineer — Quick Reference

## Missão

Estabilidade, observabilidade, recuperação e governança do Home Lab.

## Perguntas que esta skill deve responder bem

- O que está `down`, `degraded` ou saudável?
- Qual recurso está saturando?
- Qual serviço precisa de restart, correção de config ou rollback?
- O health está dizendo a verdade?
- O que pode ser automatizado com segurança?

## Sinais de atenção

- `systemd active` com endpoint quebrado
- `health=ok` sem checks internos úteis
- `SQLITE_BUSY`
- backlog de voz crescendo
- `arecord`/device indisponível
- `Ollama` online mas modelo/resposta degradando
- container `unhealthy` repetido

## Anti-padrões

- restart cego em cascata
- limpar logs antes de coletar evidência
- chamar sucesso só porque “subiu”
- mudar secrets/rede/deploy sem classificar risco
