CREATE TABLE IF NOT EXISTS tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    slug VARCHAR(64) NOT NULL UNIQUE,
    color VARCHAR(32) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS content_tags (
    content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (content_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_content_tags_tag ON content_tags(tag_id);

CREATE TABLE IF NOT EXISTS content_relations (
    id BIGSERIAL PRIMARY KEY,
    source_content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    target_content_id BIGINT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
    relation_type VARCHAR(64) NOT NULL,
    weight SMALLINT NOT NULL DEFAULT 1,
    note TEXT NOT NULL DEFAULT '',
    created_by BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (source_content_id, target_content_id, relation_type),
    CONSTRAINT chk_relation_self_ref CHECK (source_content_id <> target_content_id)
);

CREATE INDEX IF NOT EXISTS idx_content_relations_source ON content_relations(source_content_id);
CREATE INDEX IF NOT EXISTS idx_content_relations_target ON content_relations(target_content_id);
CREATE INDEX IF NOT EXISTS idx_content_relations_type ON content_relations(relation_type);
