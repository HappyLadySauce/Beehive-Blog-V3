CREATE TABLE IF NOT EXISTS review_tasks (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NULL REFERENCES content_items(id) ON DELETE SET NULL,
    revision_id BIGINT NULL REFERENCES content_revisions(id) ON DELETE SET NULL,
    submitter_user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    reviewer_user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    source_type VARCHAR(32) NOT NULL DEFAULT 'human',
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    priority SMALLINT NOT NULL DEFAULT 3,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    decided_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_review_source_type CHECK (source_type IN ('human', 'agent', 'system')),
    CONSTRAINT chk_review_status CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_review_tasks_status ON review_tasks(status, priority, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_review_tasks_content ON review_tasks(content_id, status);
CREATE INDEX IF NOT EXISTS idx_review_tasks_reviewer ON review_tasks(reviewer_user_id, status);

CREATE TABLE IF NOT EXISTS review_decisions (
    id BIGSERIAL PRIMARY KEY,
    review_task_id BIGINT NOT NULL REFERENCES review_tasks(id) ON DELETE CASCADE,
    decision VARCHAR(32) NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    decided_by BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_review_decision CHECK (decision IN ('approved', 'rejected', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_review_decisions_task ON review_decisions(review_task_id, created_at DESC);
