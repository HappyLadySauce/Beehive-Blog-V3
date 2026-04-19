-- v3 identity: users
-- 设计见 docs/v3/identity/identity-database-design.md §5.1
-- 使用 schema identity，避免与既有 public.users（v2）同名冲突

CREATE SCHEMA IF NOT EXISTS identity;

CREATE TABLE identity.users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  email VARCHAR(320) NULL,
  nickname VARCHAR(128) NULL,
  avatar_url TEXT NULL,
  role VARCHAR(32) NOT NULL DEFAULT 'member',
  status VARCHAR(32) NOT NULL DEFAULT 'active',
  last_login_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_identity_users_role CHECK (role IN ('member', 'admin')),
  CONSTRAINT chk_identity_users_status CHECK (status IN ('pending', 'active', 'disabled', 'locked'))
);

CREATE UNIQUE INDEX ux_identity_users_username ON identity.users (username);
CREATE UNIQUE INDEX ux_identity_users_email ON identity.users (email);
CREATE INDEX idx_identity_users_role_status ON identity.users (role, status);
