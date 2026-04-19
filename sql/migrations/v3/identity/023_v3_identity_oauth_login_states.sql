-- v3 identity: oauth_login_states
-- 设计见 docs/v3/identity/identity-database-design.md §5.4

CREATE TABLE identity.oauth_login_states (
  id BIGSERIAL PRIMARY KEY,
  provider VARCHAR(32) NOT NULL,
  state VARCHAR(512) NOT NULL,
  redirect_uri TEXT NOT NULL,
  client_type VARCHAR(32) NULL,
  device_id VARCHAR(128) NULL,
  code_verifier VARCHAR(255) NULL,
  requested_scopes TEXT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  consumed_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_identity_oauth_provider_state ON identity.oauth_login_states (provider, state);
CREATE INDEX idx_identity_oauth_expires_at ON identity.oauth_login_states (expires_at);
CREATE INDEX idx_identity_oauth_consumed_at ON identity.oauth_login_states (consumed_at);
