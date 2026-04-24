-- v3 content: tags
-- 标签作为 content schema 内的独立资源，通过 content_tags 绑定内容

CREATE TABLE content.tags (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  slug VARCHAR(128) NOT NULL,
  description TEXT NULL,
  color VARCHAR(32) NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_content_tags_name ON content.tags (name);
CREATE UNIQUE INDEX ux_content_tags_slug ON content.tags (slug);

CREATE TABLE content.content_tags (
  id BIGSERIAL PRIMARY KEY,
  content_id BIGINT NOT NULL REFERENCES content.items(id) ON DELETE CASCADE,
  tag_id BIGINT NOT NULL REFERENCES content.tags(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_content_content_tags_content_tag ON content.content_tags (content_id, tag_id);
CREATE INDEX idx_content_content_tags_content ON content.content_tags (content_id);
CREATE INDEX idx_content_content_tags_tag ON content.content_tags (tag_id);
