-- v3 identity: refresh_tokens
-- 设计见 docs/v3/identity/identity-database-design.md §5.6
-- 仅存 token 哈希，不存明文

CREATE TABLE identity.refresh_tokens (
  id BIGSERIAL PRIMARY KEY,
  session_id BIGINT NOT NULL REFERENCES identity.user_sessions (id) ON DELETE CASCADE,
  token_hash VARCHAR(255) NOT NULL,
  issued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  rotated_from_token_id BIGINT NULL REFERENCES identity.refresh_tokens (id) ON DELETE SET NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_identity_refresh_tokens_token_hash ON identity.refresh_tokens (token_hash);
CREATE INDEX idx_identity_refresh_tokens_session_id ON identity.refresh_tokens (session_id);
CREATE INDEX idx_identity_refresh_tokens_expires_at ON identity.refresh_tokens (expires_at);
CREATE INDEX idx_identity_refresh_tokens_revoked_at ON identity.refresh_tokens (revoked_at);
CREATE INDEX idx_identity_refresh_tokens_rotated_from ON identity.refresh_tokens (rotated_from_token_id);
