-- Create file category tables and migrate file assets from scope to category_key.
-- 创建文件分类表，并将文件资产从 scope 迁移到 category_key。
-- +migrate Up

ALTER TABLE file_assets RENAME COLUMN scope TO category_key;

CREATE TABLE IF NOT EXISTS file_categories (
    category_key VARCHAR(64) PRIMARY KEY,
    display_name VARCHAR(128) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_file_categories_key_not_empty CHECK (char_length(category_key) >= 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_file_categories_default_true
    ON file_categories (is_default)
    WHERE is_default = TRUE;

CREATE TABLE IF NOT EXISTS file_category_extensions (
    category_key VARCHAR(64) NOT NULL REFERENCES file_categories(category_key) ON DELETE CASCADE,
    extension VARCHAR(16) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (category_key, extension),
    CONSTRAINT ck_file_category_extensions_not_empty CHECK (char_length(extension) >= 2)
);

INSERT INTO file_categories (category_key, display_name, description, enabled, is_default, sort_order)
VALUES ('default', '默认类型', 'System default file category.', TRUE, TRUE, 0)
ON CONFLICT (category_key) DO UPDATE
SET display_name = EXCLUDED.display_name,
    description = EXCLUDED.description,
    enabled = TRUE,
    is_default = TRUE,
    updated_at = NOW();

INSERT INTO file_category_extensions (category_key, extension)
VALUES
    ('default', '.png'),
    ('default', '.jpg'),
    ('default', '.jpeg'),
    ('default', '.webp'),
    ('default', '.avif'),
    ('default', '.pdf')
ON CONFLICT (category_key, extension) DO NOTHING;

INSERT INTO file_categories (category_key, display_name, description, enabled, is_default, sort_order)
SELECT DISTINCT
    category_key,
    category_key,
    'Backfilled from legacy file asset scope.',
    TRUE,
    FALSE,
    100
FROM file_assets
WHERE category_key IS NOT NULL
  AND category_key <> ''
  AND category_key <> 'default'
ON CONFLICT (category_key) DO NOTHING;

ALTER TABLE file_assets
    ADD CONSTRAINT fk_file_assets_category_key
    FOREIGN KEY (category_key) REFERENCES file_categories(category_key);

DROP INDEX IF EXISTS idx_file_assets_owner_scope_status;
CREATE INDEX IF NOT EXISTS idx_file_assets_owner_category_status ON file_assets (owner_user_id, category_key, status);
