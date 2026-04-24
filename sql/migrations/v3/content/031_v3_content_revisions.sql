-- v3 content: revisions
-- 内容正文通过 revisions 保存历史版本，items.current_revision_id 指向当前有效版本

CREATE TABLE content.revisions (
  id BIGSERIAL PRIMARY KEY,
  content_id BIGINT NOT NULL REFERENCES content.items(id) ON DELETE CASCADE,
  revision_no INT NOT NULL,
  title_snapshot VARCHAR(255) NOT NULL,
  summary_snapshot TEXT NULL,
  body_markdown TEXT NOT NULL DEFAULT '',
  body_json JSONB NULL,
  editor_type VARCHAR(32) NOT NULL DEFAULT 'human',
  editor_user_id BIGINT NULL REFERENCES identity.users(id) ON DELETE SET NULL,
  editor_agent_client_id BIGINT NULL,
  change_summary TEXT NULL,
  source_type VARCHAR(32) NOT NULL DEFAULT 'manual',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_content_revisions_revision_no CHECK (revision_no > 0),
  CONSTRAINT chk_content_revisions_editor_type CHECK (editor_type IN ('human', 'agent', 'system')),
  CONSTRAINT chk_content_revisions_source_type CHECK (source_type IN ('manual', 'import_v1', 'import_markdown', 'agent_generated', 'agent_assisted'))
);

CREATE UNIQUE INDEX ux_content_revisions_content_revision_no ON content.revisions (content_id, revision_no);
CREATE INDEX idx_content_revisions_content_created ON content.revisions (content_id, created_at DESC);

ALTER TABLE content.items
  ADD CONSTRAINT fk_content_items_current_revision
  FOREIGN KEY (current_revision_id)
  REFERENCES content.revisions(id)
  ON DELETE SET NULL;
