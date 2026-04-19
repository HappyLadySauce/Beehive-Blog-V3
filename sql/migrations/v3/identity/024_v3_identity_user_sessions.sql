-- v3 identity: user_sessions
-- 设计见 docs/v3/identity/identity-database-design.md §5.5

CREATE TABLE identity.user_sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES identity.users (id) ON DELETE CASCADE,
  auth_source VARCHAR(32) NOT NULL,
  client_type VARCHAR(32) NULL,
  device_id VARCHAR(128) NULL,
  device_name VARCHAR(128) NULL,
  ip_address INET NULL,
  user_agent TEXT NULL,
  status VARCHAR(32) NOT NULL,
  last_seen_at TIMESTAMPTZ NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_identity_user_sessions_auth_source CHECK (auth_source IN ('local', 'sso')),
  CONSTRAINT chk_identity_user_sessions_status CHECK (status IN ('active', 'revoked', 'expired'))
);

CREATE INDEX idx_identity_sessions_user_status ON identity.user_sessions (user_id, status);
CREATE INDEX idx_identity_sessions_expires_at ON identity.user_sessions (expires_at);
CREATE INDEX idx_identity_sessions_device_id ON identity.user_sessions (device_id);
