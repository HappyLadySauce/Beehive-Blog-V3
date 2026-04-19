CREATE TABLE IF NOT EXISTS content_items (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(32) NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    summary TEXT NOT NULL DEFAULT '',
    body_markdown TEXT NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    visibility VARCHAR(32) NOT NULL DEFAULT 'private',
    ai_access VARCHAR(32) NOT NULL DEFAULT 'denied',
    author_user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    published_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_content_status CHECK (status IN ('draft', 'review', 'published', 'archived')),
    CONSTRAINT chk_content_visibility CHECK (visibility IN ('public', 'member', 'private')),
    CONSTRAINT chk_content_ai_access CHECK (ai_access IN ('allowed', 'denied'))
);

CREATE INDEX IF NOT EXISTS idx_content_items_type ON content_items(type);
CREATE INDEX IF NOT EXISTS idx_content_items_status ON content_items(status);
CREATE INDEX IF NOT EXISTS idx_content_items_visibility ON content_items(visibility);
CREATE INDEX IF NOT EXISTS idx_content_items_pub ON content_items(status, visibility, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_items_author ON content_items(author_user_id);

CREATE TABLE IF NOT EXISTS content_revisions (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    version INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    body_markdown TEXT NOT NULL DEFAULT '',
    change_note TEXT NOT NULL DEFAULT '',
    created_by BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (content_id, version)
);

CREATE INDEX IF NOT EXISTS idx_content_revisions_content ON content_revisions(content_id, version DESC);
