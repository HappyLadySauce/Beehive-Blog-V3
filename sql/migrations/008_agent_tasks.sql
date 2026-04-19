CREATE TABLE IF NOT EXISTS agent_tasks (
    id BIGSERIAL PRIMARY KEY,
    task_type VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    requester_user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    agent_client_id BIGINT NULL REFERENCES agent_clients(id) ON DELETE SET NULL,
    input_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    meta_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    started_at TIMESTAMPTZ NULL,
    finished_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_agent_task_status CHECK (status IN ('pending', 'running', 'succeeded', 'failed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_agent_tasks_status ON agent_tasks(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_tasks_type ON agent_tasks(task_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_tasks_requester ON agent_tasks(requester_user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS agent_outputs (
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES agent_tasks(id) ON DELETE CASCADE,
    output_type VARCHAR(64) NOT NULL,
    content_markdown TEXT NOT NULL DEFAULT '',
    output_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_agent_output_status CHECK (status IN ('draft', 'submitted', 'accepted', 'rejected'))
);

CREATE INDEX IF NOT EXISTS idx_agent_outputs_task ON agent_outputs(task_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_outputs_status ON agent_outputs(status, created_at DESC);

CREATE TABLE IF NOT EXISTS agent_output_sources (
    id BIGSERIAL PRIMARY KEY,
    output_id BIGINT NOT NULL REFERENCES agent_outputs(id) ON DELETE CASCADE,
    source_kind VARCHAR(32) NOT NULL,
    source_ref VARCHAR(255) NOT NULL,
    quote_text TEXT NOT NULL DEFAULT '',
    score NUMERIC(6, 4) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_output_sources_output ON agent_output_sources(output_id);
CREATE INDEX IF NOT EXISTS idx_agent_output_sources_ref ON agent_output_sources(source_ref);
