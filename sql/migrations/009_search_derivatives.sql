CREATE TABLE IF NOT EXISTS search_documents (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    language VARCHAR(16) NOT NULL DEFAULT 'zh',
    title TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    body_plain TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    indexed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (content_id)
);

CREATE INDEX IF NOT EXISTS idx_search_documents_content ON search_documents(content_id);
CREATE INDEX IF NOT EXISTS idx_search_documents_indexed ON search_documents(indexed_at DESC);
CREATE INDEX IF NOT EXISTS idx_search_documents_type ON search_documents((metadata->>'type'));
CREATE INDEX IF NOT EXISTS idx_search_documents_status_visibility ON search_documents((metadata->>'status'), (metadata->>'visibility'));

CREATE TABLE IF NOT EXISTS content_chunks (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    chunk_no INT NOT NULL,
    chunk_text TEXT NOT NULL,
    token_count INT NOT NULL DEFAULT 0,
    embedding_model VARCHAR(64) NOT NULL DEFAULT '',
    embedding_vector BYTEA NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (content_id, chunk_no)
);

CREATE INDEX IF NOT EXISTS idx_content_chunks_content ON content_chunks(content_id, chunk_no);

CREATE TABLE IF NOT EXISTS content_summaries (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    summary_type VARCHAR(32) NOT NULL DEFAULT 'short',
    summary_text TEXT NOT NULL DEFAULT '',
    model VARCHAR(64) NOT NULL DEFAULT '',
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (content_id, summary_type)
);

CREATE INDEX IF NOT EXISTS idx_content_summaries_content ON content_summaries(content_id);

CREATE TABLE IF NOT EXISTS rag_answer_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    query_text TEXT NOT NULL,
    answer_text TEXT NOT NULL,
    source_payload JSONB NOT NULL DEFAULT '[]'::jsonb,
    latency_ms INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rag_answer_logs_user ON rag_answer_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_rag_answer_logs_created ON rag_answer_logs(created_at DESC);
