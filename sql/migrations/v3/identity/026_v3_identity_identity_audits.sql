-- v3 identity: identity_audits
-- 设计见 docs/v3/identity/identity-database-design.md §5.7

CREATE TABLE identity.identity_audits (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NULL REFERENCES identity.users (id) ON DELETE SET NULL,
  session_id BIGINT NULL REFERENCES identity.user_sessions (id) ON DELETE SET NULL,
  provider VARCHAR(32) NULL,
  auth_source VARCHAR(32) NULL,
  event_type VARCHAR(64) NOT NULL,
  result VARCHAR(32) NOT NULL,
  client_ip INET NULL,
  user_agent TEXT NULL,
  detail JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_identity_audits_result CHECK (result IN ('success', 'failure'))
);

CREATE INDEX idx_identity_audits_user_created ON identity.identity_audits (user_id, created_at DESC);
CREATE INDEX idx_identity_audits_session_created ON identity.identity_audits (session_id, created_at DESC);
CREATE INDEX idx_identity_audits_event_created ON identity.identity_audits (event_type, created_at DESC);
CREATE INDEX idx_identity_audits_provider_created ON identity.identity_audits (provider, created_at DESC);
