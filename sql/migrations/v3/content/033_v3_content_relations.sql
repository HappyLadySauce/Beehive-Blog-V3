-- v3 content: relations
-- 内容关系用于 Studio 管理内容之间的结构化出边关系

CREATE TABLE content.relations (
  id BIGSERIAL PRIMARY KEY,
  from_content_id BIGINT NOT NULL REFERENCES content.items(id) ON DELETE CASCADE,
  to_content_id BIGINT NOT NULL REFERENCES content.items(id) ON DELETE CASCADE,
  relation_type VARCHAR(32) NOT NULL,
  weight INT NOT NULL DEFAULT 0,
  sort_order INT NOT NULL DEFAULT 0,
  metadata_json JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_content_relations_not_self CHECK (from_content_id <> to_content_id),
  CONSTRAINT chk_content_relations_type CHECK (relation_type IN ('belongs_to', 'related_to', 'derived_from', 'references', 'part_of', 'depends_on', 'timeline_of'))
);

CREATE UNIQUE INDEX ux_content_relations_from_to_type ON content.relations (from_content_id, to_content_id, relation_type);
CREATE INDEX idx_content_relations_from ON content.relations (from_content_id, sort_order ASC, weight DESC, id DESC);
CREATE INDEX idx_content_relations_to ON content.relations (to_content_id);
