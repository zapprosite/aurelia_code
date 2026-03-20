-- Habilita pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- Tabela de Agentes (Funcionários do Escritório)
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    role_description TEXT,
    capability_symbols TEXT[], -- Símbolos para o PicoLisp
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Tabela de Tarefas e Colaborações
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    run_id TEXT,
    agent_id UUID REFERENCES agents(id),
    content TEXT,
    embedding vector(1536), -- Dimensões do OpenAI
    status TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS collaborations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_agent_id UUID REFERENCES agents(id),
    target_agent_id UUID REFERENCES agents(id),
    task_id TEXT REFERENCES tasks(id),
    interaction_type TEXT, -- 'help_requested', 'review', 'handover'
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Índices vetoriais para busca semântica
CREATE INDEX ON tasks USING ivfflat (embedding vector_cosine_ops);

-- Smart Collaboration Contracts (Immune System 2026)
CREATE TABLE agent_handshakes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    requester_id UUID REFERENCES agents(id),
    responder_id UUID REFERENCES agents(id),
    task_id UUID REFERENCES tasks(id),
    credits_awarded INT DEFAULT 0,
    reputation_gain INT DEFAULT 0,
    contract_status TEXT DEFAULT 'pending', -- 'pending', 'active', 'verified', 'disputed'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE agent_reputation (
    agent_id UUID PRIMARY KEY REFERENCES agents(id),
    total_credits INT DEFAULT 100,
    reputation_score INT DEFAULT 0,
    trust_level FLOAT DEFAULT 1.0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
