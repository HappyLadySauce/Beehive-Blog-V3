-- v3 content: items
-- 设计见 docs/v3/contracts/service-contracts.md §5.4
-- 使用 schema content，避免与历史 public content 表混用

CREATE SCHEMA IF NOT EXISTS content;

CREATE TABLE content.items (
  id BIGSERIAL PRIMARY KEY,
  type VARCHAR(32) NOT NULL,
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'draft',
  visibility VARCHAR(32) NOT NULL DEFAULT 'private',
  ai_access VARCHAR(32) NOT NULL DEFAULT 'denied',
  summary TEXT NULL,
  cover_image_url TEXT NULL,
  owner_user_id BIGINT NOT NULL REFERENCES identity.users(id) ON DELETE RESTRICT,
  author_user_id BIGINT NOT NULL REFERENCES identity.users(id) ON DELETE RESTRICT,
  source_type VARCHAR(32) NOT NULL DEFAULT 'manual',
  current_revision_id BIGINT NULL,
  comment_enabled BOOLEAN NOT NULL DEFAULT true,
  is_featured BOOLEAN NOT NULL DEFAULT false,
  sort_order INT NOT NULL DEFAULT 0,
  published_at TIMESTAMPTZ NULL,
  archived_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_content_items_type CHECK (type IN ('article', 'note', 'project', 'experience', 'timeline_event', 'insight', 'portfolio', 'page')),
  CONSTRAINT chk_content_items_status CHECK (status IN ('draft', 'review', 'published', 'archived')),
  CONSTRAINT chk_content_items_visibility CHECK (visibility IN ('public', 'member', 'private')),
  CONSTRAINT chk_content_items_ai_access CHECK (ai_access IN ('allowed', 'denied')),
  CONSTRAINT chk_content_items_source_type CHECK (source_type IN ('manual', 'import_v1', 'import_markdown', 'agent_generated', 'agent_assisted'))
);

CREATE UNIQUE INDEX ux_content_items_slug ON content.items (slug);
CREATE INDEX idx_content_items_type_status_visibility ON content.items (type, status, visibility);
CREATE INDEX idx_content_items_published_at ON content.items (published_at DESC);
CREATE INDEX idx_content_items_owner_user ON content.items (owner_user_id);
CREATE INDEX idx_content_items_author_user ON content.items (author_user_id);
