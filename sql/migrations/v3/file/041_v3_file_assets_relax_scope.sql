-- Relax the scope column from a hardcoded enum to a free-form namespace string.
-- 将 scope 列从硬编码枚举放宽为自由格式的 namespace 字符串。
-- +migrate Up

ALTER TABLE file_assets DROP CONSTRAINT IF EXISTS ck_file_assets_scope;
ALTER TABLE file_assets ALTER COLUMN scope TYPE VARCHAR(64);
ALTER TABLE file_assets ADD CONSTRAINT ck_file_assets_scope_not_empty CHECK (char_length(scope) >= 1);

-- +migrate Down

ALTER TABLE file_assets DROP CONSTRAINT IF EXISTS ck_file_assets_scope_not_empty;
ALTER TABLE file_assets ALTER COLUMN scope TYPE VARCHAR(32);
ALTER TABLE file_assets ADD CONSTRAINT ck_file_assets_scope CHECK (scope IN ('avatar', 'content_cover', 'content_image', 'attachment'));
