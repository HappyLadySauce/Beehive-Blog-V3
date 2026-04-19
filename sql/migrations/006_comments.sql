CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    parent_comment_id BIGINT NULL REFERENCES comments(id) ON DELETE CASCADE,
    author_user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    author_name VARCHAR(128) NOT NULL DEFAULT '',
    author_email VARCHAR(255) NOT NULL DEFAULT '',
    body_markdown TEXT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'visible',
    moderation_note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_comment_status CHECK (status IN ('visible', 'hidden', 'deleted'))
);

CREATE INDEX IF NOT EXISTS idx_comments_content ON comments(content_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status);
