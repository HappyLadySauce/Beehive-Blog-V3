-- v3 identity: federated_identities
-- 设计见 docs/v3/identity/identity-database-design.md §5.3

CREATE TABLE identity.federated_identities (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES identity.users (id) ON DELETE CASCADE,
  provider VARCHAR(32) NOT NULL,
  provider_subject VARCHAR(255) NOT NULL,
  provider_subject_type VARCHAR(64) NOT NULL,
  unionid VARCHAR(255) NULL,
  openid VARCHAR(255) NULL,
  provider_login VARCHAR(255) NULL,
  provider_email VARCHAR(320) NULL,
  provider_display_name VARCHAR(255) NULL,
  avatar_url TEXT NULL,
  app_id_or_client_id VARCHAR(128) NULL,
  access_scope TEXT NULL,
  raw_profile JSONB NULL,
  last_login_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_identity_federated_provider_subject ON identity.federated_identities (provider, provider_subject);
CREATE INDEX idx_identity_federated_user_id ON identity.federated_identities (user_id);
CREATE INDEX idx_identity_federated_provider ON identity.federated_identities (provider);
CREATE INDEX idx_identity_federated_provider_unionid ON identity.federated_identities (provider, unionid) WHERE unionid IS NOT NULL;
