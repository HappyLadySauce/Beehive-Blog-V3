CREATE TABLE IF NOT EXISTS attachments (
    id BIGSERIAL PRIMARY KEY,
    storage_provider VARCHAR(32) NOT NULL DEFAULT 'local',
    bucket VARCHAR(128) NOT NULL DEFAULT '',
    object_key TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(128) NOT NULL DEFAULT '',
    ext VARCHAR(16) NOT NULL DEFAULT '',
    size_bytes BIGINT NOT NULL DEFAULT 0,
    checksum_sha256 VARCHAR(128) NOT NULL DEFAULT '',
    uploaded_by BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (storage_provider, bucket, object_key)
);

CREATE INDEX IF NOT EXISTS idx_attachments_mime ON attachments(mime_type);
CREATE INDEX IF NOT EXISTS idx_attachments_uploader ON attachments(uploaded_by);

CREATE TABLE IF NOT EXISTS content_attachments (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    attachment_id BIGINT NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
    usage_type VARCHAR(32) NOT NULL DEFAULT 'inline',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (content_id, attachment_id)
);

CREATE INDEX IF NOT EXISTS idx_content_attachments_content ON content_attachments(content_id, sort_order);
