-- v3 identity: credential_locals
-- 设计见 docs/v3/identity/identity-database-design.md §5.2

CREATE TABLE identity.credential_locals (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES identity.users (id) ON DELETE CASCADE,
  password_hash VARCHAR(255) NOT NULL,
  password_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_identity_credential_locals_user_id ON identity.credential_locals (user_id);
