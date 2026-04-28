CREATE TABLE IF NOT EXISTS file_assets (
    asset_id TEXT PRIMARY KEY,
    upload_id TEXT NOT NULL UNIQUE,
    owner_user_id BIGINT NOT NULL,
    scope VARCHAR(32) NOT NULL,
    visibility VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL,
    bucket VARCHAR(255) NOT NULL,
    object_key TEXT NOT NULL UNIQUE,
    public_url TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    content_type VARCHAR(128) NOT NULL,
    byte_size BIGINT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    uploaded_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_file_assets_scope CHECK (scope IN ('avatar', 'content_cover', 'content_image', 'attachment')),
    CONSTRAINT ck_file_assets_visibility CHECK (visibility IN ('public', 'private')),
    CONSTRAINT ck_file_assets_status CHECK (status IN ('pending', 'uploaded', 'deleted')),
    CONSTRAINT ck_file_assets_byte_size CHECK (byte_size > 0)
);

CREATE INDEX IF NOT EXISTS idx_file_assets_owner_scope_status ON file_assets (owner_user_id, scope, status);
CREATE INDEX IF NOT EXISTS idx_file_assets_upload_status ON file_assets (upload_id, status);
CREATE INDEX IF NOT EXISTS idx_file_assets_expires_at ON file_assets (expires_at) WHERE status = 'pending';
